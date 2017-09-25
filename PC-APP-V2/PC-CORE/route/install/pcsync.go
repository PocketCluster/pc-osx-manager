package install

import (
    "bytes"
    "encoding/binary"
    "encoding/json"
    "fmt"
    "io"
    "time"

    log "github.com/Sirupsen/logrus"
    "github.com/pkg/errors"

    "github.com/Redundancy/go-sync"
    "github.com/Redundancy/go-sync/blockrepository"
    "github.com/Redundancy/go-sync/blocksources"
    "github.com/Redundancy/go-sync/chunks"
    "github.com/Redundancy/go-sync/filechecksum"
    "github.com/Redundancy/go-sync/index"
    "github.com/Redundancy/go-sync/patcher"
    "github.com/Redundancy/go-sync/patcher/multisources"
    "github.com/Redundancy/go-sync/showpipe"
    "github.com/stkim1/pc-core/route"
    "github.com/stkim1/pc-core/utils/dblerr"
)

// reads the file headers and checks the magic string, then the semantic versioning
// return : in order of 'filesize', 'blocksize', 'blockcount', 'rootHash', 'error'
func readHeadersAndCheck(headerReader io.Reader) (int64, uint32, uint32, []byte, error) {
    var (
        bMagic                []byte = make([]byte, len(gosync.PocketSyncMagicString))
        major, minor, patch   uint16 = 0, 0, 0
        filesize              int64  = 0
        blocksize, blockcount uint32 = 0, 0
        hLen                  uint32 = 0
        rootHash              []byte = nil
    )
    // magic string
    if _, err := headerReader.Read(bMagic); err != nil {
        return 0, 0, 0, nil, errors.WithStack(err)
    } else if string(bMagic) != gosync.PocketSyncMagicString {
        return 0, 0, 0, nil, errors.New("meta header does not confirm. Not a valid meta")
    }

    // version
    for _, v := range []*uint16{&major, &minor, &patch} {
        if err := binary.Read(headerReader, binary.LittleEndian, v); err != nil {
            return 0, 0, 0, nil, errors.WithStack(err)
        }
    }
    if major != gosync.PocketSyncMajorVersion || minor != gosync.PocketSyncMinorVersion || patch != gosync.PocketSyncPatchVersion {
        return 0, 0, 0, nil, errors.Errorf("The acquired version (%v.%v.%v) does not match the tool (%v.%v.%v).",
            major, minor, patch,
            gosync.PocketSyncMajorVersion, gosync.PocketSyncMinorVersion, gosync.PocketSyncPatchVersion)
    }

    if err := binary.Read(headerReader, binary.LittleEndian, &filesize); err != nil {
        return 0, 0, 0, nil, errors.WithStack(err)
    }
    if err := binary.Read(headerReader, binary.LittleEndian, &blocksize); err != nil {
        return 0, 0, 0, nil, errors.WithStack(err)
    }
    if err := binary.Read(headerReader, binary.LittleEndian, &blockcount); err != nil {
        return 0, 0, 0, nil, errors.WithStack(err)
    }
    if err := binary.Read(headerReader, binary.LittleEndian, &hLen); err != nil {
        return 0, 0, 0, nil, errors.WithStack(err)
    }
    rootHash = make([]byte, hLen)
    if _, err := headerReader.Read(rootHash); err != nil {
        return 0, 0, 0, nil, errors.WithStack(err)
    }
    return filesize, blocksize, blockcount, rootHash, nil
}

func readIndex(indexReader io.Reader, blocksize, blockcount uint, rootChksum []byte) (*index.ChecksumIndex, error) {
    var (
        generator    = filechecksum.NewFileChecksumGenerator(blocksize)
        idx          *index.ChecksumIndex = nil
    )

    readChunks, err := chunks.CountedLoadChecksumsFromReader(
        indexReader,
        blockcount,
        generator.GetWeakRollingHash().Size(),
        generator.GetStrongHash().Size(),
    )
    if err != nil {
        return nil, errors.WithStack(err)
    }

    idx = index.MakeChecksumIndex(readChunks)
    cRootChksum, err := idx.SequentialChecksumList().RootHash()
    if err != nil {
        return nil, errors.WithStack(err)
    }
    if bytes.Compare(cRootChksum, rootChksum) != 0 {
        return nil, errors.Errorf("[ERR] mismatching checksum integrity")
    }

    return idx, nil
}

type syncActionPack struct {
    reader     io.ReadCloser
    writer     io.WriteCloser
    report     chan showpipe.PipeProgress
    msync      *multisources.MultiSourcePatcher
}

func (p *syncActionPack) close() {
    p.msync.Close()
    p.writer.Close()
    p.reader.Close()
}

func prepSync(repoList []string, syncData []byte, refChksum, imageURL string) (*syncActionPack, error) {
    var headIndexReader = bytes.NewBuffer(syncData)
    filesize, blocksize, blockcount, rootHash, err := readHeadersAndCheck(headIndexReader)
    if err != nil {
        return nil, errors.WithStack(err)
    }
    err = isTwoChksumSame(rootHash, refChksum)
    if err != nil {
        return nil, errors.WithStack(err)
    }
    chksumIdx, err := readIndex(headIndexReader, uint(blocksize), uint(blockcount), rootHash)
    if err != nil {
        return nil, errors.WithStack(err)
    }
    var (
        reader, writer, report = showpipe.PipeWithReport(uint64(filesize))
        resolver                = blockrepository.MakeKnownFileSizedBlockResolver(int64(blocksize), filesize)
        verifier                = &filechecksum.HashVerifier{
            Hash:                filechecksum.DefaultStrongHashGenerator(),
            BlockSize:           uint(blocksize),
            BlockChecksumGetter: chksumIdx,
        }
        repoSrcList []patcher.BlockRepository = nil
    )
    for rID, r := range repoList {
        repoSrcList = append(repoSrcList,
            blockrepository.NewBlockRepositoryBase(
                uint(rID),
                blocksources.NewRequesterWithTimeout(fmt.Sprintf("https://%s%s", r, imageURL), "PocketCluster/0.1.4 (OSX)", true, timeout),
                resolver,
                verifier))
    }
    msync, err := multisources.NewMultiSourcePatcher(writer, repoSrcList, chksumIdx)
    if err != nil {
        return nil, errors.WithStack(err)
    }
    return &syncActionPack{
        reader:    reader,
        writer:    writer,
        report:    report,
        msync:     msync,
    }, nil
}

func execSync(feeder route.ResponseFeeder, action *syncActionPack, stopC chan struct{}, rpPath, uaPath string) error {
    var (
        uerrC      chan error = make(chan error)
        perrC      chan error = make(chan error)
        uerr, perr      error = nil, nil
        unarchDone       bool = false
        patchDone        bool = false
    )
    defer func() {
        close(uerrC)
        close(perrC)
    }()
    go func(act *syncActionPack, errC chan error, rDir string) {
        errC <- xzUncompressor(act.reader, rDir)
        log.Debugf("execSync()->unarch() Closed")
    }(action, uerrC, uaPath)
    go func(act *syncActionPack, errC chan error) {
        errC <- act.msync.Patch()
        log.Debugf("execSync()->patcher() Closed")
    }(action, perrC)

    // wait a bit to patch action to start so we don't accidentally make requests on close BlockRepository when user
    // interruption comes in before actual patch activity. (This needs to be fixed)
    <- time.After(time.Millisecond * 100)

    for {
        select {
            // close everythign
            case <-stopC: {
                go action.close()
            }

            // patch error
            case err := <- perrC: {
                patchDone = true
                perr = err
                if unarchDone {
                    return errors.WithStack(dblerr.SummarizeErrors(perr, uerr))
                }
                // regardless of error or not, patcher should close action as it's the one to quit in normal cond.
                go action.close()
            }

            // this is emergency as unarchiving fails
            case err := <- uerrC: {
                unarchDone = true
                uerr = err
                if patchDone {
                    return errors.WithStack(dblerr.SummarizeErrors(uerr, perr))
                }
                if err != nil {
                    go action.close()
                }
            }

            // report progress
            case rpt := <- action.report: {
                if !(unarchDone || patchDone) {
                    data, err := json.Marshal(route.ReponseMessage{
                        "package-progress": {
                            "total-size":   rpt.TotalSize,
                            "received":     rpt.Received,
                            "remaining":    rpt.Remaining,
                            "speed":        rpt.Speed,
                            "done-percent": rpt.DonePercent,
                        },
                    })
                    if err != nil {
                        log.Errorf(err.Error())
                        continue
                    }
                    err = feeder.FeedResponseForPost(rpPath, string(data))
                    if err != nil {
                        log.Errorf(err.Error())
                        continue
                    }
                }
            }
        }
    }
}
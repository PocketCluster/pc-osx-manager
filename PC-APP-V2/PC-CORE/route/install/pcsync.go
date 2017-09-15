package install

import (
    "bytes"
    "encoding/binary"
    "fmt"
    "io"

    "github.com/Redundancy/go-sync"
    "github.com/Redundancy/go-sync/blockrepository"
    "github.com/Redundancy/go-sync/blocksources"
    "github.com/Redundancy/go-sync/chunks"
    "github.com/Redundancy/go-sync/filechecksum"
    "github.com/Redundancy/go-sync/index"
    "github.com/Redundancy/go-sync/patcher"
    "github.com/Redundancy/go-sync/patcher/multisources"
    "github.com/Redundancy/go-sync/showpipe"
    "github.com/pkg/errors"
)

// reads the file headers and checks the magic string, then the semantic versioning
// return : in order of 'filesize', 'blocksize', 'blockcount', 'rootHash', 'error'
func readHeadersAndCheck(r io.Reader) (int64, uint32, uint32, []byte, error) {
    var (
        bMagic                []byte = make([]byte, len(gosync.PocketSyncMagicString))
        major, minor, patch   uint16 = 0, 0, 0
        filesize              int64  = 0
        blocksize, blockcount uint32 = 0, 0
        hLen                  uint32 = 0
        rootHash              []byte = nil
    )
    // magic string
    if _, err := r.Read(bMagic); err != nil {
        return 0, 0, 0, nil, errors.WithStack(err)
    } else if string(bMagic) != gosync.PocketSyncMagicString {
        return 0, 0, 0, nil, errors.New("meta header does not confirm. Not a valid meta")
    }

    // version
    for _, v := range []*uint16{&major, &minor, &patch} {
        if err := binary.Read(r, binary.LittleEndian, v); err != nil {
            return 0, 0, 0, nil, errors.WithStack(err)
        }
    }
    if major != gosync.PocketSyncMajorVersion || minor != gosync.PocketSyncMinorVersion || patch != gosync.PocketSyncPatchVersion {
        return 0, 0, 0, nil, errors.Errorf("The acquired version (%v.%v.%v) does not match the tool (%v.%v.%v).",
            major, minor, patch,
            gosync.PocketSyncMajorVersion, gosync.PocketSyncMinorVersion, gosync.PocketSyncPatchVersion)
    }

    if err := binary.Read(r, binary.LittleEndian, &filesize); err != nil {
        return 0, 0, 0, nil, errors.WithStack(err)
    }
    if err := binary.Read(r, binary.LittleEndian, &blocksize); err != nil {
        return 0, 0, 0, nil, errors.WithStack(err)
    }
    if err := binary.Read(r, binary.LittleEndian, &blockcount); err != nil {
        return 0, 0, 0, nil, errors.WithStack(err)
    }
    if err := binary.Read(r, binary.LittleEndian, &hLen); err != nil {
        return 0, 0, 0, nil, errors.WithStack(err)
    }
    rootHash = make([]byte, hLen)
    if _, err := r.Read(rootHash); err != nil {
        return 0, 0, 0, nil, errors.WithStack(err)
    }
    return filesize, blocksize, blockcount, rootHash, nil
}

func readIndex(rd io.Reader, blocksize, blockcount uint, rootHash []byte) (*index.ChecksumIndex, error) {
    var (
        generator    = filechecksum.NewFileChecksumGenerator(blocksize)
        idx          *index.ChecksumIndex = nil
    )

    readChunks, err := chunks.CountedLoadChecksumsFromReader(
        rd,
        blockcount,
        generator.GetWeakRollingHash().Size(),
        generator.GetStrongHash().Size(),
    )
    if err != nil {
        return nil, errors.WithStack(err)
    }

    idx = index.MakeChecksumIndex(readChunks)
    cRootHash, err := idx.SequentialChecksumList().RootHash()
    if err != nil {
        return nil, errors.WithStack(err)
    }
    if bytes.Compare(cRootHash, rootHash) != 0 {
        return nil, errors.Errorf("[ERR] mismatching integrity checksum")
    }

    return idx, nil
}

type patchActionPack struct {
    reader     io.ReadCloser
    writer     io.WriteCloser
    report     chan showpipe.PipeProgress
    msync      *multisources.MultiSourcePatcher
}

func prepSync(repoList []string, syncData []byte, refChksum, imageURL string) (*patchActionPack, error) {
    filesize, blocksize, blockcount, rootHash, err := readHeadersAndCheck(bytes.NewBuffer(syncData))
    if err != nil {
        return nil, errors.WithStack(err)
    }
    err = isTwoChksumSame(rootHash, refChksum)
    if err != nil {
        return nil, errors.WithStack(err)
    }
    chksumIdx, err := readIndex(bytes.NewBuffer(syncData), uint(blocksize), uint(blockcount), rootHash)
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
    return &patchActionPack {
        reader:     reader,
        writer:     writer,
        report:     report,
        msync:      msync,
    }, nil
}

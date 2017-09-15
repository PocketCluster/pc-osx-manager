package install

import (
    "bytes"
    "io"
    "io/ioutil"
    "testing"
    "os"
    "path/filepath"
    "runtime"
    "time"

    log "github.com/Sirupsen/logrus"
    "github.com/pkg/errors"
    "github.com/Redundancy/go-sync/chunks"
    "github.com/Redundancy/go-sync/blockrepository"
    "github.com/Redundancy/go-sync/blocksources"
    "github.com/Redundancy/go-sync/filechecksum"
    "github.com/Redundancy/go-sync/merkle"
    "github.com/Redundancy/go-sync/patcher"
    "github.com/Redundancy/go-sync/patcher/multisources"
    "github.com/Redundancy/go-sync/rollsum"
    "github.com/Redundancy/go-sync/showpipe"
)

const (
    testRoutPath string = "test.routh.path"
    BLOCKSIZE           = 8
)

var (
    REFERENCE_BUFFER []byte                        = nil
    REFERENCE_BLOCKS [][]byte                      = nil
    REFERENCE_HASHES [][]byte                      = nil
    REFERENCE_RTHASH []byte                        = nil
    REFERENCE_CHKSEQ chunks.SequentialChecksumList = nil
    REFERENCE_BUFFSZ int                           = 0
    BLOCK_COUNT      int                           = 0
)

func setup() {
    _, testfile, _, _ := runtime.Caller(0)
    archfile := filepath.Join(filepath.Dir(testfile), "test.txt.tar.xz")
    data, err := ioutil.ReadFile(archfile)
    if err != nil {
        log.Panic(err.Error())
    }

    REFERENCE_BUFFER = data
    REFERENCE_BLOCKS = [][]byte{}
    REFERENCE_HASHES = [][]byte{}
    BLOCK_COUNT      = 0
    REFERENCE_BUFFSZ = len(data)
    maxLen          := len(data)
    m               := filechecksum.DefaultStrongHashGenerator()

    for i := 0; i < maxLen; i += BLOCKSIZE {
        last := i + BLOCKSIZE

        if last >= maxLen {
            last = maxLen
        }

        block := data[i:last]

        REFERENCE_BLOCKS = append(REFERENCE_BLOCKS, block)
        m.Reset()
        m.Write([]byte(block))
        REFERENCE_HASHES = append(REFERENCE_HASHES, m.Sum(nil))
    }

    BLOCK_COUNT = len(REFERENCE_BLOCKS)
    REFERENCE_CHKSEQ = buildSequentialChecksum(REFERENCE_BLOCKS, REFERENCE_HASHES, BLOCKSIZE)
    rootchksum, err := REFERENCE_CHKSEQ.RootHash()
    if err != nil {
        log.Panic(err.Error())
    }
    REFERENCE_RTHASH = rootchksum
    log.Debugf("Root Merkle Hash %v", REFERENCE_RTHASH)
}

func clean() {
    REFERENCE_BUFFER = nil
    REFERENCE_BLOCKS = nil
    REFERENCE_HASHES = nil
    REFERENCE_RTHASH = nil
    REFERENCE_CHKSEQ = nil
    BLOCK_COUNT      = 0
}

// --- test util function ---
func byteToReadSeeker() io.ReadSeeker {
    return bytes.NewReader(REFERENCE_BUFFER)
}

func buildSequentialChecksum(refBlks [][]byte, sChksums [][]byte, blocksize int) chunks.SequentialChecksumList {
    var (
        chksum = chunks.SequentialChecksumList{}
        rsum   = rollsum.NewRollsum64(uint(blocksize))
    )

    for i := 0; i < len(refBlks); i++ {
        var (
            wsum = make([]byte, blocksize)
            blk     = []byte(refBlks[i])
        )
        rsum.Reset()
        rsum.SetBlock(blk)
        rsum.GetSum(wsum)

        chksum = append(
            chksum,
            chunks.ChunkChecksum{
                ChunkOffset:    uint(i),
                WeakChecksum:   wsum,
                StrongChecksum: sChksums[i],
            })
    }
    return chksum
}

// --- test block reference ---
type testBlkRef struct{}
func (t *testBlkRef) EndBlockID() uint {
    return REFERENCE_CHKSEQ[len(REFERENCE_CHKSEQ) - 1].ChunkOffset
}

func (t *testBlkRef) MissingBlockSpanForID(blockID uint) (patcher.MissingBlockSpan, error) {
    for _, c := range REFERENCE_CHKSEQ {
        if c.ChunkOffset == blockID {
            return patcher.MissingBlockSpan{
                BlockSize:     c.Size,
                StartBlock:    c.ChunkOffset,
                EndBlock:      c.ChunkOffset,
            }, nil
        }
    }
    return patcher.MissingBlockSpan{}, errors.Errorf("[ERR] invalid missing block index %v", blockID)
}

func (t *testBlkRef) VerifyRootHash(hashes [][]byte) error {
    hToCheck, err := merkle.SimpleHashFromHashes(hashes)
    if err != nil {
        return err
    }
    if bytes.Compare(hToCheck, REFERENCE_RTHASH) != 0 {
        return errors.Errorf("[ERR] calculated root hash different from referenece")
    }
    return nil
}

// --- test feeder ---
type testFeed struct {
    path    string
    payload string
}

func (t *testFeed) FeedResponseForGet(path, payload string) error {
    t.path    = path
    t.payload = payload
    return nil
}

func (t *testFeed) FeedResponseForPost(path, payload string) error {
    t.path    = path
    t.payload = payload
    return nil
}

func (t *testFeed) FeedResponseForPut(path, payload string) error {
    t.path    = path
    t.payload = payload
    return nil
}

func (t *testFeed) FeedResponseForDelete(path, payload string) error {
    t.path    = path
    t.payload = payload
    return nil
}

// --- test action pack ---
func testActionPack(repos []patcher.BlockRepository) (*syncActionPack, error) {
    var (
        reader, writer, report = showpipe.PipeWithReport(uint64(REFERENCE_BUFFSZ))
    )
    msync, err := multisources.NewMultiSourcePatcher(
        writer,
        repos,
        &testBlkRef{},
    )
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

func Test_ExecSyncSuccess_Normal(t *testing.T) {
    log.SetLevel(log.DebugLevel)
    setup()
    var (
        stopC = make(chan struct{})
        tmpdir = os.TempDir()
        tFeeder = &testFeed{}
        repos = []patcher.BlockRepository{
            blockrepository.NewReadSeekerBlockRepository(
                0,
                byteToReadSeeker(),
                blockrepository.MakeKnownFileSizedBlockResolver(BLOCKSIZE, int64(REFERENCE_BUFFSZ)),
                blockrepository.FunctionChecksumVerifier(func(startBlockID uint, data []byte) ([]byte, error){
                    m := filechecksum.DefaultStrongHashGenerator()
                    m.Write(data)
                    return m.Sum(nil), nil
                }),
            ),
        }
    )
    defer func() {
        close(stopC)
        clean()
    }()
    act, err := testActionPack(repos)
    if err != nil {
        t.Fatal(err.Error())
    }
    err = execSync(tFeeder, act, stopC, testRoutPath, tmpdir)
    if err != nil {
        t.Fatal(err.Error())
    }
}

func Test_ExecSyncFail_With_SingleRepo(t *testing.T) {
    log.SetLevel(log.DebugLevel)
    setup()
    var (
        stopC = make(chan struct{})
        tmpdir = os.TempDir()
        tFeeder = &testFeed{}
        repos = []patcher.BlockRepository{
            blockrepository.NewBlockRepositoryBase(
                0,
                blocksources.FunctionRequester(func(start, end int64) (data []byte, err error) {
                    if start < 40 {
                        return REFERENCE_BUFFER[start:end], nil
                    }
                    return nil, &blocksources.TestError{}
                }),
                blockrepository.MakeKnownFileSizedBlockResolver(BLOCKSIZE, int64(REFERENCE_BUFFSZ)),
                blockrepository.FunctionChecksumVerifier(func(startBlockID uint, data []byte) ([]byte, error){
                    m := filechecksum.DefaultStrongHashGenerator()
                    m.Write(data)
                    return m.Sum(nil), nil
                }),
            ),
        }
    )
    defer func() {
        close(stopC)
        clean()
    }()
    act, err := testActionPack(repos)
    if err != nil {
        t.Fatal(err.Error())
    }
    err = execSync(tFeeder, act, stopC, testRoutPath, tmpdir)
    if err == nil {
        t.Fatal(err.Error())
    }
    // err should not be nil
    t.Log(err.Error())
}

func Test_ExecSyncFail_With_UserStop(t *testing.T) {
    log.SetLevel(log.DebugLevel)
    setup()
    var (
        stopC = make(chan struct{})
        progC = make(chan int64)
        ctrlC = make(chan bool)
        errC  = make(chan error)

        tmpdir = os.TempDir()
        tFeeder = &testFeed{}
        repos = []patcher.BlockRepository{
            blockrepository.NewBlockRepositoryBase(
                0,
                blocksources.FunctionRequester(func(start, end int64) (data []byte, err error) {
                    progC <- start
                    <- ctrlC
                    log.Debugf("let's hand over data for (%v)[%v:%v]", REFERENCE_BUFFSZ, start, end)
                    return REFERENCE_BUFFER[start:end], nil
                }),
                blockrepository.MakeKnownFileSizedBlockResolver(BLOCKSIZE, int64(REFERENCE_BUFFSZ)),
                blockrepository.FunctionChecksumVerifier(func(startBlockID uint, data []byte) ([]byte, error){
                    m := filechecksum.DefaultStrongHashGenerator()
                    m.Write(data)
                    return m.Sum(nil), nil
                }),
            ),
        }
    )
    defer func() {
        clean()
        close(progC)
        close(ctrlC)
        close(errC)
    }()
    act, err := testActionPack(repos)
    if err != nil {
        t.Fatal(err.Error())
    }
    go func() {
        for p := range progC {
            if p < 40 {
                ctrlC <- true
            } else {
                log.Debugf("since patcher has passed a threshold, let's stop")
                close(stopC)
            }
        }
    }()
    go func() {
       errC <- execSync(tFeeder, act, stopC, testRoutPath, tmpdir)
    }()
    select {
        case <- time.After(time.Second * 10): {
            t.Fatal("test timeout failure")
        }
        case err := <- errC: {
            if err == nil {
                t.Fatal(err.Error())
            }
            // err should not be nil
            t.Log(err.Error())
            if !multisources.IsInterruptError(err) {
                t.Fatalf("error should be user halt : %v", err.Error())
            }
        }
    }
}
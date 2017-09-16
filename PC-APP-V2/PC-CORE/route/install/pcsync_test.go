package install

import (
    "bytes"
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


// - * - * - * - * - * - * - * - * - * - * - * - TEST SET - * - * - * - * - * - * - * - * - * - * - * - * - * - * - * -
func Test_ExecSyncSuccess_Normal(t *testing.T) {
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

func Test_ExecSyncFail_With_BlockRepo(t *testing.T) {
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
        syncC = make(chan struct{})

        tmpdir = os.TempDir()
        tFeeder = &testFeed{}
        repos = []patcher.BlockRepository{
            blockrepository.NewBlockRepositoryBase(
                0,
                blocksources.FunctionRequester(func(start, end int64) (data []byte, err error) {
                    <- syncC
                    time.Sleep(time.Millisecond * 10)
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
    }()
    act, err := testActionPack(repos)
    if err != nil {
        t.Fatal(err.Error())
    }
    go func() {
        <- syncC
        time.Sleep(time.Millisecond * 200)
        log.Debugf("patcher has passed a timer mark. Let's stop")
        log.Debugf("\n\t --- PATCHER HAS PASSED A TIMER MARK. LET'S GENERATE USER STOP ---\n")
        close(stopC)
    }()

    // sync repo & stopper
    close(syncC)
    err = execSync(tFeeder, act, stopC, testRoutPath, tmpdir)
    if err == nil {
        t.Fatal(err.Error())
    }
    log.Errorf(err.Error())

    // we want to see if patcher/unarchiver loop closed succesfully
    time.Sleep(time.Second)
}

func Test_ExecSyncFail_With_Unarchive(t *testing.T) {
    log.SetLevel(log.DebugLevel)
    setup()
    var (
        stopC = make(chan struct{})
        syncC = make(chan struct{})

        tmpdir = os.TempDir()
        tFeeder = &testFeed{}
        repos = []patcher.BlockRepository{
            blockrepository.NewBlockRepositoryBase(
                0,
                blocksources.FunctionRequester(func(start, end int64) (data []byte, err error) {
                    <- syncC
                    time.Sleep(time.Millisecond * 10)
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
    }()
    act, err := testActionPack(repos)
    if err != nil {
        t.Fatal(err.Error())
    }
    go func() {
        <- syncC
        time.Sleep(time.Millisecond * 200)
        log.Debugf("\n\t --- PATCHER HAS PASSED A TIMER MARK. LET'S GENERATE IO ERROR FOR UNARCHIVER ---\n")
        act.reader.Close()
    }()

    // sync repo & stopper
    close(syncC)
    err = execSync(tFeeder, act, stopC, testRoutPath, tmpdir)
    if err == nil {
        t.Fatal(err.Error())
    }
    log.Errorf(err.Error())

    // we want to see if patcher/unarchiver loop closed succesfully
    time.Sleep(time.Second)
}
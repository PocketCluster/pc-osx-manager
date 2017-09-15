package install

import (
    "bytes"
    "io"
    "io/ioutil"
    "testing"
    "os"
    "path/filepath"
    "runtime"

    log "github.com/Sirupsen/logrus"
    "github.com/pkg/errors"
    "golang.org/x/crypto/ripemd160"
    "github.com/Redundancy/go-sync/chunks"
    "github.com/Redundancy/go-sync/blockrepository"
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
    REFERENCE_BUFFER []byte        = nil
    REFERENCE_BLOCKS [][]byte      = nil
    REFERENCE_HASHES [][]byte      = nil
    REFERENCE_RTHASH []byte        = nil
    REFERENCE_CHKSEQ chunks.SequentialChecksumList = nil
    REFERENCE_BUFSIZ int           = 0
    BLOCK_COUNT      int           = 0
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
    REFERENCE_BUFSIZ = len(data)
    maxLen          := len(data)
    m               := ripemd160.New()

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
func testActionPack() (*syncActionPack, error) {
    var (
        reader, writer, report = showpipe.PipeWithReport(uint64(REFERENCE_BUFSIZ))
        repos = []patcher.BlockRepository{
            blockrepository.NewReadSeekerBlockRepository(
                0,
                byteToReadSeeker(),
                blockrepository.MakeNullUniformSizeResolver(BLOCKSIZE),
                blockrepository.FunctionChecksumVerifier(func(startBlockID uint, data []byte) ([]byte, error){
                    m := ripemd160.New()
                    m.Write(data)
                    return m.Sum(nil), nil
                }),
            ),
        }
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

func Test_ExecSync_Normal(t *testing.T) {
    log.SetLevel(log.DebugLevel)
    setup()
    var (
        stopC = make(chan struct{})
        tmpdir = os.TempDir()
        tFeeder = &testFeed{}
    )
    defer func() {
        close(stopC)
        clean()
    }()
    act, err := testActionPack()
    if err != nil {
        t.Fatal(err.Error())
    }
    err = execSync(tFeeder, act, stopC, testRoutPath, tmpdir)
    if err != nil {
        t.Fatal(err.Error())
    }
}
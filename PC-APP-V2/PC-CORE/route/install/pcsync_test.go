package install

import (
    "bytes"
    "io"
    "testing"
    "os"

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
    REFERENCE_STRING    = "The quick brown fox jumped over the lazy dog | The quick brown fox jumped over the lazy dog"
)

var (
    REFERENCE_BUFFER *bytes.Buffer = nil
    REFERENCE_BLOCKS []string      = nil
    REFERENCE_HASHES [][]byte      = nil
    REFERENCE_RTHASH []byte        = nil
    REFERENCE_CHKSEQ chunks.SequentialChecksumList = nil
    BLOCK_COUNT      int           = 0
)

func setup() {
    REFERENCE_BUFFER = bytes.NewBufferString(REFERENCE_STRING)
    REFERENCE_BLOCKS = []string{}
    REFERENCE_HASHES = [][]byte{}
    BLOCK_COUNT      = 0

    maxLen          := len(REFERENCE_STRING)
    m               := ripemd160.New()

    log.SetLevel(log.DebugLevel)
    for i := 0; i < maxLen; i += BLOCKSIZE {
        last := i + BLOCKSIZE

        if last >= maxLen {
            last = maxLen
        }

        block := REFERENCE_STRING[i:last]

        REFERENCE_BLOCKS = append(REFERENCE_BLOCKS, block)
        m.Write([]byte(block))
        REFERENCE_HASHES = append(REFERENCE_HASHES, m.Sum(nil))
        m.Reset()
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

func stringToReadSeeker(input string) io.ReadSeeker {
    return bytes.NewReader([]byte(input))
}

func buildSequentialChecksum(refBlks []string, sChksums [][]byte, blocksize int) chunks.SequentialChecksumList {
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
        reader, writer, report = showpipe.PipeWithReport(uint64(len(REFERENCE_STRING)))
        repos = []patcher.BlockRepository{
            blockrepository.NewReadSeekerBlockRepository(
                0,
                stringToReadSeeker(REFERENCE_STRING),
                blockrepository.MakeNullUniformSizeResolver(BLOCKSIZE),
                nil),
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
    var (
        stopC = make(chan struct{})
        tmpdir = os.TempDir()
        tFeeder = &testFeed{}
    )
    act, err := testActionPack()
    if err != nil {
        t.Fatal(err.Error())
    }
    err = execSync(tFeeder, act, stopC, testRoutPath, tmpdir)
    if err != nil {
        t.Fatal(err.Error())
    }
}
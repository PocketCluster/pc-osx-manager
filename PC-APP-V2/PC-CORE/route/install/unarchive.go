package install

import (
    "archive/tar"
    "io"
    "os"
    "path/filepath"

    log "github.com/Sirupsen/logrus"
    "github.com/pkg/errors"
    "xi2.org/x/xz"
    "github.com/Redundancy/go-sync"
)

func xzUncompressor(archiveReader io.Reader, blocksize uint32, uncompPath string) error {
    var (
        xreader   *xz.Reader
        unarchive *tar.Reader
        err       error
    )

    // Check that the server actually sent compressed data
    xreader, err = xz.NewReader(archiveReader, 0)
    if err != nil {
        return errors.WithStack(err)
    }

    unarchive = tar.NewReader(xreader)
    for {
        header, err := unarchive.Next()
        if err == io.EOF {
            break
        } else if err != nil {
            return errors.WithStack(err)
        }

        path := filepath.Join(uncompPath, header.Name)
        info := header.FileInfo()
        if info.IsDir() {
            if err = os.MkdirAll(path, info.Mode()); err != nil {
                return errors.WithStack(err)
            }
            continue
        }
        file, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, info.Mode())
        if err != nil {
            return errors.WithStack(err)
        }
        written, err := io.Copy(file, unarchive)
        if err != nil {
            file.Close()
            return errors.WithStack(err)
        } else {
            log.Debugf("written %v", written)
        }
        err = file.Close()
        if err != nil {
            return errors.WithStack(err)
        }
    }

    // when reader reaches this point, it means we have some unexhausted buffer to read off.
    blksz := blocksize
    if blksz < uint32(gosync.PocketSyncDefaultBlockSize) {
        blksz = gosync.PocketSyncDefaultBlockSize
    }
    var (
        buf = make([]byte, gosync.PocketSyncDefaultBlockSize)
    )
    for {
        _, err = archiveReader.Read(buf)
        if err != nil {
            break
        }
    }

    return nil
}
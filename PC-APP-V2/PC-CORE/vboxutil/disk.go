package vboxutil

import (
    "archive/tar"
    "bytes"
    "fmt"
    "io"
    "os"
    "os/exec"
    "path/filepath"

    "github.com/pkg/errors"
)

const (
    magicString string = "boot2docker, please format-me"
)

// MachineDisk information.
type MachineDisk struct {
    // Debugging
    Debug        bool
    // VM disk image size (MB)
    DiskSize     uint
    // VMDK target path
    DiskImage    string


    // VBoxManager executable
    VBM          string
    // public key
    PubKey     []byte
}

func NewMachineDisk(baseFolder, imageName string, diskSize uint, pubKey []byte, debug bool) *MachineDisk {
    return &MachineDisk {
        Debug:        debug,
        DiskSize:     diskSize,
        DiskImage:    filepath.Join(baseFolder, fmt.Sprintf("%s.vmdk", imageName)),

        VBM:          "/usr/local/bin/VBoxManage",
        PubKey:       pubKey,
    }
}

func (m *MachineDisk) BuildDiskImage() error {
    if _, err := os.Stat(m.DiskImage); err != nil {
        if !os.IsNotExist(err) {
            return errors.WithStack(err)
        }

        buf := new(bytes.Buffer)
        tw := tar.NewWriter(buf)

        // magicString first so the automount script knows to format the disk
        file := &tar.Header{Name: magicString, Size: int64(len(magicString))}
        if err := tw.WriteHeader(file); err != nil {
            return errors.WithStack(err)
        }
        if _, err := tw.Write([]byte(magicString)); err != nil {
            return errors.WithStack(err)
        }
        // .ssh/key.pub => authorized_keys
        file = &tar.Header{Name: ".ssh", Typeflag: tar.TypeDir, Mode: 0700}
        if err := tw.WriteHeader(file); err != nil {
            return errors.WithStack(err)
        }
        file = &tar.Header{Name: ".ssh/authorized_keys", Size: int64(len(m.PubKey)), Mode: 0644}
        if err := tw.WriteHeader(file); err != nil {
            return errors.WithStack(err)
        }
        if _, err := tw.Write(m.PubKey); err != nil {
            return errors.WithStack(err)
        }
        file = &tar.Header{Name: ".ssh/authorized_keys2", Size: int64(len(m.PubKey)), Mode: 0644}
        if err := tw.WriteHeader(file); err != nil {
            return errors.WithStack(err)
        }
        if _, err := tw.Write(m.PubKey); err != nil {
            return errors.WithStack(err)
        }
        if err := tw.Close(); err != nil {
            return errors.WithStack(err)
        }

        // Create the dest dir.
        if err := os.MkdirAll(filepath.Dir(m.DiskImage), 0755); err != nil {
            return errors.WithStack(err)
        }
        // Fill in the magic string so boot2docker VM will detect this and format
        // the disk upon first boot.
        if err := makeDiskImage(m, bytes.NewReader(buf.Bytes())); err != nil {
            return errors.WithStack(err)
        }

        if m.Debug {
            fmt.Println("Initializing disk with ssh keys")
            fmt.Printf("WRITING: %s\n-----\n", buf)
        }
    }

    return nil
}

// MakeDiskImage makes a disk image at dest with the given size in MB. If r is
// not nil, it will be read as a raw disk image to convert from.
func makeDiskImage(m *MachineDisk, r io.Reader) error {
    // Convert a raw image from stdin to the MachineDisk.DiskImage VMDK image.
    sizeBytes := int64(m.DiskSize) << 20 // usually won't fit in 32-bit int (max 2GB)
    cmd := exec.Command(m.VBM, "convertfromraw", "stdin", m.DiskImage,
        fmt.Sprintf("%d", sizeBytes), "--format", "VMDK")

    if m.Debug {
        cmd.Stdout = os.Stdout
        cmd.Stderr = os.Stderr
    }

    stdin, err := cmd.StdinPipe()
    if err != nil {
        return errors.WithStack(err)
    }
    if err := cmd.Start(); err != nil {
        return errors.WithStack(err)
    }

    n, err := io.Copy(stdin, r)
    if err != nil {
        return errors.WithStack(err)
    }

    // The total number of bytes written to stdin must match sizeBytes, or
    // VBoxManage.exe on Windows will fail. Fill remaining with zeros.
    if left := sizeBytes - n; left > 0 {
        if err := ZeroFill(stdin, left); err != nil {
            return errors.WithStack(err)
        }
    }

    // cmd won't exit until the stdin is closed.
    if err := stdin.Close(); err != nil {
        return errors.WithStack(err)
    }

    return errors.WithStack(cmd.Wait())
}

// ZeroFill writes n zero bytes into w.
func ZeroFill(w io.Writer, n int64) error {
    const blocksize = 32 << 10
    zeros := make([]byte, blocksize)
    var k int
    var err error
    for n > 0 {
        if n > blocksize {
            k, err = w.Write(zeros)
        } else {
            k, err = w.Write(zeros[:n])
        }
        if err != nil {
            return err
        }
        n -= int64(k)
    }
    return nil
}

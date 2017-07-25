package vboxutil

import (
    "archive/tar"
    "bytes"
    "fmt"
    "io"
    "os"
    "os/exec"
    "path/filepath"

    log "github.com/Sirupsen/logrus"
    "github.com/pkg/errors"
    "github.com/stkim1/pc-vbox-core/crcontext/config"
)

const (
    magicString string = "pc-core, please format-me"
    DefualtCoreDiskName = "pc-core-hdd"
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
    VBM                 string
    ClusterID           string
    AuthToken           string
    UserName            string
    CoreVboxPublicKey   []byte
    CoreVboxPrivateKey  []byte
    MasterVboxPublicKey []byte
    EngineAuthCert      []byte
    EngineKeyCert       []byte
    EnginePrivateKey    []byte
}

func NewMachineDisk(baseFolder, imageName string, diskSize uint, debug bool) *MachineDisk {
    return &MachineDisk {
        Debug:               debug,
        DiskSize:            diskSize,
        DiskImage:           filepath.Join(baseFolder, fmt.Sprintf("%s.vmdk", imageName)),
        VBM:                 "/usr/local/bin/VBoxManage",
    }
}

func (m *MachineDisk) BuildCoreDiskImage() error {
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

        // cluster id
        if len(m.ClusterID) == 0 {
            return errors.Errorf("[ERR] invalid cluster id")
        }
        file = &tar.Header{Name: config.ArchivePathClusterID(), Size: int64(len(m.ClusterID)), Mode: 0600}
        if err := tw.WriteHeader(file); err != nil {
            return errors.WithStack(err)
        }
        if _, err := tw.Write([]byte(m.ClusterID)); err != nil {
            return errors.WithStack(err)
        }

        // ssh auth token
        if len(m.AuthToken) == 0 {
            return errors.Errorf("[ERR] invalid auth token")
        }
        file = &tar.Header{Name: config.ArchivePathAuthToken(), Size: int64(len(m.AuthToken)), Mode: 0600}
        if err := tw.WriteHeader(file); err != nil {
            return errors.WithStack(err)
        }
        if _, err := tw.Write([]byte(m.AuthToken)); err != nil {
            return errors.WithStack(err)
        }

        // core user name
        if len(m.UserName) == 0 {
            return errors.Errorf("[ERR] invalid user name")
        }
        file = &tar.Header{Name: config.ArchivePathUserName(), Size: int64(len(m.UserName)), Mode: 0600}
        if err := tw.WriteHeader(file); err != nil {
            return errors.WithStack(err)
        }
        if _, err := tw.Write([]byte(m.UserName)); err != nil {
            return errors.WithStack(err)
        }


        // cert dir
        file = &tar.Header{Name: config.ArchivePathCertsDir(), Typeflag: tar.TypeDir, Mode: 0700}
        if err := tw.WriteHeader(file); err != nil {
            return errors.WithStack(err)
        }

        // core vbox public key
        if len(m.CoreVboxPublicKey) == 0 {
            return errors.Errorf("[ERR] invalid core vbox public key")
        }
        file = &tar.Header{Name: config.ArchivePathCoreVboxPublicKey(), Size: int64(len(m.CoreVboxPublicKey)), Mode: 0600}
        if err := tw.WriteHeader(file); err != nil {
            return errors.WithStack(err)
        }
        if _, err := tw.Write(m.CoreVboxPublicKey); err != nil {
            return errors.WithStack(err)
        }

        // core vbox private key
        if len(m.CoreVboxPrivateKey) == 0 {
            return errors.Errorf("[ERR] invalid core vbox private key")
        }
        file = &tar.Header{Name: config.ArchivePathCoreVboxPrivateKey(), Size: int64(len(m.CoreVboxPrivateKey)), Mode: 0600}
        if err := tw.WriteHeader(file); err != nil {
            return errors.WithStack(err)
        }
        if _, err := tw.Write(m.CoreVboxPrivateKey); err != nil {
            return errors.WithStack(err)
        }

        // master vbox public key
        if len(m.MasterVboxPublicKey) == 0 {
            return errors.Errorf("[ERR] invalid master vbox public key")
        }
        file = &tar.Header{Name: config.ArchivePathMasterVboxPublicKey(), Size: int64(len(m.MasterVboxPublicKey)), Mode: 0600}
        if err := tw.WriteHeader(file); err != nil {
            return errors.WithStack(err)
        }
        if _, err := tw.Write(m.MasterVboxPublicKey); err != nil {
            return errors.WithStack(err)
        }

        // tls auth cert
        if len(m.EngineAuthCert) == 0 {
            return errors.Errorf("[ERR] invalid auth certificate")
        }
        file = &tar.Header{Name: config.ArchivePathCoreEngineAuthCert(), Size: int64(len(m.EngineAuthCert)), Mode: 0600}
        if err := tw.WriteHeader(file); err != nil {
            return errors.WithStack(err)
        }
        if _, err := tw.Write(m.EngineAuthCert); err != nil {
            return errors.WithStack(err)
        }

        // tls key certificate
        if len(m.EngineKeyCert) == 0 {
            return errors.Errorf("[ERR] invalid key certificate")
        }
        file = &tar.Header{Name: config.ArchivePathCoreEngineKeyCert(), Size: int64(len(m.EngineKeyCert)), Mode: 0600}
        if err := tw.WriteHeader(file); err != nil {
            return errors.WithStack(err)
        }
        if _, err := tw.Write(m.EngineKeyCert); err != nil {
            return errors.WithStack(err)
        }

        // tls private key
        if len(m.EnginePrivateKey) == 0 {
            return errors.Errorf("[ERR] invalid private key")
        }
        file = &tar.Header{Name: config.ArchivePathCoreEnginePrivateKey(), Size: int64(len(m.EnginePrivateKey)), Mode: 0600}
        if err := tw.WriteHeader(file); err != nil {
            return errors.WithStack(err)
        }
        if _, err := tw.Write(m.EnginePrivateKey); err != nil {
            return errors.WithStack(err)
        }


        // close archive
        if err := tw.Close(); err != nil {
            return errors.WithStack(err)
        }

        // Create the dest dir.
        if err := os.MkdirAll(filepath.Dir(m.DiskImage), 0700); err != nil {
            return errors.WithStack(err)
        }
        // Fill in the magic string so boot2docker VM will detect this and format
        // the disk upon first boot.
        if err := makeDiskImage(m, bytes.NewReader(buf.Bytes())); err != nil {
            return errors.WithStack(err)
        }

        if m.Debug {
            log.Debugf("Initializing disk with pre-generated core data size %d\n", len(buf.Bytes()))
            log.Debugf("WRITING: %s\n-----\n", buf)
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
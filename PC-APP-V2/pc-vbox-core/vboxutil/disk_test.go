package vboxutil

import (
    "io"
    "os"
    "path/filepath"

    . "gopkg.in/check.v1"
)

// Copy disk image from given source path to destination
func copyDiskImage(dst, src string) (err error) {
    // Open source disk image
    srcImg, err := os.Open(src)
    if err != nil {
        return err
    }
    defer func() {
        if ee := srcImg.Close(); ee != nil {
            err = ee
        }
    }()
    dstImg, err := os.Create(dst)
    if err != nil {
        return err
    }
    defer func() {
        if ee := dstImg.Close(); ee != nil {
            err = ee
        }
    }()
    _, err = io.Copy(dstImg, srcImg)
    return err
}


func (s *VboxUtilSuite) TestDiskCreation(c *C) {
    d := NewMachineDisk(s.dataDir, "testdisk", 20000, true)
    err := d.BuildCoreDiskImage()
    c.Assert(err, IsNil)

    diskImageClone := filepath.Join(os.Getenv("HOME"), "temp/testdisk.vmdk")
    err = copyDiskImage(diskImageClone, d.DiskImage)
    c.Assert(err, IsNil)
}
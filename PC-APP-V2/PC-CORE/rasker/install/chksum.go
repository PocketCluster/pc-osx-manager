package install

import (
    "bytes"
    "encoding/base64"

    "github.com/Redundancy/go-sync/filechecksum"
    "github.com/pkg/errors"
)

func checkMetaChksum(data []byte, refData string) error {
    if len(data) == 0 {
        return errors.Errorf("invalid data to check")
    }
    if len(refData) == 0 {
        return errors.Errorf("invalid length of reference checksum")
    }
    refChksum, err := base64.URLEncoding.DecodeString(refData)
    if err != nil {
        return errors.WithStack(err)
    }
    hasher := filechecksum.DefaultStrongHashGenerator()
    hasher.Write(data)
    if bytes.Compare(hasher.Sum(nil), refChksum) != 0 {
        return errors.Errorf("invalid checksum value")
    }
    return nil
}

func isTwoChksumSame(chksum []byte, refData string) error {
    if len(chksum) == 0 {
        return errors.Errorf("invalid length of checksum to compare")
    }
    if len(refData) == 0 {
        return errors.Errorf("invalid length of reference checksum to compare")
    }
    refChksum, err := base64.URLEncoding.DecodeString(refData)
    if err != nil {
        return errors.WithStack(err)
    }
    if bytes.Compare(chksum, refChksum) != 0 {
        return errors.Errorf("invalid checksum value")
    }
    return nil
}

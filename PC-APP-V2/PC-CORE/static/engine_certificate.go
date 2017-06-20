package main

import (
    "github.com/cloudflare/cfssl/certdb"
    "github.com/pkg/errors"

    "github.com/stkim1/pc-core/context"
    pcdefaults "github.com/stkim1/pc-core/defaults"
    "github.com/stkim1/pc-core/model"
    "github.com/stkim1/pcrypto"
)

//certificate authority generation
func buildCertAuthSigner(certRec certdb.Accessor, meta *model.ClusterMeta, country string) (*context.CertAuthBundle, error) {
    var (
        signer *pcrypto.CaSigner = nil
        prvKey []byte  = nil
        pubKey []byte  = nil
        crtPem []byte  = nil
        sshChk []byte  = nil
        err error      = nil
        caPrvRec, rerr = certRec.GetCertificate(pcdefaults.ClusterCertAuthPrivateKey, meta.ClusterUUID)
        caPubRec, uerr = certRec.GetCertificate(pcdefaults.ClusterCertAuthPublicKey, meta.ClusterUUID)
        caCrtRec, cerr = certRec.GetCertificate(pcdefaults.ClusterCertAuthCertificate, meta.ClusterUUID)
        caSshRec, serr = certRec.GetCertificate(pcdefaults.ClusterCertAuthSshCheck, meta.ClusterUUID)
    )
    if (rerr != nil || uerr != nil || cerr != nil || serr != nil) || (len(caPrvRec) == 0 || len(caPubRec) == 0 || len(caCrtRec) == 0 || len(caSshRec) == 0) {
        pubKey, prvKey, crtPem, sshChk, err = pcrypto.GenerateClusterCertificateAuthorityData(meta.ClusterDomain, country)
        if err != nil {
            return nil, errors.WithStack(err)
        }
        // save private key
        err = certRec.InsertCertificate(certdb.CertificateRecord{
            PEM:        string(prvKey),
            Serial:     pcdefaults.ClusterCertAuthPrivateKey,
            AKI:        meta.ClusterUUID,
            Status:     "good",
            Reason:     0,
        })
        if err != nil {
            return nil, errors.WithStack(err)
        }
        // save public key
        err = certRec.InsertCertificate(certdb.CertificateRecord{
            PEM:        string(pubKey),
            Serial:     pcdefaults.ClusterCertAuthPublicKey,
            AKI:        meta.ClusterUUID,
            Status:     "good",
            Reason:     0,
        })
        if err != nil {
            return nil, errors.WithStack(err)
        }
        // save certificate
        err = certRec.InsertCertificate(certdb.CertificateRecord{
            PEM:        string(crtPem),
            Serial:     pcdefaults.ClusterCertAuthCertificate,
            AKI:        meta.ClusterUUID,
            Status:     "good",
            Reason:     0,
        })
        if err != nil {
            return nil, errors.WithStack(err)
        }
        // save ssh checker
        err = certRec.InsertCertificate(certdb.CertificateRecord{
            PEM:        string(sshChk),
            Serial:     pcdefaults.ClusterCertAuthSshCheck,
            AKI:        meta.ClusterUUID,
            Status:     "good",
            Reason:     0,
        })
        if err != nil {
            return nil, errors.WithStack(err)
        }
    } else {
        prvKey = []byte(caPrvRec[0].PEM)
        pubKey = []byte(caPubRec[0].PEM)
        crtPem = []byte(caCrtRec[0].PEM)
        sshChk = []byte(caSshRec[0].PEM)
    }
    signer, err = pcrypto.NewCertAuthoritySigner(prvKey, crtPem, meta.ClusterDomain, country)
    if err != nil {
        return nil, errors.WithStack(err)
    }
    return &context.CertAuthBundle{
        CASigner:    signer,
        CAPrvKey:    prvKey,
        CAPubKey:    pubKey,
        CACrtPem:    crtPem,
        CASSHChk:    sshChk,
    }, nil
}

// host certificate
func buildHostCertificate(certRec certdb.Accessor, caSigner *pcrypto.CaSigner, hostname, clusterUUID string) (*context.HostCertBundle, error) {
    var (
        prvKey []byte  = nil
        pubKey []byte  = nil
        crtPem []byte  = nil
        sshPem []byte  = nil
        err error      = nil

        prvRec, rerr = certRec.GetCertificate(pcdefaults.MasterHostPrivateKey,  clusterUUID)
        pubRec, uerr = certRec.GetCertificate(pcdefaults.MasterHostPublicKey,   clusterUUID)
        crtRec, cerr = certRec.GetCertificate(pcdefaults.MasterHostCertificate, clusterUUID)
        sshRec, serr = certRec.GetCertificate(pcdefaults.MasterHostSshKey,      clusterUUID)
    )

    if (rerr != nil || uerr != nil || cerr != nil || serr != nil) || (len(prvRec) == 0 || len(pubRec) == 0 || len(crtRec) == 0 || len(sshRec) == 0) {
        pubKey, prvKey, sshPem, err = pcrypto.GenerateStrongKeyPair()
        if err != nil {
            return nil, errors.WithStack(err)
        }
        // we're not going to proide ip address for now
        crtPem, err = caSigner.GenerateSignedCertificate(hostname, "", prvKey)
        if err != nil {
            return nil, errors.WithStack(err)
        }
        // save private key
        err = certRec.InsertCertificate(certdb.CertificateRecord{
            PEM:        string(prvKey),
            Serial:     pcdefaults.MasterHostPrivateKey,
            AKI:        clusterUUID,
            Status:     "good",
            Reason:     0,
        })
        if err != nil {
            return nil, errors.WithStack(err)
        }
        // save public key
        err = certRec.InsertCertificate(certdb.CertificateRecord{
            PEM:        string(pubKey),
            Serial:     pcdefaults.MasterHostPublicKey,
            AKI:        clusterUUID,
            Status:     "good",
            Reason:     0,
        })
        if err != nil {
            return nil, errors.WithStack(err)
        }
        // save cert pem
        err = certRec.InsertCertificate(certdb.CertificateRecord{
            PEM:        string(crtPem),
            Serial:     pcdefaults.MasterHostCertificate,
            AKI:        clusterUUID,
            Status:     "good",
            Reason:     0,
        })
        if err != nil {
            return nil, errors.WithStack(err)
        }
        // save ssh pem
        err = certRec.InsertCertificate(certdb.CertificateRecord{
            PEM:        string(sshPem),
            Serial:     pcdefaults.MasterHostSshKey,
            AKI:        clusterUUID,
            Status:     "good",
            Reason:     0,
        })
        if err != nil {
            return nil, errors.WithStack(err)
        }
    } else {
        prvKey = []byte(prvRec[0].PEM)
        pubKey = []byte(pubRec[0].PEM)
        crtPem = []byte(crtRec[0].PEM)
        sshPem = []byte(sshRec[0].PEM)
    }
    return &context.HostCertBundle{
        PrivateKey:     prvKey,
        PublicKey:      pubKey,
        SshKey:         sshPem,
        Certificate:    crtPem,
    }, nil
}

// beacon certificate for slaves
func buildBeaconCertificate(certRec certdb.Accessor, clusterUUID string) (*context.BeaconCertBundle, error) {
    var (
        prvKey []byte  = nil
        pubKey []byte  = nil
        err error      = nil

        prvRec, rerr = certRec.GetCertificate(pcdefaults.MasterBeaconPrivateKey,  clusterUUID)
        pubRec, uerr = certRec.GetCertificate(pcdefaults.MasterBeaconPublicKey,   clusterUUID)
    )

    if (rerr != nil || uerr != nil) || (len(prvRec) == 0 || len(pubRec) == 0) {
        pubKey, prvKey, _, err = pcrypto.GenerateWeakKeyPair()
        if err != nil {
            return nil, errors.WithStack(err)
        }
        // save private key
        err = certRec.InsertCertificate(certdb.CertificateRecord{
            PEM:        string(prvKey),
            Serial:     pcdefaults.MasterBeaconPrivateKey,
            AKI:        clusterUUID,
            Status:     "good",
            Reason:     0,
        })
        if err != nil {
            return nil, errors.WithStack(err)
        }
        // save public key
        err = certRec.InsertCertificate(certdb.CertificateRecord{
            PEM:        string(pubKey),
            Serial:     pcdefaults.MasterBeaconPublicKey,
            AKI:        clusterUUID,
            Status:     "good",
            Reason:     0,
        })
        if err != nil {
            return nil, errors.WithStack(err)
        }
    } else {
        prvKey = []byte(prvRec[0].PEM)
        pubKey = []byte(pubRec[0].PEM)
    }
    return &context.BeaconCertBundle{
        PrivateKey:     prvKey,
        PublicKey:      pubKey,
    }, nil
}

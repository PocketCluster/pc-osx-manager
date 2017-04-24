package main

import (
    "database/sql"
    "time"

    "github.com/jmoiron/sqlx"
    _ "github.com/mattn/go-sqlite3"
    "github.com/cloudflare/cfssl/certdb"
    certsql "github.com/cloudflare/cfssl/certdb/sql"
    log "github.com/Sirupsen/logrus"
    "github.com/davecgh/go-spew/spew"
)

const (
    createCertTable = `CREATE TABLE certificates (
serial_number            bytea NOT NULL,
authority_key_identifier bytea NOT NULL,
ca_label                 bytea,
status                   bytea NOT NULL,
reason                   int,
expiry                   timestamp,
revoked_at               timestamp,
pem                      bytea NOT NULL,
PRIMARY KEY(serial_number, authority_key_identifier)
);`

    createOCSPTable = `CREATE TABLE ocsp_responses (
serial_number            bytea NOT NULL,
authority_key_identifier bytea NOT NULL,
body                     bytea NOT NULL,
expiry                   timestamp,
PRIMARY KEY(serial_number, authority_key_identifier),
FOREIGN KEY(serial_number, authority_key_identifier) REFERENCES certificates(serial_number, authority_key_identifier)
);`
)

func dbCreation(db *sql.DB, query string) error {
    var (
        stmt *sql.Stmt      = nil
        trans *sql.Tx       = nil
        err error           = nil
    )
    trans, err = db.Begin()
    if err != nil {
        return err
    }
    stmt, err = trans.Prepare(query)
    if err != nil {
        return err
    }
    _, err = stmt.Exec()
    if err != nil {
        return err
    }
    err = stmt.Close()
    if err != nil {
        return err
    }
    return trans.Commit()
}

func main() {
    const dbDriver string = "sqlite3"
    db, err := sql.Open(dbDriver, "testdb.db")
    if err != nil {
        log.Fatal(err.Error())
    }
    defer db.Close()

    err = dbCreation(db, createCertTable)
    if err != nil {
        log.Fatal(err.Error())
    }
    certbase := certsql.NewAccessor(sqlx.NewDb(db, dbDriver))

    expiry := time.Date(2010, time.December, 25, 23, 0, 0, 0, time.UTC)
    want := certdb.CertificateRecord{
        PEM:    "fake cert data",
        Serial: "fake serial",
        AKI:    "fake",
        Status: "good",
        Reason: 0,
        Expiry: expiry,
    }

    err = certbase.InsertCertificate(want)
    if err != nil {
        log.Fatal(err.Error())
    }

    cert, err := certbase.GetCertificate(want.Serial, want.AKI)
    if err != nil {
        log.Fatal(err)
    }
    log.Info(spew.Sdump(cert))
}
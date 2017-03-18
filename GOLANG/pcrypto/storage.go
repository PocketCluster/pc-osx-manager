package pcrypto

import (
    "database/sql"
    "fmt"

    "github.com/jmoiron/sqlx"
    // sqlite is loaded somewhere in cfssl certdb. make sure
    _ "github.com/mattn/go-sqlite3"
    "github.com/cloudflare/cfssl/certdb"
    certsql "github.com/cloudflare/cfssl/certdb/sql"
)

const (
    databaseDriver string   = "sqlite3"

    // (03/18/2017) Don't modify certificate table name for now
    certificateTable string = "certificates"
    createCertTable string  = `CREATE TABLE certificates (
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

    ocspTable string        = "ocsp_responses"
    createOCSPTable string  = `CREATE TABLE ocsp_responses (
serial_number            bytea NOT NULL,
authority_key_identifier bytea NOT NULL,
body                     bytea NOT NULL,
expiry                   timestamp,
PRIMARY KEY(serial_number, authority_key_identifier),
FOREIGN KEY(serial_number, authority_key_identifier) REFERENCES certificates(serial_number, authority_key_identifier)
);`
)

func concatErrors(errs... error) error {
    var (
        err error = nil
    )
    if len(errs) == 0 {
        return err
    }
    for _, e := range errs {
        if e != nil {
            if err == nil {
                err = e
            } else {
                err = fmt.Errorf("%v | %v", err, e)
            }
        }
    }
    return err
}

// check if every table in the list exits
func checkTablesExist(db *sql.DB, table string) (int, error) {
    var (
        counter int       = 0
        stmt *sql.Stmt    = nil
        err error         = nil
    )
    stmt, err = db.Prepare(fmt.Sprintf("SELECT count(name) FROM sqlite_master WHERE type = 'table' AND name = '%s'", table))
    if err != nil {
        return 0, err
    }
    err = stmt.QueryRow().Scan(&counter);
    err = concatErrors(err, stmt.Close())
    return counter, err
}

func upsertTable(db *sql.DB, tableQuery string) error {
    var (
        stmt *sql.Stmt      = nil
        trans *sql.Tx       = nil
        err error           = nil
    )
    trans, err = db.Begin()
    if err != nil {
        return err
    }
    stmt, err = trans.Prepare(tableQuery)
    if err == nil {
        _, err = stmt.Exec()
        err = concatErrors(err, stmt.Close())
    }
    if err == nil {
        err = trans.Commit()
    } else {
        err = concatErrors(err, trans.Rollback())
    }
    return err
}

func NewPocketCertStorage(db *sql.DB) (certdb.Accessor, error) {
    counter, err := checkTablesExist(db, certificateTable)
    if err != nil {
        return nil, err
    }
    if counter == 0 {
        err := upsertTable(db, createCertTable)
        if err != nil {
            return nil, err
        }
    }
    counter, err = checkTablesExist(db, ocspTable)
    if err != nil {
        return nil, err
    }
    if counter == 0 {
        err = upsertTable(db, createOCSPTable)
        if err != nil {
            return nil, err
        }
    }
    return certsql.NewAccessor(sqlx.NewDb(db, databaseDriver)), nil
}

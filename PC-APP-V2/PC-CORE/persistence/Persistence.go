package persistence

import "time"

const NodeCreation string = `DROP TABLE IF EXISTS "pcnode"; CREATE TABLE "pcnode" ("registered_at" timestamp, "deleted_at" timestamp, "unique_id" varchar DEFAULT NULL, "host_name" varchar DEFAULT NULL, "mac_address" varchar DEFAULT NULL, "device_type" varchar DEFAULT NULL);`

type PCNode struct {
    // Registered time
    RegisteredAt    time.Time   `db:registered_at`
    // De-registered time
    DeletedAt       time.Time   `db:deleted_at`
    // Some very unique identity
    UniqueId        string      `db:unique_id`
    // Hostname
    HostName        string      `db:host_name`
    // MAC address
    MacAddress      string      `db:mac_address`
    // Type : Vbox / RPI / Odroid / etc
    DeviceType      string      `db:device_type`
}

/*

import (
    "log"
    "time"
    "fmt"

    "upper.io/db.v2/sqlite"
)

type Birthday struct {
    Name string         `db:"name"`
    Born time.Time      `db:"born"`
}

func CreateSampleDB() {
    var settings = sqlite.ConnectionURL{
        Database: "pc-data.db", // Path to a sqlite3 database file.
    }

    sess, err := sqlite.Open(settings); if err != nil {
        log.Fatalf("db.Open(): %q\n", err)
    }
    defer sess.Close()

    // Creation & Migration should stay in
    sess.Exec(`DROP TABLE IF EXISTS "birthday";`)
    sess.Exec(`CREATE TABLE "birthday" ("name" varchar(50) DEFAULT NULL, "born" timestamp);`)

    birthdayCollection := sess.Collection("birthday")

    birthdayCollection.Insert(Birthday{
        Name: "Hayao Miyazaki",
        Born: time.Date(1941, time.January, 5, 0, 0, 0, 0, time.Local),
    })

    var birthdays []Birthday
    err = birthdayCollection.Find().All(&birthdays); if err != nil {
        log.Fatalf("res.All(): %q\n", err)
    }

    // Printing to stdout.
    for _, birthday := range birthdays {
        fmt.Printf("%s was born in %s.\n",
            birthday.Name,
            birthday.Born.Format("January 2, 2006"),
        )
    }

}
*/
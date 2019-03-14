package database

import (
	// System packages.
	"database/sql"
	"regexp"

	// Jotter packages.
	"github.com/dmiprops/jotter/modules/setting"

	// Vendor packages.
	_ "github.com/lib/pq"
)

var (
	dbConn *sql.DB
)

// AtStart create database and open connection.
func AtStart() error {
	return Connect()
}

// Connect initialises connection to database.
func Connect() error {
	// Check/create database 'jotter'.
	err := checkDatabase()
    if err != nil {
        return err
    }

	// Connect database 'jotter'
	connStr := "postgres://" + setting.CurrentAdminSettings.Database + "/jotter?sslmode=disable"
	dbConn, err = sql.Open("postgres", connStr)
	if err != nil {
		return nil
	}

	// Check/create tables.
    return checkDatabaseSchema()
}

func checkDatabase() error {
	connStr := "postgres://" + setting.CurrentAdminSettings.Database + "/postgres?sslmode=disable"
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        return err
    }
	defer db.Close()

	rows, err := db.Query("select datname from pg_database where datname = 'jotter';")
    if err != nil {
        return err
    }
	defer rows.Close()

	if !rows.Next() {
		err = createJotterDatabase(db)
		if err != nil {
			return err
		}
	}

	return nil
}

func createJotterDatabase(db *sql.DB) error {
	// Get connection user.
	reg := regexp.MustCompile(`^.+(:)`)
	usr := reg.FindString(setting.CurrentAdminSettings.Database)
	if usr == "" {
		usr = "jotter"
	} else {
		usr = usr[0:len(usr)-1]
	}

	// Create database.
	_, err := db.Exec("create database jotter with owner " + usr + " encoding = 'UTF8' connection limit = -1;")
	if err != nil {
		return err
	}
	return nil
}

func checkDatabaseSchema() error {
	rows, err := dbConn.Query("select table_name from information_schema.tables where table_schema='public';")
    if err != nil {
        return err
    }
	defer rows.Close()

	if !rows.Next() {
		err = createTables()
		if err != nil {
			return err
		}
	}

	return nil
}

func createTables() error {
	// Table 'settings'.
	_, err := dbConn.Exec(
		`create table public.settings
		(
			version varchar(20) NOT NULL,
			created timestamp NOT NULL,
			updated timestamp NOT NULL
		)
		WITHOUT OIDS;`,
	)
	if err != nil {
		return err
	}
	_, err = dbConn.Exec(
		`insert into settings (version, created, updated) values ($1, now(), now())`,
		setting.AppVer,
	)

	return nil
}
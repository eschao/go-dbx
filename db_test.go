package dbx

import (
	"flag"
	"fmt"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

const TEST_DB_FILE = "test_db.db"

// Defines test database
var tDBType = "sqlite"
var tDBName = "test_english_data"
var tDBHost = "localhost"
var tDBPort = 3306
var tDBUser = "root"
var tDBPasswd = ""
var tDatabase = NewDatabase()

func init() {
	flag.StringVar(&tDBType, "dbType", "sqlite", "database type")
	flag.StringVar(&tDBName, "dbName", "test_english_data", "database name")
	flag.StringVar(&tDBHost, "dbHost", "localhost", "database hostname")
	flag.StringVar(&tDBUser, "dbUser", "root", "database login user")
	flag.StringVar(&tDBPasswd, "dbPasswd", "", "database login password")
	flag.IntVar(&tDBPort, "dbPort", 3306, "databae port")
	flag.Parse()
}

func prepare() error {
	if tDBType == "sqlite" {
		_, err := os.Stat(TEST_DB_FILE)
		if os.IsExist(err) {
			os.Remove(TEST_DB_FILE)
		}

		os.Create(TEST_DB_FILE)
		if err := tDatabase.OpenSQLite(TEST_DB_FILE); err != nil {
			return err
		}
	} else if tDBType == "mysql" {
		err := tDatabase.OpenMySQL(tDBName, tDBUser, tDBPasswd, tDBHost, tDBPort)
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("Can't support database type: %s", tDBType)
	}

	if err := tDatabase.Register(&User{}); err != nil {
		return fmt.Errorf("Can't register table: User. Error: %s", err.Error())
	}
	if err := tDatabase.Register(&UserLogin{}); err != nil {
		return fmt.Errorf("Can't register table: UserLogin. Error: %s", err.Error())
	}
	if err := tDatabase.Register(&UserOAuth{}); err != nil {
		return fmt.Errorf("Can't register table: UserOAuth. Error: %s", err.Error())
	}
	return nil
}

func cleanUp() {
	tDatabase.Close()
}

func TestMain(m *testing.M) {
	err := prepare()
	if err != nil {
		fmt.Printf("Failed to prepare database with: %s", err.Error())
	} else {
		fmt.Printf("=== Using driver: %s\n", tDatabase.DriverName())
		m.Run()
	}
	cleanUp()
}

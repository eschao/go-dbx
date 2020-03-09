package dbx

import (
	"database/sql"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

const (
	DRIVER_SQLITE3 = "sqlite3"
	DRIVER_MYSQL   = "mysql"
	DRIVER_POSTGRE = "postgre"
)

const (
	defaultMaxMemory = 32 << 20 // 32 MB
)

var dbLogger = func(sql string) {
	fmt.Printf("SQL: %s\n", sql)
}

func SetLogger(logger func(string)) {
	dbLogger = logger
}

//type ITable interface {
//GetTableName() string
//GetRowType() reflect.Type
//}

type Column struct {
	Name            string
	FormName        string
	Index           int
	Sqlite          string
	Mysql           string
	Postgre         string
	IsPrimaryKey    bool
	IsAutoIncrement bool
}

type Table struct {
	Name    string
	Columns map[string]Column
	//RowType reflect.Type
}

func (this *Table) ColumnNames() []string {
	size := len(this.Columns)
	names := make([]string, size, size)
	i := 0
	for k, _ := range this.Columns {
		names[i] = k
		i++
	}
	return names
}

func (this *Table) ColumnIndexes() []int {
	size := len(this.Columns)
	indexes := make([]int, size, size)
	i := 0
	for _, v := range this.Columns {
		indexes[i] = v.Index
		i++
	}
	return indexes
}

func (this *Table) GetColumnsFromForm(r *http.Request) (
	[]string, []interface{}, error) {
	columns := []string{}
	values := []interface{}{}
	if err := r.ParseMultipartForm(defaultMaxMemory); err != nil {
		return columns, values, err
	}

	for _, c := range this.Columns {
		name := c.FormName
		if name == "" {
			name = c.Name
		}

		vs, ok := r.Form[name]
		if ok && len(vs) > 0 {
			columns = append(columns, name)
			values = append(values, vs[0])
		}
	}
	return columns, values, nil
}

func (this *Table) GetColumnsMapFromForm(r *http.Request) (
	map[string]interface{}, error) {
	columns := map[string]interface{}{}
	if err := r.ParseMultipartForm(defaultMaxMemory); err != nil {
		return columns, err
	}

	for _, c := range this.Columns {
		name := c.FormName
		if name == "" {
			name = c.Name
		}

		vs, ok := r.Form[name]
		if ok && len(vs) > 0 {
			columns[name] = vs[0]
		}
	}
	return columns, nil
}

func (this *Table) Parse(name string, table interface{}) error {
	v := reflect.TypeOf(table)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return fmt.Errorf("table is not a struct type: %v", v.Kind())
	}

	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		col := f.Tag.Get("column")
		if col == "" {
			col = f.Tag.Get("col")
		}
		if col == "" {
			col = f.Tag.Get("db")
		}
		if col == "" {
			continue
		}

		// get form id defined
		form := f.Tag.Get("form")

		isPrimaryKey := false
		isAutoIncrement := false
		sqlite := f.Tag.Get("sqlite")
		if sqlite == "" {
			sqlite = f.Tag.Get("sqlite3")
		}
		if sqlite != "" {
			s := strings.ToLower(sqlite)
			if strings.Contains(s, "primary key") {
				isPrimaryKey = true
			}
			if strings.Contains(s, "autoincrement") {
				isAutoIncrement = true
			}
		}

		mysql := f.Tag.Get("mysql")
		if mysql != "" {
			s := strings.ToLower(mysql)
			if sqlite != "" && isPrimaryKey != strings.Contains(s, "primary key") {
				return fmt.Errorf("column %s has different 'primary key' attribute",
					col)
			}
			if sqlite != "" && isAutoIncrement != strings.Contains(s, "auto_increment") {
				return fmt.Errorf("column %s has different 'auto-increment' attribute",
					col)
			}
		}

		postgre := f.Tag.Get("postgre")
		if sqlite == "" && mysql == "" && postgre == "" {
			return fmt.Errorf("column %s does not have sql definition", col)
		}

		if _, ok := this.Columns[col]; ok {
			return fmt.Errorf("column %s is redefined", col)
		}
		this.Columns[col] = Column{
			col, form, i, sqlite, mysql, postgre, isPrimaryKey, isAutoIncrement,
		}
	}

	if len(this.Columns) < 1 {
		return fmt.Errorf("table doesn't have column definitions")
	}

	this.Name = name
	return nil
}

func (this *Table) CreateSQL(driver string) (string, error) {
	cols := []string{}
	for k, v := range this.Columns {
		if driver == DRIVER_SQLITE3 {
			if v.Sqlite == "" {
				return "", fmt.Errorf(
					"%s column of %s table has no definition for sqlite", k, this.Name)
			}
			cols = append(cols, k+" "+v.Sqlite)
		} else if driver == DRIVER_MYSQL {
			if v.Mysql == "" {
				return "", fmt.Errorf(
					"%s column of %s table has no definition for mysql", k, this.Name)
			}
			cols = append(cols, k+" "+v.Mysql)
		} else if driver == DRIVER_POSTGRE {
			if v.Postgre == "" {
				return "", fmt.Errorf(
					"%s column of %s table has no definition for postgre", k, this.Name)
			}
			cols = append(cols, k+" "+v.Postgre)
		} else {
			return "", fmt.Errorf("unsupportted driver %s", driver)
		}
	}

	return "CREATE TABLE IF NOT EXISTS " + this.Name + "(" +
		strings.Join(cols, ",") + ")", nil
}

type Database struct {
	driver string
	db     *sql.DB
	tables map[string]Table
}

func NewDatabase() *Database {
	return &Database{tables: map[string]Table{}}
}

func (this *Database) Open(driver, dsn string) error {
	db, err := sql.Open(driver, dsn)
	if err == nil {
		this.driver = driver
		this.db = db
	}
	return err
}

func (this *Database) Close() {
	if this.db != nil {
		this.db.Close()
		this.db = nil
	}
}

func (this *Database) DB() *sql.DB {
	return this.db
}

func (this *Database) DriverName() string {
	return this.driver
}

func (this *Database) OpenSQLite(dbFile string) error {
	db, err := sql.Open(DRIVER_SQLITE3, dbFile)
	if err == nil {
		this.driver = DRIVER_SQLITE3
		this.db = db
	}
	return err
}

func (this *Database) OpenMySQL(dbName, user, passwd, host string, port int) error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", user, passwd, host, port, dbName)
	db, err := sql.Open(DRIVER_MYSQL, dsn)
	if err == nil {
		this.driver = DRIVER_MYSQL
		this.db = db
	}
	return err
}

func (this *Database) RegisterTable(name string, table interface{}) error {
	t := Table{Columns: map[string]Column{}}
	if err := t.Parse(name, table); err != nil {
		return err
	}

	if _, ok := this.tables[name]; ok {
		return fmt.Errorf("%s table is already existing", name)
	}
	this.tables[name] = t
	return nil
}

func (this *Database) GetTableSchema(name string) (Table, error) {
	table, ok := this.tables[name]
	if ok {
		return table, nil
	}
	return Table{}, nil
}

func (this *Database) CreateTables() error {
	if this.db == nil {
		return fmt.Errorf("no opened database")
	}

	for _, v := range this.tables {
		sql, err := v.CreateSQL(this.driver)
		if err != nil {
			return err
		}
		if dbLogger != nil {
			dbLogger(sql)
		}

		if _, err := this.db.Exec(sql); err != nil {
			return err
		}
	}
	return nil
}

func (this *Database) CreateTable(name string) error {
	if this.db == nil {
		return fmt.Errorf("no opened database")
	}

	t, ok := this.tables[name]
	if !ok {
		return fmt.Errorf("%s table is not registered", name)
	}
	sql, err := t.CreateSQL(this.driver)
	if err != nil {
		return err
	}
	if dbLogger != nil {
		dbLogger(sql)
	}

	_, err = this.db.Exec(sql)
	return err
}

func (this *Database) DropTable(name string) error {
	if this.db == nil {
		return fmt.Errorf("no opened database")
	}

	_, err := this.db.Exec("DROP TABLE " + name)
	return err
}

func (this *Database) T(name string) *SQLExecutor {
	t, ok := this.tables[name]
	var err error
	if !ok {
		err = fmt.Errorf("%s table is not registered", name)
	}

	return &SQLExecutor{
		table: &t, db: this.db, err: err,
		tableGetter: func(name string) *Table {
			t, _ := this.tables[name]
			return &t
		},
	}
}

func (this *Database) Begin() (*Transaction, error) {
	if this.db == nil {
		return nil, fmt.Errorf("no opened database")
	}

	tx, err := this.db.Begin()
	if err != nil {
		return nil, err
	}
	return &Transaction{db: this, tx: tx}, nil
}

//
// Database Transaction
//
type Transaction struct {
	tx *sql.Tx
	db *Database
}

func (this *Transaction) Tx() *sql.Tx {
	return this.tx
}

func (this *Transaction) Commit() error {
	if this.tx == nil {
		return fmt.Errorf("Nil transaction")
	}
	return this.tx.Commit()
}

func (this *Transaction) Rollback() error {
	if this.tx == nil {
		return fmt.Errorf("Nil transaction")
	}
	return this.tx.Rollback()
}

func (this *Transaction) T(name string) *SQLExecutor {
	t, ok := this.db.tables[name]
	var err error
	if !ok {
		err = fmt.Errorf("table %s is not registered", name)
	}

	return &SQLExecutor{
		table: &t, tx: this.tx, err: err,
		tableGetter: func(name string) *Table {
			t, _ := this.db.tables[name]
			return &t
		},
	}
}

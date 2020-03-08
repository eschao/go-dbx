package dbx

import (
	"database/sql"
	"fmt"
	"reflect"
)

// sqlFilter
type sqlFilter struct {
	where string
	args  []interface{}
}

type tableGetter func(string) *Table

// SQLExecutor
type SQLExecutor struct {
	table       *Table
	db          *sql.DB
	tx          *sql.Tx
	err         error
	tableGetter tableGetter
}

// Insert inserts given row to table
func (this *SQLExecutor) Insert(row interface{}) (sql.Result, error) {
	if this.err != nil {
		return nil, this.err
	}

	cols := ""
	vals := ""
	refs := make([]interface{}, 0, len(this.table.Columns))
	rowVal := reflect.ValueOf(row).Elem()

	for k, v := range this.table.Columns {
		if !v.IsAutoIncrement {
			cols += k + ","
			vals += "?,"
			refs = append(refs, rowVal.Field(v.Index).Interface())
		}
	}

	if cols == "" {
		return nil, fmt.Errorf("table doesn't have columns")
	}

	cols = cols[:len(cols)-1]
	vals = vals[:len(vals)-1]
	q := "INSERT INTO " + this.table.Name + "(" + cols + ") VALUES(" + vals + ")"
	if dbLogger != nil {
		dbLogger(q)
	}

	var stmt *sql.Stmt
	var err error
	if this.tx != nil {
		stmt, err = this.tx.Prepare(q)
	} else {
		stmt, err = this.db.Prepare(q)
	}
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	return stmt.Exec(refs...)
}

// CountAll counts all rows of table
func (this *SQLExecutor) CountAll() (int, error) {
	if this.err != nil {
		return 0, this.err
	}

	q := "SELECT COUNT(*) as count FROM " + this.table.Name
	if dbLogger != nil {
		dbLogger(q)
	}

	var rs *sql.Rows
	var err error
	if this.tx != nil {
		rs, err = this.tx.Query(q)
	} else {
		rs, err = this.db.Query(q)
	}

	if err == sql.ErrNoRows {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	defer rs.Close()

	var count = 0
	if rs.Next() {
		if err := rs.Scan(&count); err != nil {
			return 0, err
		}
	}
	return count, nil
}

// Count counts rows by given filter
func (this *SQLExecutor) Count(where string, args ...interface{}) (int, error) {
	if this.err != nil {
		return 0, this.err
	}

	q := "SELECT COUNT(*) as count FROM " + this.table.Name
	if where != "" {
		q += " WHERE " + where
	}
	if dbLogger != nil {
		dbLogger(q)
	}

	var rs *sql.Rows
	var err error
	if this.tx != nil {
		rs, err = this.tx.Query(q, args...)
	} else {
		rs, err = this.db.Query(q, args...)
	}

	if err == sql.ErrNoRows {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	defer rs.Close()

	var count = 0
	if rs.Next() {
		if err := rs.Scan(&count); err != nil {
			return 0, err
		}
	}
	return count, nil
}

// SelectAll selects all columns from table
func (this *SQLExecutor) SelectAll() *SQLSelector {
	return &SQLSelector{
		table: this.table, db: this.db, tx: this.tx, err: this.err,
		tableGetter: this.tableGetter,
		columns:     this.table.ColumnNames(),
		filter:      sqlFilter{args: []interface{}{}},
		sort:        sqlSort{columns: []string{}},
	}
}

// Select selects the given columns from table
func (this *SQLExecutor) Select(cols ...string) *SQLSelector {
	return &SQLSelector{
		table: this.table, db: this.db, tx: this.tx, err: this.err, columns: cols,
		tableGetter: this.tableGetter,
		filter:      sqlFilter{args: []interface{}{}},
		sort:        sqlSort{columns: []string{}},
	}
}

// Delete deletes rows by given filter
func (this *SQLExecutor) Delete(where string, args ...interface{}) error {
	if this.err != nil {
		return this.err
	}

	q := "DELETE FROM " + this.table.Name
	if where != "" {
		q += " WHERE " + where
	}
	if dbLogger != nil {
		dbLogger(q)
	}

	var err error
	if this.tx != nil {
		_, err = this.tx.Exec(q, args...)
	} else {
		_, err = this.db.Exec(q, args...)
	}
	return err
}

// Replace replaces with given row
func (this *SQLExecutor) Replace(row interface{}) (sql.Result, error) {
	cols := ""
	vals := ""
	size := len(this.table.Columns)
	refs := make([]interface{}, size, size)
	i := 0
	rowVal := reflect.ValueOf(row).Elem()
	for k, v := range this.table.Columns {
		cols += k + ","
		vals += "?,"
		refs[i] = rowVal.Field(v.Index).Addr().Interface()
		i++
	}

	if cols == "" {
		return nil, fmt.Errorf("table doesn't have columns")
	}
	cols = cols[:len(cols)-1]
	vals = vals[:len(vals)-1]
	q := "REPLACE INTO " + this.table.Name + "(" + cols + ") VALUES(" + vals + ")"
	if dbLogger != nil {
		dbLogger(q)
	}

	var stmt *sql.Stmt
	var err error
	if this.tx != nil {
		stmt, err = this.tx.Prepare(q)
	} else {
		stmt, err = this.db.Prepare(q)
	}
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	return stmt.Exec(refs...)
}

// Update updates row by given filter
func (this *SQLExecutor) Update(where string, args ...interface{}) *SQLUpdater {
	return &SQLUpdater{
		table: this.table, db: this.db, tx: this.tx, err: this.err,
		filter: sqlFilter{where: where, args: args},
	}
}

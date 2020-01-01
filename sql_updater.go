package dbx

import (
	"database/sql"
	"fmt"
	"reflect"
)

// SQLUpdater
type SQLUpdater struct {
	table   *Table
	db      *sql.DB
	tx      *sql.Tx
	err     error
	columns []string
	filter  sqlFilter
}

// Set sets columns to be updated
func (this *SQLUpdater) Set(cols ...string) *SQLUpdater {
	this.columns = cols
	return this
}

// Values updates given values to table
func (this *SQLUpdater) Values(values ...interface{}) (sql.Result, error) {
	if this.err != nil {
		return nil, this.err
	}

	// if not given columns, all columns will be updated
	size := len(this.columns)
	if size < 1 {
		return nil, fmt.Errorf("please specify columns to update")
	}

	if len(values) != size {
		return nil, fmt.Errorf("the specified columns and values are not equal")
	}

	cols := ""
	for _, n := range this.columns {
		if n != "" {
			cols += n + "=?,"
		}
	}

	if cols == "" {
		return nil, fmt.Errorf("no specified columns to update")
	}
	cols = cols[:len(cols)-1]
	q := "UPDATE " + this.table.Name + " SET " + cols
	if this.filter.where != "" {
		q += " WHERE " + this.filter.where
		values = append(values, this.filter.args...)
	}
	if DBLogger != nil {
		DBLogger(q)
	}

	if this.tx != nil {
		return this.tx.Exec(q, values...)
	} else {
		return this.db.Exec(q, values...)
	}
}

// Value updates given row to table
func (this *SQLUpdater) Value(row interface{}) (sql.Result, error) {
	if this.err != nil {
		return nil, this.err
	}

	// if not given columns, all columns will be updated
	size := len(this.columns)
	if size < 1 {
		this.columns = this.table.ColumnNames()
		size = len(this.columns)
	}

	cols := ""
	vals := []interface{}{}
	rowVal := reflect.ValueOf(row).Elem()
	for _, n := range this.columns {
		col, ok := this.table.Columns[n]
		if !ok {
			return nil, fmt.Errorf("column %s is not found", n)
		}

		if !col.IsAutoIncrement {
			cols += n + "=?,"
			vals = append(vals, rowVal.Field(col.Index).Interface())
		}
	}

	if cols == "" {
		return nil, fmt.Errorf("no specified columns to update")
	}
	cols = cols[:len(cols)-1]
	q := "UPDATE " + this.table.Name + " SET " + cols
	if this.filter.where != "" {
		q += " WHERE " + this.filter.where
		vals = append(vals, this.filter.args...)
	}
	if DBLogger != nil {
		DBLogger(q)
	}

	if this.tx != nil {
		return this.tx.Exec(q, vals...)
	} else {
		return this.db.Exec(q, vals...)
	}
}

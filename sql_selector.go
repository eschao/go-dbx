package dbx

import (
	"database/sql"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type sqlSort struct {
	op      string
	columns []string
}

func (this sqlSort) buildSQL(table string) string {
	s := ""
	for _, c := range this.columns {
		s += table + "." + c + ","
	}

	if s != "" {
		return s[:len(s)-1]
	} else {
		return ""
	}
}

// SQLSelector
type SQLSelector struct {
	table       *Table
	db          *sql.DB
	tx          *sql.Tx
	err         error
	columns     []string
	filter      sqlFilter
	limit       int
	offset      int
	sort        sqlSort
	tableGetter tableGetter
}

func (this *SQLSelector) buildColumnsSQL() string {
	s := ""
	name := this.table.Name
	for _, c := range this.columns {
		s += name + "." + c + ","
	}
	if s != "" {
		return s[:len(s)-1]
	} else {
		return ""
	}
}

func (this *SQLSelector) buildColumnsSQLRefs() (string, []int) {
	s := ""
	indexes := []int{}
	name := this.table.Name
	for _, n := range this.columns {
		c := this.table.Columns[n]
		s += name + "." + n + ","
		indexes = append(indexes, c.Index)
	}
	if s != "" {
		return s[:len(s)-1], indexes
	} else {
		return "", indexes
	}
}

// Filter set filters for select
func (this *SQLSelector) Filter(where string, args ...interface{}) *SQLSelector {
	this.filter.where = where
	this.filter.args = args
	return this
}

// Asc sorts given columns by asc
func (this *SQLSelector) Asc(cols ...string) *SQLSelector {
	this.sort.columns = cols
	this.sort.op = "ASC"
	return this
}

// Desc sorts given columns by desc
func (this *SQLSelector) Desc(cols ...string) *SQLSelector {
	this.sort.columns = cols
	this.sort.op = "DESC"
	return this
}

// Limit sets limit of query
func (this *SQLSelector) Limit(n int) *SQLSelector {
	this.limit = n
	return this
}

// Offset sets offset of query
func (this *SQLSelector) Offset(n int) *SQLSelector {
	this.offset = n
	return this
}

func (this *SQLSelector) LeftJoin(table, onLeft, onRight string) jointer {
	return func() (*SQLJointer, *sqlJoin) {
		return &SQLJointer{selector: this, joins: []sqlJoin{}}, &sqlJoin{
			op: " LEFT JOIN ", table: table, onLeft: onLeft, onRight: onRight,
			columns: []string{},
		}
	}
}

func (this *SQLSelector) LeftOuterJoin(table, onLeft, onRight string) jointer {
	return func() (*SQLJointer, *sqlJoin) {
		return &SQLJointer{selector: this, joins: []sqlJoin{}}, &sqlJoin{
			op: " LEFT OUTER JOIN ", table: table, onLeft: onLeft, onRight: onRight,
			columns: []string{},
		}
	}
}

func (this *SQLSelector) RightJoin(table, onLeft, onRight string) jointer {
	return func() (*SQLJointer, *sqlJoin) {
		return &SQLJointer{selector: this, joins: []sqlJoin{}}, &sqlJoin{
			op: " RIGHT JOIN ", table: table, onLeft: onLeft, onRight: onRight,
			columns: []string{},
		}
	}
}

func (this *SQLSelector) RightOuterJoin(table, onLeft, onRight string) jointer {
	return func() (*SQLJointer, *sqlJoin) {
		return &SQLJointer{selector: this, joins: []sqlJoin{}}, &sqlJoin{
			op: " RIGHT OUTER JOIN ", table: table, onLeft: onLeft, onRight: onRight,
			columns: []string{},
		}
	}
}

func (this *SQLSelector) FullJoin(table, onLeft, onRight string) jointer {
	return func() (*SQLJointer, *sqlJoin) {
		return &SQLJointer{selector: this, joins: []sqlJoin{}}, &sqlJoin{
			op: " FULL JOIN ", table: table, onLeft: onLeft, onRight: onRight,
			columns: []string{},
		}
	}
}

func (this *SQLSelector) FullOuterJoin(table, onLeft, onRight string) jointer {
	return func() (*SQLJointer, *sqlJoin) {
		return &SQLJointer{selector: this, joins: []sqlJoin{}}, &sqlJoin{
			op: " FULL OUTER JOIN ", table: table, onLeft: onLeft, onRight: onRight,
			columns: []string{},
		}
	}
}

func (this *SQLSelector) InnerJoin(table, onLeft, onRight string) jointer {
	return func() (*SQLJointer, *sqlJoin) {
		return &SQLJointer{selector: this, joins: []sqlJoin{}}, &sqlJoin{
			op: " INNER JOIN ", table: table, onLeft: onLeft, onRight: onRight,
			columns: []string{},
		}
	}
}

func (this *SQLSelector) buildSQL() string {
	q := "SELECT " + strings.Join(this.columns, ",") + " FROM " + this.table.Name
	if this.filter.where != "" {
		q += " WHERE " + this.filter.where
	}
	if len(this.sort.columns) > 0 {
		q += " ORDER BY " + strings.Join(this.sort.columns, ",") + " " + this.sort.op
	}
	if this.limit > 0 {
		q += " LIMIT " + strconv.Itoa(this.limit)
	}
	if this.offset > 0 {
		q += " OFFSET " + strconv.Itoa(this.offset)
	}

	if DBLogger != nil {
		DBLogger(q)
	}
	return q
}

func (this *SQLSelector) Run() (*sql.Rows, error) {
	if this.err != nil {
		return nil, this.err
	}

	q := this.buildSQL()
	if this.tx != nil {
		return this.tx.Query(q, this.filter.args...)
	} else {
		return this.db.Query(q, this.filter.args...)
	}
}

// One selects one row from table
func (this *SQLSelector) One(row interface{}) error {
	if this.err != nil {
		return this.err
	}

	size := len(this.columns)
	refs := make([]interface{}, size, size)
	rowVal := reflect.ValueOf(row).Elem()
	for i, n := range this.columns {
		col := this.table.Columns[n]
		refs[i] = rowVal.Field(col.Index).Addr().Interface()
	}

	q := this.buildSQL()
	var rs *sql.Row
	if this.tx != nil {
		rs = this.tx.QueryRow(q, this.filter.args...)
	} else {
		rs = this.db.QueryRow(q, this.filter.args...)
	}
	return rs.Scan(refs...)
}

// All selects all rows from table
func (this *SQLSelector) All(rows interface{}) error {
	if this.err != nil {
		return this.err
	}

	rowsVal := reflect.ValueOf(rows)
	if rowsVal.Kind() != reflect.Ptr || rowsVal.Elem().Kind() != reflect.Slice {
		return fmt.Errorf("rows argument must be a slice address")
	}

	size := len(this.columns)
	indexes := make([]int, size, size)
	for i, n := range this.columns {
		col := this.table.Columns[n]
		indexes[i] = col.Index
	}

	q := this.buildSQL()
	var rs *sql.Rows
	var err error
	if this.tx != nil {
		rs, err = this.tx.Query(q, this.filter.args...)
	} else {
		rs, err = this.db.Query(q, this.filter.args...)
	}
	if err != nil {
		return err
	}
	defer rs.Close()

	sliceVal := rowsVal.Elem()
	sliceVal = sliceVal.Slice(0, sliceVal.Cap())
	rowType := sliceVal.Type().Elem()
	refs := make([]interface{}, size, size)
	i := 0

	for rs.Next() {
		if sliceVal.Len() == i {
			row := reflect.New(rowType).Elem()
			//fmt.Printf("New Row: %v\n", row)
			for k, j := range indexes {
				refs[k] = row.Field(j).Addr().Interface()
			}
			if err := rs.Scan(refs...); err != nil {
				return err
			}
			sliceVal = reflect.Append(sliceVal, row)
		} else {
			row := sliceVal.Index(i)
			//fmt.Printf("Row At: %v\n", row)
			for k, j := range indexes {
				refs[k] = row.Field(j).Addr().Interface()
			}
			if err := rs.Scan(refs...); err != nil {
				return err
			}
		}
		i++
	}

	rowsVal.Elem().Set(sliceVal) //.Slice(0, i))
	return nil
}

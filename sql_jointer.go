package dbx

import (
	"database/sql"
	"fmt"
	"reflect"
	"strconv"
	//"strings"
)

type sqlJoin struct {
	op      string
	table   string
	onLeft  string
	onRight string
	where   string
	columns []string
}

type jointer func() (*SQLJointer, *sqlJoin)

func (f jointer) SelectAll() *SQLJointer {
	jointer, join := f()
	table := jointer.selector.tableGetter(join.table)
	if table != nil {
		join.columns = table.ColumnNames()
	}
	jointer.joins = append(jointer.joins, *join)
	return jointer
}

func (f jointer) Select(cols ...string) *SQLJointer {
	jointer, join := f()
	if cols != nil {
		join.columns = cols
	}
	jointer.joins = append(jointer.joins, *join)
	return jointer
}

type joinColumns struct {
	columns string
	indexes []int
}

type SQLJointer struct {
	selector *SQLSelector
	joins    []sqlJoin
}

// Filter set filters for select
func (this *SQLJointer) Filter(where string, args ...interface{}) *SQLJointer {
	this.selector.filter.where = where
	this.selector.filter.args = args
	return this
}

// Asc sorts given columns by asc
func (this *SQLJointer) Asc(cols ...string) *SQLJointer {
	this.selector.sort.columns = cols
	this.selector.sort.op = "ASC"
	return this
}

// Desc sorts given columns by desc
func (this *SQLJointer) Desc(cols ...string) *SQLJointer {
	this.selector.sort.columns = cols
	this.selector.sort.op = "DESC"
	return this
}

// Limit sets limit of query
func (this *SQLJointer) Limit(n int) *SQLJointer {
	this.selector.limit = n
	return this
}

// Offset sets offset of query
func (this *SQLJointer) Offset(n int) *SQLJointer {
	this.selector.offset = n
	return this
}

func (this *SQLJointer) LeftJoin(
	table, onLeft, onRight, onWhere string,
) jointer {
	return func() (*SQLJointer, *sqlJoin) {
		return this, &sqlJoin{
			op: " LEFT JOIN ", table: table, onLeft: onLeft, onRight: onRight,
			where: onWhere, columns: []string{},
		}
	}
}

func (this *SQLJointer) LeftOuterJoin(table, onLeft, onRight string) jointer {
	return func() (*SQLJointer, *sqlJoin) {
		return this, &sqlJoin{
			op: " LEFT OUTER JOIN ", table: table, onLeft: onLeft, onRight: onRight,
			columns: []string{},
		}
	}
}

func (this *SQLJointer) RightJoin(
	table, onLeft, onRight, onWhere string,
) jointer {
	return func() (*SQLJointer, *sqlJoin) {
		return this, &sqlJoin{
			op: " RIGHT JOIN ", table: table, onLeft: onLeft, onRight: onRight,
			where: onWhere, columns: []string{},
		}
	}
}

func (this *SQLJointer) RightOuterJoin(table, onLeft, onRight string) jointer {
	return func() (*SQLJointer, *sqlJoin) {
		return this, &sqlJoin{
			op: " RIGHT OUTER JOIN ", table: table, onLeft: onLeft, onRight: onRight,
			columns: []string{},
		}
	}
}

func (this *SQLJointer) FullJoin(table, onLeft, onRight string) jointer {
	return func() (*SQLJointer, *sqlJoin) {
		return this, &sqlJoin{
			op: " FULL JOIN ", table: table, onLeft: onLeft, onRight: onRight,
			columns: []string{},
		}
	}
}

func (this *SQLJointer) FullOuterJoin(table, onLeft, onRight string) jointer {
	return func() (*SQLJointer, *sqlJoin) {
		return this, &sqlJoin{
			op: " FULL OUTER JOIN ", table: table, onLeft: onLeft, onRight: onRight,
			columns: []string{},
		}
	}
}

func (this *SQLJointer) InnerJoin(table, onLeft, onRight string) jointer {
	return func() (*SQLJointer, *sqlJoin) {
		return this, &sqlJoin{
			op: " INNER JOIN ", table: table, onLeft: onLeft, onRight: onRight,
			columns: []string{},
		}
	}
}

func (this *SQLJointer) buildJoinSQL() (string, *[][]int, int, error) {
	selector := this.selector
	count := 0
	joinSQL := ""
	indexes := make([][]int, 0, len(this.joins))
	leftmost := selector.table.Name
	cols, tableIndexes := selector.buildColumnsSQLRefs()
	if cols != "" {
		indexes = append(indexes, tableIndexes)
		count = len(tableIndexes)
	}

	for _, join := range this.joins {
		if join.onLeft == "" || join.onRight == "" {
			return "", nil, 0, fmt.Errorf("no on columns for join")
		}

		table := selector.tableGetter(join.table)
		if table == nil {
			return "", nil, 0, fmt.Errorf("join table %s is not registered", join.table)
		}

		s := ""
		pos := make([]int, 0, len(join.columns))
		for _, name := range join.columns {
			if c, ok := table.Columns[name]; ok {
				s += join.table + "." + name + ","
				pos = append(pos, c.Index)
			} else {
				return "", nil, 0, fmt.Errorf("%s table has no column %s", join.table, name)
			}
		}

		if s != "" {
			cols += "," + s[:len(s)-1]
			count += len(pos)
			indexes = append(indexes, pos)
		}

		if joinSQL == "" {
			joinSQL = leftmost + join.op + join.table + " ON " + leftmost + "." +
				join.onLeft + "=" + join.table + "." + join.onRight
		} else {
			joinSQL = "(" + joinSQL + ")" + join.op + join.table + " ON " +
				leftmost + "." + join.onLeft + "=" + join.table + "." + join.onRight
		}

		if join.where != "" {
			joinSQL += " " + join.where
		}
	}

	sql := "SELECT " + cols + " FROM " + joinSQL
	if selector.filter.where != "" {
		sql += " WHERE " + selector.filter.where
	}
	if len(selector.sort.columns) > 0 {
		sql += " ORDER BY " + selector.sort.buildSQL(leftmost) + " " +
			selector.sort.op
	}
	if selector.limit > 0 {
		sql += " LIMIT " + strconv.Itoa(selector.limit)
	}
	if selector.offset > 0 {
		sql += " OFFSET " + strconv.Itoa(selector.offset)
	}

	if dbLogger != nil {
		dbLogger(sql)
	}

	return sql, &indexes, count, nil
}

func (this *SQLJointer) Run() (*sql.Rows, error) {
	selector := this.selector
	if selector.err != nil {
		return nil, selector.err
	}

	q, _, _, err := this.buildJoinSQL()
	if err != nil {
		return nil, err
	}

	if selector.tx != nil {
		return selector.tx.Query(q, selector.filter.args...)
	} else {
		return selector.db.Query(q, selector.filter.args...)
	}
}

func (this *SQLJointer) One(rows ...interface{}) error {
	selector := this.selector
	if selector.err != nil {
		return selector.err
	}

	q, indexes, n, err := this.buildJoinSQL()
	if err != nil {
		return err
	}
	if len(rows) != len(*indexes) {
		return fmt.Errorf("not enough rows arguments")
	}

	j := 0
	refs := make([]interface{}, n, n)
	for i, row := range rows {
		rowVal := reflect.ValueOf(row).Elem()
		for _, c := range (*indexes)[i] {
			refs[j] = rowVal.Field(c).Addr().Interface()
			j++
		}
	}

	var rs *sql.Row
	if selector.tx != nil {
		rs = selector.tx.QueryRow(q, selector.filter.args...)
	} else {
		rs = selector.db.QueryRow(q, selector.filter.args...)
	}

	return rs.Scan(refs...)
}

func (this *SQLJointer) All(rows ...interface{}) error {
	selector := this.selector
	if selector.err != nil {
		return selector.err
	}

	size := len(rows)
	rowsVals := make([]reflect.Value, size, size)
	sliceVals := make([]reflect.Value, size, size)
	rowTypes := make([]reflect.Type, size, size)
	for i, r := range rows {
		val := reflect.ValueOf(r)
		if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Slice {
			return fmt.Errorf("rows argument must be a slice address")
		}
		rowsVals[i] = val
		sliceVal := val.Elem()
		sliceVal = sliceVal.Slice(0, sliceVal.Cap())
		sliceVals[i] = sliceVal
		rowTypes[i] = sliceVal.Type().Elem()
	}

	q, indexes, n, err := this.buildJoinSQL()
	if err != nil {
		return err
	}
	if len(rows) != len(*indexes) {
		return fmt.Errorf("not enough rows arguments")
	}

	var rs *sql.Rows
	if selector.tx != nil {
		rs, err = selector.tx.Query(q, selector.filter.args...)
	} else {
		rs, err = selector.db.Query(q, selector.filter.args...)
	}
	if err != nil {
		return err
	}
	defer rs.Close()

	refs := make([]interface{}, n, n)
	rowsPt := make([]reflect.Value, size, size)

	for rs.Next() {
		k := 0
		for i, t := range rowTypes {
			p := reflect.New(t)
			row := p.Elem()
			for _, j := range (*indexes)[i] {
				refs[k] = row.Field(j).Addr().Interface()
				k++
			}
			rowsPt[i] = p
		}

		if err := rs.Scan(refs...); err != nil {
			return err
		}

		for i, _ := range sliceVals {
			sliceVals[i] = reflect.Append(sliceVals[i], rowsPt[i].Elem())
		}
	}

	for i, _ := range rowsVals {
		rowsVals[i].Elem().Set(sliceVals[i]) //.Slice(0, i))
	}
	return nil
}

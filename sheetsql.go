package sheetsql

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"io"
	"strings"

	"github.com/myoan/sheetsql/gss"
)

func init() {
	sql.Register("sheetsql", &SheetDriver{})
}

type SheetDriver struct{}

func (d *SheetDriver) Open(dsn string) (driver.Conn, error) {
	if len(dsn) <= 1 {
		return nil, errors.New("invalid dsn")
	}
	parts := strings.Split(dsn, "|")
	client, _ := gss.NewSpreadSheet(parts[0], parts[1])
	return &SheetConn{client: client}, nil
}

type SheetConn struct {
	client *gss.Sheet
}

func (c *SheetConn) Close() error {
	c.client = nil
	return nil
}

func (c *SheetConn) Prepare(query string) (driver.Stmt, error) {
	return &SheetStmt{c, query}, nil
}

func (c *SheetConn) Begin() (driver.Tx, error) {
	return nil, nil
}

type SheetStmt struct {
	conn  *SheetConn
	query string
}

func (s *SheetStmt) Close() error {
	return nil
}

func (s *SheetStmt) NumInput() int {
	return strings.Count(s.query, "?")
}

func (s *SheetStmt) Query(args []driver.Value) (driver.Rows, error) {
	tbl, _ := GetTableName(s.query)
	// exec query
	return s.GetSheetData(tbl), nil
}

func (s *SheetStmt) GetSheetData(table string) *SheetRows {
	columns := s.conn.client.GetSheetColumn(table)
	records := s.conn.client.GetSheetRecord(table, len(columns))
	return &SheetRows{s, columns, 0, records}
}

func (s *SheetStmt) Exec(args []driver.Value) (driver.Result, error) {
	return nil, errors.New("Exec does not supported")
}

type SheetRows struct {
	s       *SheetStmt
	columns []string
	index   int
	data    [][]interface{}
}

func (row *SheetRows) Close() error {
	return nil
}

func (row *SheetRows) Columns() []string {
	return row.columns
}

func (row *SheetRows) Next(dest []driver.Value) error {
	if row.index == len(row.data) {
		return io.EOF
	}
	for i, val := range row.data[row.index] {
		dest[i] = val
	}
	row.index++
	return nil
}

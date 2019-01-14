package sheetsql

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/jwt"
	sheets "google.golang.org/api/sheets/v4"
)

func getConfig(jwtJSON []byte) (*jwt.Config, error) {
	cfg, err := google.JWTConfigFromJSON(
		jwtJSON,
		"https://www.googleapis.com/auth/spreadsheets",
	)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func init() {
	sql.Register("sheetsql", &SheetDriver{})
}

type SheetDriver struct{}

func (d *SheetDriver) Open(dsn string) (driver.Conn, error) {
	if len(dsn) <= 1 {
		return nil, errors.New("invalid dsn")
	}
	parts := strings.Split(dsn, "|")

	cred, err := os.Open(parts[0])
	if err != nil {
		panic(err)
	}

	jwtJSON, err := ioutil.ReadAll(cred)
	if err != nil {
		panic(err)
	}

	config, err := getConfig(jwtJSON)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := config.Client(context.Background())
	return &SheetConn{client: client, sheetID: parts[1]}, nil
}

type SheetConn struct {
	client  *http.Client
	sheetID string
	key     string
	secret  string
	env     string
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
	srv, err := sheets.New(s.conn.client)
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
	}

	readRange := "Sheet1!A1:C"
	resp, err := srv.Spreadsheets.Values.Get(s.conn.sheetID, readRange).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve data from sheet: %v", err)
	}

	var columns []string
	for _, elem := range resp.Values[0] {
		columns = append(columns, elem.(string))
	}
	return &SheetRows{s, columns, 0, resp.Values[1:]}, nil
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

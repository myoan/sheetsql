package sheetsql

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/sheets/v4"
)

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func init() {
	sql.Register("sheetsql", &SheetDriver{})
}

type SheetDriver struct{}

func (d *SheetDriver) Open(dsn string) (driver.Conn, error) {
	fmt.Println(dsn)
	return &SheetConn{c: http.DefaultClient}, nil
}

type SheetConn struct {
	c      *http.Client
	key    string
	secret string
	env    string
}

func (c *SheetConn) Close() error {
	c.c = nil
	return nil
}

func (c *SheetConn) Prepare(query string) (driver.Stmt, error) {
	return &SheetStmt{c, query}, nil
}

func (c *SheetConn) Begin() (driver.Tx, error) {
	return nil, nil
}

type SheetStmt struct {
	c     *SheetConn
	query string
}

func (s *SheetStmt) Close() error {
	return nil
}

func (s *SheetStmt) NumInput() int {
	return 0
}

func (s *SheetStmt) Query(args []driver.Value) (driver.Rows, error) {
	b, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/spreadsheets.readonly")
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)

	srv, err := sheets.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
	}

	spreadsheetId := "xxxxxxxxxxxxxxxxxxxxxx"
	readRange := "Sheet1!A1:C"
	resp, err := srv.Spreadsheets.Values.Get(spreadsheetId, readRange).Do()
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

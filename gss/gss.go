package gss

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/jwt"
	sheets "google.golang.org/api/sheets/v4"
)

type Sheet struct {
	srv     *sheets.Service
	SheetID string
}

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

func NewSpreadSheet(file, sheetID string) (*Sheet, error) {
	cred, err := os.Open(file)
	if err != nil {
		return nil, err
	}

	jwtJSON, err := ioutil.ReadAll(cred)
	if err != nil {
		return nil, err
	}

	config, err := getConfig(jwtJSON)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
		return nil, err
	}
	client := config.Client(context.Background())

	srv, err := sheets.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
		return nil, err
	}

	return &Sheet{srv: srv, SheetID: sheetID}, nil
}

func (s *Sheet) GetSheetColumn(tbl string) []string {
	resp, err := s.srv.Spreadsheets.Values.Get(s.SheetID, columnRange(tbl)).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve data from sheet: %v", err)
	}

	var columns []string
	for _, elem := range resp.Values[0] {
		columns = append(columns, elem.(string))
	}
	return columns
}

func columnRange(table string) string {
	return table + "!A1:1"
}

func (s *Sheet) GetSheetRecord(tbl string, columns []string) [][]interface{} {
	resp, err := s.srv.Spreadsheets.Values.Get(s.SheetID, s.tableTarget(tbl, columns)).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve data from sheet: %v", err)
	}

	return resp.Values[1:]
}

func (s *Sheet) tableTarget(tbl string, columns []string) string {
	allCol := s.GetSheetColumn(tbl)
	min, max := columnToRange(allCol, columns)
	minAlp, _ := columnAlphabet(min)
	maxAlp, _ := columnAlphabet(max)
	return fmt.Sprintf("%s!%s1:%s", tbl, minAlp, maxAlp)
}

func columnToRange(col1 []string, col2 []string) (int, int) {
	if len(col2) == 1 {
		return 1, len(col1)
	}
	min := len(col1)
	max := 0
	for _, c2 := range col2 {
		for i, c1 := range col1 {
			if c1 == c2 {
				if min > i {
					min = i
				}
				if i > max {
					max = i
				}
			}
		}
	}
	return min + 1, max + 1
}

func columnAlphabet(len int) (string, error) {
	if len <= 0 {
		return "", errors.New("length OutOfRange")
	}
	val := len - 1
	div := val / 26
	mod := val % 26
	ret := string(65 + mod)
	val = div
	for val > 0 {
		div := val / 26
		mod := val % 26
		ret = string(64+mod) + ret
		val = div
	}
	return ret, nil
}

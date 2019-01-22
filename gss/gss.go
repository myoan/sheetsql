package gss

import (
	"context"
	"io/ioutil"
	"log"
	"os"

	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/jwt"
	sheets "google.golang.org/api/sheets/v4"
)

type Sheet struct {
	srv  *sheets.Service
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
	}

	return &Sheet{srv: srv, SheetID: sheetID}, nil
}

func (s *Sheet) GetSheetColumn(rng string) []string {
	resp, err := s.srv.Spreadsheets.Values.Get(s.SheetID, columnRange(rng)).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve data from sheet: %v", err)
	}

	var columns []string
	for _, elem := range resp.Values[0] {
		columns = append(columns, elem.(string))
	}
	return columns
}

func (s *Sheet) GetSheetRecord(rng string, len int) [][]interface{} {
	resp, err := s.srv.Spreadsheets.Values.Get(s.SheetID, rng+"!A1:"+columnAlphabet(len)).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve data from sheet: %v", err)
	}

	return resp.Values[1:]
}

func columnRange(table string) string {
	return table + "!A1:1"
}

func columnAlphabet(len int) string {
	return string(64 + len)
}

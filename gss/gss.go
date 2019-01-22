package gss

import (
	"context"
	"io/ioutil"
	"log"
	"net/http"
	"os"

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

func NewSpreadSheet(file string) (*http.Client, error) {
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
	return config.Client(context.Background()), nil
}

func GetSheetColumn(client *http.Client, sheetID string, rng string) []string {
	srv, err := sheets.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
	}

	resp, err := srv.Spreadsheets.Values.Get(sheetID, columnRange(rng)).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve data from sheet: %v", err)
	}

	var columns []string
	for _, elem := range resp.Values[0] {
		columns = append(columns, elem.(string))
	}
	return columns
}

func GetSheetRecord(client *http.Client, sheetID string, rng string, len int) [][]interface{} {
	srv, err := sheets.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
	}

	resp, err := srv.Spreadsheets.Values.Get(sheetID, rng+"!A1:" + columnAlphabet(len)).Do()
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

package sheetsql_test

import (
	"testing"
	sql "github.com/myoan/sheetsql"
)

func TestSimpleParse(t *testing.T) {
	actual := sql.ParseQuery("SELECT * FROM table;")
	expected := true
	if actual != expected {
		t.Errorf("got: %v\nwant: %v", actual, expected)
	}
}

package sheetsql_test

import (
	"testing"

	"github.com/myoan/sheetsql"
)

func TestGetTableName(t *testing.T) {
	validSQL := []struct {
		input  string
		output string
	}{
		{
			input: "select a, b from tbl",
		},
		{
			input: "select * from tbl",
		},
		{
			input: "insert into tbl values (1)",
		},
		{
			input: "update tbl set `id` = 1",
		},
		{
			input: "delete from tbl",
		},
	}

	for _, tcase := range validSQL {
		if tcase.output == "" {
			tcase.output = "tbl"
		}
		out, err := sheetsql.GetTableName(tcase.input)
		if err != nil {
			t.Errorf("Parse(%q) err: %v, want nil", tcase.input, err)
			continue
		}
		if out != tcase.output {
			t.Errorf("Parse(%q) = %q, want: %q", tcase.input, out, tcase.output)
		}
	}
}

func TestGetColumnFromSelect(t *testing.T) {
	validSQL := []struct {
		input  string
		output string
	}{
		{
			input:  "select * from tbl",
			output: "*",
		},
		{
			input:  "select foo, bar from tbl",
			output: "foobar",
		},
	}

	for _, tcase := range validSQL {
		out, err := sheetsql.GetColumns(tcase.input)
		if err != nil {
			t.Errorf("Parse(%q) err: %v, want nil", tcase.input, err)
			continue
		}
		var act = ""
		for _, col := range out {
			act += col.Name
		}
		if act != tcase.output {
			t.Errorf("Parse(%q) = %q, want: %q", tcase.input, out, tcase.output)
		}
	}
}

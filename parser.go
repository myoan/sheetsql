package sheetsql

import (
	"fmt"

	"github.com/xwb1989/sqlparser"
)

type Cell struct {
	column string
	raw    string
}

func (c *Cell) Format() string {
	return c.column + c.raw
}

type SheetRange struct {
	Name string
	Sp   Cell
	Ep   Cell
}

func (sr *SheetRange) ToRange() string {
	return sr.Name + "!" + sr.Sp.Format() + ":" + sr.Ep.Format()
}

func GetTableName(query string) (string, error) {
	stmt, err := sqlparser.Parse(query)
	if err != nil {
		return "", err
	}

	switch stmt := stmt.(type) {
	case *sqlparser.Select:
		result := sqlparser.String(stmt.From)
		return result, nil
	case *sqlparser.Insert:
		result := sqlparser.String(stmt.Table)
		return result, nil
	default:
		return "", nil
	}
}

func ParseQuery(query string) string {
	stmt, err := sqlparser.Parse(query)
	if err != nil {
		// Do something with the err
	}

	fmt.Println(sqlparser.String(stmt))
	// Otherwise do something with stmt
	switch stmt := stmt.(type) {
	case *sqlparser.Select:
		result := sqlparser.String(stmt.From)
		fmt.Println(result)
		return result
	case *sqlparser.Insert:
		return ""
	default:
		return ""
	}
}

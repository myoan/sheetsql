package sheetsql

import (
	"strings"

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

type Table struct {
	Name    string
	Columns *Columns
}

func (tbl *Table) ColumnNames() []string {
	result := make([]string, len(*tbl.Columns))
	for i, col := range *tbl.Columns {
		result[i] = col.Name
	}
	return result
}

type Columns []Column

type Column struct {
	Name string
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
	case *sqlparser.Update:
		result := sqlparser.String(stmt.TableExprs)
		return result, nil
	case *sqlparser.Delete:
		result := sqlparser.String(stmt.TableExprs)
		return result, nil
	default:
		return "", nil
	}
}

func GetColumns(query string) (Columns, error) {
	stmt, err := sqlparser.Parse(query)
	if err != nil {
		return nil, err
	}

	buf := sqlparser.NewTrackedBuffer(nil)
	switch stmt := stmt.(type) {
	case *sqlparser.Select:
		stmt.SelectExprs.Format(buf)
		tmp := strings.Split(buf.String(), ",")
		columns := make(Columns, len(tmp))
		for i, colName := range tmp {
			columns[i] = Column{Name: strings.TrimSpace(colName)}
		}
		return columns, nil
	default:
		return nil, nil
	}
}

func ParseQuery(query string) (*Table, error) {
	tblName, _ := GetTableName(query)
	columns, _ := GetColumns(query)

	return &Table{Name: tblName, Columns: &columns}, nil
}

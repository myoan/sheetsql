package sheetsql

import (
	"database/sql/driver"
)

type Evaluator struct {
}

func Eval(s *SheetStmt, tbl *Table) (driver.Rows, error) {
	columns := s.conn.client.GetSheetColumn(tbl.Name)
	records := s.conn.client.GetSheetRecord(tbl.Name, tbl.ColumnNames())
	return &SheetRows{s, columns, 0, records}, nil
}

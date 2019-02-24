package sheetsql

import (
	"database/sql/driver"
)

type Evaluator struct {
}

func Eval(s *SheetStmt, tbl *Table) (driver.Rows, error) {
	records := s.conn.client.GetSheetRecord(tbl.Name, len(*tbl.Columns))
	return &SheetRows{s, tbl.ColumnNames(), 0, records}, nil
}

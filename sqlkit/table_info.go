package sqlkit

import (
	"database/sql"
	"errors"
	"fmt"
)

type TableInfo struct {
	DB          *DB
	Log         FieldLogger
	Table       string
	ErrorPrefix string
	Columns     []string
	Rows        *sql.Rows
	Error       error
}

func (ti *TableInfo) SetupAndVerify(db *DB) error {
	ti.DB = db
	if ti.Table == "" {
		return errors.New("missing table name")
	}
	if ti.ErrorPrefix == "" {
		return errors.New("missing error prefix")
	}
	if ti.Columns == nil || len(ti.Columns) < 1 {
		return errors.New("missing columns list")
	}
	for _, c := range ti.Columns {
		if c == "" {
			return errors.New("invalid column name")
		}
	}
	// TODO: check that the table exists as expected on the DB
	return nil
}

func (ti *TableInfo) errorCode(etype string, code string, seq string) string {
	return fmt.Sprintf("TI:%s_%s_%s_%s", ti.ErrorPrefix, etype, code, seq)
}

func (ti *TableInfo) ExecuteStatement(stmt string, stmtType string, args ...interface{}) (resp Response) {
	if ti.Log != nil {
		ti.Log.Infof("EXECUTE_STATEMENT: %s", stmt)
	}
	_, err := ti.DB.Exec(stmt, args...)
	resp.CheckError(ti.Log, err, ti.errorCode(stmtType, "EX", "00"))
	return
}

func (ti *TableInfo) Select(stmt string, args ...interface{}) (resp Response) {
	if ti.Log != nil {
		ti.Log.Infof("SELECT: %s || WITH ARGS %v", stmt, args)
	}
	ti.Rows, ti.Error = ti.DB.Query(stmt, args...)
	resp.CheckError(ti.Log, ti.Error, ti.errorCode("QU", "SL", "00"))
	return
}

func (ti *TableInfo) FetchOnce(args ...interface{}) (resp Response) {
	if ti.Log != nil {
		ti.Log.Info("FETCH_ONCE")
	}
	resp = ti.FetchNext(args...)
	if resp.Success {
		ti.Rows.Close()
	}
	return
}

func (ti *TableInfo) FetchNext(args ...interface{}) (resp Response) {
	if ti.Log != nil {
		ti.Log.Info("FETCH_NEXT")
	}
	if ti.Rows.Next() {
		ti.Error = ti.Rows.Scan(args...)
		resp.CheckError(ti.Log, ti.Error, ti.errorCode("QU", "FN", "00"))
		if !resp.Success {
			ti.Rows.Close()
		}
	} else {
		ti.Rows.Close()
		ti.Error = ti.Rows.Err()
		resp.CheckError(ti.Log, ti.Error, ti.errorCode("QU", "FN", "01"))
		resp = Response{false, "QUFNEOC", "No more records"}
	}
	return
}

func (ti *TableInfo) InsertStatement() (stmt string) {
	nCols := 0
	stmt = "INSERT INTO " + ti.Table + " ("
	stmt, nCols = addColumns(stmt, nCols, ti.Columns, false)
	stmt, _ = addColumns(stmt, nCols, TimeAuditCols(), false)
	stmt += ") VALUES ("
	stmt += addInsertValues(nCols)
	stmt += InsertTimeAuditValues()
	stmt += ")"
	return
}

func (ti *TableInfo) UpdateStatement() (stmt string) {
	nCols := 0
	stmt = "UPDATE " + ti.Table + " SET "
	stmt, nCols = addColumns(stmt, nCols, ti.Columns, true)
	stmt += UpdateTimeAuditCols()
	return
}

func (ti *TableInfo) SelectStatementByUUID() (stmt string) {
	nCols := 0
	stmt = "SELECT "
	stmt, nCols = addColumns(stmt, nCols, ti.Columns, false)
	stmt += " FROM " + ti.Table + " WHERE uuid = ?"
	return
}

func (ti *TableInfo) SelectStatementByClause(whereClause string) (stmt string) {
	nCols := 0
	stmt = "SELECT "
	stmt, nCols = addColumns(stmt, nCols, ti.Columns, false)
	stmt += " FROM " + ti.Table + " WHERE " + whereClause
	return
}

func (ti *TableInfo) DeleteStatementByUUID() (stmt string) {
	stmt = "DELETE FROM " + ti.Table + " WHERE uuid = ?"
	return
}
func (ti *TableInfo) DeleteStatementByClause(whereClause string) (stmt string) {
	stmt = "DELETE FROM " + ti.Table + " WHERE " + whereClause
	return
}

func addColumns(curStmt string, curCols int, cols []string, updateStmt bool) (stmt string, ncols int) {
	stmt = curStmt
	ncols = curCols

	updStr := ""
	if updateStmt {
		updStr = " = ?"
	}

	for _, col := range cols {
		if ncols > 0 {
			stmt += ", "
		}
		ncols++
		stmt += col + updStr
	}
	return
}

func addInsertValues(nCols int) (stmt string) {
	for i := 0; i < nCols; i++ {
		if i > 0 {
			stmt += ", "
		}
		stmt += "?"
	}
	return
}

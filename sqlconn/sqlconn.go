package sqlconn

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
)

// SQLConn . . .
type SQLConn struct {
	db          *sql.DB
	asyncErrors chan error
}

// SQLQuery . . .
type SQLQuery struct {
	query string
	args  []interface{}
}

func getConnFromConfig() string {
	var env *os.File
	var err error
	if env, err = os.Open("config.json"); err != nil {
		panic(err)
	}

	var sqlConn struct {
		DBConn string `json:"dbConn"`
	}
	json.NewDecoder(env).Decode(&sqlConn)
	return sqlConn.DBConn
}

// Open . . .
func Open(DBAddress ...string) *SQLConn {
	DBAddress = append(DBAddress, os.Getenv("DB_CONN"))
	dbAddress := DBAddress[0]
	if dbAddress == "" {
		dbAddress = getConnFromConfig()
	}

	var sqlConn SQLConn
	var err error
	if sqlConn.db, err = sql.Open("mssql", dbAddress); sqlConn.db == nil || err != nil {
		panic(fmt.Errorf("failed to open database connection: %s", err))
	}
	return &sqlConn
}

// Close . . .
func (sqlConn SQLConn) Close() {
	sqlConn.db.Close()
}

// NewSQLQuery . . .
func NewSQLQuery(query string, args ...interface{}) SQLQuery {
	return SQLQuery{query, args}
}

// QueryRow . . .
func (sqlConn *SQLConn) QueryRow(query string, args ...interface{}) *sql.Row {
	return sqlConn.db.QueryRow(query, args...)
}

// QueryRowAsync . . .
func (sqlConn *SQLConn) QueryRowAsync(queries ...SQLQuery) (rows []*sql.Row) {
	sqlConn.asyncErrors = make(chan error)
	queryRowChs := []<-chan *sql.Row{}
	for _, query := range queries {
		queryRowChs = append(queryRowChs, sqlConn.queryRowAsync(query))
	}
	if err := sqlConn.asyncError(); err != nil {
		panic(err)
	}
	for _, queryRowCh := range queryRowChs {
		row := <-queryRowCh
		rows = append(rows, row)
	}
	return rows
}

func (sqlConn *SQLConn) queryRowAsync(query SQLQuery) <-chan *sql.Row {
	r := make(chan *sql.Row)
	go func() {
		defer close(r)
		defer func() {
			if err := recover(); err != nil && sqlConn.asyncErrors != nil {
				sqlConn.asyncErrors <- err.(error)
			}
		}()
		row := sqlConn.QueryRow(query.query, query.args...)
		if sqlConn.asyncErrors != nil {
			sqlConn.asyncErrors <- nil
		}
		r <- row
	}()
	return r
}

// Query . . .
func (sqlConn *SQLConn) Query(query string, args ...interface{}) (rows *sql.Rows) {
	var err error
	if rows, err = sqlConn.db.Query(query, args...); err != nil {
		panic(err)
	}
	return rows
}

// QueryAsync . . .
func (sqlConn *SQLConn) QueryAsync(queries ...SQLQuery) (rowsSet []*sql.Rows) {
	sqlConn.asyncErrors = make(chan error)
	queryRowChs := []<-chan *sql.Rows{}
	for _, query := range queries {
		queryRowChs = append(queryRowChs, sqlConn.queryAsync(query))
	}
	if err := sqlConn.asyncError(); err != nil {
		panic(err)
	}
	for _, queryRowCh := range queryRowChs {
		row := <-queryRowCh
		rowsSet = append(rowsSet, row)
	}
	return rowsSet
}

func (sqlConn *SQLConn) queryAsync(query SQLQuery) <-chan *sql.Rows {
	r := make(chan *sql.Rows)
	go func() {
		defer close(r)
		defer func() {
			if err := recover(); err != nil && sqlConn.asyncErrors != nil {
				sqlConn.asyncErrors <- err.(error)
			}
		}()
		rows := sqlConn.Query(query.query, query.args...)
		if sqlConn.asyncErrors != nil {
			sqlConn.asyncErrors <- nil
		}
		r <- rows
	}()
	return r
}

// Exec . . .
func (sqlConn *SQLConn) Exec(query string, args ...interface{}) (rowsAffected int64) {
	var result sql.Result
	var err error
	if result, err = sqlConn.db.Exec(query, args...); err != nil {
		panic(err)
	}
	if rowsAffected, err = result.RowsAffected(); err != nil {
		panic(err)
	}
	return rowsAffected
}

// ExecAsync . . .
func (sqlConn *SQLConn) ExecAsync(sqlQueries ...SQLQuery) (rowsAffected []int64) {
	sqlConn.asyncErrors = make(chan error)
	execChs := []<-chan int64{}
	for _, query := range sqlQueries {
		execChs = append(execChs, sqlConn.execAsync(query))
	}
	if err := sqlConn.asyncError(); err != nil {
		panic(err)
	}
	for _, queryRowCh := range execChs {
		execRowsAffected := <-queryRowCh
		rowsAffected = append(rowsAffected, execRowsAffected)
	}
	return rowsAffected
}

func (sqlConn *SQLConn) execAsync(sqlQuery SQLQuery) <-chan int64 {
	r := make(chan int64)
	go func() {
		defer close(r)
		defer func() {
			if err := recover(); err != nil && sqlConn.asyncErrors != nil {
				sqlConn.asyncErrors <- err.(error)
			}
		}()
		rowsAffected := sqlConn.Exec(sqlQuery.query, sqlQuery.args...)
		if sqlConn.asyncErrors != nil {
			sqlConn.asyncErrors <- nil
		}
		r <- rowsAffected
	}()
	return r
}

func (sqlConn *SQLConn) asyncError() error {
	err := <-sqlConn.asyncErrors
	sqlConn.asyncErrors = nil
	return err
}

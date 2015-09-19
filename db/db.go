package db

import (
	"sync"

	"github.com/jmoiron/sqlx"
)

var (
	db *sqlx.DB

	stmtSetRWMutex sync.RWMutex
	stmtSet        = make(map[string]*sqlx.Stmt) // map[query]*sqlx.Stmt

	namedStmtSetRWMutex sync.RWMutex
	namedStmtSet        = make(map[string]*sqlx.NamedStmt) // map[query]*sqlx.NamedStmt
)

func GetDB() *sqlx.DB {
	return db
}

func GetStmt(query string) (stmt *sqlx.Stmt, err error) {
	stmtSetRWMutex.RLock()
	stmt = stmtSet[query]
	stmtSetRWMutex.RUnlock()

	if stmt != nil {
		return
	}

	stmtSetRWMutex.Lock()
	defer stmtSetRWMutex.Unlock()

	if stmt = stmtSet[query]; stmt != nil {
		return
	}

	stmt, err = db.Preparex(query)
	if err != nil {
		return
	}
	stmtSet[query] = stmt
	return
}

func GetNamedStmt(query string) (stmt *sqlx.NamedStmt, err error) {
	namedStmtSetRWMutex.RLock()
	stmt = namedStmtSet[query]
	namedStmtSetRWMutex.RUnlock()

	if stmt != nil {
		return
	}

	namedStmtSetRWMutex.Lock()
	defer namedStmtSetRWMutex.Unlock()

	if stmt = namedStmtSet[query]; stmt != nil {
		return
	}

	stmt, err = db.PrepareNamed(query)
	if err != nil {
		return
	}
	namedStmtSet[query] = stmt
	return
}

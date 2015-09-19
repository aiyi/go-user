// 读取配置, 创建全局的 sqlx.DB 对象, 根据 query 创建 Stmt 对象并缓存该对象.
//
// Stmt 对象维持了"没有关闭的数据库连接"的 driver.Stmt 对象, 这个是轻资源, 并不会维持连接不关闭;
// 如果这些轻资源会有很大的负面影响, 可以考虑 LRU, 不过目前我觉得还没有必要.
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

// CloseDB 关闭数据库连接, 释放资源.
// 一般情况下没有必要调用该函数.
func CloseDB() error {
	return db.Close()
}

// CloseAllStmt 关闭所有缓存的 Stmt, 释放资源.
// 一般情况下没有必要调用该函数.
func CloseAllStmt() {
	func() {
		stmtSetRWMutex.Lock()
		defer stmtSetRWMutex.Unlock()
		for _, stmt := range stmtSet {
			stmt.Close()
		}
		stmtSet = make(map[string]*sqlx.Stmt)
	}()

	func() {
		namedStmtSetRWMutex.Lock()
		defer namedStmtSetRWMutex.Unlock()
		for _, stmt := range namedStmtSet {
			stmt.Close()
		}
		namedStmtSet = make(map[string]*sqlx.NamedStmt)
	}()
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

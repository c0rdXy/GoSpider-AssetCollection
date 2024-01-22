package mysql

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

var db *sql.DB

// InitDB 初始化数据库
func InitDB() {
	var err error
	db, err = sql.Open("mysql", "root:root@tcp(localhost:3306)/spider")
	if err != nil {
		log.Fatal(err)
	}
}

// CloseDB 关闭数据库连接
func CloseDB() {
	db.Close()
}

// InsertSource 插入公司名称
func InsertSource(company string) (int64, error) {
	result, err := db.Exec("INSERT INTO source (company) VALUES (?)", company)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

// InsertInvest 插入投资信息
func InsertInvest(sourceID int64, product string) error {
	_, err := db.Exec("INSERT INTO invest (source_id, product) VALUES (?, ?)", sourceID, product)
	return err
}

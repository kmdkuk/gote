package db

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"go.uber.org/zap"
)

var DB *sql.DB

// DB の接続情報
const (
	DRIVER_NAME = "mysql" // ドライバ名(mysql固定)
	// user:password@tcp(container-name:port)/dbname ※mysql はデフォルトで用意されているDB
	DATA_SOURCE_NAME = "docker:docker@tcp(db:3306)/network-monitoring?parseTime=true"
)

type DisconnectDate struct {
	ID      int
	StartAt time.Time
	EndAt   time.Time
}

// TODO: disconnected data store function
func RecordStartAt() {
	logger := zap.L()
	// insert
	ins, err := DB.Prepare("INSERT INTO `network-monitoring`.disconnect_dates(start_at) VALUES(?)")
	if err != nil {
		logger.Error("prepare error occurred", zap.Error(err))
	}
	now := time.Now()
	logger.Info(fmt.Sprintln(now))
	_, err = ins.Exec(now)
	if err != nil {
		logger.Error("insert error occurred", zap.Error(err))
	}
}

func RecordEndAt() {
	logger := zap.L()
	query := "select * from disconnect_dates where end_at=0;"
	rows, err := DB.Query(query)
	if err != nil {
		logger.Error("error occurred", zap.Error(err))
		return
	}
	record := DisconnectDate{}
	rows.Next()
	if err := rows.Scan(&record.ID, &record.StartAt, &record.EndAt); err != nil {
		logger.Error("error occurred", zap.Error(err))
		return
	}
	// set end_at
	upd, err := DB.Prepare("UPDATE disconnect_dates set end_at = ? where id = ?")
	if err != nil {
		logger.Error("update prepare error occurred", zap.Error(err))
		return
	}
	upd.Exec(time.Now(), record.ID)
}

func init() {
	logger := zap.L()
	var err error
	DB, err = sql.Open(DRIVER_NAME, DATA_SOURCE_NAME)
	if err != nil {
		logger.Fatal("error connecting to database", zap.Error(err))
	}
}

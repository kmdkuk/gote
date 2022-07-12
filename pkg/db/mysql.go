package db

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
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
	// insert
	ins, err := DB.Prepare("INSERT INTO `network-monitoring`.disconnect_dates(start_at) VALUES(?)")
	if err != nil {
		log.Printf("[Err] Prepare %v", err)
	}
	log.Println(time.Now())
	_, err = ins.Exec(time.Now())
	if err != nil {
		log.Printf("[Err] Insert %v", err)
	}
}

func RecordEndAt() {
	query := "select * from disconnect_dates where end_at=0;"
	rows, err := DB.Query(query)
	if err != nil {
		log.Printf("[Err] %s %v", query, err)
		return
	}
	record := DisconnectDate{}
	rows.Next()
	if err := rows.Scan(&record.ID, &record.StartAt, &record.EndAt); err != nil {
		log.Printf("[Err] %v", err)
		return
	}
	// set end_at
	upd, err := DB.Prepare("UPDATE disconnect_dates set end_at = ? where id = ?")
	if err != nil {
		log.Printf("[Err] Update Prepare: %v", err)
		return
	}
	upd.Exec(time.Now(), record.ID)
}

func init() {
	var err error
	DB, err = sql.Open(DRIVER_NAME, DATA_SOURCE_NAME)
	if err != nil {
		log.Fatal("error connecting to database: ", err)
	}
}

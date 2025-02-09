package postgresSql

import (
	"backend/internal/http-server/handlers/postPingInfo"
	"database/sql"
	"fmt"
	"log"
	"sync"
	"time"
)

type Postgres struct {
	db *sql.DB
}

var (
	pgInstance *Postgres
	pgOnce     sync.Once
)

type Data struct {
	IpAddress string
	PingTime  time.Time
	Status    string
}

func NewPg(conn string) (*Postgres, error) {
	pgOnce.Do(func() {
		db, err := sql.Open("postgres", conn)
		if err != nil {
			fmt.Errorf("unable to create conection pool: %w", err)
		}
		pgInstance = &Postgres{db: db}
	})
	if pgInstance == nil {
		return nil, fmt.Errorf("failed to initialize Postgres instance")
	}
	return pgInstance, nil
}

func (pg *Postgres) GetListIp() ([]Data, error) {
	const funcName = "GetListIp"
	str := "SELECT * FROM ip_addresses"
	rows, err := pg.db.Query(str)
	defer rows.Close()
	if err != nil {
		return []Data{}, fmt.Errorf("%s: %w", funcName, err)
	}

	var data []Data
	for rows.Next() {
		var row Data
		if err = rows.Scan(&row.IpAddress, &row.PingTime, &row.Status); err != nil {
			log.Println(err)
		}
		data = append(data, row)
	}
	return data, nil
}

func (pg *Postgres) newEntry(data Data) error {
	const funcName = "newEnty"
	dateValue := data.PingTime.Format("2006-01-02 15:04:05")
	str := "INSERT INTO ip_addresses(ip_address, ping_time, status) VALUES ($1, $2, $3)"
	_, err := pg.db.Exec(str, data.IpAddress, dateValue, data.Status)
	if err != nil {
		return fmt.Errorf("%s: %w", funcName, err)
	}
	return nil
}

func (pg *Postgres) SaveNewInfo(pingInfo []postPingInfo.PingInfo) error {
	const funcName = "postgres/SaveNewInfo"
	for i := range pingInfo {
		packetLoss := pingInfo[i].PacketLoss
		var data Data
		if packetLoss == 0.0 {
			data = Data{
				IpAddress: pingInfo[i].IPAddr,
				PingTime:  pingInfo[i].PingTime,
				Status:    "success",
			}
		} else {
			data = Data{
				IpAddress: pingInfo[i].IPAddr,
				PingTime:  pingInfo[i].PingTime,
				Status:    "failure",
			}
		}
		fmt.Println(data)
		err := pg.newEntry(data)
		if err != nil {
			return fmt.Errorf("func:%s ERROR save entry %v: %w", funcName, data, err)
		}
	}
	return nil
}

package storage

import (
	"fmt"
	"os"

	"github.com/ziutek/mymysql/mysql"
	_ "github.com/ziutek/mymysql/native"
)

type Repository interface {
	Read() []*QueueFromDB
}

const query = "SELECT * FROM proxy;"

// InitDBReponsitory should valueate connection
func InitDBReponsitory() Repository {
	var repo Repository
	repo = &QueueFromDB{}
	return repo

}

func (q *QueueFromDB) Read() []*QueueFromDB {
	conn := mysql.New(
		"tcp",
		"",
		os.Getenv("DB_HOST")+":"+os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_NAME"),
	)
	var final []*QueueFromDB
	conn.Connect()
	rows, _, _ := conn.Query(query)
	defer conn.Close()
	for row := range rows {
		fmt.Println(row)
	}

	return final
}

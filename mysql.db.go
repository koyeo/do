package do

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"log"
)

func NewMySqlDatabase(conn *Conn) *MySqlDatabase {

	db, err := xorm.NewEngine("mysql",
		fmt.Sprintf(
			"%s:%s@(%s:%d)/%s?charset=utf8",
			conn.User,
			conn.Password,
			conn.Host,
			conn.Port,
			conn.Database,
		))
	if err != nil {
		log.Fatal("Init mysql error", err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal("Ping mysql error: ", err)
	}

	return &MySqlDatabase{
		name: conn.Database,
		db:   db,
	}
}

type MySqlDatabase struct {
	name string
	db   *xorm.Engine
}

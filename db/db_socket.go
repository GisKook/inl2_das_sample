package db

import (
	"database/sql"
	"fmt"
	"github.com/giskook/inl2_das_sample/base"
	"github.com/giskook/inl2_das_sample/conf"
	_ "github.com/lib/pq"
	"log"
	"sync"
)

type DbSocket struct {
	Db                  *sql.DB
	ChargingPricesMutex sync.Mutex
	TransactionChan     chan *base.BluetoothRing
}

var G_DBSocket_Mutex sync.Mutex
var G_DBSocket *DbSocket

func NewDbSocket(db_config *conf.DBConfiguration) (*DbSocket, error) {
	defer G_DBSocket_Mutex.Unlock()
	G_DBSocket_Mutex.Lock()
	conn_string := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable", db_config.User, db_config.Passwd, db_config.Host, db_config.Port, db_config.DbName)

	log.Println(conn_string)
	db, err := sql.Open("postgres", conn_string)

	if err != nil {
		log.Println(err)
		return nil, err
	}
	log.Println("db open success")

	G_DBSocket = &DbSocket{
		Db:              db,
		TransactionChan: make(chan *base.BluetoothRing, 1024),
	}

	return G_DBSocket, nil
}

func (db_socket *DbSocket) Close() {
	db_socket.Db.Close()
}

func GetDBSocket() *DbSocket {
	defer G_DBSocket_Mutex.Unlock()
	G_DBSocket_Mutex.Lock()

	return G_DBSocket
}

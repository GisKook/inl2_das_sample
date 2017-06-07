package db

import (
	"fmt"
	//"github.com/golang/protobuf/proto"
	"github.com/giskook/inl2_das_sample/base"
	"log"
	"strings"
)

const (
	SQL_INSERT_TABLE string = "INSERT INTO t_indoor_intensity (imei, bt_id, rssi, orientation, loc, finger_id) values ('%s', %d,%d,'%s', '%s', '%s')"
)

func (db_socket *DbSocket) ProccessTransaction() {
	for {
		select {
		case trans := <-db_socket.TransactionChan:

			tx, err := db_socket.Db.Begin()
			if err != nil {
				log.Println("ProccessTransaction")
				log.Println(err)
			}
			log.Println("------")

			insert_sql := fmt.Sprintf(SQL_INSERT_TABLE, strings.ToLower(trans.RingMac), base.GetMac(trans.TagMac), trans.Rssi, trans.Orientation, base.GetLoc(trans.X, trans.Y), "gis") //trans.FingerID)
			log.Println(insert_sql)

			_, err = tx.Exec(insert_sql)
			log.Println(err)
			log.Println("to commit")
			err = tx.Commit()
			if err != nil {
				log.Println("ProccessTransaction commit error")
				log.Println(err)
			}
			//		close(transactions)
		}
	}
}

package mqtt_srv

import (
	"encoding/json"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/giskook/inl2_das_sample/base"
	"github.com/giskook/inl2_das_sample/db"
	"log"
)

var recv_func MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
	log.Println("recv from mqtt")
	log.Println(string(msg.Payload()))
	GetMqttSocket().RecvStringChan <- string(msg.Payload())
}

func (mqtt_socket *Mqtt_srv_socket) Recv() {
	for {
		select {
		case p := <-mqtt_socket.RecvStringChan:
			log.Println("mqtt recv")
			//mqtt_socket.MutexRecvChan.Lock()

			//if mqtt_socket.RecvChan == nil {
			//	mqtt_socket.RecvChan = make(chan *base.BluetoothRing, 1024)
			//}
			//		mqtt_socket.MutexRecvChan.Unlock()
			var btr base.BluetoothRing
			if err := json.Unmarshal([]byte(p), &btr); err != nil {
				log.Println(err)
			}
			//mqtt_socket.RecvChan <- &btr
			mqtt_socket.MutexRecvChan.Lock()
			db.GetDBSocket().TransactionChan <- &btr
			//		mqtt_socket.RecvChan = nil
			mqtt_socket.MutexRecvChan.Unlock()
		}
	}
}

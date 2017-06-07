package mqtt_srv

import (
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/giskook/inl2_das_sample/base"
	"github.com/giskook/inl2_das_sample/conf"
	//"github.com/giskook/inl2_das_sample/db"
	"log"
	"sync"
	"time"
)

var G_Mqtt_Socket *Mqtt_srv_socket

type Mqtt_srv_socket struct {
	Mqtt_socket    MQTT.Client
	PublishHandler MQTT.MessageHandler
	Options        *conf.MqttConfiguration
	SendChan       chan *base.BluetoothRing
	MutexRecvChan  sync.Mutex
	RecvChan       chan *base.BluetoothRing
	RecvStringChan chan string
	ticker         *time.Ticker
}

func NewMqttSocket(mqtt_opts *conf.MqttConfiguration) *Mqtt_srv_socket {
	opts := MQTT.NewClientOptions().AddBroker(mqtt_opts.Addr).SetClientID(mqtt_opts.ClientID)
	opts.SetCleanSession(true)
	opts.SetDefaultPublishHandler(f)

	mqtt_socket := MQTT.NewClient(opts)
	if token := mqtt_socket.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	log.Println("Connected to MqttServer")

	if token := mqtt_socket.Subscribe(mqtt_opts.TopicSub, 0, recv_func); token.Wait() && token.Error() != nil {
		log.Println(token.Error())
	}
	G_Mqtt_Socket = &Mqtt_srv_socket{
		Mqtt_socket:    mqtt_socket,
		PublishHandler: f,
		Options:        mqtt_opts,
		SendChan:       make(chan *base.BluetoothRing, 256),
		RecvChan:       make(chan *base.BluetoothRing, 256),
		RecvStringChan: make(chan string),
		ticker:         time.NewTicker(time.Duration(2) * time.Second),
	}

	return G_Mqtt_Socket
}

func (mqtt_socket *Mqtt_srv_socket) Send(payload *base.BluetoothRing) {
	log.Println("Mqtt add send list")
	mqtt_socket.SendChan <- payload
}

func (mqtt_socket *Mqtt_srv_socket) Proccess() {
	for {
		select {
		case send_pkg := <-mqtt_socket.SendChan:
			mqtt_socket.send(send_pkg)
			//		case <-mqtt_socket.ticker.C:
			//			mqtt_socket.MutexRecvChan.Lock()
			//			db.GetDBSocket().TransactionChan <- mqtt_socket.RecvChan
			//			mqtt_socket.RecvChan = nil
			//			mqtt_socket.MutexRecvChan.Unlock()
		}
	}
}

//func (mqtt_socket *Mqtt_srv_socket) ProccessSub() {
//	for {
//		select {
//		case <-mqtt_socket.ticker.C:
//			mqtt_socket.MutexRecvChan.Lock()
//			log.Println("batch add ")
//			db.GetDBSocket().TransactionChan <- mqtt_socket.RecvChan
//			//		mqtt_socket.RecvChan = nil
//			mqtt_socket.MutexRecvChan.Unlock()
//		}
//	}
//}
func (mqtt_socket *Mqtt_srv_socket) Close() {
	mqtt_socket.Mqtt_socket.Disconnect(0)
	close(mqtt_socket.SendChan)
	close(mqtt_socket.RecvStringChan)
	mqtt_socket.ticker.Stop()

}

func GetMqttSocket() *Mqtt_srv_socket {
	return G_Mqtt_Socket
}

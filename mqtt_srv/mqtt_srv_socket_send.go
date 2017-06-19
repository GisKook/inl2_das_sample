package mqtt_srv

import (
	"encoding/json"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/giskook/inl2_das_sample/base"
	"log"
)

var f MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
	log.Printf("TOPIC: %s\n", msg.Topic())
	log.Printf("MSG: %s\n", msg.Payload())
}

func (mqtt_socket *Mqtt_srv_socket) send(p *base.BluetoothRing) {
	payload, _ := json.Marshal(p)
	log.Println(payload)
	//log.Println(mqtt_socket.Options.TopicPub)
	if token := mqtt_socket.Mqtt_socket.Publish(p.RingMac, 0, false, payload); token.Wait() && token.Error() != nil {
		log.Println(token.Error())
	}

}

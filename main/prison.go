package main

import (
	"fmt"
	"github.com/giskook/gotcp"
	"github.com/giskook/inl2_das_sample"
	"github.com/giskook/inl2_das_sample/conf"
	"github.com/giskook/inl2_das_sample/db"
	"github.com/giskook/inl2_das_sample/event_handler"
	"github.com/giskook/inl2_das_sample/mqtt_srv"
	"github.com/giskook/inl2_das_sample/server"
	"log"
	"net"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	// read configuration
	configuration, err := conf.ReadConfig("./conf.json")

	checkError(err)
	// creates a tcp listener
	tcpAddr, err := net.ResolveTCPAddr("tcp4", ":"+configuration.Server.BindPort)
	checkError(err)
	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)

	_db, _ := db.NewDbSocket(configuration.DB)
	go _db.ProccessTransaction()

	_mqtt_srv := mqtt_srv.NewMqttSocket(configuration.Mqtt)
	go _mqtt_srv.Recv()
	go _mqtt_srv.Proccess()
	//	go _mqtt_srv.ProccessSub()

	// creates a tcp server
	config := &gotcp.Config{
		PacketSendChanLimit:    20,
		PacketReceiveChanLimit: 20,
	}
	srv := gotcp.NewServer(config, &event_handler.Callback{}, &inl2_das_sample.Pdas_Protocol{})

	// create inl2_das_sample server
	server_conf := &server.ServerConfig{
		Listener:      listener,
		AcceptTimeout: time.Duration(configuration.Server.ConnTimeout) * time.Second,
	}
	cpd_server := server.NewServer(srv, server_conf)
	server.SetServer(cpd_server)
	// starts service
	fmt.Println("listening:", listener.Addr())
	cpd_server.Start()

	// catchs system signal
	chSig := make(chan os.Signal)
	signal.Notify(chSig, syscall.SIGINT, syscall.SIGTERM)
	fmt.Println("Signal: ", <-chSig)

	// stops service
	cpd_server.Stop()
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

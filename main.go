package main

import (
	"Modbus/pkg/modbus"
	"Modbus/pkg/websockets"
	"flag"
	"log"
	"net/http"
)

func main() {
	var addr = flag.String("addr", ":8080", "http service address")
	flag.Parse()
	hub := websockets.NewHub()
	go hub.Run()

	modbusChan := make(chan []modbus.MessangeModbus)

	http.HandleFunc("/", websockets.ServeHome)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		websockets.ServeWs(hub, modbusChan, w, r)
	})
	log.Println("Откройте ваше приложение в браузере http://127.0.0.1:8080")
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}

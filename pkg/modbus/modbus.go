package modbus

import (
	"encoding/binary"
	"github.com/goburrow/modbus"
	"log"
	"time"
)

type ConnectionModbus struct {
	IpAdr string `json:"ip"`
	Port  string `json:"port"`
	//	ConnectTimeout time.Time
	SlaveId  string `json:"slaveid"`
	Function string `json:"func"`
	Adr      string `json:"adr"`
	Quantity string `json:"quantity"`
	Active   bool
}

func (con *ConnectionModbus) Connect(cancel <-chan bool, outModbus chan uint16) {

	handler := modbus.NewTCPClientHandler(con.IpAdr + ":" + con.Port)
	handler.Timeout = 10 * time.Second
	handler.SlaveId = 0x01

	//handler.Logger = log.New(os.Stdout, "test: ", log.LstdFlags)
	// Connect manually so that multiple requests are handled in one connection session
	err := handler.Connect()
	if err != nil {
		log.Println(err)
	}

	client := modbus.NewClient(handler)

	interval := time.Duration(1) * time.Second
	// create a new Ticker
	tk := time.NewTicker(interval)
	// start the ticker by constructing a loop

	for {
		select {
		case <-tk.C:
			results, err := client.ReadInputRegisters(15, 1)
			if err != nil {
				log.Println(err)
			}
			result := binary.BigEndian.Uint16(results[:2])
			//fmt.Println(result)
			outModbus <- result
		case <-cancel:
			tk.Stop()
			handler.Close()
			return
		}
	}
}

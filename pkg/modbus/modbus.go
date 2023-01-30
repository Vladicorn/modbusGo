package modbus

import (
	"encoding/binary"
	"fmt"
	"github.com/goburrow/modbus"
	"log"
	"strconv"
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

type MessangeModbus struct {
	Id    string
	Value string
}

func (con *ConnectionModbus) Connect(cancel <-chan bool, outModbus chan []MessangeModbus) {

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
	adr, err := strconv.Atoi(con.Adr)
	if err != nil {
		log.Println(err)
	}
	quantity, err := strconv.Atoi(con.Quantity)
	if err != nil {
		log.Println(err)
	}
	if err != nil {
		log.Println(err)
	}
	for {
		select {
		case <-tk.C:
			var read []byte
			fmt.Println(con.Function)
			switch con.Function {
			case "ReadCoils":
				read, err = client.ReadCoils(uint16(adr), uint16(quantity))
				if err != nil {
					log.Println(err)
				}
			case "ReadDiscreteInputs":
				read, err = client.ReadDiscreteInputs(uint16(adr), uint16(quantity))
				if err != nil {
					log.Println(err)
				}
			case "ReadHoldingRegisters":
				read, err = client.ReadHoldingRegisters(uint16(adr), uint16(quantity))
				if err != nil {
					log.Println(err)
				}
			case "ReadInputRegisters":
				read, err = client.ReadInputRegisters(uint16(adr), uint16(quantity))
				if err != nil {
					log.Println(err)
				}
			default:
				read, err = client.ReadHoldingRegisters(uint16(adr), uint16(quantity))
				if err != nil {
					log.Println(err)
				}
			}

			results := make([]MessangeModbus, 0, quantity)

			for i := 0; i < len(read)/2; i++ {
				k := i * 2
				result := MessangeModbus{}
				result.Value = fmt.Sprintf("%v", binary.BigEndian.Uint16(read[k:k+2]))
				adrInt := i + adr
				result.Id = strconv.Itoa(adrInt)
				results = append(results, result)
			}

			//fmt.Println(results)
			outModbus <- results
		case <-cancel:
			tk.Stop()
			handler.Close()
			return
		}
	}
}

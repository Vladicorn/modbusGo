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
	Type     string `json:"type"`
	Value    string `json:"value"`
	Client   modbus.Client
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
	con.Client = client
	interval := time.Duration(1) * time.Second
	// create a new Ticker
	tk := time.NewTicker(interval)
	// start the ticker by constructing a loop
	adr, err := strconv.Atoi(con.Adr)
	if err != nil {
		results := make([]MessangeModbus, 0, 1)
		result := MessangeModbus{}
		result.Value = fmt.Sprintf("%s", err)
		results = append(results, result)
		log.Println(err)
	}
	quantity, err := strconv.Atoi(con.Quantity)
	if err != nil {
		log.Println(err)
	}

	for {
		select {
		case <-tk.C:
			var read []byte
			results := make([]MessangeModbus, 0, quantity)
			switch con.Function {
			case "ReadCoils":
				if err != nil {
					result := MessangeModbus{}
					result.Value = fmt.Sprintf("%s", err)
					results = append(results, result)
				}
			case "ReadDiscreteInputs":
				if err != nil {
					result := MessangeModbus{}
					result.Value = fmt.Sprintf("%s", err)
					results = append(results, result)
				}
			case "ReadHoldingRegisters":

				read, err = client.ReadHoldingRegisters(uint16(adr), uint16(quantity))
				if err != nil {
					result := MessangeModbus{}
					result.Value = fmt.Sprintf("%s", err)
					results = append(results, result)
				}
			case "ReadInputRegisters":
				read, err = client.ReadInputRegisters(uint16(adr), uint16(quantity))
				if err != nil {
					result := MessangeModbus{}
					result.Value = fmt.Sprintf("%s", err)
					results = append(results, result)
				}
			default:

				read, err = client.ReadHoldingRegisters(uint16(adr), uint16(quantity))
				if err != nil {
					result := MessangeModbus{}
					result.Value = fmt.Sprintf("%s", err)
					results = append(results, result)
				}
			}

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

func (con *ConnectionModbus) Send() {

	adr, err := strconv.Atoi(con.Adr)
	if err != nil {
		log.Println(err)
	}
	value, err := strconv.Atoi(con.Value)
	if err != nil {
		log.Println(err)
	}
	switch con.Function {
	case "WriteSingleRegister":
		_, err = con.Client.WriteSingleRegister(uint16(adr), uint16(value))
		if err != nil {
			log.Println(err)
		}
	case "WriteSingleCoil":
		_, err = con.Client.WriteSingleCoil(uint16(adr), uint16(value))
		if err != nil {
			log.Println(err)
		}
	}

}

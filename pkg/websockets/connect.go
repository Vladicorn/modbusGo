package websockets

import (
	"Modbus/pkg/modbus"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// serveWs handles websocket requests from the peer.
func ServeWs(hub *Hub, modbusChan chan []modbus.MessangeModbus, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256)}
	client.hub.register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump(modbusChan)
	go client.readPump(modbusChan)
}

// Hub maintains the set of active clients and broadcasts messages to the
// clients.

func ServeHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "home.html")
}

var content = `
<!DOCTYPE html>
<html lang="en">
<head>
    <div>
        <h1>HTTP Modbus TCP </h1>
    </div>
<div class="main">


    <form id="form" >
        <div class="field">
            <label for="ip">Введите IP адрес</label>
            <input type="text" id="ip"  minlength="7" maxlength="15"  pattern="^((\d{1,2}|1\d\d|2[0-4]\d|25[0-5])\.){3}(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])$">
        </div>
        <div class="field">
            <label for="port">Введите порт</label>
            <input type="text" id="port"   minlength="1">
        </div>
        <div class="field">
            <label for="slaveid">Введите адрес устройства</label>
            <input type="text" id="slaveid" minlength="1"  >
        </div>
        <div class="field">
            <label for="func">Введите функцию считывания</label>
            <select id="func">
                <option value="ReadHoldingRegisters">ReadHoldingRegisters</option>
                <option value="ReadCoils">ReadCoils</option>
                <option value="ReadDiscreteInputs">ReadDiscreteInputs</option>

                <option value="ReadInputRegisters">ReadInputRegisters</option>
            </select>
        </div>
        <div class="field">
            <label for="adr">Введите адрес регистра</label>
            <input type="text" id="adr" minlength="1" >
        </div>
        <div class="field">
            <label for="quantity">Количество регистров</label>
            <input type="text" id="quantity" minlength="1" >
        </div>
        <div class="field">
            <input type="submit"  value="Начать/Остановить" />
        </div>
    </form>


    <form id="writeForm">
        <H2> Запись</H2>
        <div class="field">

            <label for="writeAdr">Введите адрес</label>
            <input type="text" id="writeAdr" minlength="1" >
        </div>
        <div class="field">
            <label for="writeFunc">Введите функцию записи</label>
            <select id="writeFunc">
                <option value="WriteSingleRegister">WriteSingleRegister</option>
                <option value="WriteSingleCoil">WriteSingleCoil</option>
            </select>
        </div>
        <div class="field">
            <label for="writeValue">Введите значение</label>
            <input type="text" id="writeValue">
        </div>

        <div class="field">
            <input type="submit" value="Записать" />
        </div>
    </form>
    </div>
    <title>Modbus Apeyron</title>
    <script type="text/javascript">
        window.onload = function () {
            var conn;
            var ip = document.getElementById("ip");
            var port = document.getElementById("port");
            var slaveid= document.getElementById("slaveid");
            var func= document.getElementById("func");
            var adr= document.getElementById("adr");
            var quantity= document.getElementById("quantity");

            var writeAdr= document.getElementById("writeAdr");
            var writeFunc= document.getElementById("writeFunc");
            var writeValue= document.getElementById("writeValue");

            var log = document.getElementById("log");
            var configModbus
            var writeModbus
            function appendLog(item) {
                var doScroll = log.scrollTop > log.scrollHeight - log.clientHeight - 1;
                log.appendChild(item);
                if (doScroll) {
                    log.scrollTop = log.scrollHeight - log.clientHeight;
                }
            }

            document.getElementById("form").onsubmit = function () {
                if (!conn) {
                    return false;
                }
                if (!ip.value) {
                    return false
                }
                configModbus = JSON.stringify({ ip: ip.value, port: port.value, slaveid: slaveid.value,func: func.value, quantity: quantity.value,adr: adr.value, type: "config" })

                conn.send(configModbus);

                return false;
            };

            document.getElementById("writeForm").onsubmit = function () {
                if (!conn) {
                    return false;
                }

                writeModbus = JSON.stringify({ adr : writeAdr.value,func: writeFunc.value, value: writeValue.value , type: "write"  })
                conn.send(writeModbus);

                return false;
            };

            if (window["WebSocket"]) {
                conn = new WebSocket("ws://" + document.location.host + "/ws");
                conn.onclose = function (evt) {
                    var item = document.createElement("div");
                    item.innerHTML = "<b>Connection closed.</b>";
                    appendLog(item);
                };
                conn.onmessage = function (evt) {
                   // console.log(evt.data)
                  //  var messages = evt.data.split('-');
                    log.innerText=evt.data;
                    var item = document.createElement("input type=\"text\"");
                    log.appendChild(item)
                    /*
                    for (var i = 0; i < messages.length; i++) {
                        var item = document.createElement("div");
                         item.innerText = messages[i];
                        log.innerText=messages[i];
                    //    log.Child
                      //  console.log(messages[i])
                        //appendLog(item);
                    }

                     */
                };
            } else {
                var item = document.createElement("div");
                item.innerHTML = "<b>Your browser does not support WebSockets.</b>";
                appendLog(item);
            }
        };
    </script>
    <style type="text/css">
        html {
            overflow: hidden;
        }

        body {
            overflow: hidden;
            padding: 0;
            margin: 0;
            width: 100%;
            height: 100%;
            background: rgb(187,187,237);
            background: linear-gradient(90deg, rgba(187,187,237,1) 3%, rgba(123,163,201,1) 62%, rgba(25,25,29,1) 100%, rgba(0,212,255,1) 100%);
        }

        #log {
            background: white;
            margin: 10px;
            height: 45%;
            padding: 0.5em 0.5em 0.5em 0.5em;
            position: absolute;
            top: 29em;
            left: 0.5em;
            right: 0.5em;
            bottom: 3em;
            overflow: auto;
        }
        .field {clear:both; text-align:right; line-height:25px;}
        .main {float:left}
/*
        #form {
            padding: 0 0.5em 0 0.5em;
            margin: 0;
            position: absolute;
            bottom: 1em;
            left: 0px;
            width: 100%;
            overflow: hidden;
        }


 */
        label {float:left; padding-right:10px;}
        input,
        label {
            margin: 1px;
        }

    </style>
</head>
<body>
<div id="log"></div>

</body>
</html>
`

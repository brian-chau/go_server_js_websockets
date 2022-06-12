package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "localhost:8080", "http service address")
var upgrader = websocket.Upgrader{} // use default options
const TextType = 1

func SendDataToScreen(c *websocket.Conn, err error, generatedDataChannel chan []byte){
	for {
		dataToSend := <- generatedDataChannel
		err = c.WriteMessage(TextType, dataToSend)
		if err != nil {
			log.Println("write:", err)
			break
		}
		time.Sleep(1 * time.Second)
	}
}

func ReadDataFromScreen(c *websocket.Conn, err error) {
	for {
		_, message, err := c.ReadMessage()
		cmd := strings.Split(string(message), "=")

		switch cmd[0] {
			case "something":
				fmt.Printf("Doing something with %s\n", cmd[1])
			case "somethingElse":
				fmt.Printf("Doing something else with %s\n", cmd[1])
			default:
				fmt.Printf("Unsupported cmd %s\n", cmd[0])
		}

		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s\n", message)
	}
}

func Handler(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()

	generatedDataChannel := make(chan []byte)

	var wg sync.WaitGroup
	wg.Add(3)
	go SendDataToScreen(c, err, generatedDataChannel)
	go ReadDataFromScreen(c, err)
	go GenerateData(generatedDataChannel)
	wg.Wait()
}

func home(w http.ResponseWriter, r *http.Request) {
	homeTemplate.Execute(w, "ws://"+r.Host+"/handler")
}

func main() {
	flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/", home)
	http.HandleFunc("/handler", Handler)

	log.Fatal(http.ListenAndServe(*addr, nil))
}

var homeTemplate = template.Must(template.ParseFiles("home.html", "script.html"))
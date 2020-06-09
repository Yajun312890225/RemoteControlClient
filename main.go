package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/url"
	"os"
	"time"

	"github.com/gorilla/websocket"
	"github.com/stianeikeland/go-rpio"
)

type Message struct {
	Type   int           `json:"type"`   //1打开 0 关闭
	During time.Duration `json:"during"` //持续时间
}

type Login struct {
	MacId string `json:"macid"`
}

var addr = flag.String("addr", "192.168.128.29:7777", "http service address")

func main() {

	if err := rpio.Open(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer rpio.Close()
	pin := rpio.Pin(24)
	pin.Output()

	u := url.URL{Scheme: "ws", Host: *addr, Path: "/ws"}
	var dialer *websocket.Dialer

	conn, _, err := dialer.Dial(u.String(), nil)
	for err != nil {
		time.Sleep(5 * time.Second)
		conn, _, err = dialer.Dial(u.String(), nil)
	}

	inter, err := net.InterfaceByName("wlan0")
	if err != nil {
		log.Fatalln(err)
	}
	conn.WriteJSON(Login{MacId: inter.HardwareAddr.String()})

	for {

		msg := Message{}
		err := conn.ReadJSON(&msg)
		if err != nil {
			fmt.Println("read:", err)
			time.Sleep(5 * time.Second)
			conn, _, err = dialer.Dial(u.String(), nil)
			for err != nil {
				time.Sleep(5 * time.Second)
				conn, _, err = dialer.Dial(u.String(), nil)
			}
			conn.WriteJSON(Login{MacId: inter.HardwareAddr.String()})
			continue
		}
		if msg.Type == 1 {
			pin.High()
			fmt.Println(msg.During)
			if msg.During != 0 {
				timer := time.NewTimer(msg.During)
				go func() {
					<-timer.C
					pin.Low()
				}()
			}
		} else if msg.Type == 2 {
			pin.Low()
		}
	}

}

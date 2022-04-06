package main

import (
	"github.com/elamre/gomberman/net/client"
	"github.com/elamre/gomberman/net/server"
	"log"
	"time"
)

const port = 50000

func main() {
	server := server.New(server.Options{
		MaxConnections: 5,
		HttpPort:       port,
		PublicIP:       "127.0.0.1",
		ICEServerURLs:  []string{"stun:127.0.0.1:3478"},
	})
	server.Start()
	go func() {
		for {
			time.Sleep(time.Second)
			for _, conn := range server.Connections() {
				if conn.IsConnected() {
					log.Println("connection connected")
					dat, d := conn.Read()
					if d {
						log.Printf("server got: % 02x", dat)
						if err := conn.Send(dat); err != nil {
							panic(err)
						}
					}
				}
			}
		}
	}()
	log.Println(" Success ?")

	c := client.New(client.Options{
		IPAddress:     "127.0.0.1:50000",
		ICEServerURLs: []string{"stun:127.0.0.1:3478"},
	})
	c.Start()
	for i := 0; i < 10; i++ {
		log.Println(c.GetLastError())
		time.Sleep(time.Second)
		if c.IsConnected() {
			log.Println("Connected")
			if err := c.Send([]byte{1, 2, 3, 4}); err != nil {
				panic(err)
			}
			dat, s := c.Read()
			if s {
				log.Println(dat)
			}
		}
	}
}

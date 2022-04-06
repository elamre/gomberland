package main

import (
	"github.com/elamre/gomberman/net/client"
	"github.com/elamre/gomberman/net/server"
	"github.com/elamre/gomberman/netcode/core"
	"github.com/elamre/gomberman/netcode/local"
	"github.com/elamre/gomberman/netcode/packet"
	"log"
	"reflect"
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

var testMap = map[reflect.Type]int{}
var mapTest = map[int]reflect.Type{}
var index = 0

func registerStruct(someStruct any) int {
	tt := reflect.TypeOf(someStruct)
	if _, ok := testMap[tt]; ok {
		panic("Already exists")
	}
	testMap[tt] = index
	mapTest[index] = tt
	index++
	log.Printf("Registered at: %d", index-1)
	return index - 1
}

func returnPackage(index int) any {
	return reflect.New(mapTest[index]).Interface()
}

type SomeInterface interface {
	SayContents()
}

type SomeStructA struct {
	Number int
}

func (s SomeStructA) SayContents() {
	log.Println("Stuct A")
}

type SomeStructB struct {
}

func (s SomeStructB) SayContents() {
	log.Println("Stuct B")
}

func main1() {
	core.RegisterPackets()
	network := local.NewFakeNetwork()
	client := local.NewLocalClient(network)
	server := local.NewLocalServer(network)

	go func() {
		for {
			pack := server.GetPacket()
			pp := *pack

			switch t := pp.(type) {
			case *packet.ChatPacket:
				log.Println("We do have A yes!")

			default:
				log.Printf("unknown type: %T", t)
			}
		}
	}()

	client.Write(packet.ChatPacket{
		Message: "Test",
	})
	time.Sleep(10 * time.Second)
	_ = client

}

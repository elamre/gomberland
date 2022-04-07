package main

import (
	"github.com/elamre/gomberman/net/webrtc_client"
	"github.com/elamre/gomberman/net/webrtc_server"
	"github.com/elamre/gomberman/netcode/core"
	"github.com/elamre/gomberman/netcode/local"
	"github.com/elamre/gomberman/netcode/packet"
	"github.com/elamre/gomberman/netcode/webrtc"
	"log"
	"time"
)

const port = 50000

type NetworkWorld struct {
}

func main1() {
	core.RegisterPackets()
	server := webrtc_server.New(webrtc_server.Options{
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
						log.Printf("webrtc_server got: % 02x", dat)
						if err := conn.Send(dat); err != nil {
							panic(err)
						}
					}
				}
			}
		}
	}()
	log.Println(" Success ?")

	c := webrtc_client.New(webrtc_client.Options{
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

func benchMark() {
	network := local.NewFakeNetwork()
	cclient := local.NewLocalClient(network)
	cserver := local.NewLocalServer(network)
	packets := 1000000
	go func() {
		index := 0
		start := time.Now()
		for {
			pack := cserver.GetPacket()
			index++
			pp := pack
			_ = pp
			if index >= packets {
				break
			}
			/*			log.Printf("%+v", pack)
						switch t := pp.(type) {
						case packet.ChatPacket:
							log.Println("Message: " + t.Message)
							log.Println("We do have A yes!")

						default:
							log.Printf("unknown type: %T", t)
						}*/
		}
		log.Printf("webrtc_server took: %d", time.Now().Sub(start).Milliseconds())
	}()

	start := time.Now()
	for i := 0; i < packets; i++ {
		cclient.Write(packet.ChatPacket{
			Message: "Test",
		})
	}
	log.Printf("Took: %vms", time.Now().Sub(start).Milliseconds())
}

func main() {
	core.RegisterPackets()

	server := webrtc_server.New(webrtc_server.Options{
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
						rawPacket := core.RawPacketFrom(dat)
						log.Printf("webrtc_server got: %+v", rawPacket)
						if err := conn.Send(dat); err != nil {
							panic(err)
						}
					}
				}
			}
		}
	}()

	time.Sleep(time.Second)
	client := webrtc.NewWebrtcClient("127.0.0.1", port)
	log.Printf("Error: %v", client.Connect())
	client.Write(packet.NewRegisterPacket("Elmar"))

	time.Sleep(3 * time.Second)
	c1 := webrtc.NewWebrtcClient("127.0.0.1", port)
	c1.Connect()
	c2 := webrtc.NewWebrtcClient("127.0.0.1", port)
	c2.Connect()
	c3 := webrtc.NewWebrtcClient("127.0.0.1", port)
	c3.Connect()
	time.Sleep(10 * time.Second)
}

func main2() {
	core.RegisterPackets()
	time.Sleep(time.Second)
	network := local.NewFakeNetwork()
	cclient := local.NewLocalClient(network)
	cserver := local.NewLocalServer(network)

	go func() {
		pack := cserver.GetPacket()
		log.Printf("Server received: %+v", pack)
		cserver.Write(packet.ConnectionPacket{
			UserId:  1,
			Action:  2,
			Message: "Registered",
		})
		pack = cserver.GetPacket()
		log.Printf("Server received: %+v", pack)
	}()

	cclient.Write(packet.NewRegisterPacket("Elmar"))
	time.Sleep(time.Second)
	pack := cclient.WaitForPacket()
	log.Printf("Client got: %+v", pack)
	time.Sleep(time.Second)
	cclient.Write(packet.ChatPacket{Message: "Registered!"})

	time.Sleep(10 * time.Second)
}

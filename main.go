package main

import (
	"github.com/elamre/gomberman/common_system/packets"
	"github.com/elamre/gomberman/lobby_system"
	packets2 "github.com/elamre/gomberman/lobby_system/packets"
	local2 "github.com/elamre/gomberman/net/local"
	webrtc2 "github.com/elamre/gomberman/net/webrtc"
	"log"
	"time"
)

const port = 50001

type NetworkWorld struct {
}

/*
func main1() {
	net_game.RegisterPackets()
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
*/
func benchMark() {
	network := local2.NewFakeNetwork()
	cclient := local2.NewLocalClient(network)
	_ = cclient
	cserver := local2.NewLocalServer(network)
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
	/*	for i := 0; i < packets; i++ {
		cclient.Write(packets.ChatPacket{
			Message: "Test",
		})
	}*/
	log.Printf("Took: %vms", time.Now().Sub(start).Milliseconds())
}

func main() {
	packets.Register()
	packets2.Register()

	serverLobby := lobby_system.NewLobbyServerSystem(webrtc2.NewWebrtcHost("127.0.0.1", port))

	go func() {
		for {
			time.Sleep(time.Second)
			serverLobby.Update()
		}
	}()
	client := webrtc2.NewWebrtcClient("127.0.0.1", port)
	client.Connect()
	for !client.IsConnected() {
	}
	clientLobby := lobby_system.NewLobbyClientSystem(client)

	clientLobby.RegisterPlayer("Elmar")
	clientLobby.SendPacket(packets2.RoomPacket{
		Action:   packets2.RoomCreateAction,
		Password: "werwe",
		Name:     "ElmaR",
	})
	go func() {
		for {
			time.Sleep(time.Second)
			clientLobby.Update()
			clientLobby.SendPacket(packets2.RoomPacket{
				Action:   packets2.RoomCreateAction,
				Password: "werwe",
				Name:     "ElmaR2",
			})
		}
	}()
	time.Sleep(10 * time.Second)
}

func main2() {
	packets.Register()
	packets2.Register()

	time.Sleep(time.Second)
	network := local2.NewFakeNetwork()
	cclient := local2.NewLocalClient(network)
	cserver := local2.NewLocalServer(network)

	go func() {
		pack := cserver.GetPacket()
		log.Printf("Server received: %+v", pack)
		cserver.Write(packets.ConnectionPacket{
			UserId:  1,
			Action:  2,
			Message: "Registered",
		})
		pack = cserver.GetPacket()
		log.Printf("Server received: %+v", pack)
	}()

	cclient.Write(packets.NewRegisterPacket("Elmar"))
	time.Sleep(time.Second)
	pack := cclient.WaitForPacket()
	log.Printf("Client got: %+v", pack)
	time.Sleep(time.Second)
	cclient.Write(packets.ChatPacket{Message: "Registered!"})

	time.Sleep(10 * time.Second)
}

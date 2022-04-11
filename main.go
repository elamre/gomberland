package main

import (
	"github.com/elamre/gomberman/common_system/common_packets"
	"github.com/elamre/gomberman/game_system"
	"github.com/elamre/gomberman/lobby_system"
	packets2 "github.com/elamre/gomberman/lobby_system/lobby_system_packets"
	local2 "github.com/elamre/gomberman/net/local"
	webrtc2 "github.com/elamre/gomberman/net/webrtc"
	"github.com/elamre/gomberman/ping_system"
	"github.com/elamre/gomberman/ping_system/ping_packets"
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
	/*	for i := 0; i < game_system_packets; i++ {
		cclient.Write(game_system_packets.ChatPacket{
			Message: "Test",
		})
	}*/
	log.Printf("Took: %vms", time.Now().Sub(start).Milliseconds())
}

var client = false

func main() {
	common_packets.Register()
	packets2.Register()
	ping_packets.Register()

	if !client {
		go func() {
			server := webrtc2.NewWebrtcHost("127.0.0.1", port)
			//server := webrtc2.NewWebrtcHost("192.168.178.43", port)
			//server := webrtc2.NewWebrtcHost("78.47.36.203", port)
			serverDelegator := NewServerDelegator(server)
			serverDelegator.RegisterSubSystem("serverlobby", lobby_system.NewLobbyServerSystem(server))
			serverDelegator.RegisterSubSystem("ping", ping_system.NewPingServerSystem())
			serverDelegator.RegisterSubSystem("game", game_system.NewGameServerSystem(nil, server, game_system.GameServerSystemOptions{TicksPerSecond: 5}))
			if false {
				client = true
			}
			for {
				serverDelegator.Update()
			}
		}()
	}
	//client := webrtc2.NewWebrtcClient("192.168.178.43", port)
	//client := webrtc2.NewWebrtcClient("78.47.36.203", port)
	client := webrtc2.NewWebrtcClient("127.0.0.1", port)
	client.Connect()
	for !client.IsConnected() {
	}
	clientDelegator := NewClientDelegator(client)
	clientLobby := lobby_system.NewLobbyClientSystem(client)
	clientDelegator.RegisterSubSystem("clientlobby", clientLobby)
	clientDelegator.RegisterSubSystem("ping", ping_system.NewPingClientSystem(client))

	clientLobby.RegisterPlayer("Elmar")
	clientLobby.OnRegisteredAction = func() {

		clientLobby.SendRoomPacket(packets2.RoomPacket{
			Action:   packets2.RoomCreateAction,
			Password: "werwe",
			Name:     "ElmaR",
		})
		clientLobby.SendRoomPacket(packets2.RoomPacket{Action: packets2.RoomReadyAction, Name: "ElmaR"})
		clientLobby.SendRoomPacket(packets2.RoomPacket{Action: packets2.RoomStartAction, Name: "ElmaR"})
	}
	go func() {
		for {
			clientDelegator.Update()

		}
	}()
	time.Sleep(10 * time.Second)
}

func main2() {
	common_packets.Register()
	packets2.Register()

	time.Sleep(time.Second)
	network := local2.NewFakeNetwork()
	cclient := local2.NewLocalClient(network)
	cserver := local2.NewLocalServer(network)

	go func() {
		pack := cserver.GetPacket()
		log.Printf("Server received: %+v", pack)
		cserver.Write(common_packets.ConnectionPacket{
			UserId:  1,
			Action:  2,
			Message: "Registered",
		})
		pack = cserver.GetPacket()
		log.Printf("Server received: %+v", pack)
	}()

	cclient.Write(common_packets.NewRegisterPacket("Elmar"))
	time.Sleep(time.Second)
	pack := cclient.WaitForPacket()
	log.Printf("Client got: %+v", pack)
	time.Sleep(time.Second)
	cclient.Write(common_packets.ChatPacket{Message: "Registered!"})

	time.Sleep(10 * time.Second)
}

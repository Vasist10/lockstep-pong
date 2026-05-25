package main

import (
	"flag"
	"fmt"
	"lockstep-pong/bridge"
	"lockstep-pong/relay"
	"lockstep-pong/render"
	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	role := flag.String("role", "server", "Role to play: server or client")
	playerID := flag.Uint("player", 1, "Player ID (1 or 2)")
	flag.Parse()

	if *role == "server" {
		server, err := relay.NewServer(9999)
		if err != nil {
			fmt.Printf("Failed to start server: %v\n", err)
			return
		}
		server.Start()
		select {}
	} else if *role == "client" {
		client, err := relay.NewClient("localhost:9999", uint8(*playerID))
		if err != nil {
			fmt.Printf("Failed to create client: %v\n", err)
			return
		}
		client.Start()
		var state bridge.GameState
		bridge.SimInit(&state)
		var currentTick uint64 = 0
		buffer := make(map[uint64]map[uint8]relay.InputPacket) //tick -> playerID -> packet
		
		onUpdate := func() {
			//send local input
			var keys uint8 = 0
			if ebiten.IsKeyPressed(ebiten.KeyW) || ebiten.IsKeyPressed(ebiten.KeyUp) {
				keys |= 1 //player 1 up
			} 
			if ebiten.IsKeyPressed(ebiten.KeyS) || ebiten.IsKeyPressed(ebiten.KeyDown) {
				keys |= 2 //player 2 down
			}
			packet := relay.InputPacket{
				Tick:     currentTick,
				PlayerID: uint8(*playerID),
				Keys:     keys,
				Checksum: 0,
			}
			//serialize
			data := relay.Serialize(packet)
			err := client.Send(data)
			if err != nil {
				fmt.Printf("Failed to send packet: %v\n", err)
			}
			//receive remote input
			DrainLoop:
			for {
				select {
				case packet := <-client.RecvChan:
					recvPacket, err := relay.Deserialize(packet.Data)
					if err != nil {
						fmt.Printf("Failed to deserialize packet: %v\n", err)
						continue
					}
					if _, ok := buffer[recvPacket.Tick]; !ok {
						buffer[recvPacket.Tick] = make(map[uint8]relay.InputPacket)
					}
					buffer[recvPacket.Tick][recvPacket.PlayerID] = recvPacket
							
				default:
					break DrainLoop			
				}
					
			}
			if len(buffer[currentTick]) == 2 {
				//advance simulation
				inputs := bridge.InputSet{
					P1Keys: buffer[currentTick][1].Keys,
					P2Keys: buffer[currentTick][2].Keys,
				}
				bridge.SimTick(&state, inputs)
				fmt.Printf("tick=%d hash=%d\n",currentTick,bridge.SimHash(&state))
				//delete old ticks from buffer
				delete(buffer, currentTick) 
				currentTick++
			}
		}
		game := render.NewPongGame(&state, onUpdate)
		ebiten.SetWindowSize(800, 600)
		ebiten.SetWindowTitle("Lockstep Pong")
		ebiten.RunGame(game)
			
	}
}

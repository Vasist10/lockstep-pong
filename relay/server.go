
package relay

import(
	"net"
	"fmt"
	"sync"

)

const MaxClients = 2

type Packet struct {
	Data []byte
	Addr *net.UDPAddr
}

type Server struct {
	Conn *net.UDPConn
	Clients []*net.UDPAddr
	ClientChans [MaxClients]chan Packet
	PacketChan chan Packet
	mu sync.Mutex
}

func NewServer(port int) (*Server, error){
	addr := net.UDPAddr{
		Port: port,
		IP: net.ParseIP("0.0.0.0"),
	}
	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		return nil, fmt.Errorf("failed to listen on UDP port %d: %w", port, err)
	}
	var clientChans [MaxClients]chan Packet
	for i := 0; i < MaxClients; i++ {
		clientChans[i] = make(chan Packet, 100)
	}

	return &Server{
		Conn: conn,
		Clients: make([]*net.UDPAddr, 0, MaxClients),
		ClientChans: clientChans,
		PacketChan: make(chan Packet, 100),
		mu: sync.Mutex{},
	}, nil
}

func (s *Server) isClient(addr *net.UDPAddr) bool {
	for _, client := range s.Clients {
		if client != nil && client.String() == addr.String() {
			return true
		}
	}
	return false
}
func (s *Server) Start(){
	//receive go routine
	go func(){
		buf := make([]byte, 1024)
		for {
			n, addr, err := s.Conn.ReadFromUDP(buf)
			if err != nil {
				fmt.Printf("Error reading UDP message: %v\n", err)
				continue
			}
			//allocate stable mem for this pckt
			data := make([]byte, n)
			//copy pckt data to stable mem
			copy(data, buf[:n])

			packet := Packet{
				Data: data,
				Addr: addr,
			}
			s.PacketChan <- packet
		}
	}()
	//dispenser go routine like broadcast
	go func(){
		for packet := range s.PacketChan{
			s.mu.Lock()
			if !s.isClient(packet.Addr) && len(s.Clients) < MaxClients {
				s.Clients = append(s.Clients, packet.Addr)
				fmt.Printf("New client connected: %s\n", packet.Addr.String())
			}else if !s.isClient(packet.Addr) {
				fmt.Printf("Max clients reached. Ignoring packet from: %s\n", packet.Addr.String())
				s.mu.Unlock()
				continue
			}
			
			clientCount := len(s.Clients)
			s.mu.Unlock()
			for i := 0;i < clientCount;i++ {
				s.ClientChans[i] <- packet
			}	
		}
		
	}()

	for i := 0; i < MaxClients; i++ {
		go func(clientIndex int){
			for packet := range s.ClientChans[clientIndex]{
				s.mu.Lock()
				if clientIndex >= len(s.Clients){
					s.mu.Unlock()
					continue
				}
				ClientAddr := s.Clients[clientIndex]
				s.mu.Unlock()

				_, err := s.Conn.WriteToUDP(packet.Data, ClientAddr)
				if err != nil {
					fmt.Printf("Error writing UDP message to %s: %v\n", ClientAddr.String(), err)
				}
			}
		}(i)
	}

}


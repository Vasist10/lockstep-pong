package relay

import (
	"fmt"
	"net"
)

type Client struct {
	Conn       *net.UDPConn
	ServerAddr *net.UDPAddr
	RecvChan   chan Packet
	PlayerID   uint8
}

func NewClient(serverAddr string, playerID uint8) (*Client, error) {
	addr, err := net.ResolveUDPAddr("udp", serverAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve UDP address: %w", err)
	}
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return nil, fmt.Errorf("failed to dial UDP: %w", err)
	}
	return &Client{
		Conn:       conn,
		ServerAddr: addr,
		RecvChan:   make(chan Packet, 100),
		PlayerID:   playerID,
	}, nil
}

func (c *Client) Start() {
	go func() {
		buf := make([]byte, 1024)
		for {
			n, err := c.Conn.Read(buf)
			if err != nil {
				fmt.Printf("Error reading from UDP connection: %v\n", err)
				continue
			}
			data := make([]byte, n)
			copy(data, buf[:n])
			packet := Packet{
				Data: data,
				Addr: c.ServerAddr,
			}
			c.RecvChan <- packet
		}

	}()
}

func (c *Client) Send(data []byte) error {
	_, err := c.Conn.Write(data)
	return err
}

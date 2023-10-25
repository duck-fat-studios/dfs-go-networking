package dfsnetworking

import "net"

type UDPComms struct {
	incomingConnection *net.UDPConn
	outgoingConnection *net.UDPConn
	readBuffer         []byte
	handler            UDPHandlerInterface
	LocalIP            net.IP
}

type UDPHandlerInterface interface {
	Handle(data []byte, addr *net.UDPAddr)
}

func NewUDPComms(localIP net.IP,  handler UDPHandlerInterface, incomingPort int, outgoingPort ...int) (*UDPComms, error) {
	inAddr := net.UDPAddr{
		Port: incomingPort,
		IP:   localIP,
	}
	incomingConnection, err := net.ListenUDP("udp", &inAddr)
	if err != nil {
		return nil, err
	}

	var outgoingConnection *net.UDPConn
	if len(outgoingPort) == 0 || outgoingPort[0] == incomingPort {
		outgoingConnection = incomingConnection
	} else {
		outAddr := net.UDPAddr{
			Port: outgoingPort[0],
			IP:   localIP,
		}
		outgoingConnection, err = net.ListenUDP("udp", &outAddr)
		if err != nil {
			return nil, err
		}
	}

	return &UDPComms{
		incomingConnection: incomingConnection,
		outgoingConnection: outgoingConnection,
		readBuffer:         make([]byte, 4096),
		handler:            handler,
		LocalIP:            localIP,
	}, nil
}

func (u *UDPComms) Run() error {
	for {
		buff := make([]byte, 4096)
		n, addr, err := u.incomingConnection.ReadFromUDP(buff)

		if err != nil {
			return err
		}

		go u.handler.Handle(buff[:n], addr)
	}
}

func (u *UDPComms) Send(data []byte, remoteAddr *net.UDPAddr) error {
	_, err := u.outgoingConnection.WriteToUDP(data, remoteAddr)
	return err
}

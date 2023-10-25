package dfsnetworking

import "net"


type UDPComms struct {
	incomingConnection *net.UDPConn
	outgoingConnection *net.UDPConn
	readBuffer []byte
	Handler func([]byte, *net.UDPAddr)
	Interface *net.Interface
	LocalIP net.IP
}


func NewUDPComms(localIP net.IP, localInPort, localOutPort int, handler func([]byte, *net.UDPAddr), iface *net.Interface) (*UDPComms, error) {
    inAddr := net.UDPAddr{
        Port: localInPort,
        IP:   localIP,
    }
    incomingConnection, err := net.ListenUDP("udp", &inAddr)
    if err != nil {
        return nil, err
    }

    outAddr := net.UDPAddr{
        Port: localOutPort,
        IP:   localIP,
    }
    outgoingConnection, err := net.ListenUDP("udp", &outAddr)
    if err != nil {
        return nil, err
    }

    return &UDPComms{
        incomingConnection: incomingConnection,
        outgoingConnection: outgoingConnection,
        readBuffer:            make([]byte, 4096),
        Handler:            handler,
        Interface:          iface,
        LocalIP:            localIP,
    }, nil
}


func (u *UDPComms) Run() error{
	for {
		buff := make([]byte, 4096)
		n, addr, err := u.incomingConnection.ReadFromUDP(buff)

		if err !=nil {
			return err
		}

		go u.Handler(buff[:n], addr)
	}
}

func (u *UDPComms) Send(data []byte, remoteAddr *net.UDPAddr) error {
    _, err := u.outgoingConnection.WriteToUDP(data, remoteAddr)
    return err
}
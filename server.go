package rtmp

import "net"

type Server struct {
}

func (srv *Server) Server(l net.Listener, handler ConnHandler) error {
	defer l.Close()

	for {
		rwc, err := l.Accept()
		if err != nil {
			continue
		}

		c := srv.newConn(rwc, handler)
		go c.Serve()
	}
}

func (srv *Server) newConn(rwc net.Conn, handler ConnHandler) *Conn {
	conn := NewConn(rwc, handler)

	return conn
}

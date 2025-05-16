package main

import (
	"fmt"
	"net"
	"time"
)

type LOG_LEVEL int32

const (
	LOG_INFO LOG_LEVEL = iota
	LOG_ERROR
	LOG_DEBUG
	LOG_SUCCESS
)

type Server struct {
	listenAddr string
	ln net.Listener
	quitch chan struct{}
}

func logger(level LOG_LEVEL, message string, args ...interface{}) {
	currentTime := time.Now()
	stringTime := currentTime.Format("2006-01-02 15:04:05")
	fmt.Printf("[%s", stringTime)
	switch level {
	case LOG_INFO:
		fmt.Printf("\033[33;1m    info\033[0m] ")
		break
	case LOG_ERROR:
		fmt.Printf("\033[91;1m   error\033[0m] ")
		break
	case LOG_DEBUG:
		fmt.Printf("\033[35;1m   debug\033[0m] ")
		break
	case LOG_SUCCESS:
		fmt.Printf("\033[92;1m success\033[0m] ")
		break
	}
	fmt.Printf(message, args...)
	fmt.Printf("\n")
}

func NewServer(listenAddr string) *Server {
	return &Server{
		listenAddr: listenAddr,
		quitch: make(chan struct{}),
	}
}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.listenAddr)
	logger(LOG_INFO, "Listening on %s", s.listenAddr)
	if err != nil {
		logger(LOG_ERROR, "%s", err)
		return err
	}
	defer ln.Close()
	s.ln = ln

	go s.acceptLoop()

	<-s.quitch

	
	return nil
}

func (s *Server) acceptLoop() {
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			logger(LOG_ERROR, "%s", err)
			continue
		}

		logger(LOG_INFO, "%s connecting..",conn.RemoteAddr())
		
		go s.readLoop(conn);
	}
}

func (s *Server) readLoop(conn net.Conn) {
	defer conn.Close()
	logger(LOG_SUCCESS, "%s connected !", conn.RemoteAddr())
	buf := make([]byte, 2049)

	//check:= false
	for {
		n, err := conn.Read(buf)
		if err != nil {
			logger(LOG_INFO, "%s disconnected (%s)", conn.RemoteAddr(), err)
			conn.Close()
			break
		}
		/*if check != true {
			check = true
			continue
		}*/

		msg := buf[:n]
		if string(msg) != "\r\n" {
			logger(LOG_DEBUG, "%s", string(msg))
		}
		conn.Write([]byte("data_back\r\n"))
	}
}

func main() {
	server := NewServer(":8080")
	server.Start()
}

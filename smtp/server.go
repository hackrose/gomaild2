package smtp

import (
	"bufio"
	log "github.com/sirupsen/logrus"
	config "github.com/spf13/viper"
	. "github.com/trapped/gomaild2/smtp/structs"
	. "github.com/trapped/gomaild2/structs"
	"net"
)

type Server struct {
	Addr string
	Port string
}

//go:generate gengen -d ./commands/ -t process.go.tmpl -o process.go

func (s *Server) Start() {
	log.WithFields(log.Fields{
		"addr": s.Addr,
		"port": s.Port,
	}).Info("Starting SMTP server")
	l, err := net.Listen("tcp", s.Addr+":"+s.Port)
	if err != nil {
		panic(err.Error())
	}
	for {
		c, _ := l.Accept()
		client := &Client{
			Conn: c,
			Rdr:  bufio.NewReadWriter(bufio.NewReader(c), bufio.NewWriter(c)),
			Data: make(map[string]interface{}),
			ID:   SessionID(12),
		}
		go accept(client)
	}
}

func accept(c *Client) {
	c.Send(Reply{Result: Ready, Message: config.GetString("server.name") + " gomaild2 ESMTP ready"})
	c.State = Connected
	log.WithFields(log.Fields{
		"id":   c.ID,
		"addr": c.Conn.RemoteAddr().String(),
	}).Info("Connected")
	for {
		if c.State == Disconnected {
			break
		}
		cmd, err := c.Receive()
		if err != nil {
			break
		}
		c.Send(Process(c, cmd))
	}
	c.Conn.Close()
	log.WithFields(log.Fields{
		"id":   c.ID,
		"addr": c.Conn.RemoteAddr().String(),
	}).Info("Disconnected")
}

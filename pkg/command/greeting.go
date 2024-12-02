package command

import (
	"bufio"
	"net"
	"time"
)

type Greeting struct {
	reply Reply
}

func (g Greeting) Send(conn *net.TCPConn) error {
	return nil
}

func (g *Greeting) Wait(conn *net.TCPConn) error {

	_ = conn.SetDeadline(time.Now().Add(time.Second))

	reader := bufio.NewReader(conn)

	data, _ := reader.ReadString('\n')
	g.reply = ParseReply(data)

	return nil
}

func (g Greeting) Name() string {
	return "server greeting"
}

func (g Greeting) Reply() Reply {
	return g.reply
}

func ServerGreeting() *Greeting {
	return &Greeting{reply: Reply{}}
}

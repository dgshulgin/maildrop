package command

import (
	"bufio"
	"net"
	"time"
)

type Pass struct {
	reply Reply
	pwd   string
}

func SendPass(pwd string) *Pass {
	return &Pass{pwd: pwd, reply: Reply{}}
}

// PASS<space>password<eol>
func (p Pass) Send(conn *net.TCPConn) error {

	cmd := p.Name() + SpaceDelim + p.pwd + EOL

	_ = conn.SetDeadline(time.Now().Add(time.Second))
	_, err := conn.Write([]byte(cmd))
	if err != nil {
		return err
	}

	return nil
}

func (p *Pass) Wait(conn *net.TCPConn) error {

	_ = conn.SetDeadline(time.Now().Add(time.Second))

	reader := bufio.NewReader(conn)

	data, _ := reader.ReadString('\n')
	p.reply = ParseReply(data)

	return nil
}

func (p Pass) Reply() Reply {
	return p.reply
}

func (p Pass) Name() string {
	return "PASS"
}

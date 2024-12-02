package command

import (
	"bufio"
	"net"
	"time"
)

type Quit struct {
	reply Reply
}

func SendQuit() *Quit {
	return &Quit{reply: Reply{}}
}

func (q Quit) Send(conn *net.TCPConn) error {

	cmd := q.Name() + EOL

	_ = conn.SetDeadline(time.Now().Add(time.Second))
	_, err := conn.Write([]byte(cmd))
	if err != nil {
		return err
	}

	return nil
}

func (q *Quit) Wait(conn *net.TCPConn) error {

	_ = conn.SetDeadline(time.Now().Add(time.Second))

	reader := bufio.NewReader(conn)
	data, _ := reader.ReadString('\n')
	q.reply = ParseReply(data)

	return nil
}

func (q Quit) Reply() Reply {
	return q.reply
}

func (q Quit) Name() string {
	return "QUIT"
}

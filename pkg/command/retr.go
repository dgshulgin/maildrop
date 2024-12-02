package command

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"strings"
	"time"
)

type Retrieve struct {
	reply     Reply
	message   Message
	requestId int // идентификатор запрашиваемого сообщения
}

func RetrieveMessageById(id int) *List {
	return &List{reply: Reply{}, requestId: id}
}

func (r Retrieve) Send(conn *net.TCPConn) error {

	cmd := fmt.Sprintf("%s %d%s", r.Name(), r.requestId, EOL)

	_ = conn.SetDeadline(time.Now().Add(time.Second))
	_, err := conn.Write([]byte(cmd))
	if err != nil {
		return err
	}

	return nil
}

/*
C: RETR 1
S: +OK 120 octets
S: <the POP3 server sends the entire message here>
S: .
*/
func (r *Retrieve) Wait(conn *net.TCPConn) error {

	_ = conn.SetDeadline(time.Now().Add(time.Second))
	reader := bufio.NewReader(conn)
	for {
		raw, _ := reader.ReadString('\n')
		if strings.HasPrefix(raw, PrefixERR) {
			return errors.New(raw) //TODO FIX
		}

		if strings.HasPrefix(raw, PrefixOK) {
			r.reply = ParseReply(raw) //убирает +OK
		}

		if strings.HasPrefix(raw, PrefixDot) {
			break
		}

		r.message.text = r.message.text + raw
	}
	r.message.id = r.requestId
	//TODO keep message size

	return nil
}

func (r Retrieve) Reply() Reply {
	return r.reply
}

func (r Retrieve) Name() string {
	return "RETR"
}

func (r Retrieve) MessageSize() int {
	return r.message.size
}

func (r Retrieve) MessageId() int {
	return r.message.id
}

func (r Retrieve) MessageText() string {
	return r.message.text
}

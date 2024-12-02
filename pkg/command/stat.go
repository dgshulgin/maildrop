package command

import (
	"bufio"
	"net"
	"strconv"
	"strings"
	"time"
)

type Stat struct {
	reply Reply
	count int //кол-во сообщений в ящике
	size  int //размер ящика
}

func SendStat() *Stat {
	return &Stat{reply: Reply{}}
}

func (s Stat) Send(conn *net.TCPConn) error {

	cmd := s.Name() + EOL

	_ = conn.SetDeadline(time.Now().Add(time.Second))
	_, err := conn.Write([]byte(cmd))
	if err != nil {
		return err
	}

	return nil
}

// +OK<space>nn<space>mm<EOL>
// -ERR message
func (s *Stat) Wait(conn *net.TCPConn) error {

	_ = conn.SetDeadline(time.Now().Add(time.Second))

	reader := bufio.NewReader(conn)
	data, _ := reader.ReadString('\n')
	s.reply = ParseReply(data)

	parts := strings.Split(s.reply.Message(), " ")
	s.count, _ = strconv.Atoi(parts[0])
	s.size, _ = strconv.Atoi(parts[1])

	return nil
}

func (s Stat) Reply() Reply {
	return s.reply
}

func (s Stat) Name() string {
	return "STAT"
}

func (s Stat) MessagesCount() int {
	return s.count
}

func (s Stat) MailboxSize() int {
	return s.size
}

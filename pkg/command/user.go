package command

import (
	"bufio"
	"net"
	"time"
)

type User struct {
	reply    Reply
	username string
}

func SendUser(username string) *User {
	return &User{username: username, reply: Reply{}}
}

// USER username\r\n
func (u User) Send(conn *net.TCPConn) error {

	cmd := u.Name() + SpaceDelim + u.username + EOL

	_ = conn.SetDeadline(time.Now().Add(time.Second))
	_, err := conn.Write([]byte(cmd))
	if err != nil {
		return err
	}

	return nil
}

// +OK message
// -ERR message
func (u *User) Wait(conn *net.TCPConn) error {

	_ = conn.SetDeadline(time.Now().Add(time.Second))

	reader := bufio.NewReader(conn)

	data, _ := reader.ReadString('\n')
	u.reply = ParseReply(data)

	return nil
}

func (u User) Reply() Reply {
	return u.reply
}

func (u User) Name() string {
	return "USER"
}

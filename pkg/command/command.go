package command

import (
	"errors"
	"net"
	"strings"
)

/*
Minimal
USER name
PASS string
QUIT

STAT
LIST [msg]
TODO RETR msg
TODO DELE msg
TODO NOOP
TODO RSET
QUIT

Optional
TODO APOP name digest
TODO TOP msg n
TODO UIDL [msg]


*/

const (
	SpaceDelim = " "
	EOL        = "\r\n"
	PrefixOK   = "+OK"
	PrefixERR  = "-ERR"
	PrefixDot  = "."
)

type Command interface {
	Send(conn *net.TCPConn) error
	Wait(conn *net.TCPConn) error
	Reply() Reply
	Name() string
}

func Serve(conn *net.TCPConn, cmd Command) (Command, error) {
	err := cmd.Send(conn)
	if err != nil {
		return nil, errors.Join(errors.New(cmd.Name()), err)
	}
	err = cmd.Wait(conn)
	if err != nil {
		return nil, errors.Join(errors.New(cmd.Name()), err)
	}
	return cmd, nil
}

type Reply struct {
	ok      bool
	message string
}

func (r Reply) IsOk() bool {
	return r.ok
}

func (r Reply) Message() string {
	return r.message
}

func ParseReply(text string) Reply {
	text = strings.TrimSuffix(text, EOL)

	msg, ok := strings.CutPrefix(text, PrefixOK)
	if ok {
		msg = strings.TrimSpace(msg)
		return Reply{ok: ok, message: msg}
	}

	return Reply{ok: false, message: text}
}

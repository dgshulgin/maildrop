package command

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

/*
LIST [msg]

Arguments:
Опциональный номер сообщения. Неразрешается указывать номер сообщения, помеченного как удаленное.

Restrictions:
Только после аутентификации пользователя.

Discussion:
Для заданного аргумента сервер возвращает строку с информацией о сообщении.
C: LIST 2
S: +OK 2 200
...

Если аргумент не задан то ответ сервера представлен несколькими строками.
1. OK+ nn mm
2. информация о сообщении вида "1 120", первое число идентификатор, второе - размер сообщения.

Если ящик пустой, сервер возвращает +OK 0 0

*/

type List struct {
	reply   Reply
	mailbox []Message

	// Значение id=0 означает запрос без параметров т.е. LIST<EOL>
	requestId int // для запроса сообщения с определенным идентификатором
}

func ListMessages() *List {
	return &List{reply: Reply{}, requestId: 0}
}

func ListMessageById(id int) *List {
	return &List{reply: Reply{}, requestId: id}
}

func (l *List) Send(conn *net.TCPConn) error {

	var cmd string
	if l.requestId > 0 {
		cmd = fmt.Sprintf("%s %d\r\n", l.Name(), l.requestId)
	} else {
		cmd = fmt.Sprintf("%s\r\n", l.Name())
	}

	_ = conn.SetDeadline(time.Now().Add(time.Second))
	_, err := conn.Write([]byte(cmd))
	if err != nil {
		return err
	}

	return nil
}

// S: +OK<space>nn<space>messages<space>(nn<space>octets)<EOL>
// S: 1 200
// S: 2 120
// S: .
// точка означает конец сообщения
// -ERR message
func (l *List) Wait(conn *net.TCPConn) error {

	_ = conn.SetDeadline(time.Now().Add(time.Second))
	reader := bufio.NewReader(conn)
	raw, _ := reader.ReadString('.')
	l.reply = ParseReply(raw) //убирает +OK

	if l.requestId == 0 {
		// заголовок nn mm
		replyReader := bufio.NewReader(strings.NewReader(l.reply.Message()))
		headerRaw, _ := replyReader.ReadString('\n')
		headerParts := strings.Split(headerRaw, " ")
		totalMessages, _ := strconv.Atoi(headerParts[0])
		//mbSize, _ := strconv.Atoi(headerParts[1])

		for i := 1; i <= totalMessages; i++ {
			msgRaw, _ := replyReader.ReadString('\n')

			msgParts := strings.Split(msgRaw, " ")

			msgId, _ := strconv.Atoi(msgParts[0])
			msgSize, _ := strconv.Atoi(msgParts[1])

			l.mailbox = append(l.mailbox, Message{id: msgId, size: msgSize})
		}
	} else {
		replyReader := bufio.NewReader(strings.NewReader(l.reply.Message()))
		msgRaw, _ := replyReader.ReadString('\n')

		msgParts := strings.Split(msgRaw, " ")

		msgId, _ := strconv.Atoi(msgParts[0])
		msgSize, _ := strconv.Atoi(msgParts[1])

		l.mailbox = append(l.mailbox, Message{id: msgId, size: msgSize})

	}

	return nil
}

func (l List) Reply() Reply {
	return l.reply
}

func (l List) Name() string {
	return "LIST"
}

func (l List) Mailbox() []Message {
	return l.mailbox
}

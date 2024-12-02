package pop3test

import (
	"maildrop/pkg/command"
	"maildrop/pkg/network"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type Pop3TestSuite struct {
	suite.Suite
	Server   string
	Port     string
	Username string
	Pass     string
}

// run suite
func TestPop3Suite(t *testing.T) {
	suite.Run(t, new(Pop3TestSuite))
}

func (suite *Pop3TestSuite) SetupTest() {

	if err := godotenv.Load("../.env"); err != nil {
		suite.T().Fatalf("Fatal: .env not found; %s\n", err.Error())
	}

	var ok bool
	suite.Server, ok = os.LookupEnv("POP3_SERVER")
	if !ok {
		suite.T().Fatalf("POP3 server undefined")
	}
	suite.Port, ok = os.LookupEnv("POP3_PORT")
	if !ok {
		suite.T().Fatalf("POP3 port undefined")
	}

	suite.Username, ok = os.LookupEnv("USER")
	if !ok {
		suite.T().Fatalf("логин не определен")
	}

	suite.Pass, ok = os.LookupEnv("PASSWORD")
	if !ok {
		suite.T().Fatalf("пароль не определен")
	}
}

// Проверяет правильность обработки ответа на приветствие сервера, которое поступает
// непосредственно после установки соедения.
func (suite *Pop3TestSuite) TestPOP3Greeting() {

	tcpConn, err := network.Connect(suite.Server, suite.Port)
	if err != nil {
		suite.T().Fatalf("%s\n", err.Error())
	}
	defer tcpConn.Close()

	g, err := command.Serve(tcpConn, command.ServerGreeting())
	if err != nil {
		suite.T().Fatalf("%s\n", err.Error())
	}
	assert.Truef(suite.T(), g.Reply().IsOk(), "server greeting failed")
}

// func TestPOP3ServerUnavailable(t *testing.T) {

// 	if err := godotenv.Load("../.env"); err != nil {
// 		t.Fatalf("Fatal: .env not found; %s\n", err.Error())
// 	}

// 	server, ok := os.LookupEnv("BAD_POP3_SERVER")
// 	if !ok {
// 		t.Fatalf("POP3 server undefined\n")
// 	}
// 	port, ok := os.LookupEnv("BAD_POP3_PORT")
// 	if !ok {
// 		t.Fatalf("POP3 port undefined\n")
// 	}

// 	tcpConn, err := network.Connect(server, port)
// 	if err != nil {
// 		t.Fatalf("%s\n", err.Error())
// 	}
// 	defer tcpConn.Close()

// 	g, err := command.Serve(tcpConn, command.ServerGreeting())
// 	if err != nil {
// 		t.Fatalf("%s\n", err.Error())
// 	}
// 	require.Truef(t, g.Reply().IsOk(), "server greeting failed")
// }

// Проверяет корректность авторизации с помощью пары команд USER/PASS.
func (suite *Pop3TestSuite) TestPOP3AuthUserPass() {

	tcpConn, err := network.Connect(suite.Server, suite.Port)
	require.NoErrorf(suite.T(), err, "connect failed")
	defer tcpConn.Close()

	g, err := command.Serve(tcpConn, command.ServerGreeting())
	assert.NoErrorf(suite.T(), err, "server greeting failed")
	assert.Truef(suite.T(), g.Reply().IsOk(), "server greeting failed")

	u, err := command.Serve(tcpConn, command.SendUser(suite.Username))
	assert.NoErrorf(suite.T(), err, "USER failed")
	assert.Truef(suite.T(), u.Reply().IsOk(), u.Reply().Message())

	p, err := command.Serve(tcpConn, command.SendPass(suite.Pass))
	assert.NoErrorf(suite.T(), err, "PASS failed")
	assert.Truef(suite.T(), p.Reply().IsOk(), p.Reply().Message())
}

// TODO Этот тест работает некорректрно, серер может вернуть OK даже если ящика нет, и вернуть ERR есть но не позволяет авторизацию через plain text.
// Проверяет корректность поведения клиентского АПИ при авторизации с неправильным именем пользователя, оно же названиe.
// В этом случае сервер возвращает -ERR login failed.
func (suite *Pop3TestSuite) TestPOP3AuthBadUser() {

	var ok bool
	suite.Username, ok = os.LookupEnv("BAD_USER")
	require.Truef(suite.T(), ok, "логин не определен")

	tcpConn, err := network.Connect(suite.Server, suite.Port)
	require.NoErrorf(suite.T(), err, "connect failed")
	defer tcpConn.Close()

	g, err := command.Serve(tcpConn, command.ServerGreeting())
	assert.NoErrorf(suite.T(), err, "server greeting failed")
	assert.Truef(suite.T(), g.Reply().IsOk(), "server greeting failed")

	u, err := command.Serve(tcpConn, command.SendUser(suite.Username))
	assert.NoErrorf(suite.T(), err, "USER failed")
	assert.Truef(suite.T(), u.Reply().IsOk(), u.Reply().Message())

	// команду PASS не посылаем
}

// Проверяет корректность поведения клиентского АПИ при авторизации с неправильным паролем.
// В этом случае сервер возвращает -ERR login failed
func (suite *Pop3TestSuite) TestPOP3AuthBadPass() {

	var ok bool
	suite.Pass, ok = os.LookupEnv("BAD_PASSWORD")
	require.Truef(suite.T(), ok, "пароль не определен")

	tcpConn, err := network.Connect(suite.Server, suite.Port)
	require.NoErrorf(suite.T(), err, "connect failed")
	defer tcpConn.Close()

	g, err := command.Serve(tcpConn, command.ServerGreeting())
	assert.NoErrorf(suite.T(), err, "server greeting failed")
	assert.Truef(suite.T(), g.Reply().IsOk(), "server greeting failed")

	u, err := command.Serve(tcpConn, command.SendUser(suite.Username))
	assert.NoErrorf(suite.T(), err, "USER failed")
	assert.Truef(suite.T(), u.Reply().IsOk(), u.Reply().Message())

	p, err := command.Serve(tcpConn, command.SendPass(suite.Pass))
	assert.NoErrorf(suite.T(), err, "PASS failed")
	assert.Truef(suite.T(), p.Reply().IsOk(), p.Reply().Message())
}

// Проверяет корректность поведение клиентского АПИ при запросе списка всех сообщений в ящике, вызов LIST без параметров
func (suite *Pop3TestSuite) TestPOP3ListAllMessages() {

	tcpConn, err := network.Connect(suite.Server, suite.Port)
	require.NoErrorf(suite.T(), err, "connect failed")
	defer tcpConn.Close()

	g, err := command.Serve(tcpConn, command.ServerGreeting())
	assert.NoErrorf(suite.T(), err, "server greeting failed")
	assert.Truef(suite.T(), g.Reply().IsOk(), "server greeting failed")

	u, err := command.Serve(tcpConn, command.SendUser(suite.Username))
	assert.NoErrorf(suite.T(), err, "USER failed")
	assert.Truef(suite.T(), u.Reply().IsOk(), u.Reply().Message())

	p, err := command.Serve(tcpConn, command.SendPass(suite.Pass))
	assert.NoErrorf(suite.T(), err, "PASS failed")
	assert.Truef(suite.T(), p.Reply().IsOk(), p.Reply().Message())

	list, err := command.Serve(tcpConn, command.ListMessages())
	assert.NoErrorf(suite.T(), err, "LIST failed")
	assert.Truef(suite.T(), list.Reply().IsOk(), list.Reply().Message())
	assert.Truef(suite.T(), len(list.(*command.List).Mailbox()) > 0, "mailbox shall not be empty")

}

// Проверяет корректность поведение клиентского АПИ при запросе сообщения с определенным индексом, гарантированно находящимся в ящике.
// Команда LIST<space>Num<EOL>, где Num порядковый номер сообщения начиная с единицы.
func (suite *Pop3TestSuite) TestPOP3ListById() {

	tcpConn, err := network.Connect(suite.Server, suite.Port)
	require.NoErrorf(suite.T(), err, "connect failed")
	defer tcpConn.Close()

	g, err := command.Serve(tcpConn, command.ServerGreeting())
	assert.NoErrorf(suite.T(), err, "server greeting failed")
	assert.Truef(suite.T(), g.Reply().IsOk(), "server greeting failed")

	u, err := command.Serve(tcpConn, command.SendUser(suite.Username))
	assert.NoErrorf(suite.T(), err, "USER failed")
	assert.Truef(suite.T(), u.Reply().IsOk(), u.Reply().Message())

	p, err := command.Serve(tcpConn, command.SendPass(suite.Pass))
	assert.NoErrorf(suite.T(), err, "PASS failed")
	assert.Truef(suite.T(), p.Reply().IsOk(), p.Reply().Message())

	list, err := command.Serve(tcpConn, command.ListMessageById(1))
	assert.NoErrorf(suite.T(), err, "LIST failed")
	assert.Truef(suite.T(), list.Reply().IsOk(), list.Reply().Message())
	assert.Truef(suite.T(), len(list.(*command.List).Mailbox()) == 1, "mailbox shall contain a single message only")
	//FIX

}

package command

type Message struct {
	id   int // порядковый номер сообщения, начиная с 1
	size int // размер сообщения, в байтах
	text string
}

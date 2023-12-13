package chat

import (
	"fmt"

	"github.com/gorilla/websocket"
)

type ChatSession struct {
	user string
	peer *websocket.Conn
}

var (
	Peers map[string]*websocket.Conn
)

const (
	chat = "%s: %s"
	left = "%s: has left the chat"
)

func init() {
	Peers = map[string]*websocket.Conn{}
}

func NewChatSession(user string, peer *websocket.Conn) *ChatSession {
	return &ChatSession{
		user: user,
		peer: peer,
	}
}

func (s *ChatSession) Start() {
	Peers[s.user] = s.peer

	go func(){
		for {
			_, msg, err := s.peer.ReadMessage() // blocking i/o: prevent infinite process
			if err != nil {
				_, ok := err.(*websocket.CloseError)
				if ok {
					s.disconnect()
				}
				return
			}

			SendToChannel(fmt.Sprintf(chat, s.user, string(msg)))
		}
	}()
}

func (s *ChatSession) disconnect() {
	RemoveUser(s.user)

	SendToChannel(fmt.Sprintf(left, s.user))
	s.peer.Close()

	delete(Peers, s.user)
}
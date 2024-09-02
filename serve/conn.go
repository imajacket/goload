package serve

import (
	"log"
	"net/http"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type connection struct {
	id     string
	conn   *websocket.Conn
	reload chan struct{}
	close  chan struct{}
}

type WsHelper struct {
	upgrader  *websocket.Upgrader
	customers map[string]*connection
	sync.RWMutex
}

// newWsHelper
// Returns the WsHelper in charge of handling websocket connections and customers.
func newWsHelper() *WsHelper {
	return &WsHelper{
		upgrader: &websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		customers: make(map[string]*connection),
	}
}

// wsHandler
// Handles websocket connections and creates corresponding customers.
func (s *WsHelper) wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	id := uuid.New().String()
	c := &connection{
		id:     id,
		conn:   conn,
		reload: make(chan struct{}),
		close:  make(chan struct{}),
	}

	s.Lock()
	s.customers[id] = c
	s.Unlock()

	defer func() {
		s.Lock()
		_ = conn.Close()
		delete(s.customers, id)
		log.Println("ended connection with", id)
		s.Unlock()
	}()

	go func() {
		for {
			msgType, _, rErr := conn.ReadMessage()
			if rErr != nil {
				c.close <- struct{}{}
			}
			if msgType != websocket.TextMessage {
				c.close <- struct{}{}
			}
		}
	}()

	log.Println("started connection with", id)
	for {
		select {
		case <-c.reload:
			wErr := conn.WriteMessage(websocket.TextMessage, []byte("reload"))
			if wErr != nil {
				return
			}
		case <-c.close:
			return
		}
	}
}

// reloadHandler
// Tells all customers to reload their page.
func (s *WsHelper) reloadHandler(w http.ResponseWriter, r *http.Request) {
	s.Lock()
	defer s.Unlock()

	for _, c := range s.customers {
		log.Println(c.id)
		c.reload <- struct{}{}
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("updating customers"))
}

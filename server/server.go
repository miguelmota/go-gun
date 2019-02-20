package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	spew "github.com/davecgh/go-spew/spew"
	ws "github.com/gorilla/websocket"
	common "github.com/miguelmota/go-gun/common"
	log "github.com/sirupsen/logrus"
)

// ErrServerAlreadyStarted ...
var ErrServerAlreadyStarted = errors.New("server already started")

// Server is the server structure
type Server struct {
	host    string
	port    uint
	peers   []*ws.Conn
	graph   map[string]interface{}
	dup     *common.Dup
	server  *http.Server
	started bool
	debug   bool
}

// Config is the server config
type Config struct {
	Port  uint
	Debug bool
}

// DefaultPort is the default port
var DefaultPort uint = 8080

var upgrader = ws.Upgrader{
	HandshakeTimeout:  0,
	ReadBufferSize:    1024,
	WriteBufferSize:   1024,
	Subprotocols:      nil,
	Error:             nil,
	CheckOrigin:       checkOrigin,
	EnableCompression: false,
}

func checkOrigin(r *http.Request) bool {
	return true
}

// NewServer returns a new server instance
func NewServer(config *Config) *Server {
	if config == nil {
		config = &Config{}
	}

	port := DefaultPort
	if config.Port != 0 {
		port = config.Port
	}

	return &Server{
		host:  fmt.Sprintf("0.0.0.0:%v", port),
		port:  port,
		peers: []*ws.Conn{},
		graph: make(map[string]interface{}),
		dup:   common.NewDup(),
		debug: config.Debug,
	}
}

// Start starts the websocket server
func (s *Server) Start() error {
	if s.started {
		return ErrServerAlreadyStarted
	}

	s.started = true
	srv := &http.Server{Addr: s.host}
	s.server = srv
	http.HandleFunc("/", s.RequestHandler)

	fmt.Printf("listening on %s\n", s.host)
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}

	return nil
}

// Stop stops the websocket server
func (s *Server) Stop() error {
	return s.server.Shutdown(context.Background())
}

// RequestHandler is the handler for incoming connections
func (s *Server) RequestHandler(w http.ResponseWriter, r *http.Request) {
	peer, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}

	s.peers = append(s.peers, peer)

	for {
		// read message from browser
		_, msg, err := peer.ReadMessage()
		if err != nil {
			log.Fatal(err)
		}

		var js map[string]interface{}
		err = json.Unmarshal(msg, &js)
		if err != nil {
			log.Fatal(err)
		}

		soul := js["#"].(string)

		if s.dup.Check(soul) {
			continue
		}

		s.dup.Track(soul)
		fmt.Printf("received: %s\n", js)
		fmt.Printf("from: %s\n", peer.RemoteAddr())

		var resp []byte
		if change, ok := js["put"]; ok {
			diff := common.Mix(change.(map[string]interface{}), s.graph)
			_ = diff
			//fmt.Println("diff", diff)

			uid := s.dup.Track(common.NewUID())
			s.logGraph()

			resp, err = json.Marshal(map[string]interface{}{
				"#": uid,
				"@": soul,
			})
			if err != nil {
				log.Fatal(err)
			}
		} else if getValue, ok := js["get"]; ok {
			ack := common.Get(getValue.(map[string]interface{}), s.graph)
			if ack != nil {
				uid := s.dup.Track(common.NewUID())
				resp, err = json.Marshal(map[string]interface{}{
					"#":   uid,
					"@":   soul,
					"put": ack,
				})
				if err != nil {
					log.Fatal(err)
				}
			}
		}

		s.logGraph()
		if err := s.Emit(resp); err != nil {
			log.Fatal(err)
		}

		if err := s.Emit(msg); err != nil {
			log.Fatal(err)
		}
	}
}

// Emit emits message to all connected peers
func (s *Server) Emit(msg []byte) error {
	return emit(s.peers, msg)
}

// RemovePeer removes a peer from the peer list
func (s *Server) RemovePeer(peer *ws.Conn) error {
	for i, pn := range s.peers {
		if pn == peer {
			s.peers = append(s.peers[:i], s.peers[i+1:]...)
			break
		}
	}

	return nil
}

// logGraph logs the graph structure
func (s *Server) logGraph() {
	spew.Dump(s.graph)
}

// emit emits message to the peer list
func emit(peers []*ws.Conn, msg []byte) error {
	for _, peer := range peers {
		if err := peer.WriteMessage(ws.TextMessage, msg); err != nil {
			return err
		}
	}

	return nil
}

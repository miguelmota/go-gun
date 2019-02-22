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
	netutil "github.com/miguelmota/go-gun/netutil"
	storage "github.com/miguelmota/go-gun/storage"
	types "github.com/miguelmota/go-gun/types"
	log "github.com/sirupsen/logrus"
)

// ErrServerAlreadyStarted ...
var ErrServerAlreadyStarted = errors.New("server already started")

// Server is the server structure
type Server struct {
	host    string
	port    uint
	peers   map[string]*ws.Conn
	graph   types.Kv
	dup     *common.Dup
	server  *http.Server
	started bool
	debug   bool
	storage storage.IStorage
}

// Config is the server config
type Config struct {
	Port  *uint
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
	if config.Port != nil {
		port = *config.Port
		if port == 0 {
			p, err := netutil.GetFreePort()
			if err != nil {
				log.Fatal(err)
			}
			port = uint(p)
		}
	}

	graph := make(types.Kv)

	return &Server{
		host:    fmt.Sprintf("0.0.0.0:%v", port),
		port:    port,
		peers:   make(map[string]*ws.Conn),
		graph:   graph,
		dup:     common.NewDup(),
		debug:   config.Debug,
		storage: storage.NewDummyKV(graph),
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

	log.Printf("listening on %s\n", s.host)
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
		log.Error(err)
		return
	}

	s.peers[peer.RemoteAddr().String()] = peer

	for {
		// read message from browser
		_, msg, err := peer.ReadMessage()
		if msg == nil {
			return
		}
		if ws.IsCloseError(err) {
			log.Error(err)
			s.RemovePeer(peer)
			return
		}
		if err != nil {
			//log.Error(err)
			return
		}

		var js types.Kv
		err = json.Unmarshal(msg, &js)
		if err != nil {
			//log.Error(err)
			continue
		}

		if js == nil {
			continue
		}

		soul, ok := js["#"].(string)
		if !ok {
			continue
		}

		if s.dup.Check(soul) {
			continue
		}

		s.dup.Track(soul)

		var resp []byte
		if change, ok := js["put"]; ok {
			diff := common.Mix(types.Kv(change.(map[string]interface{})), s.graph)
			uid := s.dup.Track(common.NewUID())

			resp, err = json.Marshal(types.Kv{
				"#": uid,
				"@": soul,
			})
			if err != nil {
				log.Error(err)
				continue
			}

			for soul, node := range diff {
				for k, v := range node.(types.Kv) {
					if k == "_" {
						continue
					}

					kstate := diff[soul].(types.Kv)["_"].(types.Kv)[">"].(types.Kv)[k]
					_ = soul
					_ = k
					_ = v
					_ = kstate
					s.storage.Put(soul, k, v, kstate)
				}
			}
		} else if getValue, ok := js["get"]; ok {
			ack := common.Get(common.IToKv(getValue), s.graph)
			if ack != nil {
				uid := s.dup.Track(common.NewUID())
				resp, err = json.Marshal(types.Kv{
					"#":   uid,
					"@":   soul,
					"put": ack,
				})

				if err != nil {
					log.Error(err)
					continue
				}
			}
		}

		if err := s.Emit(resp); err != nil {
			log.Error(err)
			continue
		}
	}
}

// Emit emits message to all connected peers
func (s *Server) Emit(msg []byte) error {
	for _, peer := range s.peers {
		if err := peer.WriteMessage(ws.TextMessage, msg); err != nil {
			if err := s.RemovePeer(peer); err != nil {
				log.Fatal(err)
			}
			continue
		}
	}

	return nil
}

// RemovePeer removes a peer from the peer list
func (s *Server) RemovePeer(peer *ws.Conn) error {
	delete(s.peers, peer.RemoteAddr().String())
	return nil
}

// logGraph logs the graph structure
func (s *Server) logGraph() {
	if s.debug {
		spew.Dump(s.graph)
	}
}

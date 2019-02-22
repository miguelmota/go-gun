package client

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/davecgh/go-spew/spew"
	ws "github.com/gorilla/websocket"
	common "github.com/miguelmota/go-gun/common"
	storage "github.com/miguelmota/go-gun/storage"
	types "github.com/miguelmota/go-gun/types"
	log "github.com/sirupsen/logrus"
)

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

// Client ...
type Client struct {
	wsendpoint string
	ws         *ws.Conn
	storage    storage.IStorage
	graph      types.Kv
}

// NewClient ...
func NewClient(wsendpoint string) *Client {
	ws, _, err := ws.DefaultDialer.Dial(wsendpoint, nil)
	if err != nil {
		log.Fatal(err)
	}

	graph := make(types.Kv)

	return &Client{
		wsendpoint: wsendpoint,
		ws:         ws,
		storage:    storage.NewDummyKV(graph),
		graph:      graph,
	}
}

// Send ...
func (c *Client) Send(msg interface{}) error {
	b, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	err = c.ws.WriteMessage(ws.TextMessage, b)
	if err != nil {
		return err
	}

	return nil
}

// Put ...
func (c *Client) Put(soul string, args ...types.Kv) {
	change := FormatPutRequest(soul, args...)

	c.Send(change)

	_, _, err := c.ws.ReadMessage()
	if err != nil {
		log.Fatal(err)
	}
}

// Get ...
func (c *Client) Get(soul string, key *string) types.Kv {
	change1 := FormatGetRequest(soul)

	c.Send(change1)

	var msg []byte
	var err error
	for {
		_, msg, err = c.ws.ReadMessage()
		if err != nil {
			log.Error(err)
			continue
		}
		if string(msg) == "" {
			continue
		}
		break
	}

	log.Printf("recv: %s", msg)

	var loaded map[string]interface{}
	err = json.Unmarshal(msg, &loaded)
	if err != nil {
		log.Fatal(err)
	}

	change := common.IToKv(loaded["put"])
	soul1 := loaded["#"].(string)
	diff := common.Mix(change, c.graph)

	resp := make(types.Kv)
	resp["@"] = soul1
	resp["#"] = common.NewUID()

	for sol, node := range diff {
		spew.Dump(node)
		for k, v := range node.(types.Kv) {
			if k == "_" {
				continue
			}

			kstate := common.GetStateOfProp(diff[sol].(types.Kv), k)
			c.storage.Put(sol, k, v, kstate)
		}
	}

	return c.storage.Get(soul, key)
}

// Close ...
func (c *Client) Close() {
	c.ws.Close()
}

// FormatPutRequest ...
func FormatPutRequest(soul string, args ...types.Kv) types.Kv {
	change := make(map[string]interface{})
	change["#"] = common.NewUID()
	change["put"] = make(types.Kv)
	change["put"].(types.Kv)[soul] = common.NewNode(soul, args...)

	js, _ := json.Marshal(change)
	fmt.Println("PUTTTT", string(js))

	return change
}

// FormatGetRequest ...
func FormatGetRequest(soul string) types.Kv {
	change := make(types.Kv)

	change["#"] = common.NewUID()
	change["get"] = make(types.Kv)
	change["get"].(types.Kv)["#"] = soul

	return change
}

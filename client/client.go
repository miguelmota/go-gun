package client

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/davecgh/go-spew/spew"
	ws "github.com/gorilla/websocket"
	common "github.com/miguelmota/go-gun/common"
	storage "github.com/miguelmota/go-gun/storage"
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
	graph      map[string]interface{}
}

// NewClient ...
func NewClient(wsendpoint string) *Client {
	ws, _, err := ws.DefaultDialer.Dial(wsendpoint, nil)
	if err != nil {
		log.Fatal(err)
	}

	return &Client{
		wsendpoint: wsendpoint,
		ws:         ws,
		storage:    storage.NewDummyKV(),
		graph:      make(map[string]interface{}),
	}
}

// Send ...
func (c *Client) Send(msg interface{}) error {
	b, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	fmt.Println("sending", string(b))

	err = c.ws.WriteMessage(ws.TextMessage, b)
	if err != nil {
		return err
	}

	/*
	   ch = format_put_request(soul, **kwargs)
	   ch_str = json.dumps(ch)
	   # print("Change: {} ".format(ch))
	   await ws.send(ch_str)
	   resp = await ws.recv()
	   # print("RESP: {} ".format(resp))
	   return resp
	*/

	return nil
}

// Put ...
func (c *Client) Put(soul string, args ...map[string]interface{}) {
	change := FormatPutRequest(soul, args...)
	fmt.Println("sending", change)

	c.Send(change)

	_, msg, err := c.ws.ReadMessage()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("put recv: %s", msg)
}

// Get ...
func (c *Client) Get(soul1 string, key *string) interface{} {
	change1 := FormatGetRequest(soul1)
	fmt.Println("getting", change1)

	c.Send(change1)

	_, msg, err := c.ws.ReadMessage()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("recv: %s", msg)

	var loaded map[string]interface{}
	err = json.Unmarshal(msg, &loaded)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("loaded")
	spew.Dump(loaded)

	change := make(map[string]interface{})
	if change1, ok := loaded["put"].(map[string]interface{}); ok {
		change = change1
	}
	soul := loaded["#"].(string)

	fmt.Println("LOAD SOUL", soul)

	fmt.Println("GAPH")
	fmt.Println("CHANGE", change)
	spew.Dump(c.graph)
	diff := common.Mix(change, c.graph)
	fmt.Println("GAPH")
	spew.Dump(c.graph)

	resp := make(map[string]interface{})
	resp["@"] = soul
	resp["#"] = common.NewUID()

	fmt.Println("DIFF", diff)
	spew.Dump(diff)

	_ = resp

	for sol, node := range diff {
		fmt.Println("NODE")
		spew.Dump(node)
		for k, v := range node.(map[string]interface{}) {
			if k == "_" {
				continue
			}
			fmt.Println("YOOO", k)

			kstate := diff[sol].(map[string]interface{})["_"].(map[string]interface{})[">"].(map[string]interface{})[k]

			fmt.Println("PUT", sol, k, v)

			c.storage.Put(sol, k, v, kstate)
		}
	}
	fmt.Println("GET", soul, key)

	return c.storage.Get(soul, key)

	/*
	   ch = format_get_request(soul)
	   ch_str = json.dumps(ch)
	   # print("Change: {} ".format(ch))
	   await ws.send(ch_str)
	   resp = await ws.recv()
	   loaded = json.loads(resp)
	   # print("RESP: {} ".format(resp))
	   change = loaded['put']
	   # print("CHANGE IS: ", change)
	   soul = loaded['#']
	   diff = ham_mix(change, self.backend)

	   resp = {'@':soul, '#':newuid(), 'ok':True}
	   # print("DIFF:", diff)

	   for soul, node in diff.items():
	       for k, v in node.items():
	           if k == "_":
	               continue
	           kstate = diff[soul]['_']['>'][k]
	           print("KSTATE: ", kstate)
	           self.backend.put(soul, k, v, kstate)
	   return self.backend.get(soul, key)
	*/
}

// FormatPutRequest ...
func FormatPutRequest(soul string, args ...map[string]interface{}) map[string]interface{} {
	change := make(map[string]interface{})
	change["#"] = common.NewUID()
	change["put"] = make(map[string]interface{})
	change["put"].(map[string]interface{})[soul] = common.NewNode(soul, args...)

	js, _ := json.Marshal(change)
	fmt.Println("PUTTTT", string(js))

	return change
}

// FormatGetRequest ...
func FormatGetRequest(soul string) map[string]interface{} {
	change := make(map[string]interface{})
	change["#"] = common.NewUID()
	change["get"] = make(map[string]interface{})
	change["get"].(map[string]interface{})["#"] = soul

	return change
}

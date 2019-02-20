package server

import (
	"testing"
	"time"

	log "github.com/sirupsen/logrus"
)

func TestServer(t *testing.T) {
	srv := NewServer(&Config{})
	go func() {
		err := srv.Start()
		if err != nil {
			log.Fatal(err)
		}
	}()

	time.Sleep(1 * time.Second)

	err := srv.Stop()
	if err != nil {
		t.Error(err)
	}
}

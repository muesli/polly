package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"time"
)

// Handler is a clever config parser
type Handler struct {
	destroyTimeoutHandler chan bool
	fileName              string
	lastModificationTime  time.Time
	lastReadValid         bool
	notificationChannel   chan interface{}
	rescanInterval        uint16
	ticker                *time.Ticker
	unmarshalStruct       interface{}
}

// NewHandler returns a new config handler
func NewHandler(fileName string, unmarshalStruct interface{}, notificationChannel chan interface{}) *Handler {
	if _, err := os.Stat(fileName); err != nil {
		return nil
	}
	ch := Handler{
		fileName:             fileName,
		lastModificationTime: time.Unix(0, 0),
	}
	if notificationChannel != nil {
		ch.notificationChannel = notificationChannel
	}
	ch.unmarshalStruct = unmarshalStruct
	ch.destroyTimeoutHandler = make(chan bool)
	ch.Rescan()
	ch.SetRescanInterval(3)
	return &ch
}

// SetRescanInterval controls how often the config should be rescanned
func (ch *Handler) SetRescanInterval(interval uint16) {
	if ch.ticker != nil {
		ch.ticker.Stop()
		ch.destroyTimeoutHandler <- true
	}

	if interval == 0 {
		ch.ticker = nil
		return
	}

	ch.ticker = time.NewTicker(time.Second * time.Duration(interval))
	go ch.handleTimeout()
}

// Rescan parses the config file again
func (ch *Handler) Rescan() {
	ch.lastReadValid = false
	var fi os.FileInfo
	var err error
	if fi, err = os.Stat(ch.fileName); err != nil {
		return
	}
	if fi.ModTime().Equal(ch.lastModificationTime) {
		return
	}
	ch.lastModificationTime = fi.ModTime()
	var contents []byte
	if contents, err = ioutil.ReadFile(ch.fileName); err != nil {
		return
	}
	if err = json.Unmarshal(contents, ch.unmarshalStruct); err != nil {
		return
	}
	ch.lastReadValid = true
	if ch.notificationChannel != nil {
		ch.notificationChannel <- ch.unmarshalStruct
	}
}

// LastReadValid returns whether the last parse succeeded
func (ch *Handler) LastReadValid() bool {
	return ch.lastReadValid
}

// CurrentData returns the current config data
func (ch *Handler) CurrentData() interface{} {
	return ch.unmarshalStruct
}

func (ch *Handler) handleTimeout() {
	for {
		select {
		case _ = <-ch.destroyTimeoutHandler:
			return
		case _ = <-ch.ticker.C:
			ch.Rescan()
		}
	}
}

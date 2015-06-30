package main

import (
	"bufio"
	"net"
	"strconv"
	"strings"

	log "github.com/Sirupsen/logrus"
)

const (
	DEFAULT_BIND = ":6543"
)

type Event struct {
	Key       string
	Timestamp int64
	Value     float64
}

func main() {

	log.SetLevel(log.DebugLevel)

	server, err := net.Listen("tcp", DEFAULT_BIND)
	if server == nil {
		log.Panic("couldn't start listening: %v", err)
	}
	conns := clientConns(server)
	for {
		go handleConn(<-conns)
	}
}

func clientConns(listener net.Listener) chan net.Conn {
	ch := make(chan net.Conn)
	i := 0
	go func() {
		for {
			client, err := listener.Accept()
			if client == nil {
				log.Info("couldn't accept: %v", err)
				continue
			}
			i++
			log.Infof("%d: %v <-> %v\n", i, client.LocalAddr(), client.RemoteAddr())
			ch <- client
		}
	}()
	return ch
}

func handleConn(client net.Conn) {
	b := bufio.NewReader(client)
	for {
		raw, err := b.ReadBytes('\n')
		if err != nil { // EOF, or worse
			break
		}
		line := string(raw[:])
		f := strings.Fields(line)

		t, err := strconv.ParseInt(f[1], 10, 0)
		if err != nil {
			log.Warnf("Failed to parse timestamp in \"%s\"", line)
		}
		v, err := strconv.ParseFloat(f[2], 64)
		if err != nil {
			log.Warnf("Failed to parse value in \"%s\"", line)
		}

		e := Event{Key: f[0], Timestamp: t, Value: v}

		log.Debug(e)
		log.Infof("key=>%s\ttimestamp=>%v\tvalue=>%v", e.Key, e.Timestamp, e.Value)
	}
}

package cmd

import (
	"bufio"
	"net"
	"strconv"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/gregbuehler/murmur/modules/setting"
)

type Event struct {
	Source    string
	Key       string
	Timestamp int64
	Value     float64
}

var CmdServer = cli.Command{
	Name:        "server",
	Usage:       "Start a murmur server",
	Description: "Starts a murmur server instance",
	Action:      runServer,
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "verbose, v",
			Usage: "Verbose output",
		},
		cli.StringFlag{
			Name:   "host",
			Value:  ":6543",
			Usage:  "Endpoint for listening",
			EnvVar: "MURMUR_HOST",
		},
		cli.StringFlag{
			Name:   "dbpath, p",
			Value:  "/var/murmur/",
			Usage:  "Database location",
			EnvVar: "MURMUR_DBPATH",
		},
		cli.IntFlag{
			Name:   "interval, i",
			Value:  5000,
			Usage:  "Recording interval (in milliseconds)",
			EnvVar: "MURMUR_INTERVAL",
		},
		cli.Float64Flag{
			Name:   "deadband, d",
			Value:  0.1,
			Usage:  "value delta threshold ",
			EnvVar: "MURMUR_DEADBAND",
		},
	},
}

func runServer(ctx *cli.Context) {
	setting.Verbose = ctx.Bool("verbose")
	setting.Host = ctx.String("host")
	setting.DbPath = ctx.String("dbpath")
	setting.Interval = int64(ctx.Int("interval"))
	setting.Deadband = ctx.Float64("deadband")

	if setting.Verbose {
		log.SetLevel(log.DebugLevel)
	}

	log.Infof("Starting server on %s", setting.Host)
	server, err := net.Listen("tcp", setting.Host)
	if server == nil {
		log.Panicf("couldn't start listening: %v", err)
	}
	conns := clientConns(server)
	for {
		go handleConn(<-conns)
	}
}

func clientConns(listener net.Listener) chan net.Conn {
	log.Debugf("Initializing client connections")
	ch := make(chan net.Conn)
	i := 0
	go func() {
		for {
			client, err := listener.Accept()
			if client == nil {
				log.Debugf("couldn't accept: %v", err)
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
	log.Debugf("Handling client connection")
	b := bufio.NewReader(client)
	for {
		raw, err := b.ReadBytes('\n')
		if err != nil { // EOF, or worse
			break
		}
		line := string(raw[:])
		f := strings.Fields(line)
		if len(f) != 4 {
			log.Errorf("Failed to parse event: %s", line)
			continue
		} else {
			t, err := strconv.ParseInt(f[2], 10, 0)
			if err != nil {
				log.Warnf("Failed to parse timestamp: %s", line)
				continue
			}
			v, err := strconv.ParseFloat(f[3], 64)
			if err != nil {
				log.Warnf("Failed to parse value in: %s", line)
				continue
			}

			e := Event{Source: f[0], Key: f[1], Timestamp: t, Value: v}

			log.Debug(e)
			log.Infof("source=>%s\tkey=>%s\ttimestamp=>%v\tvalue=>%v", e.Source, e.Key, e.Timestamp, e.Value)
		}
	}
}

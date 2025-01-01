package main

import (
	"bufio"
	"flag"
	"log"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "localhost:8080", "http service address")
var path = flag.String("path", "/echo", "ws route")

func main() {

	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: *addr, Path: *path}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("recv: %s", message)
		}
	}()

	var ticker *time.Ticker

	stdin := make(chan string)
	scanner := bufio.NewScanner(os.Stdin)

	go func() {
		for scanner.Scan() {
			stdin <- scanner.Text()
		}
		if scanner.Err() != nil {
			log.Println(err)
		}
	}()

	closeTicker := func() {
		if ticker != nil {
			ticker.Stop()
			ticker = nil

		}
	}

	startTicker := func(millis int) {
		if ticker == nil {
			ticker = time.NewTicker(time.Millisecond * time.Duration(millis))
		}
	}

	for {
		select {
		case <-done:
			return
		case str := <-stdin:
			if args := strings.Split(str, " "); len(args) == 2 && args[0] == "-ticker" {
				dur, err := strconv.Atoi(args[1])
				if err != nil {
					log.Println(err)
					continue
				}
				closeTicker()
				startTicker(dur)
				continue
			}
			err := c.WriteMessage(websocket.TextMessage, []byte(str))
			if err != nil {
				log.Println("write:", err)
				return
			}
		case <-interrupt:
			log.Println("interrupt")
			if ticker != nil {
				closeTicker()
				continue
			}
			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
				log.Println("connection was forcefully closed")
			}
			return
		default:
			if ticker == nil {
				continue
			}
			if t, ok := <-ticker.C; ok {
				{
					err := c.WriteMessage(websocket.TextMessage, []byte(t.String()))
					if err != nil {
						log.Println("write:", err)
						return
					}
				}
			}
		}
	}
}

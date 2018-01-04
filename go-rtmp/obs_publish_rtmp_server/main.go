package main

import (
	"fmt"
	"log"
	"sync"

	"github.com/nareix/joy4/av/avutil"
	"github.com/nareix/joy4/av/pubsub"
	"github.com/nareix/joy4/format"
	"github.com/nareix/joy4/format/rtmp"
)

func init() {
	format.RegisterAll()
}

type Channel struct {
	que *pubsub.Queue
}

func main() {
	rtmp.Debug = false
	server := &rtmp.Server{}
	l := &sync.RWMutex{}
	channels := map[string]*Channel{}
	server.HandlePlay = func(conn *rtmp.Conn) {
		log.Println("server.HandlePlay")
		l.RLock()
		ch := channels[conn.URL.Path]
		l.RUnlock()
		if ch != nil {
			cursor := ch.que.Latest()
			avutil.CopyFile(conn, cursor)
		}
	}

	server.HandlePublish = func(conn *rtmp.Conn) {
		log.Println("server.HandlePublish")
		streams, _ := conn.Streams()
		l.Lock()
		ch := channels[conn.URL.Path]
		fmt.Println("=====================")
		log.Println(conn.URL.Path)
		fmt.Println("=====================")
		if ch == nil {
			ch = &Channel{}
			ch.que = pubsub.NewQueue()
			ch.que.WriteHeader(streams)
			channels[conn.URL.Path] = ch
		} else {
			ch = nil
		}
		l.Unlock()
		if ch == nil {
			return
		}
		avutil.CopyPackets(ch.que, conn)
		l.Lock()
		delete(channels, conn.URL.Path)
		l.Unlock()
		ch.que.Close()
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf(err.Error())
	}
}

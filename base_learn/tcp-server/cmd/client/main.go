package main

import (
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/lucasepe/codename"
	"github.com/yuyuyue/tcp-server-demo1/frame"
	"github.com/yuyuyue/tcp-server-demo1/packet"
)

func main() {
	var wg sync.WaitGroup
	var num int = 5

	wg.Add(num)

	for i := 0; i < num; i++ {
		go func(i int) {
			defer wg.Done()
			startClient(i)
		}(i + 1)
	}
	wg.Wait()
}

func startClient(i int) {
	quit := make(chan struct{})
	done := make(chan struct{})
	conn, err := net.Dial("tcp", ":8888")
	if err != nil {
		fmt.Println("client conn error", err)
	}
	defer conn.Close()
	fmt.Printf("[client %d]: dial ok", i)
	rng, err := codename.DefaultRNG()
	if err != nil {
		panic(err)
	}
	frameStream := frame.NewFrameStream()
	var counter int
	go func() {
		// handle ack
		for {
			select {
			case <-quit:
				done <- struct{}{}
			default:
			}
			conn.SetDeadline(time.Now().Add(5 * time.Second))
			ack, err := frameStream.Decode(conn)

			if err != nil {
				if e, ok := err.(net.Error); ok {
					if e.Timeout() {
						continue
					}
					panic(err)
				}
			}
			p, err := packet.Decode(ack)
			submitAck, ok := p.(*packet.SubmitAck)
			if !ok {
				panic("not submitack")
			}
			fmt.Printf("[client %d]: the result of submit ack[%s] is %d\n", i, submitAck.ID, submitAck.Result)
		}
	}()
	for {
		// handle
		counter++
		id := fmt.Sprintf("%08d", counter)
		payload := codename.Generate(rng, 4)

		s := &packet.Submit{
			ID:      id,
			Payload: []byte(payload),
		}
		framePayload, err := packet.Encode(s)
		if err != nil {
			panic(err)
		}
		fmt.Printf("[client %d]: send submit id = %s, payload=%s, frame length = %d\n",
			i, s.ID, s.Payload, len(framePayload)+4)
		err = frameStream.Encode(conn, framePayload)
		if err != nil {
			panic(err)
		}

		time.Sleep(1 * time.Second)
		if counter >= 10 {
			quit <- struct{}{}
			<-done
			fmt.Printf("[client %d]: exit ok\n", i)
			return
		}
	}
}

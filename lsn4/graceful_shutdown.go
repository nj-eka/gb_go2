//Пример программы, которая при получении в канал сигнала SIGTERM останавливается
//не позднее, чем за одну секунду (установить таймаут).

package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

const (
	workTimeout      = 3 * time.Second
	workHeartbeat    = workTimeout / 10
	interruptTimeout = 1 * time.Second
	waitExitTimeout  = 1 * time.Second
)

func init() {
	log.SetFlags(log.Lmicroseconds | log.Lshortfile)
	log.SetOutput(os.Stdout)
}

func doWorkRight(parentCtx context.Context, workName string) context.Context {
	ctx, cancel := context.WithTimeout(parentCtx, workTimeout)
	go func() {
		defer cancel()
		workTicker := time.NewTicker(workHeartbeat)
		defer workTicker.Stop()
		for {
			select {
			case <-ctx.Done():
				log.Printf("%s stop working gracefully: %v\n", workName, ctx.Err())
				return
			case <-workTicker.C:
				log.Println(workName, "is working...")
			}
		}
	}()
	return ctx
}

func doWorkWrong(workName string) <-chan struct{} {
	done := make(chan struct{})
	go func() {
		log.Println(workName, workTimeout)
		defer close(done)
		workTicker := time.NewTicker(workHeartbeat)
		workStop := time.After(workTimeout)
		defer workTicker.Stop()
		for {
			select {
			case <-workStop:
				log.Printf("%s stop working: all done\n", workName)
				return
			case <-workTicker.C:
				log.Println(workName, "is working...")
			}
		}
	}()
	return done
}

func allClosed(streams ...<-chan struct{}) <-chan struct{} {
	result := make(chan struct{})
	go func() {
		defer close(result)
		wg := sync.WaitGroup{}
		wg.Add(len(streams))
		defer wg.Wait()
		for _, stream := range streams {
			go func(stream <-chan struct{}) {
				defer wg.Done()
				for range stream {
				}
			}(stream)
		}
	}()
	return result
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	log.Println("start works")
	w1Ctx, w2Ctx := doWorkRight(ctx, "goody-1"), doWorkRight(ctx, "goody-2")
	w1Done, w2Done := doWorkWrong("poor-1"), doWorkWrong("poor-2")
	interruptTimer := time.AfterFunc(interruptTimeout, func() {
		log.Println("SIGINT -> ")
		if err := syscall.Kill(os.Getpid(), syscall.SIGINT); err != nil {
			log.Println("error: %v", err)
		}
	})
	defer interruptTimer.Stop()

	select {
	case <-ctx.Done():
		log.Println("main stopping:", ctx.Err()) // prints "context canceled"
		stop()                                   // stop receiving signal notifications as soon as possible.
		<-time.After(waitExitTimeout)
		log.Println("Timeout emergency exit fired.")
		os.Exit(1)
	case <-allClosed(w1Ctx.Done(), w2Ctx.Done(), w1Done, w2Done):
		log.Println("all works stopped - all done.")
		os.Exit(0)
	}
}

//go run graceful_shutdown.go
//19:51:00.153522 graceful_shutdown.go:89: start works
//19:51:00.454886 graceful_shutdown.go:42: goody-2 is working...
//19:51:00.455073 graceful_shutdown.go:61: poor-1 is working...
//19:51:00.455149 graceful_shutdown.go:42: goody-1 is working...
//19:51:00.455168 graceful_shutdown.go:61: poor-2 is working...
//19:51:00.755556 graceful_shutdown.go:42: goody-2 is working...
//19:51:00.755612 graceful_shutdown.go:42: goody-1 is working...
//19:51:00.755612 graceful_shutdown.go:61: poor-2 is working...
//19:51:00.755621 graceful_shutdown.go:61: poor-1 is working...
//19:51:01.055625 graceful_shutdown.go:42: goody-2 is working...
//19:51:01.055660 graceful_shutdown.go:61: poor-1 is working...
//19:51:01.055678 graceful_shutdown.go:42: goody-1 is working...
//19:51:01.055679 graceful_shutdown.go:61: poor-2 is working...
//19:51:01.154456 graceful_shutdown.go:93: SIGINT ->
//19:51:01.154706 graceful_shutdown.go:102: main stopping: context canceled
//19:51:01.154713 graceful_shutdown.go:39: goody-1 stop working gracefully: context canceled
//19:51:01.154717 graceful_shutdown.go:39: goody-2 stop working gracefully: context canceled
//19:51:01.355016 graceful_shutdown.go:61: poor-2 is working...
//19:51:01.355076 graceful_shutdown.go:61: poor-1 is working...
//19:51:01.654597 graceful_shutdown.go:61: poor-1 is working...
//19:51:01.655149 graceful_shutdown.go:61: poor-2 is working...
//19:51:01.955146 graceful_shutdown.go:61: poor-1 is working...
//19:51:01.955186 graceful_shutdown.go:61: poor-2 is working...
//19:51:02.156759 graceful_shutdown.go:105: Timeout emergency exit fired.
//exit status 1

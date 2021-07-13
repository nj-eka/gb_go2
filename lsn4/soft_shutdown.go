package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"
)

const (
	workTimeout = 2 * time.Second
	waitExitTimeout = 1 * time.Second
)

var done = make(chan struct{})

func doWork(outerCtx context.Context, done chan <- struct{}) {
	localCtx, cancel := context.WithTimeout(context.Background(), workTimeout)
	defer cancel()
	defer close(done)
	select {
	case <-outerCtx.Done():
		fmt.Println("External stop goroutine.")
	case <-localCtx.Done():
		fmt.Println("Internal stop goroutine.")
	}
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	// sending SIGINT signal to itself simulation:
	//p, err := os.FindProcess(os.Getpid())
	//if err != nil {
	//	log.Fatal(err)
	//}
	//if err := p.Signal(os.Interrupt); err != nil {
	//	log.Fatal(err)
	//}

	go doWork(ctx, done)

	select {
	case <-ctx.Done():
		fmt.Println("Interrupt signal is fired. Stopping with ", ctx.Err()) // prints "context canceled"
		stop()                 // stop receiving signal notifications as soon as possible.
		select {
		case <-time.After(waitExitTimeout) :
			fmt.Println("Timeout emergency exit fired.")
			os.Exit(2)
		case <-done:
			fmt.Println("Soft exit - done.")
			os.Exit(1)
		}
	case <-done:
		fmt.Println("Work is done.")
		os.Exit(0)
	}
}
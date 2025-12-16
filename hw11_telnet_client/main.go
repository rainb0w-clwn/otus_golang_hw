package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var errLog = log.New(os.Stderr, "", 0)

func main() {
	timeout := flag.Duration("timeout", 10*time.Second, "timeout")
	flag.Parse()
	if len(flag.Args()) != 2 {
		errLog.Println("args validation error")
		return
	}
	host := flag.Arg(0)
	port := flag.Arg(1)
	rootCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT)
	defer stop()
	client := NewTelnetClient(net.JoinHostPort(host, port), *timeout, os.Stdin, os.Stdout)
	err := client.Connect()
	if err != nil {
		errLog.Println(fmt.Errorf("client connection error: %w", err))
		return
	}
	defer func() {
		_ = client.Close()
	}()
	go func() {
		_ = client.Send()
		stop()
	}()
	go func() {
		_ = client.Receive()
		stop()
	}()
	<-rootCtx.Done()
}

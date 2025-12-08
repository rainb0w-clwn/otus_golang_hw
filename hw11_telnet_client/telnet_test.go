package main

import (
	"bytes"
	"io"
	"net"
	"sync"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestTelnetClient(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		l, err := net.Listen("tcp", "127.0.0.1:")
		require.NoError(t, err)
		defer func() { require.NoError(t, l.Close()) }()

		var wg sync.WaitGroup
		wg.Add(2)

		go func() {
			defer wg.Done()

			in := &bytes.Buffer{}
			out := &bytes.Buffer{}

			timeout, err := time.ParseDuration("10s")
			require.NoError(t, err)

			client := NewTelnetClient(l.Addr().String(), timeout, io.NopCloser(in), out)
			require.NoError(t, client.Connect())
			defer func() { require.NoError(t, client.Close()) }()

			in.WriteString("hello\n")
			err = client.Send()
			require.NoError(t, err)

			err = client.Receive()
			require.NoError(t, err)
			require.Equal(t, "world\n", out.String())
		}()

		go func() {
			defer wg.Done()

			conn, err := l.Accept()
			require.NoError(t, err)
			require.NotNil(t, conn)
			defer func() { require.NoError(t, conn.Close()) }()

			request := make([]byte, 1024)
			n, err := conn.Read(request)
			require.NoError(t, err)
			require.Equal(t, "hello\n", string(request)[:n])

			n, err = conn.Write([]byte("world\n"))
			require.NoError(t, err)
			require.NotEqual(t, 0, n)
		}()

		wg.Wait()
	})
	t.Run("no server", func(t *testing.T) {
		client := NewTelnetClient(
			"0.0.0.0:8000",
			time.Second*10,
			io.NopCloser(&bytes.Buffer{}),
			&bytes.Buffer{},
		)
		err := client.Connect()
		require.Error(t, err)
	})
	t.Run("incorrect usage", func(t *testing.T) {
		client := NewTelnetClient("0.0.0.0:8000", time.Second*10, io.NopCloser(&bytes.Buffer{}), &bytes.Buffer{})
		sendErr := client.Send()
		require.ErrorIs(t, sendErr, ErrConnectionEmpty)
		receiveErr := client.Receive()
		require.ErrorIs(t, receiveErr, ErrConnectionEmpty)
		closeErr := client.Close()
		require.ErrorIs(t, closeErr, ErrConnectionEmpty)
	})
	t.Run("eof test", func(t *testing.T) {
		l, err := net.Listen("tcp", "0.0.0.0:")
		require.NoError(t, err)
		defer func() { require.NoError(t, l.Close()) }()
		wg := sync.WaitGroup{}
		wg.Add(2)
		go func() {
			defer wg.Done()
			in := &bytes.Buffer{}
			client := NewTelnetClient(l.Addr().String(), time.Second*10, io.NopCloser(in), &bytes.Buffer{})
			require.NoError(t, client.Connect())
			defer func() { require.NoError(t, client.Close()) }()
			in.WriteString("hello\u0000")
			err = client.Send()
			require.NoError(t, err)
		}()
		go func() {
			defer wg.Done()
			conn, err := l.Accept()
			require.NoError(t, err)
			require.NotNil(t, conn)
			defer func() { require.NoError(t, conn.Close()) }()
			request := make([]byte, 1024)
			n, err := conn.Read(request)
			require.NoError(t, err)
			require.Equal(t, "hello\u0000", string(request)[:n])
			_, err = conn.Read(request)
			require.Equal(t, err, io.EOF)
		}()
		wg.Wait()
	})
	t.Run("connection timeout", func(t *testing.T) {
		l, err := net.Listen("tcp", "0.0.0.0:")
		require.NoError(t, err)
		defer func() { require.NoError(t, l.Close()) }()
		client := NewTelnetClient(l.Addr().String(), time.Nanosecond*1, io.NopCloser(&bytes.Buffer{}), &bytes.Buffer{})
		require.Error(t, client.Connect())
	})
	t.Run("send/receive to closed connection", func(t *testing.T) {
		l, err := net.Listen("tcp", "0.0.0.0:")
		require.NoError(t, err)
		in := &bytes.Buffer{}
		client := NewTelnetClient(l.Addr().String(), time.Second*10, io.NopCloser(in), &bytes.Buffer{})
		require.NoError(t, client.Connect())
		require.NoError(t, l.Close())
		in.WriteString("hello\n")
		require.ErrorIs(t, client.Send(), syscall.ECONNRESET)
	})
	t.Run("send/receive to closed connection", func(t *testing.T) {
		l, err := net.Listen("tcp", "0.0.0.0:")
		require.NoError(t, err)
		client := NewTelnetClient(l.Addr().String(), time.Second*10, io.NopCloser(&bytes.Buffer{}), &bytes.Buffer{})
		require.NoError(t, client.Connect())
		require.NoError(t, l.Close())
		require.Error(t, client.Receive(), syscall.ECONNRESET)
	})
}

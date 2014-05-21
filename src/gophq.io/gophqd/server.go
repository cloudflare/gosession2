package gophqd

import (
	"gophq.io/proto"
	"gophq.io/tls"
	"log"
	"net"
	"sync"
)

type Server struct {
	sync.Mutex

	wg sync.WaitGroup

	TLSConfig *tls.TLSConfig
}

func (this *Server) Serve(l net.Listener) error {
	for {
		// TODO make this Accept loop more robust
		// see Server.Serve in $GOROOT/src/pkg/net/http/server.go

		// See also: https://code.google.com/p/go/issues/detail?id=6163

		c, err := l.Accept()
		if err != nil {
			return err
		}

		// read about Type switch at:
		// http://golang.org/ref/spec#Switch_statements
		switch conn := c.(type) {
		case *net.TCPConn:
			// TODO set socket buffer sizes
			// http://golang.org/pkg/net/#TCPConn.SetReadBuffer

			// switch socket to use TLS if configured
			c = this.TLSConfig.Server(c)

		case *net.UnixConn:
			log.Printf("using a unix socket [%T]", conn)

			// Maybe check credentials of connecting process:
			// http://golang.org/src/pkg/syscall/creds_test.go

			// Or pass some file descriptors as the first
			// step in the protocol to allow future passing
			// of file descriptors to shared memory
			// http://golang.org/src/pkg/syscall/passfd_test.go

			// The other piece of the puzzle is here:
			// https://github.com/jsgilmore/shm
		}

		// TODO track these goroutines with a WaitGroup so
		// that Server can be shut down in an orderly fashion.
		// Especially useful for testing.

		this.wg.Add(1)
		go func() {
			defer this.wg.Done()
			err := this.handleClient(c)
			if err != nil {
				log.Printf("handleClient: %v", err)
			}
		}()
	}
}

func (this *Server) Terminate() {
	// TODO Signal shutdown to all client goroutines

	// TODO Don't wait more than a configurable amount
	// of time for shutdown, since shutdown might be
	// due to SIGTERM, which will be followed by SIGKILL
	// if it takes too long

	// TODO Wait on the WaitGroup
	this.wg.Wait()
}

func (this *Server) handleClient(c net.Conn) error {
	// TODO could recover panics from the
	// client handler here, but it's debatable
	// whether doing so is a good idea

	// remember to close your Conn
	defer c.Close()

	// TODO Read/write deadlines could be a good idea
	// when dealing with potentially slow clients
	// http://golang.org/pkg/net/#Conn

	for {
		b, err := proto.ReadRequestOrResponse(c)
		if err != nil {
			return err
		}

		var anyReq proto.Request
		err = proto.Decode(b, &anyReq)
		if err != nil {
			return err
		}

		switch req := anyReq.Body.(type) {
		case *proto.ProduceRequest:
			response, err := this.handleProduceRequest(req)
			if err != nil {
				return err
			}
			b, err = proto.Encode(&proto.Response{response})
		case *proto.FetchRequest:
			response, err := this.handleFetchRequest(req)
			if err != nil {
				return err
			}
			b, err = proto.Encode(&proto.Response{response})
		}
		if err != nil {
			return err
		}

		_, err = c.Write(b)
		if err != nil {
			return err
		}
	}
}

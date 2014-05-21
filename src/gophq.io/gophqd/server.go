package gophqd

import (
	"gophq.io/proto"
	"log"
	"net"
	"sync"
)

type Server struct {
	wg sync.WaitGroup

	DataDir string
}

func (this *Server) DoMaintenance() error {
	for {

	}
	return nil
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
			err = this.handleProduceRequest(req)
		case *proto.FetchRequest:
			err = this.handleFetchRequest(req)
		}
		if err != nil {
			return err
		}
	}
}

func (this *Server) handleProduceRequest(req *proto.ProduceRequest) error {
	log.Printf("%+v", req)
	for _, msgBlock := range req.MsgSet.Messages {
		msg := msgBlock.Msg
		log.Printf("%+v: %x -> %x", msg, msg.Key, msg.Value)
	}
	return nil
}

func (this *Server) handleFetchRequest(req *proto.FetchRequest) error {
	return nil
}

// TODO write message to disk... with crc32 hash

// demo crc32 vs crc32c

// A log for a topic named "my_topic" with two partitions consists
// of two directories (namely my_topic_0 and my_topic_1) populated
// with data files containing the messages for that topic. The format
// of the log files is a sequence of "log entries""; each log entry
// is a 4 byte integer N storing the message length which is followed
// by the N message bytes. Each message is uniquely identified by
// a 64-bit integer offset giving the byte position of the start of
// this message in the stream of all messages ever sent to that topic
// on that partition. The on-disk format of each message is given
// below. Each log file is named with the offset of the first message
// it contains. So the first file created will be 00000000000.kafka,
// and each additional file will have an integer name roughly S bytes
// from the previous file where S is the max log file size given
// in the configuration.

// The exact binary format for messages is versioned and maintained
// as a standard interface so message sets can be transfered between
// producer, broker, and client without recopying or conversion when
// desirable. This format is as follows:

// On-disk format of a message
// message length : 4 bytes (value: 1+4+n)
// "magic" value  : 1 byte
// crc            : 4 bytes
// payload        : n bytes

// Writes

// The log allows serial appends which always go to the last file.
// This file is rolled over to a fresh file when it reaches a configurable size (say 1GB).
// The log takes two configuration parameter M which gives the number of messages to write
// before forcing the OS to flush the file to disk, and S which gives a number of seconds
// after which a flush is forced. This gives a durability guarantee of losing at most M
// messages or S seconds of data in the event of a system crash.

// Reads

// Reads are done by giving the 64-bit logical offset of a message and an S-byte max
// chunk size. This will return an iterator over the messages contained in the S-byte
// buffer. S is intended to be larger than any single message, but in the event of an
// abnormally large message, the read can be retried multiple times, each time doubling
// the buffer size, until the message is read successfully. A maximum message and
// buffer size can be specified to make the server reject messages larger than some
// size, and to give a bound to the client on the maximum it need ever read to get a
// complete message. It is likely that the read buffer ends with a partial message,
// this is easily detected by the size delimiting.

// The actual process of reading from an offset requires first locating the log
// segment file in which the data is stored, calculating the file-specific offset
// from the global offset value, and then reading from that file offset. The search
// is done as a simple binary search variation against an in-memory range
// maintained for each file.

// The log provides the capability of getting the most recently written message to
// allow clients to start subscribing as of "right now". This is also useful in the
// case the consumer fails to consume its data within its SLA-specified number of
// days. In this case when the client attempts to consume a non-existant offset it is
// given an OutOfRangeException and can either reset itself or fail as appropriate
// to the use case.

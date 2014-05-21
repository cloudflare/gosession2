package gophqd

import (
	"bufio"
	"gophq.io/proto"
	"log"
	"os"
)

func (this *Server) handleProduceRequest(req *proto.ProduceRequest) (*proto.ProduceResponse, error) {
	response := &proto.ProduceResponse{
		Topic: req.Topic,
		Err:   proto.NoError,
	}

	const flags = os.O_CREATE | os.O_APPEND | os.O_WRONLY
	f, err := os.OpenFile(req.Topic+".dat", flags, 0666)
	if err != nil {
		log.Printf("OpenFile: %v", err)
		return nil, err
	}
	defer f.Close()

	offset, err := f.Seek(0, os.SEEK_CUR)
	if err != nil {
		log.Printf("Seek: %v", err)
		return nil, err
	}
	response.Offset = offset

	// TODO a pool of these that get reused using
	// bufio.Writer#Reset could help a lot
	bufw := bufio.NewWriterSize(f, 128<<10)
	defer bufw.Flush()

	for _, msgBlock := range req.MsgSet.Messages {
		msg := msgBlock.Msg
		log.Printf("%+v: %x -> %x", msg, msg.Key, msg.Value)

		// TODO calculate and pass the real message offset
		err = msg.Write(0, bufw)
		if err != nil {
			// an error from Write probably indicates an I/O error
			// that should be escalated
			log.Printf("Write: %v", err)
			return nil, err
		}
	}

	return response, nil
}

package gophqd

import (
	"gophq.io/proto"
	"log"
	"os"
)

func (this *Server) handleFetchRequest(req *proto.FetchRequest) (*proto.FetchResponse, error) {
	log.Printf("%+v", req)

	response := &proto.FetchResponse{Topic: req.Topic}

	// TODO topic name validation
	f, err := os.Open(req.Topic + ".dat")
	if err != nil {
		log.Printf("unknown topic or partition: %v", req.Topic)
		response.Err = proto.UnknownTopicOrPartition
		return response, nil
	}
	defer f.Close()

	if req.FetchOffset == proto.LatestOffset {
		// TODO not implemented yet
		response.Err = proto.OffsetOutOfRange
		return response, nil
	}

	if req.FetchOffset == proto.EarliestOffset {
		// TODO topic only has one file, so we can cheat
		// a real server with rolling files needs to do more here
		req.FetchOffset = 0
	}

	offset, err := f.Seek(req.FetchOffset, os.SEEK_SET)
	if err != nil {
		log.Printf("Seek: %v: %v", req.Topic, err)
		return nil, err
	}
	if offset != req.FetchOffset {
		// Only exact offsets allowed.
		response.Err = proto.OffsetOutOfRange
		return response, nil
	}

	return &proto.FetchResponse{}, nil
}

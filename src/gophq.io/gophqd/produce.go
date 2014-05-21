package gophqd

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

	for _, msgBlock := range req.MsgSet.Messages {
		msg := msgBlock.Msg
		log.Printf("%+v: %x -> %x", msg, msg.Key, msg.Value)

		// TODO write messages to file and index
	}

	return response, nil
}

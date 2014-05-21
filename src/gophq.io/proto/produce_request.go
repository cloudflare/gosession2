package proto

type ProduceRequest struct {
	Topic  string
	MsgSet MessageSet
}

func (p *ProduceRequest) encode(pe packetEncoder) error {
	err := pe.putString(p.Topic)
	if err != nil {
		return err
	}
	pe.push(&lengthField{})
	err = p.MsgSet.encode(pe)
	if err != nil {
		return err
	}
	return pe.pop()
}

func (p *ProduceRequest) decode(pd packetDecoder) error {
	topic, err := pd.getString()
	if err != nil {
		return err
	}
	p.Topic = topic

	msgSetSize, err := pd.getInt32()
	if err != nil {
		return err
	}

	msgSetDecoder, err := pd.getSubset(int(msgSetSize))
	if err != nil {
		return err
	}
	err = (&p.MsgSet).decode(msgSetDecoder)

	return err
}

func (p *ProduceRequest) AddMessage(key, value Encoder) {
	var kb []byte
	var vb []byte
	if key != nil {
		kb, _ = key.Encode()
	}
	if value != nil {
		vb, _ = value.Encode()
	}
	p.MsgSet.addMessage(&Message{Key: kb, Value: vb})
}

func (p *ProduceRequest) key() int16 {
	return 0
}

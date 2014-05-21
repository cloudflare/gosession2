package proto

type FetchResponse struct {
	Topic  string
	Err    KError
	MsgSet MessageSet
}

func (fr *FetchResponse) encode(pe packetEncoder) error {
	err := pe.putString(fr.Topic)
	if err != nil {
		return err
	}

	pe.putInt16(int16(fr.Err))

	pe.push(&lengthField{})

	err = fr.MsgSet.encode(pe)
	if err != nil {
		return err
	}
	return pe.pop()
}

func (fr *FetchResponse) decode(pd packetDecoder) error {
	var err error

	fr.Topic, err = pd.getString()
	if err != nil {
		return err
	}

	tmp, err := pd.getInt16()
	if err != nil {
		return err
	}
	fr.Err = KError(tmp)

	msgSetSize, err := pd.getInt32()
	if err != nil {
		return err
	}

	msgSetDecoder, err := pd.getSubset(int(msgSetSize))
	if err != nil {
		return err
	}
	err = (&fr.MsgSet).decode(msgSetDecoder)

	return err
}

func (fr *FetchResponse) AddMessage(key, value Encoder, offset int64) {
	var kb []byte
	var vb []byte
	if key != nil {
		kb, _ = key.Encode()
	}
	if value != nil {
		vb, _ = value.Encode()
	}
	fr.MsgSet.addMessageOffset(&Message{Key: kb, Value: vb}, offset)
}

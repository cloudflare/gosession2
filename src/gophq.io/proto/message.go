package proto

type Message struct {
	Key   []byte      // the message key, may be nil
	Value []byte      // the message contents
	Set   *MessageSet // the message set a message might wrap
}

func (m *Message) encode(pe packetEncoder) error {
	pe.push(&crc32Field{})

	err := pe.putBytes(m.Key)
	if err != nil {
		return err
	}

	err = pe.putBytes(m.Value)
	if err != nil {
		return err
	}

	return pe.pop()
}

func (m *Message) decode(pd packetDecoder) (err error) {
	err = pd.push(&crc32Field{})
	if err != nil {
		return err
	}

	m.Key, err = pd.getBytes()
	if err != nil {
		return err
	}

	m.Value, err = pd.getBytes()
	if err != nil {
		return err
	}

	return pd.pop()
}

func (m *Message) decodeSet() (err error) {
	pd := realDecoder{raw: m.Value}
	m.Set = &MessageSet{}
	return m.Set.decode(&pd)
}

func NewMessage(key, value []byte) *Message {
	var kb []byte
	var vb []byte
	if key != nil {
		kb, _ = ByteEncoder(key).Encode()
	}
	if value != nil {
		vb, _ = ByteEncoder(value).Encode()
	}
	return &Message{Key: kb, Value: vb}
}

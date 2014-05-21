package proto

import (
	"encoding/binary"
	"io"
)

type Message struct {
	Crc   uint32      // message CRC
	Key   []byte      // the message key, may be nil
	Value []byte      // the message contents
	Set   *MessageSet // the message set a message might wrap
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

// Kafka's wire protocol uses big endian values, which
// makes sense for Java, but this means that many values
// need to be converted when writing so that they can
// be sent as-is using sendfile(2) when reading.

func (m *Message) Write(offset int64, w io.Writer) error {
	msgLen := 8 + 4 + 4 + 4 + len(m.Key) + 4 + len(m.Value)

	// offset
	err := binary.Write(w, binary.BigEndian, offset)
	if err != nil {
		return err
	}

	// message length
	err = binary.Write(w, binary.BigEndian, int32(msgLen))
	if err != nil {
		return err
	}

	// Crc
	err = binary.Write(w, binary.BigEndian, m.Crc)
	if err != nil {
		return err
	}

	// length of key
	err = binary.Write(w, binary.BigEndian, int32(len(m.Key)))
	if err != nil {
		return err
	}

	// key
	_, err = w.Write(m.Key)
	if err != nil {
		return err
	}

	// length of value
	err = binary.Write(w, binary.BigEndian, int32(len(m.Value)))
	if err != nil {
		return err
	}

	// value
	_, err = w.Write(m.Value)
	return err
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
	crc := &crc32Field{}
	err = pd.push(crc)
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

	err = pd.pop()
	if err != nil {
		return err
	}
	m.Crc = crc.crc

	return nil
}

func (m *Message) decodeSet() (err error) {
	pd := realDecoder{raw: m.Value}
	m.Set = &MessageSet{}
	return m.Set.decode(&pd)
}

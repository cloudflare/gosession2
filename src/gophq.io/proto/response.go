package proto

import (
	"fmt"
	"io"
	"log"
)

type Response struct {
	Body encoder
}

func (r *Response) encode(pe packetEncoder) (err error) {
	pe.push(&lengthField{})
	err = r.Body.encode(pe)
	if err != nil {
		return err
	}
	return pe.pop()
}

type lengthHeader struct {
	length int32
}

func (r *lengthHeader) decode(pd packetDecoder) (err error) {
	r.length, err = pd.getInt32()
	if err != nil {
		return err
	}
	const maxMessageSize = 512 << 20
	if r.length <= 0 || r.length > maxMessageSize {
		return DecodingError{Info: fmt.Sprintf("Message too large or too small. Got %d", r.length)}
	}
	return nil
}

func ReadRequestOrResponse(r io.Reader) ([]byte, error) {
	var header [4]byte
	_, err := io.ReadFull(r, header[:])
	if err != nil {
		return nil, err
	}

	decodedHeader := lengthHeader{}
	err = Decode(header[:], &decodedHeader)
	if err != nil {
		return nil, err
	}

	log.Printf("header done, need to read %v", decodedHeader.length)

	// TODO it would be good to recycle this allocation
	buf := make([]byte, decodedHeader.length)
	_, err = io.ReadFull(r, buf)
	if err != nil {
		return nil, err
	}

	return buf, nil
}

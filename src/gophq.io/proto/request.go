package proto

import (
	"fmt"
)

type requestBody interface {
	encoder
	decoder
	key() int16
}

type Request struct {
	Body requestBody
}

func (r *Request) encode(pe packetEncoder) (err error) {
	pe.push(&lengthField{})
	pe.putInt16(r.Body.key())
	err = r.Body.encode(pe)
	if err != nil {
		return err
	}
	return pe.pop()
}

func (r *Request) decode(pd packetDecoder) error {
	key, err := pd.getInt16()
	if err != nil {
		return err
	}

	switch key {
	case 0:
		r.Body = &ProduceRequest{}
	case 1:
		r.Body = &FetchRequest{}
	default:
		return DecodingError{Info: fmt.Sprintf("invalid key: %d", key)}
	}

	return r.Body.decode(pd)
}

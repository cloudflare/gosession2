package proto

import (
	"time"
)

type FetchRequest struct {
	Topic       string
	MinBytes    int32
	MaxBytes    int32
	MaxWaitTime time.Duration
	FetchOffset int64
}

func (f *FetchRequest) encode(pe packetEncoder) error {
	err := pe.putString(f.Topic)
	if err != nil {
		return err
	}
	pe.putInt32(f.MinBytes)
	pe.putInt32(f.MaxBytes)
	pe.putInt64(int64(f.MaxWaitTime))
	pe.putInt64(f.FetchOffset)
	return nil
}

func (f *FetchRequest) decode(pd packetDecoder) error {
	var err error
	f.Topic, err = pd.getString()
	if err != nil {
		return err
	}

	f.MinBytes, err = pd.getInt32()
	if err != nil {
		return err
	}

	f.MaxBytes, err = pd.getInt32()
	if err != nil {
		return err
	}

	tmp, err := pd.getInt64()
	if err != nil {
		return err
	}
	f.MaxWaitTime = time.Duration(tmp)

	f.FetchOffset, err = pd.getInt64()
	if err != nil {
		return err
	}

	return err
}

func (f *FetchRequest) key() int16 {
	return 1
}

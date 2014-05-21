package proto

type ProduceResponse struct {
	Topic  string
	Err    KError
	Offset int64
}

func (pr *ProduceResponse) encode(pe packetEncoder) error {
	err := pe.putString(pr.Topic)
	if err != nil {
		return err
	}
	pe.putInt16(int16(pr.Err))
	pe.putInt64(pr.Offset)
	return nil
}

func (pr *ProduceResponse) decode(pd packetDecoder) error {
	topic, err := pd.getString()
	if err != nil {
		return err
	}
	pr.Topic = topic

	tmp, err := pd.getInt16()
	if err != nil {
		return err
	}
	pr.Err = KError(tmp)

	pr.Offset, err = pd.getInt64()
	if err != nil {
		return err
	}

	return nil
}

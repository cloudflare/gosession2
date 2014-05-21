package proto

type MessageBlock struct {
	Offset int64
	Msg    *Message
}

// Messages convenience helper which returns either all the
// messages that are wrapped in this block
func (msb *MessageBlock) Messages() []*MessageBlock {
	if msb.Msg.Set != nil {
		return msb.Msg.Set.Messages
	}
	return []*MessageBlock{msb}
}

func (msb *MessageBlock) encode(pe packetEncoder) error {
	pe.putInt64(msb.Offset)
	pe.push(&lengthField{})
	err := msb.Msg.encode(pe)
	if err != nil {
		return err
	}
	return pe.pop()
}

func (msb *MessageBlock) decode(pd packetDecoder) (err error) {
	msb.Offset, err = pd.getInt64()
	if err != nil {
		return err
	}

	pd.push(&lengthField{})
	if err != nil {
		return err
	}

	msb.Msg = new(Message)
	err = msb.Msg.decode(pd)
	if err != nil {
		return err
	}

	return pd.pop()
}

type MessageSet struct {
	Messages   []*MessageBlock
	Incomplete bool
}

func (ms *MessageSet) encode(pe packetEncoder) error {
	for i := range ms.Messages {
		err := ms.Messages[i].encode(pe)
		if err != nil {
			return err
		}
	}
	return nil
}

// TODO a sync.Cache containing []*MessageBlock could help here

func (ms *MessageSet) decode(pd packetDecoder) (err error) {
	ms.Messages = nil

	for pd.remaining() > 0 {
		msb := new(MessageBlock)
		err = msb.decode(pd)
		switch err {
		case nil:
			ms.Messages = append(ms.Messages, msb)
		case InsufficientData:
			ms.Incomplete = true
			return nil
		default:
			return err
		}
	}

	return nil
}

func (ms *MessageSet) addMessage(msg *Message) {
	block := new(MessageBlock)
	block.Msg = msg
	ms.Messages = append(ms.Messages, block)
}

func (ms *MessageSet) addMessageOffset(msg *Message, offset int64) {
	block := new(MessageBlock)
	block.Msg = msg
	block.Offset = offset
	ms.Messages = append(ms.Messages, block)
}

func (ms *MessageSet) Clear() {
	ms.Messages = ms.Messages[:0]
}

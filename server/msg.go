package server

type Msg struct {
	From    uint64
	To      uint64
	Content string
}

type MsgQueue struct {
	Queues []Msg
	Lock   bool
	IsNil  bool
}

func (this *MsgQueue) AddMsg(msg Msg) (err error) {
	err = checkMsg(msg)
	if err != nil {
		return
	}
	for this.Lock {
	}
	this.Lock = true
	this.Queues = append(this.Queues, msg)
	this.IsNil = false
	this.Lock = false
	return
}

func (this *MsgQueue) PopMsg() (msg Msg, err error) {
	for this.Lock {
	}
	this.Lock = true
	if this.IsNil {
		err = errors.New("msgQueue is nil")
		this.Lock = false
		return
	}
	msg = this.Queues[0]
	if len(this.Queues) == 1 {
		this.IsNil = true
		this.Queues = nil
	} else {
		this.Queues = this.Queues[1:]
	}
	this.Lock = false
	return
}

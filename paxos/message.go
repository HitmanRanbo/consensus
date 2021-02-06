package paxos

type MessageType int

const (
	Prepare MessageType = iota + 1
	Promise
	Propose
	Accept
)

type Proposal struct {
	id    int
	value string
}

func (p Proposal) GetId() int {
	return p.id
}

func (p Proposal) GetValue() string {
	return p.value
}

func NewMessage(from, to, n int, typ MessageType, proposal Proposal) Message {
	return Message{
		from:     from,
		to:       to,
		n:        n,
		typ:      typ,
		proposal: proposal,
	}
}

type Message struct {
	from     int
	to       int
	n        int
	typ      MessageType
	proposal Proposal
}

func (msg Message) GetFrom() int {
	return msg.from
}

func (msg Message) GetTo() int {
	return msg.to
}

func (msg Message) GetN() int {
	return msg.n
}

func (msg Message) GetType() MessageType {
	return msg.typ
}

func (msg Message) GetProposal() Proposal {
	return msg.proposal
}

package paxos

import (
	"log"
	"time"
)

func NewAcceptor(client Client) *Acceptor {
	return &Acceptor{client: client}
}

type Acceptor struct {
	client           Client
	acceptedProposal Proposal
	n                int
	learners         []Learner
}

func (a *Acceptor) Run() {
	for {
		msg, err := a.client.Recv(time.Second)
		if err != nil {
			log.Printf("failed to recive. err is %s", err.Error())
			continue
		}

		switch msg.GetType() {
		case Prepare:
			if a.acceptPrepare(msg) {
				a.n = msg.GetN()
				a.send(msg.GetFrom(), Promise)
			}
		case Propose:
			if a.acceptPropose(msg) {
				a.acceptedProposal = msg.GetProposal()
				for _, learner := range a.learners {
					a.send(learner.client.GetId(), Accept)
				}
			}
		default:
			log.Panicf("acceptor: %d unexpected message type: %v", a.GetClientId(), msg.typ)
		}
	}
}

// Process a PREPARE(n) message
// If the n is not greater than the greatest n the acceptor has seen, reject it.
// Or if the acceptor has accepted any proposal, notice it
func (a *Acceptor) acceptPrepare(msg Message) bool {
	return msg.GetN() > a.n
}

// If the ID is the highest number it has processed
// then accept the proposal and propagate the value to the proposer and to all the learners.
func (a *Acceptor) acceptPropose(msg Message) bool {
	return msg.GetN() == a.n
}

func (a *Acceptor) GetClientId() int {
	return a.client.GetId()
}

func (a *Acceptor) send(to int, typ MessageType) {
	log.Printf("accepter: %d send a %v request with propsal: %v. current n is: %d", a.GetClientId(), typ, a.acceptedProposal, a.n)
	err := a.client.Send(NewMessage(a.GetClientId(), to, a.n, typ, a.acceptedProposal))
	if err != nil {
		log.Fatal(err)
	}
}

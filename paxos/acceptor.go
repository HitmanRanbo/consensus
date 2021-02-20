package paxos

import (
	"log"
	"time"
)

func NewAcceptor(client Client, learners ...Learner) *Acceptor {
	return &Acceptor{client: client, learners: learners}
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
			log.Printf("acceptor: %d failed to recive. err is %s", a.GetClientId(), err.Error())
			continue
		}

		switch msg.GetType() {
		case Prepare:
			// Process a PREPARE(n) message
			// If the n is not greater than the greatest n the acceptor has seen, reject it.
			// Or if the acceptor has accepted any proposal, notice it
			if msg.GetN() > a.n {
				a.n = msg.GetN()
				a.send(msg.GetFrom(), Promise)
			}
		case Propose:
			// If the ID is the highest number it has processed
			// then accept the proposal and propagate the value to the proposer and to all the learners.
			if msg.GetN() == a.n {
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

func (a *Acceptor) GetClientId() int {
	return a.client.GetId()
}

func (a *Acceptor) send(to int, typ MessageType) {
	log.Printf("accepter: %d send request to %d with propsal: %v. request type is: %d current n is: %d", a.GetClientId(), to, a.acceptedProposal, typ, a.n)
	err := a.client.Send(NewMessage(a.GetClientId(), to, a.n, typ, a.acceptedProposal))
	if err != nil {
		log.Fatal(err)
	}
}

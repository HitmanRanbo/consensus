package paxos

import (
	"log"
	"time"
)

type Learner struct {
	client           Client
	acceptedProposal Proposal
}

func NewLearner(client Client) *Learner {
	return &Learner{client: client}
}

func (l *Learner) Learn() Proposal {
	for {
		msg, err := l.client.Recv(time.Hour)
		if err != nil {
			log.Printf("failed to recive. err is %s", err.Error())
			continue
		}

		if msg.GetType() != Accept {
			log.Panicf("learner: %d unexpected message type: %v", l.client.GetId(), msg.typ)
			continue
		}

		l.acceptedProposal = msg.GetProposal()
		return  l.acceptedProposal
	}
}
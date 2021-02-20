package paxos

import (
	"log"
	"time"
)

type Learner struct {
	client      Client
	learnFrom   int
	chosenValue string
	valueCount  map[string]int
}

func NewLearner(client Client) *Learner {
	return &Learner{client: client, valueCount: make(map[string]int)}
}

func (l *Learner) AddLearnFrom() {
	l.learnFrom++
}

func (l *Learner) Learn() string {
	for {
		msg, err := l.client.Recv(time.Hour)
		if err != nil {
			log.Printf("failed to recive. err is %s", err.Error())
			continue
		}

		if msg.GetType() != Accept {
			log.Panicf("learner: %d unexpected message type: %v", l.client.GetId(), msg.typ)
		}

		if l.choose(msg.GetProposal().GetValue()) {
			return l.chosenValue
		}
	}
}

func (l *Learner) choose(value string) bool {
	l.valueCount[value]++
	if l.valueCount[value] >= l.learnFrom/2+1 {
		l.chosenValue = value
		return true
	}
	return false
}

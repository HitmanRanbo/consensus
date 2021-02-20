package paxos

import (
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
	"time"
)

func TestProposer(t *testing.T) {
	testWith2ProposerAnd3Acceptor(t, randomInterval())

}

func randomInterval() time.Duration {
	rand.Seed(time.Now().UnixNano())
	interval := rand.Int() % 9
	if interval == 0 || interval == 3 {
		interval++
	}
	return time.Duration(interval)*time.Millisecond
}


// Test with 2 proposer and 3 acceptor.
// It takes 1 millisecond for proposer to send a request.
// So every proposer need 3 millisecond second to finish prepare process, and another 2 millisecond to proposal a majority
// If the interval is greater than 3 millisecond, the first proposer will proposal the majority before the later one finish its prepare
// If the interval is between 1 millisecond and 3 millisecond, the later proposer will finish its prepare before the first one proposal a majority
func testWith2ProposerAnd3Acceptor(t *testing.T, interval time.Duration) {
	t.Logf("run proposer test with interval: %v", interval)
	network := NewEagerNetwork()

	l1 := NewLearner(network.NewClient(3001))
	l2 := NewLearner(network.NewClient(3002))
	l3 := NewLearner(network.NewClient(3003))

	a1 := NewAcceptor(network.NewClient(2001), l1, l2, l3)
	a2 := NewAcceptor(network.NewClient(2002), l2, l1, l3)
	a3 := NewAcceptor(network.NewClient(2003), l3, l2, l1)

	p1 := NewProposer(network.NewClient(1001), a1, a2, a3)
	p2 := NewProposer(network.NewClient(1002), a3, a1, a2)

	go a1.Run()
	go a2.Run()
	go a3.Run()

	go p1.Run("Yes")
	time.Sleep(interval)
	go p2.Run("No")

	if interval < 4 * time.Millisecond {
		assert.Equal(t, "No", l1.Learn())
		assert.Equal(t, "No", l2.Learn())
		assert.Equal(t, "No", l3.Learn())
	} else {
		assert.Equal(t, "Yes", l1.Learn())
		assert.Equal(t, "Yes", l2.Learn())
		assert.Equal(t, "Yes", l3.Learn())
	}
}


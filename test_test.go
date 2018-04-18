package raft

//
// Raft tests.
//
// we will use the original test_test.go to test your code for grading.
// so, while you can modify this code to help you debug, pleasp
// test with the original before submitting.
//

import (
	// "encoding/binary"
	"fmt"
	"os"
	"testing"
	"time"
)

// import "math/rand"

// import "sync/atomic"
// import "sync"

// The tester generously allows solutions to complete elections in one second
// (much more than the paper's range of timeouts).
const RaftElectionTimeout = 1000 * time.Millisecond
const CommandBreak = 50 * time.Millisecond

var responseCount int = 0

func (rf *Raft) ReceiveResponse(response interface{}) {
	responseCount += 1
}

// PASSED
func TestStartCommand(t *testing.T) {
	servers := 30
	cfg := make_config(t, servers, false)
	defer cfg.cleanup()

	var filename = "./data"

	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	fmt.Printf("Initial election ...\n")

	leader := cfg.checkOneLeader()
	fmt.Printf("Leader: %d\n", leader)

	// time.Sleep(2000 * time.Millisecond)

	start := time.Now()

	var index int

	var count int = 1
	for i := 0; i < count; i++ {
		// fmt.Printf("i: %d\n", i)
		var ok bool
		command := i

		// 生成数字签名
		cmdBytes := GetBytes(command)
		sig := signature(cmdBytes)

		time.Sleep(CommandBreak)

		index, _, ok = cfg.rafts[leader].Start(command, sig)

		if ok {
			fmt.Println("Index:", index)
		} else {
			fmt.Println("Failed.")
		}
	}

	time.Sleep(2 * RaftElectionTimeout)

	n, _ := cfg.nCommitted(index)

	for server := range cfg.rafts {
		rf := *cfg.rafts[server]
		fmt.Printf("server: %d, commitIndex: %d\n", rf.me, rf.commitIndex)
		// PrintSortedMap(rf.m)
	}

	fmt.Println("Committed number:", n)

	end := time.Now()
	elapsed := end.Sub(start)
	timeused := fmt.Sprintf("%v\n", elapsed-2*RaftElectionTimeout-time.Duration(count)*CommandBreak)
	if _, err := f.WriteString(timeused); err != nil {
		panic(err)
	}
	//fmt.Printf("command count: %d, time elapsed: %v\n", count, elapsed-2*RaftElectionTimeout-100*CommandBreak)

	// end := time.Now()
	// elapsed := end.Sub(start)
	// fmt.Printf("command count: %d, time elapsed: %v\n", count, elapsed)

}

// PASSED
// func TestInitialElection2A(t *testing.T) {
// 	servers := 6
// 	cfg := make_config(t, servers, false)
// 	defer cfg.cleanup()

// 	fmt.Printf("Test (2A): initial election ...\n")

// 	// is a leader elected?
// 	leader := cfg.checkOneLeader()

// 	command := 20
// 	// 生成数字签名
// 	cmdBytes, _ := GetBytes(command)
// 	sig := signature(cmdBytes)

// 	index, _, ok := cfg.rafts[leader].Start(command, sig)
// 	fmt.Printf("index: %d, ok: %v\n", index, ok)

// 	// does the leader+term stay the same if there is no network failure?
// 	term1 := cfg.checkTerms()
// 	time.Sleep(2 * RaftElectionTimeout)
// 	term2 := cfg.checkTerms()
// 	if term1 != term2 {
// 		fmt.Printf("warning: term changed even though there were no failures")
// 	}

// 	fmt.Printf("  ... Passed\n")
// }

// // PASSED
// func TestReElection2A(t *testing.T) {
// 	servers := 6
// 	cfg := make_config(t, servers, false)
// 	defer cfg.cleanup()

// 	fmt.Printf("Test: election after network failure ...\n")

// 	leader1 := cfg.checkOneLeader()
// 	fmt.Printf("leader1: %d\n", leader1)

// 	// if the leader disconnects, a new one should be elected.
// 	cfg.disconnect(leader1)
// 	tmpLeader := cfg.checkOneLeader()
// 	fmt.Printf("tmpLeader: %d\n", tmpLeader)

// 	// if the old leader rejoins, that shouldn't
// 	// disturb the old leader.
// 	cfg.connect(leader1)
// 	leader2 := cfg.checkOneLeader()
// 	fmt.Printf("leader2: %d\n", leader2)

// 	fmt.Printf("  ... Passed\n")
// }

// // PASSED
// func TestBasicAgree2B(t *testing.T) {
// 	servers := 6
// 	cfg := make_config(t, servers, false)
// 	defer cfg.cleanup()

// 	fmt.Printf("Test (2B): basic agreement ...\n")

// 	iters := 1
// 	for index := 1; index < iters+1; index++ {
// 		// fmt.Println("....")
// 		nd, _ := cfg.nCommitted(index)
// 		if nd > 0 {
// 			t.Fatalf("some have committed before Start()")
// 		}

// 		xindex := cfg.one(index*100, servers)
// 		if xindex != index {
// 			t.Fatalf("got index %v but expected %v", xindex, index)
// 		}
// 	}

// 	fmt.Printf("  ... Passed\n")
// }

// PASSED
// func TestFailAgree2B(t *testing.T) {
// 	servers := 6
// 	cfg := make_config(t, servers, false)
// 	defer cfg.cleanup()

// 	fmt.Printf("Test (2B): agreement despite follower disconnection ...\n")

// 	ret := cfg.one(101, servers)
// 	fmt.Printf("ret: %d\n", ret)

// 	// follower network disconnection
// 	leader := cfg.checkOneLeader()
// 	fmt.Printf("leader: %d\n", leader)
// 	cfg.disconnect((leader + 1) % servers)

// 	// agree despite one disconnected server?
// 	cfg.one(102, servers-1)
// 	cfg.one(103, servers-1)
// 	time.Sleep(RaftElectionTimeout)
// 	cfg.one(104, servers-1)
// 	cfg.one(105, servers-1)

// 	// re-connect
// 	cfg.connect((leader + 1) % servers)

// 	// agree with full set of servers?
// 	cfg.one(106, servers)
// 	time.Sleep(RaftElectionTimeout)
// 	cfg.one(107, servers)

// 	fmt.Printf("  ... Passed\n")
// }

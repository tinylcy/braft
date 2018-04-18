package raft

//
// this is an outline of the API that raft must expose to
// the service (or tester). see comments below for
// each of these functions for more details.
//
// rf = Make(...)
//   create a new Raft server.
// rf.Start(command interface{}) (index, term, isleader)
//   start agreement on a new log entry
// rf.GetState() (term, isLeader)
//   ask a Raft for its current term, and whether it thinks it is leader
// ApplyMsg
//   each time a new entry is committed to the log, each Raft peer
//   should send an ApplyMsg to the service (or tester)
//   in the same server.
//

import (
	"braft/labrpc"
	//"net"
	"bytes"
	"encoding/gob"
	"fmt"
	"math/rand"
	"os"
	// "strconv"
	"sync"
	"time"
)

// import "bytes"

var AddedTime int64
var CommittedTime int64

//
// as each Raft peer becomes aware that successive log entries are
// committed, the peer should send an ApplyMsg to the service (or
// tester) on the same server, via the applyCh passed to Make().
//
type ApplyMsg struct {
	Index       int
	Command     interface{}
	UseSnapshot bool   // ignore for lab2; only used in lab3
	Snapshot    []byte // ignore for lab2; only used in lab3
}

type Status int

const (
	FOLLOWER  Status = iota // value --> 0
	CANDIDATE               // value --> 1
	LEADER                  // value --> 2
)

const quorum = 19

const HEARTBEAT_TIME int = 50

//
// A Go object implementing a single Raft peer.
//
type Raft struct {
	mu            sync.Mutex          // Lock to protect shared access to this peer's state
	peers         []*labrpc.ClientEnd // RPC end points of all peers
	persister     *Persister          // Object to hold this peer's persisted state
	me            int                 // this peer's index into peers[]
	currentLeader int

	// Your data here (2A, 2B, 2C).
	// Look at the paper's Figure 2 for a description of what
	// state a Raft server must maintain.
	logs []Entry

	commitIndex int
	lastApplied int

	nextIndex []int

	currentTerm int
	votedFor    int

	status                 Status
	beLeader               chan bool
	timeToCommit           chan bool
	grantVote              chan bool
	getHeartBeat           chan bool
	getCommit              chan bool
	getAppendEntriesCommit chan bool

	applyCh   chan ApplyMsg
	voteCount int

	m map[AppendEntriesCommitKey]int

	lastHeartbeat time.Time
}

type Entry struct {
	Index     int
	Term      int
	Command   interface{}
	Signature []byte
}

func (entry Entry) String() string {
	return fmt.Sprintf("Index: %d, Term: %d, Command: %v",
		entry.Index, entry.Term, entry.Command)
}

// return currentTerm and whether this server
// believes it is the leader.
func (rf *Raft) GetState() (int, bool) {

	var term int
	var isLeader bool
	// Your code here (2A).
	rf.lock()
	term = rf.currentTerm
	isLeader = (rf.status == LEADER)
	rf.unLock()

	return term, isLeader
}

func (rf *Raft) getLastEntry() Entry {
	return rf.logs[len(rf.logs)-1]
}

//
// save Raft's persistent state to stable storage,
// where it can later be retrieved after a crash and restart.
// see paper's Figure 2 for a description of what should be persistent.
//
func (rf *Raft) persist() {
	// w := new(bytes.Buffer)
	// e := gob.NewEncoder(w)
	// e.Encode(rf.currentTerm)
	// e.Encode(rf.votedFor)
	// e.Encode(rf.logs)
	// data := w.Bytes()
	// rf.persister.SaveRaftState(data)
}

//
// restore previously persisted state.
//
func (rf *Raft) readPersist(data []byte) {

	// if data == nil || len(data) < 1 {
	// 	// bootstrap without any state?
	// 	return
	// }
	// r := bytes.NewBuffer(data)
	// d := gob.NewDecoder(r)
	// d.Decode(&rf.currentTerm)
	// d.Decode(&rf.votedFor)
	// d.Decode(&rf.logs)
}

//
// 	预投票阶段，请求其他节点给自己投票，但是还需要 Committed证明 阶段
//
type PreRequestVoteArgs struct {
	Term                  int // Candidate Term
	CandidateId           int // Candidate Id
	LastCommittedLogTerm  int // Candidate当前最后一个已经 Committed 日志项的 Term（可能伪造）
	LastCommittedLogIndex int // Candidate当前最后一个已经 Committed 日志项的 Index（可能伪造）
}

//
// 预投票响应，包含待证明日志项的 Index 和 Term。
//
type PreRequestVoteReply struct {
	Success                       bool
	ReceiverId                    int
	ReceiverLastCommittedLogTerm  int
	ReceiverLastCommittedLogIndex int
}

//
// example RequestVote RPC arguments structure.
// field names must start with capital letters!
//
type RequestVoteArgs struct {
	Term        int
	CandidateId int

	// Committed 证明
	LastCommittedLogTerm  int    // 对方节点最后一个已经 Committed 日志项的 Term
	LastCommittedLogIndex int    // 对方节点最后一个已经 Committed 日志项的 Index
	LastCommittedLogHash  string // 对方节点最后一个已经 Committed 日志项的 Hash
}

//
// example RequestVote RPC reply structure.
// field names must start with capital letters!
//
type RequestVoteReply struct {
	Term        int
	VoteGranted bool
}

type AppendEntriesArgs struct {
	Term     int
	LeaderId int

	LogEntry Entry

	PrevLogIndex int
	PrevLogTerm  int
	IsHeartbeat  bool
}

type AppendEntriesReply struct {
	Term    int
	Success bool
}

type AppendEntriesCommitArgs struct {
	Term       int
	PeerId     int    // 确认消息发送者Id
	Signature  []byte // PeerId 的数字签名，用于防止拜占庭节点伪造确认消息
	EntryIndex int
	EntryTerm  int
	EntryHash  string
}

type AppendEntriesCommitReply struct {
	Term    int
	Success bool
}

type InstallSnapshotArgs struct {
	Term     int
	LeaderId int

	LastIncludedIndex int
	LastIncludedTerm  int

	Data []byte
}

type InstallSnapshotReply struct {
	Term    int
	Success bool
}

func (args *AppendEntriesArgs) String() string {
	return fmt.Sprintf("Term: %d, LeaderId: %d, entry: %v, PrevLogIndex: %d, PrevLogTerm: %d",
		args.Term, args.LeaderId, args.LogEntry, args.PrevLogIndex, args.PrevLogTerm)
}

func (args *PreRequestVoteArgs) String() string {
	return fmt.Sprintf("CandidateId: %d, LastCommittedLogTerm: %d, LastCommittedLogIndex: %d",
		args.CandidateId, args.LastCommittedLogTerm, args.LastCommittedLogIndex)
}

func (args *RequestVoteArgs) String() string {
	return fmt.Sprintf("CandidateId: %d, Term: %d, LastCommittedLogTerm: %d, LastCommittedLogIndex: %d, LastCommittedLogHash: %s",
		args.CandidateId, args.Term, args.LastCommittedLogTerm, args.LastCommittedLogIndex, args.LastCommittedLogHash)
}

func (rf *Raft) PreRequestVoteArgs(args PreRequestVoteArgs, reply *PreRequestVoteReply) {
	rf.lock()
	defer rf.unLock()

	if args.Term < rf.currentTerm {
		reply.Success = false
		return
	}

	if args.Term > rf.currentTerm {
		// if rf.status == LEADER {
		// 	fmt.Printf("PreRequestVoteArgs - LEADER [%d] becomes FOLLOWER, Current Time: %v\n", rf.me, time.Now().UnixNano()/1000000)
		// }
		rf.status = FOLLOWER
		rf.votedFor = -1
		rf.voteCount = 0
	}

	rf.currentTerm = args.Term

	commitIndex := rf.commitIndex
	if commitIndex >= len(rf.logs) {
		rf.commitIndex = len(rf.logs) - 1
	}
	commitTerm := rf.logs[rf.commitIndex].Term

	if rf.commitIndex == 0 || commitTerm == 0 {
		reply.Success = true
		reply.ReceiverId = rf.me
		reply.ReceiverLastCommittedLogTerm = 0
		reply.ReceiverLastCommittedLogIndex = 0
		return
	}

	//fmt.Printf("server: %d, commitTerm: %d, commitIndex: %d args.Term: %d\n",
	//	rf.me, commitTerm, rf.commitIndex, args.Term)

	if (args.LastCommittedLogTerm > rf.logs[rf.commitIndex].Term) ||
		(args.LastCommittedLogTerm == rf.logs[rf.commitIndex].Term && args.LastCommittedLogIndex >= rf.commitIndex) {
		// fmt.Printf("serve %d received PreRequestVoteArgs: %v\n", rf.me, args)
		reply.ReceiverId = rf.me
		reply.ReceiverLastCommittedLogTerm = rf.logs[rf.commitIndex].Term
		reply.ReceiverLastCommittedLogIndex = rf.commitIndex
		reply.Success = true
		return
	}

	reply.Success = false

}

func (rf *Raft) sendPreRquestVoteArgs(server int, args PreRequestVoteArgs, reply *PreRequestVoteReply) bool {
	ok := rf.peers[server].Call("Raft.PreRequestVoteArgs", args, reply)

	if !ok {
		return ok
	}

	if reply.Success {
		// if reply.ReceiverLastCommittedLogTerm == 0 || reply.ReceiverLastCommittedLogIndex == 0 {
		requestVoteArgs := RequestVoteArgs{}
		requestVoteArgs.CandidateId = rf.me
		requestVoteArgs.Term = rf.currentTerm
		requestVoteArgs.LastCommittedLogTerm = reply.ReceiverLastCommittedLogTerm
		requestVoteArgs.LastCommittedLogIndex = reply.ReceiverLastCommittedLogIndex
		requestVoteArgs.LastCommittedLogHash, _ = SHA256(rf.logs[reply.ReceiverLastCommittedLogIndex])

		requestVoteReply := RequestVoteReply{}

		go func() {
			// fmt.Printf("Server [%d] begin sendRequestVote: Current Time: %v\n", rf.me, time.Now().UnixNano()/1000000)
			rf.sendRequestVote(reply.ReceiverId, requestVoteArgs, &requestVoteReply)
		}()
		//}
	}
	// TODO
	return ok
}

//
// example RequestVote RPC handler.
//
func (rf *Raft) RequestVote(args RequestVoteArgs, reply *RequestVoteReply) {
	rf.lock()
	defer rf.unLock()
	defer rf.persist()
	fmt.Printf("Server [%d] received RequestVote from [%d], currentTerm: %d, args.term: %d, LastCommittedLogTerm: %d, LastCommittedLogIndex: %d\n",
		rf.me, args.CandidateId, rf.currentTerm, args.Term, args.LastCommittedLogTerm, args.LastCommittedLogIndex)

	// 如果投票请求的 Term 比 currentTerm 小，直接返回拒绝投票。
	if args.Term < rf.currentTerm {
		reply.Term = rf.currentTerm
		reply.VoteGranted = false
		return
	}

	// 如果投票消息中的 Term 大于 CurrentTerm，变更状态为 Follower。
	// 如果投票消息中的 Term 等于 CurrentTerm，需要比较 Index 来决定是否投票。
	if args.Term > rf.currentTerm {
		rf.currentTerm = args.Term
		rf.status = FOLLOWER
		rf.votedFor = -1
	}

	reply.Term = args.Term

	// 如果已经给其他的 Candidate 投过票了，直接返回拒绝投票。
	if rf.votedFor != -1 && rf.votedFor != args.CandidateId {
		reply.VoteGranted = false
		reply.Term = args.Term
		return
	}

	//////////////////// Begin Committed 证明 ////////////////////
	if rf.commitIndex == 0 {
		reply.Term = args.Term
		reply.VoteGranted = true
		rf.votedFor = args.CandidateId
		rf.grantVote <- true
		return
	}

	hash, _ := SHA256(rf.logs[rf.commitIndex])
	if args.LastCommittedLogTerm == rf.logs[rf.commitIndex].Term &&
		args.LastCommittedLogIndex == rf.commitIndex && args.LastCommittedLogHash == hash {
		// fmt.Printf("Server: %d vote for [%d]\n", rf.me, args.CandidateId)
		reply.Term = args.Term
		reply.VoteGranted = true
		rf.votedFor = args.CandidateId
		rf.grantVote <- true
		return
	} else {
		reply.Term = args.Term
		reply.VoteGranted = false
	}
	//////////////////// End Committed 证明 ////////////////////

	// // 当前节点最后一个日志项的 Index 和 Term。
	// receiverLastIndex := rf.getLastEntry().Index
	// receiverLastTerm := rf.getLastEntry().Term

	// // 如果投票请求中的 Term 比当前节点 LastTerm 大，投票。
	// // 如果 Term 相同，但是请求中的 Index 更大，投票。
	// if (args.LastLogTerm > receiverLastTerm) ||
	// 	(args.LastLogTerm == receiverLastTerm && args.LastLogIndex >= receiverLastIndex) {
	// 	rf.grantVote <- true // 一旦投了票，重置计时器：向 grantVote Channel 中写入一条消息
	// 	rf.votedFor = args.CandidateId
	// 	reply.VoteGranted = true
	// } else {
	// 	reply.VoteGranted = false
	// }

}

func (rf *Raft) sendRequestVote(server int, args RequestVoteArgs, reply *RequestVoteReply) bool {
	// 发送 RPC 请求。
	ok := rf.peers[server].Call("Raft.RequestVote", args, reply)
	rf.lock()
	defer rf.unLock()
	//fmt.Printf("%d send request to %d\n", rf.me, index)

	// 如果请求发送失败，直接返回。
	if !ok {
		return ok
	}

	// 如果当前节点不是出于 Candidate 状态，直接返回
	if rf.status != CANDIDATE {
		return ok
	}

	// 如果 RPC 请求中的 Term 不等于 currentTerm，直接返回。
	if args.Term != rf.currentTerm {
		return ok
	}

	// 如果收到的响应中的 Term 大于 currentTerm，更新 currentTerm，并变更状态为 Follower。
	if reply.Term > rf.currentTerm {
		rf.currentTerm = reply.Term
		rf.votedFor = -1 // 重置 votedFor，表示自己还没投过票
		rf.status = FOLLOWER
		rf.persist()
		return ok
	}

	// 如果成功收到投票
	if reply.VoteGranted && rf.status == CANDIDATE {
		rf.voteCount++
		// 如果收到大部分节点的投票，向 beLeader Channel 中发送一个消息
		if rf.voteCount >= quorum {
			//println(strconv.Itoa(rf.me) + " get leader ----------")
			rf.beLeader <- true
		}
	}
	return ok
}

func (rf *Raft) AppendEntries(args AppendEntriesArgs, reply *AppendEntriesReply) {
	rf.lock()
	defer rf.unLock()

	if !args.IsHeartbeat {
		// fmt.Printf("AppendEntries - server: %d, entry: %v\n", rf.me, args.LogEntry)

		if !verifySignature(GetBytes(args.LogEntry.Command), args.LogEntry.Signature) {
			rf.status = CANDIDATE
			rf.getHeartBeat <- true // Change status immediately
			fmt.Printf("FOLLOWER %d becomes CANDIDATE..., Current Time: %v\n", rf.me, time.Now().UnixNano()/1000000)
			return
		}
	}

	if args.Term < rf.currentTerm {
		rf.status = CANDIDATE
		rf.voteCount = 0
		fmt.Printf("AppendEntries - Server [%d] becomes CANDIDATE.\n", rf.me)
		reply.Success = false
		reply.Term = rf.currentTerm
		return
	}

	if args.Term > rf.currentTerm {
		rf.status = FOLLOWER
		fmt.Printf("AppendEntries - Server [%d] becomes FOLLOWER, Current Time: %v\n", rf.me, time.Now().UnixNano()/1000000)
		rf.votedFor = -1
		rf.currentTerm = args.Term
		reply.Term = args.Term
	}

	rf.getHeartBeat <- true
	rf.lastHeartbeat = time.Now()

	if args.IsHeartbeat {
		reply.Success = true
		reply.Term = args.Term
		return
	}

	lastLogIndex := rf.getLastEntry().Index
	lastLogTerm := rf.getLastEntry().Term

	if args.PrevLogIndex == lastLogIndex && args.PrevLogTerm == lastLogTerm {
		rf.logs = append(rf.logs, args.LogEntry)
		AddedTime = time.Now().UnixNano() / 1000000
		reply.Success = true
		reply.Term = args.Term

		// fmt.Printf("Server [%d] peers length: %d\n", rf.me, len(rf.peers))
		// broadcastAppendEntriesCommit(rf, args.LogEntry.Index, args.LogEntry.Term)
		for i := range rf.peers {
			if i != rf.me && rf.status != CANDIDATE {
				appendEntriesCommitArgs := AppendEntriesCommitArgs{Term: rf.currentTerm, PeerId: rf.me, EntryIndex: args.LogEntry.Index, EntryTerm: args.LogEntry.Term}
				appendEntriesCommitReply := AppendEntriesCommitReply{}
				go func(server int) {
					rf.sendAppendEntriesCommit(server, appendEntriesCommitArgs, &appendEntriesCommitReply)
					// fmt.Printf("Server [%d] wants to send append entries commit[index: %d] to server [%d]\n", rf.me, appendEntriesCommitArgs.EntryIndex, server)
				}(i)
			}
		}

	} else {
		reply.Success = false
		reply.Term = args.Term
	}

}

func (rf *Raft) sendAppendEntries(server int, args AppendEntriesArgs, reply *AppendEntriesReply) bool {
	// For Byzantine Fault Tolerant Test
	// if server == 5 {
	// 	args.LogEntry.Command = 666
	// }

	ok := rf.peers[server].Call("Raft.AppendEntries", args, reply)
	rf.lock()
	defer rf.unLock()

	// 如果 RPC 请求失败，直接返回。
	if !ok {
		return ok
	}

	// 如果当前节点不是 Leader，直接返回。
	if rf.status != LEADER {
		return ok
	}

	// 如果 args 中的 Term 不等于 currentTerm，直接返回。
	if args.Term != rf.currentTerm {
		return ok
	}

	if !args.IsHeartbeat {
		//fmt.Printf("%d send entry to %d \n", rf.me, server)
	}

	// 如果响应中的 Term 大于 Leader 的 currentTerm，Leader 设置 Term，并变更状态为 Follower，然后返回。
	if reply.Term > rf.currentTerm {
		rf.currentTerm = reply.Term
		rf.status = FOLLOWER
		rf.votedFor = -1
		rf.persist()
		return ok
	}

	return ok
}

func (rf *Raft) AppendEntriesCommit(args AppendEntriesCommitArgs, reply *AppendEntriesCommitReply) {

	rf.lock()
	defer rf.unLock()

	// fmt.Printf("AppendEntriesCommit - Server[%d], index: %d, sender: %d\n", rf.me, args.EntryIndex, args.PeerId)

	key := AppendEntriesCommitKey{Term: args.EntryTerm, Index: args.EntryIndex}
	_, ok := rf.m[key]
	if ok == true {
		rf.m[key] += 1
	} else {
		rf.m[key] = 1
	}

	// fmt.Printf("Server [%d] rf.m[key]: %d\n", rf.me, rf.m[key])
	// if rf.m[key] >= quorum {
	// 	fmt.Printf("Server[%d], key.Index = %d, commitIndex = %d\n", rf.me, key.Index, rf.commitIndex)
	// }
	if rf.m[key] >= quorum && key.Index > rf.commitIndex {
		rf.commitIndex = key.Index
		// fmt.Printf("Server [%d]: Index %d Committed.\n", rf.me, key.Index)
		rf.timeToCommit <- true
		CommittedTime = time.Now().UnixNano() / 1000000
		rf.WriteLineToFile(fmt.Sprintf("%d\n", CommittedTime-AddedTime))
	}

	reply.Success = true
}

func (rf *Raft) sendAppendEntriesCommit(server int, args AppendEntriesCommitArgs, reply *AppendEntriesCommitReply) bool {
	ok := rf.peers[server].Call("Raft.AppendEntriesCommit", args, reply)

	// 如果 RPC 请求失败，直接返回。
	if !ok {
		return ok
	}

	return ok
}

func (rf *Raft) Start(command interface{}, sig []byte) (int, int, bool) {
	// fmt.Println("signature:", sig)
	index := -1
	term := -1
	isLeader := true

	rf.lock()
	if rf.status == LEADER {
		index = rf.getLastEntry().Index + 1
		entry := Entry{index, rf.currentTerm, command, sig}
		rf.logs = append(rf.logs, entry)
		AddedTime = time.Now().UnixNano() / 1000000

		for i := range rf.peers {
			if i != rf.me && rf.status != CANDIDATE {
				appendEntriesCommitArgs := AppendEntriesCommitArgs{Term: rf.currentTerm, PeerId: rf.me, EntryIndex: entry.Index, EntryTerm: entry.Term}
				appendEntriesCommitReply := AppendEntriesCommitReply{}
				go func(server int) {
					rf.sendAppendEntriesCommit(server, appendEntriesCommitArgs, &appendEntriesCommitReply)
					// fmt.Printf("Server [%d] wants to send append entries commit[index: %d] to server [%d]\n", rf.me, appendEntriesCommitArgs.EntryIndex, server)
				}(i)
			}
		}

		rf.persist()

	} else {
		isLeader = false
	}

	term = rf.currentTerm
	rf.unLock()
	return index, term, isLeader
}

func election(rf *Raft) {
	rf.lock()
	// 自增 currentTerm
	rf.currentTerm++
	// 给自己投一票
	rf.votedFor = rf.me
	rf.persist()
	rf.unLock()

	go func() {
		broadcastPreRequestVote(rf)
	}()
}

func broadcastPreRequestVote(rf *Raft) {
	rf.lock()
	defer rf.unLock()
	for i := range rf.peers {
		if i != rf.me && rf.status == CANDIDATE {
			commitIndex := rf.commitIndex
			if commitIndex >= len(rf.logs) {
				rf.commitIndex = len(rf.logs) - 1
			}
			// fmt.Printf("Server [%d]: commitIndex: %d, log length: %d\n", rf.me, commitIndex, len(rf.logs))
			args := PreRequestVoteArgs{Term: rf.currentTerm, CandidateId: rf.me, LastCommittedLogTerm: rf.logs[rf.commitIndex].Term, LastCommittedLogIndex: rf.commitIndex}
			reply := PreRequestVoteReply{}
			go func(server int) {
				rf.sendPreRquestVoteArgs(server, args, &reply)
			}(i)
		}
	}
}

func broadcastRequestVote(rf *Raft) {
	rf.lock()
	defer rf.unLock()
	rf.voteCount = 1
	for i := range rf.peers {
		if i != rf.me && rf.status == CANDIDATE {
			requestVoteArgs := RequestVoteArgs{Term: rf.currentTerm, CandidateId: rf.me,
				LastCommittedLogIndex: rf.getLastCommittedLog().Index, LastCommittedLogTerm: rf.getLastCommittedLog().Term}
			result := RequestVoteReply{}
			go func(server int) {
				rf.sendRequestVote(server, requestVoteArgs, &result)
			}(i)
		}
	}
}

func broadcastAppendEntries(rf *Raft) {
	rf.lock()
	defer rf.unLock()
	// fmt.Printf("broadcastAppendEntries - rf.logs: %v\n", rf.logs)
	for i := range rf.peers {
		if i != rf.me && rf.status == LEADER {

			appendEntriesArgs := AppendEntriesArgs{}
			appendEntriesArgs.Term = rf.currentTerm
			appendEntriesArgs.LeaderId = rf.me

			nextIndex := rf.nextIndex[i]

			if rf.getLastEntry().Index >= nextIndex {
				appendEntriesArgs.LogEntry = rf.logs[nextIndex]
				appendEntriesArgs.IsHeartbeat = false
				// fmt.Printf("leader: %d, receiver: %d, last log index: %d, nextIndex: %d, entry: %v\n", rf.me, i, rf.getLastEntry().Index, nextIndex, appendEntriesArgs.entry)
			} else {
				appendEntriesArgs.IsHeartbeat = true
			}

			// fmt.Printf("PrevLogIndex: %d, len(rf.logs): %d\n", nextIndex-1, len(rf.logs))
			appendEntriesArgs.PrevLogIndex = nextIndex - 1
			appendEntriesArgs.PrevLogTerm = rf.logs[appendEntriesArgs.PrevLogIndex].Term

			result := AppendEntriesReply{}

			go func(server int) {
				// fmt.Printf("broadcastAppendEntries - LEADER [%d] send AppendEntries to FOLLOWER [%d]\n", rf.me, server)
				rf.sendAppendEntries(server, appendEntriesArgs, &result)
				if result.Success {
					if !appendEntriesArgs.IsHeartbeat {
						rf.nextIndex[server] = rf.nextIndex[server] + 1
					}
				} else {
					if rf.nextIndex[server] > 1 {
						rf.nextIndex[server] = rf.nextIndex[server] - 1
					}
				}
			}(i)
		}
	}
}

func broadcastAppendEntriesCommit(rf *Raft, index int, term int) {
	// rf.lock()
	// defer rf.unLock()
	for i := range rf.peers {
		if i != rf.me && rf.status != CANDIDATE {
			args := AppendEntriesCommitArgs{Term: rf.currentTerm, PeerId: rf.me, EntryIndex: index, EntryTerm: term}
			reply := AppendEntriesCommitReply{}
			go func(server int) {
				rf.sendAppendEntriesCommit(server, args, &reply)
			}(i)
		}
	}
}

func Make(peers []*labrpc.ClientEnd, me int, persister *Persister, applyCh chan ApplyMsg) *Raft {
	rf := &Raft{}
	rf.peers = peers
	rf.persister = persister
	rf.me = me

	// Your initialization code here (2A, 2B, 2C).
	rf.currentTerm = 0
	rf.votedFor = -1

	rf.status = FOLLOWER
	rf.beLeader = make(chan bool, 1)
	rf.timeToCommit = make(chan bool, 1)
	rf.grantVote = make(chan bool, 1)
	rf.getHeartBeat = make(chan bool, 1)
	rf.getCommit = make(chan bool, 1)
	rf.getAppendEntriesCommit = make(chan bool, 1)

	rf.m = make(map[AppendEntriesCommitKey]int)

	rf.logs = make([]Entry, 1)
	rf.commitIndex = 0
	rf.lastApplied = 0
	//rf.nextIndex = make([]int, 0)
	//rf.matchIndex = [len(rf.peers)]int{}

	rf.lastHeartbeat = time.Now()

	rf.applyCh = applyCh

	var filename = fmt.Sprintf("%d_data", rf.me)
	//if _, err := os.Stat(filename); err == nil {
	if err := os.Remove(filename); err != nil {
		// panic(err)
	}
	// }

	if _, err := os.Create(filename); err != nil {
		panic(err)
	}

	// initialize from state persisted before a crash
	rf.readPersist(persister.ReadRaftState())
	go func(rf *Raft) {
		for {
			switch rf.status {
			case FOLLOWER:
				select {
				case <-time.After(getRandomExpireTime()):
					rf.lock()
					rf.votedFor = -1
					fmt.Printf("Server [%d] Timeout....., Status: %v, Current Time: %v\n", rf.me, rf.status, time.Now().UnixNano()/1000000)

					now := time.Now()
					fmt.Printf("Server [%d] heartbeat gap: %v\n", rf.me, now.Sub(rf.lastHeartbeat))
					rf.lastHeartbeat = now

					rf.status = CANDIDATE
					rf.persist()
					rf.unLock()
				case <-rf.getHeartBeat:
				// 一旦投了票，重置计时器
				case <-rf.grantVote:
				}
			case LEADER:
				broadcastAppendEntries(rf)
				time.Sleep(time.Duration(HEARTBEAT_TIME) * time.Millisecond)
			case CANDIDATE:
				election(rf)
				select {
				// Candidate 在经过 election timeout 后还是 Candidate。
				case <-time.After(getRandomExpireTime()):
				case <-rf.getHeartBeat:
					rf.lock()
					rf.status = FOLLOWER
					rf.votedFor = -1
					rf.persist()
					rf.unLock()
				case <-rf.beLeader:
					rf.lock()
					rf.status = LEADER
					fmt.Printf("Server [%d] becomes LEADER... Current Time: %v\n", rf.me, time.Now().UnixNano()/1000000)
					rf.nextIndex = make([]int, len(rf.peers))

					for i := 0; i < len(rf.peers); i++ {
						rf.nextIndex[i] = rf.getLastEntry().Index + 1
					}

					rf.unLock()
				}
			}
		}

	}(rf)

	go func(rf *Raft) {
		for {
			select {
			case <-rf.timeToCommit:
				rf.lock()

				for rf.lastApplied < rf.commitIndex && rf.lastApplied+1 < len(rf.logs) {
					//println(lastAppliedIndex)
					rf.applyCommand(rf.logs[rf.lastApplied+1])
				}
				rf.unLock()
			}
		}
	}(rf)

	// go func(rf *Raft) {
	// 	for {
	// 		select {
	// 		case <-rf.getAppendEntriesCommit:
	// 			fmt.Println("...")
	// 		}
	// 	}
	// }()

	return rf
}

func (rf *Raft) Kill() {
	// Your code here, if desired.
}

func getRandomExpireTime() time.Duration {
	return time.Duration(rand.Int63n(300-150)+150) * time.Millisecond
	// return time.Duration(200) * time.Millisecond
}

func (rf *Raft) lock() {
	rf.mu.Lock()
}

func (rf *Raft) unLock() {
	rf.mu.Unlock()
}

func (rf *Raft) applyCommand(entry Entry) {
	applyCh := ApplyMsg{}
	applyCh.Index = entry.Index
	applyCh.Command = entry.Command
	rf.applyCh <- applyCh
	rf.lastApplied++
	rf.sendResponseToClient()
}

func (rf *Raft) getLastCommittedLog() Entry {
	rf.lock()
	defer rf.unLock()
	return rf.logs[rf.commitIndex]
}

func (rf *Raft) containsLogEntry(term int, index int) bool {
	for i := range rf.logs {
		entry := rf.logs[i]
		if entry.Term == term && entry.Index == index {
			return true
		}
	}
	return false
}

func GetBytes(key interface{}) []byte {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(key)
	if err != nil {
		panic(err)
	}
	return buf.Bytes()
}

func (rf *Raft) sendResponseToClient() {
	// TODO
	rf.ReceiveResponse(rf.me)
}

func (rf *Raft) WriteLineToFile(line string) {
	var filename = fmt.Sprintf("%d_data", rf.me)

	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	if _, err := f.WriteString(line); err != nil {
		panic(err)
	}

}

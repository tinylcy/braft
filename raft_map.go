package raft

import (
	"fmt"
	"sort"
)

type AppendEntriesCommitKey struct {
	Term  int
	Index int
	Hash  string
}

type CommitKeySlice []AppendEntriesCommitKey

func (s CommitKeySlice) Len() int {
	return len(s)
}

func (s CommitKeySlice) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func (s CommitKeySlice) Less(i, j int) bool {
	if s[i].Term < s[j].Term {
		return true
	}
	if s[i].Index < s[j].Index {
		return true
	}

	return s[i].Hash < s[j].Hash
}

func (key AppendEntriesCommitKey) less(elem AppendEntriesCommitKey) bool {
	if key.Term < elem.Term {
		return true
	}
	if key.Index < elem.Index {
		return true
	}
	return key.Hash < elem.Hash
}

func PrintSortedMap(m map[AppendEntriesCommitKey]int) {
	var keys []AppendEntriesCommitKey
	for k := range m {
		keys = append(keys, k)
	}
	sort.Sort(CommitKeySlice(keys))
	for _, k := range keys {
		fmt.Printf("[%v-%v] ", k, m[k])
	}
	fmt.Println("\n")
}

func GetMapSortedKeys(m map[AppendEntriesCommitKey]int) []AppendEntriesCommitKey {
	var keys []AppendEntriesCommitKey
	for k := range m {
		keys = append(keys, k)
	}
	sort.Sort(CommitKeySlice(keys))
	return keys
}

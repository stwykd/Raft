package raft

import (
	"bytes"
	"encoding/gob"
	"io/ioutil"
	"log"
)

func (rf *Raft) encode() []byte {
	w := new(bytes.Buffer)
	e := gob.NewEncoder(w)
	if e.Encode(rf.currentTerm) != nil ||
		e.Encode(rf.votedFor) != nil ||
		e.Encode(rf.logIndex) != nil ||
		e.Encode(rf.commitIndex) != nil ||
		e.Encode(rf.lastApplied) != nil ||
		e.Encode(rf.log) != nil {
		log.Fatalf("[%v] error while marshaling raft state", rf.me)
	}

	data := w.Bytes()
	return data
}

func (rf *Raft) decode(data []byte) {
	if data == nil || len(data) == 0 {
		return
	}
	d := gob.NewDecoder(bytes.NewBuffer(data))
	currentTerm, votedFor, logIndex, commitIndex, lastApplied := 0, 0, 0, 0, 0
	if d.Decode(&currentTerm) != nil ||
		d.Decode(&votedFor) != nil ||
		d.Decode(&logIndex) != nil ||
		d.Decode(&commitIndex) != nil ||
		d.Decode(&lastApplied) != nil ||
		d.Decode(&rf.log) != nil {
		log.Fatalf("[%v] error while marshaling raft state", rf.me)
	}
	rf.currentTerm, rf.votedFor, rf.logIndex, rf.commitIndex, rf.lastApplied =
		currentTerm, votedFor, logIndex, commitIndex, lastApplied
}

func (rf *Raft) persist() {
	err := ioutil.WriteFile("raft_state", rf.encode(), 0777)
	if err != nil {
		log.Fatalf("[%v] error while persisting raft state: %s", rf.me, err.Error())
	}
	log.Printf("[%v] state persisted to disk", rf.me)
}

func (rf *Raft) recover() {
	data, err := ioutil.ReadFile("raft_state")
	if err != nil {
		log.Fatalf("[%v] error while recovering raft state: %s", rf.me, err.Error())
	}
	rf.decode(data)
	log.Printf("[%v] state recovered from disk", rf.me)
}

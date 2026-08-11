package raft

import (
	"RaftLike/config"
	"RaftLike/pkg/logger"
	"context"
	"sync"
	"sync/atomic"
	"time"
)

type Raft struct {
	state           atomic.Int32 // 0=Follower, 1=Candidate, 2=Leader
	logger          *logger.Logger
	config          *config.Config
	heartbeater     *Heartbeater
	electioner      *Electioner
	heartbeatCancel context.CancelFunc
	electionCancel  context.CancelFunc
	baseCtx         context.Context
	wg              sync.WaitGroup
}

type State int

const (
	Follower State = iota
	Candidate
	Leader
)

func NewRaft(state State, logger *logger.Logger, config *config.Config, context context.Context) *Raft {
	raft := &Raft{atomic.Int32{}, logger, config, nil, nil, nil, nil, context, sync.WaitGroup{}}

	switch state {
	case Follower:
		raft.becomeFollower()
	case Candidate:
		raft.becomeCandidate()
	case Leader:
		raft.becomeLeader()
	default:
		logger.Error("Unknown state")
	}

	return raft
}

func (r *Raft) becomeLeader() {
	r.state.Store(int32(Leader))
	if r.heartbeater != nil {
		r.heartbeater.Stop(r.heartbeatCancel)
		r.wg.Done()
	}
	r.heartbeater = NewHeartbeater(r.logger, time.Duration(r.config.RAFT.Heartbet)*time.Millisecond)
	ctx, cancel := context.WithCancel(r.baseCtx)
	r.heartbeatCancel = cancel
	go r.heartbeater.Run(ctx)
	r.wg.Add(1)
}
func (r *Raft) becomeFollower() {
	r.state.Store(int32(Follower))
	if r.heartbeater != nil {
		r.heartbeater.Stop(r.heartbeatCancel)
		r.wg.Done()
	}
	r.electioner = NewElectioner(r.logger, time.Duration(r.config.RAFT.Election)*time.Millisecond)
	ctx, cancel := context.WithCancel(r.baseCtx)
	r.electionCancel = cancel
	go r.electioner.Run(ctx)
	r.wg.Add(1)
}
func (r *Raft) becomeCandidate() {
	r.state.Store(int32(Candidate))
	if r.heartbeater != nil {
		r.heartbeater.Stop(r.heartbeatCancel)
		r.wg.Done()
	}
	if r.electioner != nil {
		r.electioner.Stop(r.electionCancel)
		r.wg.Done()
	}
}
func (r *Raft) Stop() {
	if r.heartbeater != nil {
		r.heartbeater.Stop(r.heartbeatCancel)
		r.wg.Done()
	}
	if r.electioner != nil {
		r.electioner.Stop(r.electionCancel)
		r.wg.Done()
	}
	r.wg.Wait()
}

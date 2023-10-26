package ece428_mp2

import (
	"ece428_mp2/internal/network/code"
	"time"
)

type GossipGMConfig struct {
	ListenPort  int
	THeartbeat  time.Duration
	TFail       time.Duration
	TCleanup    time.Duration
	OutPerRound int
	Mode        code.RunMode

	// TODO: Support Suspicion Mode
	DropRate int // rate%
	TSuspect time.Duration
	TConfirm time.Duration
}

func GetDefaultGossipGMConfig() *GossipGMConfig {
	config := GossipGMConfig{
		ListenPort:  10088,
		THeartbeat:  1500 * time.Millisecond,
		TFail:       5 * time.Second,
		TCleanup:    6 * time.Second,
		OutPerRound: 4,
		Mode:        code.PureGossip,
	}

	return &config
}

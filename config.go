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
		THeartbeat:  2 * time.Second,
		TFail:       6 * time.Second,
		TCleanup:    5 * time.Second,
		OutPerRound: 3,
		Mode:        code.PureGossip,
	}

	return &config
}

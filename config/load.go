package config

import (
	"ece428_mp2/pkg/network/code"
	"fmt"
	"github.com/spf13/viper"
	"os"
	"time"
)

const (
	configFileName = "gossip"
	configFilePath = "./config"

	ListenPortKeyName               = "listen_port"
	IntroducerKeyName               = "introducer"
	THeartbeatKeyName               = "t_heartbeat"
	TFailKeyName                    = "t_fail"
	TCleanupKeyName                 = "t_cleanup"
	NumOfGossipPerRoundKeyName      = "fan_out"
	EnableSuspicionByDefaultKeyName = "default_suspicion"
	TSuspectKeyName                 = "t_suspect"
	TConfirmKeyName                 = "t_confirm"
	ServerListKeyName               = "server_list"
	DropRateKeyName                 = "drop_rate"
)

func MustLoadGossipFDConfig() {
	viper.SetConfigName(configFileName)
	viper.AddConfigPath(configFilePath)

	if err := viper.ReadInConfig(); err != nil {
		panic("can not load config of client")
	}

	if err := CheckClientConfig(); err != nil {
		panic(err)
	}
}

func CheckClientConfig() error {
	if !viper.IsSet(ListenPortKeyName) ||
		!viper.IsSet(IntroducerKeyName) ||
		!viper.IsSet(THeartbeatKeyName) ||
		!viper.IsSet(TFailKeyName) ||
		!viper.IsSet(TCleanupKeyName) ||
		!viper.IsSet(NumOfGossipPerRoundKeyName) ||
		!viper.IsSet(EnableSuspicionByDefaultKeyName) ||
		!viper.IsSet(TSuspectKeyName) ||
		!viper.IsSet(TConfirmKeyName) ||
		!viper.IsSet(ServerListKeyName) ||
		!viper.IsSet(DropRateKeyName) {
		return fmt.Errorf("missing config")
	}

	return nil
}

func GetIntroducerHost() string {
	return viper.GetString(IntroducerKeyName)
}

func GetListenPort() int {
	return viper.GetInt(ListenPortKeyName)
}

func GetTHeartbeat() time.Duration {
	return time.Second * time.Duration(viper.GetInt(THeartbeatKeyName))
}

func GetTFail() time.Duration {
	return time.Second * time.Duration(viper.GetInt(TFailKeyName))
}

func GetTCleanup() time.Duration {
	return time.Second * time.Duration(viper.GetInt(TCleanupKeyName))
}

func GetNumOfGossipPerRound() int {
	return viper.GetInt(NumOfGossipPerRoundKeyName)
}

func GetSelfHostName() string {
	hostname, err := os.Hostname()
	if err != nil {
		return ""
	}
	return hostname
	// hostname format fa23-cs425-48XX.cs.illinois.edu
}

func GetDefaultRunMode() code.RunMode {
	enableSuspicionByDefault := viper.GetBool(EnableSuspicionByDefaultKeyName)
	if enableSuspicionByDefault {
		return code.GossipWithSuspicion
	}
	return code.PureGossip
}

func GetTSuspect() time.Duration {
	return time.Second * time.Duration(viper.GetInt(TSuspectKeyName))
}

func GetTConfirm() time.Duration {
	return time.Second * time.Duration(viper.GetInt(TConfirmKeyName))
}

func GetServerList() []string {
	return viper.GetStringSlice(ServerListKeyName)
}

func GetDropRate() int {
	return viper.GetInt(DropRateKeyName)
}

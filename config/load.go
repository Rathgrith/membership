package config

import (
	"fmt"
	"github.com/spf13/viper"
	"time"
)

const (
	configFileName = "gossip"
	configFilePath = "./config"

	ListenPortKeyName = "listen_port"
	IntroducerKeyName = "introducer"
	THeartbeatKeyName = "t_heartbeat"
	TFailKeyName      = "t_fail"
	TCleanupKeyName   = "t_cleanup"
)

func MustLoadGossipConfig() {
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
		!viper.IsSet(TCleanupKeyName) {
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

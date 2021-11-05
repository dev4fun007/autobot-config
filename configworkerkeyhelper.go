package config

import "github.com/dev4fun007/autobot-common"

const (
	KeySeparator = "_"
)

func GetConfigWorkerKey(name string, strategyType common.StrategyType) string {
	return name + KeySeparator + string(strategyType)
}

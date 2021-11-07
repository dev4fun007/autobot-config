package config

import (
	"github.com/dev4fun007/autobot-common"
	"github.com/rs/zerolog/log"
	"sync"
)

const (
	RegistryServiceTag = "RegistryService"
)

type WorkerRegistryService struct {
	*sync.Mutex
	workerMap       map[string]common.Worker
	activeWorkerMap map[string]common.Worker
}

func NewWorkerRegistryService() WorkerRegistryService {
	return WorkerRegistryService{
		Mutex:           &sync.Mutex{},
		workerMap:       make(map[string]common.Worker),
		activeWorkerMap: make(map[string]common.Worker),
	}
}

func (receiver WorkerRegistryService) GetActiveWorkers() map[string][]common.Worker {
	marketWorkerMap := make(map[string][]common.Worker)
	for _, val := range receiver.activeWorkerMap {
		if l, ok := marketWorkerMap[val.GetBaseConfig().Market]; ok {
			l = append(l, val)
			marketWorkerMap[val.GetBaseConfig().Market] = l
		} else {
			l = make([]common.Worker, 0, 8)
			l = append(l, val)
			marketWorkerMap[val.GetBaseConfig().Market] = l
		}
	}
	return marketWorkerMap
}

func (receiver WorkerRegistryService) GetRegisteredWorker(configName string, strategyType common.StrategyType) common.Worker {
	key := GetConfigWorkerKey(configName, strategyType)
	return receiver.workerMap[key]
}

/**
Register commonconfig and it's data channel
*/
func (receiver WorkerRegistryService) RegisterConfigWorker(worker common.Worker, strategyType common.StrategyType) {
	receiver.Lock()
	defer receiver.Unlock()
	key := GetConfigWorkerKey(worker.GetBaseConfig().Name, strategyType)
	receiver.workerMap[key] = worker
	// add worker to active worker map if it is active
	if worker.GetBaseConfig().IsActive {
		receiver.activeWorkerMap[key] = worker
	}
	log.Debug().Str(common.LogComponent, RegistryServiceTag).
		Str("config-name", worker.GetBaseConfig().Name).
		Str("strategy-type", string(strategyType)).
		Msg("config worker registered")
}

/**
Update commonconfig and it's data channel in the registry
*/
func (receiver WorkerRegistryService) UpdateConfigWorkerRegistry(worker common.Worker, strategyType common.StrategyType) {
	receiver.Lock()
	defer receiver.Unlock()
	key := GetConfigWorkerKey(worker.GetBaseConfig().Name, strategyType)
	receiver.workerMap[key] = worker

	if worker.GetBaseConfig().IsActive {
		// add worker to active worker map if it is active
		receiver.activeWorkerMap[key] = worker
	} else {
		// delete from the other map - not checking as delete is a no-op, why waste a get check?
		delete(receiver.activeWorkerMap, key)
	}
	log.Debug().Str(common.LogComponent, RegistryServiceTag).
		Str("config-name", worker.GetBaseConfig().Name).
		Str("strategy-type", string(strategyType)).
		Msg("config worker updated")
}

/**
Delete commonconfig and it's data channel from the registry
*/
func (receiver WorkerRegistryService) RemoveConfigWorkerFromRegistry(name string, strategyType common.StrategyType) {
	receiver.Lock()
	defer receiver.Unlock()
	key := GetConfigWorkerKey(name, strategyType)
	if w, ok := receiver.workerMap[key]; ok {
		w.Close()
		w = nil
	}
	delete(receiver.activeWorkerMap, key)
	delete(receiver.workerMap, key)
	log.Debug().Str(common.LogComponent, RegistryServiceTag).
		Str("config-name", name).
		Str("strategy-type", string(strategyType)).
		Msg("config worker removed")
}

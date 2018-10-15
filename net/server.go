package sys

import (
	"sync"

	"bitbucket.org/MatthiasLeuth/chat_server/internal/sys/systruct"
	"bitbucket.org/MatthiasLeuth/chat_server/internal/module/usermod"
	"bitbucket.org/MatthiasLeuth/chat_server/internal/module/lockmod"
	"bitbucket.org/MatthiasLeuth/chat_server/internal/module/usermod/user_badger"
	"bitbucket.org/MatthiasLeuth/chat_server/internal/module/lockmod/lock_badger"
)

var server *serverSetup
var serverCreateMutex = &sync.Mutex{}

type serverSetup struct {
	Config          *systruct.Config
	UserGateCreator usermod.GateCreator
	LogGateCreator  lockmod.GateCreator
}

// NewServerSetup creates a 'serverSetup'
func NewServerSetup(config *systruct.Config) *serverSetup {
	serverCreateMutex.Lock()
	defer serverCreateMutex.Unlock()

	if server != nil && config.OneServerSetup {
		return server
	}

	server = &serverSetup{
		Config:          config,
		UserGateCreator: usermodCreator(config),
		LogGateCreator:  lockmodCreator(config),
	}

	return server
}

func ServerSetup() *serverSetup {
	return server
}

func usermodCreator(config *systruct.Config) usermod.GateCreator {
	creator, err := user_badger.NewGateCreator(config.SessionTTL)
	if err != nil {
		panic(err)
	}
	return creator
}

func lockmodCreator(config *systruct.Config) lockmod.GateCreator {
	creator, err := lock_badger.NewGateCreator(config.LockExpireTime)
	if err != nil {
		panic(err)
	}

	return creator
}

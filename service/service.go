package service

import (
	"github.com/MGYOSBEL/janus/pkg/log"
)

func Run() {
	logger := log.New("ZapLogger", "debug")

	logger.Debugf("Hola Service")
}

package client

import (
	"Anthole/pkg/common"
	"github.com/sirupsen/logrus"
)

func AntClientStart(conf *common.AntHoleConfig) error {
	logrus.Info("Client is Starting...")
	clientWorker := StartClientWorkerInstance()
	err := clientWorker.Run(conf)
	if err != nil {
		return err
	}
	return nil
}

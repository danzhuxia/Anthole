package client

import "Anthole/pkg/common"

func AntClientStart(conf *common.AntHoleConfig) error {
	clientWorker := StartClientWorkerInstance()
	err := clientWorker.Run(conf)
	if err != nil {
		return err
	}
	return nil
}

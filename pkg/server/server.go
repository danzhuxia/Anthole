package server

// func AntServerStart(conf *common.AntHoleConfig) (err error) {
// 	// 启动worker
// 	for i, item := range conf.Client.Services {
// 		go func(id, port int) {
// 			CreateWorker(id, port)
// 		}(i, item.RemotePort)
// 	}

// 	master := StartMasterInstance()
// 	err = master.Run(conf.Server.Port)
// 	return
// }

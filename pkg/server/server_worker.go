package server

import (
	"Anthole/pkg/common"
	"fmt"
	"net"
	"strconv"
	"sync"

	"github.com/sirupsen/logrus"
)

// const (
// 	optionNameConfigFilePath = "config-dir"
// )

// socker池
var RequestConnPool sync.Map

func CreateWorker(index int, port int) error {
	listen, err := net.Listen("tcp", strconv.Itoa(port))
	if err != nil {
		return err
	}
	defer listen.Close()
	for {
		if requestConn, err := listen.Accept(); err == nil {
			go func(serviceId int, conn net.Conn) {
				reqId := common.GenerateID()
				sendOpenCommand(reqId, uint16(serviceId))
				poolKey := fmt.Sprintf("%d_%d", serviceId, reqId)
				RequestConnPool.Store(poolKey, conn)
				buf := make([]byte, common.TransmissionDataLength)
				for {
					//网络数据流读入 buffer
					count, err := conn.Read(buf)
					//数据读尽、读取错误 关闭 socket 连接
					if count == 0 || err != nil {
						logrus.Debug("request socket断开")
						conn.Close()
						RequestConnPool.Delete(poolKey)
						// 通知client断开socket
						sendCloseCommand(reqId, uint16(serviceId))
						break
					}

					// 将读取到的buffer写到Anthole Client
					data := common.Transmission{
						RequestId:  reqId,
						Data:       buf,
						DataLength: uint16(count),
						ServiceId:  uint16(serviceId),
					}
					_, err = StartMasterInstance().SendDataToClient(data)
					// 写入client错误，关闭调request socket
					if err != nil {
						logrus.Debug("request socket断开")
						requestConn.Close()
						RequestConnPool.Delete(poolKey)
						sendCloseCommand(reqId, uint16(serviceId))
					}
				}
			}(index, requestConn)
		} else {
			logrus.Errorf("accept error:%v", err)
			return fmt.Errorf("accept error:%s", err.Error())
			// log.Fatalf("accept error:%v", err)
		}
	}
}

func CloseRequestSocket(poolKey string) {
	conn, ok := RequestConnPool.Load(poolKey)
	if ok {
		(conn.(net.Conn)).Close()
		RequestConnPool.Delete(poolKey)
	}
}

func sendDataToRequest(data common.Transmission) {
	poolKey := fmt.Sprintf("%d_%d", data.ServiceId, data.RequestId)
	conn, ok := RequestConnPool.Load(poolKey)
	if ok {
		logrus.Println("发送数据给Request")
		(conn.(net.Conn)).Write(data.GetData())
	}
}

func sendCloseCommand(reqId uint64, serviceId uint16) {
	data := []byte(fmt.Sprintf("%d,%d", reqId, serviceId))
	tdata := common.Transmission{
		RequestId:  common.CloseSocket,
		ServiceId:  common.CloseSocket,
		DataLength: uint16(len(data)),
	}

	paddingLen := common.TransmissionDataLength - len(data)

	padding := make([]byte, paddingLen)

	tdata.Data = append(data, padding...)

	StartMasterInstance().SendDataToClient(tdata)
}

//通知client打开本地socket
func sendOpenCommand(reqId uint64, serviceId uint16) {
	data := []byte(fmt.Sprintf("%d,%d", reqId, serviceId))
	transmissonData := common.Transmission{
		RequestId:  reqId,
		ServiceId:  serviceId,
		DataLength: uint16(len(data)),
		Data:       data,
	}
	paddingLen := common.TransmissionDataLength - len(data)

	padding := make([]byte, paddingLen)

	transmissonData.Data = append(data, padding...)

	StartMasterInstance().SendDataToClient(transmissonData)
}

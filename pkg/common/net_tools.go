package common

import (
	"fmt"
	"net"

	"github.com/sirupsen/logrus"
)

func GetDataFromConn(conn net.Conn, dataLength int) ([]byte, error) {
	totalDataLen := 0
	totalData := make([]byte, dataLength)
	for {
		//读取 conn 中的数据
		buf := make([]byte, dataLength-totalDataLen)
		cnt, err := conn.Read(buf)
		logrus.Debugf("读取数据长度 %d", cnt)

		if cnt == 0 || err != nil {
			logrus.Info("与服务端连接断开")
			return nil, fmt.Errorf("与服务端通讯错误:%v", err.Error())
		}
		totalData = append(totalData[:totalDataLen], buf...)
		totalDataLen += cnt
		if totalDataLen == dataLength {
			break
		}
	}
	return totalData, nil
}

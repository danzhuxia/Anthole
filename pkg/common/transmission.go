package common

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/binary"

	"github.com/sirupsen/logrus"
)

const (
	CloseSocket = 1 // 断开本地socket通知
	OpenSocket  = 3 // 打开本地socket通知
	VerifyKey   = 2 // 验证客户端秘钥

	TransmissionPackageLength = 4112 // 8字节的RequestID， 2字节的serviceid  2字节dataLength  4096字节的数据包长度
	TransmissionDataLength    = 4096 // 传输的数据包长度
)

var iv = []byte("a01020c73554d64!")

type Transmission struct {
	RequestId  uint64
	ServiceId  uint16
	DataLength uint16
	Data       []byte
}

func (t *Transmission) GetData() []byte {
	return t.Data[:t.DataLength]
}

func (t *Transmission) ConvertToByte() []byte {
	requestId := make([]byte, 8)
	binary.LittleEndian.PutUint64(requestId, uint64(t.RequestId))

	serviceId := make([]byte, 2)
	binary.LittleEndian.PutUint16(serviceId, uint16(t.ServiceId))

	dataLength := make([]byte, 2)
	binary.LittleEndian.PutUint16(dataLength, uint16(t.DataLength))

	ret := append(requestId, serviceId...)
	ret = append(ret, dataLength...)
	ret = append(ret, t.Data...)

	// 校验数据包长度
	if len(t.Data) < TransmissionDataLength {
		padding := make([]byte, TransmissionDataLength-len(t.Data))
		ret = append(ret, padding...)
	}

	// 数据包加密
	encrData, _ := Encrypt(ret, []byte(AntConf.Common.Token))

	return encrData
}

func TodoTransmission(encrData []byte) Transmission {

	// 数据包解密
	data, err := Decrypt(encrData, []byte(AntConf.Common.Token))
	if err != nil {
		logrus.Errorf("数据包解密失败！ %v", err)
	}
	requestId := binary.LittleEndian.Uint64(data[0:8])
	serviceId := binary.LittleEndian.Uint16(data[8:10])
	dataLength := binary.LittleEndian.Uint16(data[10:12])

	tdata := data[12:]
	ret := Transmission{
		RequestId:  requestId,
		ServiceId:  serviceId,
		DataLength: dataLength,
		Data:       tdata,
	}

	return ret
}

//数据加密
func Encrypt(text []byte, key []byte) ([]byte, error) {

	md5Ctx := md5.New()
	md5Ctx.Write(key)
	cipherStr := md5Ctx.Sum(nil)

	//生成cipher.Block 数据块
	block, err := aes.NewCipher(cipherStr)
	if err != nil {
		logrus.Println("错误 -" + err.Error())
		return []byte{}, err
	}
	//填充内容，如果不足16位字符
	blockSize := block.BlockSize()
	originData := pad(text, blockSize)
	//加密方式
	blockMode := cipher.NewCBCEncrypter(block, iv)
	//加密，输出到[]byte数组
	crypted := make([]byte, len(originData))
	blockMode.CryptBlocks(crypted, originData)
	return crypted, nil
}

func pad(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

//数据解密
func Decrypt(data []byte, key []byte) ([]byte, error) {

	md5Ctx := md5.New()
	md5Ctx.Write(key)
	cipherStr := md5Ctx.Sum(nil)

	//生成密码数据块cipher.Block
	block, _ := aes.NewCipher(cipherStr)
	//解密模式
	blockMode := cipher.NewCBCDecrypter(block, iv)
	//输出到[]byte数组
	origin_data := make([]byte, len(data))
	blockMode.CryptBlocks(origin_data, data)
	//去除填充,并返回
	return unpad(origin_data), nil
}

func unpad(ciphertext []byte) []byte {
	length := len(ciphertext)
	//去掉最后一次的padding
	unpadding := int(ciphertext[length-1])
	return ciphertext[:(length - unpadding)]
}

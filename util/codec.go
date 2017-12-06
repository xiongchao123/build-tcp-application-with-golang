package util

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
)

func Encode(message string) ([]byte, error) {
	// 读取消息的长度
	var length int32 = int32(len(message))
	var pkg *bytes.Buffer = new(bytes.Buffer)
	// 写入消息头
	err := binary.Write(pkg, binary.LittleEndian, length)
	if err != nil {
		return nil, err
	}
	// 写入消息实体
	err = binary.Write(pkg, binary.LittleEndian, []byte(message))
	if err != nil {
		return nil, err
	}

	return pkg.Bytes(), nil
}

func Decode(reader *bufio.Reader, package_length_offset int, order binary.ByteOrder) ([]byte, error) {
	// 读取消息的长度
	lengthByte, err := reader.Peek(package_length_offset + 4)
	if(err != nil){
		return []byte{}, err
	}
	var length uint32
	bytesBuffer := lengthByte[package_length_offset:]
	lengthBuff := bytes.NewBuffer(bytesBuffer)
	binary.Read(lengthBuff, order, &length)
	fmt.Println("buffer:",uint32(reader.Buffered()))
	fmt.Println("lenght:",length)
	if uint32(reader.Buffered()) < length {
		return []byte{}, nil
	}
	// 读取消息真正的内容
	pack := make([]byte, length)
	_, err = reader.Read(pack)
	if err != nil {
		return []byte{}, err
	}
	return pack, nil
}

func DecodeBak(reader *bufio.Reader) (string, error) {
	// 读取消息的长度
	lengthByte, _ := reader.Peek(4)
	lengthBuff := bytes.NewBuffer(lengthByte)
	var length int32
	err := binary.Read(lengthBuff, binary.LittleEndian, &length)
	if err != nil {
		return "", err
	}
	if int32(reader.Buffered()) < length+4 {
		return "", err
	}

	// 读取消息真正的内容
	pack := make([]byte, int(4+length))
	_, err = reader.Read(pack)
	if err != nil {
		return "", err
	}
	return string(pack[4:]), nil
}

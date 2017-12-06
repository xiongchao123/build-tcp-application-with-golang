package util

import (
	"strings"
	"strconv"
	"bytes"
	"encoding/binary"
	"reflect"
	"unsafe"
)

//解包
func (p *Protocol) UnPack(msg []byte) []interface{} {
	la := len(p.Format)
	ret := make([]interface{}, la)
	if la > 0 {
		for i := 0; i < la; i++ {
			switch {
			case strings.Contains(p.Format[i], "a"): //字符串空白以 NULL 字符填满
				num, _ := strconv.Atoi(strings.TrimLeft(p.Format[i], "a"))
				ret[i] = string(msg[0:num])
				msg = msg[num:len(msg)]
			case strings.Contains(p.Format[i], "A"): //字符串空白以 SPACE 字符 (空格) 填满
				num, _ := strconv.Atoi(strings.TrimLeft(p.Format[i], "A"))
				ret[i] = string(msg[0:num])
				msg = msg[num:len(msg)]
			case p.Format[i] == "c": //有符号字符
				ret[i] = BytesToInt(msg[0:1], "int8", GetOrder())
				msg = msg[1:len(msg)]
			case p.Format[i] == "C": //无符号字符
				ret[i] = BytesToInt(msg[0:1], "uint8", GetOrder())
				msg = msg[1:len(msg)]
			case p.Format[i] == "s": //有符号短整数 (16位，主机字节序)
				ret[i] = BytesToInt(msg[0:2], "int16", GetOrder())
				msg = msg[2:len(msg)]
			case p.Format[i] == "S": //无符号短整数 (16位，主机字节序)
				ret[i] = BytesToInt(msg[0:2], "uint16", GetOrder())
				msg = msg[2:len(msg)]
			case p.Format[i] == "n": //无符号短整数 (16位，大端字节序)
				ret[i] = BytesToInt(msg[0:2], "int16",binary.BigEndian)
				msg = msg[2:len(msg)]
			case p.Format[i] == "v": //无符号短整数 (16位，小端字节序)
				ret[i] = BytesToInt(msg[0:2], "uint16", binary.LittleEndian)
				msg = msg[2:len(msg)]
			case p.Format[i] == "l": //有符号长整数 (32位，主机字节序)
				ret[i] = BytesToInt(msg[0:4], "int16", GetOrder())
				msg = msg[4:len(msg)]
			case p.Format[i] == "L": //无符号长整数 (32位，主机字节序)
				ret[i] = BytesToInt(msg[0:4], "uint16", GetOrder())
				msg = msg[4:len(msg)]
			}
		}
	}
	return ret
}

func BytesToHex(b []byte){

}

func BytesToInt(b []byte, t string, order binary.ByteOrder) interface{} {
	bytesBuffer := bytes.NewBuffer(b)
	switch t {
	case "int8":
		var data int8
		binary.Read(bytesBuffer, order, &data)
		return data
	case "uint8":
		var data uint8
		binary.Read(bytesBuffer, order, &data)
		return data
	case "int16":
		var data int16
		binary.Read(bytesBuffer, order, &data)
		return data
	case "uint16":
		var data uint16
		binary.Read(bytesBuffer, order, &data)
		return data
	case "int32":
		var data int32
		binary.Read(bytesBuffer, order, &data)
		return data
	case "uint32":
		var data uint32
		binary.Read(bytesBuffer, order, &data)
		return data
	default:
		return nil
	}
}

func BytesToString(b []byte) string {
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	sh := reflect.StringHeader{bh.Data, bh.Len}
	return *(*string)(unsafe.Pointer(&sh))
}



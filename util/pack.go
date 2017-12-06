package util

import (
	"bytes"
	"encoding/binary"
	"fmt"
	_ "os"
	"strconv"
	"strings"
)

type Protocol struct {
	Format []string
}

func HexToB(str string) ([]byte) {
	slen := len(str)
	bHex := make([]byte, len(str)/2)
	ii := 0
	for i := 0; i < len(str); i = i + 2 {
		if slen != 1 {
			ss := string(str[i]) + string(str[i+1])
			bt, _ := strconv.Atoi(ss)
			bHex[ii] = byte(bt)
			ii = ii + 1
			slen = slen - 2
		}
	}
	return bHex
}

func HexToTenB(str string) ([]byte) {
	slen := len(str)
	bHex := make([]byte, len(str)/2)
	ii := 0
	for i := 0; i < len(str); i = i + 2 {
		if slen != 1 {
			ss := string(str[i]) + string(str[i+1])
			bt, _ := strconv.ParseInt(ss, 16, 32)
			bHex[ii] = byte(bt)
			ii = ii + 1
			slen = slen - 2
		}
	}
	return bHex
}

func BytetoH(b []byte) (H string) {
	H = fmt.Sprintf("%x", b)
	return
}

//封包
func (p *Protocol) Pack(args ...interface{}) []byte {
	la := len(args)
	ls := len(p.Format)
	ret := []byte{}
	if ls > 0 && la > 0 && ls == la {
		for i := 0; i < ls; i++ {
			switch {
			case strings.Contains(p.Format[i], "a"): //将字符串空白以 NULL 字符填满
				num, _ := strconv.Atoi(strings.TrimLeft(p.Format[i], "a"))
				if (num == 0) {
					num = 1
				}
				ret = append(ret, []byte(fmt.Sprintf("%s%s", args[i].(string), strings.Repeat("\x00", num-len(args[i].(string)))))...)
				break
			case strings.Contains(p.Format[i], "A"): //将字符串空白以 SPACE 字符 (空格) 填满
				num, _ := strconv.Atoi(strings.TrimLeft(p.Format[i], "A"))
				if (num == 0) {
					num = 1
				}
				ret = append(ret, []byte(fmt.Sprintf("%s%s", args[i].(string), strings.Repeat("\x20", num-len(args[i].(string)))))...)
				break
			case p.Format[i] == "c": //有符号字符
				data := int8(args[i].(int))
				ret = append(ret, IntToBytes(data, GetOrder())...)
				break
			case p.Format[i] == "C": //无符号字符
				data := uint8(args[i].(int))
				ret = append(ret, IntToBytes(data, GetOrder())...)
				break
			case p.Format[i] == "s": //有符号短整数 (16位，主机字节序)
				data := int16(args[i].(int))
				ret = append(ret, IntToBytes(data, GetOrder())...)
				break
			case p.Format[i] == "S": //无符号短整数 (16位，主机字节序)
				data := uint16(args[i].(int))
				ret = append(ret, IntToBytes(data, GetOrder())...)
				break
			case p.Format[i] == "n": //无符号短整数 (16位, 大端字节序)
				data := uint16(args[i].(int))
				ret = append(ret, IntToBytes(data, binary.BigEndian)...)
				break
			case p.Format[i] == "v": // 无符号短整数 (16位, 小端字节序)
				data := uint16(args[i].(int))
				ret = append(ret, IntToBytes(data, binary.LittleEndian)...)
				break
			case p.Format[i] == "l": //有符号长整数 (32位，主机字节序)
				data := int32(args[i].(int))
				ret = append(ret, IntToBytes(data, GetOrder())...)
				break
			case p.Format[i] == "L": // 无符号长整数 (32位，主机字节序)
				data := uint32(args[i].(int))
				ret = append(ret, IntToBytes(data, GetOrder())...)
				break
			case p.Format[i] == "N": //无符号长整数 (32位, 大端字节序)
				data := uint32(args[i].(int))
				ret = append(ret, IntToBytes(data, binary.BigEndian)...)
				break
			case p.Format[i] == "V": //无符号长整数 (32位, 小端字节序)
				data := uint32(args[i].(int))
				ret = append(ret, IntToBytes(data, binary.LittleEndian)...)
				break
			}
		}
	}
	return ret
}

func IntToBytes(data interface{}, order binary.ByteOrder) []byte {
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, order, data)
	gbyte := bytesBuffer.Bytes()
	return gbyte
}


func (p *Protocol) Size() int {
	size := 0
	ls := len(p.Format)
	if ls > 0 {
		for i := 0; i < ls; i++ {
			if p.Format[i] == "H" {
				size = size + 2
			} else if p.Format[i] == "I" {
				size = size + 4
			} else if strings.Contains(p.Format[i], "s") {
				num, _ := strconv.Atoi(strings.TrimRight(p.Format[i], "s"))
				size = size + num
			}
		}
	}
	return size
}

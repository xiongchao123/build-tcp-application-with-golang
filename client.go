package main

import (
	"encoding/json"
	"bufio"
	"log"
	"fmt"
	"net"
	"bytes"
	"./util"
	"time"
	"encoding/binary"
	"hash/crc32"
)

var quitSemaphore = make(chan bool, 1)

type ReqBody struct {
	Code string `json:"code"`
	Test string `json:"test"`
}

func main() {
	var tcpAddr *net.TCPAddr
	tcpAddr, err := net.ResolveTCPAddr("tcp", "10.10.83.231:9503")
	if (err != nil ) {
		log.Panic(err.Error())
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if (err != nil ) {
		log.Panic(err.Error())
	}
	defer conn.Close()
	fmt.Println("connected!")

	go onMessageRecived(conn)

	var b bytes.Buffer
	body := ReqBody{"test", "hello"}
	req_body, err := json.Marshal(body)
	if (err != nil) {
		log.Fatal(err.Error())
	}
	p := new(util.Protocol)
	p.Format = []string{"C", "C", "C", "C", "V"}
	head_crc_af:=p.Pack(1, 5, 0, 0, 45+len(req_body))
	b.Write(head_crc_af)
	p.Format = []string{"V"}
	b.Write(p.Pack(int(crc32.ChecksumIEEE(head_crc_af))))
	uuid := util.RandStr(33)
	b.WriteString(uuid)
	b.Write(req_body)
	conn.Write(b.Bytes())
	//xintiao qingqiu
	go heartReq(conn,uuid)
	<-quitSemaphore
}

func heartReq(conn *net.TCPConn,uuid string){
	var b bytes.Buffer
	p := new(util.Protocol)
	p.Format = []string{"C", "C", "C", "C", "V"}
	head_crc_af:=p.Pack(1, 3, 0, 0, 45)
	p.Format = []string{"V"}
	b.Write(head_crc_af)
	b.Write(p.Pack(int(crc32.ChecksumIEEE(head_crc_af))))
	b.WriteString(uuid)
	tick1 := time.Tick(time.Second * 5)
	for {
		select {
		case <-tick1:
			conn.Write(b.Bytes())
		}
	}
}

func onMessageRecived(conn *net.TCPConn) {
	reader := bufio.NewReader(conn)
	for {
		fmt.Println("start:")
			msg, err := util.Decode(reader, 4, binary.LittleEndian)
			if err != nil {
				log.Println(err)
				quitSemaphore <- true
				break
			}
			if (0 != len(msg)) {
				fmt.Println("msg lenght:", len(msg))
				body:=(string(msg[12:]))
				fmt.Println(len(body))
				fmt.Println(body)
			}
		time.Sleep(time.Second)
	}
}

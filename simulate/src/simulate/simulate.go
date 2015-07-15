package main

import (
	"log"
	"misc/packet"
	"net"
	"os"
	"time"
)

var (
	seqid = uint32(0)
)

const (
	DEFAULT_AGENT_HOST = "127.0.0.1:8888"
)

func checkErr(err error) {
	if err != nil {
		log.Println(err)
		panic("error occured in protocol module")
	}
}
func main() {
	host := DEFAULT_AGENT_HOST
	if env := os.Getenv("AGENT_HOST"); env != "" {
		host = env
	}
	addr, err := net.ResolveTCPAddr("tcp", host)
	if err != nil {
		log.Println(err)
		os.Exit(-1)
	}
	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		log.Println(err)
		os.Exit(-1)
	}
	defer conn.Close()
	for i := 0; i < 100; i++ {
		send_proto(conn, Code["heart_beat_req"], nil)
	}
}

func send_proto(conn net.Conn, p int16, info interface{}) {
	seqid++
	payload := packet.Pack(p, info, nil)
	writer := packet.Writer()
	writer.WriteU16(uint16(len(payload)) + 4)
	writer.WriteU32(seqid)
	writer.WriteRawBytes(payload)
	conn.Write(writer.Data())
	time.Sleep(time.Second)
}

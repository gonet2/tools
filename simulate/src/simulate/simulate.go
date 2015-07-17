package main

import (
	"log"
	"math/rand"
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

	//heart_beat_req
	send_proto(conn, Code["heart_beat_req"], nil)

	//proto_ping_req
	p1 := auto_id{
		F_id: rand.Int31(),
	}
	send_proto(conn, Code["proto_ping_req"], p1)

	//get_seed_req
	p2 := seed_info{
		F_client_send_seed:    rand.Int31(),
		F_client_receive_seed: 0,
	}
	send_proto(conn, Code["get_seed_req"], p2)

	//user_login_req
	p3 := user_login_info{
		F_login_way:          0,
		F_open_udid:          "udid",
		F_client_certificate: "qwertyuiopasdfgh",
		F_client_version:     1,
		F_user_lang:          "en",
		F_app_id:             "com.yrhd.lovegame",
		F_os_version:         "android4.4",
		F_device_name:        "simulate",
		F_device_id:          "device_id",
		F_device_id_type:     1,
		F_login_ip:           "127.0.0.1",
	}
	send_proto(conn, Code["user_login_req"], p3)
}

func send_proto(conn net.Conn, p int16, info interface{}) {
	seqid++
	payload := packet.Pack(p, info, nil)
	writer := packet.Writer()
	writer.WriteU16(uint16(len(payload)) + 4)
	writer.WriteU32(seqid)
	writer.WriteRawBytes(payload)
	conn.Write(writer.Data())
	log.Printf("%#v", writer.Data())
	time.Sleep(time.Second)
}

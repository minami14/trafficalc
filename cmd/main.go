package main

import (
	"fmt"
	"github.com/jessevdk/go-flags"
	"os"
)

type options struct {
	Client   int `short:"c" long:"client" description:"Number of clients connecting to the same room at the same time" default:"100"`
	Trackers int `short:"t" long:"tracker" description:"Number of trackers per client including hmd and controller" default:"10"`
	Fps      int `short:"f" long:"fps" description:"Number of times to synchronize transforms per second" default:"90"`
}

const (
	_  = iota
	kb = 1 << (iota * 10)
	mb
	gb
)

const (
	tcp         = 20
	ip          = 20
	eth         = 20
	metadata    = tcp + ip + eth
	f32         = 4
	i32         = 4
	vector3     = f32 * 3
	quaternion  = f32 * 4
	objectId    = i32
	clientId    = i32
	transform   = vector3 + quaternion + objectId
	rpcTarget   = 1
	messageType = 1
	messageSize = 2
	bitsPerByte = 8
)

func main() {
	var opt options
	if _, err := flags.Parse(&opt); err != nil {
		os.Exit(1)
	}

	transforms := transform*opt.Trackers + 1
	payload := messageSize + rpcTarget + messageType + transforms
	data := metadata + payload
	bps := data * opt.Fps * bitsPerByte
	clientUpKbps := float64(bps) / float64(kb)
	serverDownMbps := float64(bps) * float64(opt.Client) / float64(mb)
	payload = messageSize + clientId + messageType + transforms
	data = metadata + payload
	bps = data * opt.Fps * opt.Client * bitsPerByte
	clientDownMbps := float64(bps) / float64(mb)
	serverUpGbps := float64(bps) * float64(opt.Client) / float64(gb)

	fmt.Printf("client down\t%.2f Mbps\n", clientDownMbps)
	fmt.Printf("client up\t%.2f Kbps\n", clientUpKbps)
	fmt.Printf("server down\t%.2f Mbps\n", serverDownMbps)
	fmt.Printf("server up\t%.2f Gbps\n", serverUpGbps)
}

package pub

import "net"

type Event struct {
	Name   string // 时间名称，比如创建，登陆、离线、聊天、心跳等
	Data   []byte
	Conn   *net.TCPConn
	Player *Player
}

const (
	EventNameAccept     = "accept"
	EventNameDisconnect = "disconnect"
	EventNameIO         = "io"
)

type EventProc func(e *Event)

func EventLoop() {
	for {
		select {
		case ev := <-event:
			dispatchEvent(ev)
		}
	}
}

func RegEvent(name string, proc EventProc) {
	callBakcs[name] = append(callBakcs[name], proc)
}

func PutEvent(ev *Event) {
	event <- ev
}

var event chan *Event
var callBakcs map[string][]EventProc

func init() {
	callBakcs = map[string][]EventProc{}
	event = make(chan *Event, 100)
}

func dispatchEvent(e *Event) {
	if procs, ok := callBakcs[e.Name]; ok {
		for _, v := range procs {
			v(e)
		}
	}
}

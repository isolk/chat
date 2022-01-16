package main

import (
	"chat/pub"
	"flag"
	"fmt"
	"log"
	"net"
)

func initConnection() {
	attr := flag.String("addr", "127.0.0.1:8888", "输入服务器地址,比如  127.0.0.1:8888")
	flag.Parse()
	raddr, err := net.ResolveTCPAddr("tcp", *attr)
	if err != nil {
		log.Fatalf("wrong addr err=%v", err)
	}
	con, err := net.DialTCP("tcp", nil, raddr)
	if err != nil {
		log.Fatalf("dail faild err=%v", err)
	}
	pl = &pub.Player{
		Con: con,
	}
	fmt.Printf("connected to %s,welcome!\n", con.RemoteAddr())
	go pl.Loop() // 开始监听服务器发来的消息,收到消息后，会放入全局的事件管理器中。
}

func initEvent() {
	// 客户端只有服务器发来的消息，比如请求的返回等
	pub.RegEvent(pub.EventNameIO, IoEventHandler)
	pub.RegEvent(pub.EventNameDisconnect, DisconnectHandler)
	go pub.EventLoop()
}

func main() {
	// 先初始化服务器连接
	initConnection()
	// 初始化事件管理器以及事件注册
	initEvent()
	// 开始循环准备接受消息
	inputLoop()
	fmt.Println("聊天结束")
}

package main

import (
	"bufio"
	"chat/pub"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

// 初始化事件框架，就是注册几个事件处理函数，基本大量的逻辑都是在 IOEventHandler中，也就是不同玩家的网络消息。
// 采用事件主要是基于以下考虑
// - 当玩家数量多了之后，大部分的开销实际上是在网络io上，一个玩家发消息，就需要同步给所有人，如果同时10个，100个人发的话，网络瓶颈很快就上去了
// - 事件框架虽然是单线程（协程）处理，其基本的cpu压力都来自于屏蔽词过滤，聊天内容不长的话，剩下的都没有什么压力，所以即便是单线程问题也不大。
// - 后面考虑扩展的话，可以将屏蔽词过滤这部分的逻辑拆分出去，再搞个协程处理，就可以承载更大量的聊天内容。
func initEvent() {
	pub.RegEvent(pub.EventNameAccept, NewAcceptHandler)
	pub.RegEvent(pub.EventNameIO, IOEventHandler)
	pub.RegEvent(pub.EventNameDisconnect, DisconnectHandler)
	go pub.EventLoop()
}

// 开始网络循环，此循环逻辑处理比较简单，就是不断获取新的连接，然后将其放入事件框架中。
func netLoop() {
	attr := flag.String("addr", "127.0.0.1:8888", "输入服务器地址,比如  127.0.0.1:8888")
	flag.Parse()

	laddr, err := net.ResolveTCPAddr("tcp", *attr)
	if err != nil {
		log.Fatalf("wrong addr err=%v\n", err)
	}

	listener, err := net.ListenTCP("tcp", laddr)
	if err != nil {
		log.Fatalf("failed to listen,err=%v\n", err)
	}

	for {
		con, err := listener.Accept()
		if err != nil {
			listener.Close()
			log.Fatalf("failed to accept connection,err=%v\n", err)
		}
		e := &pub.Event{Name: pub.EventNameAccept, Conn: con.(*net.TCPConn)}
		pub.PutEvent(e)
	}
}

//  初始化屏蔽词树
func initWordFilter() {
	f, err := os.Open("dic.txt")
	if err != nil {
		log.Fatalf("failed to init word filter,err=%v", err)
	}
	pub.RootTire = &pub.TireNode{}
	reader := bufio.NewReader(f)
	for {
		buf, pre, err := reader.ReadLine()
		if pre {
			continue
		}
		if err == nil {
			word := string(buf)
			invalid := false
			for i := 0; i < len(word); i++ {
				if !(word[i] <= 'Z' && word[i] >= 'A' ||
					word[i] <= 'z' && word[i] >= 'a') {
					invalid = true
				}
			}
			if !invalid {
				pub.RootTire.Insert(string(word), 0)
			}
		} else if err == io.EOF {
			break
		} else {
			log.Fatalf("failed to read dic file,err=%v", err)
		}
	}
	fmt.Println("word filter init ok ")
}

func main() {
	// 启动事件循环
	initEvent()
	// 初始化过滤器
	initWordFilter()
	// 开始监听网络事件
	netLoop()

	fmt.Println("服务器关闭")
}

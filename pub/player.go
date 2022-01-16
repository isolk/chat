package pub

import (
	"encoding/binary"
	"fmt"
	"net"
)

type Player struct {
	Name      string
	LoginTime int64
	Room      int

	Con *net.TCPConn
}

func (p *Player) Send(obj Message) {
	bytes, err := MessagePack(obj)
	if err != nil {
		fmt.Println("err pack", err)
		return
	}
	go func() {
		for n := 0; n < len(bytes); {
			n, err := p.Con.Write(bytes)
			if err != nil {
				e := &Event{
					Name:   EventNameDisconnect,
					Conn:   p.Con,
					Player: p,
				}
				PutEvent(e)
				return
			}
			bytes = bytes[n:]
		}
	}()
}

// 这块儿将每个链接上的client都当作player，即便还没有登陆，这样可以方便简化代码。
// 这是一个单独的网络协程，每个连接一个协程，可能会阻塞，但是不影响主协程。
func (p *Player) Loop() {
	buf := make([]byte, 1024*1024)

	// 消息格式
	// magic(chat),datalen(2),data(type2,data)
	for {
		readLenth := 0
		for readLenth < headerLen {
			n, err := p.Con.Read(buf[readLenth:headerLen])
			if err != nil {
				p.fireDisconnectEvent()
				return
			}
			readLenth += n
		}
		if magic := string(buf[0:len(magic)]); magic != magic {
			fmt.Printf("message magic error,magic=%s\n", magic)
			p.fireDisconnectEvent()
			return
		}
		dataLen := binary.BigEndian.Uint16(buf[magicLen:])
		for readLenth < headerLen+int(dataLen) {
			n, err := p.Con.Read(buf[readLenth : headerLen+int(dataLen)])
			if err != nil {
				p.fireDisconnectEvent()
				return
			}
			readLenth += n
		}
		// 有了完整的消息包了，构建事件。
		ev := &Event{
			Name:   EventNameIO,
			Player: p,
			Data:   buf[headerLen : headerLen+int(dataLen)],
		}
		PutEvent(ev)
	}
}

func (p *Player) fireDisconnectEvent() {
	e := &Event{
		Name:   EventNameDisconnect,
		Conn:   p.Con,
		Player: p,
	}
	PutEvent(e)
}

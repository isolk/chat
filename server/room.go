package main

import (
	"chat/pub"
	"time"
)

func init() {
	roomMgr.rooms[1] = &Room{id: 1, players: make(map[string]*pub.Player)}
}

type Room struct {
	id            int
	players       map[string]*pub.Player
	content       []pub.Content
	recentContens []pub.Word
}

// 底层采用了单开个协程去处理，因为大并发下存在阻塞的可能，所以不能同步去处理。事件框架协程必须保证不能阻塞。
// 
func (m *Room) SendToAll(req pub.Message) {
	for _, v := range m.players {
		v.Send(req)
	}
}

// 添加文本时，会判断是否超过50条，以及更新高频词库
func (r *Room) AppendContent(id int, playerName, conent string) bool {
	c := pub.Content{
		Data:     conent,
		Player:   playerName,
		ChatTime: time.Now().Unix(),
	}

	r.content = append(r.content, c)
	if len(r.content) > 50 {
		r.content = r.content[1:len(r.content)]
	}

	words := []string{}
	last := 0
	for i := 0; i < len(conent); i++ {
		if conent[i] == ' ' || conent[i] == ',' || conent[i] == '.' {
			if i != 0 && last <= i {
				words = append(words, conent[last:i])
				last = i + 1
			}
		} else if i == len(conent)-1 {
			words = append(words, conent[last:i+1])
		}
	}
	for _, v := range words {
		r.recentContens = append(r.recentContens, pub.Word{Time: c.ChatTime, Data: v})
	}
	return true
}

// o(n)的时间复杂度,o(n)的空间复杂度
func (r *Room) GetTopWord() string {
	deleteK := 0
	wantT := time.Now().Unix() - 60*10
	for i := 0; i < len(r.recentContens); i++ {
		if r.recentContens[i].Time > wantT {
			break
		}
		deleteK = i
	}
	r.recentContens = r.recentContens[deleteK:]
	if len(r.recentContens) == 0 {
		return ""
	}

	dic := make(map[string]int, len(r.recentContens))
	for _, v := range r.recentContens {
		dic[v.Data] += 1
	}
	maxK := ""
	max := 0
	for k, v := range dic {
		if v > max {
			max = v
			maxK = k
		}
	}
	return maxK
}

type RoomManager struct {
	rooms map[int]*Room
}

func (r *RoomManager) SearchContent(id int) []pub.Content {
	if r.rooms[id] == nil {
		return nil
	}
	return r.rooms[id].content
}

func (r *RoomManager) GetRoom(id int) *Room {
	return r.rooms[id]
}

func (r *RoomManager) AddPlayer(id int, pl *pub.Player) *Room {
	if r.rooms[id] == nil {
		r.rooms[id] = &Room{players: map[string]*pub.Player{}, id: id}
	}
	room := r.rooms[id]
	room.players[pl.Name] = pl
	return room
}

func (r *RoomManager) RemovePlayer(pl *pub.Player) {
	if r.rooms[pl.Room] == nil {
		return
	}

	room := r.rooms[pl.Room]
	delete(room.players, pl.Name)
	// 房间没人，就把房间释放掉，防止越累计越多
	// if len(room.players) == 0 {
	// 	delete(r.rooms, room.id)
	// }
}

var roomMgr = RoomManager{rooms: map[int]*Room{}}

package main

import (
	"chat/pub"
	"encoding/binary"
	"fmt"
	"net"
	"time"
)

var playerMg PlyaerManager

func init() {
	playerMg = PlyaerManager{
		loginPlayers: map[string]*pub.Player{},
		readyPlayers: map[net.Conn]*pub.Player{},
		names:        map[string]string{},
	}
}

const (
	salt = "oqBNUoChvNzLspYR76cc3htHuXImSK9h"
)

type PlyaerManager struct {
	loginPlayers map[string]*pub.Player
	readyPlayers map[net.Conn]*pub.Player
	names        map[string]string
}

func (m *PlyaerManager) Login(pl *pub.Player, req *pub.LoginRequest) {
	if m.IsOnline(pl) {
		return
	}

	fmt.Printf("player login. name=%s,pwd=%s\n", req.Name, req.PassWord)
	if req.Name == "" || req.PassWord == "" {
		return
	}

	reqPwd := pub.MD5String(req.PassWord + salt)
	if pwd, ok := m.names[req.Name]; ok && reqPwd != pwd {
		msg := &pub.LoginResponse{
			ErrCode: 1,
		}
		pl.Send(msg)
		return
	}

	//  初始化新登陆玩家状态
	pl.Name = req.Name
	pl.LoginTime = time.Now().Unix()
	pl.Room = 1
	m.names[pl.Name] = reqPwd
	m.loginPlayers[pl.Name] = pl
	delete(m.readyPlayers, pl.Con)

	//  更新房间信息
	roomMgr.AddPlayer(1, pl)

	msg := &pub.LoginResponse{
		Room: pl.Room,
	}
	room := roomMgr.GetRoom(pl.Room)
	for _, v := range room.players {
		msg.OnlinePlyaers = append(msg.OnlinePlyaers, v.Name)
	}
	msg.RecentContent = roomMgr.SearchContent(pl.Room)
	pl.Send(msg)
}

func (m *PlyaerManager) Chat(pl *pub.Player, req *pub.ChatRequest) {
	if !m.IsOnline(pl) {
		return
	}

	room := roomMgr.GetRoom(pl.Room)
	if room == nil {
		return
	}

	req.Content = pub.RootTire.ReplaceSentence(req.Content)
	room.AppendContent(pl.Room, pl.Name, req.Content)
	res := &pub.ChatSyncRequest{
		Player:  pl.Name,
		Content: req.Content,
		Time:    time.Now().Unix(),
	}
	room.SendToAll(res)
}

func (m *PlyaerManager) Select(pl *pub.Player, req *pub.SelectRequest) {
	if !m.IsOnline(pl) {
		return
	}

	roomMgr.RemovePlayer(pl)
	pl.Room = req.Room
	newRoom := roomMgr.AddPlayer(pl.Room, pl)

	msg := &pub.LoginResponse{
		Room: pl.Room,
	}
	for _, v := range newRoom.players {
		msg.OnlinePlyaers = append(msg.OnlinePlyaers, v.Name)
	}
	msg.RecentContent = roomMgr.SearchContent(pl.Room)
	pl.Send(msg)
	fmt.Printf("%s select room [%d]\n", pl.Name, pl.Room)
}

func (m *PlyaerManager) StatusRequests(pl *pub.Player, req *pub.StatusRequest) {
	res := &pub.StatusResponse{}
	if !m.IsOnline(pl) {
		res.Err = 1
	} else {
		res.LoginTime = pl.LoginTime
		res.Room = pl.Room
	}
	pl.Send(res)
}

func (m *PlyaerManager) PopularRequests(pl *pub.Player, req *pub.PopularRequest) {
	res := &pub.PopularResponse{Room: req.Room}
	room := roomMgr.GetRoom(req.Room)
	if room == nil {
		res.Err = 1
	} else {
		res.Word = room.GetTopWord()
	}
	pl.Send(res)
}

func (m *PlyaerManager) IsOnline(pl *pub.Player) bool {
	return m.loginPlayers[pl.Name] != nil
}

func IOEventHandler(e *pub.Event) {
	messageType := binary.BigEndian.Uint16(e.Data[:pub.MessageTypeSize])
	e.Data = e.Data[pub.MessageTypeSize:]
	switch messageType {
	case 1:
		req := &pub.LoginRequest{}
		pub.MessageUnpack(req, e.Data)
		playerMg.Login(e.Player, req)
	case 3:
		req := &pub.ChatRequest{}
		pub.MessageUnpack(req, e.Data)
		playerMg.Chat(e.Player, req)
	case 5:
		req := &pub.SelectRequest{}
		pub.MessageUnpack(req, e.Data)
		playerMg.Select(e.Player, req)
	case 6:
		req := &pub.StatusRequest{}
		pub.MessageUnpack(req, e.Data)
		playerMg.StatusRequests(e.Player, req)
	case 8:
		req := &pub.PopularRequest{}
		pub.MessageUnpack(req, e.Data)
		playerMg.PopularRequests(e.Player, req)
	}
}

func NewAcceptHandler(e *pub.Event) {
	p := &pub.Player{
		Con: e.Conn,
	}
	playerMg.readyPlayers[e.Conn] = p
	// 每个新连接都认为是一个玩家。此时处于未登陆状态。但是这个会持续不断接受信息了
	fmt.Printf("ip %s connect\n", e.Conn.RemoteAddr())
	go p.Loop()
}

func DisconnectHandler(e *pub.Event) {
	pl := e.Player
	if e.Player.Name != "" {
		fmt.Printf("player %s logout\n", e.Player.Name)
	} else {
		fmt.Printf("player %s logout\n", e.Conn.RemoteAddr())
	}

	delete(playerMg.loginPlayers, e.Player.Name)
	delete(playerMg.readyPlayers, e.Player.Con)
	e.Conn.Close()
	roomMgr.RemovePlayer(pl)
}

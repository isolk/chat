package main

import (
	"chat/pub"
	"encoding/binary"
	"fmt"
	"os"
	"time"
)

var pl *pub.Player

func IoEventHandler(e *pub.Event) {
	messageType := binary.BigEndian.Uint16(e.Data[:pub.MessageTypeSize])
	e.Data = e.Data[pub.MessageTypeSize:]
	switch messageType {
	case 2:
		req := &pub.LoginResponse{}
		pub.MessageUnpack(req, e.Data)
		LoginResponse(req)
	case 4:
		req := &pub.ChatSyncRequest{}
		pub.MessageUnpack(req, e.Data)
		ChatSync(req)
	case 7:
		req := &pub.StatusResponse{}
		pub.MessageUnpack(req, e.Data)
		StatusResponse(req)
	case 9:
		req := &pub.PopularResponse{}
		pub.MessageUnpack(req, e.Data)
		PopularResponse(req)
	}
}
func DisconnectHandler(e *pub.Event) {
	fmt.Println("服务器已关闭，退出聊天")
	os.Exit(0)
}

func LoginResponse(res *pub.LoginResponse) {
	if res.ErrCode == 1 {
		fmt.Printf("密码错误，请提供正确的密码。或者换个角色登陆。\n")
	} else {
		pl.Room = res.Room
		fmt.Printf("进入房间成功，房间号%d，当前房间在线人员%v\n", res.Room, res.OnlinePlyaers)
		for _, v := range res.RecentContent {
			fmt.Printf("%s(%d)[%v]:%s\n", v.Player, pl.Room, v.ChatTime, v.Data)
		}
	}
}

func ChatSync(res *pub.ChatSyncRequest) {
	if bench {
		totalChan <- 1
	}

	fmt.Printf("%s(%d)[%v]:%s\n", res.Player, pl.Room, res.Time, res.Content)
}

func StatusResponse(res *pub.StatusResponse) {
	if res.Err != 0 {
		fmt.Printf("该玩家%s不在线\n", res.Name)
		return
	}
	fmt.Printf("status: [loginTime:%d,onlineTime:%d,room:%d]:\n", res.LoginTime, time.Now().Unix()-res.LoginTime, res.Room)
}

func PopularResponse(res *pub.PopularResponse) {
	if res.Err != 0 {
		fmt.Printf("room[%d]不存在\n", res.Room)
	} else {
		fmt.Printf("popular_word[%d]:%s\n", res.Room, res.Word)
	}
}

package main

import (
	"bufio"
	"chat/pub"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

func inputLoop() {
	for {
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		text = text[0 : len(text)-1]
		inputs := strings.Split(text, " ")
		if len(inputs) != 0 {
			parseInput(text, inputs[0], inputs[1:]...)
		}
	}
}
func parseInput(text string, cmd string, params ...string) {
	switch cmd {
	case ":login":
		if len(params) == 2 {
			login(params[0], params[1])
		}
	case ":select":
		if len(params) == 1 {
			k := params[0]
			selectRoom(k)
		}
	case "/stats":
		if len(params) == 1 {
			status(params[0])
		}
	case "/popular":
		if len(params) == 1 {
			popular(params[0])
		}
	case ":bench":
		if len(params) == 2 {
			benchMark(params[0], params[1])
		}
	default:
		chat(text)
	}
}

func login(name, password string) {
	req := &pub.LoginRequest{
		Name:     name,
		PassWord: password,
	}
	pl.Send(req)
}

func chat(content string) {
	req := &pub.ChatRequest{
		Content: content,
	}
	pl.Send(req)
}

func selectRoom(name string) {
	id, err := strconv.Atoi(name)
	if err != nil || id <= 0 {
		fmt.Println("房间号必须大于0", err)
		return
	}
	req := &pub.SelectRequest{Room: id}
	pl.Send(req)
}

func status(name string) {
	req := &pub.StatusRequest{Name: name}
	pl.Send(req)
}

func popular(name string) {
	id, err := strconv.Atoi(name)
	if err != nil || id <= 0 {
		fmt.Println("房间号必须大于0", err)
		return
	}
	req := &pub.PopularRequest{Room: id}
	pl.Send(req)
}

var bench bool
var totalChan = make(chan int)

func benchMark(lenStr, totalStr string) {
	lenth, err := strconv.Atoi(lenStr)
	if err != nil {
		fmt.Println("len 错误", err)
		return
	}

	total, err := strconv.Atoi(totalStr)
	if err != nil {
		fmt.Println("total 错误", err)
		return
	}

	p := "ab gh ijk lmn op qr stu vwx yz"
	contntBytes := []byte{}
	for i := 0; i < lenth; i++ {
		ranI := rand.Int() % len(p)
		contntBytes = append(contntBytes, p[ranI])
	}
	req := &pub.ChatRequest{
		Content: string(contntBytes),
	}

	bench = true
	for i := 0; i < total; i++ {
		pl.Send(req)
	}

	begin := time.Now()
	i := total
	for i > 0 {
		select {
		case <-totalChan:
			i--
		}
	}
	bench = false
	timeCost := time.Since(begin).Milliseconds() + 1
	fmt.Printf("bench mark over,timeCost=%dms,sendCount=%d,qps=%d\n", timeCost, total, total/int(timeCost)*1000)
}

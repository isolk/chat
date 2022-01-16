package pub

type Content struct {
	Player   string
	Data     string
	ChatTime int64
}

type Word struct {
	Time int64
	Data string
}

// 登陆
type LoginRequest struct {
	Name     string
	PassWord string
}

func (*LoginRequest) MsgID() uint16 {
	return 1
}

// 登陆响应
type LoginResponse struct {
	ErrCode       int
	Room          int
	OnlinePlyaers []string
	RecentContent []Content
}

func (*LoginResponse) MsgID() uint16 {
	return 2
}

// 聊天
type ChatRequest struct {
	Content string
}

func (*ChatRequest) MsgID() uint16 {
	return 3
}

// 聊天同步
type ChatSyncRequest struct {
	Player  string
	Time    int64
	Content string
}

func (*ChatSyncRequest) MsgID() uint16 {
	return 4
}

// 选择房间
type SelectRequest struct {
	Room int
}

func (*SelectRequest) MsgID() uint16 {
	return 5
}

// 状态查询
type StatusRequest struct {
	Name string
}

func (*StatusRequest) MsgID() uint16 {
	return 6
}

type StatusResponse struct {
	Name      string
	Err       int
	LoginTime int64
	Room      int
}

func (*StatusResponse) MsgID() uint16 {
	return 7
}

// 流行词查询
type PopularRequest struct {
	Room int
}

func (*PopularRequest) MsgID() uint16 {
	return 8
}

type PopularResponse struct {
	Room int
	Err  int
	Word string
}

func (*PopularResponse) MsgID() uint16 {
	return 9
}

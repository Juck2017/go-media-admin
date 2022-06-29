package core

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
	"suitbim.com/go-media-admin/utils"
)

var (
	dispatcher *Dispatcher
)

// Dispatcher 管理端的主结构体，职责：
// 1. 维护与Dispatcher（remote media center）之间的信号来往（websocket协议）
// 2. 维护负责给Dispatcher推流（tcp）的client列表
type Dispatcher struct {
	Channels          []Channel           // 从数据库获取的通道列表
	ChannelMap        map[string]*Channel // 从数据库获取的通道映射信息
	ctx               context.Context     // 上下文
	cancel            context.CancelFunc  // 取消触发的方法
	viewersManager    *WSServer           // 观看管理服务
	signalManager     *WSServer           // 通信管理服务
	signalConnections []*websocket.Conn   // 通信连接情况
}

// 初始化一个dispatcher管理端
func NewDispatcher() *Dispatcher {
	disp := Dispatcher{}
	disp.signalConnections = make([]*websocket.Conn, 0)
	disp.Channels = LoadChannels()
	disp.ChannelMap = make(map[string]*Channel)
	for i := 0; i < len(disp.Channels); i++ {
		// 以通道编码为key把所有通道数据放入channelMap中
		disp.ChannelMap[disp.Channels[i].Code] = &disp.Channels[i]
	}
	// 上下文传递,用来控制通道状态
	disp.ctx, disp.cancel = context.WithCancel(context.Background())
	dispatcher = &disp
	return dispatcher
}

func (disp *Dispatcher) Start() {
	disp.signalManager = &WSServer{Port: G.Dispatcher.Signal.Port}
	disp.viewersManager = &WSServer{Port: G.Dispatcher.Viewer.Port}
	go disp.StartTCPServer()
	go disp.signalManager.Start(disp.OnPusherManagerEnter)
	disp.viewersManager.Start(disp.OnViewerEnter)
}

// StartTCPServer Dispatcher开始工作
func (disp *Dispatcher) StartTCPServer() {
	receiver := TCPServer{}
	receiver.Start(disp.ctx, disp)
}

// OnViewerEnter 处理来自前端流消费者的事件
func (disp *Dispatcher) OnViewerEnter(c *gin.Context) {
	// 升级get请求为websocket协议
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	chCode := c.Query("code")
	G.Logger.Infof("Viewer[%s] require channel [%s]", c.ClientIP(), chCode)
	ch, ok := disp.ChannelMap[chCode]
	var v *Viewer
	// 如果在数据库存在该通道
	if ok {
		// 将viewer添加值dispatch的对应通道的观看者列表中(登记)
		v = &Viewer{conn: ws, vid: uuid.Must(uuid.NewV4(), nil).String()}
		ch.Viewers = append(ch.Viewers, v)
		G.Logger.Infof("Channel[%s] remain [%d] Viewers", ch.Code, len(ch.Viewers))
		G.Logger.Infof("Channel[%s] holding Pusher[%v]", ch.Code, ch.pusher)
		if ch.pusher == nil {
			// 如果该通道目前还没有Pusher提供服务
			// 调用信号管理员发送推流消息给远端的Pusher管理员
			disp.SendSignal(ch, &Message{Sign: SIG_OPEN, Chan: ch.Code}, 1)
		}
	} else {
		// 对于不在已配置的通道范围内的viewer，直接断开连接，避免浪费服务器资源
		_ = ws.Close()
		return
	}
	// 最后需要给浏览器正常的响应（OK，101）；
	// 否则浏览器得到的Response Status会是Finished，可能会导致无限重发
	_ = ws.WriteMessage(1, []byte("connected"))

	// 监听Viewer断开连接
	if _, _, err = ws.ReadMessage(); err != nil {
		G.Logger.Errorf("Viewer[%s][%s] occur error [%v]", chCode, c.ClientIP(), err)
		ch.OnViewerLeaved(v, err)
	}
}

// OnPusherManagerEnter 处理与Pusher之间的信号交互
func (disp *Dispatcher) OnPusherManagerEnter(c *gin.Context) {
	// 升级get请求为websocket协议
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	//disp.signalManager.conn = ws
	disp.signalConnections = append(disp.signalConnections, ws)
}

// SendSignal 发送消息给Pusher Managers
func (disp *Dispatcher) SendSignal(channel *Channel, msg *Message, times int) {
	G.Logger.Infof("Channel[%s] Notify Pusher Manager do [%s] [%d]Times: ", channel.Code, msg.Sign, times)
	// 通知所有已跟Dispatcher取得联系的pusher manager
	for _, remotePusherManager := range disp.signalConnections {
		if err := remotePusherManager.WriteMessage(
			websocket.BinaryMessage,
			[]byte(utils.Struct2JsonString(*msg)),
		); err != nil {
			G.Logger.Error("Notify Pusher Manager Error: ", err)
		}
	}

	// 如果是关闭信号，可以直接中断推流连接，对端检测到错误后主动停止推流动作
	if msg.Sign == SIG_CLOS && channel.pusher != nil {
		_ = channel.pusher.Close()
	}
}

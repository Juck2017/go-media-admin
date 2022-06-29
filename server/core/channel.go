package core

import (
	"github.com/gorilla/websocket"
	"net"
)

type Viewer struct {
	conn *websocket.Conn
	vid  string
}

// Channel 通道
type Channel struct {
	Id      uint      `json:"id"`      // 通道id
	Code    string    `json:"code"`  // 通道编码
	Status  int       `json:"status"`// 通道状态
	Viewers []*Viewer // 观看者列表
	pusher  net.Conn  // tcp推流连接
}

// OnPushingData Pusher有流推送过来时
func (ch *Channel) OnPushingData(data []byte) {
	if len(ch.Viewers) == 0 {
		// 流没有了消费者，则断开与pusher之间的连接，使其停止推流，并且结束当前goroutine
		G.Logger.Errorf("Puser[%s][%s] has none Viewers, will Closing !", ch.pusher.RemoteAddr().String(), ch.Code)
		_ = ch.pusher.Close()
		// 将当前通道的push提供者置为空
		ch.pusher = nil
		return
	} else {
		// 否则将流推给每个前端浏览者viewer
		for _, viewer := range ch.Viewers {
			err := viewer.conn.WriteMessage(websocket.BinaryMessage, data)
			if err != nil {
				ch.OnViewerLeaved(viewer, err)
			}
		}
	}
}

// OnViewerLeaved 推流给Viewer失败时的回调
func (ch *Channel) OnViewerLeaved(v *Viewer, err error) {
	// 如果推流给viewer失败，则认为viewer已经关闭
	G.Logger.Infof("Viewer[%s] ch=[%s] leaved with error %v", v.conn.RemoteAddr(), ch.Code, err)
	// 从channel将对应的viewer删除
	var vIdx = -1
	for i := 0; i < len(ch.Viewers); i++ {
		if ch.Viewers[i].vid == v.vid {
			vIdx = i
			break
		}
	}
	if vIdx > -1 {
		ch.Viewers = append(ch.Viewers[:vIdx], ch.Viewers[vIdx+1:]...)
	}
	G.Logger.Infof("Channel[%s] remain [%d] Viewers", ch.Code, len(ch.Viewers))
	if len(ch.Viewers) == 0 {
		// 需要检测该viewer请求的通道是否还有其他viewer还在使用流
		// 如果该通道已经没有消费者了，就可以通知PusherManager关闭该通道的推流动作了
		dispatcher.SendSignal(ch, &Message{Sign: SIG_CLOS, Chan: ch.Code}, 1)
		ch.pusher = nil
	}
}

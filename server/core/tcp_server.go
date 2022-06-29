package core

import (
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"net"
)

// TCPServer 用于与Dispatcher建立TCP通道、控制推流
type TCPServer struct {
}

type Pusher struct {
	conn net.Conn
}

func (p *Pusher) Receive(buf []byte) int {
	n, err := p.conn.Read(buf)
	if err != nil {
		G.Logger.Error(err)
		return 0
	}
	return n
}

func (p *Pusher) RemoteAddr() net.Addr {
	return p.conn.RemoteAddr()
}

// 开启tcp服务端用来接收推送的流
func (svr *TCPServer) Start(ctx context.Context, disp *Dispatcher) {
	listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", G.Dispatcher.Stream.Port))
	if err != nil {
		G.Logger.Error(err)
		panic(err)
	}

	G.Logger.Infof("Receiver has started tcp server at port:%d", G.Dispatcher.Stream.Port)

	for {
		select {
		case <-ctx.Done():
			G.Logger.Error("TCPServer has been canceled.")
			return
		default:
		}
		conn, err := listener.Accept()
		if err != nil {
			G.Logger.Error("listener.Accept Error: ", err)
			continue
		}

		go OnPusherEnter(ctx, conn, disp) // 启动一个goroutine处理连接
	}
}

// 当有一个pusher进入时的处理
func OnPusherEnter(ctx context.Context, conn net.Conn, disp *Dispatcher) {
	G.Logger.Infof("Pusher[%s] enter", conn.RemoteAddr().String())
	defer conn.Close()
	var buf [512000]byte
	var mCh *Channel
	for {
		select {
		case <-ctx.Done():
			G.Logger.Error("TCPServer has been canneled.")
			return
		default:
		}

		//读取包头
		if _, err := io.ReadFull(conn, buf[:2]); err != nil {
			G.Logger.Infof("Pusher[%s] leaved", conn.RemoteAddr().String())
			_ = conn.Close()
			return
		}
		bodyLen := binary.LittleEndian.Uint16(buf[:2])
		// 读包体
		if _, err := io.ReadFull(conn, buf[:bodyLen]); err != nil {
			G.Logger.Infof("Pusher[%s] leaved", conn.RemoteAddr().String())
			_ = conn.Close()
			return
		}

		if bodyLen < 20 {
			// 包长度小于10，是告知通道号的
			chCode := BytesToString(buf[:bodyLen])
			ch, ok := disp.ChannelMap[chCode]
			if ok {
				ch.pusher = conn
				mCh = ch // 通过 pusher发来的通道号 匹配到 channel实例
				G.Logger.Infof("Puser[%s][%s] matched", conn.RemoteAddr().String(), chCode)
			} else {
				G.Logger.Errorf("Puser[%s][%s] can not matched", conn.RemoteAddr().String(), chCode)
			}
		} else {
			// 包长度大于10，才是实际的流媒体数据
			mCh.OnPushingData(buf[:bodyLen])
		}
	}
}

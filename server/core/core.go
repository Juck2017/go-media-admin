package core

import (
	"suitbim.com/go-media-admin/common"
)

// 定义核心包下会用到的共享变量
var (
	// G 全局配置对象引用
	G = *common.GlobalConf
)

// 定义核心包下会用到的常量
const (
	// SIG_OPEN 信号`开启推流`
	SIG_OPEN = "open"
	// SIG_CLOS 信号`关闭推流`
	SIG_CLOS = "close"
	// SIG_ACTS 信号`当前激活通道`
	SIG_ACTS = "actives"
)

// Message 用于解构从Dispatcher发来的消息
type Message struct {
	Sign    string   `json:"signal"`   // 信号,开启/关闭
	Chan    string   `json:"channel"`  // 通道编码
	Viewers []string `json:"Viewers"`  // 消费者
}

type IPusher interface {
	Pull() (data []byte, err error)
	Close() error
}

// JSON 仿JSON
type JSON map[string]interface{}

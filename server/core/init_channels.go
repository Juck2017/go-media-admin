package core

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

// LoadChannels 加载所有通道的配置信息
func LoadChannels() []Channel {
	db, err := sql.Open("sqlite3", G.App.Db.Path)
	if err != nil {
		G.Logger.Errorf("Open Database [%s] Failed: %v", G.App.Db.Path, err)
		panic(err)
	}

	defer db.Close()

	//查询所有通道数据
	rows, err := db.Query("SELECT id, code, status FROM channel")
	if err != nil {
		G.Logger.Errorf("Query Channels Failed: %v", err)
		panic(err)
	}
	// 定义通道数组存放多个通道数据
	var channels = make([]Channel, 0)

	for rows.Next() {
		var c Channel
		_ = rows.Scan(&c.Id, &c.Code, &c.Status)
		channels = append(channels, c)
	}
	G.Logger.Infof("Load all channel success")
	return channels
}

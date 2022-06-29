# go-media-pusher

#### 介绍
go-media-pusher 针对复杂的网络情况，就近部署于可连接摄像头或录像机的网络坏境，根据go-media-dispather的指令进行实时视频流采集、转码和推流。


#### 软件架构
软件架构说明


#### 安装教程
```bash
go mod tidy
go build
build.sh 压缩打包
tools/upx 压缩打包用到的工具
```

#### 使用说明

1. banner.txt     配置启动日志打印文本banner
2. db/data.db     配置本pusher所能完成的视频流通道列表
3. config.yml     主配置文件：配置应用运行参数
4. 上述3个配置文件须和编译后的可执行文件同目录放置（RTSP协议）
5. banner配置网址https://www.bootschool.net/ascii,字体为bell

#### 参与贡献

1.  Fork 本仓库
2.  新建 Feat_xxx 分支
3.  提交代码
4.  新建 Pull Request


#### 其他说明
1. 前端视频解码库 viewer/jsmpeg.min.js中，new WebSocket时，原版会多传递第二个protocol参数，值未null，这会使得服务端不能正常处理（可能是个小bug），将第二个参数去掉后就可以正常工作了

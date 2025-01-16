package router

import (
	"Fire/service"
	"log"
	"net"
)

func StartVideoStreamServer() {
	listener, err := net.Listen("tcp", "192.168.31.232:5000")
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		videoStreamService := service.NewVideoStreamService(conn)
		go videoStreamService.SendVideoStream("/home/jhq/GolandProjects/video.mp4") // 指定视频文件路径
	}
}

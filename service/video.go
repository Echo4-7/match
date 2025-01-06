package service

import (
	"io"
	"log"
	"net"
	"os"
)

type VideoStreamService struct {
	conn net.Conn
}

func NewVideoStreamService(conn net.Conn) *VideoStreamService {
	return &VideoStreamService{conn: conn}
}

func (v *VideoStreamService) SendVideoStream(videoPath string) {
	defer v.conn.Close()

	videoFile, err := os.Open(videoPath)
	if err != nil {
		log.Println("无法打开视频文件:", err)
		return
	}
	defer videoFile.Close()

	buffer := make([]byte, 1024*1024*10) // 1MB的缓冲区
	for {
		n, err := videoFile.Read(buffer)
		if err != nil {
			if err != io.EOF {
				log.Println("读取视频文件时出错:", err)
			}
			break
		}
		_, err = v.conn.Write(buffer[:n])
		if err != nil {
			log.Println("发送视频数据时出错:", err)
			break
		}
	}
}

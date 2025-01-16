package service

import (
	"bufio"
	"fmt"
	"io"
	l "log"
	"net"
	"os/exec"
)

type VideoStreamService struct {
	conn net.Conn
}

func NewVideoStreamService(conn net.Conn) *VideoStreamService {
	return &VideoStreamService{conn: conn}
}

func (v *VideoStreamService) SendVideoStream(videoPath string) {
	// 每秒拆分的帧数
	fps := 1

	// 构建 FFmpeg 命令
	cmd := exec.Command("ffmpeg", "-i", videoPath, "-vf", fmt.Sprintf("fps=%d", fps), "-f", "image2pipe", "-vcodec", "mjpeg", "-")

	// 创建管道
	pipeReader, pipeWriter := io.Pipe()
	cmd.Stdout = pipeWriter

	// 启动 FFmpeg 命令
	err := cmd.Start()
	if err != nil {
		l.Fatalf("Failed to start FFmpeg command: %v", err)
	}

	// 从管道读取数据并发送到服务器
	reader := bufio.NewReader(pipeReader)
	for {
		// 读取一帧数据
		frame, err := reader.ReadBytes('\n')

		if err != nil {
			if err == io.EOF {
				break // FFmpeg 完成输出
			}
			l.Fatalf("Failed to read frame: %v", err)
		}

		// 发送数据
		_, err = v.conn.Write(frame)
		fmt.Println(frame)
		if err != nil {
			l.Fatalf("Failed to send frame: %v", err)
		}
	}
	// 关闭管道写端
	pipeWriter.Close()

	// 等待 FFmpeg 命令完成
	err = cmd.Wait()
	if err != nil {
		l.Fatalf("FFmpeg command failed: %v", err)
	}
}

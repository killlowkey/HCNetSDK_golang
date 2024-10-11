package main

import (
	"alarm/internal/Device"
	"fmt"
	"log"
	"runtime"
	"time"
)

func init() {
	Device.InitHikSDK()
}

func main() {
	fmt.Println(runtime.GOOS)
	info := Device.DeviceInfo{
		IP:       "172.16.19.3",
		UserName: "admin",
		Password: "itc123456",
		Port:     8000,
	}
	device := Device.NewHKDevice(info)
	if userId, deviceInfo, err := device.Login(); err != nil {
		log.Fatal("登录失败", err.Error())
	} else {
		log.Printf("成功登录，用户id=%d, 设备信息=%v\n", userId, deviceInfo)
	}
	defer func() {
		_ = device.Logout()
		Device.HKExit()
	}()

	if id, err := device.Play(); err != nil {
		log.Fatal("播放失败", err.Error())
	} else {
		log.Printf("成功播放, id=%d\n", id)
	}

	filepath := "./images/" + time.Now().Format("20060102150405") + ".jpeg"
	if err := device.Capture(filepath); err != nil {
		log.Fatal("抓拍失败", err.Error())
	} else {
		log.Printf("成功抓拍, 图片地址=%s\n", filepath)
	}

	// device.SetAlarmCallBack()
	// device.StartListenAlarmMsg()
	// time.Sleep(time.Second * 100)
	// device.StopListenAlarmMsg()
}

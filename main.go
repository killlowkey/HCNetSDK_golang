package main

import (
	"alarm/internal/Device"
	"fmt"
	"log"
	"runtime"
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
	if userId, info, err := device.Login(); err != nil {
		log.Fatal("登录失败", err.Error())
	} else {
		log.Printf("成功登录，用户id=%d, 设备信息=%v\n", userId, info)
	}
	defer func() {
		device.Logout()
		Device.HKExit()
	}()

	if id, err := device.Play(); err != nil {
		log.Fatal("播放失败", err.Error())
	} else {
		log.Printf("成功播放, id=%d\n", id)
	}

	if path, err := device.Capture(); err != nil {
		log.Fatal("抓拍失败", err.Error())
	} else {
		log.Printf("成功抓拍, 图片地址=%s\n", path)
	}

	// device.SetAlarmCallBack()
	// device.StartListenAlarmMsg()
	// time.Sleep(time.Second * 100)
	// device.StopListenAlarmMsg()
}

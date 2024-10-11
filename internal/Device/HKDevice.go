package Device

/*
#cgo CFLAGS:  -I../../include
#cgo LDFLAGS: -L$../../build  -lHCCore -lhpr -lhcnetsdk
#include <stdio.h>
#include <stdlib.h>
#include "HCNetSDK.h"
extern void AlarmCallBack(LONG lCommand, NET_DVR_ALARMER *pAlarmer, char *pAlarmInfo, DWORD dwBufLen, void* pUser);
*/
import "C"
import (
	"errors"
	"fmt"
	"unsafe"
)

// #cgo LDFLAGS 需要根据不同的平台修改 hcnetsdk（linux) / HCNETSDK(Windows)

//export AlarmCallBack
func AlarmCallBack(command C.LONG, alarm *C.NET_DVR_ALARMER, info *C.char, len C.DWORD, user unsafe.Pointer) {
	fmt.Println("receive alarm")
}

type HKDevice struct {
	ip          string
	port        int
	username    string
	password    string
	loginId     int
	alarmHandle int
}

// InitHikSDK hk sdk init
func InitHikSDK() {
	C.NET_DVR_Init()
	C.NET_DVR_SetConnectTime(2000, 5)
	C.NET_DVR_SetReconnect(10000, 1)
}

// HKExit hk sdk clean
func HKExit() {
	C.NET_DVR_Cleanup()
}

// NewHKDevice new hk-device instance
func NewHKDevice(info DeviceInfo) Device {
	return &HKDevice{
		ip:       info.IP,
		port:     info.Port,
		username: info.UserName,
		password: info.Password,
	}
}

// Login hk device loin
func (device *HKDevice) Login() (int, *NET_DVR_DEVICEINFO_V30, error) {
	// init data
	var deviceInfoV30 C.NET_DVR_DEVICEINFO_V30
	ip := C.CString(device.ip)
	usr := C.CString(device.username)
	passwd := C.CString(device.password)
	defer func() {
		C.free(unsafe.Pointer(ip))
		C.free(unsafe.Pointer(usr))
		C.free(unsafe.Pointer(passwd))
	}()

	device.loginId = int(C.NET_DVR_Login_V30(ip, C.WORD(device.port), usr, passwd,
		(*C.NET_DVR_DEVICEINFO_V30)(unsafe.Pointer(&deviceInfoV30)),
	))
	if device.loginId < 0 {
		return -1, nil, device.HKErr("login")
	}

	var deviceInfo = new(NET_DVR_DEVICEINFO_V30)
	deviceInfo.Convert(deviceInfoV30)
	return device.loginId, deviceInfo, nil
}

// Logout hk device logout
func (device *HKDevice) Logout() error {
	C.NET_DVR_Logout_V30(C.LONG(device.loginId))
	if err := device.HKErr("NVRLogout"); err != nil {
		return err
	}
	return nil
}

func (device *HKDevice) SetAlarmCallBack() error {
	if C.NET_DVR_SetDVRMessageCallBack_V30(C.MSGCallBack(C.AlarmCallBack), C.NULL) != C.TRUE {
		return device.HKErr(device.ip + ":set alarm callback")
	}
	return nil
}
func (device *HKDevice) StartListenAlarmMsg() error {
	var struAlarmParam C.NET_DVR_SETUPALARM_PARAM
	struAlarmParam.dwSize = C.ulong(unsafe.Sizeof(struAlarmParam)) //Windows -> C.ulong
	struAlarmParam.byAlarmInfoType = 0

	device.alarmHandle = int(C.NET_DVR_SetupAlarmChan_V41(C.long(device.loginId), &struAlarmParam)) // Windows -> C.long
	if device.alarmHandle < 0 {
		return device.HKErr("setup alarm chan")
	}
	return nil
}

func (device *HKDevice) StopListenAlarmMsg() error {
	if C.NET_DVR_CloseAlarmChan_V30(C.long(device.alarmHandle)) != C.TRUE { //Windows  C.long
		return device.HKErr("stoop alarm chan")
	}
	return nil
}

// Play 播放视频
// uid:摄像头登录成功的id
// 返回播放视频标识 pid
func (device *HKDevice) Play() (int64, error) {
	var pDetectInfo C.NET_DVR_CLIENTINFO
	pDetectInfo.lChannel = C.LONG(1)
	pid := C.NET_DVR_RealPlay_V30(C.LONG(device.loginId), (*C.NET_DVR_CLIENTINFO)(unsafe.Pointer(&pDetectInfo)), nil, nil, C.BOOL(1))
	if int64(pid) < 0 {
		if err := isErr("Play"); err != nil {
			return -1, err
		}
		return -1, errors.New("播放失败")
	}

	return int64(pid), nil
}

// Capture 抓拍
func (device *HKDevice) Capture(filepath string) error {
	var jpegpara C.NET_DVR_JPEGPARA
	var lChannel uint32 = 1
	c_path := C.CString(filepath)
	defer C.free(unsafe.Pointer(c_path))
	msgId := C.NET_DVR_CaptureJPEGPicture(C.LONG(device.loginId), C.LONG(lChannel),
		(*C.NET_DVR_JPEGPARA)(unsafe.Pointer(&jpegpara)),
		c_path,
	)

	if int64(msgId) < 0 {
		if err := isErr("Capture"); err != nil {
			return err
		}
		return errors.New("抓拍失败")
	}
	return nil
}

// 是否有错误
func isErr(oper string) error {
	errno := int64(C.NET_DVR_GetLastError())
	if errno > 0 {
		reMsg := fmt.Sprintf("%s摄像头失败,失败代码号：%d", oper, errno)
		return errors.New(reMsg)
	}
	return nil
}

// HKErr Detect success of operation
func (device *HKDevice) HKErr(operation string) error {
	errno := int64(C.NET_DVR_GetLastError())
	if errno > 0 {
		reMsg := fmt.Sprintf("%s:%s摄像头失败,失败代码号：%d", device.ip, operation, errno)
		return errors.New(reMsg)
	}
	return nil
}

// CToGo 将C结构体转换为Go结构体
func (n *NET_DVR_DEVICEINFO_V30) Convert(cStruct C.NET_DVR_DEVICEINFO_V30) {
	// 复制设备地址
	serialNumber := C.GoBytes(unsafe.Pointer(&cStruct.sSerialNumber), C.int(len(cStruct.sSerialNumber)))
	n.SSerialNumber = string(serialNumber)
	n.ByAlarmInPortNum = *(*byte)(unsafe.Pointer(&cStruct.byAlarmInPortNum))
	n.ByAlarmOutPortNum = *(*byte)(unsafe.Pointer(&cStruct.byAlarmOutPortNum))
	n.ByDiskNum = *(*byte)(unsafe.Pointer(&cStruct.byDiskNum))
	n.ByDVRType = *(*byte)(unsafe.Pointer(&cStruct.byDVRType))
	n.ByChanNum = *(*byte)(unsafe.Pointer(&cStruct.byChanNum))
	n.ByStartChan = *(*byte)(unsafe.Pointer(&cStruct.byStartChan))
	n.ByAudioChanNum = *(*byte)(unsafe.Pointer(&cStruct.byAudioChanNum))
	n.ByIPChanNum = *(*byte)(unsafe.Pointer(&cStruct.byIPChanNum))
	n.ByZeroChanNum = *(*byte)(unsafe.Pointer(&cStruct.byZeroChanNum))
	n.ByMainProto = *(*byte)(unsafe.Pointer(&cStruct.byMainProto))
	n.BySubProto = *(*byte)(unsafe.Pointer(&cStruct.bySubProto))
	n.BySupport = *(*byte)(unsafe.Pointer(&cStruct.bySupport))
	n.BySupport1 = *(*byte)(unsafe.Pointer(&cStruct.bySupport1))
	n.BySupport2 = *(*byte)(unsafe.Pointer(&cStruct.bySupport2))
	n.WDevType = *(*uint16)(unsafe.Pointer(&cStruct.wDevType))
	n.BySupport3 = *(*byte)(unsafe.Pointer(&cStruct.bySupport3))
	n.ByMultiStreamProto = *(*byte)(unsafe.Pointer(&cStruct.byMultiStreamProto))
	n.ByStartDChan = *(*byte)(unsafe.Pointer(&cStruct.byStartDChan))
	n.ByStartDTalkChan = *(*byte)(unsafe.Pointer(&cStruct.byStartDTalkChan))
	n.ByHighDChanNum = *(*byte)(unsafe.Pointer(&cStruct.byHighDChanNum))
	n.BySupport4 = *(*byte)(unsafe.Pointer(&cStruct.bySupport4))
	n.ByLanguageType = *(*byte)(unsafe.Pointer(&cStruct.byLanguageType))
	n.ByVoiceInChanNum = *(*byte)(unsafe.Pointer(&cStruct.byVoiceInChanNum))
	n.ByStartVoiceInChanNo = *(*byte)(unsafe.Pointer(&cStruct.byStartVoiceInChanNo))
	n.BySupport5 = *(*byte)(unsafe.Pointer(&cStruct.bySupport5))
	n.BySupport6 = *(*byte)(unsafe.Pointer(&cStruct.bySupport6))
	n.ByMirrorChanNum = *(*byte)(unsafe.Pointer(&cStruct.byMirrorChanNum))
	n.WStartMirrorChanNo = *(*uint16)(unsafe.Pointer(&cStruct.wStartMirrorChanNo))
	n.BySupport7 = *(*byte)(unsafe.Pointer(&cStruct.bySupport7))
	n.ByRes2 = *(*byte)(unsafe.Pointer(&cStruct.byRes2))
}

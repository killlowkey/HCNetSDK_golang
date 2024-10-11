package Device

type Device interface {
	// Login 登录
	Login() (int, *NET_DVR_DEVICEINFO_V30, error)
	// Logout 注销
	Logout() error
	SetAlarmCallBack() error
	StartListenAlarmMsg() error
	StopListenAlarmMsg() error
	// Play 播放视频
	Play() (int64, error)
	// Capture 抓图
	Capture() (string, error)
}
type DeviceInfo struct {
	IP       string
	Port     int
	UserName string
	Password string
}

type NET_DVR_DEVICEINFO_V30 struct {
	SSerialNumber        string // 序列号
	ByAlarmInPortNum     byte   // 报警输入个数
	ByAlarmOutPortNum    byte   // 报警输出个数
	ByDiskNum            byte   // 硬盘个数
	ByDVRType            byte   // 设备类型, 1:DVR 2:ATM DVR 3:DVS ......
	ByChanNum            byte   // 模拟通道个数
	ByStartChan          byte   // 起始通道号,例如DVS-1,DVR - 1
	ByAudioChanNum       byte   // 语音通道数
	ByIPChanNum          byte   // 最大数字通道个数，低位
	ByZeroChanNum        byte   // 零通道编码个数
	ByMainProto          byte   // 主码流传输协议类型 0-private, 1-rtsp, 2-同时支持private和rtsp
	BySubProto           byte   // 子码流传输协议类型 0-private, 1-rtsp, 2-同时支持private和rtsp
	BySupport            byte   // 能力，位与结果为0表示不支持，1表示支持
	BySupport1           byte   // 能力集扩展，位与结果为0表示不支持，1表示支持
	BySupport2           byte   // 能力集扩展，位与结果为0表示不支持，1表示支持
	WDevType             uint16 // 设备型号
	BySupport3           byte   // 能力集扩展，位与结果为0表示不支持，1表示支持
	ByMultiStreamProto   byte   // 是否支持多码流，按位表示，0-不支持，1-支持，bit1-码流3，bit2-码流4，bit7-主码流，bit-8子码流
	ByStartDChan         byte   // 起始数字通道号，0表示无效
	ByStartDTalkChan     byte   // 起始数字对讲通道号，区别于模拟对讲通道号，0表示无效
	ByHighDChanNum       byte   // 数字通道个数，高位
	BySupport4           byte   // 能力集扩展，位与结果为0表示不支持，1表示支持
	ByLanguageType       byte   // 语言类型
	ByVoiceInChanNum     byte   // 音频输入通道数
	ByStartVoiceInChanNo byte   // 音频输入起始通道号 0表示无效
	BySupport5           byte   // 能力集扩展，位与结果为0表示不支持，1表示支持
	BySupport6           byte   // 能力集扩展，位与结果为0表示不支持，1表示支持
	ByMirrorChanNum      byte   // 镜像通道个数，<录播主机中用于表示导播通道>
	WStartMirrorChanNo   uint16 // 起始镜像通道号
	BySupport7           byte   // 能力集扩展，位与结果为0表示不支持，1表示支持
	ByRes2               byte   // 保留
}

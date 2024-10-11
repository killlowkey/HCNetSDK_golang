package Device

type Device interface {
	Login() (int, error)
	Logout() error
	SetAlarmCallBack() error
	StartListenAlarmMsg() error
	StopListenAlarmMsg() error
	Play() (int64, error)
	Capture() (string, error)
}
type DeviceInfo struct {
	IP       string
	Port     int
	UserName string
	Password string
}

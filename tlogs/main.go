package tlogs

func init() {
	logDisable()
	logConfigLoad()
}

// LogType : 日志类型
type LogType int

const (
	Info LogType = iota
	Warn
	Error
)

// DoLog :执行日志
func DoLog(lt LogType, v ...interface{}) {
	switch lt {
	case Info:
		LogMain.Info(v...)
	case Warn:
		LogMain.Warn(v...)
	case Error:
		LogMain.Error(v...)
	}
}

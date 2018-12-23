package tlogs

func init() {
	// 初始化日志
	logDisable()
	logConfigLoad()
}

// LogType : 日志类型
type LogType int

const (
	Debug LogType = iota
	Info
	Warn
	Error
)

// DoLog :执行日志
func DoLog(lt LogType, v ...interface{}) {
	switch lt {
	case Debug:
		LogMain.Debug(v...)
	case Info:
		LogMain.Info(v...)
	case Warn:
		LogMain.Warn(v...)
	case Error:
		LogMain.Error(v...)
	}
}

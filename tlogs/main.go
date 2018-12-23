package tlogs

func init() {
	// 初始化日志
	LogDisable()
	LogConfigLoad()
}

func LogDebug(v ...interface{}) {
	LogMain.Debug(v...)
}

func LogInfo(v ...interface{}) {
	LogMain.Info(v...)
}

func LogError(v ...interface{}) {
	LogMain.Error(v...)
}

func LogWarn(v ...interface{}) {
	LogMain.Warn(v...)
}

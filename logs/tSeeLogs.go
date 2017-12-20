package logs

import (
	"log"

	seelog "github.com/cihub/seelog"
)

// Global var base on seelog framework
var LogMain seelog.LoggerInterface

// local log configration
var strLogConfig string = `
<seelog >
    <outputs formatid="all">
        <file path="./log/tcplog.log"/>
        <filter levels="info">
          <console formatid="fmtinfo"/>
        </filter>
        <filter levels="error" formatid="fmterror">
          <console/>
          <file path="./log/errors.log"/>
        </filter>
		<filter levels="critical">
			<file path="./log/critical.log" formatid="critical"/>
			<smtp formatid="criticalemail" senderaddress="12679333@qq.com" sendername="logTest"  hostname="smtp.qq.com" hostport="587"
				username="12679333@qq.com" password="dn65523213">
				<recipient address="david.hu@makedoctor.com" />
			</smtp>
		</filter>
    </outputs>
    <formats>
        <format id="fmtinfo" format="[%Level] [%Time] %Msg%n"/>
        <format id="fmterror" format="[%LEVEL] [%Time] [%FuncShort @ %File.%Line] %Msg%n"/>
		<format id="critical" format="[%LEVEL] [%Time] [%FuncShort @ %File.%Line] %Msg%n"/>
        <format id="all" format="[%Level] [%Time] %Msg%n"/>
        <format id="criticalemail" format="Critical error on our server!\n    %Time %Date %RelFile %Func %Msg \nSent by Seelog"/>
    </formats>
</seelog>
`

// load the config file
func LogConfigLoad() {
	//	logger, err := seelog.LoggerFromConfigAsFile("./config/seelog-main.xml")
	logger, err := seelog.LoggerFromConfigAsBytes([]byte(strLogConfig))
	if err != nil {
		log.Println("初始化日志出错")
		log.Println(err)
		return
	}
	UserLogger(logger)
}

// 使用指定日志对象
func UserLogger(newLogger seelog.LoggerInterface) {
	LogMain = newLogger
}

// 初始化全局变量Logger为seeLog的禁用状态，主要为防止Logger被多次初始化
func LogDisable() {
	LogMain = seelog.Disabled
}

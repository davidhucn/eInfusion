package httpOperate

import (
    "eInfusion/comm"
    "golang.org/x/net/websocket"
    "net/http"
    "os"
    "time"
)

//错误处理函数
func checkErr(err error, extra string) bool {
    if err != nil {
        formatStr := " Err : %s\n";
        if extra != "" {
            formatStr = extra + formatStr;
        }

        fmt.Fprintf(os.Stderr, formatStr, err.Error());
        return true;
    }

    return false;
}

func svrConnHandler(conn *websocket.Conn) {
    request := make([]byte, 128);
    defer conn.Close();
    for {
        readLen, err := conn.Read(request)
        if checkErr(err, "Read") {
            break;
        }

        //socket被关闭了
        if readLen == 0 {
            fmt.Println("Client connection close!");
            break;
        } else {
            //输出接收到的信息
            fmt.Println(string(request[:readLen]))

            time.Sleep(time.Second);
            //发送
            conn.Write([]byte("World !"));
        }

        request = make([]byte, 128);
    }
}

func main() {
    http.Handle("/echo", websocket.Handler(svrConnHandler));
    err := http.ListenAndServe(":6666", nil);
    checkErr(err, "ListenAndServe");
    fmt.Println("Func finish.");
}
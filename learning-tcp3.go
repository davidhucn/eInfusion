 /** 三步搭建服务端
        1 定义任意名称struct的数据结构，必须包含Pmap、Phost两个
          字段，其中Phost为服务端ip+port拼接的字符串，Pmap为自定
          义数据包类型与数据包名称的映射。
        2 实例化对象为字段赋值，实现对应已定义`包名称`的数据包处
          理方法，方法名必为"P[包名称]",如type包的处理方法为Ptype
          。方法中请定义数据处理逻辑,输入输入皆为[]byte类型。
        3 stpro.New()传入实例化的对象，如无报错则服务端开始监听，
          并按照你所定义的逻辑处理数据包，返回响应数据。
    **/
    package main

    import (
        "fmt"
        "stpro"
    )

    type Server struct {
        Phost string
        Pmap  map[uint8]string
    }

    func (m Server) Ptype(in []byte) (out []byte) {
        fmt.Printf("客户端发来type包:%s\n", in)
        /** process... **/
        bytes := []byte("hello1")
        return bytes
    }

    func (m Server) Pname(in []byte) (out []byte) {
        fmt.Printf("客户端发来name包:%s\n", in)
        /** process... **/
        bytes := []byte("hello2")
        return bytes
    }

    func main() {
        m := Model{
            Phost: ":9091",
            Pmap:  make(map[uint8]string),
        }
        m.Pmap[0x01] = "type"
        m.Pmap[0x02] = "name"
        err := stpro.New(m)
        if err != nil {
            fmt.Println(err)
        }
   }
3.client端
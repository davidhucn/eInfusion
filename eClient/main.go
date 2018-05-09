package main

import (
	"encoding/json"
	"fmt"
	//	"net/url"
	//	"io/ioutil"
	"bytes"
	"net/http"
)

func main() {
	//	resp, _ := http.Get("https://www.cnblogs.com/hitfire/articles/6427033.html")
	//	defer resp.Body.Close()
	//	body, _ := ioutil.ReadAll(resp.Body)
	//	fmt.Println(string(body))

	song := make(map[string]interface{})
	song["name"] = "李白"
	song["timelength"] = 128
	song["author"] = "李荣浩"

	bytesData, err := json.Marshal(song)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	reader := bytes.NewReader(bytesData)
	url := "http://127.0.0.1/r:7779"
	request, err := http.NewRequest("POST", url, reader)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	client := http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Println("error")
		return
	}

	var dd []byte
	n, _ := resp.Body.Read(dd)
	fmt.Print("n:", n)
	fmt.Print("dd:", dd)
	//	respBytes, err := ioutil.ReadAll(resp.Body)
	//	if err != nil {
	//		fmt.Println(err.Error())
	//		return
	//	}
	//	//byte数组直接转成string，优化内存
	//	str := (*string)(unsafe.Pointer(&respBytes))
	//	fmt.Println(*str)

}

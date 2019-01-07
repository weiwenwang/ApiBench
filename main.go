package main

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"strings"
	"net/url"
	"log"
	"time"
	"sync/atomic"
	"sync"
)

func main() {
	var sum int64
	a := sync.WaitGroup{}
	for j := 0; j < 200; j++ {
		a.Add(1)
		go func() {
			for i := 0; i < 10; i++ {
				t1 := time.Now().UnixNano()
				httpHandle("GET",
					"https://yxpt-alpha.wochacha.cn/public/?s=App.EnterFake.WarningIndex&limit=30&current=1", "")
				t2 := time.Now().UnixNano()
				atomic.AddInt64(&sum, (t2-t1)/1e6)
				fmt.Println((t2 - t1) / 1e6)
			}
			a.Done()
		}()
	}
	a.Wait()
	fmt.Println("sum:", sum/200/10)
}

//http请求
func httpHandle(method, urlVal, data string) {
	client := &http.Client{}
	var req *http.Request

	if data == "" {
		urlArr := strings.Split(urlVal, "?")
		if len(urlArr) == 2 {
			urlVal = urlArr[0] + "?" + getParseParam(urlArr[1])
		}
		req, _ = http.NewRequest(method, urlVal, nil)
	} else {
		req, _ = http.NewRequest(method, urlVal, strings.NewReader(data))
	}

	//添加cookie，key为X-Xsrftoken，value为df41ba54db5011e89861002324e63af81
	//可以添加多个cookie
	cookie1 := &http.Cookie{Name: "token",
		Value: "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpYXQiOjE1NDA4MDQ2MTcsImV4cCI6MTU0MzM5NjYxNywiaWQiOjEsInVzZXJuYW1lIjoiYWRtaW4xMjMiLCJpc19zdXBlcnVzZXIiOjEsInl4cHRfcm9sZSI6MSwibGFzdF9sb2dpbiI6IjIwMTgtMTAtMjkgMTQ6NTI6NTIiLCJsb2dpblR5cGUiOiJwYXNzd29yZCJ9.kJjTa4gGRbDSs1RnQFwLppwBuFB_fMQiWI0__-airp4",
		HttpOnly: true}
	req.AddCookie(cookie1)

	//添加header，key为X-Xsrftoken，value为b6d695bbdcd111e8b681002324e63af81
	req.Header.Add("token", "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpYXQiOjE1NDA4MDQ2MTcsImV4cCI6MTU0MzM5NjYxNywiaWQiOjEsInVzZXJuYW1lIjoiYWRtaW4xMjMiLCJpc19zdXBlcnVzZXIiOjEsInl4cHRfcm9sZSI6MSwibGFzdF9sb2dpbiI6IjIwMTgtMTAtMjkgMTQ6NTI6NTIiLCJsb2dpblR5cGUiOiJwYXNzd29yZCJ9.kJjTa4gGRbDSs1RnQFwLppwBuFB_fMQiWI0__-airp4")

	resp, err := client.Do(req)

	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	ioutil.ReadAll(resp.Body)
	//b, _ := ioutil.ReadAll(resp.Body)

	//fmt.Println(string(b))
}

//将get请求的参数进行转义
func getParseParam(param string) string {
	return url.PathEscape(param)
}

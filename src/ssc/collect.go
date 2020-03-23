package ssc

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"server/client"
	"server/verify"
	"strings"
	"syscall"
	"time"
)

func Collect(path string) {
	file, err := os.Open(path)
	if err != nil {
		err = errors.New("打开文件失败:" + err.Error())
		return
	}

	defer func() {
		file.Close()
	}()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		err = errors.New("读取文件失败:" + err.Error())

		return
	}

	dataStr := string(data)
	dataSplit := strings.Split(dataStr, "\r\n")
	//headers := make(map[string]string, 0)
	//headers["Content-Type"] = "application/json; charset=utf-8"
	//阻塞主线程
	chExit := make(chan os.Signal)
	signal.Notify(chExit, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGKILL)

	for _, url := range dataSplit {
		go func(url string) {
			t := time.NewTimer(10 * time.Second)
			for {
				select {
				case <-t.C:
					response, _ := client.HttpRequest("GET", url, nil, nil, 0)
					code := response.StatusCode
					if code == http.StatusOK {
						respBody, _ := verify.GetZipData(response)
						//respBody, _ := ioutil.ReadAll(response.Body)
						//var smsRet sms
						//json.Unmarshal(respBody, &smsRet)
						fmt.Printf("%s 结果:%s \n", url, string(respBody))
					}
				}
			}

		}(url)
	}
	s := <-chExit
	fmt.Print(s)
}

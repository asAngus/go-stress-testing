package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"server/client"
	"server/verify"
	"strings"
	"syscall"
	"time"
)

func main() {
	var (
		pix string
		url string
		num uint64
	)
	mChan := make(chan string, 10)
	flag.StringVar(&pix, "p", "", "号码前缀")
	flag.StringVar(&url, "url", "http://dx.uat.huanqiujr.com", "地址")
	flag.Parse()
	//file, err := os.Open(path)
	//if err != nil {
	//	err = errors.New("打开文件失败:" + err.Error())
	//	return
	//}

	//defer func() {
	//	file.Close()
	//}()
	//
	//data, err := ioutil.ReadAll(file)
	//if err != nil {
	//	err = errors.New("读取文件失败:" + err.Error())
	//
	//	return
	//}

	//dataStr := string(data)

	//arrmbls := strings.Split(dataStr, ",")
	headers := make(map[string]string, 0)
	headers["Content-Type"] = "application/json; charset=utf-8"
	chExit := make(chan os.Signal)
	//d := fmt.Sprintf("%06v", 10)
	fmt.Println(1 << 8)
	//pix := "1381000"
	pLen := len(pix)
	sLen := 11 - pLen
	switch sLen {
	case 1:
		num = 10
	case 2:
		num = 100
	case 3:
		num = 1000
	case 4:
		num = 10000
	case 5:
		num = 100000
	}
	s1 := pix + "%" + fmt.Sprintf(`0%dv`, sLen)
	signal.Notify(chExit, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGKILL)
	go func() {
		fmt.Println("[receiver] Received a sync signal and wait a second...")
		time.Sleep(time.Second)
		for {
			if elem, ok := <-mChan; ok {
				appendFile(pix+".txt", elem)
				fmt.Println("[receiver] Received:", elem)
			} else {
				break
			}
		}
		fmt.Println("[receiver] Stopped.")
	}()

	var i uint64
	for {
		mbl := fmt.Sprintf(s1, i)
		fmt.Println(mbl)
		go func(mbl string) {
			//-X POST -H "Content-Type:application/json" -d '"title":"comewords","content":"articleContent"'
			//发送验证码

			body := strings.NewReader(fmt.Sprintf(`{"mblNo": "%s"}`, mbl))

			response, _ := client.HttpRequest("POST", url+"/v1/edu/register/student/captchaSms", body, headers, 0)
			code := response.StatusCode
			if code == http.StatusOK {
				respBody, _ := verify.GetZipData(response)
				//respBody, _ := ioutil.ReadAll(response.Body)
				//var smsRet sms
				//json.Unmarshal(respBody, &smsRet)
				fmt.Printf("短信手机号码:%s,结果:%s,  \n", mbl, string(respBody))
			}
			//注册
			body = strings.NewReader(fmt.Sprintf(`{"mblNo": "%s","captcha": "2035","nick": "nike_%s","password": "123456789"}`, mbl, mbl))

			response, _ = client.HttpRequest("POST", url+"/v1/edu/register/student", body, headers, 0)
			code = response.StatusCode
			if code == http.StatusOK {
				respBody, _ := verify.GetZipData(response)
				//respBody, _ := ioutil.ReadAll(response.Body)
				//var smsRet sms
				//json.Unmarshal(respBody, &smsRet)
				fmt.Printf("注册手机号码:%s,结果:%s \n", mbl, string(respBody))
				mChan <- mbl + ","
			}
		}(mbl)
		time.Sleep(5 * time.Millisecond)
		i = i + 1
		if i > num {
			break
		}
	}
	//主进程阻塞直到有退出信号
	s := <-chExit
	fmt.Print(s)
	//time.Sleep(1000000000000)
}
func appendFile(path string, content string) {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()
	_, err = f.WriteString(content)
	if err != nil {
		log.Println(err)
	}
}

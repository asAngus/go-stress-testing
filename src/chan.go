//测试chan
package main

import "fmt"
import "time"

var strChan = make(chan string, 3)

func main() {
	syncChan1 := make(chan struct{}, 1) //接收同步变量
	syncChan2 := make(chan struct{}, 2) //主线程启动了两个goruntime线程，
	//等这两个goruntime线程结束后主线程才能结束

	//用于演示接受操作
	go func() {
		<-syncChan1 //表示可以开始接收数据了，否则等待
		fmt.Println("[receiver] Received a sync signal and wait a second...")
		time.Sleep(time.Second)
		for {
			if elem, ok := <-strChan; ok {
				fmt.Println("[receiver] Received:", elem)
			} else {
				break
			}
		}
		fmt.Println("[receiver] Stopped.")
		syncChan2 <- struct{}{}
	}()

	//用于演示发送操作
	go func() {
		for i, elem := range []string{"a", "b", "c", "d"} {
			fmt.Println("[sender] Sent:", elem)
			strChan <- elem
			if (i+1)%3 == 0 {
				syncChan1 <- struct{}{}
				fmt.Println("[sender] Sent a sync signal. wait 1 secnd...")
				time.Sleep(time.Second)
			}
		}
		fmt.Println("[sender] wait 2 seconds...")
		time.Sleep(time.Second)
		close(strChan)
		syncChan2 <- struct{}{}
	}()

	//主线程等待发送线程和接收线程结束后再结束
	fmt.Println("[main] waiting...")
	<-syncChan2
	<-syncChan2
	fmt.Println("[main] stoped")
}

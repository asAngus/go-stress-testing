package main

import (
	"fmt"
	"time"
)

var mapChan = make(chan map[string]int, 1)

func main() {
	syncChan := make(chan struct{}, 2)

	//用于演示接收操作
	go func() {
		for {
			if elem, ok := <-mapChan; ok {
				elem["count"]++ //每次接收到chan里面的map对象后，将key为count的值加1
			} else {
				break
			}
		}
		fmt.Println("[receiveder] stoped.")
		syncChan <- struct{}{}
	}()

	//用于演示发送操作
	go func() {
		countMap := make(map[string]int)
		for i := 0; i < 5; i++ {
			mapChan <- countMap
			time.Sleep(time.Second)
			fmt.Println("[sender] the count map:", countMap)
		}
		fmt.Println("[sender] stop chan.")
		close(mapChan)
		fmt.Println("[sender] stoped.")
		syncChan <- struct{}{}
	}()
	fmt.Println("[main] waiting...")
	<-syncChan
	<-syncChan
	fmt.Println("[main] stoped.")
}

/**
* Created by GoLand.
* User: link1st
* Date: 2019-08-21
* Time: 15:43
 */

package golink

import (
	"fmt"
	"heper"
	"model"
	"server/client"
	"sync"
	"time"
)

const (
	firstTime    = 1 * time.Second // 连接以后首次请求数据的时间
	intervalTime = 1 * time.Second // 发送数据的时间间隔
)

var (
	// 请求完成以后是否保持连接
	keepAlive bool
)

func init() {
	keepAlive = true
}

// web socket go link
func WebSocket(studentId string, courseId string, chanId uint64, ch chan<- *model.RequestResults, totalNumber uint64, wg *sync.WaitGroup, request *model.Request, ws *client.WebSocket) {

	defer func() {
		wg.Done()
	}()

	// fmt.Printf("启动协程 编号:%05d \n", chanId)

	defer func() {
		ws.Close()
	}()

	var (
		i uint64
	)

	// 暂停60秒
	t := time.NewTimer(firstTime)
	for {
		select {
		case <-t.C:
			t.Reset(intervalTime)

			// 请求
			webSocketRequest(studentId, courseId, chanId, ch, i, request, ws)

			// 结束条件
			i = i + 1
			if i >= totalNumber {
				goto end
				//t.Reset(intervalTime)
			}
		}
	}

end:
	t.Stop()

	if keepAlive {
		// 保持连接
		chWaitFor := make(chan int, 0)
		<-chWaitFor
	}

	return
}

// 请求
func webSocketRequest(studentId string, courseId string, chanId uint64, ch chan<- *model.RequestResults, i uint64, request *model.Request, ws *client.WebSocket) {

	var (
		startTime = time.Now()
		isSucceed = false
		errCode   = model.HttpOk
	)

	// 需要发送的数据
	//seq := fmt.Sprintf("%d_%d", chanId, i)
	studentInJson := fmt.Sprintf(`{"method":"StudentIn","params":{"studentId":"%s","courseId":"%s","secret":"%s"}}`, studentId, courseId, `testin#123452!`)
	//openRoomJson := fmt.Sprintf(`{"method":"methodOpenClassRoom80293jrljs908e024oqiod-0-=123-129-3./[;[ iuao90","params":{"courseId":"11719"}}`);
	//{"method":\"StudentIn\",\"params\":{\"studentId\":\"" + studentId + "\",\"courseId\":\"" + courseId + "\"" + ",\"secret\":\"" + secret + "\"}}
	//err := ws.Write([]byte(`{"seq":"` + seq + `","cmd":"ping","data":{}}`))
	//if()
	err := ws.Write([]byte(studentInJson))

	if err != nil {
		errCode = model.RequestErr // 请求错误
	} else {
		var stopChan = make(chan bool)
		ping(studentId, stopChan, request, ws)
		var add bool
		var errCount int64
		for {
			// time.Sleep(1 * time.Second)
			msg, err := ws.Read()
			if err != nil {
				errCode = model.ParseError
				fmt.Println("读取数据 失败~")
				errCount++
				if errCode > 3 {
					break
				}
			} else {
				// fmt.Println(msg)
				//_, _ = request.VerifyWebSocket(request, seq, msg)
				if request.GetDebug() {
					fmt.Println("response:", string(msg))
				}
				errCode = model.HttpOk
				isSucceed = true
			}
			if !add {
				requestTime := uint64(heper.DiffNano(startTime))

				requestResults := &model.RequestResults{
					Time:      requestTime,
					IsSucceed: isSucceed,
					ErrCode:   errCode,
				}
				requestResults.SetId(chanId, i)
				ch <- requestResults
				add = true
			}

		}
		stopChan <- true
	}

}

//发送心跳消息
func ping(studentId string, stopChan chan bool, request *model.Request, ws *client.WebSocket) {
	ticker := time.NewTicker(60 * time.Second)

	//定时发送心跳
	go func(ws *client.WebSocket) {
		pingJson := `{"method":"Ping"}`
		for {
			select {
			case <-ticker.C:
				//发送心跳
				err := ws.Write([]byte(pingJson))
				if request.GetDebug() {
					if err != nil {
						fmt.Println("心跳发送失败")
					} else {
						fmt.Println("发送心跳成功" + studentId)
					}
				}
			case <-stopChan:
				// 处理完成
				fmt.Println("发送心跳结束" + studentId)
				return
			}
		}
	}(ws)
}

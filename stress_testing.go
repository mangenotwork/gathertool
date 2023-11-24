/*
*	Description : 接口压力测试， 并输出结果 TODO 测试
*	Author 		: ManGe
*	Mail 		: 2912882908@qq.com
**/

package gathertool

import (
	"log"
	"net"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

type stateCodeData struct {
	Code    int
	ReqTime int64
}

// StressUrl 压力测试一个url
type StressUrl struct {
	Url    string
	Method string
	Sum    int64
	Total  int
	TQueue TodoQueue

	// 请求时间累加
	sumReqTime int64

	// 测试结果
	avgReqTime time.Duration

	// 接口传入的json
	JsonData string

	// 接口传入类型
	ContentType string

	// 是否重试
	isRetry bool

	stateCodeList    []*stateCodeData
	stateCodeListMux *sync.Mutex
}

// NewTestUrl 实例化一个新的url压测
func NewTestUrl(url, method string, sum int64, total int) *StressUrl {
	return &StressUrl{
		Url:              url,
		Method:           method,
		Sum:              sum,
		Total:            total,
		TQueue:           NewQueue(),
		sumReqTime:       int64(0),
		isRetry:          false,
		stateCodeList:    make([]*stateCodeData, 0),
		stateCodeListMux: &sync.Mutex{},
	}
}

// SetJson 设置json
func (s *StressUrl) SetJson(str string) {
	s.JsonData = str
}

func (s *StressUrl) OpenRetry() {
	s.isRetry = true
}

// Run 运行压测
func (s *StressUrl) Run(vs ...any) {
	//解析可变参
	var (
		succeedFunc  SucceedFunc
		n            int64
		wg           sync.WaitGroup
		reqTimeout   ReqTimeOut
		reqTimeoutMs ReqTimeOutMs
		header       http.Header
	)
	for _, v := range vs {
		switch vv := v.(type) {
		// 使用方传入了 header
		case SucceedFunc:
			succeedFunc = vv
		case ReqTimeOut:
			reqTimeout = vv
		case ReqTimeOutMs:
			reqTimeoutMs = vv
		case http.Header:
			header = vv
		case *http.Header:
			header = *vv
		}
	}

	//初始化队列
	for n = 0; n < s.Sum; n++ {
		_ = s.TQueue.Add(&Task{Url: s.Url})
	}
	log.Println("总执行次数： ", s.TQueue.Size())
	var count int64 = 0
	for job := 0; job < s.Total; job++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			for {
				var (
					ctx = &Context{}
				)
				if s.TQueue.IsEmpty() {
					break
				}
				task := s.TQueue.Poll()
				if task == nil {
					continue
				}
				// 定义适用于压力测试的client
				t := http.DefaultTransport.(*http.Transport).Clone()
				t.MaxIdleConns = s.Total * 2
				t.MaxIdleConnsPerHost = s.Total * 2
				t.DisableKeepAlives = true
				t.DialContext = (&net.Dialer{
					Timeout:   3 * time.Second,
					KeepAlive: 3 * time.Second,
				}).DialContext
				t.IdleConnTimeout = 3 * time.Second
				t.ExpectContinueTimeout = 1 * time.Second
				client := http.Client{
					Transport: t,
					Timeout:   1 * time.Second,
				}
				switch s.Method {
				case "get", "Get", "GET":
					ctx = NewGet(task.Url, client, succeedFunc, reqTimeout, reqTimeoutMs, header)
				case "post", "Post", "POST":
					ctx = NewPost(task.Url, []byte(s.JsonData), s.ContentType, client, succeedFunc, reqTimeout, reqTimeoutMs, header)
				default:
					log.Println("暂时不支持的 Method.")
				}
				if ctx == nil {
					continue
				}
				ctx.JobNumber = i
				if !s.isRetry {
					ctx.CloseRetry()
				}
				ctx.Do()

				atomic.AddInt64(&s.sumReqTime, int64(ctx.Ms))
				atomic.AddInt64(&count, int64(1))
				s.stateCodeListMux.Lock()
				if ctx.Resp != nil {
					s.stateCodeList = append(s.stateCodeList, &stateCodeData{
						Code:    ctx.Resp.StatusCode,
						ReqTime: int64(ctx.Ms),
					})
				} else {
					//请求错误
					s.stateCodeList = append(s.stateCodeList, &stateCodeData{
						Code:    -1,
						ReqTime: int64(ctx.Ms),
					})
				}
				s.stateCodeListMux.Unlock()
			}

		}(job)
	}
	wg.Wait()
	Info("执行次数 : ", count)

	var (
		maxTime int64 = 0
		minTime int64 = 9999999999
	)

	fb := make(map[int]int)
	for _, v := range s.stateCodeList {
		if v.ReqTime >= maxTime {
			maxTime = v.ReqTime
		}
		if v.ReqTime <= minTime {
			minTime = v.ReqTime
		}

		if _, ok := fb[v.Code]; ok {
			fb[v.Code]++
		} else {
			fb[v.Code] = 1
		}
		//s.sumReqTime = s.sumReqTime + v.ReqTime
	}

	Info("状态码分布: ", fb)
	avg := float64(s.sumReqTime) / float64(s.Sum)
	avg = avg / (1000 * 1000)
	Info("平均用时： ", avg, "ms")
	Info("最高用时: ", float64(maxTime)/(1000*1000), "ms")
	Info("最低用时: ", float64(minTime)/(1000*1000), "ms")
	Info("执行完成！！！")
}

package gathertool

import (
	"log"
	"sync"
	"sync/atomic"
	"time"
)


type stateCodeData struct {
	Code int
	ReqTime int64
}

// StressUrl 压力测试一个url
type StressUrl struct {
	Url string
	Method string
	Sum int64
	Total int
	TQueue TodoQueue

	// 请求时间累加
	sumReqTime int64

	// 测试结果
	avgReqTime time.Duration

	// 接口传入的json
	JsonData string

	stateCodeList []*stateCodeData
	stateCodeListMux *sync.Mutex
}

// NewTestUrl 实例化一个新的url压测
func NewTestUrl(url, method string, sum int64, total int) *StressUrl {
	return &StressUrl{
		Url : url,
		Method : method,
		Sum : sum,
		Total : total,
		TQueue : NewQueue(),
		sumReqTime: int64(0),
		stateCodeList: make([]*stateCodeData,0),
		stateCodeListMux: &sync.Mutex{},
	}
}

// Run 运行压测
func (s *StressUrl) Run(vs ...interface{}){
	//解析可变参
	var (
		succeedFunc SucceedFunc
		n int64
		wg sync.WaitGroup
	    reqTimeout ReqTimeOut
		reqTimeoutms ReqTimeOutMs
	)
	for _, v := range vs {
		//log.Println("参数： ", v)
		switch vv := v.(type) {
		// 使用方传入了 header
		case SucceedFunc:
			succeedFunc = vv
			//log.Println("成功的方法", vv)
		case ReqTimeOut:
			reqTimeout = vv
		case ReqTimeOutMs:
			reqTimeoutms = vv
		}
	}

	//初始化队列
	for n=0; n<s.Sum; n++{
		s.TQueue.Add(&Task{Url: s.Url})
	}
	log.Println("总执行次数： ", s.TQueue.Size())

	var count int64 = 0

	for job:=0; job<s.Total; job++{
		wg.Add(1)
		go func(i int){
			//log.Println("启动第",i ,"个任务; ")
			defer wg.Done()
			for {

				var (
					ctx = &Context{}
					err error
				)

				if s.TQueue.IsEmpty(){
					break
				}

				task := s.TQueue.Poll()
				if task == nil{
					continue
				}

				//log.Println("第",i,"个任务取的值： ", task)
				switch s.Method {
					case "get","Get","GET":
						ctx, err = Get(task.Url, succeedFunc, reqTimeout, reqTimeoutms)
					case "post","Post","POST":
						ctx, err = PostJson(task.Url, s.JsonData, succeedFunc, reqTimeout, reqTimeoutms)
					default:
						log.Println("未知 Method.")
				}

				//ctx, err := Get(task.Url, succeedFunc, reqTimeout, reqTimeoutms)
				if err != nil {
					log.Println(err)
					return
				}

				if ctx == nil {
					continue
				}

				ctx.JobNumber = i
				ctx.Do()
				atomic.AddInt64(&s.sumReqTime, int64(ctx.Ms))
				atomic.AddInt64(&count, int64(1))

				s.stateCodeListMux.Lock()
				if ctx.Resp != nil{
					s.stateCodeList = append(s.stateCodeList, &stateCodeData{
						Code: ctx.Resp.StatusCode,
						ReqTime: int64(ctx.Ms),
					})
				}else{
					//请求错误
					s.stateCodeList = append(s.stateCodeList, &stateCodeData{
						Code: -1,
						ReqTime: int64(ctx.Ms),
					})
				}
				s.stateCodeListMux.Unlock()

			}

		}(job)
	}
	wg.Wait()

	log.Println("执行次数 : ", count)

	fb := make(map[int]int,0)
	for _, v := range s.stateCodeList{
		//if int64(k) >= s.Sum{
		//	break
		//}
		if _,ok := fb[v.Code]; ok{
			fb[v.Code]++
		}else{
			fb[v.Code] = 1
		}
		//s.sumReqTime = s.sumReqTime + v.ReqTime
	}
	log.Println("状态码分布: ", fb)

	avg := float64(s.sumReqTime)/float64(s.Sum)
	avg = avg/(1000*1000)
	log.Println("平均用时： ", avg,"ms")

	log.Println("执行完成！！！")

}
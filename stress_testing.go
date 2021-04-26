package gathertool

import (
	"log"
	"sync"
	"sync/atomic"
	"time"
)

type stateCodeMap struct {
	m map[int]int64
	mux sync.Mutex
}

// StressUrl 压力测试一个url
type StressUrl struct {
	Url string
	Method string
	Sum int64
	Total int
	TQueue TodoQueue
	TimeOut int64

	// 请求时间累加
	sumReqTime int64

	// 测试结果
	avgReqTime time.Duration

	// state code
	stateCodeMap *stateCodeMap

}

// NewTestUrl 实例化一个新的url压测
func NewTestUrl(url, method string, sum int64, total int, timeOut int64) *StressUrl {
	return &StressUrl{
		Url : url,
		Method : method,
		Sum : sum,
		Total : total,
		TimeOut : timeOut,
		TQueue : NewQueue(),
		sumReqTime: int64(0),
		stateCodeMap: &stateCodeMap{
			m: make(map[int]int64),
			mux: sync.Mutex{},
		},
	}
}

// Run 运行压测
func (s *StressUrl) Run(vs ...interface{}){
	//解析可变参
	var (
		succeedFunc SucceedFunc
		n int64
		wg sync.WaitGroup
	)
	for _, v := range vs {
		log.Println("参数： ", v)
		switch vv := v.(type) {
		// 使用方传入了 header
		case SucceedFunc:
			succeedFunc = vv
			log.Println("成功的方法", vv)
		}
	}

	//初始化队列
	for n=0; n<s.Sum; n++{
		s.TQueue.Add(&Task{Url: s.Url})
	}

	for job:=0; job<s.Total; job++{
		wg.Add(1)
		go func(i int){
			log.Println("启动第",i ,"个任务; ")
			defer wg.Done()
			for {
				if s.TQueue.IsEmpty(){
					break
				}
				task := s.TQueue.Poll()
				if task == nil{
					continue
				}
				//log.Println("第",i,"个任务取的值： ", task)

				ctx, err := Get(task.Url, succeedFunc, ReqTimeOutMs(s.TimeOut))
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


				s.stateCodeMap.mux.Lock()
				if ctx.Resp != nil{
					s.stateCodeMap.m[ctx.Resp.StatusCode]++
				}else{
					//请求错误
					s.stateCodeMap.m[-1]++
				}
				s.stateCodeMap.mux.Unlock()


			}
			log.Println("第",i ,"个任务结束！！")
		}(job)
	}
	wg.Wait()
	avg := float64(s.sumReqTime)/float64(s.Sum)
	avg = avg/(1000*1000)
	log.Println("平均用时： ", avg,"ms")
	log.Println("状态码分布: ", s.stateCodeMap)
	log.Println("执行完成！！！")

}
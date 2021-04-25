package gathertool

import (
	"log"
	"sync"
)

// StressUrl 压力测试一个url
type StressUrl struct {
	Url string
	Method string
	Sum int64
	Total int
	TQueue TodoQueue
}

func NewTestUrl(url, method string, sum int64, total int) *StressUrl {
	return &StressUrl{
		Url : url,
		Method : method,
		Sum : sum,
		Total : total,
		TQueue : NewQueue(),
	}
}

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
				log.Println("第",i,"个任务取的值： ", task)
				ctx, err := Get(task.Url, succeedFunc)
				ctx.JobNumber = i
				if err != nil {
					log.Println(err)
					return
				}
				ctx.Do()
			}
			log.Println("第",i ,"个任务结束！！")
		}(job)
	}
	wg.Wait()
	log.Println("执行完成！！！")

}
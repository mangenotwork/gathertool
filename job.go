/*
*	Description : 并发工作任务  TODO 测试
*	Author 		: ManGe
*	Mail 		: 2912882908@qq.com
**/

package gathertool

import (
	"errors"
	"net/http"
	"runtime"
	"sync"
	"time"
)

// TodoQueue 任务队列
type TodoQueue interface {
	Add(task *Task) error // 向队列中添加元素
	Poll() *Task          // 移除队列中最前面的元素
	Clear() bool          // 清空队列
	Size() int            // 获取队列的元素个数
	IsEmpty() bool        // 判断队列是否为空
	Print()               // 打印
}

// Queue 队列
type Queue struct {
	mux  *sync.Mutex
	list []*Task
}

// Task 任务对象
type Task struct {
	Url       string
	JsonParam string
	HeaderMap *http.Header
	Data      map[string]any // 上下文传递的数据
	Urls      []*ReqUrl      // 多步骤使用
	Type      string         // "", "upload", "do"
	SavePath  string
	SaveDir   string
	FileName  string
	once      *sync.Once
}

// NewTask 新建任务
func NewTask() *Task {
	return &Task{
		Data:      make(map[string]any),
		Urls:      make([]*ReqUrl, 0),
		once:      &sync.Once{},
		HeaderMap: &http.Header{},
	}
}

func (task *Task) GetDataStr(key string) string {
	v, ok := task.Data[key]
	if ok {
		if vStr, yes := v.(string); yes {
			return vStr
		}
	}
	return ""
}

func (task *Task) AddData(key string, value any) *Task {
	task.once.Do(func() {
		task.Data = make(map[string]any)
	})
	task.Data[key] = value
	return task
}

func (task *Task) SetUrl(urlStr string) *Task {
	task.Url = urlStr
	return task
}

func (task *Task) SetJsonParam(jsonStr string) *Task {
	task.JsonParam = jsonStr
	return task
}

func CrawlerTask(url, jsonParam string, vs ...any) *Task {
	header := &http.Header{}

	for _, v := range vs {
		switch vv := v.(type) {
		case http.Header:
			for key, values := range vv {
				for _, value := range values {
					header.Add(key, value)
				}
			}
		case *http.Header:
			for key, values := range *vv {
				for _, value := range values {
					header.Add(key, value)
				}
			}
		}
	}

	return &Task{
		Url:       url,
		JsonParam: jsonParam,
		HeaderMap: header,
		Data:      make(map[string]any),
		Type:      "do",
	}
}

// ReqUrl 单个请求地址对象
type ReqUrl struct {
	Url    string
	Method string
	Params map[string]any
}

// NewQueue 新建一个队列
func NewQueue() TodoQueue {
	list := make([]*Task, 0)
	return &Queue{list: list, mux: &sync.Mutex{}}
}

// Add 向队列中添加元素
func (q *Queue) Add(task *Task) error {
	q.mux.Lock()
	defer q.mux.Unlock()
	q.list = append(q.list, task)
	return nil
}

// Poll 移除队列中最前面的额元素
func (q *Queue) Poll() *Task {
	q.mux.Lock()
	defer q.mux.Unlock()
	if q.IsEmpty() {
		return nil
	}

	first := q.list[0]
	q.list = q.list[1:]
	return first
}

func (q *Queue) Clear() bool {
	if q.IsEmpty() {
		return false
	}
	for i := 0; i < q.Size(); i++ {
		q.list[i].Url = ""
	}
	q.list = nil
	return true
}

func (q *Queue) Size() int {
	return len(q.list)
}

func (q *Queue) IsEmpty() bool {
	if len(q.list) == 0 {
		return true
	}
	return false
}

func (q *Queue) Print() {
	Info(q.list)
}

// UploadQueue 下载队列
type UploadQueue struct {
	mux  *sync.Mutex
	list []*Task
}

// NewUploadQueue 新建一个下载队列
func NewUploadQueue() TodoQueue {
	list := make([]*Task, 0)
	return &UploadQueue{list: list, mux: &sync.Mutex{}}
}

// Add 向队列中添加元素
func (q *UploadQueue) Add(task *Task) error {
	if task.Url == "" {
		return errors.New("task url is null")
	}
	if task.SavePath == "" && (task.SaveDir == "" || task.FileName == "") {
		return errors.New("save path is null")
	}
	q.mux.Lock()
	defer q.mux.Unlock()
	q.list = append(q.list, task)
	return nil
}

// Poll 移除队列中最前面的额元素
func (q *UploadQueue) Poll() *Task {
	q.mux.Lock()
	defer q.mux.Unlock()
	if q.IsEmpty() {
		return nil
	}

	first := q.list[0]
	q.list = q.list[1:]
	return first
}

func (q *UploadQueue) Clear() bool {
	if q.IsEmpty() {
		return false
	}
	for i := 0; i < q.Size(); i++ {
		q.list[i].Url = ""
	}
	q.list = nil
	return true
}

func (q *UploadQueue) Size() int {
	return len(q.list)
}

func (q *UploadQueue) IsEmpty() bool {
	if len(q.list) == 0 {
		return true
	}
	return false
}

func (q *UploadQueue) Print() {
	Info(q.list)
}

// TODO 使用双向链表实现队列

// StartJob 开始运行并发
func StartJob(jobNumber int, queue TodoQueue, f func(task *Task)) {
	var wg sync.WaitGroup
	for job := 0; job < jobNumber; job++ {
		wg.Add(1)
		go func(i int) {
			Info("启动第", i, "个任务")
			defer wg.Done()
			for {
				if queue.IsEmpty() {
					break
				}
				task := queue.Poll()
				Info("第", i, "个任务取的值： ", task)
				f(task)
			}
			Info("第", i, "个任务结束！！")
		}(job)
	}
	wg.Wait()
	Info("执行完成！！！")
}

// Err2Retry 设置遇到错误执行 Retry 事件
type Err2Retry bool

// StartJobGet 并发执行Get,直到队列任务为空
// jobNumber 并发数，
// queue 全局队列，
func StartJobGet(jobNumber int, queue TodoQueue, vs ...any) {

	var (
		client    *http.Client
		succeed   SucceedFunc
		retry     RetryFunc
		failed    FailedFunc
		err2Retry Err2Retry
		sleep     Sleep
	)

	runtime.GOMAXPROCS(runtime.NumCPU())

	for _, v := range vs {
		switch vv := v.(type) {
		case *http.Client:
			client = vv
		case SucceedFunc:
			succeed = vv
		case FailedFunc:
			failed = vv
		case RetryFunc:
			retry = vv
		case Err2Retry:
			err2Retry = vv
		case Sleep:
			sleep = vv
			Info("设置 sleep = ", sleep)
		}
	}

	var wg sync.WaitGroup
	for job := 0; job < jobNumber; job++ {
		wg.Add(1)
		go func(i int) {
			Info("启动第", i, "个任务")
			defer wg.Done()
			for {
				if queue.IsEmpty() {
					break
				}
				task := queue.Poll()
				Info("第", i, "个任务取的值： ", task)
				ctx := NewGet(task.Url, task)
				if client != nil {
					ctx.Client = client
				}
				if succeed != nil {
					ctx.SetSucceedFunc(succeed)
				}
				if retry != nil {
					ctx.SetRetryFunc(retry)
				}
				if failed != nil {
					ctx.SetFailedFunc(failed)
				}
				if err2Retry {
					ctx.OpenErr2Retry()
				}
				if sleep != 0 {
					time.Sleep(time.Duration(sleep) * time.Second)
				}
				switch task.Type {
				case "", "do":
					ctx.Do()
				case "upload":
					if task.SavePath == "" {
						task.SavePath = task.SaveDir + task.FileName
					}
					ctx.Upload(task.SavePath)
				default:
					ctx.Do()
				}
			}
			Info("第", i, "个任务结束！！")
		}(job)
	}
	wg.Wait()
	Info("执行完成！！！")
}

// StartJobPost 开始运行并发Post
func StartJobPost(jobNumber int, queue TodoQueue, vs ...any) {
	var (
		client    *http.Client
		succeed   SucceedFunc
		retry     RetryFunc
		failed    FailedFunc
		err2Retry Err2Retry
		sleep     Sleep
	)

	for _, v := range vs {
		switch vv := v.(type) {
		case *http.Client:
			client = vv
		case SucceedFunc:
			succeed = vv
		case FailedFunc:
			failed = vv
		case RetryFunc:
			retry = vv
		case Err2Retry:
			err2Retry = vv
		case Sleep:
			sleep = vv
			Info("设置 sleep = ", sleep)
		}
	}

	var wg sync.WaitGroup
	for job := 0; job < jobNumber; job++ {
		wg.Add(1)
		go func(i int) {
			Info("启动第", i, "个任务")
			defer wg.Done()
			for {
				if queue.IsEmpty() {
					break
				}
				task := queue.Poll()
				Info("第", i, "个任务取的值： ", task, task.HeaderMap)
				ctx := NewPost(task.Url, []byte(task.JsonParam), "application/json;", task)
				if client != nil {
					ctx.Client = client
				}
				if succeed != nil {
					ctx.SetSucceedFunc(succeed)
				}
				if retry != nil {
					ctx.SetRetryFunc(retry)
				}
				if failed != nil {
					ctx.SetFailedFunc(failed)
				}
				if err2Retry {
					ctx.OpenErr2Retry()
				}
				if sleep != 0 {
					ctx.sleep = sleep
				}
				switch task.Type {
				case "", "do":
					ctx.Do()
				case "upload":
					if task.SavePath == "" {
						task.SavePath = task.SaveDir + task.FileName
					}
					ctx.Upload(task.SavePath)
				default:
					ctx.Do()
				}

			}
			Info("第", i, "个任务结束！！")
		}(job)
	}
	wg.Wait()
	Info("执行完成！！！")
}

// CPUMax 多核执行
func CPUMax() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

// TODO: 优化并重构

/*
	Description : 任务队列
	Author : ManGe
	Mail : 2912882908@qq.com
	Github : https://github.com/mangenotwork/gathertool
*/

package gathertool

import (
	"errors"
	"log"
	"net/http"
	"sync"
)

// TodoQueue 任务队列
type TodoQueue interface {
	Add(task *Task) error //向队列中添加元素
	Poll() *Task          //移除队列中最前面的元素
	Clear() bool          //清空队列
	Size() int            //获取队列的元素个数
	IsEmpty() bool        //判断队列是否是空
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
	Data      map[string]interface{} // 上下文传递的数据
	Urls      []*ReqUrl              // 多步骤使用
	Type      string                 // "", "upload", "do"
	SavePath  string
	SaveDir   string
	FileName  string
	once      *sync.Once
}

// NewTask 新建任务
func NewTask() *Task {
	return &Task{
		Data:      make(map[string]interface{}),
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

func (task *Task) AddData(key string, value interface{}) *Task {
	task.once.Do(func() {
		task.Data = make(map[string]interface{})
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

func CrawlerTask(url, jsonParam string, vs ...interface{}) *Task {
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
		Data:      make(map[string]interface{}, 0),
		Type:      "do",
	}
}

// ReqUrl 单个请求地址对象
type ReqUrl struct {
	Url    string
	Method string
	Params map[string]interface{}
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
		return errors.New("Task url is null.")
	}
	if task.SavePath == "" && (task.SaveDir == "" || task.FileName == "") {
		return errors.New("save path is null.")
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
	log.Println(q.list)
}

// todo 使用双向链表实现队列

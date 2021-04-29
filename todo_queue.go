/*
	Description : 任务队列
	Author : ManGe
	Version : v0.1
	Date : 2021-04-24
*/

package gathertool

import (
	"errors"
	"fmt"
	"log"
	"sync"
)

type TodoQueue interface {
	Add(task *Task) error  //向队列中添加元素
	Poll()   *Task  //移除队列中最前面的元素
	Clear()  bool   //清空队列
	Size()  int     //获取队列的元素个数
	IsEmpty() bool  //判断队列是否是空
	Print() // 打印
}

// 队列
type Queue struct {
	mux *sync.Mutex
	list []*Task
}

// 任务对象
type Task struct {
	Url string
	// 上下文传递的数据
	Data map[string]interface{}
	Urls []*ReqUrl // 多步骤使用
	Type string // "", "upload", "do"
	SavePath string
	SaveDir string
	FileName string
}

// 单个请求地址对象
type ReqUrl struct {
	Url string
	Method  string
	Params  map[string]interface{}
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
	q.list = append(q.list,task)
	return nil
}

// Poll 移除队列中最前面的额元素
func (q *Queue) Poll() *Task {
	q.mux.Lock()
	defer q.mux.Unlock()
	if q.IsEmpty() {
		fmt.Println("queue is empty!")
		return nil
	}

	first := q.list[0]
	q.list = q.list[1:]
	return first
}

func (q *Queue) Clear() bool {
	if q.IsEmpty() {
		fmt.Println("queue is empty!")
		return false
	}
	for i:=0 ; i< q.Size() ; i++ {
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
	log.Println(q.list)
}


// 下载队列
type UploadQueue struct {
	mux *sync.Mutex
	list []*Task
}

// NewQueue 新建一个队列
func NewUploadQueue() TodoQueue {
	list := make([]*Task, 0)
	return &UploadQueue{list: list, mux: &sync.Mutex{}}
}

// Add 向队列中添加元素
func (q *UploadQueue) Add(task *Task) error {
	if task.Url == "" {
		return errors.New("Task url is null.")
	}
	if task.SavePath == "" && (task.SaveDir == "" || task.FileName == ""){
		return errors.New("save path is null.")
	}
	q.mux.Lock()
	defer q.mux.Unlock()
	q.list = append(q.list,task)
	return nil
}

// Poll 移除队列中最前面的额元素
func (q *UploadQueue) Poll() *Task {
	q.mux.Lock()
	defer q.mux.Unlock()
	if q.IsEmpty() {
		fmt.Println("queue is empty!")
		return nil
	}

	first := q.list[0]
	q.list = q.list[1:]
	return first
}

func (q *UploadQueue) Clear() bool {
	if q.IsEmpty() {
		fmt.Println("queue is empty!")
		return false
	}
	for i:=0 ; i< q.Size() ; i++ {
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
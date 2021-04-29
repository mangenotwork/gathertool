package main

import (
	gt "github.com/mangenotwork/gathertool"
	"log"
)

var fileLinks = map[string]string{
	"https://jyzd.bfsu.edu.cn/uploadfile/bfsu/front/default/upload_file_35256.pdf":"/home/mange/Desktop/upload_file_1.pdf",
	"http://www.cuc.edu.cn/_upload/article/files/c8/44/741ea34245acbc7b8aa9b7660c0f/579130b8-3974-4c9d-99a3-19f4e4e095d1.pdf":"/home/mange/Desktop/upload_file_2.pdf",
	"https://job.buct.edu.cn/_upload/article/files/a0/b9/84fd06284b078e8051163a221546/2e302909-e81e-4aea-aa94-c5f37179cab0.pdf":"/home/mange/Desktop/upload_file_3.pdf",
	"http://kr.shanghai-jiuxin.com/file/2020/1031/774218be86d832f359637ab120eba52d.jpg":"/home/mange/Desktop/upload_file_1.jpg",
	"http://job.szu.edu.cn/UpLoadFile/InfoCenter/14247/File/20210119100014406.pdf":"/home/mange/Desktop/upload_file_4.pdf",
}

func main(){
	// 普通下载
	err := gt.Upload("https://jyzd.bfsu.edu.cn/uploadfile/bfsu/front/default/upload_file_35256.pdf", "/home/mange/Desktop/upload_file_3.pdf")
	log.Println(err)
	//for k,v := range fileLinks{
	//	gt.Upload(k,v)
	//}


	//// 并发下载
	//queue := gt.NewUploadQueue()
	//for k,v := range fileLinks{
	//	queue.Add(&gt.Task{
	//		Url: k,
	//		SavePath: v,
	//		Type: "upload",
	//	})
	//}
	//gt.StartJobGet(5, queue)

}
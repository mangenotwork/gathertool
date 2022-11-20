package main

import (
	gt "github.com/mangenotwork/gathertool"
	"log"
)

func main() {
	getSchoolList()

}

func getSchoolList() {
	caseUrl := "https://api.eol.cn/gkcx/api/?access_token=&admissions=&central=&department=&dual_class=&f211=&f985=&is_doublehigh=&is_dual_class=&keyword=&nature=&page=1&province_id=&ranktype=&request_type=1&school_type=&signsafe=&size=20&sort=view_total&top_school_id=[]&type=&uri=apidata/api/gk/school/lists"
	//ctx := gt.PostJson(caseUrl, "")//.SetSucceedFunc(getSchoolListSuueed)
	//log.Println(ctx.String())

	gt.PostJson(caseUrl, "", gt.SucceedFunc(getSchoolListSuueed))

}

func getSchoolListSuueed(ctx *gt.Context) {
	//log.Println("ctx = ", ctx)
	log.Println(ctx.RespBodyString())
	log.Println(ctx.CheckReqMd5())
	log.Println(ctx.CheckMd5())
}

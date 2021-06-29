package main

import (
	gt "github.com/mangenotwork/gathertool"
	"log"
)

func main(){
	getSchoolList()

}



func getSchoolList(){
	caseUrl := "https://api.eol.cn/gkcx/api/?access_token=&admissions=&central=&department=&dual_class=&f211=&f985=&is_doublehigh=&is_dual_class=&keyword=&nature=&page=1&province_id=&ranktype=&request_type=1&school_type=&signsafe=&size=20&sort=view_total&top_school_id=[]&type=&uri=apidata/api/gk/school/lists"
	cxt, err := gt.PostJson(caseUrl, "")
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(cxt.String())
	log.Println(cxt.CheckReqMd5())
	log.Println(cxt.CheckMd5())
}
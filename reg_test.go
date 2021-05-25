package gathertool

import (
	"testing"
)

func TestReg(t *testing.T){

	//list := RegFindAll(`<option(.*?)</option>`, txt)
	//log.Println(list)


	//list := RegHtmlInput(txt)
	//list := RegHtmlA(txt)
	//list := RegHtmlTr(txt)
	//
	//for k,v := range list {
	//	log.Println(k, " ---> ", v)
	//}
}

var txt = `<div class="width1000">

<table class="table-top yxk-filter">
    <form action="/zsgs/zhangcheng/listVerifedZszc.do" id="form1"></form>
        <input type="hidden" name="method" value="index">
        <tbody><tr>
            <td class="right-85">院校所在地：</td>
            <td>
<a href="aaaadddddddd">ffdg
fdgdg</a>
                <select name="ssdm" class="ch-hide">
                        <option value="" selected="selected">全部</option>
                        <option value="11">
                            北京
                        </option>
                        <option value="" selected="selected">全部</option>
                        <option value="12">
                            天津
                        </option>
                        <option value="" selected="selected">全部</option>
                        <option value="13">
                            河北
                        </option>
                        <option value="" selected="selected">全部</option>
                        <option value="14">
                            山西
                        </option>
                        <option value="" selected="selected">全部</option>
                        <option value="15">
                            内蒙古
                        </option>
                        <option value="" selected="selected">全部</option>
                        <option value="21">
                            辽宁
                        </option>
                        <option value="" selected="selected">全部</option>
                        <option value="22">
                            吉林
                        </option>
                        <option value="" selected="selected">全部</option>
                        <option value="23">
                            黑龙江
                        </option>
                        <option value="" selected="selected">全部</option>
                        <option value="31">
                            上海
                        </option>
                        <option value="" selected="selected">全部</option>
                        <option value="32">
                            江苏
                        </option>
                        <option value="" selected="selected">全部</option>
                        <option value="33">
                            浙江
                        </option>
                        <option value="" selected="selected">全部</option>
                        <option value="34">
                            安徽
                        </option>
                        <option value="" selected="selected">全部</option>
                        <option value="35">
                            福建
                        </option>
                        <option value="" selected="selected">全部</option>
                        <option value="36">
                            江西
                        </option>
                        <option value="" selected="selected">全部</option>
                        <option value="37">
                            山东
                        </option>
                        <option value="" selected="selected">全部</option>
                        <option value="41">
                            河南
                        </option>
                        <option value="" selected="selected">全部</option>
                        <option value="42">
                            湖北
                        </option>
                        <option value="" selected="selected">全部</option>
                        <option value="43">
                            湖南
                        </option>
                        <option value="" selected="selected">全部</option>
                        <option value="44">
                            广东
                        </option>
                        <option value="" selected="selected">全部</option>
                        <option value="45">
                            广西
                        </option>
                        <option value="" selected="selected">全部</option>
                        <option value="46">
                            海南
                        </option>
                        <option value="" selected="selected">全部</option>
                        <option value="50">
                            重庆
                        </option>
                        <option value="" selected="selected">全部</option>
                        <option value="51">
                            四川
                        </option>
                        <option value="" selected="selected">全部</option>
                        <option value="52">
                            贵州
                        </option>
                        <option value="" selected="selected">全部</option>
                        <option value="53">
                            云南
                        </option>
                        <option value="" selected="selected">全部</option>
                        <option value="54">
                            西藏
                        </option>
                        <option value="" selected="selected">全部</option>
                        <option value="61">
                            陕西
                        </option>
                        <option value="" selected="selected">全部</option>
                        <option value="62">
                            甘肃
                        </option>
                        <option value="" selected="selected">全部</option>
                        <option value="63">
                            青海
                        </option>
                        <option value="" selected="selected">全部</option>
                        <option value="64">
                            宁夏
                        </option>
                        <option value="" selected="selected">全部</option>
                        <option value="65">
                            新疆
                        </option>
                        <option value="" selected="selected">全部</option>
                        <option value="81">
                            香港
                        </option>
                        <option value="" selected="selected">全部</option>
                        <option value="91">
                            澳门
                        </option>
                        <option value="" selected="selected">全部</option>
                        <option value="71">
                            台湾
                        </option>
                </select>
            <td>
                <tr class="ch-hide" name="yxls">
                    <option value="" selected="selected">全部</option>
                    <option value="moe">教育部</option>
                    <option value="min">其他部委</option>
                    <option value="loc">地方</option>
                    <option value="army">军校</option>
                </tr>
                <span class="yxk-all js-option js-ls selected" data-id="">全部</span>
                <span class="yxk-option js-option " data-id="moe">教育部</span>
                <span class="yxk-option js-option " data-id="min">其他部委</span>
                <span class="yxk-option js-option " data-id="loc">地方</span>
                <span class="yxk-option js-option " data-id="army">军校</span>
            </td>
        </tr>
        <tr>
            <td class="right-85">院校类型：</td>
            <td>
                <select class="ch-hide" name="yxlx">
                    <option value="">院校类型</option>
                    <option value="01">综合</option>
                    <option value="02">工科</option>
                    <option value="03">农业</option>
                    <option value="04">林业</option>
                    <option value="05">医药</option>
                    <option value="06">师范</option>
                    <option value="07">语言</option>
                    <option value="08">财经</option>
                    <option value="09">政法</option>
                    <option value="10">体育</option>
                    <option value="11">艺术</option>
                    <option value="12">民族</option>
                </select>
                <span class="yxk-all js-option js-ls selected" data-id="">全部</span>
                <span class="yxk-option js-option " data-id="01">综合</span>
                <span class="yxk-option js-option " data-id="02">工科</span>
                <span class="yxk-option js-option " data-id="03">农业</span>
                <span class="yxk-option js-option " data-id="04">林业</span>
                <span class="yxk-option js-option " data-id="05">医药</span>
                <span class="yxk-option js-option " data-id="06">师范</span>
                <span class="yxk-option js-option " data-id="07">语言</span>
                <span class="yxk-option js-option " data-id="08">财经</span>
                <span class="yxk-option js-option " data-id="09">政法</span>
                <span class="yxk-option js-option " data-id="10">体育</span>
                <span class="yxk-option js-option " data-id="11">艺术</span>
                <span class="yxk-option js-option " data-id="12">民族</span>
            </td>

        </tr>
        <tr>
            <td class="right-85">学历层次：</td>
            <td>
                <select class="ch-hide" name="xlcc">
                    <option value="" selected="selected">全部</option>
                    <option value="bk">本科</option>
                    <option value="gzzk">高职(专科)</option>
                </select>
                <span class="yxk-all js-option js-ls selected" data-id="">全部</span>
                <span class="yxk-option js-option " data-id="bk">本科</span>
                <span class="yxk-option js-option " data-id="gzzk">高职(专科)</span>
            </td>
        </tr>
        <tr>
            <td class="right-85">院校特性：</td>
            <td class="yxk-xz">
                <label class="ch-check-label ">
                    <input type="checkbox" name="zgsx" id="ylxx" value="ylxx">一流大学建设高校
                </label>
                <label class="ch-check-label ">
                    <input type="checkbox" name="zgsx" id="ylxk" value="ylxk">一流学科建设高校
                </label>
                <label class="ch-check-label ">
                    <input type="checkbox" name="zgsx" id="yjsy" value="yjsy">研究生院
                </label>
                <label class="ch-check-label ">
                    <input type="radio" name="yxjbz" id="yxjbz" value="2">民办高校
                </label>
                <label class="ch-check-label ">
                    <input type="radio" name="yxjbz" id="yxjbz" value="3">独立学院
                </label>
                <label class="ch-check-label ">
                    <input type="radio" name="yxjbz" id="yxjbz" value="4">中外合作办学
                </label>
                <label class="ch-check-label ">
                    <input type="radio" name="yxjbz" id="yxjbz" value="5">内地与港澳台地区合作办学
                </label>
            </td>
        </tr>
    <tr>
        <td colspan="2">
            <div class="input-group">
                <form name="yxmcSearhc" action="/zsgs/zhangcheng/listVerifedZszc.do" method="post">
                    <input type="hidden" name="method" id="method" value="listInfoByYxmc">
                    直接搜索院校 <input name="yxmc" id="yxmc" value="请输入院校名称搜索" size="18" type="text" onfocus="setYxmcFocus('请输入院校名称搜索')" onblur="setYxmcBlue('请输入院校名称搜索')">
                    <button type="submit" class="button_grey" onclick="return checkYxmc()" style="background: #CCCCCC;color: #333333;">
                    搜索</button>
                </form>
            </div>
        </td>
    </tr>
</tbody></table>
<div class="info">
                    <div class="hd">
                        <a href="https://movie.douban.com/subject/1291546/" class="">
                            <span class="title">霸王别姬</span>
                                <span class="other">&nbsp;/&nbsp;再见，我的妾  /  Farewell My Concubine</span>
                        </a>


`


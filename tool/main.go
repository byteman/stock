package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"encoding/json"
	"os"
	"bufio"
	"strings"
	"io"
	"gopkg.in/ini.v1"
	"time"
	"math/rand"
)

type PEGINFO struct {
	PEG struct {
		Dates  []string  `json:"dates"`
		Values []float32 `json:"values"`
	} `json:"PEG"`
}
type PEINFO struct{
	PE struct {
		Dates  []string  `json:"dates"`
		Values []float32 `json:"values"`
	} `json:"市盈率"`
}
type PBINFO struct{
	PB struct {
		Dates  []string  `json:"dates"`
		Values []float32 `json:"values"`
	} `json:"市净率"`
}

type RATIOINFO struct {
	RATIO struct {
		Dates  []string  `json:"dates"`
		Values []float32 `json:"values"`
	} `json:"近4季净利润"`
}
type ROICINFO struct{
	ROIC struct {
		Dates  []string  `json:"dates"`
		Values []float32 `json:"values"`
	} `json:"ROIC"`
}
func Request(url string, result interface{})error{
	resp, err := http.Get(url)

	if err != nil {
		// handle error
		fmt.Println("Get ",err)
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
		fmt.Println("ReadAll=",err)
		return err
	}
	if err:= json.Unmarshal(body,result);err!=nil{
		fmt.Println(string(body))
		fmt.Println("Unmarshal=",err)

		return err
	}
	return nil

}
func trim(str string)string  {
	str = strings.Replace(str, " ", "", -1)
	// 去除换行符
	str = strings.Replace(str, "\n", "", -1)

	str = strings.Replace(str, "\r", "", -1)
	return str
}
func RequestPE(code string)(value float32,err error){
	var pe PEINFO
	err=Request("https://api.wayougou.com/api/ratios/stock?code="+code+"&name=%E5%B8%82%E7%9B%88%E7%8E%87&chart=%E5%B8%82%E7%9B%88%E7%8E%87&span=10",&pe)
	if err!=nil{
		fmt.Println(err)
		return 0,err
	}
	if len(pe.PE.Values) == 0{
		return 0,fmt.Errorf("没有找到%s的PE数据",code)
	}
	index:=len(pe.PE.Values)-1
	return pe.PE.Values[index],nil
}
func RequestPB(code string)(value float32,err error){
	var pb PBINFO

	err=Request("https://api.wayougou.com/api/ratios/stock?code="+code+"&name=%E5%B8%82%E5%87%80%E7%8E%87&chart=%E5%B8%82%E5%87%80%E7%8E%87&span=10",&pb)
	if err!=nil{
		fmt.Println(err)
		return 0,err
	}
	if len(pb.PB.Values) == 0{
		return 0,fmt.Errorf("没有找到%s的PB数据",code)
	}
	index:=len(pb.PB.Values)-1
	return pb.PB.Values[index],nil
}
func RequestPEG(code string)(value float32,err error){
	var peg PEGINFO
	err=Request("https://api.wayougou.com/api/ratios/stock?code="+code+"&name=PEG&chart=PEG",&peg)
	if err!=nil{
		fmt.Println(err)
		return 0,err
	}
	if len(peg.PEG.Values) == 0{
		return 0,fmt.Errorf("没有找到%s的PEG数据",code)
	}
	index:=len(peg.PEG.Values)-1
	return peg.PEG.Values[index],nil
}
func RequestROCI(code string)(value float32,err error){
	var roci ROICINFO
	err=Request("https://api.wayougou.com/api/ratios/stock?code="+code+"&name=ROIC&chart=ROIC",&roci)
	if err!=nil{
		fmt.Println(err)
		return 0,err
	}
	if len(roci.ROIC.Values) == 0{
		return 0,fmt.Errorf("没有找到%s的ROIC数据",code)
	}
	//fmt.Println(roci.ROIC.Values)
	index:=len(roci.ROIC.Values)-1
	return roci.ROIC.Values[index],nil
}
//统计后一个季度大于前一个季度的总数目
func GetUpTotal(values []float32)int {
	var total = 0
	if len(values) == 0{
		return 0
	}
	if len(values) == 1{
		return 0
	}
	for i:=0; i< len(values)-1;i++{

		if values[i+1] > values[i]{
			total++
		}
	}
	return total
}
//统计最后一个季度连续大于前面多少个季度
func GetLastUpTotal(values []float32)int {
	var total = 0
	if len(values) == 0{
		return 0
	}
	if len(values) == 1{
		return 0
	}
	for i:=len(values)-1; i > 0;i--{

		if values[i] > values[i-1]{
			total++
		}else{
			break
		}
	}
	return total
}
func RequestRATIOS(code string)(value string,err error){
	var ratio RATIOINFO
	err=Request("https://api.wayougou.com/api/ratios/stock?code="+code+"&name=%E5%87%80%E5%88%A9%E6%B6%A6&chart=%E8%BF%914%E5%AD%A3%E5%87%80%E5%88%A9%E6%B6%A6",&ratio)
	if err!=nil{
		fmt.Println(err)
		return "0.0",err
	}
	if len(ratio.RATIO.Values) == 0{
		return "0.0",fmt.Errorf("没有找到%s的RATIO数据",code)
	}
	var values []float32
	if len(ratio.RATIO.Values) >= 21{
		pos:=len(ratio.RATIO.Values) - 21
		values = ratio.RATIO.Values[pos:]
	}else{
		values = ratio.RATIO.Values
	}
	zs:=GetUpTotal(values)
	xs:=GetLastUpTotal(values)
	return fmt.Sprintf("%d.%02d",zs,xs),nil

}


//https://wayougou.com/stock/002415/净利润/近4季净利润
var target_file *os.File
var codes []string //股票代码列表
var base = 2
var min_rand = 1
var max_rand = 5
func initConfig()error{

	f, err := os.Open("wyg.txt")
	if err != nil {
		fmt.Println(err)
		return err
	}


	target_file, err = os.Create("wyg1.txt") //创建文件
	if err!=nil{
		fmt.Println(err)
		return err
	}
	defer f.Close()
	codes = nil
	rd := bufio.NewReader(f)
	fmt.Println("开始加载股票代码")
	for {
		line, err := rd.ReadString('\n') //以'\n'为结束符读入一行
		//fmt.Println("line=",line)
		if err != nil || io.EOF == err {
			fmt.Println(err)
			break
		}
		items := strings.Split(line, " ")
		//fmt.Println(items,len(items))
		//fmt.Println(items[0],items[1],items[2])
		if len(items) < 1 {
			fmt.Println(line + " length < 1")
			continue
		}
		fmt.Println("code=",items[0])
		codes=append(codes,items[0])
	}
	fmt.Println("加载股票代码完成，总共",len(codes),"只股票")

	config,err:=ini.Load("tool.ini")
	if err==nil {
		base = config.Section("config").Key("base").MustInt(5)
		min_rand = config.Section("config").Key("min").MustInt(1)
		max_rand = config.Section("config").Key("max").MustInt(5)
	}
	fmt.Printf("提取时间%d秒,最小随机时间%d秒，最大随机时间%d秒\n",base,min_rand,max_rand)
	return nil
}
func parseThread()error  {

return nil

}
func ReadThenWriteStock(code string) (err error) {
	var pe,pb,peg,roic float32
	var ratio string

	if roic,err=RequestROCI(code);err !=nil{
		fmt.Println(err)
		return
	}
	if pe,err=RequestPE(code);err!=nil{
		return
	}
	if pb,err=RequestPB(code);err!=nil{
		return
	}
	if peg,err=RequestPEG(code);err!=nil{
		return
	}
	if ratio,err=RequestRATIOS(code);err!=nil{
		return
	}
	fmt.Println("pb  =",pb)
	fmt.Println("pe  =",pe)
	fmt.Println("peg =",peg)
	fmt.Println("roic=",roic)
	fmt.Println("ratio=",ratio)
	var pbs = fmt.Sprintf("%f",pb)
	if len(pbs) > 4{
		pbs=pbs[:4]
	}
	var pes = fmt.Sprintf("%f",pe)
	if len(pes) > 4{
		pes=pes[:4]
	}
	var pegs = fmt.Sprintf("%f",peg)
	if len(pegs) > 4{
		pegs=pegs[:4]
	}
	var roics = fmt.Sprintf("%f",roic)
	if len(roics) > 4{
		roics=roics[:4]
	}

	_, err = target_file.WriteString(fmt.Sprintf("%s,%s,%s,%s,%s,%s\r\n",code,pes,pbs,pegs,roics,ratio)) //写入文件(字节数组)

	return err
}
func getRandomTime() time.Duration {
	rand.Seed(time.Now().Unix())
	diff:=max_rand -min_rand
	randm:=0
	if diff != 0{
		randm=rand.Intn(diff)
	}

	//fmt.Println("randm=",randm)
	nr:=(randm + min_rand)+base
	return time.Duration(nr)*time.Second
}
func main() {

	if err:=initConfig();err!=nil{
		fmt.Println("初始化失败:",err)
		return
	}

	var length = len(codes)
	for index,code:=range codes{
		t:=getRandomTime()
		fmt.Println("休息",t)
		time.Sleep(t)
		fmt.Printf("正在采集第%d只股票[股票代码:%s][总进度:%d%%]\n",index,code,int(index*100/length))
		ReadThenWriteStock(code)
	}

	//ReadThenWriteStock("000028")

	if target_file!=nil{
		target_file.Close()
	}
	fmt.Println("完成采集")
}
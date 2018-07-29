package controllers

import (
	"github.com/gin-gonic/gin"
	"os"
	"bufio"
	"fmt"
	"io"
	"strings"
	"sync"
	"strconv"
	"net/http"
	"github.com/axgle/mahonia"
	"io/ioutil"
)

var mutex sync.Mutex

type StockInfo struct{
	Value float64 `json:"value"`
	Name string `json:"name"`
	Code string `json:"code"`
}
var stocks map[string]StockInfo

func trim(str string)string  {
	str = strings.Replace(str, " ", "", -1)
	// 去除换行符
	str = strings.Replace(str, "\n", "", -1)

	str = strings.Replace(str, "\r", "", -1)
	return str
}
func loadStock()error  {
	f, err := os.Open("appstockdata.txt")
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer f.Close()
	stocks = make(map[string]StockInfo)
	rd := bufio.NewReader(f)
	for {
		line, err := rd.ReadString('\n') //以'\n'为结束符读入一行
		//fmt.Println("line=",line)
		if err != nil || io.EOF == err {
			fmt.Println(err)
			break
		}
		items:=strings.SplitN(line,"  ",3)
		//fmt.Println(items,len(items))
		//fmt.Println(items[0],items[1],items[2])
		if len(items) != 3{
			fmt.Println(line + " lenght not == 2")
			continue
		}
		mutex.Lock()
		var info StockInfo
		info.Code = trim(items[0])
		info.Name = trim(items[1])

		v:=trim(items[2])
		v1, err := strconv.ParseFloat(v, 32)
		if err!= nil{
			fmt.Println(err)
			continue
		}
		info.Value = v1
		stocks[info.Code] = info
		mutex.Unlock()
	}
	fmt.Println("items count=",len(stocks))
	return nil
}
type StockResult struct{
	Value float64 `json:"value"`
	Name string `json:"name"`
	Code string `json:"code"`
	Error int `json:"error"`
	Message string `json:"message"`
}
func queryStock(c *gin.Context)  {
	id:=c.Param("id")


	if value, ok := stocks[id]; ok {
		c.JSON(200,&StockResult{
			Code:id,
			Value:value.Value,
			Name:value.Name,

			Error:0,
			Message:"",
		})
	}else{
		c.JSON(200,&StockResult{
			Code:id,
			Value:value.Value,
			Name:value.Name,
			Error:1,
			Message:"not exist",
		})
	}
}
func ConvertToByte(src string, srcCode string, targetCode string) []byte {
	srcCoder := mahonia.NewDecoder(srcCode)
	srcResult := srcCoder.ConvertString(src)
	tagCoder := mahonia.NewDecoder(targetCode)
	_, cdata, _ := tagCoder.Translate([]byte(srcResult), true)
	return cdata
}
func convGbk2Utf8(file string)error{


	data,err:=ioutil.ReadFile(file)
	if err!=nil{
		return err
	}
	utf8 := ConvertToByte(string(data), "gbk", "utf8")

	return ioutil.WriteFile(file,[]byte(utf8),0666)


}
func uploadFile(c *gin.Context)  {
	file, _ := c.FormFile("file")
	fmt.Println(file.Filename)

	c.SaveUploadedFile(file,"appstockdata.txt")
	convGbk2Utf8("appstockdata.txt")
	loadStock()
	c.String(http.StatusOK, fmt.Sprintf("'%s' uploaded!", file.Filename))
}
func init()  {
	loadStock()
}
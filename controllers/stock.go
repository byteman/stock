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
)

var mutex sync.Mutex
var stocks map[string]float32
func loadStock()error  {
	f, err := os.Open("stock.txt")
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer f.Close()
	stocks = make(map[string]float32)
	rd := bufio.NewReader(f)
	for {
		line, err := rd.ReadString('\n') //以'\n'为结束符读入一行

		if err != nil || io.EOF == err {
			fmt.Println(err)
			break
		}
		items:=strings.Split(line,"  ")
		//fmt.Println(items,len(items))
		//fmt.Println(items[0],items[1],items[2])
		if len(items) != 2{
			fmt.Println("lenght not 2")
			continue
		}
		mutex.Lock()
		v:=strings.Trim(items[1],"\r\n")
		v1, err := strconv.ParseFloat(v, 32)
		if err!= nil{
			fmt.Println(err)
			continue
		}
		stocks[items[0]] = float32(v1)
		fmt.Println("key=",items[0],"value=",float32(v1))
		mutex.Unlock()
		//fmt.Println(line)
	}
	fmt.Println("items count=",len(stocks))
	return nil
}
type StockResult struct{
	Code string `json:"code"`
	Value float32 `json:"value"`
	Error int `json:"error"`
	Message string `json:"message"`
}
func queryStock(c *gin.Context)  {
	id:=c.Param("id")


	if value, ok := stocks[id]; ok {
		c.JSON(200,&StockResult{
			Code:id,
			Value:value,
			Error:0,
			Message:"",
		})
	}else{
		c.JSON(200,&StockResult{
			Code:id,
			Value:value,
			Error:1,
			Message:"not exist",
		})
	}
}

func init()  {
	loadStock()
}
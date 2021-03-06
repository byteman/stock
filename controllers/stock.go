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
	"sort"
	"gopkg.in/dgrijalva/jwt-go.v3"
	"stock/models"
)

var mutex sync.Mutex

type StockInfo struct{
	Value float64 `json:"value"` //黑马系数 //3
	Name string `json:"name"` //股票名称  //1
	Code string `json:"code"` //股票代码  //2
	Track float64 `json:"track"` //短线跟踪=多空 //4
	Step float64 `json:"step"` //峰值系数 //5
	RO float64 `json:"ro"` //ROCI基本盘系数 //8
	HGT float64 `json:"hgt"` //沪港通. //7
	PE float64 `json:"pe"` //
	PB float64 `json:"pb"`
	PEG float64 `json:"peg"`
	ShortQuery float64 `json:"short_query"` //短线查询 6列
}
var stocks map[string]StockInfo
//var namestocks map[string]StockInfo

type StockSlice []StockInfo
type StockSliceRO []StockInfo


var stockArr StockSlice
var stockArrRO StockSliceRO
var StockSliceHGT StockSlice
var stockIndex = 0
var stockRoIndex = 0
func (a StockSlice) Len() int {         // 重写 Len() 方法
	return len(a)
}
func (a StockSlice) Swap(i, j int){     // 重写 Swap() 方法
	//fmt.Println(i,j)
	a[i], a[j] = a[j], a[i]
}
func (a StockSliceRO) Less(i, j int) bool {    // 重写 Less() 方法， 从大到小排序
	return a[j].ShortQuery < a[i].ShortQuery
}
func (a StockSliceRO) Len() int {         // 重写 Len() 方法
	return len(a)
}
func (a StockSliceRO) Swap(i, j int){     // 重写 Swap() 方法
	//fmt.Println(i,j)
	a[i], a[j] = a[j], a[i]
}
func (a StockSlice) Less(i, j int) bool {    // 重写 Less() 方法， 从大到小排序
	return a[j].Value < a[i].Value
}
func trim(str string)string  {
	str = strings.Replace(str, " ", "", -1)
	// 去除换行符
	str = strings.Replace(str, "\n", "", -1)

	str = strings.Replace(str, "\r", "", -1)
	return str
}
func toFloat(str string)(v float64,err error )  {
	tmp:= trim(str)
	v1, err := strconv.ParseFloat(tmp, 32)
	if err!= nil{
		return 0,err
	}
	return v1,err
}
func loadStock()error  {
	f, err := os.Open("appstockdata.txt")
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer f.Close()
	stocks = make(map[string]StockInfo)
	stockArr=nil
	StockSliceHGT = nil
	stockArrRO = nil
	stockIndex = 0
	stockRoIndex = 0
	rd := bufio.NewReader(f)
	for {
		line, err := rd.ReadString('\n') //以'\n'为结束符读入一行
		//fmt.Println("line=",line)
		if err != nil || io.EOF == err {
			fmt.Println(err)
			break
		}
		items:=strings.Split(line,",")
		//fmt.Println(items,len(items))
		//fmt.Println(items[0],items[1],items[2])
		if len(items) < 8{
			fmt.Println(line + " length < 8")
			continue
		}
		mutex.Lock()
		var info StockInfo
		info.Code = trim(items[0])
		info.Name = trim(items[1])
		info.Value, err = toFloat(items[2])
		if err!= nil{
			fmt.Println("value convert error ",err)
			continue
		}
		info.Track, err = toFloat(items[3])
		if err!= nil{
			fmt.Println("track convert error " , err)
			continue
		}
		info.Step, err = toFloat(items[4])
		if err!= nil{
			fmt.Println("step convert error" ,err)
			continue
		}
		info.ShortQuery, err = toFloat(items[5])
		if err!= nil{
			fmt.Println("ShortQuery convert error" ,err)
			continue
		}
		info.RO, err = toFloat(items[7])
		if err!= nil{
			fmt.Println("step convert error" ,err)
			continue
		}
		info.HGT, err = toFloat(items[6])
		if err!= nil{
			fmt.Println("step convert error" ,err)
			continue
		}
		info.PE, err = toFloat(items[8])
		if err!= nil{
			fmt.Println("PE convert error" ,err)
			continue
		}
		info.PB, err = toFloat(items[9])
		if err!= nil{
			fmt.Println("PB convert error" ,err)
			continue
		}
		info.PEG, err = toFloat(items[10])
		if err!= nil{
			fmt.Println("PEG convert error" ,err)
			continue
		}
		stocks[info.Code] = info
		stocks[info.Name] = info
		stockArr = append(stockArr,info)
		stockArrRO =append(stockArrRO,info)
		if info.HGT == 1{
			StockSliceHGT=append(StockSliceHGT,info)
		}
		mutex.Unlock()
	}

	sort.Sort(stockArr)    // 按照 黑马系数 的逆序排序
	sort.Sort(stockArrRO)    // 按照 短线系数 的逆序排序
	sort.Sort(StockSliceHGT) //按照黑马系数排序
	//fmt.Println(stockArr)
	index := sort.Search(len(stockArr), func(i int) bool { return stockArr[i].Value < 100 })
	if index < len(stockArr) {
		stockIndex = index
	}
	index = sort.Search(len(stockArrRO), func(i int) bool { return stockArrRO[i].ShortQuery < 10 })
	if index < len(stockArrRO) {
		stockRoIndex = index
	}
	fmt.Println("< 100 StockIndex=",stockIndex)
	fmt.Println("< 10  ShortIndex=",stockRoIndex)
	fmt.Println("items count=",len(stocks))
	//fmt.Println(stockArr)
	return nil
}
type StockResult struct{
	Value float64 `json:"value"`
	Name string `json:"name"`
	Code string `json:"code"`
	Track float64 `json:"track"` //短线跟踪
	Step float64 `json:"step"` //峰值系数
	RO float64 `json:"ro"`
	Error int `json:"error"`
	Message string `json:"message"`
}
var hmacSampleSecret = []byte("secret key")

func checkPermission(c *gin.Context,level int)(err error){
	tokenString := c.GetHeader("token")
	var user *models.User = nil
	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return hmacSampleSecret, nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		username:=claims["id"].(string)

		user,err=models.GetUserByName(username)
		if err!=nil{
			return  err
		}

		fmt.Printf("payType=%d,level=%d\n",user.PayType,level)

		if err:=user.CheckPermission(level);err!=nil{
			return  err
		}
		if user!=nil{
			models.AddLog(user.ID,level)
		}
	} else {
		return fmt.Errorf("无法获取用户信息")
	}
	return nil
}
func advanceQueryStock(c *gin.Context){
	//score:=DefaultQueryInt(c,"score",100)

	if err:=checkPermission(c,2);err!=nil{
		SuccessQueryResponse(c,gin.H{
			"error":1,
			"message":fmt.Sprintf("%s",err),
		})
		return
	}

	c.IndentedJSON(200,gin.H{
		"error":0,
		"message":"ok",
		"data":stockArr[:stockIndex],
	})
}
//短线查询
func ShortQueryStock(c *gin.Context){
	//score:=DefaultQueryInt(c,"score",100)

	if err:=checkPermission(c,1);err!=nil{
		SuccessQueryResponse(c,gin.H{
			"error":1,
			"message":fmt.Sprintf("%s",err),
		})
		return
	}

	c.IndentedJSON(200,gin.H{
		"error":0,
		"message":"ok",
		"data":stockArrRO[:stockRoIndex],
	})
}
//互通查询
func HGTQueryStock(c *gin.Context){
	//score:=DefaultQueryInt(c,"score",100)

	if err:=checkPermission(c,2);err!=nil{
		SuccessQueryResponse(c,gin.H{
			"error":1,
			"message":fmt.Sprintf("%s",err),
		})
		return
	}

	c.IndentedJSON(200,gin.H{
		"error":0,
		"message":"ok",
		"data":StockSliceHGT,
	})
}
func basicQueryStock(c *gin.Context)  {


	if err:=checkPermission(c,1);err!=nil{
		SuccessQueryResponse(c,gin.H{
			"error":1,
			"message":fmt.Sprintf("%s",err),
		})
		return
	}

	id:=c.Param("id")

	//strconv.ParseInt()

	if value, ok := stocks[id]; ok {
		c.JSON(200,&StockResult{
			Code:id,
			Value:value.Value,
			Name:value.Name,
			Track:value.Track,
			Step:value.Step,
			RO:value.RO,
			Error:0,
			Message:"",
		})
	}else{
		c.JSON(200,&StockResult{
			Code:id,
			Value:value.Value,
			Name:value.Name,
			Track:value.Track,
			Step:value.Step,
			RO:value.RO,
			Error:1,
			Message:"该股票不存在",
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



	if file.Filename == "appstockdata.txt"{
		c.SaveUploadedFile(file,file.Filename)
		convGbk2Utf8(file.Filename)
		loadStock()
	}else if file.Filename == "stockhelp.html"{
		c.SaveUploadedFile(file,"static/dist/"+file.Filename)
		convGbk2Utf8("static/dist/" + file.Filename)
	}else {
		c.SaveUploadedFile(file,"static/dist/upload/"+file.Filename)
	}

	c.String(http.StatusOK, fmt.Sprintf("'%s' uploaded!", file.Filename))
}
func init()  {
	loadStock()
}

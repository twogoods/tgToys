package express

import (
	. "TgToys/utils"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	. "github.com/twogoods/golib/gohttp"
)

const (
	//see http://www.kdniao.com/api-track
	APPID  = "your id"
	APIKEY = "your apikey"

	ShipperInfoRequestType = "2002"
	ExpressInfoRequestType = "1002"

	Blue    = "blue"
	Red     = "red"
	Yellow  = "yellow"
	Green   = "green"
	Purple  = "purple"
	Skyblue = "skyblue"
)

type ShipperResult struct {
	EBusinessID, LogisticCode, Code string
	Success                         bool
	Shippers                        []ShipperInfo
}

type ShipperInfo struct {
	ShipperCode, ShipperName string
}

type ExpressInfo struct {
	Success bool
	State   string
	Reason  string
	Traces  []Trace
}

type Trace struct {
	AcceptTime, AcceptStation, Remark string
}

func encrypt(requestData string, AppKey string) string {
	if AppKey != "" {
		return Base64Encode(MD5(requestData + AppKey))
	}
	return Base64Encode(MD5(requestData))
}

func getShipperCode(logisticCode string) ([]byte, error) {
	requestData := "{'LogisticCode':'" + logisticCode + "'}"
	return sendPost(requestData, ShipperInfoRequestType)
}

func getExpressDetail(logisticCode string, shipperCode string) ([]byte, error) {
	requestData := "{'ShipperCode':'" + shipperCode + "','LogisticCode':'" + logisticCode + "'}"
	return sendPost(requestData, ExpressInfoRequestType)
}

func sendPost(requestData string, requestType string) ([]byte, error) {
	postbody := FormBodyBuilder().AddParam("requestData", Urlencode(requestData)).
		AddParam("EBusinessID", APPID).
		AddParam("RequestType", requestType).
		AddParam("DataSign", Urlencode(encrypt(requestData, APIKEY))).
		AddParam("DataType", "2").Build()
	req, _ := RequestBuilder().Url("http://api.kdniao.cc/Ebusiness/EbusinessOrderHandle.aspx").Post(postbody).Build()
	client := HttpClientBuilder().Build()
	resp, err := client.Execute(req)
	if err == nil {
		bytes, _ := resp.BodyByte()
		return bytes, nil
	} else {
		return nil, err
	}
}

func selectShipper(shippers []ShipperInfo) string {
	s := "当前单号对应如下快递公司:  "
	for i, v := range shippers {
		s += fmt.Sprint(colorize(strconv.Itoa(i+1), Skyblue), ".", v.ShipperName, " ")
	}
	fmt.Println(s)
	var index int
	var err error
	for {
		fmt.Print(colorize("请输入快递公司序号或q退出程序:", Green))
		var command string
		fmt.Scanf("%s", &command)
		if command == "q" || command == "Q" {
			os.Exit(0)
		}
		index, err = strconv.Atoi(command)
		if err != nil || index < 1 || len(shippers) < index {
			fmt.Println(colorize("输入数据有误,请检查重新输入...", Red))
			continue
		}
		break
	}
	return shippers[index-1].ShipperCode
}

func GetExpressInfo(logisticCode string) {
	b, err := getShipperCode(logisticCode)
	if err != nil {
		fmt.Println(colorize("异常："+err.Error(), Red))
		return
	}
	var shipperResult ShipperResult
	json.Unmarshal(b, &shipperResult)
	if shipperResult.Success {
		if len(shipperResult.Shippers) == 0 {
			fmt.Println(colorize("异常：无法识别的单号，请检查单号填写是否有误", Red))
			return
		}
	} else {
		fmt.Println(colorize("抱歉暂无查询记录...", Red))
		return
	}
	fmt.Println(shipperResult)
	var shipperCode string
	if len(shipperResult.Shippers) == 1 {
		shipperCode = shipperResult.Shippers[0].ShipperCode
	} else {
		shipperCode = selectShipper(shipperResult.Shippers)
	}

	b, err = getExpressDetail(logisticCode, shipperCode)
	if err != nil {
		fmt.Println(colorize("异常："+err.Error(), Red))
		return
	}
	var expressInfo ExpressInfo
	json.Unmarshal(b, &expressInfo)
	drawView(shipperResult.Shippers[0], expressInfo)

}
func drawView(shipperInfo ShipperInfo, expressInfo ExpressInfo) {
	state := ""
	switch expressInfo.State {
	case "2":
		state = colorize("【在途中】", Blue)
	case "3":
		state = colorize("【已签收】", Green)
	case "4":
		state = colorize("【问题件】", Red)
	default:
		state = colorize("-", Yellow)
	}
	printDelay("")
	printDelay(fmt.Sprint(colorize("-------", Yellow), colorize(shipperInfo.ShipperName, Skyblue), colorize("----------------", Yellow),
		state, colorize("----------------------", Yellow)))
	if len(expressInfo.Traces) == 0 {
		if expressInfo.Reason == "" {
			fmt.Println(colorize("此单无物流信息", Red))
		} else {
			fmt.Println(colorize(expressInfo.Reason, Red))
		}
		return
	}
	for i, t := range expressInfo.Traces {
		if len(expressInfo.Traces)-1 == i {
			drawLastItem(t, Green)
		} else {
			drawItem(t, Purple)
		}
	}
	printDelay(colorize("-------------------------------------------------------------------", Yellow))
	printDelay("")

}
func drawItem(trace Trace, color string) {
	printDelay(colorize("                         │ │", color))
	printDelay(colorize("                         .-. ", color))
	printDelay(colorize("   "+trace.AcceptTime, "skyblue") + colorize("   '-'  ", color) +
		colorize(trace.AcceptStation, "skyblue"))
}

func drawLastItem(trace Trace, color string) {
	printDelay(colorize("                         │ │", Purple))
	printDelay("\033[5;32;1m" + "                         .-. " + "\033[0m")
	printDelay(colorize("   "+trace.AcceptTime, "skyblue") + "\033[5;32;1m" + "   '-'  " + "\033[0m" +
		colorize(trace.AcceptStation, "skyblue"))
}

func printDelay(txt string) {
	time.Sleep(300 * time.Millisecond)
	fmt.Println(txt)
}

func colorize(text string, status string) string {
	var out string
	switch status {
	case Blue:
		out = "\033[34;1m" // Blue
	case Red:
		out = "\033[31;1m" // Red
	case Yellow:
		out = "\033[33;1m" // Yellow
	case Green:
		out = "\033[32;1m" // Green
	case Purple:
		out = "\033[35;1m" // Purple
	case Skyblue:
		out = "\033[36;1m" // skyGreen
	default:
		out = "\033[0m" // Default
	}
	return out + text + "\033[0m"
}

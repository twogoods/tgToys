package weather

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	. "github.com/twogoods/golib/gohttp"
)

const (

	// see http://apistore.baidu.com/apiworks/servicedetail/478.html
	APIKEY = "your apikey"

	Blue    = "blue"
	Red     = "red"
	Yellow  = "yellow"
	Green   = "green"
	Purple  = "purple"
	Skyblue = "skyblue"

	MainLine1 = "                                                         ┌─────────────┐                                                       "
	MainLine2 = "┌────────────────────────────────────────────────────────────\033[36;1mWeather\033[0m────────────────────────────┬───────────────────────────────┐"
	MainLine3 = "│                              \033[36;1mNow\033[0m                       └──────┬──────┘    %s          │          %s           │"
	MainLine4 = "├───────────────────────────────────────────────────────────────┼───────────────────────────────┼───────────────────────────────┤"

	MainLine5 = "├───────────────────────────────┬───────────────────────────────┼───────────────────────────────┼───────────────────────────────┤"

	MainLine6 = "└───────────────────────────────────────────────────────────────┴───────────────────────────────┴───────────────────────────────┘"

	MainLine7 = "┌───────────────────────────────────────────────────────────────────────────────────────────────┬───────────────────────────────┐"
	MainLine8 = "│           %s          │           %s          │           %s          │          %s           │"
	MainLine9 = "├───────────────────────────────┼───────────────────────────────┼───────────────────────────────┼───────────────────────────────┤"

	MainLine10 = "└───────────────────────────────┴───────────────────────────────┴───────────────────────────────┴───────────────────────────────┘"
)

var client *Httpclient

var iconMap map[string][]string
var code2iconMap map[string]string

func init() {

	client = HttpClientBuilder().Build()

	code2iconMap = make(map[string]string)

	code2iconMap["100"] = "iconSunny"
	code2iconMap["101"] = "iconVeryCloudy"
	code2iconMap["102"] = "iconCloudy"
	code2iconMap["103"] = "iconPartlyCloudy"
	code2iconMap["104"] = "iconCloudy"
	code2iconMap["200"] = "iconSunny"
	code2iconMap["201"] = "iconSunny"
	code2iconMap["202"] = "iconSunny"
	code2iconMap["203"] = "iconSunny"
	code2iconMap["204"] = "iconSunny"
	code2iconMap["205"] = "iconCloudy"
	code2iconMap["206"] = "iconCloudy"
	code2iconMap["207"] = "iconCloudy"
	code2iconMap["208"] = "iconCloudy"
	code2iconMap["209"] = "iconCloudy"

	code2iconMap["300"] = "iconLightShowers"
	code2iconMap["301"] = "iconHeavyShowers"
	code2iconMap["302"] = "iconThunderyShowers"
	code2iconMap["303"] = "iconThunderyHeavyRain"
	code2iconMap["304"] = "iconThunderyHeavyRain"
	code2iconMap["305"] = "iconLightRain"
	code2iconMap["306"] = "iconLightRain"
	code2iconMap["307"] = "iconHeavyRain"
	code2iconMap["308"] = "iconHeavyRain"
	code2iconMap["309"] = "iconLightRain"
	code2iconMap["310"] = "iconHeavyRain"
	code2iconMap["311"] = "iconHeavyRain"
	code2iconMap["312"] = "iconHeavyRain"
	code2iconMap["313"] = "iconLightSleet"

	code2iconMap["400"] = "iconLightSnow"
	code2iconMap["401"] = "iconLightSnow"
	code2iconMap["402"] = "iconHeavySnow"
	code2iconMap["403"] = "iconThunderySnowShowers"
	code2iconMap["404"] = "iconLightSleet"
	code2iconMap["405"] = "iconLightSleet"
	code2iconMap["406"] = "iconLightSleet"
	code2iconMap["407"] = "iconLightSnowShowers"
	code2iconMap["500"] = "iconFog"
	code2iconMap["501"] = "iconFog"
	code2iconMap["502"] = "iconFog"

	iconMap = make(map[string][]string)
	//晴
	iconSunny := []string{
		"\033[38;5;226m    \\   /    \033[0m",
		"\033[38;5;226m     .-.     \033[0m",
		"\033[38;5;226m  ― (   ) ―  \033[0m",
		"\033[38;5;226m     '-'     \033[0m",
		"\033[38;5;226m    /   \\    \033[0m"}
	//晴间多云
	iconPartlyCloudy := []string{
		"\033[38;5;226m   \\  /\033[0m      ",
		"\033[38;5;226m _ /\"\"\033[38;5;250m.-.    \033[0m",
		"\033[38;5;226m   \\_\033[38;5;250m(   ).  \033[0m",
		"\033[38;5;226m   /\033[38;5;250m(___(__) \033[0m",
		"             "}
	//少云
	iconCloudy := []string{
		"             ",
		"\033[38;5;250m     .--.    \033[0m",
		"\033[38;5;250m  .-(    ).  \033[0m",
		"\033[38;5;250m (___.__)__) \033[0m",
		"             "}

	iconMap["iconSunny"] = iconSunny
	iconMap["iconPartlyCloudy"] = iconPartlyCloudy
	iconMap["iconCloudy"] = iconCloudy
	//多云
	iconVeryCloudy := []string{
		"             ",
		"\033[38;5;240;1m     .--.    \033[0m",
		"\033[38;5;240;1m  .-(    ).  \033[0m",
		"\033[38;5;240;1m (___.__)__) \033[0m",
		"             "}

	iconMap["iconVeryCloudy"] = iconVeryCloudy
	//阵雨
	iconLightShowers := []string{
		"\033[38;5;226m _`/\"\"\033[38;5;250m.-.    \033[0m",
		"\033[38;5;226m  /\\_\033[38;5;250m(   ).  \033[0m",
		"\033[38;5;226m   /\033[38;5;250m(___(__) \033[0m",
		"\033[38;5;111m     ‘ ‘ ‘ ‘ \033[0m",
		"\033[38;5;111m    ‘ ‘ ‘ ‘  \033[0m"}

	iconMap["iconLightShowers"] = iconLightShowers
	//强阵雨
	iconHeavyShowers := []string{
		"\033[38;5;226m _`/\"\"\033[38;5;240;1m.-.    \033[0m",
		"\033[38;5;226m  /\\_\033[38;5;240;1m(   ).  \033[0m",
		"\033[38;5;226m   /\033[38;5;240;1m(___(__) \033[0m",
		"\033[38;5;21;1m   ‚‘‚‘‚‘‚‘  \033[0m",
		"\033[38;5;21;1m   ‚’‚’‚’‚’  \033[0m"}

	iconMap["iconHeavyShowers"] = iconHeavyShowers
	//雷阵雨
	iconThunderyShowers := []string{
		"\033[38;5;226m _`/\"\"\033[38;5;250m.-.    \033[0m",
		"\033[38;5;226m  /\\_\033[38;5;250m(   ).  \033[0m",
		"\033[38;5;226m   /\033[38;5;250m(___(__) \033[0m",
		"\033[38;5;228;5m    ⚡\033[38;5;111;25m‘ ‘\033[38;5;228;5m⚡\033[38;5;111;25m‘ ‘ \033[0m",
		"\033[38;5;111m    ‘ ‘ ‘ ‘  \033[0m"}

	iconMap["iconThunderyShowers"] = iconThunderyShowers
	//强雷阵雨
	iconThunderyHeavyRain := []string{
		"\033[38;5;240;1m     .-.     \033[0m",
		"\033[38;5;240;1m    (   ).   \033[0m",
		"\033[38;5;240;1m   (___(__)  \033[0m",
		"\033[38;5;21;1m  ‚‘\033[38;5;228;5m⚡\033[38;5;21;25m‘‚\033[38;5;228;5m⚡\033[38;5;21;25m‚‘   \033[0m",
		"\033[38;5;21;1m  ‚’‚’\033[38;5;228;5m⚡\033[38;5;21;25m’‚’   \033[0m"}

	iconMap["iconThunderyHeavyRain"] = iconThunderyHeavyRain

	//小雨
	iconLightRain := []string{
		"\033[38;5;250m     .-.     \033[0m",
		"\033[38;5;250m    (   ).   \033[0m",
		"\033[38;5;250m   (___(__)  \033[0m",
		"\033[38;5;111m    ‘ ‘ ‘ ‘  \033[0m",
		"\033[38;5;111m   ‘ ‘ ‘ ‘   \033[0m"}

	iconMap["iconLightRain"] = iconLightRain
	//大雨
	iconHeavyRain := []string{
		"\033[38;5;240;1m     .-.     \033[0m",
		"\033[38;5;240;1m    (   ).   \033[0m",
		"\033[38;5;240;1m   (___(__)  \033[0m",
		"\033[38;5;21;1m  ‚‘‚‘‚‘‚‘   \033[0m",
		"\033[38;5;21;1m  ‚’‚’‚’‚’   \033[0m"}

	iconMap["iconHeavyRain"] = iconHeavyRain

	//阵雪
	iconLightSnowShowers := []string{
		"\033[38;5;226m _`/\"\"\033[38;5;250m.-.    \033[0m",
		"\033[38;5;226m  /\\_\033[38;5;250m(   ).  \033[0m",
		"\033[38;5;226m   /\033[38;5;250m(___(__) \033[0m",
		"\033[38;5;255m     *  *  * \033[0m",
		"\033[38;5;255m    *  *  *  \033[0m"}

	iconMap["iconLightSnowShowers"] = iconLightSnowShowers

	//雨夹雪
	iconLightSleet := []string{
		"\033[38;5;250m     .-.     \033[0m",
		"\033[38;5;250m    (   ).   \033[0m",
		"\033[38;5;250m   (___(__)  \033[0m",
		"\033[38;5;111m    ‘ \033[38;5;255m*\033[38;5;111m ‘ \033[38;5;255m*  \033[0m",
		"\033[38;5;255m   *\033[38;5;111m ‘ \033[38;5;255m*\033[38;5;111m ‘   \033[0m"}

	iconMap["iconLightSleet"] = iconLightSleet
	//TODO
	iconLightSleetShowers := []string{
		"\033[38;5;226m _`/\"\"\033[38;5;250m.-.    \033[0m",
		"\033[38;5;226m  /\\_\033[38;5;250m(   ).  \033[0m",
		"\033[38;5;226m   /\033[38;5;250m(___(__) \033[0m",
		"\033[38;5;111m     ‘ \033[38;5;255m*\033[38;5;111m ‘ \033[38;5;255m* \033[0m",
		"\033[38;5;255m    *\033[38;5;111m ‘ \033[38;5;255m*\033[38;5;111m ‘  \033[0m"}

	iconMap["iconLightSleetShowers"] = iconLightSleetShowers

	//暴雪
	iconThunderySnowShowers := []string{
		"\033[38;5;226m _`/\"\"\033[38;5;250m.-.    \033[0m",
		"\033[38;5;226m  /\\_\033[38;5;250m(   ).  \033[0m",
		"\033[38;5;226m   /\033[38;5;250m(___(__) \033[0m",
		"\033[38;5;255m     *\033[38;5;228;5m⚡\033[38;5;255;25m *\033[38;5;228;5m⚡\033[38;5;255;25m * \033[0m",
		"\033[38;5;255m    *  *  *  \033[0m"}

	iconMap["iconThunderySnowShowers"] = iconThunderySnowShowers
	//小雪
	iconLightSnow := []string{
		"\033[38;5;250m     .-.     \033[0m",
		"\033[38;5;250m    (   ).   \033[0m",
		"\033[38;5;250m   (___(__)  \033[0m",
		"\033[38;5;255m    *  *  *  \033[0m",
		"\033[38;5;255m   *  *  *   \033[0m"}

	iconMap["iconLightSnow"] = iconLightSnow
	//大雪
	iconHeavySnow := []string{
		"\033[38;5;240;1m     .-.     \033[0m",
		"\033[38;5;240;1m    (   ).   \033[0m",
		"\033[38;5;240;1m   (___(__)  \033[0m",
		"\033[38;5;255;1m   * * * *   \033[0m",
		"\033[38;5;255;1m  * * * *    \033[0m"}

	iconMap["iconHeavySnow"] = iconHeavySnow
	//雾
	iconFog := []string{
		"             ",
		"\033[38;5;251m _ - _ - _ - \033[0m",
		"\033[38;5;251m  _ - _ - _  \033[0m",
		"\033[38;5;251m _ - _ - _ - \033[0m",
		"             "}

	iconMap["iconFog"] = iconFog

	iconUnknown := []string{
		"    .-.      ",
		"     __)     ",
		"    (        ",
		"     `-’     ",
		"      •      "}

	iconMap["iconUnknown"] = iconUnknown

}

type ApiResult struct {
	Status         string `json:"status"`
	Now            WeatherNow
	Aqi            CityApi
	Daily_forecast []WeatherForecast
}

type AirQuality struct {
	Aqi  string //空气质量指数
	Co   string //一氧化碳1小时平均值(ug/m³)
	No2  string //二氧化氮1小时平均值(ug/m³)
	O3   string //臭氧1小时平均值(ug/m³)
	Pm10 string //PM10 1小时平均值(ug/m³)
	Pm25 string //PM2.5 1小时平均值(ug/m³)
	Qlty string //空气质量类别
	So2  string //二氧化硫1小时平均值(ug/m³)
}
type CityApi struct {
	City AirQuality
}

type Wind struct {
	Deg string //风向（360度）
	Dir string //风向(北风)
	Sc  string //风力
	Spd string //风速（kmph）
}

type Cond struct {
	Code, Txt     string
	Code_d, Txt_d string //白天天气状况代码，描述
	Code_n, Txt_n string
}
type Temp struct {
	Max, Min string
}

type WeatherItem struct {
	Date string //日期
	Fl   string //体感温度
	Hum  string //相对湿度（%）
	Pcpn string //降水量（mm）
	Pres string //气压
	Tmp  string //温度
	Vis  string //能见度（km）
	Cond Cond
	Wind Wind
}

type WeatherNow struct {
	WeatherItem
	Tmp string //温度
}

type WeatherForecast struct {
	WeatherItem
	Tmp Temp //温度
}

type Suggestion struct {
}

type SuggestionItem struct {
	Brf string
	Txt string
}
type IpResult struct {
	Showapi_res_code  int
	Showapi_res_error string
	Showapi_res_body  IpBody
}
type IpBody struct {
	City string
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
func tmpColor(tmp string) string {
	t, err := strconv.Atoi(tmp)
	if err == nil {
		switch {
		case t <= 0:
			return colorize(tmp, Skyblue)
		case 0 < t && t <= 15:
			return colorize(tmp, Blue)
		case 15 < t && t <= 25:
			return colorize(tmp, Green)
		case 25 < t && t < 32:
			return colorize(tmp, Yellow)
		case t >= 32:
			return colorize(tmp, Red)
		}
	}
	return colorize(tmp, Red)
}

func weatherColor(code string, txt string) string {
	t, err := strconv.Atoi(code)
	if err == nil {
		switch {
		case t == 100 || (t >= 200 && t <= 204):
			return colorize(txt, Yellow)
		case t < 300:
			return colorize(txt, Purple)
		case 300 <= t && t < 400:
			return colorize(txt, Blue)
		case 400 <= t && t < 500:
			return colorize(txt, Skyblue)
		case t >= 500:
			return colorize(txt, Purple)
		}
	}
	return colorize(txt, Green)
}

func sendGet(city string) (ApiResult, error) {
	getreq, _ := RequestBuilder().Header("apikey", APIKEY).Url("http://apis.baidu.com/heweather/weather/free?city=" + city).Build()
	resp, err := client.Execute(getreq)
	var result ApiResult
	if err == nil {
		s, _ := resp.BodyString()
		s = s[strings.Index(s, "[")+1 : strings.LastIndex(s, "]")]
		json.Unmarshal([]byte(s), &result)
		return result, nil
	}
	return result, err
}
func getexternalIP() (string, error) {
	getreq, _ := RequestBuilder().Url("http://myexternalip.com/raw").Build()
	resp, err := client.Execute(getreq)
	if err == nil {
		str, _ := resp.BodyString()
		str = str[:len(str)-1]
		return str, nil
	}
	return "", err
}

func getCityWithIP(ip string) (string, error) {

	getreq, _ := RequestBuilder().Header("apikey", APIKEY).Url("http://apis.baidu.com/showapi_open_bus/ip/ip?ip=" + ip).Build()
	resp, err := client.Execute(getreq)
	if err == nil {
		var result IpResult
		b, _ := resp.BodyByte()
		json.Unmarshal(b, &result)
		cityname := result.Showapi_res_body.City
		return cityname, nil
	}
	return "", err
}

func printDelay(txt string) {
	time.Sleep(300 * time.Millisecond)
	fmt.Println(txt)
}

func GetWeather(city ...string) {
	var cityName string
	if len(city) == 0 {
		ip, err := getexternalIP()
		if err != nil {
			fmt.Println(colorize("获取公网ip出错---"+err.Error(), Red))
			return
		}
		cityName, err = getCityWithIP(ip)
		if err != nil || cityName == "" {
			fmt.Println(colorize("通过ip查询城市失败，请手动输入城市...", Red))
			return
		}
		fmt.Println(colorize("定位 : "+cityName+" ...", Green))
	} else {
		cityName = city[0]
	}
	r, err := sendGet(cityName)
	if err != nil {
		fmt.Println(colorize("异常：获取天气出错...", Red))
		return
	}
	if r.Status != "ok" {
		fmt.Println(colorize(r.Status, Red))
		return
	}

	printDelay(MainLine1)
	printDelay(MainLine2)
	printDelay(fmt.Sprintf(MainLine3, colorize(r.Daily_forecast[0].Date, Red), colorize(r.Daily_forecast[1].Date, Red)))
	printDelay(MainLine4)
	drawNow(r.Now, r.Daily_forecast[:2], r.Aqi)
	printDelay(MainLine5)
	printDelay(fmt.Sprintf(MainLine8, colorize(r.Daily_forecast[2].Date, Red), colorize(r.Daily_forecast[3].Date, Red),
		colorize(r.Daily_forecast[4].Date, Red), colorize(r.Daily_forecast[5].Date, Red)))
	printDelay(MainLine9)
	drawItem(r.Daily_forecast[2:6])
	printDelay(MainLine10)

}
func windDir(dir string) string {
	if dir == "无持续风向" {
		return "无风向"
	}
	return dir
}
func drawNow(now WeatherNow, items []WeatherForecast, aqi CityApi) {

	name, ok := code2iconMap[now.Cond.Code]
	if !ok {
		name = "iconUnknown"
	}
	icon := iconMap[name]
	line0 := fmt.Sprintf("│%s  %s  \t\t\t", icon[0], weatherColor(now.Cond.Code, now.Cond.Txt))
	line1 := fmt.Sprintf("│%s  %s°C \t\t\t", icon[1], tmpColor(now.Tmp))
	line2 := fmt.Sprintf("│%s  %s%s级 %s kmph \t", icon[2], now.Wind.Dir, now.Wind.Sc, now.Wind.Spd)
	line3 := fmt.Sprintf("│%s  相对湿度:%s%% \t\t", icon[3], now.Hum)
	line4 := fmt.Sprintf("│%s  气压:%s hPa \t\t", icon[4], now.Pres)

	line0 += fmt.Sprintf("空气质量指数:%s \t", "30")
	line1 += fmt.Sprintf("空气质量:%s \t\t", aqi.City.Qlty)
	line2 += fmt.Sprintf("PM2.5:%s ug/m³ \t\t", aqi.City.Pm25)
	line3 += fmt.Sprintf("降水量:%s mm \t\t", now.Pcpn)
	line4 += fmt.Sprintf("能见度:%s km \t\t", now.Vis)

	drawItem(items, line0, line1, line2, line3, line4)
}

func drawItem(items []WeatherForecast, lines ...string) {
	var line0, line1, line2, line3, line4 string
	for i, line := range lines {
		switch i {
		case 0:
			line0 = line
		case 1:
			line1 = line
		case 2:
			line2 = line
		case 3:
			line3 = line
		case 4:
			line4 = line
		}
	}
	for _, item := range items {
		name, ok := code2iconMap[item.Cond.Code_d]
		if !ok {
			name = "iconUnknown"
		}
		icon := iconMap[name]
		line0 += fmt.Sprintf("│%s  %s  \t\t", icon[0], weatherColor(item.Cond.Code_d, item.Cond.Txt_d))
		line1 += fmt.Sprintf("│%s  %s~%s°C    \t", icon[1], tmpColor(item.Tmp.Min), tmpColor(item.Tmp.Max))
		line2 += fmt.Sprintf("│%s  %s%s级 \t", icon[2], windDir(item.Wind.Dir), item.Wind.Sc)
		line3 += fmt.Sprintf("│%s  相对湿度:%s%% \t", icon[3], item.Hum)
		line4 += fmt.Sprintf("│%s  气压:%s hPa \t", icon[4], item.Pres)
	}
	printDelay(line0 + "│")
	printDelay(line1 + "│")
	printDelay(line2 + "│")
	printDelay(line3 + "│")
	printDelay(line4 + "│")
}

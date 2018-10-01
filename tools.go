package tools

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/axgle/mahonia"
)

/*全局调试控制器*/
var debugbool = true

const (
	base64Str = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
)

/**
 * http请求参数
 */
type Req struct {
	Url      string            //请求地址
	Timeout  time.Duration     //超时时间 毫秒数
	Encoding string            //请求编码
	Headers  map[string]string //请求头
	Data     io.Reader         //请求的数据
}

/**
 * 时间戳格式化
 */
type DateTime struct {
	Time   int64
	Format string
}

/**
 * [Get http GET请求]
 * @作者    como
 * @时间    2018-09-23
 * @版本    1.0.0
 * @param {[type]}    params ReqGet) (data []byte,res *http.Response, err error [description]
 */
func Get(params Req) (data []byte, res *http.Response, err error) {
	if params.Url == "" {
		return
	}
	if params.Encoding == "" {
		params.Encoding = "utf8"
	}
	if len(params.Headers) == 0 {
		headers := make(map[string]string)
		headers["User-Agent"] = "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.87 Safari/537.36"
		params.Headers = headers
	}
	if params.Timeout == 0 {
		params.Timeout = 3000
	}
	client := &http.Client{
		Timeout: params.Timeout * time.Second,
	}
	var req *http.Request
	req, err = http.NewRequest("GET", params.Url, nil)
	if err != nil {
		return
	}
	for key, value := range params.Headers {
		req.Header.Add(key, value)
	}
	res, err = client.Do(req)
	defer res.Body.Close()
	if err != nil {
		return
	}
	if res.StatusCode == 200 {
		var tmpdata []byte
		tmpdata, err = ioutil.ReadAll(res.Body)
		if err != nil {
			return
		}
		enc := mahonia.NewDecoder(params.Encoding)
		data = []byte(enc.ConvertString(string(tmpdata)))
	}
	return
}

/**
 * [Post 发送post请求]
 * @作者    como
 * @时间    2018-09-24
 * @版本    1.0.0
 * @param {[type]}    req Req) (data []byte, res *http.Response, err error [description]
 */
func Post(params Req) (data []byte, res *http.Response, err error) {
	if params.Url == "" {
		return
	}
	if params.Encoding == "" {
		params.Encoding = "utf8"
	}

	if len(params.Headers) == 0 {
		headers := make(map[string]string)
		headers["User-Agent"] = "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.87 Safari/537.36"
		params.Headers = headers
	}
	if _, ok := params.Headers["Content-Type"]; !ok {
		params.Headers["Content-Type"] = "application/x-www-form-urlencoded"
	}
	client := &http.Client{
		Timeout: params.Timeout * time.Second,
	}
	var req *http.Request
	req, err = http.NewRequest("POST", params.Url, params.Data)
	if err != nil {
		return
	}
	for key, value := range params.Headers {
		req.Header.Add(key, value)
	}
	res, err = client.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()
	if res.StatusCode == 200 {
		var tmpdata []byte
		tmpdata, err = ioutil.ReadAll(res.Body)
		if err != nil {
			return
		}
		enc := mahonia.NewDecoder(params.Encoding)
		data = []byte(enc.ConvertString(string(tmpdata)))
	}
	return
}

/**
 * [Json_encode 任意类型转成json格式]
 * @作者    como
 * @时间    2018-09-24
 * @版本    1.0.0
 * @param {[type]}    params interface{}) (data []byte, err error [description]
 */
func Json_encode(params interface{}) (data []byte, err error) {
	data, err = json.Marshal(params)
	return
}

/**
 * [Time 获取秒级单位时间]
 * @作者    como
 * @时间    2018-09-24
 * @版本    1.0.0
 * @param {[type]}    ) (t int64 [description]
 */
func Time() (t int64) {
	//设置时区
	t = time.Now().Unix()
	return
}

/**
 * [Date 时间戳转日期型]
 * @作者    como
 * @时间    2018-09-24
 * @版本    1.0.0
 * @param {[type]}    dtime DateTime) (str string [description]
 */
func Date(dtime DateTime) (str string) {
	if dtime.Time == 0 {
		dtime.Time = Time()
	}
	if dtime.Format == "" {
		dtime.Format = "Y-m-d H:i:s"
	}
	Y := "2006"
	m := "01"
	d := "02"
	H := "15"
	i := "04"
	s := "05"
	timeLayout := ""
	for _, val := range dtime.Format {
		switch string(val) {
		case "Y":
			timeLayout += Y
		case "m":
			timeLayout += m
		case "d":
			timeLayout += d
		case "H":
			timeLayout += H
		case "i":
			timeLayout += i
		case "s":
			timeLayout += s
		default:
			timeLayout += string(val)
		}
	}
	str = time.Unix(dtime.Time, 0).Format(timeLayout)
	return
}

/**
 * [Md5 md5加密码]
 * @作者    como
 * @时间    2018-09-24
 * @版本    1.0.0
 * @param {[type]}    params []byte) (str string [description]
 */
func Md5(params []byte) (str string) {
	h := md5.New()
	h.Write(params)
	cipherStr := h.Sum(nil)
	str = hex.EncodeToString(cipherStr)
	return
}

/**
 * [Base64_encode 转base64字符串]
 * @作者    como
 * @时间    2018-09-24
 * @版本    1.0.0
 * @param {[type]}    src []byte) (str string [description]
 */
func Base64_encode(src []byte) (str string) {
	coder := base64.NewEncoding(base64Str)
	str = coder.EncodeToString(src)
	return
}

/**
 * [Base64_decode base64解码]
 * @作者    como
 * @时间    2018-09-24
 * @版本    1.0.0
 * @param {[type]}    str string) (data []byte, err error [description]
 */
func Base64_decode(str string) (data []byte) {
	coder := base64.NewEncoding(base64Str)
	var err error
	data, err = coder.DecodeString(str)
	if err != nil {
		return
	}
	return
}

/**
 * [Mt_rand 生成区间随机数]
 * @作者    como
 * @时间    2018-09-24
 * @版本    1.0.0
 * @param {[type]}    min [description]
 * @param {[type]}    max int64)        (int64 [description]
 */
func Mt_rand(min, max int) (num int) {
	rand.Seed(time.Now().UnixNano())
	for {
		tmp := rand.Intn(max)
		if tmp >= min {
			num = tmp
			break
		}
	}
	return
}

/**
 * [Iconv 编码转换]
 * @作者    como
 * @时间    2018-09-26
 * @版本    1.0.0
 * @param {[type]}    tocode string  [description]
 * @param {[type]}    instr  string) (outstr       string [description]
 */
func Iconv(tocode string, instr string) (outstr string) {
	enc := mahonia.NewDecoder(tocode)
	outstr = enc.ConvertString(instr)
	return
}

/**
 * [Trim 去掉字符串左右空格]
 * @作者    como
 * @时间    2018-09-26
 * @版本    1.0.0
 * @param {[type]}    str string) (data string [description]
 */
func Trim(str string) (data string) {
	tmp := bytes.Trim([]byte(str), " ")
	data = string(tmp)
	return
}

/**
 * [File_get_contents 实现了php的file_put_contens]
 * @作者    como
 * @时间    2018-09-26
 * @版本    1.0.0
 * @param {[type]}    path string) (str string, err error [description]
 */
func File_get_contents(path string) (str string, err error) {
	var exists bool
	exists, err = File_exists(path)
	if err != nil {
		return
	}
	if exists {
		str, err = Fread(path)
		return
	}
	req := Req{
		Url:      path,
		Encoding: "utf8",
	}
	data, _, err1 := Get(req)
	if err1 != nil {
		err = err1
		return
	}
	str = string(data)
	return
}

/**
 * [File_exists 判断文件是否存在]
 * @作者    como
 * @时间    2018-10-01
 * @版本    1.0.0
 * @param {[type]}    path string) (bool, error [description]
 */
func File_exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

/**
 * [File_put_contens 实现了php的file_put_contens]
 * @作者    como
 * @时间    2018-09-26
 * @版本    1.0.0
 * @param {[type]}    path [description]
 * @param {[type]}    data string)       (err error [description]
 */
func File_put_contens(path, data string) (err error) {
	var f *os.File
	f, err = os.Create(path)
	if err != nil {
		return
	}
	defer f.Close()
	_, err = f.Write([]byte(data))
	return
}

/**
 * [Fread 文件读路函数]
 * @作者    como
 * @时间    2018-09-26
 * @版本    1.0.0
 * @param {[type]}    path string) (str string, err error [description]
 */
func Fread(path string) (str string, err error) {
	var f *os.File
	f, err = os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()
	var data []byte
	data, err = ioutil.ReadAll(f)
	if err != nil {
		return
	}
	str = string(data)
	return
}

/**
 * [Explode 字符串拆分成数组]
 * @Author    como
 * @DateTime  2018-09-16
 * @version   [1.0.1]
 * @param     {[type]}   spt string  [description]
 * @param     {[type]}   str string) (data         []string [description]
 */
func Explode(spt string, str string) (data []string) {
	data = strings.Split(str, spt)
	return
}

/**
 * [Implode 字符串数组合并成字符串]
 * @Author    como
 * @DateTime  2018-09-16
 * @version   [1.0.1]
 * @param     {[type]}   spt  string    [description]
 * @param     {[type]}   data []string) (str          string [description]
 */
func Implode(spt string, data []string) (str string) {
	str = strings.Join(data, spt)
	return
}

/**
 * [Strtotime 日期时间转时间戳]
 * @Author    como
 * @DateTime  2018-09-27
 * @version   [1.0.1]
 * @param     {[type]}   str string) (data int64 [description]
 */
func Strtotime(str string) (data int64) {
	timeLayout := "2006-01-02 15:04:05"
	loc, _ := time.LoadLocation("Local")
	//使用模板在对应时区转化为time.time类型
	theTime, _ := time.ParseInLocation(timeLayout, str, loc)
	data = theTime.Unix()
	return
}

/**
 * [debug 全局调试函数]
 * @作者     como
 * @时间     2018-09-24
 * @版本     1.0.0
 * @param  {[type]}    v ...interface{} [description]
 * @return {[type]}      [description]
 */
func debug(v ...interface{}) {
	if debugbool {
		log.Println(v)
	}
}

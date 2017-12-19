/*
	This package is used to translate texts with youdao's api.
	Response:
		`{
			"web":[
				{
					"value":["你好世界","举个例子","开始"],
					"key":"hello world"
				},
				{
					"value":["会写的人多了去了"],
					"key":"Hello   World"
				},
				{
					"value":["凯蒂猫气球世界"],
					"key":"Hello Kitty World"
				}
			],
			"query":"hello world",
			"translation":["hello world"],
			"errorCode":"0",
			"dict":{"url":"yddict://m.youdao.com/dict?le=eng&q=hello+world"},
			"webdict":{"url":"http://m.youdao.com/dict?le=eng&q=hello+world"},
			"basic":{
				"speech":"hello world",
				"uk-speech":"hello world",
				"us-speech":"hello world",
				"explains":["你好世界"]
			},
			"l":"en2zh-CHS"
		}`
	If you want to use this package you must sign up a Acount for YouDao Dic API,
	there have a test key:
	AppKey: "2ccac2276928012f",
	SecKey: "tm0aC9BVe2qZq4DuHRR9p5KdEA7y6l1Y"

*/

package youdao

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
)

var (
	AppKey = ""
	Seckey = ""
	Salt  = "53"
)

const (
	URL         = "http://openapi.youdao.com/api"
	LENGTHLIMIT = 5000
)

var (
	Trace   *log.Logger
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
)

func init() {
	file, err := os.OpenFile("log_Youdao.txt",
		os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to open error log file:", err)
	}
	// Discard 成功调用，但是什么事都不做
	Trace = log.New(ioutil.Discard,
		"TRACE: ",
		log.Ldate|log.Ltime|log.Lshortfile)
	Info = log.New(file,
		"INFO: ",
		log.Ldate|log.Ltime|log.Lshortfile)
	Warning = log.New(os.Stdout,
		"WARNING: ",
		log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(io.MultiWriter(file, os.Stderr),
		"ERROR: ",
		log.Ldate|log.Ltime|log.Lshortfile)
}

type Language int

const (
	Chinese Language = iota
	Japanese
	English
	Korean
	French
	Russian
	Portuguese
	Spanish
	Others
)

func (l Language) String() string {
	switch l {
	case Chinese:
		return "zh-CHS"
	case English:
		return "EN"
	case Japanese:
		return "ja"
	case French:
		return "fr"
	case Korean:
		return "ko"
	case Portuguese:
		return "pt"
	case Spanish:
		return "es"
	case Russian:
		return "ru"
	}
	return "auto"
}

type (
	youdaoDicReq struct {
		Query  string
		From   Language
		To     Language
		AppKey string
		Salt   string
		Sign   string
	}

	youdaoDicResp struct {
		ErrorCode   string   `json:"errorCode"`
		Query       string   `json:"query"`
		Translation []string `json:"translation"`
		Basic       struct {
			Explains []string `json:"explains"`
		} `json:"basic"`
		Web []struct {
			Key   string `json:"key"`
			Value string `json:"value"`
		} `json:"web"`
		Dict struct {
			Url string `json:"url"`
		} `json:"dict"`
		Webdict struct {
			Url string `json:"url"`
		} `json:"webdict"`
		Lan string `json:"l"`
	}
)

func New(query string, from Language, to Language) *youdaoDicReq {
	if AppKey == "" || Seckey == "" {
		Error.Println("Please run Youdao.Config to config appkey and seckey")
		panic("Please run Youdao.Config to config appkey and seckey")
	}
	ymd5 := fmt.Sprintf("%x", md5.Sum([]byte(AppKey+query+Salt+Seckey)))
	return &youdaoDicReq{
		Query:  query,
		From:   from,
		To:     to,
		AppKey: AppKey,
		Salt:   Salt,
		Sign:   ymd5,
	}
}

func (reqData *youdaoDicReq) Request() (string, error) {
	result := youdaoDicResp{}
	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {

	}
	// create url parameters
	var urlPara = req.URL.Query()
	urlPara.Add("q", reqData.Query)
	urlPara.Add("from", reqData.From.String())
	urlPara.Add("to", reqData.To.String())
	urlPara.Add("salt", reqData.Salt)
	urlPara.Add("sign", reqData.Sign)
	urlPara.Add("appKey", reqData.AppKey)
	// parameters url encode
	req.URL.RawQuery = urlPara.Encode()
	Info.Println("URL REQ: ", req.URL.Host, req.URL.Path, req.URL.RawQuery)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		Error.Println("Req fialed", err)
		return "", error(err)
	}
	defer resp.Body.Close()
	// decode jsreon to struct
	err2 := json.NewDecoder(resp.Body).Decode(&result)
	if err2 != nil {
		Error.Println("Decode to json failed", err2)
		return "", err2
	}
	Info.Println("RESP: ", result)
	return strings.Join(result.Translation, "."), result.checkErrorCode()
}

func (reqData *youdaoDicResp) checkErrorCode()  error{
	switch reqData.ErrorCode {
	case "0":
		return nil
	case "103":
		Error.Println("[103] Text is too long")
		return fmt.Errorf("Text is too long")
	case "108":
		Error.Println("[108] Unvaliable AppKey")
		return fmt.Errorf("Unvaliable AppKey")
	case "111":
		Error.Println("[111] Account fail")
		return fmt.Errorf("Account fail")
	case "202":
		Error.Println("[202] MD5 check error")
		return fmt.Errorf("MD5 check error")
	default:
		Error.Println("[%s] Unknow Error", reqData.ErrorCode)
		return fmt.Errorf("Unknow Error")

	}
	return nil
}

var tailSign = []rune{'.', '?', '!', '。', '？', '！'}

func SetTail(tail []rune) {
	tailSign = tail
}

func isTail(r rune) bool {
	for _, i := range tailSign {
		if r == i {
			return true
		}
	}
	return false
}

type info struct {
	index int
	str   string
}

func transLimLen(text []rune, from Language, to Language, limitLen int) string {
	textNum := len(text)/limitLen + 1
	results := make(chan *info, textNum)
	wg := sync.WaitGroup{}
	wg.Add(textNum)
	for i := 0; i < textNum; i++ {
		go func(index int) {
			tail := (index + 1) * limitLen
			if tail > len(text) {
				tail = len(text)
			}
			str, err := New(string(text[index*limitLen:tail]), from, to).Request()
			if err != nil {
			}
			results <- &info{index: index, str: str}
			wg.Done()
		}(i)
	}
	wg.Wait()
	close(results)
	resslice := make([]string, textNum)
	for i := range results {
		resslice[i.index] = i.str
	}
	return strings.Join(resslice, " ")
}

func TranslateTexts(texts []string, from Language, to Language) (result string) {
	resultschan := make(chan *info, len(texts))
	wgroot := sync.WaitGroup{}
	wgroot.Add(2)

	wg := sync.WaitGroup{}
	wg.Add(len(texts))

	go func() {
		for i, t := range texts {
			go func(index int, str string) {
				strs := transLimLen([]rune(str), from, to, LENGTHLIMIT)
				resultschan <- &info{index: index, str: strs}
				wg.Done()
			}(i, t)
		}
		wgroot.Done()
	}()

	go func() {
		resultslice := make([]string, len(texts))
		for i := range resultschan {
			resultslice[i.index] = i.str
		}
		result = strings.Join(resultslice, ".")
		wgroot.Done()
	}()
	wg.Wait()
	close(resultschan)
	wgroot.Wait()
	return
}

func Translate(text string, from Language, to Language, hastexts bool) (result string) {
	if hastexts {
		texts := strings.FieldsFunc(strings.TrimSpace(text), isTail)
		return TranslateTexts(texts, from, to)
	}
	return transLimLen([]rune(text), from, to, LENGTHLIMIT)
}

func Config(appkey string, secKey string) {
	AppKey = appkey
	Seckey = secKey
}

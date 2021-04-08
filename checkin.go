package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

var (
	_url            = "https://zhiyou.smzdm.com/user/checkin/jsonp_checkin"
	chanify_url		= "https://api.chanify.net/v1/sender/"
	smzdm_cookie    = ""
	chanify_token	= ""
	default_headers = map[string]string{
		"Accept":     "*/*",
		"Host":       "zhiyou.smzdm.com",
		"Referer":    "https://www.smzdm.com/",
		"User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.4324.192 Safari/537.36",
	}
)

type checkinType struct {
	ErrorCode int    `json:"error_code,omitempty"`
	ErrorMsg  string `json:"error_msg,omitempty"`
	Data      struct {
		AddPoint                  int    `json:"add_point,omitempty"`
		CheckinNum                string `json:"checkin_num,omitempty"`
		Point                     int    `json:"point,omitempty"`
		Exp                       int    `json:"exp,omitempty"`
		Gold                      int    `json:"gold,omitempty"`
		Prestige                  int    `json:"prestige,omitempty"`
		Rank                      int    `json:"rank,omitempty"`
		Slogan                    string `json:"slogan,omitempty"`
		Cards                     string `json:"cards,omitempty"`
		CanContract               int    `json:"can_contract,omitempty"`
		ContinueCheckinDays       int    `json:"continue_checkin_days,omitempty"`
		ContinueCheckinRewardShow bool   `json:"continue_checkin_reward_show,omitempty"`
	} `json:"data,omitempty"`
}

func initCheck() {
	if os.Getenv("SMZDM_COOKIE") == "" {
		panic("SMZDM_COOKIE 为空")
	}
	if os.Getenv("CHANIFY_TOKEN") == "" {
		fmt.Println("CHANIFY_TOKEN 未设置，无失败通知")
	} else {
		chanify_token = os.Getenv("CHANIFY_TOKEN")
	}
}
func main() {
	initCheck()
	req, _ := http.NewRequest("GET", _url, nil)
	for k, v := range default_headers {
		req.Header.Set(k, v)
	}
	req.Header.Set("Cookie", os.Getenv("SMZDM_COOKIE"))
	// req.Header.Set("Cookie", test_cookie)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("clirnt.Do()", err)
		return
	}
	// byteresp, _ := ioutil.ReadAll(resp.Body)
	var ct checkinType
	var errs error
	if errs = json.NewDecoder(resp.Body).Decode(&ct); errs != nil {
		switch et := err.(type) {
		case *json.UnmarshalTypeError:
			log.Printf("UnmarshalTypeError: Value[%s] Type[%v]\n", et.Value, et.Type)
		case *json.InvalidUnmarshalError:
			log.Printf("InvalidUnmarshalError: Type[%v]\n", et.Type)
		default:
			log.Println(errs)
		}
	}

	switch ct.ErrorCode {
	case 0:
		log.Println("张大妈签到完毕!", ct.ErrorCode, ct.Data.Slogan)
	default:
		s := fmt.Sprintf("张大妈签到失败 %s ErrCode:%d,ErrMsg:%s", time.Now().Format("2006-01-02"), ct.ErrorCode, ct.ErrorMsg)
		log.Println(s)
		Send(s)
	}

}

func Send(msg string) {
	if len(chanify_token) < 5 {
		log.Println("未设置chanify_token，不发送通知")
		return
	}
	v := url.Values{}
	v.Add("text", msg)
	req, _ := http.NewRequest(http.MethodPost, chanify_url+chanify_token, strings.NewReader(v.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	_, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
	}

}

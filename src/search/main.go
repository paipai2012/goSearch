/*****************************************************************************
 *
 *  author : 彭东江
 *  email  : pengdongjiang@gmail.com
 *  description : 应用入口 对外提供http服务
 *
******************************************************************************/
package main

import (
	"config"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"service"
	"strconv"
)

const (
	URL_SEARCH  uint64 = 1
	URL_UPDATE  uint64 = 2
	URL_CREATE  uint64 = 3
	URL_CONTROL uint64 = 4
	URL_SHOW    uint64 = 5
)

//定义http服务框架
type HttpService struct {
}

//实现http服务入口
func (this *HttpService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	parms, err := this.parseArgs(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, MakeErrorResult(-1, err.Error()))
		return
	}
	requestUrl := r.RequestURI
	_, reqType, err := this.ParseURL(requestUrl)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, MakeErrorResult(-1, err.Error()))
		goto END
	}
	switch reqType {
	case URL_SEARCH:
		fmt.Println(r.Method)
		fmt.Println(parms)
		res, err := service.Search("yyyy")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			io.WriteString(w, MakeErrorResult(-2, err.Error()))
			goto END
		}
		fmt.Println(res)
		io.WriteString(w, res)
	}
END:
	return
}

//解析http服务url路径参数
func (this *HttpService) ParseURL(url string) (int, uint64, error) {
	urlPattern := "/v(\\d)/(_search|_update|_contrl|_create|_show|_debug|_status|_load)\\?"
	urlRegexp, err := regexp.Compile(urlPattern)
	if err != nil {
		return -1, 0, err
	}
	matchs := urlRegexp.FindStringSubmatch(url)
	if matchs == nil {
		return -1, 0, errors.New("URL ERROR")
	}
	versionNum, _ := strconv.ParseInt(matchs[1], 10, 8)
	version := int(versionNum)
	resource := matchs[2]
	if resource == "_search" {
		return version, URL_SEARCH, nil
	}
	if resource == "_update" {
		return version, URL_UPDATE, nil
	}
	if resource == "_create" {
		return version, URL_CREATE, nil
	}
	if resource == "_show" {
		return version, URL_SHOW, nil
	}
	if resource == "_contrl" {
		return version, URL_CONTROL, nil
	}
	return -1, 0, errors.New("Error")
}

func main() {
	config.LoadConfig()
	// http.HandleFunc("/", hello)
	httpService := &HttpService{}
	http.ListenAndServe(":8081", httpService)
}

func hello(rw http.ResponseWriter, req *http.Request) {
	io.WriteString(rw, "hello world")
}

//错误信息输出类
func MakeErrorResult(errcode int, errmsg string) string {
	data := map[string]interface{}{
		"error_code": errcode,
		"message":    errmsg,
	}
	result, err := json.Marshal(data)
	if err != nil {
		return fmt.Sprintf("{\"error_code\":%v,\"message\":\"%v\"}", errcode, errmsg)
	}
	return string(result)
}

func (this *HttpService) parseArgs(r *http.Request) (map[string]string, error) {
	err := r.ParseForm()
	if err != nil {
		return nil, err
	}

	res := make(map[string]string)
	for k, v := range r.Form {
		res[k] = v[0]
	}
	return res, nil
}

package main

import (
	"config"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
)

const (
	URL_SEARCH  uint64 = 1
	URL_UPDATE  uint64 = 2
	URL_CREATE  uint64 = 3
	URL_CONTROL uint64 = 4
	URL_SHOW    uint64 = 5
)

type HttpService struct {
}

func (this *HttpService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	requestUrl := r.RequestURI
	_, reqType, err := this.ParseURL(requestUrl)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, MakeErrorResult(-1, err.Error()))
		goto END
	}
	fmt.Printf("%v|", err)
	fmt.Printf("%v|", reqType)
END:
	return
}

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

func MakeErrorResult(errcode int, errmsg string) string {
	data := [string]interface{}{
		"error_code": errcode,
		"message":    errmsg,
	}
	result, err := json.Marshal(data)
	if err != nil {
		return fmt.Sprintf("{\"error_code\":%v,\"message\":\"%v\"}", errcode, errmsg)
	}
	return string(result)
}

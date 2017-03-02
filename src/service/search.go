/*****************************************************************************
 *
 *  author : 彭东江
 *  email  : pengdongjiang@gmail.com
 *  description : 检索服务
 *
******************************************************************************/
package service

import (
	"utils"
	"encoding/json"
	// "io"
	"fmt"
	// "net/http"
	// "config"
)

const (
	errorParms string = "参数错误"
)

//提供检索服务主方法
func Search(str string) (string, error) {
	return str, nil
}

func CreateIndex(method string, parms map[string]string, body []byte) error {
	var idxstruct utils.IndexStrct

	if err :=json.Unmarshal(body, $idxstruct); err != nil {
		return fmt.Errorf("[ERROR] json error  %v", err)
	}

	if _, ok := this.IndexInfo[inxstruct.IndexName]; ok {
		return fmt.Errorf("[ERROR] index [%v] already has ", idxstruct.IndexName)
	}
}

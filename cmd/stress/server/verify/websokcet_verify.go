// Package verify 校验
package verify

import (
	"encoding/json"
	"fmt"

	"github.com/BeCrafter/commander/cmd/stress/model"
	"github.com/fatih/color"
)

// WebSocketResponseJSON 返回数据结构体，返回值为json
type WebSocketResponseJSON struct {
	Seq      string `json:"seq"`
	Cmd      string `json:"cmd"`
	Response struct {
		Code    int         `json:"code"`
		CodeMsg string      `json:"codeMsg"`
		Data    interface{} `json:"data"`
	} `json:"response"`
}

// WebSocketJSON 通过返回的Body 判断
// 返回示例: {"seq":"1566276523281-585638","cmd":"heartbeat","response":{"code":200,"codeMsg":"Success","data":null}}
// code 取body中的返回code
func WebSocketJSON(request *model.Request, seq string, msg []byte) (code int, isSucceed bool) {
	responseJSON := &WebSocketResponseJSON{}
	err := json.Unmarshal(msg, responseJSON)
	if err != nil {
		code = model.ParseError
		color.New(color.FgRed).Printf("请求结果 json.Unmarshal msg:%s err:%v", string(msg), err)
	} else {

		if seq != responseJSON.Seq {
			code = model.ParseError
			color.New(color.FgRed).Println("请求和返回seq不一致 ~请求:", seq, responseJSON.Seq, string(msg))
		} else {
			code = responseJSON.Response.Code
			// body 中code返回200为返回数据成功
			if code == 200 {
				isSucceed = true
			} else {
				resp, _ := json.Marshal(responseJSON)
				color.New(color.FgRed).Printf("请求结果 %s 的 code != 200\n", string(resp))
			}
		}
	}
	// 开启调试模式
	if request.GetDebug() {
		seqstr := color.New(color.FgHiGreen).Sprintf("Seq: %s", seq)
		msgstr := color.New(color.FgHiCyan).Sprintf("Msg: %s", string(msg))
		fmt.Printf("\n请求结果: \n %v\n %v \n\n", seqstr, msgstr)
	}
	return
}

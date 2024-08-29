package helper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/google/go-cmp/cmp"
	"github.com/yudai/gojsondiff"
	"github.com/yudai/gojsondiff/formatter"
)

type Request struct {
	Method string
	Host   []string
	Url    string
	Header []string
	Data   []string
	Retry  int
	Debug  bool
	Quiet  bool
	Sort   bool
}

func (r *Request) check() error {
	if len(r.Host) < 2 {
		return fmt.Errorf("host is not enough")
	}
	if r.Method != "GET" && r.Method != "POST" {
		return fmt.Errorf("method must be GET or POST")
	}
	return nil
}

func (r *Request) dataConvert() string {
	if len(r.Data) == 0 {
		return ""
	}

	var tmp []map[string]interface{}
	if err := json.Unmarshal([]byte(fmt.Sprintf("[%v]", strings.Join(r.Data, ","))), &tmp); err == nil {
		v, _ := json.Marshal(tmp[0])
		return string(v)
	}

	data := make(url.Values)
	for _, v := range r.Data {
		kv := strings.SplitN(v, "=", 2)
		if len(kv) == 1 {
			data[kv[0]] = []string{""}
		} else if len(kv) > 1 {
			data[kv[0]] = []string{kv[1]}
		}
	}
	return data.Encode()
}

func (r *Request) Run() ([]byte, []byte) {
	if err := r.check(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	data := r.dataConvert()

	url1 := r.Host[0] + r.Url
	url2 := r.Host[1] + r.Url

	str1, err1 := r.sendRequestWithRetry(url1, data)
	if err1 != nil {
		color.New(color.FgHiRed).Printf("Error fetching result from %s: %v\n", url1, err1)
		os.Exit(2)
	}

	str2, err2 := r.sendRequestWithRetry(url2, data)
	if err2 != nil {
		color.New(color.FgHiRed).Printf("Error fetching result from %s: %v\n", url2, err2)
		os.Exit(2)
	}

	if r.Debug {
		fmt.Printf("\n\n# ======================== Debug Ouput ========================== #\n\n")
		fmt.Printf("First  request response: %s\n", string(str1))
		fmt.Printf("Second request response: %s\n", string(str2))
		fmt.Printf("\n# =============================================================== #\n\n\n")
	}

	return str1, str2
}

// sendRequestWithRetry 发送HTTP请求并重试
func (r *Request) sendRequestWithRetry(url string, data string) (ret []byte, err error) {
	for i := 0; i < r.Retry; i++ {
		ret, err = r.sendRequest(url, data)
		if err == nil {
			break
		}
		time.Sleep(time.Millisecond * 100)
	}
	return ret, err
}

// SendRequest 发送HTTP请求
func (r *Request) sendRequest(url string, data string) ([]byte, error) {
	var req *http.Request
	var err error

	if r.Debug {
		color.New(color.FgYellow).Printf("Sending request to: %s params: %s\n", url, data)
	}

	switch r.Method {
	case "POST":
		req, err = http.NewRequest(r.Method, url, bytes.NewBufferString(data))
		if err != nil {
			return nil, err
		}
	default:
		if len(data) > 0 {
			if strings.Contains(url, "?") {
				url += "&" + data
			} else {
				url += "?" + data
			}
		}
		req, err = http.NewRequest(r.Method, url, nil)
		if err != nil {
			return nil, err
		}
	}

	for _, v := range r.Header {
		kv := strings.SplitN(v, ":", 2)
		if len(kv) == 1 {
			req.Header.Add(strings.TrimSpace(kv[0]), "")
		} else if len(kv) > 1 {
			req.Header.Add(strings.TrimSpace(kv[0]), strings.TrimSpace(kv[1]))
		}
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if r.Sort {
		body = SortMapValuesByBytes(body)
	}

	if !r.Debug {
		var res map[string]interface{}
		json.Unmarshal(body, &res)
		if _, has := res["traceId"]; has {
			delete(res, "traceId")
			body, _ = json.Marshal(res)
		}
	}

	return body, nil
}

func (r *Request) JsonDiff(data1, data2 []byte, format string) string {

	differ := gojsondiff.New()
	d, err := differ.Compare(data1, data2)
	if err != nil {
		color.New(color.FgHiRed).Printf("Failed to unmarshal file: %s\n", err.Error())
		os.Exit(3)
	}

	if d.Modified() || !r.Quiet {
		if !r.Quiet {
			color.New(color.FgHiCyan).Printf("\nThe JSON objects result:\n\n")
		} else {
			color.New(color.FgHiRed).Printf("\nThe JSON objects are different:\n\n")
		}
	} else {
		color.New(color.FgHiGreen).Printf("\nThe JSON objects are the same.\n\n")
		os.Exit(0)
	}

	// Output the result
	var diffString string
	if format == "delta" {
		formatter := formatter.NewDeltaFormatter()
		diffString, _ = formatter.Format(d)
	} else {
		var aJson map[string]interface{}
		json.Unmarshal(data1, &aJson)

		config := formatter.AsciiFormatterConfig{
			ShowArrayIndex: true,
			Coloring:       true,
		}

		formatter := formatter.NewAsciiFormatter(aJson, config)
		diffString, _ = formatter.Format(d)
	}

	return diffString
}

func (r *Request) CmpDiff(data1, data2 []byte) string {
	var obj1 map[string]interface{}
	var obj2 map[string]interface{}

	if err1 := json.Unmarshal(data1, &obj1); err1 != nil {
		color.New(color.FgHiRed).Printf("Failed to unmarshal file1: %s\n", err1.Error())
		os.Exit(3)
	}
	if err2 := json.Unmarshal(data2, &obj2); err2 != nil {
		color.New(color.FgHiRed).Printf("Failed to unmarshal file2: %s\n", err2.Error())
		os.Exit(3)
	}

	if len(obj1) == 0 || len(obj2) == 0 {
		color.New(color.FgHiRed).Printf("\nresult is empty.\n\n")
		os.Exit(1)
	}

	res := cmp.Diff(obj1, obj2)
	if len(res) == 0 {
		color.New(color.FgHiGreen).Printf("\nThe JSON objects are the same.\n\n")
		os.Exit(0)
	}

	return res
}

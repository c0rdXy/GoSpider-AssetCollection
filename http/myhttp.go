package http

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

var (
	// 设置默认请求头
	defaultHeaders = map[string]string{
		"User-Agent":      "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:121.0) Gecko/20100101 Firefox/121.0",
		"Accept":          "application/json, text/plain, */*",
		"Accept-Language": "zh-CN,zh;q=0.8,zh-TW;q=0.7,zh-HK;q=0.5,en-US;q=0.3,en;q=0.2",
		"Accept-Encoding": "gzip, deflate",
		"Content-Type":    "application/json",
		"Version":         "TYC-Web",
		"X-Tycid":         "b2329190b5e911eebc20b30d1fc4016a",
		"Origin":          "https://www.tianyancha.com",
		"Referer":         "https://www.tianyancha.com/",
		"Sec-Fetch-Dest":  "empty",
		"Sec-Fetch-Mode":  "cors",
		"Sec-Fetch-Site":  "same-site",
		"Te":              "trailers",
		"Connection":      "close",
	}
)

// MyHTTPPost 发送 HTTP POST 请求
func MyHTTPPost(url string, payload interface{}) ([]byte, error) {
	// 将 payload 转换为 JSON 格式
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	// 创建 HTTP 请求
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	// 设置请求头
	setRequestHeaders(req)

	// 发送请求
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	//body, err := io.ReadAll(resp.Body)
	//fmt.Println(resp.Body)
	//if err != nil {
	//	return nil, err
	//}

	body, err := checkGzip(resp)
	if err != nil {
		return nil, err
	}

	return body, nil
}

// checkGzip 检查是否有 gzip 压缩
func checkGzip(res *http.Response) ([]byte, error) {
	// 是否有 gzip
	gzipFlag := false
	for k, v := range res.Header {
		if strings.ToLower(k) == "content-encoding" && strings.ToLower(v[0]) == "gzip" {
			gzipFlag = true
		}
	}

	var content []byte
	if gzipFlag {

		// 创建 gzip.Reader
		gr, err := gzip.NewReader(res.Body)
		defer gr.Close()
		if err != nil {
			fmt.Println(err.Error())
		}
		content, _ = io.ReadAll(gr)
		return content, err
	} else {
		content, _ = io.ReadAll(res.Body)
		return content, nil
	}
	return nil, nil
}

// setRequestHeaders 设置请求头
func setRequestHeaders(req *http.Request) {
	header := req.Header
	for key, value := range defaultHeaders {
		header.Set(key, value)
	}
}

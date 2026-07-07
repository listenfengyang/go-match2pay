package utils

import (
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
)

// RestyRequest HTTP 请求日志
type RestyRequest struct {
	Method  string      `json:"method"`
	URL     string      `json:"url"`
	Headers http.Header `json:"headers"`
	Body    interface{} `json:"body"`
}

// RestyResponse HTTP 响应日志
type RestyResponse struct {
	StatusCode int         `json:"status_code"`
	Status     string      `json:"status"`
	Headers    http.Header `json:"headers"`
	Body       string      `json:"body"`
	ReceivedAt time.Time   `json:"received_at"`
	Rtt        int64       `json:"rtt"`
}

// RestyLog HTTP 请求/响应日志
type RestyLog struct {
	Request  RestyRequest  `json:"request"`
	Response RestyResponse `json:"response"`
}

// GetRestyLog 从 resty.Response 中提取日志信息
func GetRestyLog(resp *resty.Response) RestyLog {
	reqHeaders := map[string][]string(resp.Request.Header)
	delete(reqHeaders, "User-Agent")

	return RestyLog{
		Request: RestyRequest{
			Method:  resp.Request.Method,
			URL:     resp.Request.URL,
			Headers: reqHeaders,
			Body:    resp.Request.Body,
		},
		Response: RestyResponse{
			StatusCode: resp.StatusCode(),
			Status:     resp.Status(),
			Headers:    resp.Header(),
			Body:       resp.String(),
			ReceivedAt: resp.ReceivedAt(),
			Rtt:        resp.Time().Milliseconds(),
		},
	}
}

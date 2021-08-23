package httptool

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	qlog "vault/log"
)

// 嵌入ResponseWriter
type ResponseWithRecorder struct {
	 http.ResponseWriter
	 statusCode int
	 body bytes.Buffer
}

// 这里一定要用指针方法，否则对自身结构修改不生效
func (r *ResponseWithRecorder)Write(data []byte) (int, error) {
	r.body.Write(data)
	return r.ResponseWriter.Write(data)
}

func (r *ResponseWithRecorder)WriteHeader(statusCode int) {
	r.statusCode = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}

type Mid_logger struct {}

func (m *Mid_logger)mid_handle(next http.Handler) http.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {

		// 构建自定义responseWriter
		wr := &ResponseWithRecorder{
			w,
			http.StatusOK,
			bytes.Buffer{},
		}

		// 调用子模块handler 计算响应时长
		timeStart := time.Now()
		next.ServeHTTP(wr, r)
		timeElapsed := time.Since(timeStart)

		// 处理请求&响应
		requestHeader, err := json.Marshal(r.Header)
		if err != nil {
			qlog.HttpAccess.Println("read request header failed:", err)
		}
		requestBodyBuf := make([]byte, r.ContentLength)
		r.Body.Read(requestBodyBuf)

		// 日志写入
		accessLogString := fmt.Sprintf("%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v",
			r.RemoteAddr,
			r.Method,
			r.URL.Path,
			string(requestHeader),
			string(requestBodyBuf),
			wr.statusCode,
			timeElapsed,
			wr.Header(),
			wr.body.String())
		qlog.HttpAccess.Println(accessLogString)

		fmt.Println("Logger Middle Ware Work Done.")
	}
	return http.HandlerFunc(handler)
}
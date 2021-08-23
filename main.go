package main

/*
上传文件：curl -v -F "file=@/Users/workspace_sun/Desktop/Document/2channel.mp4" "http://127.0.0.1:8081/2channel.mp4"
*/

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
	"httptool"
	qlog "log"
)

const (
	//rootPath = "/root/vault"
	rootPath = "/Users/workspace_sun/Desktop"
)

func main() {

	// 服务中间件
	router := httptool.NewRouter()
	router.Use(&httptool.Mid_timer{})
	router.Use(&httptool.Mid_logger{})
	router.Add("/", http.HandlerFunc(objectHandler))

	// 监听http服务
	mux := http.NewServeMux()

	// router绑定到mux
	router.BindMux(mux)

	server := &http.Server{
		Addr:         "127.0.0.1:8081",
		Handler:      mux,
		ReadTimeout:  10,
		WriteTimeout: 10,
		ErrorLog:     qlog.HttpError,
	}
	server.ListenAndServe()
}

 func objectHandler(w http.ResponseWriter, r *http.Request) {
 	m := r.Method
 	if m == http.MethodPost {
		post(w, r)
		return
	}
	if m == http.MethodGet {
		get(w, r)
		return
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
 }

 func post(w http.ResponseWriter, r *http.Request) {

 	// 匹配"/"
 	if r.URL.EscapedPath() == "/" {
		fmt.Fprint(w, "Show me your money.")
		return
	}

	// 读取表单文件
	r.ParseMultipartForm(1024*1024*1024)
 	file, _, err := r.FormFile("file")
	 if err != nil {
		 log.Println("parse file err", err)
		 w.WriteHeader(http.StatusBadGateway)
		 return
	 }
	 defer file.Close()

 	// 拼接file path
	filePath := rootPath + r.URL.EscapedPath()

	// path解析
	pathItems := strings.Split(filePath, "/")
	fileName := "/" + pathItems[len(pathItems) - 1]
	dirPath := strings.TrimRight(filePath, fileName)
	fmt.Println("target path:", filePath)
	fmt.Println("file name:" , fileName)
	fmt.Println("dir path:" , dirPath)

	// "/"结尾抛error
	 if fileName == "" {
		 w.Header().Set("Content-Type", "application/json")
		 w.WriteHeader(http.StatusBadRequest)
		 w.Write([]byte("path cannot be ended with '/'"))
		 return
	 }

	// 判断文件夹路径是否存在
	if isExist(dirPath) == false {
		// 递归创建文件夹
		err := os.MkdirAll(dirPath, os.ModePerm)
		if err != nil {
			log.Println("create dir failed.", err)
			w.WriteHeader(http.StatusBadGateway)
			return
		}
	}

	// 判断文件路径是否存在  如果存在先删除
	 if isExist(filePath) {
		 err := os.Remove(filePath)
		 if err != nil {
		 	log.Println("remove old file failed.", err)
			w.WriteHeader(http.StatusBadGateway)
			 return
		 }
	 }

	// 创建文件
	 f, err := os.Create(filePath)
	 if err != nil {
		 log.Println("create file failed.", err)
		 w.WriteHeader(http.StatusBadGateway)
		 return
	 }
	 defer f.Close()
	 io.Copy(f, file)
 }

func get(w http.ResponseWriter, r *http.Request) {

	// 匹配"/"
	if r.URL.EscapedPath() == "/" {
		fmt.Fprint(w, "What you want from me.")
		return
	}

	f, e := os.Open(rootPath + r.URL.EscapedPath())
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	defer f.Close()
	http.ServeContent(w, r, "", time.Now(), f)
}

func isExist(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}



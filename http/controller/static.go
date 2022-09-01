package controller

import (
	"github.com/go-chinese-site/dreamgo/global"
	"github.com/go-chinese-site/dreamgo/route"
	"net/http"
	"strings"
)

// 静态文件控制器
type StaticController struct{}

func (self StaticController) RegisterRoutes() {
	route.HandleFunc("/static/", self.Default)
}

// Default 以/static/开头的URL为静态文件，使用 http.FileServer 直接处理
func (StaticController) Default(w http.ResponseWriter, r *http.Request) {
	reqURI := r.RequestURI
	//以/结尾的URL，直接返回404
	if strings.HasSuffix(reqURI, "/") {
		http.NotFound(w, r)
	} else {
		fileHandler := http.StripPrefix("/static/", http.FileServer(http.Dir(global.App.ProjectRoot+"/static")))
		fileHandler.ServeHTTP(w, r)
	}
}

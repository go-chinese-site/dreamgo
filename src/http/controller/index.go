// Copyright 2017 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// https://studygolang.com
// Author: polaris	polaris@studygolang.com

package controller

import (
	"datasource"
	"net/http"
	"route"
	"view"
)

var defaults = map[string]bool{
	"/":           true,
	"/index.html": true,
	"/index.htm":  true,
}

// IndexController 首页 controller
type IndexController struct{}

// RegisterRoute 注册路由
func (self IndexController) RegisterRoute() {
	route.HandleFunc("/", self.Home)
}

// Home 首页
func (IndexController) Home(w http.ResponseWriter, r *http.Request) {
	if _, ok := defaults[r.RequestURI]; !ok {
		http.NotFound(w, r)
		return
	}
	posts := datasource.DefaultDataSourcer.PostList()
	//	io.WriteString(w, "你说这是不是个玩笑")
	view.Render(w, r, "index.html", map[string]interface{}{"posts": posts})
}

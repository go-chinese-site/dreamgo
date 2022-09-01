// Copyright 2017 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// https://studygolang.com
// Author: polaris	polaris@studygolang.com

package controller

import (
	"github.com/go-chinese-site/dreamgo/datasource"
	"github.com/go-chinese-site/dreamgo/route"
	"github.com/go-chinese-site/dreamgo/view"
	"net/http"
)

type ArchiveController struct{}

// RegisterRoute 注册路由
func (self ArchiveController) RegisterRoute() {
	route.HandleFunc("/archives", self.List)
}

// List 处理归档列表请求
func (ArchiveController) List(w http.ResponseWriter, r *http.Request) {

	// 从数据源查询归档列表
	yearArchives := datasource.DefaultDataSourcer.PostArchive()
	// 渲染模板archives.html，并传入数据
	view.Render(w, r, "archives.html", map[string]interface{}{"archives": yearArchives})
}

// Copyright 2017 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// https://studygolang.com
// Author: polaris	polaris@studygolang.com

package controller

import (
	"github.com/go-chinese-site/dreamgo/route"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/go-chinese-site/dreamgo/datasource"
	"github.com/go-chinese-site/dreamgo/util"
	"github.com/go-chinese-site/dreamgo/view"
)

type PostController struct{}

// RegisterRoute register route
func (self PostController) RegisterRoute() {
	route.HandleFunc("/post/", self.Detail)
}

// Detail 处理文件详情请求
func (PostController) Detail(w http.ResponseWriter, r *http.Request) {
	// 获取文章文件名，即文章的路径
	filename := filepath.Base(r.RequestURI)
	if strings.HasSuffix(filename, ".md") {
		// 处理markdown
		datasource.DefaultDataSourcer.ServeMarkdown(w, r, filename)

	} else if strings.HasSuffix(filename, ".html") {
		// 根据路径查找文件
		post, err := datasource.DefaultDataSourcer.FindPost(util.Filename(filename))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// 渲染模板single.html，并传入数据
		view.Render(w, r, "single.html", map[string]interface{}{
			"post": post,
		})
	} else {
		// 返回404
		http.NotFound(w, r)
	}
}

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
	"net/url"
	"path/filepath"
	"sort"
)

type TagController struct{}

// RegisterRoute 注册路由
func (self TagController) RegisterRoute() {
	route.HandleFunc("/tag/", self.Detail)
	route.HandleFunc("/tags", self.List)
}

// Detail 处理标签详情请求
func (TagController) Detail(w http.ResponseWriter, r *http.Request) {
	// 从URL中获取标签名
	reqUrl, _ := url.ParseRequestURI(r.RequestURI)
	tagName := filepath.Base(reqUrl.Path)

	// 根据标签名查询标签
	tag := datasource.DefaultDataSourcer.FindTag(tagName)

	if tag != nil {
		// 渲染模板tag.html，并传入数据
		view.Render(w, r, "tag.html", map[string]interface{}{"tag": tag})
	} else {
		// 返回404
		http.NotFound(w, r)
	}
}

// List 处理标签列表请求
func (TagController) List(w http.ResponseWriter, r *http.Request) {

	// 从数据源获取标签列表
	tags := datasource.DefaultDataSourcer.TagList()
	// 按文章数量倒序排序

	sort.Slice(tags, func(i, j int) bool {
		return len(tags[i].Posts) > len(tags[j].Posts)
	})
	// 渲染模板tags.html，并传入数据
	view.Render(w, r, "tags.html", map[string]interface{}{"tags": tags})
}

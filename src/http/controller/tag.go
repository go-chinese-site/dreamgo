// Copyright 2017 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// https://studygolang.com
// Author: polaris	polaris@studygolang.com

package controller

import (
	"datasource"
	"net/http"
	"net/url"
	"path/filepath"
	"route"
	"view"
)

type TagController struct{}

// RegisterRoute 注册路由
func (self TagController) RegisterRoute() {
	route.HandleFunc("/tag/", self.Detail)
}

func (TagController) Detail(w http.ResponseWriter, r *http.Request) {
	reqUrl, _ := url.ParseRequestURI(r.RequestURI)
	tagName := filepath.Base(reqUrl.Path)

	tag := datasource.DefaultGithub.FindTag(tagName)
	if tag != nil {
		view.Render(w, r, "tag.html", map[string]interface{}{"tag": tag})
	} else {
		http.NotFound(w, r)
	}
}

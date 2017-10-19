// Copyright 2017 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// https://studygolang.com
// Author: polaris	polaris@studygolang.com

package controller

import (
	"net/http"
	"path/filepath"
	"route"
	"strings"

	"datasource"
	"util"
	"view"
)

type PostController struct{}

// RegisterRoute register route
func (self PostController) RegisterRoute() {
	route.HandleFunc("/post/", self.Detail)
}

func (PostController) Detail(w http.ResponseWriter, r *http.Request) {
	filename := filepath.Base(r.RequestURI)
	if strings.HasSuffix(filename, ".md") {

		datasource.DefaultDataSourcer.ServeMarkdown(w, r, filename)

	} else if strings.HasSuffix(filename, ".html") {
		post, err := datasource.DefaultDataSourcer.FindPost(util.Filename(filename))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		view.Render(w, r, "single.html", map[string]interface{}{
			"post": post,
		})
	} else {
		http.NotFound(w, r)
	}
}

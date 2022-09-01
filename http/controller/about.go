// Copyright 2017 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// https://studygolang.com
// Author: tk103331	tk103331@gmail.com

package controller

import (
	"github.com/go-chinese-site/dreamgo/datasource"
	"github.com/go-chinese-site/dreamgo/logger"
	"github.com/go-chinese-site/dreamgo/route"
	"github.com/go-chinese-site/dreamgo/view"
	"net/http"
)

type AboutController struct{}

func (self AboutController) RegisterRoutes() {
	route.HandleFunc("/about", self.Detail)
}

func (AboutController) Detail(w http.ResponseWriter, r *http.Request) {
	about, err := datasource.DefaultDataSourcer.AboutPost()
	if err == nil {
		view.Render(w, r, "about.html", map[string]interface{}{"about": about})
	} else {
		logger.Instance().Error("get about.md error " + err.Error())
		http.NotFound(w, r)
	}
}

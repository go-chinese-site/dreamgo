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

type ArchiveController struct{}

// RegisterRoute 注册路由
func (self ArchiveController) RegisterRoute() {
	route.HandleFunc("/archives", self.List)
}

func (ArchiveController) List(w http.ResponseWriter, r *http.Request) {
	yearArchives := datasource.DefaultDataSourcer.PostArchive()

	view.Render(w, r, "archives.html", map[string]interface{}{"archives": yearArchives})
}

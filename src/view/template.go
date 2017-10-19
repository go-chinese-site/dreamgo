// Copyright 2017 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// https://studygolang.com
// Author: polaris	polaris@studygolang.com

package view

import (
	"config"
	"html/template"
	"net/http"
	"time"

	"global"
)

// funcMap is the customize template functions
var funcMap = template.FuncMap{
	"noescape": func(s string) template.HTML {
		return template.HTML(s)
	},
	"formatTime": func(t time.Time, layout string) string {
		return t.Format(layout)
	},
}

// Render 渲染模板并输出
func Render(w http.ResponseWriter, r *http.Request, htmlFile string, data map[string]interface{}) {
	if data == nil {
		data = make(map[string]interface{})
	}
	data["app"] = global.App
	data["site_name"] = config.YamlConfig.Get("setting.site_name").String()
	data["title"] = config.YamlConfig.Get("setting.title").String()
	data["subtitle"] = config.YamlConfig.Get("setting.subtitle").String()

	// 加载布局模板layout.html
	tpl, err := template.New("layout.html").Funcs(funcMap).
		ParseFiles(global.App.TemplateDir+"layout.html", global.App.TemplateDir+htmlFile)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// 加载seo关键词和描述
	if seoTpl := tpl.Lookup("seo"); seoTpl == nil {
		seoKeywords := config.YamlConfig.Get("seo.keywords").String()
		seoDescription := config.YamlConfig.Get("seo.description").String()

		tpl.Parse(`{{define "seo"}}
			<meta name="keywords" content="` + seoKeywords + `">
			<meta name="description" content="` + seoDescription + `">
		{{end}}`)
	}
	startTime := r.Context().Value("start_time").(time.Time)
	data["response_time"] = time.Since(startTime)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(http.StatusOK)
	// 渲染模板，并输出到w
	err = tpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Copyright 2017 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// https://studygolang.com
// Author: polaris	polaris@studygolang.com

package route

import (
	"context"
	"net/http"
	"time"
)

func HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	DefaultBlogMux.HandleFunc(pattern, handler)
}

// BlogMux 路由处理器，扩展http.ServeMux
type BlogMux struct {
	*http.ServeMux
}

// DefaultBlogMux 默认路由处理器
var DefaultBlogMux = NewBlogMux()

func NewBlogMux() *BlogMux {
	return &BlogMux{ServeMux: http.DefaultServeMux}
}

// ServeHTTP 路由分发方法，封装 http.DefaultServeMux.ServeHTTP()
func (this *BlogMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 创建上下文，并写入start_time
	ctx := context.WithValue(r.Context(), "start_time", time.Now())
	// 使用上下文
	r = r.WithContext(ctx)
	// 调用http.DefaultServeMux的路由分发方法
	this.ServeMux.ServeHTTP(w, r)
}

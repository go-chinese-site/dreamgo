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

type BlogMux struct {
	*http.ServeMux
}

var DefaultBlogMux = NewBlogMux()

func NewBlogMux() *BlogMux {
	return &BlogMux{ServeMux: http.DefaultServeMux}
}

func (this *BlogMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := context.WithValue(r.Context(), "start_time", time.Now())
	r = r.WithContext(ctx)

	this.ServeMux.ServeHTTP(w, r)
}

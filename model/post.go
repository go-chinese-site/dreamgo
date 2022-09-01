// Copyright 2017 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// https://studygolang.com
// Author: polaris	polaris@studygolang.com

package model

import "time"

// 文章
type Post struct {
	Content string `yaml:"content"`
	*Meta
}

type Meta struct {
	Title   string   `yaml:"title"`
	Path    string   `yaml:"path"`
	PubTime string   `yaml:"pub_time"`
	Tags    []string `yaml:"tags"`

	PostTime time.Time `yaml:"post_time"`
}

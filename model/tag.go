// Copyright 2017 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// https://studygolang.com
// Author: tk103331	tk103331@gmail.com
package model

// 标签
type Tag struct {
	Name  string  `yaml:"name"`
	Posts []*Post `yaml:"posts"`
}

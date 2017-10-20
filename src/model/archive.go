// Copyright 2017 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// https://studygolang.com
// Author: polaris	polaris@studygolang.com

package model

// YearArchive 归档
type YearArchive struct {
	Year int `yaml:"year"`

	MonthArchives []*MonthArchive `yaml:"month_archive"`
}

type MonthArchive struct {
	Month int `yaml:"month"`

	Posts []*Post `yaml:"posts"`
}

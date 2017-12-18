// Copyright 2017 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// https://studygolang.com
// Author: polaris	polaris@studygolang.com

package datasource_test

import (
	"datasource"
	"global"
	"os"
	"strings"

	"testing"
)

// DefaultGithub git数据源结构体实例
var DefaultGithub *datasource.GithubRepo

func setup() {
	cwd, _ := os.Getwd()
	pos := strings.LastIndex(cwd, "src")
	global.App.ProjectRoot = cwd[:pos]
	DefaultGithub = datasource.NewGithub()
}

func TestGenIndexYaml(t *testing.T) {
	setup()

	DefaultGithub.GenIndexYaml()
}

func TestGenArchiveYaml(t *testing.T) {
	setup()

	DefaultGithub.GenArchiveYaml()
}

func TestGenTagsYaml(t *testing.T) {
	setup()

	DefaultGithub.GenTagsYaml()
}

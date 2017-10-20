// Copyright 2017 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// https://studygolang.com
// Author: polaris	polaris@studygolang.com

// global 全局信息
package global

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"
)

// Build 构建信息，从 git 仓库获取
var Build string

type app struct {
	Name      string
	Build     string
	Version   string
	BuildDate time.Time

	ProjectRoot string
	TemplateDir string

	Copyright string

	LaunchTime time.Time

	Host string
	Port string

	locker sync.Mutex
}

// App is the App Info
var App = &app{}

var showVersion = flag.Bool("version", false, "Print version of this binary")

func init() {
	App.Name = os.Args[0]
	App.Version = "V1.0.0"
	App.Build = Build
	App.LaunchTime = time.Now()
	// 查找可执行程序的路径
	binaryPath, err := exec.LookPath(os.Args[0])
	if err != nil {
		panic(err)
	}
	// 获取可执行程序的绝对路径
	binaryPath, err = filepath.Abs(binaryPath)
	if err != nil {
		panic(err)
	}
	// 获取可执行程序的文件信息
	fileInfo, err := os.Stat(binaryPath)
	if err != nil {
		panic(err)
	}
	// 构建时间为可执行程序的修改时间
	App.BuildDate = fileInfo.ModTime()
	App.Copyright = fmt.Sprintf("%d", time.Now().Year())
}

// InitPath 初始化相关路径，包括项目根目录、模板目录
func (this *app) InitPath() {
	App.setProjectRoot()

	App.SetTemplateDir("default")
}

// Uptime calculates the duration of lauching
func (this *app) Uptime() time.Duration {
	this.locker.Lock()
	defer this.locker.Unlock()
	return time.Now().Sub(this.LaunchTime)
}

func (this *app) setProjectRoot() {
	curFilename := os.Args[0]

	binaryPath, err := exec.LookPath(curFilename)
	if err != nil {
		panic(err)
	}

	binaryPath, err = filepath.Abs(binaryPath)
	if err != nil {
		panic(err)
	}

	projectRoot := filepath.Dir(filepath.Dir(binaryPath))

	this.ProjectRoot = projectRoot + "/"
}

// SetTemplateDir 设置模板目录
func (this *app) SetTemplateDir(theme string) {
	this.TemplateDir = this.ProjectRoot + "template/theme/" + theme + "/"
}

// PrintVersion prints current version info
func PrintVersion(w io.Writer) {
	if !flag.Parsed() {
		flag.Parse()
	}

	if showVersion == nil || !*showVersion {
		return
	}

	fmt.Fprintf(w, "Binary: %s\n", App.Name)
	fmt.Fprintf(w, "Version: %s\n", App.Version)
	fmt.Fprintf(w, "Build: %s\n", App.Build)
	fmt.Fprintf(w, "Compile date: %s\n", App.BuildDate.Format("2006-01-02 15:04:05"))
	os.Exit(0)
}

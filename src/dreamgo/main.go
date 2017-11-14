// Copyright 2017 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// https://studygolang.com
// Author: polaris	polaris@studygolang.com

package main

import (
	"config"
	"datasource"
	"flag"
	"global"
	"http/controller"
	"log"
	"logger"
	"math/rand"
	"net/http"
	"route"
	"strings"
	"time"
)

var configFile string

func init() {
	rand.Seed(time.Now().Unix())

	flag.StringVar(&configFile, "config", "config/env.yml", "The config file. Default is $ProjectRoot/config/env.yml")
}

func main() {
	// 日志
	logger := logger.Init("dreamgo")
	logger.Info("main ... ")
	// 解析命令行参数
	flag.Parse()
	// 初始化程序路径
	global.App.InitPath()

	if strings.HasPrefix(configFile, "/") { //以/开头为绝对路径，直接解析
		config.Parse(configFile)
	} else { // 相对路径，以程序根目录为基础解析
		config.Parse(global.App.ProjectRoot + configFile)
	}
	datasource.Init()
	go updateGitDataSource()
	// 设置模板目录，默认为default
	global.App.SetTemplateDir(config.YamlConfig.MustValue("theme", "default"))
	// 从配置文件中获取监听IP和端口
	global.App.Host = config.YamlConfig.Get("listen.host").String()
	global.App.Port = config.YamlConfig.Get("listen.port").String()

	addr := global.App.Host + ":" + global.App.Port
	// 注册路由
	controller.RegisterRoutes()
	// 启动监听，使用封装的 route.DefaultBlogMux 处理http请求
	log.Fatal(http.ListenAndServe(addr, route.DefaultBlogMux))
}

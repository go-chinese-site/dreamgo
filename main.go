// Copyright 2017 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// https://studygolang.com
// Author: polaris	polaris@studygolang.com

package main

import (
	"flag"
	"github.com/go-chinese-site/dreamgo/config"
	"github.com/go-chinese-site/dreamgo/datasource"
	"github.com/go-chinese-site/dreamgo/global"
	"github.com/go-chinese-site/dreamgo/http/controller"
	"github.com/go-chinese-site/dreamgo/logger"
	"github.com/go-chinese-site/dreamgo/route"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

var configFile string

func init() {
	rand.Seed(time.Now().Unix())

	flag.StringVar(&configFile, "conf", "conf/env.yml", "The conf file. Default is $ProjectRoot/conf/env.yml")
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

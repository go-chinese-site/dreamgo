// Copyright 2017 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// https://studygolang.com
// Author: polaris	polaris@studygolang.com

package main

import (
	"config"
	"flag"
	"http/controller"
	"log"
	"math/rand"
	"net/http"
	"route"
	"strings"
	"time"

	"global"
)

var configFile string

func init() {
	rand.Seed(time.Now().Unix())

	flag.StringVar(&configFile, "config", "config/env.yml", "The config file. Default is $ProjectRoot/config/env.yml")
}

func main() {
	flag.Parse()

	if strings.HasPrefix(configFile, "/") {
		config.Parse(configFile)
	} else {
		config.Parse(global.App.ProjectRoot + configFile)
	}

	global.App.InitPath()

	go updateGitDataSource()

	global.App.SetTemplateDir(config.YamlConfig.MustValue("theme", "default"))

	global.App.Host = config.YamlConfig.Get("listen.host").String()
	global.App.Port = config.YamlConfig.Get("listen.port").String()

	addr := global.App.Host + ":" + global.App.Port

	controller.RegisterRoutes()

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(global.App.ProjectRoot+"/static"))))

	log.Fatal(http.ListenAndServe(addr, route.DefaultBlogMux))
}

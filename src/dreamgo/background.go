// Copyright 2017 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// https://studygolang.com
// Author: polaris	polaris@studygolang.com

package main

import (
	"config"
	"datasource"
	"log"
	"os"
	"os/exec"

	"github.com/robfig/cron"

	"global"
	"util"
)

func updateGitDataSource() {
	typ := config.YamlConfig.Get("datasource.type").String()
	if typ != datasource.TypeGit {
		return
	}
	// 检查文章目录(data/post/)是否存在,不存在则克隆远程仓库
	gitRepoDir := global.App.ProjectRoot + datasource.PostDir
	if !util.Exist(gitRepoDir) {
		if err := os.MkdirAll(gitRepoDir, os.ModePerm); err != nil {
			panic(err)
		}

		cloneRepo(gitRepoDir)
	}

	gitFolder := gitRepoDir + ".git"
	for {
		if util.Exist(gitFolder) {
			break
		}

		cloneRepo(gitRepoDir)
	}
	// 解析仓库文件，生成首页、归档、标签数据
	datasource.DefaultGithub.GenIndexYaml()
	datasource.DefaultGithub.GenArchiveYaml()
	datasource.DefaultGithub.GenTagsYaml()

	// 定时每天自动更新仓库，并生成首页、归档、标签数据
	c := cron.New()
	c.AddFunc("@daily", func() {
		datasource.DefaultGithub.Pull(gitRepoDir)
		datasource.DefaultGithub.GenIndexYaml()
		datasource.DefaultGithub.GenArchiveYaml()
		datasource.DefaultGithub.GenTagsYaml()
	})
	c.Start()
}

// 使用git clone命令克隆文章仓库
func cloneRepo(gitRepoDir string) {
	cmdName := "git"
	pullArgs := []string{"clone", config.YamlConfig.Get("datasource.url").String(), "."}

	cmd := exec.Command(cmdName, pullArgs...)
	cmd.Dir = gitRepoDir

	if err := cmd.Run(); err != nil {
		log.Printf("error clone master at %s: %v", gitRepoDir, err)
		return
	}
}

// 更新mysql数据
func updateMysqlDataSource() {
	typ := config.YamlConfig.Get("datasource.type").String()
	if typ != datasource.TypeMysql {
		return
	}
	// 检查文章目录(data/post/)是否存在，不存在则连接mysql生成
	mysqlRepoDir := global.App.ProjectRoot + datasource.PostDir
	if !util.Exist(mysqlRepoDir) {
		if err := os.MkdirAll(mysqlRepoDir, os.ModePerm); err != nil {
			panic(err)
		}
		// 待实现
	}

}

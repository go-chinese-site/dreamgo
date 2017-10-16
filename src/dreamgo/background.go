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

	datasource.DefaultGithub.GenIndexYaml()
	datasource.DefaultGithub.GenArchiveYaml()
	datasource.DefaultGithub.GenTagsYaml()

	c := cron.New()

	c.AddFunc("@daily", func() {
		datasource.DefaultGithub.Pull(gitRepoDir)
		datasource.DefaultGithub.GenIndexYaml()
		datasource.DefaultGithub.GenArchiveYaml()
		datasource.DefaultGithub.GenTagsYaml()
	})

	c.Start()
}

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

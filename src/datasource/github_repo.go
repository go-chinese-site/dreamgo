// Copyright 2017 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// https://studygolang.com
// Author: polaris	polaris@studygolang.com

package datasource

import (
	"global"
	"io/ioutil"
	"log"
	"model"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"sort"
	"strings"
	"time"
	"util"

	"github.com/pkg/errors"
	"github.com/russross/blackfriday"
	yaml "gopkg.in/yaml.v2"
)

const (
	// PostDir is the directory of storing posts
	PostDir = "data/post/"

	IndexFile   = "index.yaml"
	ArchiveFile = "archive.yaml"
)

type GithubRepo struct{}

var DefaultGithub = &GithubRepo{}

func (self GithubRepo) PostList() []*model.Post {
	in, err := ioutil.ReadFile(global.App.ProjectRoot + PostDir + IndexFile)
	if err != nil {
		return nil
	}
	posts := make([]*model.Post, 0)
	err = yaml.Unmarshal(in, &posts)
	if err != nil {
		return nil
	}

	return posts
}

func (self GithubRepo) PostArchive() []*model.YearArchive {
	in, err := ioutil.ReadFile(global.App.ProjectRoot + PostDir + ArchiveFile)
	if err != nil {
		return nil
	}
	yearArchives := make([]*model.YearArchive, 0)
	err = yaml.Unmarshal(in, &yearArchives)
	if err != nil {
		return nil
	}

	return yearArchives
}

func (self GithubRepo) ServeMarkdown(w http.ResponseWriter, r *http.Request, filename string) {
	http.ServeFile(w, r, global.App.ProjectRoot+PostDir+util.Filename(filename)+"/post.md")
}

var titleReg = regexp.MustCompile(`^#\s(.+)`)

func (self GithubRepo) FindPost(path string) (*model.Post, error) {
	postDir := global.App.ProjectRoot + PostDir + path

	post, err := self.genOnePost(postDir, path)
	if err == nil {
		post.Content, err = replaceCodeParts(blackfriday.MarkdownCommon([]byte(post.Content)))
	}

	return post, err
}

func (self GithubRepo) Pull(gitRepoDir string) error {
	cmdName := "git"
	pullArgs := []string{"pull", "origin", "master"}

	cmd := exec.Command(cmdName, pullArgs...)
	cmd.Dir = gitRepoDir

	if err := cmd.Run(); err != nil {
		log.Printf("error pulling master at %s: %v", gitRepoDir, err)
		return err
	}

	return nil
}

func (self GithubRepo) GenIndexYaml() {
	posts := self.fetchPosts()

	length := 20
	if len(posts) < length {
		length = len(posts)
	}

	buf, err := yaml.Marshal(posts[:length])
	if err != nil {
		log.Printf("gen index yaml error:%v\n", err)
		return
	}

	indexYaml := global.App.ProjectRoot + PostDir + IndexFile
	ioutil.WriteFile(indexYaml, buf, 0777)
}

func (self GithubRepo) GenArchiveYaml() {
	posts := self.fetchPosts()

	yearArchiveMap := make(map[int]*model.YearArchive)

	for _, post := range posts {
		post.Content = ""

		year := post.PostTime.Year()
		month := int(post.PostTime.Month())

		if yearArchive, ok := yearArchiveMap[year]; ok {
			monthExists := false
			for _, monthArchive := range yearArchive.MonthArchives {
				if monthArchive.Month == month {
					monthArchive.Posts = append(monthArchive.Posts, post)
					monthExists = true
					break
				}
			}

			if !monthExists {
				yearArchive.MonthArchives = append(yearArchive.MonthArchives, &model.MonthArchive{
					Month: month,
					Posts: []*model.Post{post},
				})
			}

		} else {
			monthArchive := &model.MonthArchive{
				Month: month,
				Posts: []*model.Post{post},
			}
			yearArchive = &model.YearArchive{
				Year:          year,
				MonthArchives: []*model.MonthArchive{monthArchive},
			}

			yearArchiveMap[year] = yearArchive
		}
	}

	yearArchives := make([]*model.YearArchive, 0, len(yearArchiveMap))
	for _, yearArchive := range yearArchiveMap {
		yearArchives = append(yearArchives, yearArchive)
	}

	sort.Slice(yearArchives, func(i, j int) bool {
		return yearArchives[i].Year > yearArchives[j].Year
	})

	buf, err := yaml.Marshal(yearArchives)
	if err != nil {
		log.Printf("gen archives yaml error:%v\n", err)
		return
	}

	archiveYaml := global.App.ProjectRoot + PostDir + ArchiveFile
	ioutil.WriteFile(archiveYaml, buf, 0777)
}

func (self GithubRepo) fetchPosts() []*model.Post {
	var (
		posts = make([]*model.Post, 0, 31)

		post *model.Post
		err  error
	)

	postDir := global.App.ProjectRoot + PostDir
	names := util.ScanDir(postDir)
	for _, name := range names {
		if util.IsFile(postDir + name) {
			continue
		}

		if name == ".git" {
			continue
		}

		post, err = self.genOnePost(postDir+name, name)
		if err != nil {
			continue
		}

		pos := strings.Index(post.Content, `<!--more-->`)
		if pos > 0 {
			post.Content = post.Content[:pos]
		}

		posts = append(posts, post)
	}

	sort.Slice(posts, func(i, j int) bool {
		return posts[i].PubTime > posts[j].PubTime
	})

	return posts
}

func (self GithubRepo) genOnePost(postDir, path string) (*model.Post, error) {
	markdown, err := ioutil.ReadFile(postDir + "/post.md")
	if err != nil {
		return nil, errors.Wrap(err, "read post.md error")
	}

	var meta = &model.Meta{}
	metaBytes, err := ioutil.ReadFile(postDir + "/meta.yml")
	if err == nil {
		err = yaml.Unmarshal(metaBytes, meta)
		if err != nil {
			return nil, errors.Wrap(err, "yaml unmarshal meta.yml error")
		}

		meta.PostTime = self.parsePubTime(meta.PubTime)
	} else {
		meta.Path = path + ".html"
		fileInfo, _ := os.Stat(postDir + "/post.md")
		meta.PostTime = fileInfo.ModTime()
		meta.PubTime = meta.PostTime.Format("2006-01-02 15:04")
		matches := titleReg.FindStringSubmatch(string(markdown))
		if len(matches) > 2 {
			meta.Title = matches[1]
		} else {
			meta.Title = path
		}
	}

	post := &model.Post{
		Content: string(markdown),
		Meta:    meta,
	}

	return post, nil
}

func (self GithubRepo) parsePubTime(pubTime string) time.Time {
	layouts := []string{
		"2006-01-02 15:04:05",
		"2006-01-02 15:04",
		"2006年01月02 15:04:05",
		"2006年01月02 15:04",
	}

	for _, layout := range layouts {

		t, err := time.ParseInLocation(layout, pubTime, time.Local)
		if err != nil {
			continue
		}

		return t
	}

	return time.Now()
}

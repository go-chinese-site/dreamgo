// Copyright 2017 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// https://studygolang.com
// Author: polaris	polaris@studygolang.com

package datasource

import (
	"config"
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
	"github.com/robfig/cron"
	"github.com/russross/blackfriday"
	yaml "gopkg.in/yaml.v2"
)

const (
	// PostDir 文章存放目录
	PostDir = "data/post/"

	// IndexFile 首页数据文件
	IndexFile = "index.yaml"
	// ArchiveFile 归档数据文件
	ArchiveFile = "archive.yaml"
	// TagsFile 标签数据文件
	TagsFile = "tags.yaml"
	// FriendFile 友情链接数据文件
	FriendFile = "friends.yaml"
)

// GithubRepo git数据源结构体
type GithubRepo struct{}

// NewGithub 创建git数据源实例，相当于构造方法
func NewGithub() *GithubRepo {
	return &GithubRepo{}
}

// PostList 读取文章列表
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

// PostArchive 读取归档列表
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

// ServeMarkdown 处理查看 Markdown 请求
func (self GithubRepo) ServeMarkdown(w http.ResponseWriter, r *http.Request, filename string) {
	http.ServeFile(w, r, global.App.ProjectRoot+PostDir+util.Filename(filename)+"/post.md")
}

var titleReg = regexp.MustCompile(`^#\s(.+)`)

// FindPost 根据路径查找文章
func (self GithubRepo) FindPost(path string) (*model.Post, error) {
	postDir := global.App.ProjectRoot + PostDir + path

	post, err := self.genOnePost(postDir, path)
	if err == nil {
		post.Content, err = replaceCodeParts(blackfriday.MarkdownCommon([]byte(post.Content)))
	}

	return post, err
}

// Pull 使用 git pull origin master 命令从远程仓库更新文章
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

// GenIndexYaml 生成首页数据文件index.yaml
func (self GithubRepo) GenIndexYaml() {
	posts := self.fetchPosts()
	// 首页最多显示20篇文章
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

// GenArchiveYaml 生成归档数据文件archive.yaml
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

// GenTagsYaml 生成标签数据文件tags.yaml
func (self GithubRepo) GenTagsYaml() {
	allPosts := self.fetchPosts()
	tagMap := make(map[string][]*model.Post)
	// 遍历所有文章对象，分析出标签数据
	for _, post := range allPosts {
		post.Content = ""
		for _, tag := range post.Tags {
			posts, ok := tagMap[tag]
			if !ok {
				posts = make([]*model.Post, 0)
			}
			posts = append(posts, post)
			tagMap[tag] = posts
		}
	}
	// 组装标签列表
	tags := make([]*model.Tag, 0)
	for tag, posts := range tagMap {
		sort.Slice(posts, func(i, j int) bool {
			return posts[i].PubTime > posts[j].PubTime
		})
		tags = append(tags, &model.Tag{Name: tag, Posts: posts})
	}
	// 按文件数量倒序排序
	sort.Slice(tags, func(i, j int) bool {
		return len(tags[i].Posts) > len(tags[j].Posts)
	})

	buf, err := yaml.Marshal(tags)
	if err != nil {
		log.Printf("gen tags yaml error:%v\n", err)
		return
	}

	tagsYaml := global.App.ProjectRoot + PostDir + TagsFile
	ioutil.WriteFile(tagsYaml, buf, 0777)
}

// fetchPosts 读取所有文章数据，遍历目录，解析每个目录中的meta.yaml和post.md
func (self GithubRepo) fetchPosts() []*model.Post {
	var (
		posts = make([]*model.Post, 0, 31)

		post *model.Post
		err  error
	)
	// 遍历 data/post 下的目录
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
	// 按发布时间倒序排序
	sort.Slice(posts, func(i, j int) bool {
		return posts[i].PubTime > posts[j].PubTime
	})

	return posts
}

// genOnePost 解析meta.yaml和post.md文件生成model.Post对象
func (self GithubRepo) genOnePost(postDir, path string) (*model.Post, error) {
	// 从post.md中读取文章内容
	markdown, err := ioutil.ReadFile(postDir + "/post.md")
	if err != nil {
		return nil, errors.Wrap(err, "read post.md error")
	}
	// 从meta.yml文件读取文章信息
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

// parsePubTime 解析发布时间
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

// TagList 读取标签列表
func (self GithubRepo) TagList() []*model.Tag {
	in, err := ioutil.ReadFile(global.App.ProjectRoot + PostDir + TagsFile)
	if err != nil {
		return nil
	}
	tags := make([]*model.Tag, 0)
	err = yaml.Unmarshal(in, &tags)
	if err != nil {
		return nil
	}

	return tags
}

// FindTag 通过标签名查找标签
func (self GithubRepo) FindTag(tagName string) *model.Tag {
	tags := self.TagList()
	for _, tag := range tags {
		if tag.Name == tagName {
			return tag
		}
	}
	return nil
}

// AboutPost 获取关于页
func (self GithubRepo) AboutPost() (*model.Post, error) {
	// 从 about.md 中读取关于内容
	postDir := global.App.ProjectRoot + PostDir
	markdown, err := ioutil.ReadFile(postDir + "/about.md")
	if err != nil {
		return nil, errors.Wrap(err, "read about.md error")
	}
	// 关于页不需要 meta.yml
	var meta = &model.Meta{}
	post := &model.Post{
		Content: string(markdown),
		Meta:    meta,
	}
	return post, nil
}

// UpdateDataSource 更新数据
func (self GithubRepo) UpdateDataSource() {
	// 检查文章目录(data/post/)是否存在,不存在则克隆远程仓库
	gitRepoDir := global.App.ProjectRoot + PostDir
	if !util.Exist(gitRepoDir) {
		if err := os.MkdirAll(gitRepoDir, os.ModePerm); err != nil {
			panic(err)
		}

		self.cloneRepo(gitRepoDir)
	}

	gitFolder := gitRepoDir + ".git"
	for {
		if util.Exist(gitFolder) {
			break
		}

		self.cloneRepo(gitRepoDir)
	}
	// 解析仓库文件，生成首页、归档、标签数据
	self.GenIndexYaml()
	self.GenArchiveYaml()
	self.GenTagsYaml()

	// 定时每天自动更新仓库，并生成首页、归档、标签数据
	c := cron.New()
	c.AddFunc("@daily", func() {
		self.Pull(gitRepoDir)
		self.GenIndexYaml()
		self.GenArchiveYaml()
		self.GenTagsYaml()
	})
	c.Start()
}

// 使用git clone命令克隆文章仓库
func (self GithubRepo) cloneRepo(gitRepoDir string) {
	cmdName := "git"
	pullArgs := []string{"clone", config.YamlConfig.Get("datasource.url").String(), "."}

	cmd := exec.Command(cmdName, pullArgs...)
	cmd.Dir = gitRepoDir

	if err := cmd.Run(); err != nil {
		log.Printf("error clone master at %s: %v", gitRepoDir, err)
		return
	}
}

// GetFriends 友情链接
func (self GithubRepo) GetFriends() ([]*model.Friend, error) {
	// 从friends.yaml 中读取友情链接内容
	in, err := ioutil.ReadFile(global.App.ProjectRoot + PostDir + FriendFile)

	if err != nil {
		return nil, errors.Wrap(err, "read friends.yaml error")
	}

	friends := make([]*model.Friend, 0)
	err = yaml.Unmarshal(in, &friends)
	if err != nil {
		return nil, errors.Wrap(err, "Unmarshal friends.yaml error")
	}
	return friends, nil
}

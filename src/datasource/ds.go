// Copyright 2017 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// https://studygolang.com
// Author: polaris	polaris@studygolang.com

package datasource

import (
	"bytes"
	"config"
	"model"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/pkg/errors"
	"github.com/sourcegraph/syntaxhighlight"
)

// 数据源类型
const (
	TypeGit   = "git"
	TypeMysql = "mysql"
)

// DataSourcer 数据源接口
type DataSourcer interface {
	PostList() []*model.Post
	PostArchive() []*model.YearArchive
	ServeMarkdown(w http.ResponseWriter, r *http.Request, filename string)
	FindPost(path string) (*model.Post, error)
	TagList() []*model.Tag
	FindTag(tagName string) *model.Tag
	AboutPost() (*model.Post, error)
	UpdateDataSource()
	GetFriends() ([]*model.Friend, error)
}

// DefaultDataSourcer 默认数据源
var DefaultDataSourcer DataSourcer

// Init 数据源初始化
func Init() {

	dataSourcerType := config.YamlConfig.Get("datasource.type").String()
	switch dataSourcerType {
	case "git":
		DefaultDataSourcer = NewGithub()
	case "mongodb":
		DefaultDataSourcer = NewMongoDB()
	case "mysql":
		DefaultDataSourcer = NewMysql(config.YamlConfig.Get("datasource.mysqlAddr").String())
	default:
		DefaultDataSourcer = NewGithub()
	}
	go DefaultDataSourcer.UpdateDataSource()
}

func replaceCodeParts(htmlFile []byte) (string, error) {
	byteReader := bytes.NewReader(htmlFile)
	doc, err := goquery.NewDocumentFromReader(byteReader)
	if err != nil {
		return "", errors.Wrap(err, "error while parsing html")
	}

	// find code-parts via css selector and replace them with highlighted versions
	doc.Find("code[class*=\"language-\"]").Each(func(i int, s *goquery.Selection) {
		oldCode := s.Text()
		formatted, _ := syntaxhighlight.AsHTML([]byte(oldCode))
		s.SetHtml(string(formatted))
	})
	new, err := doc.Html()
	if err != nil {
		return "", errors.Wrap(err, "error while generating html")
	}

	// replace unnecessarily added html tags
	new = strings.Replace(new, "<html><head></head><body>", "", 1)
	new = strings.Replace(new, "</body></html>", "", 1)
	return new, nil
}

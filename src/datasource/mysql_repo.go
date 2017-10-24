package datasource

import (
	"model"
	"net/http"
)

type MysqlRepo struct{}

var DefaultMysql = NewMysql()

func NewMysql() *MysqlRepo {
	return &MysqlRepo{}
}

// PostList 读取文章列表
func (self *MysqlRepo) PostList() []*model.Post {
	return nil
}

// PostArchive 读取归档列表
func (self *MysqlRepo) PostArchive() []*model.YearArchive {
	return nil
}

// ServeMarkdown 处理查看 Markdown 请求
func (self *MysqlRepo) ServeMarkdown(w http.ResponseWriter, r *http.Request, filename string) {

}

// FindPost 根据路径查找文章
func (self *MysqlRepo) FindPost(path string) (*model.Post, error) {
	return nil, nil
}

// TagList 读取标签列表
func (self *MysqlRepo) TagList() []*model.Tag {
	return nil
}

// FindTag 通过标签名查找标签
func (self *MysqlRepo) FindTag(tagName string) *model.Tag {
	return nil
}

// AboutPost 获取关于页
func (self *MysqlRepo) AboutPost() (*model.Post, error) {
	return nil, nil
}

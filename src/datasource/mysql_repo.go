package datasource

import (
	"database/sql"
	"fmt"
	"global"
	"io/ioutil"
	"log"
	"model"
	"net/http"
	"os"
	"sort"
	"strconv"
	"time"
	"util"

	_ "github.com/go-sql-driver/mysql" // data source
	"github.com/pkg/errors"
	"github.com/robfig/cron"
	"github.com/russross/blackfriday"
	"gopkg.in/yaml.v2"
)

// MysqlRepo mysql 数据源结构体
type MysqlRepo struct {
	db                    *sql.DB
	selectTag             *sql.Stmt
	selectArticleById     *sql.Stmt
	selectArticleIndex    *sql.Stmt
	selectArticleTagsById *sql.Stmt
	selectArticleArchives *sql.Stmt
	selectArticlesByTag   *sql.Stmt
	selectFriends         *sql.Stmt
}

type articleInfo struct {
	Id      int64  `json:"id"`
	Title   string `json:"title"`
	PubTime int64  `json:"pub_time"`
	Content string `json:"content"`
}

type tagInfo struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

type friendInfo struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
	Link string `json:"link"`
	Logo string `json:"logo"`
}

// NewMysql 创建mysql数据源实例，相当于构造方法
func NewMysql(dbParams string) *MysqlRepo {
	db, err := sql.Open("mysql", dbParams)
	if err != nil {
		log.Fatalf("Couldn't connect to database: %s", err)
	}

	return &MysqlRepo{
		db:                    db,
		selectTag:             prepare(db, "SELECT * FROM `tag`"),
		selectArticleById:     prepare(db, "SELECT * FROM `article` WHERE `id`= ?"),
		selectArticleIndex:    prepare(db, "SELECT * FROM `article` ORDER BY `pub_time` DESC LIMIT 20"),
		selectArticleTagsById: prepare(db, "SELECT t.`name` FROM `article_tag` at LEFT JOIN `tag` t ON at.`tag_id`=t.`id` WHERE `article_id`= ?"),
		selectArticleArchives: prepare(db, "SELECT `id`,`title`,`pub_time` FROM `article`"),
		selectArticlesByTag:   prepare(db, "SELECT a.`id`,a.`title`,a.`pub_time` FROM `article` a LEFT JOIN `article_tag` at ON a.`id`=at.`article_id` WHERE at.`tag_id`=?"),
		selectFriends:         prepare(db, "SELECT * FROM `friend_link`"),
	}
}

func prepare(db *sql.DB, sql string) *sql.Stmt {
	stmt, err := db.Prepare(sql)
	if err != nil {
		log.Fatalf("Prepare SQL '%s' failed: %s", sql, err)
	}
	return stmt
}

// PostList 读取文章列表
func (self *MysqlRepo) PostList() []*model.Post {
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
func (self *MysqlRepo) PostArchive() []*model.YearArchive {
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
func (self *MysqlRepo) ServeMarkdown(w http.ResponseWriter, r *http.Request, filename string) {
	http.ServeFile(w, r, global.App.ProjectRoot+PostDir+util.Filename(filename)+"/post.md")
}

// FindPost 根据路径查找文章
func (self *MysqlRepo) FindPost(path string) (*model.Post, error) {
	id, err := strconv.Atoi(path)
	if err != nil {
		log.Printf("Invalid path :%s\n", err)
		return nil, fmt.Errorf("文章不存在")
	}
	row := self.selectArticleById.QueryRow(id)
	info := articleInfo{}
	row.Scan(&info.Id, &info.Title, &info.PubTime, &info.Content)

	post := self.genOnePost(info)
	rows, err := self.selectArticleTagsById.Query(info.Id)
	if err != nil {
		log.Printf("Query article tags error:%s", err)
		return nil, fmt.Errorf("文章不存在")
	}
	var tags []string
	for rows.Next() {
		var tagName string
		err = rows.Scan(&tagName)
		if err != nil {
			log.Printf("Scan tag error: %s\n", err)
		}
		tags = append(tags, tagName)
	}
	post.Tags = tags

	post.Content, err = replaceCodeParts(blackfriday.MarkdownCommon([]byte(post.Content)))

	return post, err
}

// TagList 读取标签列表
func (self *MysqlRepo) TagList() []*model.Tag {
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
func (self *MysqlRepo) FindTag(tagName string) *model.Tag {
	tags := self.TagList()
	for _, tag := range tags {
		if tag.Name == tagName {
			return tag
		}
	}
	return nil
}

// AboutPost 获取关于页
func (self *MysqlRepo) AboutPost() (*model.Post, error) {
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

// GenIndexYaml 生成首页数据文件index.yaml
func (self *MysqlRepo) GenIndexYaml() {
	// 首页最多显示20篇文章
	var posts []*model.Post
	rows, err := self.selectArticleIndex.Query()
	if err != nil {
		log.Fatalf("query index error:%s", err)
	}
	for rows.Next() {
		info := articleInfo{}
		err = rows.Scan(&info.Id, &info.Title, &info.PubTime, &info.Content)
		if err != nil {
			log.Println("scan error", err)
		}
		// post.Content, err = replaceCodeParts(blackfriday.MarkdownCommon([]byte(post.Content)))
		posts = append(posts, self.genOnePost(info))
	}

	buf, err := yaml.Marshal(posts)
	if err != nil {
		log.Printf("gen index yaml error: %v\n", err)
		return
	}
	indexYaml := global.App.ProjectRoot + PostDir + IndexFile
	ioutil.WriteFile(indexYaml, buf, 0777)
}

func (self *MysqlRepo) parsePubTime(pubTime int64) string {
	t := time.Unix(pubTime, 0).In(time.Local)
	return t.Format("2006-01-02 15:04:05")
}

// GenArchiveYaml 生成归档数据文件archive.yaml
func (self *MysqlRepo) GenArchiveYaml() {
	var posts []*model.Post
	rows, err := self.selectArticleArchives.Query()
	for rows.Next() {
		info := articleInfo{}
		err = rows.Scan(&info.Id, &info.Title, &info.PubTime)
		if err != nil {
			log.Println("query error", err)
		}
		posts = append(posts, self.genOnePost(info))
	}

	yearArchiveMap := make(map[int]*model.YearArchive)

	for _, post := range posts {

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
func (self *MysqlRepo) GenTagsYaml() {
	tagMap := make(map[string][]*model.Post)
	tagRows, err := self.selectTag.Query()
	if err != nil {
		log.Fatalf("query tag error:%s", err)
	}
	for tagRows.Next() {
		info := tagInfo{}
		err = tagRows.Scan(&info.Id, &info.Name)
		if err != nil {
			log.Println("scan error", err)
		}
		articleRows, err := self.selectArticlesByTag.Query(info.Id)
		if err != nil {
			log.Fatalf("query tag articles error:%s", err)
		}
		for articleRows.Next() {
			article := articleInfo{}
			err = articleRows.Scan(&article.Id, &article.Title, &article.PubTime)
			if err != nil {
				log.Println("query error", err)
			}
			tagMap[info.Name] = append(tagMap[info.Name], self.genOnePost(article))
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

// genOnePost 组装一个post
func (self *MysqlRepo) genOnePost(info articleInfo) *model.Post {
	return &model.Post{
		Content: info.Content,
		Meta: &model.Meta{
			Title:    info.Title,
			Path:     fmt.Sprintf("%d.html", info.Id),
			PubTime:  self.parsePubTime(info.PubTime),
			PostTime: time.Unix(info.PubTime, 0).In(time.Local),
		},
	}
}

// GenFriendsYaml 生成友情链接数据文件friends.yaml
func (self *MysqlRepo) GenFriendsYaml() {
	rows, err := self.selectFriends.Query()
	if err != nil {
		log.Fatalf("query friend error:%s", err)
	}
	var friends []*model.Friend
	for rows.Next() {
		info := friendInfo{}
		err = rows.Scan(&info.Id, &info.Name, &info.Link, &info.Logo)
		if err != nil {
			log.Println("scan error", err)
		}
		// post.Content, err = replaceCodeParts(blackfriday.MarkdownCommon([]byte(post.Content)))
		friends = append(friends, &model.Friend{Name: info.Name, Link: info.Link, Logo: info.Logo})
	}
	buf, err := yaml.Marshal(friends)
	if err != nil {
		log.Printf("gen friends yaml error:%v\n", err)
		return
	}
	friendsYaml := global.App.ProjectRoot + PostDir + FriendFile
	ioutil.WriteFile(friendsYaml, buf, 0777)
}

// UpdateDataSource 更新mysql数据
func (self *MysqlRepo) UpdateDataSource() {
	// 检查文章目录(data/post/)是否存在，不存在则连接mysql生成
	mysqlRepoDir := global.App.ProjectRoot + PostDir
	if !util.Exist(mysqlRepoDir) {
		if err := os.MkdirAll(mysqlRepoDir, os.ModePerm); err != nil {
			panic(err)
		}
	}
	// 解析仓库文件，生成首页、归档、标签数据
	self.GenIndexYaml()
	self.GenArchiveYaml()
	self.GenTagsYaml()
	self.GenFriendsYaml()

	// 定时每天自动更新仓库，并生成首页、归档、标签数据
	c := cron.New()
	c.AddFunc("@daily", func() {
		self.GenIndexYaml()
		self.GenArchiveYaml()
		self.GenTagsYaml()
		self.GenFriendsYaml()
	})
	c.Start()
}

// GetFriends 友情链接
func (self *MysqlRepo) GetFriends() ([]*model.Friend, error) {
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

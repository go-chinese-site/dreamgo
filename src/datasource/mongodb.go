package datasource

import (
	"config"
	"log"
	"model"
	"net/http"
	"sort"
	"time"

	"github.com/russross/blackfriday"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// MongoDB 数据源结构体
type MongoDB struct {
	session *mgo.Session
	addr    string
	db      string
}

// NewMongoDB 创建MongoDB数据源实例，相当于构造方法
func NewMongoDB() *MongoDB {
	addr := config.YamlConfig.Get("datasource.monogdbaddr").String()
	db := config.YamlConfig.Get("datasource.monogdbdb").String()
	if len(addr) <= 0 || len(addr) <= 0 {
		log.Fatalf("get mongodb addr or db failed addr [%s] db[%s]\n", addr, db)
	}
	mongoDBDialInfo := &mgo.DialInfo{
		Addrs:     []string{addr},
		Timeout:   10 * time.Second,
		PoolLimit: 4096,
	}
	session, err := mgo.DialWithInfo(mongoDBDialInfo)
	if err != nil {
		log.Fatalf("dial mongodb failed err:%s\n", err)
	}
	session.SetMode(mgo.Monotonic, true)
	return &MongoDB{session: session, addr: addr, db: db}
}

func (self *MongoDB) sessionclone() *mgo.Session {
	if self.session == nil {
		var err error
		mongoDBDialInfo := &mgo.DialInfo{
			Addrs:     []string{self.addr},
			Timeout:   10 * time.Second,
			PoolLimit: 4096,
		}
		self.session, err = mgo.DialWithInfo(mongoDBDialInfo)
		if err != nil {
			log.Fatalf("err:%s", err)
		}
		self.session.SetMode(mgo.Monotonic, true)
	}
	return self.session.Clone()
}

// PostList 读取文章列表
func (self MongoDB) PostList() []*model.Post {

	s := self.sessionclone()
	defer s.Close()
	posts := make([]*model.Post, 0)
	c := s.DB(self.db).C("index")
	// 根据meta.pubtime 逆序排序并 取出20个
	err := c.Find(nil).Sort("-meta.pubtime").Limit(20).All(&posts)
	if err != nil {
		log.Printf("get list failed from mongodb err: %s\n", err)
		return nil
	}
	return posts
}

// PostArchive 归档
func (self MongoDB) PostArchive() []*model.YearArchive {
	// 目前先从mongodb中将所有的文章都取出来 在进行处理
	s := self.sessionclone()
	defer s.Close()
	posts := make([]*model.Post, 0)
	c := s.DB(self.db).C("index")
	// 根据meta.pubtime 逆序排序并 取出20个
	err := c.Find(nil).Sort("-meta.pubtime").All(&posts)
	if err != nil {
		log.Printf("get list failed from mongodb err: %s\n", err)
		return nil
	}

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

	return yearArchives
}

// ServeMarkdown 处理Markdown
func (self MongoDB) ServeMarkdown(w http.ResponseWriter, r *http.Request, filename string) {
	//TODO
	//	http.ServeFile(w, r, global.App.ProjectRoot+PostDir+util.Filename(filename)+"/post.md")
}

// FindPost 根据路径查找文章
func (self MongoDB) FindPost(path string) (*model.Post, error) {
	var post *model.Post
	s := self.sessionclone()
	defer s.Close()

	c := s.DB(self.db).C("index")

	err := c.Find(bson.M{"meta.path": path + ".html"}).One(&post)
	if err != nil {
		log.Printf("Find post failed from mongodb err:%s\n", err)
		return post, err
	}
	post.Content, err = replaceCodeParts(blackfriday.MarkdownCommon([]byte(post.Content)))
	return post, err
}

// TagList 标签列表
func (self MongoDB) TagList() []*model.Tag {
	// 目前先从mongodb中将所有的文章都取出来 在进行处理
	s := self.sessionclone()
	defer s.Close()
	allPosts := make([]*model.Post, 0)
	c := s.DB(self.db).C("index")
	// 根据meta.pubtime 逆序排序并 取出20个
	err := c.Find(nil).Sort("-meta.pubtime").All(&allPosts)
	if err != nil {
		log.Printf("get list failed from mongodb err: %s\n", err)
		return nil
	}
	tagMap := make(map[string][]*model.Post)
	//遍历所有文章对象，分析出标签数据
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
	//组装标签列表
	tags := make([]*model.Tag, 0)
	for tag, posts := range tagMap {
		sort.Slice(posts, func(i, j int) bool {
			return posts[i].PubTime > posts[j].PubTime
		})
		tags = append(tags, &model.Tag{Name: tag, Posts: posts})
	}
	//按文件数量倒序排序
	sort.Slice(tags, func(i, j int) bool {
		return len(tags[i].Posts) > len(tags[j].Posts)
	})

	return tags
}

// FindTag 查找标签
func (self MongoDB) FindTag(tagName string) *model.Tag {
	tags := self.TagList()
	for _, tag := range tags {
		if tag.Name == tagName {
			return tag
		}
	}
	return nil
}

// AboutPost 关于
func (self MongoDB) AboutPost() (*model.Post, error) {
	var meta = &model.Meta{}
	post := &model.Post{
		Content: string(""),
		Meta:    meta,
	}
	return post, nil
}

// UpdateDataSource 更新数据
func (self MongoDB) UpdateDataSource() {
}

// GetFriends 友情链接
func (self MongoDB) GetFriends() ([]*model.Friend, error) {
	var friends = []*model.Friend{
		{Name: "go语言中文网", Link: "https://studygolang.com"},
	}

	return friends, nil
}

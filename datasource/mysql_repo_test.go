package datasource_test

import (
	"datasource"
	"global"
	"os"
	"strings"
	"testing"
)

var DefaultMysql *datasource.MysqlRepo

func Init() {
	cwd, _ := os.Getwd()
	pos := strings.LastIndex(cwd, "src")
	global.App.ProjectRoot = cwd[:pos]
	DefaultMysql = datasource.NewMysql("dreamgo:123456@tcp(127.0.0.1:3306)/dreamgo")
}

func TestGenMysqlIndexYaml(t *testing.T) {
	Init()

	DefaultMysql.GenIndexYaml()
}

func TestGenMysqlArchiveYaml(t *testing.T) {
	Init()

	DefaultMysql.GenArchiveYaml()
}

func TestGenMysqlTagsYaml(t *testing.T) {
	Init()

	DefaultMysql.GenTagsYaml()
}

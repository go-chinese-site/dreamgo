package datasource_test

import (
	"datasource"
	"global"
	"os"
	"strings"
	"testing"
)

func Init() {
	cwd, _ := os.Getwd()
	pos := strings.LastIndex(cwd, "src")
	global.App.ProjectRoot = cwd[:pos]
}

func TestGenMysqlIndexYaml(t *testing.T) {
	Init()

	datasource.DefaultMysql.GenIndexYaml()
}

func TestGenMysqlArchiveYaml(t *testing.T) {
	Init()

	datasource.DefaultMysql.GenArchiveYaml()
}

func TestGenMysqlTagsYaml(t *testing.T) {
	Init()

	datasource.DefaultMysql.GenTagsYaml()
}

// Copyright 2017 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// https://studygolang.com
// Author: polaris	polaris@studygolang.com

package config

import (
	"github.com/go-chinese-site/cfg"
)

// YamlConfig stores the config content
var YamlConfig *cfg.YamlConfig

// Parse parses the configFile into YamlConfig
func Parse(configFile string) {
	var err error
	YamlConfig, err = cfg.ParseYaml(configFile)
	if err != nil {
		panic(err)
	}
}

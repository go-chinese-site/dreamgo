// Copyright 2017 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// https://studygolang.com
// Author: polaris	polaris@studygolang.com

package controller

// RegisterRoutes 注册路由
func RegisterRoutes() {
	new(PostController).RegisterRoute()    // 注册文章相关路由
	new(ArchiveController).RegisterRoute() // 注册归档相关路由
	new(IndexController).RegisterRoute()   // 注册首页相关路由
	new(TagController).RegisterRoute()     // 注册标签相关路由
	new(AboutController).RegisterRoutes()  // 注册关于页面路由
	new(StaticController).RegisterRoutes() // 注册静态文件路由
}

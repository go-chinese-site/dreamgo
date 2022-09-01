package controller

import (
	"github.com/go-chinese-site/dreamgo/datasource"
	"github.com/go-chinese-site/dreamgo/logger"
	"github.com/go-chinese-site/dreamgo/route"
	"github.com/go-chinese-site/dreamgo/view"
	"net/http"
)

type FriendsController struct{}

func (self FriendsController) RegisterRoutes() {
	route.HandleFunc("/friends", self.Detail)
}

func (FriendsController) Detail(w http.ResponseWriter, r *http.Request) {
	friends, err := datasource.DefaultDataSourcer.GetFriends()
	if err == nil {
		view.Render(w, r, "friends.html", map[string]interface{}{"friends": friends})
	} else {
		logger.Instance().Error("get friends.yaml error " + err.Error())
		http.NotFound(w, r)
	}
}

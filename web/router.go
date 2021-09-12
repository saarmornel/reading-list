package web

import (
	"github.com/julienschmidt/httprouter"
)

func SetupRoutes(router *httprouter.Router) {
	router.POST("/api/signin", Signin)
	router.POST("/api/signup", Signup)
	router.POST("/api/logout", Logout)
	router.DELETE("/api/bookmarks/:url", Auth(deleteBookmark))
	router.GET("/api/user", Auth(GetMyUser))
	router.GET("/api/bookmarks", Auth(GetMyBookmarks))
	router.POST("/api/bookmarks", Auth(CreateBookmark))
	router.GET("/api/user/:username", GetUser)
	router.GET("/api/bookmarks/:username", GetBookmarks)

}

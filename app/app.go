package app

import (
	"github.com/gin-gonic/gin"
	"myblog-api/app/controller/admin"
	"myblog-api/app/controller/front"
	"myblog-api/app/db/mysql"
	"myblog-api/app/db/redis"
	"myblog-api/app/loger"
	"myblog-api/app/middleware"
	"myblog-api/app/protocol"
	"net/http"
)

type App struct {
	engine *gin.Engine
}

func Default() *App {
	app := &App{}
	app.engine = gin.Default()
	return app
}

func (this *App) Run() {
	loger.Default()
	redis.Default()
	mysql.Default()
	this.SetAccessLog()
	this.SetCors()
	this.setFront()
	this.setAdmin()
	this.set404()
}

func (this *App) GetEngin() *gin.Engine {
	return this.engine
}
func (this *App) SetAccessLog() {
	this.engine.Use(middleware.AddAccessLog())
}
func (this *App) SetCors() {
	this.engine.Use(middleware.Cors())
}

func (this *App) setFront() {
	article_front_ctrl := front.Articles{}
	this.engine.GET("/articles", article_front_ctrl.GetList)
	this.engine.GET("/categories", article_front_ctrl.GetCategories)
	this.engine.GET("/article/:id", article_front_ctrl.Show)
}

func (this *App) setAdmin() {
	//后台管理
	login_admin_ctrl := admin.Login{}
	article_admin_ctrl := admin.Articles{}
	user_admin_ctrl := admin.User{}
	this.engine.POST("/adapi/login", login_admin_ctrl.Login)
	authorized := this.engine.Group("/adapi")
	authorized.Use(middleware.CheckAuth())
	{
		authorized.POST("/article", article_admin_ctrl.Add)
		authorized.PUT("/article/:id", article_admin_ctrl.Update)
		authorized.GET("/articles", article_admin_ctrl.GetList)
		authorized.GET("/categories", article_admin_ctrl.GetCategories)
		authorized.DELETE("/article/:id", article_admin_ctrl.Delete)
		authorized.GET("/article/:id", article_admin_ctrl.Show)
		authorized.POST("/article/cache", article_admin_ctrl.DeleteCache)
		authorized.GET("/user", user_admin_ctrl.Show)
	}
}

func (this *App) set404() {
	this.engine.NoRoute(func(context *gin.Context) {
		resp := protocol.Resp{Ret: 404, Msg: "page not exists!", Data: ""}
		//返回404状态码
		context.JSON(http.StatusNotFound, resp)
	})
}

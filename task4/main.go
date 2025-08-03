package main

import (
	"task4/blog"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	db, err := gorm.Open(mysql.Open("root:123456@tcp(127.0.0.1:3306)/gorm?charset=utf8mb4&parseTime=True&loc=Local"))
	if err != nil {
		panic(err)
	}
	blog.Run(db)

	router := gin.Default()
	userGroupRouter := router.Group("/user")
	//注册
	userGroupRouter.POST("/register", func(c *gin.Context) {
		blog.Register(c, db)
	})
	//登录
	userGroupRouter.POST("/login", func(c *gin.Context) {
		blog.Login(c, db)
	})
	//帖子路由组
	postProtectedGroupRouter := router.Group("/post")
	postProtectedGroupRouter.Use(blog.Authenticate()).POST("/create", func(c *gin.Context) {
		blog.CreatePost(c, db)
	})
	postProtectedGroupRouter.PUT("/update", func(c *gin.Context) {
		blog.UpdatePost(c, db)
	})
	postProtectedGroupRouter.DELETE("/:id", func(c *gin.Context) {
		blog.DeletePost(c, db)
	})
	postpublicGroupRouter := router.Group("/post")
	postpublicGroupRouter.GET("/list", func(c *gin.Context) {
		blog.GetPosts(c, db)
	})

	//评论路由组
	commentProtectedGroupRouter := router.Group("/comment").Use(blog.Authenticate())
	commentProtectedGroupRouter.POST("/create", func(c *gin.Context) {
		blog.CreateComment(c, db)
	})
	commentPublicGroupRouter := router.Group("/comment")
	commentPublicGroupRouter.GET("/list/:postId", func(c *gin.Context) {
		blog.GetComments(c, db)
	})

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	router.Run() // 监听并在 0.0.0.0:8080 上启动服务
}

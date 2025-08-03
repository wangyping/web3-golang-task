package blog

// 项目初始化
// 创建一个新的 Go 项目，使用 go mod init 初始化项目依赖管理。
// 安装必要的库，如 Gin 框架、GORM 以及数据库驱动（如 MySQL 或 SQLite）。
// 数据库设计与模型定义
// 设计数据库表结构，至少包含以下几个表：
// users 表：存储用户信息，包括 id 、 username 、 password 、 email 等字段。
// posts 表：存储博客文章信息，包括 id 、 title 、 content 、 user_id （关联 users 表的 id ）、 created_at 、 updated_at 等字段。
// comments 表：存储文章评论信息，包括 id 、 content 、 user_id （关联 users 表的 id ）、 post_id （关联 posts 表的 id ）、 created_at 等字段。
// 使用 GORM 定义对应的 Go 模型结构体。

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"io"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	gorm.Model
	Username string `gorm:"type:varchar(50);uniqueIndex;not null"`
	Password string `gorm:"type:varchar(255);not null"`
	Email    string `gorm:"type:varchar(100);not null"`

	// 关联关系
	Posts    []Post    `gorm:"foreignKey:UserID"`
	Comments []Comment `gorm:"foreignKey:UserID"`
}

// Post 博客文章模型
type Post struct {
	gorm.Model
	Title   string `gorm:"type:varchar(255);not null;" binding:"required"`
	Content string `gorm:"type:text" binding:"required"`
	UserID  uint   `gorm:"not null"`

	// 关联关系
	Comments []Comment `gorm:"foreignKey:PostID"`
}

// Comment 评论模型
type Comment struct {
	gorm.Model
	Content string `gorm:"type:text;not null" binding:"required"`
	UserID  uint   `gorm:"not null" binding:"required"`
	PostID  uint   `gorm:"not null" binding:"required"`
}

// 实际项目中应从环境变量读取
var jwtSecret = []byte("test-jwt-secret")

// GenerateToken 生成 JWT token
func GenerateToken(user *User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(time.Minute * 10).Unix(),
	})
	return token.SignedString(jwtSecret)
}

// ParseToken 解析并验证 JWT token
func ParseToken(tokenString string) (uint, string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		return 0, "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// 验证并转换 id 和 username
		idFloat, idOk := claims["id"].(float64)
		username, nameOk := claims["username"].(string)
		if !idOk || !nameOk {
			return 0, "", fmt.Errorf("invalid claims")
		}

		// 验证 token 有效期为 10 分钟
		timestamp, expOk := claims["exp"].(float64)
		if !expOk {
			return 0, "", fmt.Errorf("missing exp claim")
		}

		expTime := time.Unix(int64(timestamp), 0)
		now := time.Now()

		// 检查是否过期或超过 10 分钟有效期
		if now.After(expTime) {
			return 0, "", fmt.Errorf("token expired")
		}

		// 检查是否在 10 分钟有效期内 (允许一定误差)
		if expTime.Sub(now) > 10*time.Minute {
			return 0, "", fmt.Errorf("token validity exceeds 10 minutes")
		}

		return uint(idFloat), username, nil

	}
	return 0, "", fmt.Errorf("invalid token")
}

func Run(db *gorm.DB) {
	r := gin.Default()

	// 添加数据库注入中间件
	r.Use(func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	})
	// db.AutoMigrate(&User{}, &Post{}, &Comment{})
	// password, err := bcrypt.GenerateFromPassword([]byte("test1234"), bcrypt.DefaultCost)
	// if err != nil {
	// 	panic(err)
	// }

	// user := User{Username: "Alice", Password: string(password)}
	// db.Debug().Create(&user)
	// post := Post{Title: "Hello World", UserID: user.ID}
	// db.Debug().Create(&post)
	// comment1 := Comment{Content: "Nice post!", PostID: post.ID, UserID: user.ID}
	// db.Debug().Create(&comment1)
	// comment2 := Comment{Content: "Second post!", PostID: post.ID, UserID: user.ID}
	// db.Debug().Create(&comment2)
	// comment3 := Comment{Content: "Recommend post", PostID: post.ID, UserID: user.ID}
	// db.Debug().Create(&comment3)

	// user2 := User{Username: "Bob", Password: string(password)}
	// db.Debug().Create(&user2)
	// post2 := Post{Title: "Hello everybody", UserID: user2.ID}
	// db.Debug().Create(&post2)
	// comment4 := Comment{Content: "First Comment!", PostID: post2.ID, UserID: user2.ID}
	// db.Debug().Create(&comment4)
	// comment5 := Comment{Content: "Second Comment!", PostID: post2.ID, UserID: user2.ID}
	// db.Debug().Create(&comment5)

}

// 用户认证与授权
// 实现用户注册和登录功能，用户注册时需要对密码进行加密存储，登录时验证用户输入的用户名和密码。
// 使用 JWT（JSON Web Token）实现用户认证和授权，用户登录成功后返回一个 JWT，后续的需要认证的接口需要验证该 JWT 的有效性。

func Register(c *gin.Context, db *gorm.DB) {
	var user User
	// 绑定JSON数据到user对象中
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	user.Password = string(hashedPassword)
	if err := db.Debug().Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "User registered successfully"})
}

func Login(c *gin.Context, db *gorm.DB) {
	var user User
	// 绑定JSON数据到user对象中
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var dbUser User
	if err := db.Where("username = ?", user.Username).First(&dbUser).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// 验证密码
	// 第一个参数是已经加密过的密码
	// 第二个参数是明文密码
	if err := bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}
	// 生成jwt token
	token, err := GenerateToken(&dbUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}
	// 密码验证成功并返回token
	c.JSON(http.StatusOK, gin.H{"code": 200, "token": token})

}

// 认证jwt的有效性
func ValidJwt(tokenString string) (bool, *User, error) {
	id, username, err := ParseToken(tokenString)
	if err != nil {
		return false, nil, err
	}
	user := User{Username: username}
	user.ID = id
	return id > 0 && username != "", &user, errors.New("invalid token")

}

// 认证jwt的有效性
func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		authorization := c.GetHeader("Authorization")
		if authorization == "" {
			c.JSON(401, gin.H{
				"message": "请先登录",
			})
			c.Abort()
			return
		}
		valid, usr, err := ValidJwt(authorization)
		if valid {
			c.Set("user", usr)
			c.Next()
		} else {
			log.Println("authenticate error:", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "Unauthorized", "error": err})
		}
	}
}

// 实现文章的创建功能，只有已认证的用户才能创建文章，创建文章时需要提供文章的标题和内容。
func CreatePost(c *gin.Context, db *gorm.DB) {
	post := Post{}
	if err := c.ShouldBindJSON(&post); err != nil {
		log.Printf("创建文章参数绑定失败: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误",
			"error":   err.Error(),
		})
		return
	}
	user, ok := c.Get("user")
	if user == nil || !ok {
		log.Println("创建文章时无法获取用户信息")
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "用户未认证",
		})
		return
	}
	post.UserID = user.(*User).ID
	if err := db.Debug().Create(&post).Error; err != nil {
		log.Printf("创建文章失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "创建文章失败",
			"error":   err.Error(),
		})
		return
	}
	log.Printf("用户 %d 创建了文章 %d", user.(*User).ID, post.ID)
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "文章创建成功",
		"data":    post,
	})

}

// 实现文章的读取功能，支持获取所有文章列表和单个文章的详细信息。
func GetPosts(c *gin.Context, db *gorm.DB) {
	posts := []Post{}
	if err := db.Debug().Preload("Comments").Find(&posts).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "查询失败", "error": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "查询成功", "data": posts})
}

// 实现文章的更新功能，只有文章的作者才能更新自己的文章。
func UpdatePost(c *gin.Context, db *gorm.DB) {
	post := Post{}
	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "invalid request"})
		return
	}
	if post.ID < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "invalid post id"})
		return
	}
	user, ok := c.Get("user")
	if user == nil || !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "cannot get user"})
		return
	}
	dbPost := Post{}
	if err := db.Debug().Where("id = ?", post.ID).First(&dbPost).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "查询失败", "error": err})
		return
	}
	if dbPost.UserID != user.(*User).ID {
		fmt.Println(user.(*User).ID)
		fmt.Println(post.UserID)
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "不是文章作者，没有权限更新"})
		return
	}
	if err := db.Debug().Model(&post).Updates(map[string]interface{}{
		"title":   post.Title,
		"content": post.Content,
	}).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "更新文章失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "更新文章成功"})
}

// 实现文章的删除功能，只有文章的作者才能删除自己的文章。
func DeletePost(c *gin.Context, db *gorm.DB) {
	var post Post
	dbPost := Post{}
	postId := c.Param("id")
	if err := db.Debug().Where("id = ?", postId).First(&dbPost).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "查询失败", "error": err})
		return
	}
	user, ok := c.Get("user")
	if user == nil || !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "cannot get user"})
		return
	}
	if dbPost.UserID != user.(*User).ID {
		fmt.Println(user.(*User).ID)
		fmt.Println(post.UserID)
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "不是文章作者，没有权限删除"})
		return
	}
	if err := db.Delete(&dbPost).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "删除失败", "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "删除成功"})
}

// 评论功能
// 实现评论的创建功能，已认证的用户可以对文章发表评论。
func CreateComment(c *gin.Context, db *gorm.DB) {
	var comment Comment
	if err := c.ShouldBindJSON(&comment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误", "error": err.Error()})
		return
	}
	if err := db.Debug().Where("id = ?", comment.PostID).First(&Post{}).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "文章不存在"})
		return
	}

	if err := db.Debug().Create(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建评论失败", "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "创建评论成功", "data": comment})
}

// 实现评论的读取功能，支持获取某篇文章的所有评论列表。
func GetComments(c *gin.Context, db *gorm.DB) {
	var post *Post
	if err := db.Debug().Preload("Comments").Where("id = ?", c.Param("postId")).Find(&post).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "文章不存在",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "获取评论列表成功", "data": post.Comments})
}

// 错误处理与日志记录
// 对可能出现的错误进行统一处理，如数据库连接错误、用户认证失败、文章或评论不存在等，返回合适的 HTTP 状态码和错误信息。
// 使用日志库记录系统的运行信息和错误信息，方便后续的调试和维护。
// ErrorHandler 统一错误处理中间件
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// 处理发生的错误
		if len(c.Errors) > 0 {
			err := c.Errors.Last()
			log.Printf("请求错误: %v", err.Err)

			// 根据错误类型返回不同的HTTP状态码
			switch err.Type {
			case gin.ErrorTypePrivate:
				c.JSON(http.StatusInternalServerError, gin.H{
					"code":    500,
					"message": "服务器内部错误",
					"error":   "内部系统错误",
				})
			case gin.ErrorTypePublic:
				c.JSON(http.StatusBadRequest, gin.H{
					"code":    400,
					"message": "请求参数错误",
					"error":   err.Error(),
				})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{
					"code":    500,
					"message": "未知错误",
				})
			}
			return
		}
	}
}

// BusinessError 自定义业务错误类型
type BusinessError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Err     error  `json:"error,omitempty"`
}

func (e *BusinessError) Error() string {
	return e.Message
}

// 常见错误类型定义
var (
	ErrUserNotFound = &BusinessError{
		Code:    404,
		Message: "用户不存在",
	}

	ErrPostNotFound = &BusinessError{
		Code:    404,
		Message: "文章不存在",
	}

	ErrUnauthorized = &BusinessError{
		Code:    401,
		Message: "未授权访问",
	}

	ErrForbidden = &BusinessError{
		Code:    403,
		Message: "没有权限执行此操作",
	}
)

func setupLogger() {
	// 创建日志文件
	logFile, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("无法创建日志文件:", err)
	}

	// 将日志同时输出到文件和控制台
	gin.DefaultWriter = io.MultiWriter(logFile, os.Stdout)

	// 设置日志格式
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// 中间件
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从 cookie 中获取会话状态
		_, err := c.Cookie("session")
		if err != nil {
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}
		c.Next()
	}
}

// 处理函数
func ShowLoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", gin.H{})
}

func Login(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	if username == os.Getenv("USERNAME") && password == os.Getenv("PASSWORD") {
		// 设置 cookie
		c.SetCookie("session", "true", 3600, "/", "", false, true)
		c.Redirect(http.StatusFound, "/")
		return
	}

	c.HTML(http.StatusUnauthorized, "login.html", gin.H{
		"error": "账号或密码错误",
	})
}

func Logout(c *gin.Context) {
	// 删除 cookie
	c.SetCookie("session", "", -1, "/", "", false, true)
	c.Redirect(http.StatusFound, "/login")
}

func getFiles() ([]string, error) {
	var files []string
	err := filepath.Walk("./storage", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			files = append(files, filepath.Base(path))
		}
		return nil
	})
	return files, err
}

func ShowFilePage(c *gin.Context) {
	files, err := getFiles()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "index.html", gin.H{
			"error": "读取失败",
		})
		return
	}

	c.HTML(http.StatusOK, "index.html", gin.H{
		"files": files,
	})
}

func AddFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.HTML(http.StatusBadRequest, "index.html", gin.H{
			"error": "未选择文件",
		})
		return
	}

	// 保存文件
	dst := filepath.Join("./storage", file.Filename)
	if err := c.SaveUploadedFile(file, dst); err != nil {
		c.HTML(http.StatusInternalServerError, "index.html", gin.H{
			"error": "保存失败",
		})
		return
	}

	c.Redirect(http.StatusFound, "/")
}

func main() {
	// 加载 .env 文件
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	r := gin.Default()

	// 加载模板
	r.LoadHTMLGlob("templates/*")

	// 静态文件服务
	r.Static("/f", "./storage")

	// 设置上传文件大小限制 (默认 32 MB)
	r.MaxMultipartMemory = 8 << 20 // 8 MB

	// 登录相关路由
	r.GET("/login", ShowLoginPage)
	r.POST("/login", Login)
	r.GET("/logout", Logout)

	// 需要认证的路由
	authorized := r.Group("/")
	authorized.Use(AuthRequired())
	{
		authorized.GET("/", ShowFilePage)
		authorized.POST("/add", AddFile)
	}

	r.Run(":8080")
}

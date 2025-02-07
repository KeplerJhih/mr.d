package handlers

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func ShowLoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", gin.H{})
}

func Login(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	if username == os.Getenv("USERNAME") && password == os.Getenv("PASSWORD") {
		c.Set("session", true)
		c.Redirect(http.StatusFound, "/")
		return
	}

	c.HTML(http.StatusUnauthorized, "login.html", gin.H{
		"error": "Invalid credentials",
	})
}

func Logout(c *gin.Context) {
	c.Set("session", false)
	c.Redirect(http.StatusFound, "/login")
}

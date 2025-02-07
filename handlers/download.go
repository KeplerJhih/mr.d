package handlers

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func ShowDownloadPage(c *gin.Context) {
	files, err := getDownloadableFiles()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "download.html", gin.H{
			"error": "Error reading files",
		})
		return
	}

	c.HTML(http.StatusOK, "download.html", gin.H{
		"files": files,
	})
}

func getDownloadableFiles() ([]string, error) {
	var files []string
	err := filepath.Walk("./downloads", func(path string, info os.FileInfo, err error) error {
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

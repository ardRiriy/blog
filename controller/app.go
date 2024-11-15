package controller

import (
	"blog_server/cache"
	"blog_server/db"
	"blog_server/model"
	"bytes"
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"os/exec"
	"strings"

	"github.com/gin-gonic/gin"
)

func GetArticleFromName(c *gin.Context) {
	name := c.Params.ByName("name")

	if val, ok := cache.Cache.Get(name); ok {
		// 完成形のHTMLを直接返す
		c.Data(http.StatusOK, "text/html; charset=utf-8", val.([]byte))
		return
	}

	article := new(model.Article)
	err := db.DB.NewSelect().
		Model(article).
		Where("name = ?", name).
		Limit(1).
		Scan(c.Request.Context())
	if err != nil {
		if err == sql.ErrNoRows {
			// 404
			// TODO: ちゃんと404用のページを作って返す
			c.JSON(http.StatusNotFound, gin.H{
				"error": fmt.Sprintf("Article {} not found", name),
			})
			return
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal Server Error.",
			})
			return
		}
	}

	// 対応するarticleをparseする
	htmlContent, err := exec.Command("armp", article.FilePath).Output()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal Server Error.",
		})

		// TODO: 対応するレコードを削除する
		return
	}

	splitedPath := strings.Split(article.FilePath, "/")
	title := strings.TrimSuffix(splitedPath[len(splitedPath)-1], ".md")
	// 完成形のHTMLを生成
	var buf bytes.Buffer
	err = c.MustGet("tmpl").(*template.Template).ExecuteTemplate(&buf, "article.tmpl", gin.H{
		"title":    template.HTML(title),
		"subtitle": template.HTML(article.Subtitle),
		"date":     template.HTML(article.UpdatedAt.Format("2006/01/02")),
		"content":  template.HTML(htmlContent),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal Server Error.",
		})
		return
	}
	html := buf.Bytes()
	cache.Cache.Add(name, html)
	c.Data(http.StatusOK, "text/html; charset=utf-8", html)
}

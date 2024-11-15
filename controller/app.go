package controller

import (
	"blog_server/cache"
	"blog_server/db"
	"blog_server/model"
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"os/exec"

	"github.com/gin-gonic/gin"
)

func GetArticleFromName(c *gin.Context) {
	fmt.Println("0")
	name := c.Params.ByName("name")

	if val, ok := cache.Cache.Get(name); ok {
		c.HTML(http.StatusOK, "article.tmpl", gin.H{
			"content": template.HTML(val.([]byte)),
		})
		return
	}

	fmt.Println("1")

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

	cache.Cache.Add(name, htmlContent)

	c.HTML(http.StatusOK, "article.tmpl", gin.H{
		"content": template.HTML(htmlContent),
	})
}

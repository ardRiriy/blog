package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
)

type Article struct {
	bun.BaseModel `bun:"table:articles"`

	Name      string    `bun:",pk"` // urlSuffixを持つ(key)
	FilePath  string    `bun:"file_path,notnull"`
	Tags      []string  `bun:"type:text[]"`
	CreatedAt time.Time `bun:"created_at,notnull,default:current_timestamp"`
	UpdatedAt time.Time `bun:"updated_at,notnull,default:current_timestamp"`
}

var (
	db *bun.DB
)

func initDB() error {
	dbUsername := os.Getenv("PSQL_USERNAME")
	password := os.Getenv("PSQL_PASSWORD")
	dsn := fmt.Sprintf("postgres://%s:%s@localhost:5432/blog?sslmode=disable", dbUsername, password)
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	db = bun.NewDB(sqldb, pgdialect.New())

	// SQL文のロギングを有効にする
	db.AddQueryHook(bundebug.NewQueryHook(
		bundebug.WithVerbose(true),
	))

	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %v", err)
	}

	return nil
}

func getArticleFromName(c *gin.Context) {
	name := c.Params.ByName("name")

	article := new(Article)
	err := db.NewSelect().
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

	// titleはファイル名
	arr := strings.Split(article.Name, "/")
	title := arr[len(arr)-1]
	// 対応するarticleをparseする
	htmlContent, err := exec.Command("armp", article.FilePath).Output()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal Server Error.",
		})
		return
	}
	c.HTML(http.StatusOK, "article.tmpl", gin.H{
		"title":   title,
		"content": template.HTML(htmlContent),
	})
}

func main() {
	if err := initDB(); err != nil {
		fmt.Errorf("Failed to initialize database:")
		fmt.Errorf("%v", err)
		os.Exit(1)
	}

	defer db.Close()

	engine := gin.Default()
	engine.Static("/assets", "./assets")
	engine.LoadHTMLGlob("templates/*")

	engine.GET("/hc", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "Service working.",
		})
	})
	engine.GET("/article/:name", getArticleFromName)
	engine.Run(":3000")
}

package model

import (
	"blog_server/db"
	"bytes"
	"context"
	"html/template"
	"log"
	"strings"
	"time"

	"github.com/uptrace/bun"
)

type Article struct {
	bun.BaseModel `bun:"table:articles"`

	Name      string    `bun:",pk"` // urlSuffixを持つ(key)
	FilePath  string    `bun:"file_path,notnull"`
	Subtitle  string    `bun:"subtitle"`
	CreatedAt time.Time `bun:"created_at,notnull,default:current_timestamp"`
	UpdatedAt time.Time `bun:"updated_at,notnull,default:current_timestamp"`
}

func FetchArticles(size int) []Article {
	ctx := context.Background()

	var articles []Article

	// データベースクエリの実行
	err := db.DB.NewSelect().
		Model(&articles).
		OrderExpr("updated_at DESC").
		Limit(size).
		Scan(ctx)
	if err != nil {
		log.Printf("Failed to fetch articles: %v", err)
		return nil
	}

	return articles
}

func ToListElement(article *Article) (string, error) {
	// テンプレート文字列
	tmpl := `
<article class="card">
    <div class="card-content">
        <header>
            <h2>{{ .Title }}</h2>
        </header>
        <p>{{ .Subtitle }}</p>
    </div>
</article>
`

	// テンプレートをパース
	t, err := template.New("article").Parse(tmpl)
	if err != nil {
		return "", err
	}
	arr := strings.Split(article.FilePath, "/")
	title := strings.TrimSuffix(arr[len(arr)-1], ".md")
	// データをmapで渡す
	data := map[string]interface{}{
		"Title":    title,
		"Subtitle": article.Subtitle,
	}

	// バッファに出力を保存
	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return "", err
	}

	// 生成したHTMLを文字列として返す
	return buf.String(), nil
}

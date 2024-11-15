package main

import (
	"bufio"
	"context"
	"database/sql"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
)

type Article struct {
	bun.BaseModel `bun:"table:articles"`

	Name      string    `bun:",pk"` // urlSuffixを持つ(key)
	FilePath  string    `bun:"file_path,notnull"`
	Subtitle  string    `bun:"subtitle"`
	CreatedAt time.Time `bun:"created_at,notnull,default:current_timestamp"`
	UpdatedAt time.Time `bun:"updated_at,notnull,default:current_timestamp"`
}

func main() {
	root, ok := os.LookupEnv("KNOWLEDGES")
	if !ok {
		fmt.Println("Environment variable KNOWLEDGES must be set.")
		os.Exit(1)
	}

	articles := []Article{}

	dbUsername := os.Getenv("PSQL_USERNAME")
	password := os.Getenv("PSQL_PASSWORD")
	dsn := fmt.Sprintf("postgres://%s:%s@localhost:5432/blog?sslmode=disable", dbUsername, password)
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	db := bun.NewDB(sqldb, pgdialect.New())
	// SQL文のロギングを有効にする
	db.AddQueryHook(bundebug.NewQueryHook(
		bundebug.WithVerbose(true),
	))

	filepath.Walk(root, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			fmt.Println("Failed to access %s: %v", path, err)
			return err
		}

		if info.IsDir() {
			return nil
		}

		if filepath.Ext(path) == ".md" {
			file, err := os.Open(path)
			if err != nil {
				fmt.Println("Failed to open %s:", path)
				fmt.Println(err)
				return err
			}
			scanner := bufio.NewScanner(file)
			scanner.Scan()
			firstLine := strings.TrimSpace(scanner.Text())
			t1, b1 := strings.CutPrefix(firstLine, "<!-- url:")
			t2, b2 := strings.CutSuffix(t1, "-->")
			urlSuffix := strings.TrimSpace(t2)

			if b1 && b2 {
				// 公開対象なのでpathのurlの組をDBに保存
				// initializeではinsertする。空のDBに対して実行しているので重複しないはず。
				scanner.Scan()
				secontLine := strings.TrimSpace(scanner.Text())
				subtitle := ""
				t3, b3 := strings.CutPrefix(secontLine, "<!-- subtitle:")
				t4, b4 := strings.CutSuffix(t3, "-->")
				if b3 && b4 {
					subtitle = strings.TrimSpace(t4)

				}
				article := Article{
					Name:     urlSuffix,
					FilePath: path,
					Subtitle: subtitle,
				}
				articles = append(articles, article)
			}
		}

		return nil
	})

	_, err := db.NewInsert().Model(&articles).Exec(context.Background())
	if err != nil {
		fmt.Println("Unknown err: ")
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("All Done.")

}

package controller

import (
	"blog_server/cache"
	"blog_server/server"
	"fmt"
	"log"
	"os"
	"os/exec"
	"sync"

	"github.com/gin-gonic/gin"
)

func GithubWebhook(c *gin.Context) {
	// TODO: GitHubからのリクエストであることを確認する

	// MEMO: 可読性の観点からメンバ関数したい
	server.StartServerMaintenance(server.Sv)

	var wg sync.WaitGroup

	wg.Add(1)
	// PULLして変更があったファイルに前処理を行う
	go func() {
		defer wg.Done()
		if err := pullChanges(); err != nil {
			log.Printf("Failed to pull changes: %v", err)
			// MEMO: Discordに通知する機構を入れる
		} else {
			log.Println("Git pull completed")
		}
	}()

	// キャッシュをクリアする処理
	wg.Add(1)
	go func() {
		defer wg.Done()
		cache.Cache.Purge()
		log.Println("Cache cleared")
	}()

	wg.Wait()
	server.EndServerMaintenance(server.Sv)

	// webhookなので返却不要のはず？
}

func pullChanges() error {
	repo := os.Getenv("KNOWLEDGES")
	if repo == "" {
		return fmt.Errorf("Environment variable KNOWLEGES is not set.")
	}

	cmd := exec.Command("git", "-C", repo, "pull")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("git pull failed %v because of &s", err, string(output))
	}

	log.Printf("success git pull")
	return nil
}

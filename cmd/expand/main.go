package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func main() {
	// knowledges内に含まれる画像をassets/以下に移動させる
	// 性質的に数がそこまで多くならないのでこれで妥協
	root, ok := os.LookupEnv("KNOWLEDGES")
	if !ok {
		fmt.Println("Environment variable KNOWLEDGES must be set.")
		os.Exit(1)
	}
	sourceDir := root + "/pictures" // 指定されたディレクトリに変更
	// コピー先ディレクトリ
	destinationDir := "../../assets/pictures"

	// コピー先ディレクトリを作成
	if err := os.MkdirAll(destinationDir, 0755); err != nil {
		fmt.Println("ディレクトリ作成エラー:", err)
		return
	}

	// コピー元ディレクトリ内のファイルを処理
	err := filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		// エラー処理
		if err != nil {
			return err
		}

		// ディレクトリはスキップ
		if info.IsDir() {
			return nil
		}

		// ファイルの名前をdestinationDirにコピーする
		destPath := filepath.Join(destinationDir, info.Name())
		if err := copyFile(path, destPath); err != nil {
			return err
		}
		fmt.Printf("ファイル %s を %s にコピーしました。\n", path, destPath)
		return nil
	})

	if err != nil {
		fmt.Println("ファイルコピー中にエラーが発生しました:", err)
	}
}

// ファイルをコピーする関数
func copyFile(src, dst string) error {
	// ソースファイルを開く
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	// デスティネーションファイルを作成
	destinationFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destinationFile.Close()

	// ファイルをコピー
	_, err = io.Copy(destinationFile, sourceFile)
	return err
}

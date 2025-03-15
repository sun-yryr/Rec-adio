# recadio - Rec-adio コマンドラインツール

このディレクトリには、Rec-adioのコマンドラインツールの実装が含まれています。

## 概要

recadioは、Rec-adioのコアサービス（recadiod）と通信するためのコマンドラインインターフェースを提供します。以下の機能を提供します：

- 録音予約の追加・削除・一覧表示
- スケジュールの管理
- 録音済みファイルの管理
- 設定の変更
- プラグインの管理
- 演者とプラットフォームアカウントの管理

## 実装方法

コマンドラインツールは、以下のコンポーネントで構成されています：

1. **コマンド処理**: コマンドラインの引数を解析し、適切な処理を実行します
2. **APIクライアント**: コアサービスのAPIと通信します
3. **設定管理**: ローカルの設定ファイルを管理します
4. **出力フォーマッター**: コマンドの実行結果を適切な形式で出力します

## サンプルコード

```go
package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/username/rec-adio/pkg/client"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "recadio",
		Short: "Rec-adio CLI tool",
		Long:  `Command line interface for Rec-adio, a tool for recording radio and live streams.`,
	}

	// クライアントの初期化
	apiClient, err := client.New("http://localhost:8080")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// スケジュールコマンド
	var scheduleCmd = &cobra.Command{
		Use:   "schedule",
		Short: "Manage recording schedules",
	}

	// スケジュール一覧表示コマンド
	var listScheduleCmd = &cobra.Command{
		Use:   "list",
		Short: "List recording schedules",
		Run: func(cmd *cobra.Command, args []string) {
			schedules, err := apiClient.ListSchedules()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}

			for _, s := range schedules {
				fmt.Printf("%s: %s (%s)\n", s.ID, s.Title, s.StartTime)
			}
		},
	}

	scheduleCmd.AddCommand(listScheduleCmd)
	rootCmd.AddCommand(scheduleCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

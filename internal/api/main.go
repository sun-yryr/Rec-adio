package api

import (
	"context"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
)

var App *fiber.App

// New はAPIサーバーを初期化し、指定されたコンテキストがキャンセルされたときに
// サーバーをシャットダウンします
func New(ctx context.Context) {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	// ゴルーチンでサーバーを起動
	go func() {
		if err := app.Listen(":3000"); err != nil {
			log.Printf("APIサーバーの起動に失敗しました: %v", err)
		}
	}()

	// コンテキストのキャンセルを監視
	go func() {
		<-ctx.Done()
		log.Println("APIサーバーをシャットダウンしています...")

		// シャットダウンのタイムアウトを設定
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := app.ShutdownWithContext(shutdownCtx); err != nil {
			log.Printf("APIサーバーのシャットダウンに失敗しました: %v", err)
		} else {
			log.Println("APIサーバーをシャットダウンしました")
		}
	}()

	App = app
}

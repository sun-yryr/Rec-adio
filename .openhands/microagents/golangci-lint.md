---
name: golangci-lint-guidelines
type: knowledge
agent: CodeActAgent
triggers:
- lint
- linter
- golangci-lint
---

# golangci-lint 推奨スタイルガイド

このドキュメントは、`golangci-lint` によって指摘される一般的な問題を防ぎ、一貫性のあるコードベースを維持するためのスタイルガイドラインを提供します。

## エラー処理

-   **エラーパッケージ**: `github.com/cockroachdb/errors` パッケージを使用してエラーをラップし、スタックトレース情報を含めることを推奨します。
-   **エラー変数**: 頻繁に使用するエラーは、パッケージレベルの変数として定義し、再利用してください。
    ```go
    var (
        errTest  = errors.New("test error")
        errTest1 = errors.New("test error 1")
    )
    ```

## コンテキスト

-   **テストでのコンテキスト**: テスト関数内では、`context.Background()` の代わりに `t.Context()` を使用して、テストのライフサイクルに紐付いたコンテキストを取得してください。
    ```go
    // 良い例
    ctx := logger.WithLogger(t.Context(), testLogger)

    // 悪い例
    // ctx := logger.WithLogger(context.Background(), testLogger)
    ```
-   **ロガーの引き渡し**: コンテキスト経由でロガーを渡す場合、関数やメソッドの先頭で `logger.FromContext(ctx)` を使ってロガーを取得し、それを利用してください。ハンドラ内でロガーが正しく渡されているか（例：`assert.NotEmpty(t, log.Core())`）を確認することも有効です。

## レスポンスとアサーション

-   **Getterメソッド**: protobufなどで生成された構造体のフィールドにアクセスする場合、直接フィールドを参照するのではなく、提供されている `GetXxx()` 形式のGetterメソッドを使用してください (例: `resp.GetStatus()` のように `resp.Status` の代わりに)。
-   **アサーションの順序**: `require.NoError(t, err)` のようなクリティカルな事前条件チェックは、他のアサーションよりも先に行い、エラーが発生した場合は早期にテストを失敗させてください。

## コードスタイル

-   **変数宣言**: 複数の変数を宣言する場合は、`var` キーワードでグループ化してください。
    ```go
    // 良い例
    var (
        capturedSubject string
        capturedMessage []byte
    )

    // 悪い例
    // var capturedSubject string
    // var capturedMessage []byte
    ```
-   **ファイルパーミッション**: `os.WriteFile` などでファイルのパーミッションを指定する際は、`0600` のような10進数表記ではなく、`0o600` のような8進数表記を使用してください。
-   **未使用の関数引数**: 関数シグネチャに含まれるものの、関数内で使用しない引数がある場合は、`_` (アンダースコア) を使って明示的に無視してください。
    ```go
    // 良い例
    checkFunc: func(_ context.Context) error {
        return errTest
    }

    // 悪い例
    // checkFunc: func(ctx context.Context) error {
    //     return errTest
    // }
    ```
-   **コメント**: コメント文の末尾には、句読点（`.` や `。`）を記述してください。
    ```go
    // 良い例
    // mockBroker はブローカーのモック。

    // 悪い例
    // mockBroker はブローカーのモック
    ```
-   **構造体の初期化**: `eventutil.NewEventService` のような複数の引数を取る関数で構造体を初期化する場合、可読性向上のために引数を複数行に分けて記述することを検討してください。
    ```go
    // 良い例
    service := eventutil.NewEventService[recording.RequestedEvent](
        mockBroker,
        logger,
        "test.subject",
    )
    ```
-   **エラーのアンラップ**: `errors.Is` や `errors.As` を使用する際、エラーが複数回ラップされている可能性を考慮し、必要に応じて `errors.Unwrap` を複数回呼び出して原因となるエラーを比較してください。
    ```go
    // 例
    assert.True(t, errors.Is(errors.Unwrap(errors.Unwrap(err)), publishErr))
    ```

---
name: testify
type: knowledge
agent: CodeActAgent
triggers:
- test
- テスト
- testify
- テスト作成
---

# Testifyライブラリによるテストスタイルガイド

## 基本方針

このリポジトリのテストコードでは、`github.com/stretchr/testify` パッケージを使用して簡潔で読みやすいテストを記述します。

## 使用するアサーション

- `assert`: 一般的なアサーションに使用（テスト継続）
- `require`: 重要な条件のアサーションに使用（失敗時テスト中断）

## テストコード作成ルール

### 基本構造

```go
func TestXxx(t *testing.T) {
    t.Parallel() // 可能な限り並列実行を有効に

    // テスト準備
    // ...

    // テスト実行
    // ...

    // 結果検証
    assert.Equal(t, expected, actual)
    require.NoError(t, err)
}
```

### アサーションの使い分け

- `assert.Equal(t, expected, actual)` - 値の一致検証（テスト継続）
- `assert.True(t, condition)` - 条件が真であることを検証
- `assert.FileExists(t, path)` - ファイルの存在を検証
- `require.NoError(t, err)` - エラーがないことを検証（エラーがあれば中断）
- `require.ErrorAs(t, err, &target)` - エラーの型を検証

### テーブル駆動テスト

```go
tests := []struct {
    name     string
    input    string
    expected string
    isValid  bool
}{
    {
        name:     "valid case",
        input:    "test",
        expected: "result",
        isValid:  true,
    },
    // ...
}

for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
        t.Parallel()
        
        // テスト実行
        // ...
        
        // 結果検証 - 早期リターンを使用してネストを減らす
        if !tt.isValid {
            assert.Error(t, err)

            return
        }
        
        assert.Equal(t, tt.expected, actual)
    })
}
```

### 早期リターンの使用

ネストが深くならないよう、条件分岐では以下のように早期リターンを使用して可読性を高めます：

```go
// 悪い例（ネストが深い）
if tt.isValid {
    if err == nil {
        assert.Equal(t, tt.expected, actual)
    } else {
        t.Errorf("unexpected error: %v", err)
    }
} else {
    assert.Error(t, err)
}

// 良い例（早期リターンでネストを減らす）
if !tt.isValid {
    assert.Error(t, err)

    return
}

require.NoError(t, err)
assert.Equal(t, tt.expected, actual)
```

## 従来のテストコードからの移行方針

### if文によるエラー検証の代替

- 古い: `if err != nil { t.Fatalf("error: %v", err) }`
- 新しい: `require.NoError(t, err)`

### 値の比較の代替

- 古い: `if actual != expected { t.Errorf("expected %v, got %v", expected, actual) }`
- 新しい: `assert.Equal(t, expected, actual)`

### ファイル存在確認の代替

- 古い: `if _, err := os.Stat(path); os.IsNotExist(err) { t.Error("file not found") }`
- 新しい: `assert.FileExists(t, path)`

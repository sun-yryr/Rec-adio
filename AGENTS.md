## コマンドやパス

Swiftly を利用しています。そのため、 where swift の結果が swiftly を向いていれば正常です。

間違っても MacOS 標準の swift や XCode の swift ツールチェーンを使わないようにしてください。

### サンドボックス環境での実行

Codex / Claude など、サンドボックス付きのエージェント実行環境では、SwiftPM の既定パス (`~/Library`, `~/.cache`) と内部 sandbox が衝突しやすいです。
以下の設定で、キャッシュ類を `.build` 配下へ寄せて実行してください。

```bash
mkdir -p \
  .build/swiftpm-cache \
  .build/swiftpm-config \
  .build/swiftpm-security \
  .build/xdg-cache \
  .build/clang-module-cache \
  .build/scratch

env \
  PATH=/Users/sun-yryr/.swiftly/bin:$PATH \
  XDG_CACHE_HOME=$PWD/.build/xdg-cache \
  CLANG_MODULE_CACHE_PATH=$PWD/.build/clang-module-cache \
  swift test \
    --disable-sandbox \
    --cache-path $PWD/.build/swiftpm-cache \
    --config-path $PWD/.build/swiftpm-config \
    --security-path $PWD/.build/swiftpm-security \
    --scratch-path $PWD/.build/scratch \
    --manifest-cache local
```

`swift build` / `swift package resolve` も同様に上記オプションを付けて実行すること。

### よく使うもの

**パッケージの解決・新しい依存を追加した時など**

```
swift package resolve
```

**ビルド**

```
swift build
```

**テスト**

```
swift test
```

**Buf**

```
# ファイルの生成
buf generate

# FileDescriptorSet の生成
buf build -o Sources/App/Resources/Reflection/recoto.pb --as-file-descriptor-set
```

### その他

わからない場合は各コマンドの help を見ること。

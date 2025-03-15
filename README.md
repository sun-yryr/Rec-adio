# Rec-adio

Rec-adioは、様々なプラットフォーム（radiko、超A&G+、Twitter Space、BiliBili Liveなど）のラジオやライブ配信を自動的に録音するためのソフトウェアです。

## 特徴

- シンプルで直感的なコマンドラインインターフェース
- プラグイン式のアーキテクチャによる拡張性
- 複数のプラットフォームに対応
- 軽量で効率的な実装
- クロスプラットフォーム対応（Windows、macOS、Linux）

## サポートするプラットフォーム

- 超A&G+（手動予約、将来的に番組表からの予約も対応予定）
- Twitter Space（将来的に対応予定）
- BiliBili Live（将来的に対応予定）
- Radiko（将来的に対応予定）
- YouTube Live（将来的に対応予定）

## アーキテクチャ

Rec-adio v4は、「dockerd+docker」のような関係性を持つクライアント・サーバーアーキテクチャを採用しています：

1. **コアサービス（recadiod）**：バックグラウンドで動作し、実際の録音処理を担当
2. **コマンドラインツール（recadio）**：コアサービスとやり取りするためのインターフェース

## インストール

```bash
# 準備中
```

## 使い方

```bash
# 準備中
```

## 開発

### 必要条件

- Go 1.24以上
- Git

### セットアップ

```bash
# リポジトリのクローン
git clone https://github.com/sun-yryr/Rec-adio.git
cd Rec-adio

# 依存関係のインストール
go mod download

# ビルド
go build -o bin/recadiod ./cmd/recadiod
go build -o bin/recadio ./cmd/recadio
```

## 貢献

新しいプラットフォームの追加や機能の改善など、プロジェクトへの貢献を歓迎します。詳細は[開発者ガイド](docs/dev/README.md)を参照してください。

## ライセンス

このプロジェクトはMITライセンスの下で公開されています。

## 過去のバージョン

- [v3](v3/README.md) - Python実装

[![agqr url check](https://github.com/sun-yryr/Rec-adio/actions/workflows/agqr-check.yml/badge.svg)](https://github.com/sun-yryr/Rec-adio/actions/workflows/agqr-check.yml)
[![main-ci](https://github.com/sun-yryr/Rec-adio/actions/workflows/main-ci.yml/badge.svg)](https://github.com/sun-yryr/Rec-adio/actions/workflows/main-ci.yml)

> 録音中に予約の変更等が行えるjob方式を採用した新バージョン開発中  
https://github.com/sun-yryr/Rec-adio/tree/v4

# radio
オタク！ラジオを聞き逃すな！ってことで作るradio録音

## 録音対応プラットフォーム
- Radiko
- 超A&G+
- 音泉
- 響

## できること
- 設定ファイルの作成と読み込み
- 番組表の取得，キーワード録音
- 複数番組の同時録音
- systemctlでの自動起動，再起動
- 保存先の選択
    - DropBox
    - Object-Storage(Swift)
    - ローカル
    - rcloneによる同期
- 録音番組情報をMysqlに登録
- 録音完了をLINE Notifyで通知

# Require

- pipenv
- ffmpeg
- rtmpdump

# 実行方法(録音ツール)

## 設定
```bash
$ git clone https://github.com/sun-yryr/Rec-adio.git
$ cd Rec-adio
$ pipenv install
$ pipenv run python Setup.py
$ nano ./conf/config.json
# キーワードの設定を好みに合わせてください。正規表現なので .+ とかやると全部取れるはずです。
```

## コマンドから実行する

```bash
$ pipenv run start
```

## 起動を自動化する

`rec_adio.service` を環境に合わせて設定する。  
変更する箇所
- User ユーザー名
- WorkingDirectory `hogehoge/Rec-adio`になるように絶対パスで書く
- ExecStart `which pipenv`で出力されたパス + `run start`にする

```bash
$ sudo mv ./rec_adio.service /etc/systemd/system/
$ sudo systemctl daemon-reload
$ sudo systemctl enable rec_adio.service
$ sudo systemctl start rec_adio.service
```

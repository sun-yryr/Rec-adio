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
- 録音番組情報をMysqlに登録
- 録音完了をLINE Notifyで通知

# 設定方法
Pipenvが必要です
```
pip install pipenv --user
```

Pipenvインストール済みの人
```
sudo apt install -y ffmpeg rtmpdump
git clone https://github.com/sun-yryr/Rec-adio.git
cd Rec-adio
pipenv install
cp ./conf/example_config.json ./conf/config.json
nano ./conf/config.json
設定なので好みに合わせてください。キーワードは正規表現なので .+ とかやると全部取れるはずです。

nano ./rec_adio.service
ユーザー名，パスを環境に合わせて下さい（下記参照）

sudo mv ./rec_adio.service /etc/systemd/system/
sudo systemctl daemon-reload
sudo systemctl enable rec_adio.service
sudo systemctl start rec_adio.service
```

### 設定ファイル
- **all** 主に全体向け
    - `savedir: string` 保存先。指定がなければ `Rec-adio/savefile` が作成される。
    - `dbx_token: string` DropBoxのアクセストークン。任意
    - `line_token: string` Notifyのアクセストークン。任意(?)
    - `keywords: Array<string>` 録音するキーワード。正規表現が使える。
- **swift** オブジェクトストレージ用
    - `tenantid: string` テナントID
    - `username: string`
    - `password: string`
    - `identityUrl: string` 認証用Url。https~
    - `objectStorageUrl: string` 保存用Url。https~
- **mysql** Mysql用
    - `hostname: string`
    - `port: string` default = 3306
    - `username: string`
    - `password: string`
    - `database: string` データベース名

### Systemd serviceファイル
変更する箇所
- User ユーザー名
- WorkingDirectory `hogehoge/Rec-adio`になるように絶対パスで書く
- ExecStart `which pipenv`で出力されたパス + `run start`にする
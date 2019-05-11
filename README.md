# radio
オタク！ラジオを聞き逃すな！ってことで作るradio録音

## できること
- 設定ファイルの作成と読み込み
- 番組表の取得，キーワード録音
- 複数番組の同時録音
- systemctlでの自動起動，再起動
- 出来たらdropに保存（アクセストークンは設定ファイル）
- 録音完了を何らかの方法で通知する（line）

# 設定方法 centOS7の場合
参考ページ : 
```
git clone https://github.com/sun-yryr/Rec-adio.git
cd Rec-adio
cp ./conf/example_config.json ./conf/config.json
nano ./conf/config.json

# 設定なので適当にいじってください

sudo mv ./rec_adio.service /etc/systemd/system/
sudo systemctl daemon-reload
sudo systemctl enable rec_adio.service
sudo systemctl start rec_adio.service
```

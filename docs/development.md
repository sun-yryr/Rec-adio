# 開発ガイド

## gRPC経由でNATSイベント送信を確認する手順

この手順では、gRPCのStartFromURLメソッドを呼び出し、NATSにイベントが送信されていることを確認します。

### 1. NATSのクライアントを準備する

[natscli - GitHub](https://github.com/nats-io/natscli) が利用できます。  
READMEに従ってインストールして下さい。

```sh
$ nats context add recoto --server 127.0.0.1:4222 --select
$ nats sub 'recoto.*.*.*'
21:52:16 Subscribing on recoto.*.*.*
```

### 2. サーバーを起動する

```sh
go run ./cmd/recotod/main.go
```

### 3. grpcurlでgRPCリクエストを送信

以下のコマンドを実行し、gRPCサーバーにリクエストを送信します。

```sh
grpcurl -plaintext \
  -d '{"url": "YOUR_URL", "title": "YOUR_TITLE", "duration": 3600}' \
  localhost:8080 \
  recoto.recording.v1.RecordingService/StartFromURL
```

- `YOUR_URL`, `YOUR_TITLE`, `duration`は適宜変更してください。
- `localhost:8080`はgRPCサーバーのアドレスです。

### 4. 受信を確認

natscli側に表示されているはずです。

```sh
$ nats sub 'recoto.*.*.*'
21:52:16 Subscribing on recoto.*.*.*
[#1] Received on "recoto.recording.requested.v1"
{"recordingId":"5cfd8e3f-d0ba-4e13-89a9-5c35c060e74a","status":"requested","source":{"kind":"url","id":"YOUR_URL","meta":null},"output":"data/output/YOUR_TITLE.mp3","duration":3600,"timestamp":"2025-05-06T21:52:20.753822+09:00"}
```

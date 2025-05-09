#!/bin/bash

# TLS用ディレクトリを作成
mkdir -p tls
cd tls

# redis.key: 秘密鍵
# redis.crt: 自己署名証明書
# ca.crt: 開発用なのでredis.crtをコピーしてCAとして使用
openssl req -x509 -newkey rsa:2048 -sha256 -nodes \
  -keyout redis.key \
  -out redis.crt \
  -days 365 \
  -subj "/CN=localhost"

cp redis.crt ca.crt

# ファイルの確認メッセージ
echo "✔ TLS証明書生成完了:"
echo "  - tls/redis.crt"
echo "  - tls/redis.key"
echo "  - tls/ca.crt"

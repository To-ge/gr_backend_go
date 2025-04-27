# gr_backend_go
卒業研究のバックエンド（言語：Go）のプログラム

## 概要
私の大学では毎年一度ハイブリッドロケットの打ち上げを行なっています。
私はこのプロジェクトに関わり、「リアルタイムロケット軌道表示(RTTD)システムの開発と評価」という研究テーマに取り組みました。
ロケットに搭載した人工衛星で取得した位置情報をリアルタイムにブラウザのマップUI上に表示することができます。
このリポジトリはRTTDシステムのバックエンドを担っています。

## 使用技術

| Category          | Technology Stack             |
| ----------------- | ---------------------------- |
| Main              | Go, gRPC                     |
| Infrastructure    | Amazon Web Services          |
| Database          | PostgreSQL                   |
| Environment setup | Docker                       |
| CI/CD             | GitHub Actions               |
| etc.              | Echo, Gorm                   |

## インフラ構成図

![インフラ構成図](/docs/image/aws-architecture.png)
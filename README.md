# blog

ハッカソン用に公開中。ハッカソン終了後は非公開にし運用します。

## setup

db version: `psql (PostgreSQL) 16.3`

[golang-migrate](https://github.com/golang-migrate/migrate) version: `4.18.1`

まず、適当なユーザとデータベースを作成します。データベース名は`blog`にしてください(そうでない場合はinitializeを書き換えて対応するか、環境変数に対応させるかしてください。あなたですよ)。次に、以下のコマンドを実行するとDBがmigrateされます。
シェルスクリプト内の環境変数は適宜置き換えること。
```bash
$ migrate --path db/migrations --database "postgresql://${PSQL_USERNAME}:${PSQL_PASSWORD}@localhost:5432/${db_name}?sslmode=disable" -verbose up 
```

もし、テーブルを削除したい場合は以下のようにすればよいです。
```bash
$ migrate --path db/migrations --database "postgresql://${PSQL_USERNAME}:${PSQL_PASSWORD}@localhost:5432/${db_name}?sslmode=disable" -verbose down 
```

## initialize

## 構成
DBクライアントにはbunを使用しています。
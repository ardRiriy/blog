# blog

ハッカソン用に公開中。ハッカソン終了後は非公開にし運用します。

## setup

db version: `psql (PostgreSQL) 16.3`

[golang-migrate](https://github.com/golang-migrate/migrate) version: `4.18.1`
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

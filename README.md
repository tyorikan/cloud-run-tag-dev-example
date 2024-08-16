# Cloud Run Tag Dev Example
Cloud Run で Pull Request 毎の環境を払い出すデモ

## Launch API
```
go run cmd/api/main.go
```

## Cloud Build pipelines
1. cloudbuild_pr.yaml（no-traffic で Cloud Run デプロイ & タグ発行）  
PR が作成されたら実行

2. cloudbuild_rm_run_tag.yaml  
Branch が削除されたら実行（GitHub Actions から Cloud Build を呼び出し、タグ削除）

3. cloudbuild.yaml  
main ブランチに merge (or push) されたら実行（Cloud Run へのデプロイ via Cloud Deploy）

## Setup
[チュートリアル](tutorial.md)を参照

# Cloud Run Tag Dev Example

## Launch API
```
go run cmd/api/main.go
```

## Cloud Build pipelines
* cloudbuild.yaml  
main ブランチに push されたら実行

* cloudbuild_pr.yaml  
PR が作成されたら実行

* cloudbuild_rm_run_tag.yaml  
Branch が削除されたら実行（GitHub Actions からの呼び出し）
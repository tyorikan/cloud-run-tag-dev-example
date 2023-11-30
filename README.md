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
main ブランチに push されたら実行（Cloud Run へのデプロイ via Cloud Deploy）

## Setup

1. [Workload Identity Federation](https://github.com/google-github-actions/auth) の設定  
2. [Cloud Build トリガーの設定](#cloud-build-pipelines)  
3. GitHub Actions で利用する secret の登録  
`GCP_PROJECT_NUMBER`:  Google Cloud プロジェクト番号  
`GCP_SA_ID`: Workload Identity Federation で利用するサービスアカウント

### IAM で付与するロール
```
{ PROJECT_NUMBER }-compute@developer.gserviceaccount.com
Cloud Deploy ランナー
Cloud Run デベロッパー
サービス アカウント ユーザー

{ PROJECT_NUMBER }@cloudbuild.gserviceaccount.com
Cloud Build サービス アカウント
Cloud Deploy オペレーター
Cloud Run 管理者
サービス アカウント ユーザー
```

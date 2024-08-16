## Google Cloud Project の作成、環境変数定義
自身が Owner 権限（または Editor + IAM 管理者）を持つプロジェクトを作成
```
export PROJECT_ID={Google Cloud Project ID}
export PROJECT_NUMBER={Google Cloud Project Number}
export GITHUB_ACCOUNT={自身の GitHub アカウント}
```

## API の有効化
```bash
gcloud services enable artifactregistry.googleapis.com run.googleapis.com cloudbuild.googleapis.com clouddeploy.googleapis.com compute.googleapis.com iam.googleapis.com
```

## IAM の準備
1. Cloud Build 用のサービスアカウントの作成
```bash
gcloud iam service-accounts create cloud-build-runner 
```
2. Cloud Run 用のサービスアカウントの作成
```bash
gcloud iam service-accounts create demo-backend-api
```

### Role の付与
1. Cloud Deploy で利用するデフォルト SA
```bash
gcloud projects add-iam-policy-binding $PROJECT_ID --member serviceAccount:${PROJECT_NUMBER}-compute@developer.gserviceaccount.com --role=roles/clouddeploy.jobrunner
gcloud projects add-iam-policy-binding $PROJECT_ID --member serviceAccount:${PROJECT_NUMBER}-compute@developer.gserviceaccount.com --role=roles/clouddeploy.releaser
gcloud projects add-iam-policy-binding $PROJECT_ID --member serviceAccount:${PROJECT_NUMBER}-compute@developer.gserviceaccount.com --role=roles/run.developer
gcloud projects add-iam-policy-binding $PROJECT_ID --member serviceAccount:${PROJECT_NUMBER}-compute@developer.gserviceaccount.com --role=roles/iam.serviceAccountUser
```
2. Cloud Build で利用する　SA
```bash
gcloud projects add-iam-policy-binding $PROJECT_ID --member serviceAccount:cloud-build-runner@${PROJECT_ID}.iam.gserviceaccount.com --role=roles/cloudbuild.builds.builder
gcloud projects add-iam-policy-binding $PROJECT_ID --member serviceAccount:cloud-build-runner@${PROJECT_ID}.iam.gserviceaccount.com --role=roles/clouddeploy.operator
gcloud projects add-iam-policy-binding $PROJECT_ID --member serviceAccount:cloud-build-runner@${PROJECT_ID}.iam.gserviceaccount.com --role=roles/run.admin
gcloud projects add-iam-policy-binding $PROJECT_ID --member serviceAccount:cloud-build-runner@${PROJECT_ID}.iam.gserviceaccount.com --role=roles/iam.serviceAccountUser
```

## GitHub の準備
1. このリポジトリを自分のアカウント以下に Fork
https://github.com/tyorikan/cloud-run-tag-dev-example
2. Secrets の設定 
Settings -> Secrets and variables -> Actions 画面で、Variables タブの選択

| Name | Value |
-------|-------- 
| CLOUD_BUILD_REGION | asia-northeast1 |
| CLOUD_BUILD_TRIGGER_NAME | demo-backend-api-remove-cloud-run-tag |
| GCP_PROJECT_NUMBER | {Google Cloud Project Number} |
| GCP_SA_ID | cloud-build-runner@{GOOGLE Cloud Project ID}.iam.gserviceaccount.com |
| WORKLOAD_IDENTITY_POOL | github-actions-pool |
| WORKLOAD_IDENTITY_PROVIDER | github-actions-provider |

## Workload Idenitty 連携の準備
1. IAM -> Workload Identity 連携へ移動し、プロバイダを追加
```
ID プール名：github-actions-pool
プロバイダ：OIDC
プロバイダ名：github-actions-provider
発行元：https://token.actions.githubusercontent.com
オーディエンス：デフォルト
プロバイダ属性：
google.subject=assertion.sub
attribute.repository_owner=assertion.repository_owner
```
2. GitHub Actions から Cloud Build を呼び出すため、Cloud Build で利用する SA に対し、Workload Identity ユーザーの権限を追加
```bash
gcloud iam service-accounts add-iam-policy-binding cloud-build-runner@${PROJECT_ID}.iam.gserviceaccount.com \
    --role=roles/iam.workloadIdentityUser \
    --member="principalSet://iam.googleapis.com/projects/${PROJECT_NUMBER}/locations/global/workloadIdentityPools/github-actions-pool/attribute.repository_owner/${GITHUB_ACCOUNT}"
```

## Cloud Build の準備
1. Cloud Build -> リポジトリ -> ホスト接続を作成、で GitHub と接続
2. Cloud Build -> リポジトリ -> リポジトリをリンク、から Fork したリポジトリをリンク
3. 環境変数にセット
```bash
export GITHUB_HOST=...
export GITHUB_REPO=...
```

### トリガーの作成
1. demo-backend-api-pull-request
```bash
cat <<EOF > ./pr-trigger.yaml
description: Build and deploy to Cloud Run service demo-backend-api on pull request
filename: cloudbuild_pr.yaml
includeBuildLogs: INCLUDE_BUILD_LOGS_WITH_STATUS
name: demo-backend-api-pull-request
repositoryEventConfig:
  pullRequest:
    branch: .*
    commentControl: COMMENTS_ENABLED_FOR_EXTERNAL_CONTRIBUTORS_ONLY
  repository: projects/${PROJECT_ID}/locations/asia-northeast1/connections/${GITHUB_HOST}/repositories/${GITHUB_REPO}
  repositoryType: GITHUB
serviceAccount: projects/${PROJECT_ID}/serviceAccounts/cloud-build-runner@${PROJECT_ID}.iam.gserviceaccount.com
EOF

gcloud builds triggers import --source=./pr-trigger.yaml --region asia-northeast1
```
2. demo-backend-api-push-main
```bash
cat <<EOF > ./main-trigger.yaml
description: Build and deploy to Cloud Run service demo-backend-api on push to main
filename: cloudbuild.yaml
includeBuildLogs: INCLUDE_BUILD_LOGS_WITH_STATUS
name: demo-backend-api-push-main
repositoryEventConfig:
  push:
    branch: ^main$
  repository: projects/${PROJECT_ID}/locations/asia-northeast1/connections/${GITHUB_HOST}/repositories/${GITHUB_REPO}
  repositoryType: GITHUB
serviceAccount: projects/${PROJECT_ID}/serviceAccounts/cloud-build-runner@${PROJECT_ID}.iam.gserviceaccount.com
EOF

gcloud builds triggers import --source=./main-trigger.yaml --region asia-northeast1
```
3. demo-backend-api-remove-cloud-run-tag	
```bash
cat <<EOF > ./rm-tag-trigger.yaml
description: Remove Cloud Run Tags
gitFileSource:
  path: cloudbuild_rm_run_tag.yaml
  repository: projects/${PROJECT_ID}/locations/asia-northeast1/connections/${GITHUB_HOST}/repositories/${GITHUB_REPO}
  revision: refs/heads/main
name: demo-backend-api-remove-cloud-run-tag
serviceAccount: projects/${PROJECT_ID}/serviceAccounts/cloud-build-runner@${PROJECT_ID}.iam.gserviceaccount.com
sourceToBuild:
  ref: refs/heads/main
  repository: projects/${PROJECT_ID}/locations/asia-northeast1/connections/${GITHUB_HOST}/repositories/${GITHUB_REPO}
EOF

gcloud builds triggers import --source=./rm-tag-trigger.yaml --region asia-northeast1
```

## Cloud Run サービスの準備
(TODO: Cloud Run サービスが存在しない場合は作成するよう CI Pipeline を修正)  
サンプルコンテナを利用して仮サービスを作成（コストはかからない）
1. demo-backend-api-dev	
2. demo-backend-api-prod
```bash
gcloud config set run/region asia-northeast1
gcloud config set run/platform managed

gcloud run deploy demo-backend-api-dev --image=us-docker.pkg.dev/cloudrun/container/hello --allow-unauthenticated
gcloud run deploy demo-backend-api-prod --image=us-docker.pkg.dev/cloudrun/container/hello --allow-unauthenticated
```

## 試してみる
1. GitHub で Pull Request (to main branch) を作成しましょう
2. main ブランチにマージしましょう
3. ブランチを削除しましょう
4. Cloud Deploy で dev to prod にプロモートしてみましょう

# SAML
IdP-initiated で SAML 認証する SP のサンプルアプリケーションです

# 注意
このサンプルアプリケーションは未実装の部分があり、不完全なアプリケーションなので注意が必要です。

- CSRF 対策などの不備
- SingleLogout の Response 未署名（仕様では MUST）

# 設計
[design.md](https://github.com/ksrnnb/saml-impl/blob/main/design.md)

# 使い方

## 1. 前提条件
- Go 1.20 以上がインストールされていること

## 2. 起動
下記コマンドを実行すると、docker で keycloak が立ち上がり、ローカルで golang のアプリケーションが立ち上がります。

```bash
make run
```

## 3. IdP に SP のメタデータを登録
まずは SP ( http://localhost:3000 ) にログインして、必要となるメタデータを確認します。

![SP メタデータ画面](https://user-images.githubusercontent.com/48155865/233210659-15f6b24e-e879-470f-a31a-3922fff65f25.png)


次に別タブを開いて、 keycloak ( http://localhost:8080 ) にログインして、必要となるメタデータを登録します。


### IdP ログイン
下記情報で admin console にログインします。

- username: admin
- password: admin

### クライアントの作成
Clients > Create client でクライアントを作成します。

- Client type: SAML
- Client ID: SP で取得した SP Entity ID

### Valid redirect URIs の設定
作成したクライアントの設定画面から、 Valid redirect URIs に SP の URL を設定します。今回は `http://localhost:3000/*` とします。
![Valid redirect URIs](https://user-images.githubusercontent.com/48155865/233211085-736267b9-24c4-40bd-bcc6-e4f58a6e6477.png)

### SSO URL の設定
IDP-Initiated SSO URL name を `38azqp4z` にします。すると、SSO URL が下に表示されます。ここで Save ボタンを押しておきます。

![SSO URL](https://user-images.githubusercontent.com/48155865/233210971-0852285a-8dd5-4e37-8f26-67062fd21164.png)


### NameID の設定
NameID format を email に設定、Force POST binding を On にしておきます。

![NameID](https://user-images.githubusercontent.com/48155865/233211385-4e30eac7-6ec0-4e42-a961-420bda21fabb.png)

### クライアントリクエストの署名設定
Keys タブを選択し、Client signature required を Off にしておきます。この設定により、 SP-initiated のときに署名を不要にします。

### ACS URL の設定
Advanced タブを選択し、 Assertion Consumer Service POST Binding URL に、SP のページから取得した ACS URL を設定します。

### IdP ユーザーのメールアドレスの設定
左タブの Users を選択し、ユーザーの設定画面に遷移して、ユーザーを登録します。

今回は、SP に[2人のユーザー](https://github.com/ksrnnb/saml-impl/blob/2f6d33898c1eedd37e0c7d8023c4a8c5563a96ef/model/user.go#L5-L11)が存在するので、それぞれのメールアドレスで登録しておきます。

## 4. SP に IdP のメタデータを登録
IdP の設定が完了したので、次は SP にメタデータを登録します。

### IdP のメタデータ取得
まずは IdP のメタデータを Realm settings のページでダウンロードします。

![メタデータダウンロードリンク](https://user-images.githubusercontent.com/48155865/202046891-6cb65962-2f36-442f-bb5e-30658eedf144.png)

### メタデータ登録
ダウンロードした xml ファイルを、SP の設定ページからアップロードします。アップロードすると、フォームに自動入力されるようになっています。

![メタデータ設定画面](https://user-images.githubusercontent.com/48155865/202049355-5959d732-c0d2-4a58-9a69-ef25bea09351.png)

## 5. ログアウトする
フローをわかりやすくするため、SP と IdP 両方ログアウトしておきます。

SP 側は、トップページから [2. SSO ログイン](http://localhost:3000/ssologin) にアクセスすると、ログアウトボタンが表示されます。

## 6. IdP-initiated で SAML 認証する
IdP-initiated の SSO URL にアクセスします。手順通りに行っていれば、SSO URL は http://localhost:8080/realms/master/protocol/saml/clients/38azqp4z になっています。

SAML 認証後、下の画面のように「SAML 認証に成功しました」というメッセージが出れば、SAML 認証が問題なく実行できています。

![SAML認証後の画面](https://user-images.githubusercontent.com/48155865/233212172-290d5b69-a97c-4d0d-a857-d1aecc4478a2.png)

## 8. IdP-initiated でログアウトする
[IdP のページ](http://localhost:8080/admin/master/console) に遷移して、ログアウトします。Single Logout の設定が正しくできていれば、SP もログアウトされます。

## 9. SP-initiated で SAML 認証する

SP-initiated の認証ページは、トップページの「SAML SSO でログインする」リンクから表示できます。
![ログイン画面](https://user-images.githubusercontent.com/48155865/234134503-760f3cb5-43a8-412c-aa18-688c764f2151.png)

ページを表示すると、会社IDの入力が求められるので、表示されている表の ID を入力します。正しいIDを入力してからボタンをクリックすると、 SP-initiated による SAML 認証が開始します。

![SP-initiated 認証画面](https://user-images.githubusercontent.com/48155865/234134609-38db58f6-6f56-4628-83e4-f7c1b4ad34be.png)

## 8. 立ち下げ

```bash
make down
```

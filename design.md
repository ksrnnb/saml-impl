# 概要
会社に所属しているユーザーが SAML 認証でログインし、IdP 起点でログアウトができるようにする。今回は SP を実装し、IdP は Keycloak を Docker で利用する。

```mermaid
graph LR
  User-->SP["SP <br> localhost:3000"]
  SP-->IdP["Keycloak <br> localhost:8080"]
  SP-->db["DB <br> sqlite"]
```

## 実装すること
- IdP メタデータの登録・更新
- SP-initiated SSO
- IdP-initiated SSO
- IdP-initiated Single Logout

## 実装しないこと
- ユーザーの追加・更新・削除
- 会社の追加・更新・削除
- Single Logout の署名（仕様では MUST）
- SP-initiated Single Logout
- persistent id への対応
- Just In Time Provisioning
- SAML 機能の有効化・無効化

# 詳細
## DB テーブル設計
会社に複数のユーザーが所属する。

今回は簡単のため、SAML 認証を有効にするかどうかのフラグはもたない。

```mermaid
erDiagram
    companies ||--|{ users : "1:N"
    companies ||--|o idp_metadatas : "1:1"

    companies {
        string id
        string name
    }

    users {
        string id
        string company_id
        string password
        string email
        int user_type
    }

    idp_metadatas {
        int id
        string name
        string description
        float price
    }
```

## 認証
SP-initiated と IdP-initiated の両方に対応する。

### SP-initiated SSO

1. SAML SSO でログイン
2. IdP SSO URL にリダイレクトするレスポンスを返す。必要であれば、RelayState も返す。このとき、セッションを利用して Request の ID を保存する。
3. そのままリダイレクトする。
4. IdP の認証画面を返す。認証済みの場合はスキップされる。
5. 認証する。
6. HTML レスポンスを返す。HTML には JavaScript のコードが含まれる。
7. JavaScript が実行され、HTTP POST Binding で SAMLResponse を ACS URL に送信する。
    - クロスオリジンではあるが、Cookie は送信される。IdP から送られる HTML レスポンスは、JavaScript で form を submit するようになっているので、[単純リクエスト](https://developer.mozilla.org/ja/docs/Web/HTTP/CORS#%E5%8D%98%E7%B4%94%E3%83%AA%E3%82%AF%E3%82%A8%E3%82%B9%E3%83%88)に相当する。
    - したがって、2番で保存した Request の ID を取得し、 InResponseTo を検証できる。
8. RelayState に応じてリダイレクト先を決めて、レスポンスを返す。

```mermaid
sequenceDiagram
  autonumber
  participant UA as User Agent
  participant SP
  participant IdP

  UA->>SP: SAML SSO でログイン
  SP->>UA: IdP SSO URL にリダイレクト with SAMLRequest
  UA->>IdP: SAMLRequest
  IdP->>UA: 認証画面
  UA->>IdP: 認証
  IdP->>UA: HTML with SAMLResponse
  UA->>SP: SAML Response
  SP->>UA: 任意のページにリダイレクト
```
### IdP-initiated SSO
SP-initiated とは異なり、3番の SAMLResponse を受け取った時、InResponseTo の検証ができない。これは CSRF 攻撃のリスクとなるので注意。[ritou さんの記事](https://zenn.dev/ritou/articles/9366cc534860e5) に分かりやすく記載されている。

```mermaid
sequenceDiagram
  autonumber
  participant UA as User Agent
  participant SP
  participant IdP

  UA->>IdP: SAML SSO で SP にログイン
  IdP->>UA: HTML with SAMLResponse
  UA->>SP: SAML Response
  SP->>UA: SP の任意のページにリダイレクト
```

### パスワードログインの制御
ユーザータイプが Normal の場合は、パスワードによるログインを制限する。 Admin はパスワードによるログインを許可する。

## ログアウト
今回は IdP-initiated のみ対応する。また、仕様では署名が MUST であるが、簡単のため今回は対応しない。

ログアウト時、ユーザーの情報からセッションキーを推測することはできないので、Invalidation によりセッションを無効化する。

### IdP-initiated Single Logout

```mermaid
sequenceDiagram
  autonumber
  participant User as User Agent
  participant SP
  participant IdP

  User->>IdP: ログアウト
  IdP->>User: HTML with SAML LogoutRequest

  User->>SP: SAML LogoutRequest
  SP->>SP: ユーザーのセッション終了
  SP->>User: SAML LogoutResponse

  User->>IdP: SAML LogoutResponse
  IdP->>IdP: ユーザーのセッション終了
  IdP->>User: ログアウト成功
```

### セッションの invalidation
セッションを Invalidation することで、Single Logout を実現する。

#### Invalidation
ログアウト時に Key Value Store (KVS) に Invalidation のキーを追加する。キーは `invalidate:userId:<user_id>` とする。

```mermaid
sequenceDiagram
  autonumber
  participant SP
  participant kvs as KVS

  SP->>kvs: Invalidation のキー追加
  kvs->>SP: ok 
```

#### Invalidation 後
Invalidation の後では、認証が必要なページにアクセスするとログインページにリダイレクトされる。

```mermaid
sequenceDiagram
  autonumber
  participant User
  participant SP
  participant kvs as KVS

  User->>SP: 認証が必要なページにアクセス
  SP->>kvs: セッションキーを検索
  kvs->>SP: ユーザーID
  SP->>kvs: 対象のユーザーの Invalidation を検索
  kvs->>SP: Invalidation されている
  SP->>User: ログインページにリダイレクト
```

#### Invalidation 後のログイン
ログインした時に Invalidation のキーを削除する。

```mermaid
sequenceDiagram
  autonumber
  User->>SP: ログイン
  SP->>kvs: セッションキーの追加
  kvs->>SP: ok
  SP->>kvs: Invalidation のキー削除
  kvs->>SP: ok
  SP->>User: SetCookie: <セッションキー>
```

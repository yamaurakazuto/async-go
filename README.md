# 銀行アプリ (PoC)

## 構成
- フロントエンド: TypeScript (静的ファイル)
- バックエンド: Go
- データベース: MySQL

## 起動方法

### 1. MySQL
以下の例では `banking` データベースを利用します。

```sql
CREATE DATABASE banking;

CREATE TABLE users (
  id BIGINT AUTO_INCREMENT PRIMARY KEY,
  email VARCHAR(255) NOT NULL UNIQUE,
  password_hash VARCHAR(255) NOT NULL
);

CREATE TABLE accounts (
  id BIGINT AUTO_INCREMENT PRIMARY KEY,
  user_id BIGINT NOT NULL,
  balance DECIMAL(12, 2) NOT NULL DEFAULT 0,
  CONSTRAINT fk_accounts_user FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE TABLE transfers (
  id BIGINT AUTO_INCREMENT PRIMARY KEY,
  from_account_id BIGINT NOT NULL,
  to_account_number VARCHAR(64) NOT NULL,
  amount DECIMAL(12, 2) NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT fk_transfers_account FOREIGN KEY (from_account_id) REFERENCES accounts(id)
);
```

> **Note**
> 現時点では `password_hash` を平文パスワードとして扱っています。

### 2. API サーバー

```bash
export MYSQL_DSN="root:password@tcp(127.0.0.1:3306)/banking?parseTime=true"

go run ./cmd/server
```

### 3. フロントエンド
TypeScript を更新した場合は、下記でビルドします。

```bash
npm install
npm run build
```

その後、`http://localhost:8080` にアクセスしてください。

## API

- `POST /api/login`
  - body: `{ "email": "...", "password": "..." }`
- `GET /api/balance?user_id=1`
- `POST /api/transfer`
  - body: `{ "from_user_id": 1, "to_account_number": "123-456", "amount": 1000 }`

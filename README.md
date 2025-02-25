# Cursor Reset Service

Cursor Reset Service 是一个用于管理 Cursor 账号的服务，包含前后端一体化的解决方案。

## 部署说明

### 使用 Docker Compose 本地部署

1. 确保已安装 Docker 和 Docker Compose
2. 克隆本仓库
3. 创建 `.env` 文件并配置数据库连接信息：

```bash
DATABASE_URL="postgresql://username:password@host:port/dbname"
PORT=8080
```

4. 在项目根目录运行：

```bash
docker-compose up -d
```

5. 访问 http://localhost:8080 即可使用管理界面
6. 默认管理员账号：admin，密码：admin123

### 在 Zeabur 等云平台部署

1. 在 Zeabur 创建新项目
2. 连接 GitHub 仓库
3. 在 Zeabur 控制台中配置以下环境变量：
   - `DATABASE_URL`: PostgreSQL 数据库连接 URL（格式：postgresql://username:password@host:port/dbname）
   - `PORT`: 应用监听端口（可选，默认 8080）
   - `ADMIN_USERNAME`: 管理员用户名（可选，默认 admin）
   - `ADMIN_PASSWORD`: 管理员密码（可选，默认 admin123）
4. Zeabur 会自动检测 `docker-compose.yml` 和 `Dockerfile`，并进行构建
5. 部署完成后，可以通过分配的域名访问服务

## 环境变量配置

应用通过环境变量进行配置，所有配置项如下：

- `DATABASE_URL`: PostgreSQL 数据库连接 URL（必需，格式：postgresql://username:password@host:port/dbname）
- `PORT`: 应用监听端口（默认 8080）
- `ADMIN_USERNAME`: 管理员用户名（默认 admin）
- `ADMIN_PASSWORD`: 管理员密码（默认 admin123）
- `JWT_ADMIN_SECRET`: 管理员 JWT 密钥（默认 admin-secret-key）
- `JWT_APP_SECRET`: 应用 JWT 密钥（默认 app-secret-key）

## API 接口

### 管理员接口

- `POST /api/login`: 管理员登录
- `POST /api/activation-codes`: 创建激活码
- `GET /api/activation-codes`: 获取激活码列表
- `GET /api/activation-codes/:id`: 获取激活码详情
- `PUT /api/activation-codes/:id/status`: 更新激活码状态

### 应用接口

- `POST /api/app/activate`: 激活码激活
- `GET /api/app/account`: 获取账号列表
- `POST /api/app/account`: 上传账号信息
- `GET /api/app/code-info`: 获取激活码信息

## 客户端使用示例

### 激活码激活

```bash
curl --location --request POST 'http://localhost:8080/api/app/activate' \
--header 'Content-Type: application/json' \
--data-raw '{
    "code": "YOUR_ACTIVATION_CODE"
}'
```

### 上传账号信息

```bash
curl --location --request POST 'http://localhost:8080/api/app/account' \
--header 'Authorization: Bearer YOUR_APP_TOKEN' \
--header 'Content-Type: application/json' \
--data-raw '{
    "email": "user@example.com",
    "emailPassword": "emailPassword",
    "cursorPassword": "cursorPassword",
    "accessToken": "cursor_access_token",
    "refreshToken": "cursor_refresh_token"
}'
```

### 获取账号列表

```bash
curl --location --request GET 'http://localhost:8080/api/app/account' \
--header 'Authorization: Bearer YOUR_APP_TOKEN'
```

### 获取激活码信息

```bash
curl --location --request GET 'http://localhost:8080/api/app/code-info' \
--header 'Authorization: Bearer YOUR_APP_TOKEN'
``` 
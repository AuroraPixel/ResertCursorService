# 多阶段构建

# 前端构建阶段
FROM node:18-alpine AS frontend-builder
WORKDIR /app/web

# 复制前端代码
COPY web/package*.json ./
RUN npm install --legacy-peer-deps

COPY web/ ./
RUN npm run build

# 后端构建阶段
FROM golang:1.21-alpine AS backend-builder
WORKDIR /app

# 安装依赖
RUN apk add --no-cache gcc musl-dev

# 复制 Go 模块定义
COPY go.mod go.sum ./
RUN go mod download

# 复制后端代码
COPY . .
# 复制前端构建产物
COPY --from=frontend-builder /app/web/dist ./web/dist

# 编译后端
RUN CGO_ENABLED=1 GOOS=linux go build -a -o cursor-reset-service .

# 最终运行阶段
FROM alpine:latest
WORKDIR /app

# 安装运行时依赖
RUN apk add --no-cache ca-certificates tzdata

# 设置时区
ENV TZ=Asia/Shanghai

# 复制编译好的二进制文件
COPY --from=backend-builder /app/cursor-reset-service .
# 复制前端构建产物
COPY --from=frontend-builder /app/web/dist ./web/dist

# 定义环境变量
ENV PORT=8080 \
    DATABASE_URL="" \
    ADMIN_USERNAME="" \
    ADMIN_PASSWORD="" \
    JWT_ADMIN_SECRET="" \
    JWT_APP_SECRET=""

# 暴露端口
EXPOSE ${PORT}

# 启动应用
CMD ["./cursor-reset-service"] 
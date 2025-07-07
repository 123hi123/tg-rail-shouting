# 台铁竹北站列车监控 Telegram Bot

这是一个监控台铁竹北站列车信息的后端服务，会在下班时间（6点后）每30分钟检查一次列车信息，并通过Telegram Bot推送到指定聊天。

## 功能特点

- 🚄 实时监控竹北站列车信息
- 📱 通过Telegram Bot推送消息
- ⏰ 可配置的监控时间段（默认6点后）
- 🎯 按到达时间排序显示列车
- 🗺️ 显示每个班次的途经站点信息
- 📊 显示到富岗站的完整路线信息

## 环境要求

- Go 1.21+
- Telegram Bot Token
- TDX API 账号和密钥（可选 - 不填写将使用免费API）

## 配置说明

1. 复制 `.env.example` 为 `.env`
2. 填写以下配置：
   - `TELEGRAM_BOT_TOKEN`: Telegram Bot Token（必填）
   - `TELEGRAM_CHAT_ID`: 接收消息的Telegram Chat ID（必填）
   - `TDX_CLIENT_ID`: TDX API客户端ID（可选 - 用于提升API限制）
   - `TDX_CLIENT_SECRET`: TDX API客户端密钥（可选 - 用于提升API限制）

## 使用方法

```bash
# 安装依赖
go mod tidy

# 运行服务
go run main.go
```

## 获取必要的API密钥

### TDX API 密钥（可选）
**免费使用：** 程序默认使用免费API，每日限制50次请求，无需注册。

**注册使用：** 如需更高频率使用：
1. 访问 https://tdx.transportdata.tw/
2. 注册账号并完成邮箱验证
3. 在会员中心获取API密钥

## 重要信息

- **竹北站ID**: 1180
- **监控时间**: 6点后到11点
- **检查频率**: 30分钟一次
- **API限制**: 免费使用每日50次请求
- **方向设置**: 1=北上，0=南下

### Telegram Bot Token
1. 与 @BotFather 聊天
2. 创建新的Bot: `/newbot`
3. 获取Bot Token

### Telegram Chat ID
1. 与你的Bot聊天，发送任意消息
2. 访问: `https://api.telegram.org/bot<YOUR_BOT_TOKEN>/getUpdates`
3. 从返回的JSON中获取chat_id
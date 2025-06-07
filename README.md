---
title: "1"
language_tabs:
  - shell: Shell
  - http: HTTP
  - javascript: JavaScript
  - ruby: Ruby
  - python: Python
  - php: PHP
  - java: Java
  - go: Go
toc_footers: []
includes: []
search: true
code_clipboard: true
highlight_theme: darkula
headingLevel: 2
generator: "@tarslib/widdershins v4.0.30"

---

# 1

Base URLs:

* <a href="http://dev-cn.your-api-server.com">开发环境: http://dev-cn.your-api-server.com</a>

# Authentication

# 用户相关

## POST 获取验证码

POST /captcha

- 成功时，接口会向指定邮箱发送6位数字验证码，验证码5分钟内有效。
- 失败时，`code` 字段非0，`message` 字段包含错误原因。
- 邮箱模板文件路径为 `./resources/template/email/captcha.html`，模板中需包含 `{{CODE}}` 占位符用于插入验证码。

### 典型错误码

| code | message                        | 说明                                         |
|------|--------------------------------|----------------------------------------------|
| 0    | Captcha sent                   | 验证码发送成功                               |
| 1    | Invalid request or missing email| 请求参数错误或缺少邮箱                       |
| 2    | Failed to load or send email    | 邮件模板加载或发送失败                       |

> Body 请求参数

```json
{
  "email": "abcxiaoyao1234@163.com"
}
```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|Content-Type|header|string| 是 |none|
|body|body|object| 否 |none|
|» email|body|string| 是 |none|

> 返回示例

> 验证码发送成功

```json
{
  "code": 0,
  "message": "Captcha sent"
}
```

```json
{
  "code": 1,
  "message": "Missing email"
}
```

```json
{
  "code": 2,
  "message": "Failed to send email"
}
```

```json
{
  "code": 3,
  "message": "Too many requests for this email, please try again later"
}
```

> 400 Response

```json
{
  "code": 1,
  "message": "Invalid request"
}
```

> 429 Response

```json
{
  "code": 3,
  "message": "Too many requests for this email, please try again later"
}
```

> 500 Response

```json
{
  "code": 2,
  "message": "Failed to send email"
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|验证码发送成功|Inline|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|请求参数错误（如邮箱缺失或格式错误）|Inline|
|429|[Too Many Requests](https://tools.ietf.org/html/rfc6585#section-4)|同一邮箱或同一IP短时间内请求过多|Inline|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|服务器内部错误（如模板加载失败、邮件发送失败）|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» code|integer|true|none||none|
|» message|string|true|none||none|

状态码 **400**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» code|integer|true|none||none|
|» message|string|true|none||none|

状态码 **429**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» code|integer|true|none||none|
|» message|string|true|none||none|

状态码 **500**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» code|integer|true|none||none|
|» message|string|true|none||none|

## POST 注册

POST /register

### 状态码说明

| code | message                        | 说明                       |
|------|--------------------------------|----------------------------|
| 0    | Register success               | 注册成功                   |
| 1    | Invalid request / Missing fields / Username or email already exists | 请求参数错误/缺少字段/用户名或邮箱已存在 |
| 2    | Database connection error / Database error / Password encryption failed / Failed to generate userId / Register failed | 服务器内部错误             |
| 4    | Invalid or expired captcha     | 验证码无效或已过期         |

### 说明

- 用户名只能用字母数字下划线
- 所有字段均需去除首尾空格后校验。
- 邮箱验证码通过 `VerifyCaptcha(email, captcha)` 校验。
- 用户名或邮箱已存在时，注册失败。
- 密码使用 bcrypt 加密存储。
- 注册成功后返回 code=0。

> Body 请求参数

```json
{
  "username": "ayndpa",
  "password": "abc134625",
  "email": "abcxiaoyao1234@163.com",
  "captcha": "065074"
}
```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|body|body|object| 否 |none|
|» username|body|string| 是 |none|
|» password|body|string| 是 |none|
|» email|body|string| 是 |none|
|» captcha|body|string| 是 |none|

> 返回示例

> 注册成功

```json
{
  "code": 0,
  "message": "Register success"
}
```

```json
{
  "code": 1,
  "message": "Missing fields"
}
```

```json
{
  "code": 4,
  "message": "Invalid or expired captcha"
}
```

```json
{
  "code": 1,
  "message": "Username or email already exists"
}
```

```json
{
  "code": 2,
  "message": "Database connection error"
}
```

```json
{
  "code": 3,
  "message": "Too many requests for this email, please try again later"
}
```

> 400 Response

```json
{
  "code": 1,
  "message": "Missing email"
}
```

> 401 Response

```json
{
  "code": 4,
  "message": "Invalid or expired captcha"
}
```

> 409 Response

```json
{
  "code": 1,
  "message": "Username or email already exists"
}
```

> 500 Response

```json
{
  "code": 2,
  "message": "Database connection error"
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|注册成功|Inline|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|请求参数缺失或格式错误|Inline|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|验证码错误或已过期|Inline|
|409|[Conflict](https://tools.ietf.org/html/rfc7231#section-6.5.8)|用户名或邮箱已存在|Inline|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|服务器内部错误，如数据库或加密失败|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» code|integer|true|none||none|
|» message|string|true|none||none|

状态码 **400**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» code|integer|true|none||none|
|» message|string|true|none||none|

状态码 **401**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» code|integer|true|none||none|
|» message|string|true|none||none|

状态码 **409**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» code|integer|true|none||none|
|» message|string|true|none||none|

状态码 **500**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» code|integer|true|none||none|
|» message|string|true|none||none|

## POST 登录

POST /login

### 状态码说明

| code | message                                      | 说明                                   |
|------|----------------------------------------------|----------------------------------------|
| 0    | Login success                                | 登录成功，返回 token                   |
| 1    | Invalid request / Missing fields / User not found / Database connection error / Database error | 请求参数错误/缺少字段/用户不存在/数据库错误 |
| 2    | Incorrect password                           | 密码错误                               |
| 3    | Token generation failed                      | Token 生成失败                         |
| 4    | Ban check failed                             | 封禁状态检查失败                       |
| 5    | User is banned[: BanReason]                  | 用户被封禁，附带封禁原因（如有）        |

### 说明

- 接口路径：`/user/login`，请求方法：`POST`，请求体为 JSON，包含 `username`（可为用户名、邮箱或纯数字ID）和 `password` 字段。
- 所有字段会去除首尾空格后校验，缺失或为空直接返回错误。
- 支持用户名、邮箱、纯数字ID三种方式登录。
- 登录时会校验用户是否存在、密码是否正确、是否被封禁。
- 登录成功返回 JWT Token。
- 被封禁用户会返回封禁原因（如有）。
- 密码采用 bcrypt 加密校验。

> Body 请求参数

```json
{
  "username": "test",
  "password": "test"
}
```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|Content-Type|header|string| 是 |none|
|body|body|object| 否 |none|
|» username|body|string| 是 |none|
|» password|body|string| 是 |none|

> 返回示例

> 登录成功

```json
{
  "code": 0,
  "message": "Login success",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

```json
{
  "code": 1,
  "message": "User not found"
}
```

> 400 Response

```json
{
  "code": 1,
  "message": "Missing fields"
}
```

> 401 Response

```json
{
  "code": 1,
  "message": "User not found"
}
```

> 500 Response

```json
{
  "code": 1,
  "message": "Database error"
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|登录成功|Inline|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|请求参数错误（如缺少字段、格式错误）|Inline|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|用户不存在或密码错误|Inline|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|服务器内部错误或数据库错误|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» code|integer|true|none||none|
|» message|string|true|none||none|
|» token|string|true|none||none|

状态码 **400**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» code|integer|true|none||none|
|» message|string|true|none||none|

状态码 **401**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» code|integer|true|none||none|
|» message|string|true|none||none|

状态码 **500**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» code|integer|true|none||none|
|» message|string|true|none||none|

# 数据模型


# AI Commit 规范

所有 AI 助手在为本项目生成提交时，必须严格遵守此规范。

---

## 提交格式

每条提交必须包含 **Title**（标题）和 **Description**（描述）两部分。

### Title

```
<type>: <subject>
```

- `type` 和 `:` 之间**无空格**，`:` 后跟一个**英文空格**
- 整行**不超过 72 个字符**
- 使用**祈使句**（如 "添加"、"修复"、"更新"，不要用 "添加了"、"修复了"）
- 结尾**不加句号**

### Description

Description 必须用中文书写，清晰描述本次提交的**修改内容**和**修改原因**。

```
1. 修改了哪些文件 / 模块
2. 为什么做这个修改
3. 如果有需要注意的事项（如破坏性变更），也在这里说明
```

---

## Type 类型

| Type       | 说明                 |
| ---------- | -------------------- |
| `feat`     | 新功能               |
| `fix`      | 修复 Bug             |
| `docs`     | 仅文档变更           |
| `style`    | 代码格式调整         |
| `refactor` | 代码重构             |
| `perf`     | 性能优化             |
| `test`     | 添加或修改测试       |
| `build`    | 构建或依赖变更       |
| `chore`    | 杂项                 |
| `revert`   | 回退之前的提交       |

---

## 提交示例

### ✅ 正确示例

```
feat: 添加 DeepSeek token 余额查询接口

1. 新增 `internal/deepseek/client.go`，实现 DeepSeek API 调用
2. 在 `/api/balance/deepseek` 路由注册查询端点
3. 原因是用户需要实时查看 DeepSeek 平台的 token 余额
```

```
fix: 修复 MiMo 接口返回格式解析错误

1. 修改 `internal/mimo/parser.go`，调整 JSON 字段映射
2. 原因是上游 API 响应格式发生了变更
3. 旧字段 `remaining` 已改为 `remaining_tokens`
```

```
docs: 更新 README 中的已接入平台列表

1. 新增 XiaoMi MiMo 到平台列表
2. 同步更新了对应的 API 文档示例
```

### ❌ 错误示例

| 错误写法                        | 问题                           |
| ------------------------------- | ------------------------------ |
| `update something`              | 缺少规范 type，描述模糊        |
| `feat: 添加了新的余额查询功能`   | 「添加了」非祈使句             |
| `fix: 修bug`                    | subject 太简短，无 description |
| 只有标题没有 description        | 缺少必要的修改说明             |

---

---

## AI 身份标识（Co-author）

所有 AI 提交的 commit **必须**在底部添加 `Co-authored-by` 声明，格式如下：

```
Co-authored-by: <模型名称> <模型名称@公司域名>
```

### 命名规则

| 字段     | 格式                                                       |
| -------- | ---------------------------------------------------------- |
| 名称     | 使用 AI 模型的完整产品名称（如 `DeepSeek V4 Flash`）       |
| 邮箱     | 全小写 + 连字符拼接的模型名称，`@` 后跟对应公司域名        |

### 已注册的 AI 身份

| AI 模型                  | Co-author 格式                                                    |
| ------------------------ | ----------------------------------------------------------------- |
| DeepSeek V4 Flash        | `DeepSeek V4 Flash <deepseek-v4-flash@deepseek.com>`              |

> 后续如有新 AI 模型加入，按此规则扩展即可。

### 示例

```
docs: 更新 README 中的已接入平台列表

1. 新增 XiaoMi MiMo 到平台列表

Co-authored-by: DeepSeek V4 Flash <deepseek-v4-flash@deepseek.com>
```

---

## 强制规则

1. **禁止**使用 `update`、`modify`、`change` 作为 type — 必须使用上方表格中的规范 type
2. **禁止**提交信息为空或仅有 "update"、"wip" 等无意义内容
3. **禁止**省略 Description — 每次提交都必须写清楚「改了哪些文件」和「为什么改」
4. **禁止**省略 `Co-authored-by` — AI 提交必须声明 AI 身份
5. 如果一次修改涉及多个 type，选择**最主要的一个**作为 type

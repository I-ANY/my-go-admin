# 聊天机器人 HTML/Markdown 功能演示

## 功能概述

聊天机器人现在支持HTML和Markdown语法渲染，提供更丰富的内容展示能力。

## 支持的格式

### 1. 纯文本

普通文本内容，直接显示。

### 2. HTML

支持HTML标签的内容，如：

```html
<h2>标题</h2>
<p>段落内容</p>
<strong>粗体</strong>
<a href="#">链接</a>
```

### 3. Markdown

支持完整的Markdown语法，包括：

#### 标题

```markdown
# 一级标题

## 二级标题

### 三级标题
```

#### 文本格式

```markdown
**粗体文本** _斜体文本_ ~~删除线文本~~ `行内代码`
```

#### 代码块

````markdown
```javascript
function hello() {
  console.log('Hello World!');
}
```
````

````

#### 列表
```markdown
- 无序列表项1
- 无序列表项2

1. 有序列表项1
2. 有序列表项2
````

#### 任务列表

```markdown
- [x] 已完成任务
- [ ] 未完成任务
```

#### 引用

```markdown
> 这是一个引用块
```

#### 表格

```markdown
| 列1   | 列2   | 列3   |
| ----- | ----- | ----- |
| 内容1 | 内容2 | 内容3 |
```

#### 链接和图片

```markdown
[链接文本](URL) ![图片描述](图片URL)
```

#### 分割线

```markdown
---
```

## 正在思考中功能演示

### 功能说明

当用户发送消息后，聊天机器人会立即显示"正在思考中..."的提示，并伴有动画效果。这个提示会在收到后端回复后自动替换为实际内容。

### 演示步骤

1. 打开聊天机器人
2. 发送一条消息（如："请帮我分析一个复杂的问题"）
3. 观察立即出现的"正在思考中..."提示
4. 等待几秒钟后，提示会被实际回复替换

### 特点

- **即时反馈**：发送消息后立即显示思考提示
- **动画效果**：带有优雅的点点动画
- **自动替换**：收到回复后自动移除提示
- **错误处理**：请求失败时也会移除提示

## 自动类型检测

组件会自动检测消息内容类型：

- 包含HTML标签的内容 → HTML类型
- 包含Markdown语法的内容 → Markdown类型
- 其他内容 → 文本类型

## 使用示例

### 在代码中使用

```typescript
// 发送Markdown内容
addMessage('# 欢迎使用\n这是一个**Markdown**示例', false, 'markdown');

// 发送HTML内容
addMessage('<h2>HTML标题</h2><p>这是<strong>HTML</strong>内容</p>', false, 'html');

// 自动检测类型
addMessage('**粗体文本**和`代码`', false); // 自动检测为markdown
```

### 在机器人回复中使用

```typescript
const chatBotOptions: ChatBotOptions = {
  botName: '智能助手',

  defaultReply: async (input: string) => {
    if (input.includes('帮助')) {
      return `# 帮助信息

## 可用命令
- **help** - 显示帮助信息
- **markdown** - 显示Markdown示例
- **table** - 显示表格示例

## 代码示例
\`\`\`javascript
function hello() {
  console.log('Hello World!');
}
\`\`\`

> 这是一个引用示例

- [x] 已完成任务
- [ ] 待完成任务`;
    }

    return '我支持Markdown语法，试试输入"帮助"查看示例！';
  },
};
```

## 样式特性

- 响应式设计，适配不同屏幕尺寸
- 代码高亮显示
- 表格样式美观
- 引用块有特殊样式
- 链接有悬停效果
- 图片自适应大小

## 安全特性

- HTML内容经过安全过滤
- 防止XSS攻击
- 自动转义危险标签

## 性能优化

- 使用showdown库进行Markdown解析
- 懒加载渲染
- 缓存解析结果

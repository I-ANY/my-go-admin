# 聊天机器人插件

一个简单易用的 Vue 3 聊天机器人组件插件。

## 功能特性

- 🎯 简单易用的聊天界面
- 🤖 可自定义的机器人回复逻辑
- 🎨 支持自定义头像和样式
- ⚡ 响应式设计，支持移动端
- 🔧 插件化安装，全局可用
- 📱 支持键盘快捷键（Enter 发送）
- 🔒 临时关闭功能，支持状态重置
- ❌ 头像上的关闭按钮，悬停显示
- 🖱️ 支持拖动头像到任意位置（窗体打开后不可拖动）
- 👆 双击头像打开聊天窗口
- 📝 支持HTML和Markdown语法渲染
- 🔍 自动检测内容类型（文本/HTML/Markdown）
- 🎨 丰富的Markdown样式支持（表格、代码块、列表等）
- 💭 智能的"正在思考中"提示，提升用户体验
- ↵ 支持换行符显示（\n转换为<br>标签）
- 🔄 混合内容支持（文本、HTML、Markdown格式）

## 安装使用

### 1. 注册插件

在 `main.ts` 中注册插件：

```typescript
import { createApp } from 'vue';
import App from './App.vue';
import ChatBotPlugin from '@/chat-bot';

const app = createApp(App);

// 注册插件（可选配置）
app.use(ChatBotPlugin, {
  botName: '智能助手',
  botAvatar: '/path/to/bot-avatar.png',
  userAvatar: '/path/to/user-avatar.png',
  placeholder: '输入消息...',
  replyDelay: 1000,
});

app.mount('#app');
```

### 2. 在组件中使用

```vue
<template>
  <div>
    <h1>聊天机器人示例</h1>
    <ChatBot :options="chatBotOptions" />
  </div>
</template>

<script setup lang="ts">
  import type { ChatBotOptions } from '@/chat-bot';

  const chatBotOptions: ChatBotOptions = {
    botName: '我的助手',
    botAvatar: '/images/bot.png',
    userAvatar: '/images/user.png',
    placeholder: '请输入您的问题...',
    replyDelay: 800,

    // 自定义回复逻辑
    defaultReply: async (input: string) => {
      // 这里可以调用真实的API
      const response = await fetch('/api/chat', {
        method: 'POST',
        body: JSON.stringify({ message: input }),
      });
      const data = await response.json();
      return data.reply;
    },

    // 消息回调
    onMessage: (message) => {
      console.log('新消息:', message);
    },

    // 关闭回调
    onClose: () => {
      console.log('机器人已关闭');
    },

    // 重新打开回调
    onReopen: () => {
      console.log('机器人已重新打开');
    },

    // 登录检查配置（可选）
    requireLogin: true, // 默认true，设置为false可允许未登录用户使用

    // 关闭按钮配置（可选）
    showCloseButton: true, // 默认true，设置为false可隐藏关闭按钮
  };
</script>
```

## 配置选项

### ChatBotOptions

| 属性 | 类型 | 默认值 | 说明 |
| --- | --- | --- | --- |
| `botName` | `string` | `'智能助手'` | 机器人名称 |
| `botAvatar` | `string` | `'https://via.placeholder.com/40x40?text=Bot'` | 机器人头像URL |
| `userAvatar` | `string` | `'https://via.placeholder.com/40x40?text=User'` | 用户头像URL |
| `placeholder` | `string` | `'输入消息...'` | 输入框占位符 |
| `replyDelay` | `number` | `1000` | 回复延迟（毫秒） |
| `defaultReply` | `(input: string) => string \| Promise<string>` | 内置简单回复 | 自定义回复逻辑 |
| `onMessage` | `(message: MessageItem) => void` | - | 消息回调函数 |
| `requireLogin` | `boolean` | `true` | 是否要求用户登录才能使用 |
| `showCloseButton` | `boolean` | `true` | 是否显示关闭按钮 |
| `onClose` | `() => void` | - | 关闭机器人时的回调函数 |
| `onReopen` | `() => void` | - | 重新打开机器人时的回调函数 |

### MessageItem

```typescript
interface MessageItem {
  id: string; // 消息唯一ID
  content: string; // 消息内容
  type: BotMessageType; // 消息类型：'text' | 'html' | 'markdown' | 'image' | 'loading'
  isUser: boolean; // 是否是用户消息
  timestamp: number; // 时间戳
  renderOptions?: {
    allowHtml?: boolean; // 是否允许HTML标签
    markdownOptions?: {
      tables?: boolean; // 是否支持表格
      strikethrough?: boolean; // 是否支持删除线
      tasklists?: boolean; // 是否支持任务列表
      ghCodeBlocks?: boolean; // 是否支持GitHub风格代码块
      emoji?: boolean; // 是否支持emoji
      underline?: boolean; // 是否支持下划线
      // ... 更多Markdown选项
    };
  };
}
```

## HTML和Markdown支持

聊天机器人现在支持HTML和Markdown语法渲染，提供更丰富的内容展示能力。

### 自动类型检测

组件会自动检测消息内容类型：

- **文本**：普通文本内容
- **HTML**：包含HTML标签的内容
- **Markdown**：包含Markdown语法的内容

### 支持的Markdown语法

- **标题**：`# 一级标题`、`## 二级标题`等
- **粗体**：`**粗体文本**`
- **斜体**：`*斜体文本*`
- **删除线**：`~~删除的文本~~`
- **代码**：`` `行内代码` ``
- **代码块**：

  ````markdown
  ```javascript
  console.log('Hello World');
  ```
  ````

  ```

  ```

- **链接**：`[链接文本](URL)`
- **图片**：`![图片描述](图片URL)`
- **列表**：

  ```markdown
  - 无序列表项

  1. 有序列表项
  ```

- **任务列表**：
  ```markdown
  - [x] 已完成任务
  - [ ] 未完成任务
  ```
- **引用**：`> 引用文本`
- **表格**：
  ```markdown
  | 列1   | 列2   | 列3   |
  | ----- | ----- | ----- |
  | 内容1 | 内容2 | 内容3 |
  ```
- **分割线**：`---`

### 使用示例

```typescript
// 发送Markdown内容
addMessage('# 欢迎使用\n这是一个**Markdown**示例', false, 'markdown');

// 发送HTML内容
addMessage('<h2>HTML标题</h2><p>这是<strong>HTML</strong>内容</p>', false, 'html');

// 发送带换行的文本内容
addMessage('第一行\n第二行\n第三行', false, 'text');

// 自动检测类型
addMessage('**粗体文本**和`代码`', false); // 自动检测为markdown
addMessage('普通文本\n包含换行符', false); // 自动检测为text，换行符会被转换为<br>标签
```

### 消息格式支持

聊天机器人现在支持多种消息格式：

#### 1. 纯文本格式

- 支持换行符 `\n`，会自动转换为 `<br>` 标签
- 适合发送简单的文本消息

#### 2. HTML格式

- 支持完整的HTML标签
- 适合发送富文本内容
- 包含安全过滤，防止XSS攻击

#### 3. Markdown格式

- 支持完整的Markdown语法
- 包括标题、列表、代码块、表格等
- 自动转换为HTML显示

#### 4. 混合内容

- 支持在Markdown中包含换行符
- 支持在HTML中包含换行符
- 自动检测内容类型并正确渲染
- **重要修复**：所有消息类型（text、html、markdown）都支持换行符 `\n` 自动转换为 `<br>` 标签

#### 5. 内容类型检测优先级

系统按照以下优先级自动检测内容类型：

1. **Markdown 语法** - 优先处理为 `markdown` 类型
2. **换行符 (\n)** - 其次处理为 `text` 类型

这意味着包含 Markdown 语法的内容会优先被识别为 Markdown 类型，确保 Markdown 格式能够正确渲染。

#### 6. 换行符优化处理

系统对换行符进行智能处理：

- **单个换行符 `\n`** - 转换为 `<br>` 标签，不使用 `<p>` 标签包裹
- **连续换行符 `\n\n`** - 转换为段落分隔 `</p><p>`，使用 `<p>` 标签包裹
- **智能段落判断** - 只有包含段落分隔的内容才使用 `<p>` 标签
- **去除首尾空白** - 自动去除内容首尾的空白字符，避免末尾多余换行

这样可以避免消息看起来有多余的换行，提供更自然的显示效果。

## 正在思考中功能

聊天机器人现在支持智能的"正在思考中"提示功能，当用户发送消息后，机器人会立即显示"正在思考中..."的提示，并伴有动画效果，直到收到后端回复后自动替换为实际回复内容。

### 功能特点

- **即时反馈**：用户发送消息后立即显示思考提示
- **动画效果**：带有优雅的点点动画，提升用户体验
- **自动替换**：收到回复后自动移除思考提示并显示实际内容
- **错误处理**：如果请求失败，也会移除思考提示并显示错误信息

### 使用示例

```typescript
// 在自定义回复逻辑中，思考提示会自动显示
const chatBotOptions: ChatBotOptions = {
  defaultReply: async (input: string) => {
    // 这里可以调用真实的API，思考提示会自动显示
    const response = await fetch('/api/chat', {
      method: 'POST',
      body: JSON.stringify({ message: input }),
    });
    const data = await response.json();
    return data.reply;
  },
};
```

### 样式自定义

思考提示的样式可以通过CSS自定义：

```css
/* 自定义思考提示文字样式 */
.thinking-text {
  color: #6b7280 !important;
  font-style: italic !important;
}

/* 自定义思考动画点样式 */
.thinking-dots span {
  background: #667eea !important;
  animation-duration: 1.2s !important;
}
```

## 工具函数

插件还提供了一些实用的工具函数：

```typescript
import {
  generateMessageId,
  createMessage,
  getSimpleReply,
  simulateApiDelay,
  formatTimestamp,
} from '@/chat-bot/chat-bot';

// 生成唯一消息ID
const id = generateMessageId();

// 创建消息项
const message = createMessage('Hello', true);

// 获取简单回复
const reply = getSimpleReply('你好', '助手');

// 模拟API延迟
await simulateApiDelay(1000);

// 格式化时间戳
const time = formatTimestamp(Date.now());
```

## 样式自定义

组件使用 scoped 样式，如需自定义样式，可以通过 CSS 变量或覆盖样式：

```css
/* 自定义聊天容器样式 */
.chat-bot-container {
  --primary-color: #1890ff;
  --border-radius: 8px;
  --shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
}

/* 自定义消息样式 */
.message-content {
  background: var(--primary-color) !important;
}
```

## 注意事项

1. **登录检查**：默认情况下，只有已登录的用户才能使用聊天机器人。可以通过设置 `requireLogin: false` 来禁用此功能
2. 头像URL如果加载失败，会自动显示文字头像
3. 支持键盘 Enter 键发送消息
4. 消息会自动滚动到底部
5. 组件会自动处理加载状态和错误状态
6. 支持完全关闭机器人，关闭后机器人完全隐藏，可通过API重新打开
7. **关闭按钮**：机器人头像右上角有一个关闭按钮，鼠标悬停时显示，点击可完全关闭机器人（窗口内的关闭按钮已隐藏）
8. **拖动功能**：机器人头像支持拖动到任意位置，但只有在窗体收起状态下才能拖动，窗体打开后拖动功能会自动禁用
9. **打开方式**：双击机器人头像打开聊天窗口，单击不会触发任何操作

## 开发

如需扩展功能，可以：

1. 修改 `types.ts` 添加新的配置选项
2. 在 `chat-bot.ts` 中添加新的工具函数
3. 在 `ChatBot.vue` 中添加新的UI功能
4. 在 `index.ts` 中注册新的全局方法

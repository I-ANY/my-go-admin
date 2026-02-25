# 正在思考中功能说明

## 功能概述

为聊天机器人添加了智能的"正在思考中"提示功能，当用户发送消息后，机器人会立即显示"正在思考中..."的提示，并伴有动画效果，直到收到后端回复后自动替换为实际回复内容。

## 实现原理

### 1. 消息发送流程

```typescript
// 1. 用户发送消息
addMessage(input, true);

// 2. 立即添加思考中提示
const thinkingMessage = addMessage('正在思考中...', false, 'loading');

// 3. 调用后端API
const replyContent = await getBotReply(input);

// 4. 移除思考中提示
const thinkingIndex = messages.value.findIndex((msg) => msg.id === thinkingMessage.id);
if (thinkingIndex !== -1) {
  messages.value.splice(thinkingIndex, 1);
}

// 5. 添加实际回复
addMessage(replyContent, false);
```

### 2. 错误处理

如果API调用失败，也会移除思考中提示并显示错误信息：

```typescript
try {
  const replyContent = await getBotReply(input);
  // 成功处理...
} catch (error) {
  // 移除思考中提示
  const thinkingIndex = messages.value.findIndex((msg) => msg.id === thinkingMessage.id);
  if (thinkingIndex !== -1) {
    messages.value.splice(thinkingIndex, 1);
  }

  // 显示错误信息
  addMessage('抱歉，我现在无法回答您的问题', false);
}
```

## 界面展示

### 思考中提示的HTML结构

```html
<div v-else-if="msg.type === 'loading'" class="loading-content">
  <div class="thinking-text">{{ msg.content }}</div>
  <div class="thinking-dots"> <span></span><span></span><span></span> </div>
</div>
```

### CSS样式

```css
.loading-content {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  gap: 8px;
  padding: 8px 12px;
}

.thinking-text {
  color: #6b7280;
  font-size: 14px;
  font-style: italic;
}

.thinking-dots {
  display: flex;
  gap: 3px;
}

.thinking-dots span {
  width: 6px;
  height: 6px;
  animation: loading 1s infinite;
  border-radius: 50%;
  background: #9ca3af;
}

.thinking-dots span:nth-child(2) {
  animation-delay: 0.2s;
}

.thinking-dots span:nth-child(3) {
  animation-delay: 0.4s;
}
```

## 使用方法

### 自动启用

该功能已自动集成到聊天机器人中，无需额外配置。当用户发送消息时，会自动显示思考中提示。

### 自定义样式

可以通过CSS自定义思考提示的样式：

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

### 测试功能

可以使用测试页面验证功能：

1. 打开 `test-html-markdown.vue` 测试页面
2. 点击"测试思考中提示"按钮
3. 或直接发送包含"复杂"或"分析"关键词的消息
4. 观察思考中提示的显示和替换过程

## 技术细节

### 消息类型

- 使用 `loading` 类型的消息来显示思考中提示
- 消息内容为 "正在思考中..."
- 包含动画效果和文字提示

### 动画效果

- 使用CSS动画实现点点闪烁效果
- 三个点依次闪烁，营造思考的感觉
- 动画持续时间为1秒，无限循环

### 性能优化

- 使用 `findIndex` 和 `splice` 高效移除消息
- 避免不必要的DOM操作
- 保持界面响应性

## 兼容性

- ✅ 向后兼容，不影响现有功能
- ✅ 自动启用，无需配置
- ✅ 支持所有消息类型
- ✅ 支持错误处理

## 更新日志

### v2.1.0

- 新增"正在思考中"提示功能
- 优化用户体验，减少等待焦虑
- 添加动画效果和错误处理
- 更新文档和测试用例

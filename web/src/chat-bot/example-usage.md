# 聊天机器人功能使用示例

## 基本使用

```vue
<template>
  <div>
    <ChatBot ref="chatBotRef" :options="chatBotOptions" />
    <button @click="reopenBot" v-if="isBotClosed">重新打开机器人</button>
  </div>
</template>

<script setup lang="ts">
  import { ref } from 'vue';
  import type { ChatBotOptions } from '@/chat-bot';

  const chatBotRef = ref();
  const isBotClosed = ref(false);

  const chatBotOptions: ChatBotOptions = {
    botName: '智能助手',

    // 关闭回调
    onClose: () => {
      console.log('机器人已关闭');
      isBotClosed.value = true;
      // 可以在这里执行一些清理操作
      // 比如保存聊天记录、发送统计数据等
    },

    // 重新打开回调
    onReopen: () => {
      console.log('机器人已重新打开');
      isBotClosed.value = false;
    },

    // 显示关闭按钮（默认true）
    showCloseButton: true,
  };

  // 重新打开机器人
  const reopenBot = () => {
    if (chatBotRef.value) {
      chatBotRef.value.reopenBot();
    }
  };
</script>
```

## 隐藏关闭按钮

```vue
<script setup lang="ts">
  const chatBotOptions: ChatBotOptions = {
    botName: '智能助手',

    // 隐藏关闭按钮
    showCloseButton: false,
  };
</script>
```

## 通过键盘快捷键重新打开

```vue
<template>
  <div>
    <ChatBot ref="chatBotRef" :options="chatBotOptions" />
    <div class="hint">按 Ctrl+Shift+B 重新打开机器人</div>
  </div>
</template>

<script setup lang="ts">
  import { ref, onMounted, onUnmounted } from 'vue';

  const chatBotRef = ref();

  const chatBotOptions: ChatBotOptions = {
    botName: '智能助手',
    onClose: () => console.log('机器人已关闭'),
    onReopen: () => console.log('机器人已重新打开'),
  };

  // 键盘快捷键监听
  const handleKeydown = (event: KeyboardEvent) => {
    if (event.ctrlKey && event.shiftKey && event.key === 'B') {
      event.preventDefault();
      if (chatBotRef.value) {
        chatBotRef.value.reopenBot();
      }
    }
  };

  onMounted(() => {
    document.addEventListener('keydown', handleKeydown);
  });

  onUnmounted(() => {
    document.removeEventListener('keydown', handleKeydown);
  });
</script>
```

## HTML和Markdown支持示例

### 发送Markdown内容

```vue
<script setup lang="ts">
  import { ref } from 'vue';
  import type { ChatBotOptions } from '@/chat-bot';

  const chatBotRef = ref();

  const chatBotOptions: ChatBotOptions = {
    botName: '智能助手',

    // 自定义回复逻辑，返回Markdown内容
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

      if (input.includes('markdown')) {
        return `# Markdown 语法示例

## 文本格式
**粗体文本** 和 *斜体文本*
~~删除线文本~~

## 代码
行内代码：\`console.log('hello')\`

代码块：
\`\`\`javascript
function greet(name) {
  return \`Hello, \${name}!\`;
}
\`\`\`

## 链接和图片
[访问官网](https://example.com)
![示例图片](https://via.placeholder.com/300x200)

## 列表
1. 有序列表项1
2. 有序列表项2

- 无序列表项1
- 无序列表项2`;
      }

      if (input.includes('table')) {
        return `# 表格示例

| 功能 | 状态 | 说明 |
|------|------|------|
| HTML支持 | ✅ | 支持HTML标签渲染 |
| Markdown支持 | ✅ | 支持Markdown语法 |
| 表格支持 | ✅ | 支持表格渲染 |
| 代码高亮 | ✅ | 支持代码块高亮 |
| 任务列表 | ✅ | 支持任务列表 |

> 表格功能完全支持！`;
      }

      return '我支持Markdown语法，试试输入"markdown"或"table"查看示例！';
    },
  };
</script>
```

### 发送HTML内容

```vue
<script setup lang="ts">
  const chatBotOptions: ChatBotOptions = {
    botName: '智能助手',

    defaultReply: async (input: string) => {
      if (input.includes('html')) {
        return `<h2 style="color: #667eea;">HTML 内容示例</h2>
<p>这是<strong>HTML</strong>格式的内容，支持各种HTML标签。</p>
<ul>
  <li>列表项1</li>
  <li>列表项2</li>
  <li>列表项3</li>
</ul>
<blockquote style="border-left: 4px solid #667eea; padding-left: 10px;">
  这是一个引用块
</blockquote>
<a href="https://example.com" target="_blank">外部链接</a>`;
      }

      return '我支持HTML语法，试试输入"html"查看示例！';
    },
  };
</script>
```

### 自动类型检测

```vue
<script setup lang="ts">
  const chatBotOptions: ChatBotOptions = {
    botName: '智能助手',

    defaultReply: async (input: string) => {
      // 组件会自动检测内容类型
      if (input.includes('检测')) {
        return '**这是Markdown内容**，包含`代码`和[链接](https://example.com)';
      }

      if (input.includes('html标签')) {
        return '<strong>这是HTML内容</strong>，包含<a href="#">链接</a>';
      }

      return '这是普通文本内容';
    },
  };
</script>
```

</script>
```

## 关闭时的状态重置

当用户点击关闭按钮时，机器人会：

1. 触发 `onClose` 回调函数
2. 完全隐藏机器人（包括右下角的机器人头部）
3. 重置所有状态：
   - 收起聊天窗口
   - 清空所有消息
   - 重置输入框
   - 清除加载状态
   - 重置首次打开标记

## 重新打开机器人

机器人关闭后，可以通过以下方式重新打开：

1. 调用组件的 `reopenBot()` 方法
2. 触发 `onReopen` 回调函数
3. 机器人会重新显示，状态完全重置

## 在关闭回调中执行自定义逻辑

```vue
<script setup lang="ts">
  const chatBotOptions: ChatBotOptions = {
    botName: '智能助手',

    onClose: () => {
      // 保存聊天记录
      const chatHistory = messages.value;
      localStorage.setItem('chatHistory', JSON.stringify(chatHistory));

      // 发送统计数据
      analytics.track('chatbot_closed', {
        sessionDuration: Date.now() - sessionStartTime,
        messageCount: messages.value.length,
      });

      // 显示确认消息
      ElMessage.success('聊天已结束，下次见！');
    },
  };
</script>
```

## 条件性关闭

```vue
<script setup lang="ts">
  const chatBotOptions: ChatBotOptions = {
    botName: '智能助手',

    onClose: () => {
      // 检查是否有未保存的重要信息
      const hasImportantInfo = messages.value.some(
        (msg) => msg.content.includes('重要') || msg.content.includes('紧急'),
      );

      if (hasImportantInfo) {
        ElMessageBox.confirm('检测到重要信息，确定要关闭聊天吗？', '确认关闭', {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'warning',
        })
          .then(() => {
            // 用户确认关闭
            console.log('用户确认关闭聊天');
          })
          .catch(() => {
            // 用户取消关闭
            console.log('用户取消关闭聊天');
          });
      }
    },
  };
</script>
```

## 拖动功能

机器人头像支持拖动到任意位置，提供更好的用户体验：

### 交互行为

- **打开窗口**：双击机器人头像打开聊天窗口
- **拖动功能**：只有在窗体收起状态下才能拖动头像
- **不可拖动状态**：窗体打开后，拖动功能自动禁用
- **拖动范围**：头像会被限制在浏览器视窗范围内
- **拖动反馈**：拖动时头像会放大并显示抓取光标

### 操作方式

1. **打开窗口**：双击机器人头像打开聊天窗口
2. **开始拖动**：在收起状态下，鼠标按下头像开始拖动
3. **拖动过程**：移动鼠标，头像会跟随鼠标移动
4. **结束拖动**：松开鼠标，头像停在当前位置
5. **触摸支持**：在移动设备上支持触摸拖动

### 拖动状态样式

拖动时头像会显示以下视觉反馈：

- 光标变为抓取状态（grabbing）
- 头像轻微放大（scale: 1.1）
- 阴影效果增强

### 注意事项

- 双击头像打开聊天窗口，单击不会触发任何操作
- 拖动功能只在窗体收起状态下有效
- 窗体打开后，拖动功能自动禁用
- 拖动位置会自动保存，下次打开时保持在上次拖动的位置
- 支持响应式设计，在不同屏幕尺寸下都能正常拖动

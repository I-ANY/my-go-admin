// 聊天机器人核心逻辑和工具函数

import type { MessageItem, BotMessageType } from './types';

/**
 * 生成唯一消息ID
 */
export const generateMessageId = (): string => {
  return `msg_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`;
};

/**
 * 创建消息项
 */
export const createMessage = (
  content: string,
  isUser: boolean,
  type: BotMessageType = 'text',
): MessageItem => {
  return {
    id: generateMessageId(),
    content,
    type,
    isUser,
    timestamp: Date.now(),
  };
};

/**
 * 简单的关键词匹配回复逻辑
 */
export const getSimpleReply = (input: string, botName: string = '智能助手'): string => {
  const lowerInput = input.toLowerCase();
  if (lowerInput.startsWith('帮助') || lowerInput.startsWith('help')) {
    return `# 帮助信息

我是${botName}，可以帮您：

## 基础功能
- **简单对话交流** - 和我聊天
- **HTML/Markdown支持** - 支持富文本内容

## 测试命令
- 输入"markdown"查看Markdown示例
- 输入"html"查看HTML示例
- 输入"table"查看表格示例
- 输入"code"查看代码块示例

有什么问题都可以问我哦！`.trim();
  }

  if (lowerInput.includes('markdown')) {
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
- 无序列表项2`.trim();
  }

  if (lowerInput.includes('html')) {
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
<a href="https://example.com" target="_blank">外部链接</a>`.trim();
  }

  if (lowerInput.includes('table')) {
    return `# 表格示例

| 功能 | 状态 | 说明 |
|------|------|------|
| HTML支持 | ✅ | 支持HTML标签渲染 |
| Markdown支持 | ✅ | 支持Markdown语法 |
| 表格支持 | ✅ | 支持表格渲染 |
| 代码高亮 | ✅ | 支持代码块高亮 |
| 任务列表 | ✅ | 支持任务列表 |

> 表格功能完全支持！`.trim();
  }

  if (lowerInput.includes('code')) {
    return `# 代码块示例

## JavaScript 代码
\`\`\`javascript
function fibonacci(n) {
  if (n <= 1) return n;
  return fibonacci(n - 1) + fibonacci(n - 2);
}

console.log(fibonacci(10));
\`\`\`

## Python 代码
\`\`\`python
def factorial(n):
    if n == 0:
        return 1
    return n * factorial(n - 1)

print(factorial(5))
\`\`\``.trim();
  }

  // 默认回复
  return `已收到您的消息："${input}"，我会尽快为您处理～`.trim();
};

/**
 * 模拟API延迟
 */
export const simulateApiDelay = (delay: number = 1000): Promise<void> => {
  return new Promise((resolve) => setTimeout(resolve, delay));
};

/**
 * 格式化时间戳
 */
export const formatTimestamp = (timestamp: number): string => {
  const date = new Date(timestamp);
  const now = new Date();
  // const diff = now.getTime() - timestamp;

  // 如果是今天
  if (date.toDateString() === now.toDateString()) {
    return date.toLocaleTimeString('zh-CN', {
      hour: '2-digit',
      minute: '2-digit',
    });
  }

  // 如果是昨天
  const yesterday = new Date(now);
  yesterday.setDate(yesterday.getDate() - 1);
  if (date.toDateString() === yesterday.toDateString()) {
    return `昨天 ${date.toLocaleTimeString('zh-CN', {
      hour: '2-digit',
      minute: '2-digit',
    })}`;
  }

  // 其他日期
  return date.toLocaleDateString('zh-CN');
};

<template>
  <div class="test-container">
    <h1>聊天机器人 HTML/Markdown 功能测试</h1>

    <div class="test-buttons">
      <button @click="testText">测试纯文本</button>
      <button @click="testHTML">测试HTML</button>
      <button @click="testMarkdown">测试Markdown</button>
      <button @click="testTable">测试表格</button>
      <button @click="testCode">测试代码块</button>
      <button @click="testThinking">测试思考中提示</button>
    </div>

    <ChatBot ref="chatBotRef" :options="chatBotOptions" />
  </div>
</template>

<script setup lang="ts">
  import { ref } from 'vue';
  import ChatBot from './ChatBot.vue';
  import type { ChatBotOptions } from './types';

  const chatBotRef = ref();

  const chatBotOptions: ChatBotOptions = {
    botName: '测试助手',
    placeholder: '输入消息进行测试...',
    replyDelay: 500,

    defaultReply: async (input: string) => {
      // 根据输入返回不同类型的测试内容
      if (input.includes('html')) {
        return `<h2 style="color: #667eea;">HTML 测试内容</h2>
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

      if (input.includes('markdown')) {
        return `# Markdown 测试内容

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
        return `# 表格测试

| 功能 | 状态 | 说明 |
|------|------|------|
| HTML支持 | ✅ | 支持HTML标签渲染 |
| Markdown支持 | ✅ | 支持Markdown语法 |
| 表格支持 | ✅ | 支持表格渲染 |
| 代码高亮 | ✅ | 支持代码块高亮 |
| 任务列表 | ✅ | 支持任务列表 |

> 表格功能完全支持！`;
      }

      if (input.includes('code')) {
        return `# 代码块测试

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
\`\`\`

## CSS 代码
\`\`\`css
.chat-bot {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 8px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}
\`\`\``;
      }

      // 测试思考中功能
      if (input.includes('复杂') || input.includes('分析')) {
        // 模拟一个需要较长时间处理的请求
        await new Promise((resolve) => setTimeout(resolve, 3000));
        return `# 复杂问题分析结果

经过深入分析，我为您提供了以下见解：

## 主要发现
1. **问题复杂度**: 这是一个多维度的问题，需要考虑多个因素
2. **解决方案**: 建议采用分步骤的方法来处理
3. **风险评估**: 需要仔细评估各种可能的风险

## 详细建议
- 首先进行数据收集和分析
- 然后制定详细的执行计划
- 最后进行实施和监控

> 这个分析花费了一些时间，但希望能为您提供有价值的参考！`;
      }

      return '这是普通文本回复。试试输入包含"html"、"markdown"、"table"或"code"的内容来测试不同格式！';
    },
  };

  // 测试函数
  const testText = () => {
    if (chatBotRef.value) {
      chatBotRef.value.addMessage('这是纯文本测试消息', true);
    }
  };

  const testHTML = () => {
    if (chatBotRef.value) {
      chatBotRef.value.addMessage('请显示HTML内容', true);
    }
  };

  const testMarkdown = () => {
    if (chatBotRef.value) {
      chatBotRef.value.addMessage('请显示Markdown内容', true);
    }
  };

  const testTable = () => {
    if (chatBotRef.value) {
      chatBotRef.value.addMessage('请显示表格', true);
    }
  };

  const testCode = () => {
    if (chatBotRef.value) {
      chatBotRef.value.addMessage('请显示代码块', true);
    }
  };

  const testThinking = () => {
    if (chatBotRef.value) {
      // 模拟一个需要较长时间处理的请求
      chatBotRef.value.addMessage('请帮我分析一个复杂的问题，这可能需要一些时间', true);
    }
  };
</script>

<style scoped>
  .test-container {
    max-width: 1200px;
    margin: 0 auto;
    padding: 20px;
  }

  .test-buttons {
    display: flex;
    flex-wrap: wrap;
    margin-bottom: 20px;
    gap: 10px;
  }

  .test-buttons button {
    padding: 8px 16px;
    transition: all 0.2s;
    border: 1px solid #ddd;
    border-radius: 4px;
    background: #f8f9fa;
    cursor: pointer;
  }

  .test-buttons button:hover {
    border-color: #adb5bd;
    background: #e9ecef;
  }
</style>

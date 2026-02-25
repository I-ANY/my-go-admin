<template>
  <div class="test-container">
    <h1>消息格式测试</h1>

    <div class="test-section">
      <h2>测试按钮</h2>
      <div class="test-buttons">
        <button @click="testTextWithNewlines">测试换行文本</button>
        <button @click="testHTML">测试HTML</button>
        <button @click="testMarkdown">测试Markdown</button>
        <button @click="testMixedContent">测试混合内容</button>
      </div>
    </div>

    <div class="test-section">
      <h2>ChatBot 组件</h2>
      <ChatBot ref="chatBotRef" :options="chatBotOptions" />
    </div>
  </div>
</template>

<script setup lang="ts">
  import { ref } from 'vue';
  import ChatBot from './ChatBot.vue';
  import type { ChatBotOptions } from './types';

  const chatBotRef = ref();

  const chatBotOptions: ChatBotOptions = {
    botName: '格式测试助手',
    placeholder: '输入消息进行测试...',
    replyDelay: 300,
    requireLogin: false, // 不需要登录

    defaultReply: async (input: string) => {
      if (input.includes('换行')) {
        return `这是第一行\n这是第二行\n这是第三行\n\n空行后的内容`;
      }

      if (input.includes('html')) {
        return `<h3>HTML 测试</h3>
<p>这是<strong>粗体</strong>和<em>斜体</em>文本。</p>
<ul>
  <li>列表项 1</li>
  <li>列表项 2</li>
  <li>列表项 3</li>
</ul>
<blockquote>这是一个引用块</blockquote>`;
      }

      if (input.includes('markdown')) {
        return `# Markdown 测试

## 文本格式
**粗体文本** 和 *斜体文本*
~~删除线文本~~

## 代码
行内代码：\`console.log('hello')\`

代码块：
\`\`\`javascript
function test() {
  console.log('Hello World');
}
\`\`\`

## 列表
1. 有序列表项
2. 另一个列表项

- 无序列表项
- 另一个无序列表项`;
      }

      if (input.includes('混合')) {
        return `# 混合内容测试

这是普通文本，包含换行符\n这是第二行\n这是第三行

**粗体文本** 和 *斜体文本*

\`\`\`javascript
// 代码块
function mixed() {
  return "混合内容";
}
\`\`\`

> 引用块中的内容

- 列表项1
- 列表项2`;
      }

      return `请测试以下格式：
1. 输入"换行"测试换行符
2. 输入"html"测试HTML格式
3. 输入"markdown"测试Markdown格式
4. 输入"混合"测试混合内容`;
    },
  };

  // 测试函数
  const testTextWithNewlines = () => {
    if (chatBotRef.value) {
      chatBotRef.value.addMessage('请显示带换行的文本', true);
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

  const testMixedContent = () => {
    if (chatBotRef.value) {
      chatBotRef.value.addMessage('请显示混合内容', true);
    }
  };
</script>

<style scoped>
  .test-container {
    max-width: 1200px;
    margin: 0 auto;
    padding: 20px;
  }

  .test-section {
    margin-bottom: 30px;
  }

  .test-section h2 {
    margin-bottom: 15px;
    color: #333;
  }

  .test-buttons {
    display: flex;
    flex-wrap: wrap;
    gap: 10px;
    margin-bottom: 20px;
  }

  .test-buttons button {
    padding: 10px 16px;
    transition: all 0.2s;
    border: 1px solid #ddd;
    border-radius: 6px;
    background: #f8f9fa;
    font-size: 14px;
    cursor: pointer;
  }

  .test-buttons button:hover {
    border-color: #667eea;
    background: #667eea;
    color: white;
  }
</style>

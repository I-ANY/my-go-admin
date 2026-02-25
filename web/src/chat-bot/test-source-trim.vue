<template>
  <div class="test-container">
    <h1>消息发送源头修复测试</h1>

    <div class="test-section">
      <h2>修复说明</h2>
      <div class="fix-info">
        <p><strong>修复内容：</strong></p>
        <ul>
          <li>在 <code>addMessage</code> 调用时使用 <code>trim()</code></li>
          <li>修复 <code>getSimpleReply</code> 中多行字符串的末尾换行</li>
          <li>确保所有机器人回复都没有末尾空白字符</li>
        </ul>
      </div>
    </div>

    <div class="test-section">
      <h2>测试按钮</h2>
      <div class="test-buttons">
        <button @click="testHelp">测试帮助命令</button>
        <button @click="testMarkdown">测试Markdown命令</button>
        <button @click="testHtml">测试HTML命令</button>
        <button @click="testTable">测试表格命令</button>
        <button @click="testCode">测试代码命令</button>
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
    botName: '源头修复测试助手',
    placeholder: '输入消息进行测试...',
    replyDelay: 300,
    requireLogin: false,

    defaultReply: async (input: string) => {
      // 测试多行字符串，确保没有末尾换行
      if (input.includes('测试')) {
        return `这是第一行
这是第二行
这是第三行
这是第四行`;
      }

      return `收到您的消息："${input}"
我会尽快为您处理～`;
    },
  };

  // 测试函数
  const testHelp = () => {
    if (chatBotRef.value) {
      chatBotRef.value.addMessage('帮助', true);
    }
  };

  const testMarkdown = () => {
    if (chatBotRef.value) {
      chatBotRef.value.addMessage('markdown', true);
    }
  };

  const testHtml = () => {
    if (chatBotRef.value) {
      chatBotRef.value.addMessage('html', true);
    }
  };

  const testTable = () => {
    if (chatBotRef.value) {
      chatBotRef.value.addMessage('table', true);
    }
  };

  const testCode = () => {
    if (chatBotRef.value) {
      chatBotRef.value.addMessage('code', true);
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

  .fix-info {
    padding: 15px;
    border-left: 4px solid #667eea;
    border-radius: 8px;
    background: #f8f9fa;
  }

  .fix-info ul {
    margin: 10px 0 0 20px;
  }

  .fix-info li {
    margin: 5px 0;
    line-height: 1.5;
  }

  .fix-info code {
    padding: 2px 4px;
    border-radius: 3px;
    background: #e9ecef;
    font-family: monospace;
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

<template>
  <div class="test-container">
    <h1>去除P标签包裹测试</h1>

    <div class="test-section">
      <h2>修复说明</h2>
      <div class="fix-info">
        <p><strong>修复内容：</strong></p>
        <ul>
          <li>只有包含连续换行符（段落分隔）时才使用 &lt;p&gt; 标签</li>
          <li>单行文本或只有单个换行符的文本不使用 &lt;p&gt; 标签</li>
          <li>避免消息看起来有多余的换行</li>
        </ul>
      </div>
    </div>

    <div class="test-section">
      <h2>测试按钮</h2>
      <div class="test-buttons">
        <button @click="testSingleLine">测试单行文本</button>
        <button @click="testSingleNewline">测试单个换行</button>
        <button @click="testMultipleNewlines">测试多个换行</button>
        <button @click="testParagraphs">测试段落分隔</button>
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
    botName: 'P标签修复测试助手',
    placeholder: '输入消息进行测试...',
    replyDelay: 300,
    requireLogin: false,

    defaultReply: async (input: string) => {
      if (input.includes('单行')) {
        return '这是一行简单的文本消息';
      }

      if (input.includes('单个换行')) {
        return '这是第一行\n这是第二行';
      }

      if (input.includes('多个换行')) {
        return '这是第一行\n这是第二行\n这是第三行';
      }

      if (input.includes('段落')) {
        return '这是第一段\n\n这是第二段\n\n这是第三段';
      }

      return '请测试以下格式：\n1. 输入"单行"测试单行文本\n2. 输入"单个换行"测试单个换行\n3. 输入"多个换行"测试多个换行\n4. 输入"段落"测试段落分隔';
    },
  };

  // 测试函数
  const testSingleLine = () => {
    if (chatBotRef.value) {
      chatBotRef.value.addMessage('请显示单行文本', true);
    }
  };

  const testSingleNewline = () => {
    if (chatBotRef.value) {
      chatBotRef.value.addMessage('请显示单个换行', true);
    }
  };

  const testMultipleNewlines = () => {
    if (chatBotRef.value) {
      chatBotRef.value.addMessage('请显示多个换行', true);
    }
  };

  const testParagraphs = () => {
    if (chatBotRef.value) {
      chatBotRef.value.addMessage('请显示段落分隔', true);
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

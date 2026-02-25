<template>
  <div class="test-container">
    <h1>内容类型检测优先级测试</h1>

    <div class="test-section">
      <h2>优先级说明</h2>
      <div class="priority-info">
        <p><strong>优先级顺序：</strong></p>
        <ol>
          <li><strong>Markdown 语法</strong> - 优先处理为 markdown 类型</li>
          <li><strong>换行符 (\n)</strong> - 其次处理为 text 类型</li>
        </ol>
      </div>
    </div>

    <div class="test-section">
      <h2>测试按钮</h2>
      <div class="test-buttons">
        <button @click="testNewlinePriority">测试换行符优先级</button>
        <button @click="testMarkdownPriority">测试Markdown优先级</button>
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
    botName: '优先级测试助手',
    placeholder: '输入消息进行测试...',
    replyDelay: 300,
    requireLogin: false,

    defaultReply: async (input: string) => {
      if (input.includes('换行')) {
        return `这是第一行\n这是第二行\n这是第三行\n\n空行后的内容`;
      }

      if (input.includes('markdown')) {
        return `# Markdown 测试\n\n**粗体文本** 和 *斜体文本*\n\n- 列表项1\n- 列表项2`;
      }

      if (input.includes('混合')) {
        return `# 混合内容测试\n\n这是普通文本，包含换行符\n这是第二行\n这是第三行\n\n**粗体文本** 和 *斜体文本*\n\n<blockquote>引用块中的内容</blockquote>\n\n- 列表项1\n- 列表项2`;
      }

      return `请测试以下格式：\n1. 输入"换行"测试换行符优先级\n2. 输入"markdown"测试Markdown优先级\n3. 输入"混合"测试混合内容`;
    },
  };

  // 测试函数
  const testNewlinePriority = () => {
    if (chatBotRef.value) {
      chatBotRef.value.addMessage('请显示带换行的文本', true);
    }
  };

  const testMarkdownPriority = () => {
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

  .priority-info {
    padding: 15px;
    border-left: 4px solid #667eea;
    border-radius: 8px;
    background: #f8f9fa;
  }

  .priority-info ol {
    margin: 10px 0 0 20px;
  }

  .priority-info li {
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

<template>
  <div class="test-container">
    <h1>换行符优化测试</h1>

    <div class="test-section">
      <h2>优化说明</h2>
      <div class="optimization-info">
        <p><strong>优化内容：</strong></p>
        <ul>
          <li><code>\n\n</code> - 连续换行符转换为段落分隔</li>
          <li><code>\n</code> - 单个换行符转换为 &lt;br&gt; 标签</li>
          <li>自动添加 &lt;p&gt; 标签包裹内容</li>
        </ul>
      </div>
    </div>

    <div class="test-section">
      <h2>测试按钮</h2>
      <div class="test-buttons">
        <button @click="testSingleNewline">测试单个换行</button>
        <button @click="testDoubleNewline">测试连续换行</button>
        <button @click="testMixedNewlines">测试混合换行</button>
        <button @click="testMarkdownWithNewlines">测试Markdown换行</button>
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
    botName: '换行优化测试助手',
    placeholder: '输入消息进行测试...',
    replyDelay: 300,
    requireLogin: false,

    defaultReply: async (input: string) => {
      if (input.includes('单个')) {
        return `这是第一行\n这是第二行\n这是第三行`;
      }

      if (input.includes('连续')) {
        return `这是第一段\n\n这是第二段\n\n这是第三段`;
      }

      if (input.includes('混合')) {
        return `这是第一行\n这是第二行\n\n这是新段落\n这是新段落的第二行\n\n这是最后一个段落`;
      }

      if (input.includes('markdown')) {
        return `# Markdown 标题\n\n这是第一段内容\n这是第一段的第二行\n\n## 子标题\n\n这是第二段内容\n\n- 列表项1\n- 列表项2`;
      }

      return `请测试以下格式：\n1. 输入"单个"测试单个换行\n2. 输入"连续"测试连续换行\n3. 输入"混合"测试混合换行\n4. 输入"markdown"测试Markdown换行`;
    },
  };

  // 测试函数
  const testSingleNewline = () => {
    if (chatBotRef.value) {
      chatBotRef.value.addMessage('请显示单个换行', true);
    }
  };

  const testDoubleNewline = () => {
    if (chatBotRef.value) {
      chatBotRef.value.addMessage('请显示连续换行', true);
    }
  };

  const testMixedNewlines = () => {
    if (chatBotRef.value) {
      chatBotRef.value.addMessage('请显示混合换行', true);
    }
  };

  const testMarkdownWithNewlines = () => {
    if (chatBotRef.value) {
      chatBotRef.value.addMessage('请显示Markdown换行', true);
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

  .optimization-info {
    padding: 15px;
    border-left: 4px solid #667eea;
    border-radius: 8px;
    background: #f8f9fa;
  }

  .optimization-info ul {
    margin: 10px 0 0 20px;
  }

  .optimization-info li {
    margin: 5px 0;
    line-height: 1.5;
  }

  .optimization-info code {
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

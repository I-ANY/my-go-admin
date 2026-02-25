<template>
  <div class="test-container">
    <h1>表格显示测试</h1>

    <div class="test-section">
      <h2>测试说明</h2>
      <div class="fix-info">
        <p><strong>测试内容：</strong></p>
        <ul>
          <li>测试简单表格显示</li>
          <li>测试表格HTML标签是否正确</li>
          <li>验证Markdown转换是否正常</li>
        </ul>
      </div>
    </div>

    <div class="test-section">
      <h2>测试按钮</h2>
      <div class="test-buttons">
        <button @click="testSimpleTable">测试简单表格</button>
        <button @click="testComplexTable">测试复杂表格</button>
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
    botName: '表格显示测试助手',
    placeholder: '输入消息进行测试...',
    replyDelay: 300,
    requireLogin: false,

    defaultReply: async (input: string) => {
      if (input.includes('简单')) {
        return `| 姓名 | 年龄 | 职业 |
|------|------|------|
| 张三 | 25 | 工程师 |
| 李四 | 30 | 设计师 |
| 王五 | 28 | 产品经理 |`;
      }

      if (input.includes('复杂')) {
        return `# 复杂表格示例

| 项目 | 状态 | 进度 | 负责人 | 截止日期 |
|------|------|------|--------|----------|
| 功能开发 | 进行中 | 75% | 张三 | 2024-01-15 |
| 测试验证 | 待开始 | 0% | 李四 | 2024-01-20 |
| 文档编写 | 已完成 | 100% | 王五 | 2024-01-10 |

> 项目进度表`;
      }

      if (input.includes('混合')) {
        return `**项目状态报告**

| 模块 | 完成度 | 问题数 | 优先级 |
|------|--------|--------|--------|
| 用户管理 | 90% | 2 | 高 |
| 权限控制 | 60% | 5 | 中 |
| 数据统计 | 30% | 8 | 低 |

- 总体进度：60%
- 预计完成时间：2024年2月`;
      }

      return `请测试以下格式：\n1. 输入"简单"测试简单表格\n2. 输入"复杂"测试复杂表格\n3. 输入"混合"测试混合内容`;
    },
  };

  // 测试函数
  const testSimpleTable = () => {
    if (chatBotRef.value) {
      chatBotRef.value.addMessage('请显示简单表格', true);
    }
  };

  const testComplexTable = () => {
    if (chatBotRef.value) {
      chatBotRef.value.addMessage('请显示复杂表格', true);
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

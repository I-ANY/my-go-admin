<template>
  <div class="test-container">
    <h1>Markdown表格修复测试</h1>

    <div class="test-section">
      <h2>修复说明</h2>
      <div class="fix-info">
        <p><strong>修复内容：</strong></p>
        <ul>
          <li>清理HTML中的多余空白字符</li>
          <li>移除表格前面的多余空行</li>
          <li>确保表格显示紧凑，没有多余换行</li>
        </ul>
      </div>
    </div>

    <div class="test-section">
      <h2>测试按钮</h2>
      <div class="test-buttons">
        <button @click="testTable">测试表格</button>
        <button @click="testTableWithText">测试表格+文本</button>
        <button @click="testMixedContent">测试混合内容</button>
        <button @click="testComplexTable">测试复杂表格</button>
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
    botName: '表格修复测试助手',
    placeholder: '输入消息进行测试...',
    replyDelay: 300,
    requireLogin: false,

    defaultReply: async (input: string) => {
      if (input.includes('表格')) {
        return `# 测试表格

| 功能 | 状态 | 说明 |
|------|------|------|
| HTML支持 | ✅ | 支持HTML标签渲染 |
| Markdown支持 | ✅ | 支持Markdown语法 |
| 表格支持 | ✅ | 支持表格渲染 |
| 代码高亮 | ✅ | 支持代码块高亮 |
| 任务列表 | ✅ | 支持任务列表 |

> 表格功能完全支持！`;
      }

      if (input.includes('文本')) {
        return `这是表格前面的文本

| 列1 | 列2 | 列3 |
|-----|-----|-----|
| 数据1 | 数据2 | 数据3 |
| 数据4 | 数据5 | 数据6 |

这是表格后面的文本`;
      }

      if (input.includes('混合')) {
        return `# 混合内容测试

**粗体文本** 和 *斜体文本*

| 项目 | 状态 | 描述 |
|------|------|------|
| 功能A | 完成 | 基础功能实现 |
| 功能B | 进行中 | 正在开发中 |
| 功能C | 计划中 | 待开发 |

- 列表项1
- 列表项2
- 列表项3`;
      }

      if (input.includes('复杂')) {
        return `# 复杂表格示例

## 用户信息表

| 用户ID | 姓名 | 年龄 | 邮箱 | 状态 |
|--------|------|------|------|------|
| 001 | 张三 | 25 | zhangsan@example.com | 活跃 |
| 002 | 李四 | 30 | lisi@example.com | 离线 |
| 003 | 王五 | 28 | wangwu@example.com | 活跃 |

## 统计信息

| 指标 | 数值 | 单位 | 趋势 |
|------|------|------|------|
| 用户数 | 1000 | 人 | ↗️ |
| 活跃度 | 85% | 百分比 | ↗️ |
| 响应时间 | 200ms | 毫秒 | ↘️ |

> 数据统计完成`;
      }

      return `请测试以下格式：\n1. 输入"表格"测试简单表格\n2. 输入"文本"测试表格+文本\n3. 输入"混合"测试混合内容\n4. 输入"复杂"测试复杂表格`;
    },
  };

  // 测试函数
  const testTable = () => {
    if (chatBotRef.value) {
      chatBotRef.value.addMessage('请显示表格', true);
    }
  };

  const testTableWithText = () => {
    if (chatBotRef.value) {
      chatBotRef.value.addMessage('请显示表格+文本', true);
    }
  };

  const testMixedContent = () => {
    if (chatBotRef.value) {
      chatBotRef.value.addMessage('请显示混合内容', true);
    }
  };

  const testComplexTable = () => {
    if (chatBotRef.value) {
      chatBotRef.value.addMessage('请显示复杂表格', true);
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

<template>
  <div
    class="chat-bot-widget"
    :class="{ collapsed: isCollapsed, zoomed: isZoomed }"
    :style="widgetStyle"
    v-if="isUserLoggedIn && !isClosed"
  >
    <!-- 聊天窗口（展开状态显示） -->
    <div class="chat-window" v-show="!isCollapsed">
      <!-- 窗口头部 -->
      <div class="window-header">
        <div class="bot-info">
          <div class="bot-avatar">
            <div class="bot-face-small">
              <div class="bot-eyes-small">
                <div class="eye-small"></div>
                <div class="eye-small"></div>
              </div>
            </div>
          </div>
          <div class="bot-details">
            <div class="bot-name">{{ botName }}</div>
            <div class="bot-status">在线</div>
          </div>
        </div>
        <div class="window-controls">
          <button
            class="control-btn zoom-btn"
            @click="toggleZoom"
            :title="isZoomed ? '还原' : '放大'"
          >
            <svg v-if="!isZoomed" width="12" height="12" viewBox="0 0 12 12">
              <path d="M5 2v3H2v2h3v3h2V7h3V5H7V2H5z" fill="currentColor" />
            </svg>
            <svg v-else width="12" height="12" viewBox="0 0 12 12">
              <path d="M2 5h8v2H2V5z" fill="currentColor" />
            </svg>
          </button>
          <button
            class="control-btn minimize-btn"
            @click="toggleCollapse"
            title="最小化"
            v-show="!isZoomed"
          >
            <svg width="12" height="12" viewBox="0 0 12 12">
              <path d="M2 6h8" stroke="currentColor" stroke-width="1.5" fill="none" />
            </svg>
          </button>
          <!-- 窗口上的关闭按钮已隐藏，只保留头像上的关闭按钮 -->
        </div>
      </div>

      <!-- 消息区域 -->
      <div ref="messagesContainer" class="messages-container">
        <div
          v-for="msg in messages"
          :key="msg.id"
          :class="['message-item', { 'user-message': msg.isUser }]"
        >
          <div class="message-content">
            <div v-if="msg.type === 'text'" class="text-content selectable-text">
              <div v-html="renderMessageContent(msg)"></div>
            </div>

            <div
              v-else-if="msg.type === 'markdown'"
              class="markdown-content selectable-text"
              v-html="renderMessageContent(msg)"
            ></div>
            <div
              v-else-if="msg.type === 'image'"
              class="image-content"
              v-html="renderMessageContent(msg)"
            ></div>
            <div v-else-if="msg.type === 'loading'" class="loading-content">
              <div class="thinking-text">{{ msg.content }}</div>
              <div class="thinking-dots"> <span></span><span></span><span></span> </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 输入区域 -->
      <div class="input-container">
        <textarea
          v-model="inputText"
          :placeholder="placeholder"
          @keydown="handleKeydown"
          rows="1"
        ></textarea>
        <button @click="handleSend" :disabled="isLoading" class="send-btn">
          <svg v-if="!isLoading" width="16" height="16" viewBox="0 0 16 16">
            <path d="M1 8l7-7v4.5c4.5 0 8 3.5 8 8-4.5 0-8 3.5-8 8v-4.5L1 8z" fill="currentColor" />
          </svg>
          <div v-else class="loading-spinner"></div>
        </button>
      </div>
    </div>

    <!-- 机器人头部（固定在右下角） -->
    <div
      class="bot-head"
      @dblclick="toggleCollapse"
      @mousedown="startDrag"
      @touchstart="startDrag"
      :class="{
        'has-notification': hasUnreadMessages,
        dragging: isDragging && isCollapsed,
      }"
      v-show="!isZoomed"
    >
      <div class="bot-face">
        <div class="bot-eyes">
          <div class="eye left-eye"></div>
          <div class="eye right-eye"></div>
        </div>
        <div class="bot-mouth"></div>
      </div>
      <div class="notification-dot" v-if="hasUnreadMessages"></div>

      <!-- 头像上的关闭按钮 -->
      <button
        v-if="showCloseButton"
        class="bot-close-btn"
        @click.stop="handleClose"
        title="关闭机器人"
      >
        <svg width="14" height="14" viewBox="0 0 14 14">
          <path d="M2 2l10 10m0-10l-10 10" stroke="currentColor" stroke-width="2" fill="none" />
        </svg>
      </button>

      <!-- 打开提示 -->
      <div class="zoom-hint" v-if="isCollapsed">
        <span>双击打开|关闭窗口</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
  import { ref, computed, onMounted, watch, nextTick, readonly } from 'vue';
  import { useUserStore } from '@/store/modules/user';
  import type { MessageItem, ChatBotPluginOptions, BotMessageType } from './types';
  import { generateMessageId, getSimpleReply } from './chat-bot';
  import { renderMessageContent, detectContentType } from './utils';

  // 定义 props
  const props = defineProps<{
    options: ChatBotPluginOptions;
  }>();

  // 用户store
  const userStore = useUserStore();

  // 响应式状态
  const messages = ref<MessageItem[]>([]);
  const inputText = ref('');
  const messagesContainer = ref<HTMLElement | null>(null);
  const isLoading = ref(false);
  const isCollapsed = ref(true); // 默认收起状态
  const hasUnreadMessages = ref(false);
  const isZoomed = ref(false);
  const isFirstOpen = ref(true); // 标记是否首次打开
  const isClosed = ref(false); // 标记是否已关闭

  // 拖动相关状态
  const isDragging = ref(false);
  const dragStartX = ref(0);
  const dragStartY = ref(0);
  const widgetPosition = ref({ x: 20, y: 20 }); // 默认位置（右下角）
  const initialPosition = ref({ x: 0, y: 0 }); // 拖动开始时的位置
  const dragStartTime = ref(0); // 记录拖动开始时间，用于区分单击和双击

  // 计算属性
  const botName = computed(() => props.options.botName || '智能助手');
  // const botAvatar = computed(
  //   () => props.options.botAvatar || 'https://via.placeholder.com/40x40?text=Bot',
  // );
  // const userAvatar = computed(
  //   () => props.options.userAvatar || 'https://via.placeholder.com/40x40?text=User',
  // );
  const placeholder = computed(
    () => props.options.placeholder || '输入消息... (Enter发送，Shift+Enter换行)',
  );
  const replyDelay = computed(() => props.options.replyDelay || 1000);
  const showCloseButton = computed(() => props.options.showCloseButton !== false); // 默认显示关闭按钮

  // 计算widget样式（包含位置）
  const widgetStyle = computed(() => {
    if (isZoomed.value) {
      return {}; // 放大时使用默认定位
    }
    return {
      right: `${widgetPosition.value.x}px`,
      bottom: `${widgetPosition.value.y}px`,
      willChange: isDragging.value ? 'right, bottom' : 'auto', // 优化拖动性能
      transform: isDragging.value ? 'translateZ(0)' : 'none', // 启用硬件加速
      transition: isDragging.value ? 'none' : 'all 0.1s ease', // 拖动时禁用过渡
    };
  });

  // 检查用户是否已登录
  const isUserLoggedIn = computed(() => {
    // 如果配置为不需要登录，则直接返回true
    if (props.options.requireLogin === false) {
      return true;
    }

    const token = userStore.getToken;
    const isExpired = userStore.getTokenIsExpire;
    return token && !isExpired;
  });

  // 初始化默认消息 - 移除自动显示欢迎消息
  onMounted(() => {
    // 默认不显示欢迎消息，等用户第一次点击时再显示
  });

  // 拖动相关方法
  const startDrag = (event: MouseEvent | TouchEvent) => {
    // 如果窗体已打开，不允许拖动
    if (!isCollapsed.value) {
      return;
    }

    // 记录拖动开始时间
    dragStartTime.value = Date.now();

    event.preventDefault();
    isDragging.value = true;

    const clientX = 'touches' in event ? event.touches[0].clientX : event.clientX;
    const clientY = 'touches' in event ? event.touches[0].clientY : event.clientY;

    // 记录拖动开始时的位置
    initialPosition.value = { ...widgetPosition.value };

    // 计算鼠标相对于头像的偏移
    dragStartX.value = clientX - (window.innerWidth - widgetPosition.value.x);
    dragStartY.value = clientY - (window.innerHeight - widgetPosition.value.y);

    // 添加全局事件监听
    document.addEventListener('mousemove', onDrag);
    document.addEventListener('touchmove', onDrag, { passive: false });
    document.addEventListener('mouseup', stopDrag);
    document.addEventListener('touchend', stopDrag);
  };

  const onDrag = (event: MouseEvent | TouchEvent) => {
    if (!isDragging.value) return;

    event.preventDefault();
    const clientX = 'touches' in event ? event.touches[0].clientX : event.clientX;
    const clientY = 'touches' in event ? event.touches[0].clientY : event.clientY;

    // 直接计算头像位置，让头像跟随鼠标
    const newX = window.innerWidth - clientX + dragStartX.value;
    const newY = window.innerHeight - clientY + dragStartY.value;

    // 限制在视窗范围内
    const maxX = window.innerWidth - 60; // 60px是头像宽度
    const maxY = window.innerHeight - 60; // 60px是头像高度

    // 直接更新位置，确保流畅跟随
    widgetPosition.value = {
      x: Math.max(0, Math.min(newX, maxX)),
      y: Math.max(0, Math.min(newY, maxY)),
    };
  };

  const stopDrag = () => {
    // 检查是否是快速点击（可能是双击的一部分），如果是则不执行拖动
    const dragDuration = Date.now() - dragStartTime.value;
    if (dragDuration < 200) {
      // 200ms内的点击可能是双击的一部分
      isDragging.value = false;
      return;
    }

    isDragging.value = false;

    // 移除全局事件监听
    document.removeEventListener('mousemove', onDrag);
    document.removeEventListener('touchmove', onDrag);
    document.removeEventListener('mouseup', stopDrag);
    document.removeEventListener('touchend', stopDrag);
  };

  // 切换收起/展开状态
  const toggleCollapse = () => {
    // 检查登录状态
    if (!isUserLoggedIn.value) {
      return;
    }

    isCollapsed.value = !isCollapsed.value;
    if (!isCollapsed.value) {
      hasUnreadMessages.value = false;
      // 如果是首次打开，显示欢迎消息
      if (isFirstOpen.value) {
        isFirstOpen.value = false;
        addMessage('你好！我是你的智能问答助手，有什么可以帮你的吗？\n help 查看帮助', false);
      }
      nextTick(() => {
        scrollToBottom();
      });
    }
  };

  // 切换缩放状态
  const toggleZoom = () => {
    // 检查登录状态
    if (!isUserLoggedIn.value) {
      return;
    }

    isZoomed.value = !isZoomed.value;
    if (isZoomed.value && isCollapsed.value) {
      // 如果是在最小化状态下放大，先展开聊天窗口
      isCollapsed.value = false;
      hasUnreadMessages.value = false;
      // 如果是首次打开，显示欢迎消息
      if (isFirstOpen.value) {
        isFirstOpen.value = false;
        addMessage('你好！我是你的智能问答助手，有什么可以帮你的吗？\n help 查看帮助', false);
      }
      nextTick(() => {
        scrollToBottom();
      });
    }
  };

  // 处理关闭机器人
  const handleClose = () => {
    // 触发关闭回调
    if (props.options.onClose) {
      props.options.onClose();
    }

    // 完全关闭机器人
    isClosed.value = true;

    // 重置状态
    isCollapsed.value = true;
    isZoomed.value = false;
    hasUnreadMessages.value = false;
    messages.value = [];
    inputText.value = '';
    isLoading.value = false;
    isFirstOpen.value = true;
  };

  // 重新打开机器人
  const reopenBot = () => {
    isClosed.value = false;

    // 触发重新打开回调
    if (props.options.onReopen) {
      props.options.onReopen();
    }
  };

  // 添加消息
  const addMessage = (content: string, isUser: boolean, type?: BotMessageType): MessageItem => {
    // 如果没有指定类型，自动检测内容类型
    const detectedType = type || detectContentType(content);

    const newMsg: MessageItem = {
      id: generateMessageId(),
      content,
      type: detectedType,
      isUser,
      timestamp: Date.now(),
      renderOptions: {
        allowHtml: detectedType === 'markdown',
        markdownOptions:
          detectedType === 'markdown'
            ? {
                tables: true,
                strikethrough: true,
                tasklists: true,
                ghCodeBlocks: true,
                emoji: true,
                underline: true,
                simpleLineBreaks: true,
                ghCompatibleHeaderId: true,
                parseImgDimensions: true,
                simplifiedAutoLink: true,
              }
            : undefined,
      },
    };
    messages.value.push(newMsg);

    // 如果是机器人消息且当前是收起状态，显示未读提示
    if (!isUser && isCollapsed.value) {
      hasUnreadMessages.value = true;
    }

    // 触发消息回调
    if (props.options.onMessage) {
      props.options.onMessage(newMsg);
    }

    // 自动滚动到底部
    if (!isCollapsed.value) {
      scrollToBottom();
    }

    return newMsg;
  };

  // 处理键盘事件
  const handleKeydown = (event: KeyboardEvent) => {
    if (event.key === 'Enter') {
      if (event.shiftKey) {
        // Shift+Enter: 换行
        return;
      } else {
        // Enter: 发送消息
        event.preventDefault();
        handleSend();
      }
    }
  };

  // 处理发送消息
  const handleSend = async () => {
    // 检查登录状态
    if (!isUserLoggedIn.value) {
      return;
    }

    const input = inputText.value.trim();
    if (!input || isLoading.value) return;

    // 添加用户消息
    addMessage(input, true);

    // 清空输入框
    inputText.value = '';
    isLoading.value = true;

    // 添加"正在思考中"的提示消息
    const thinkingMessage = addMessage('正在思考中...', false, 'loading');

    try {
      // 获取机器人回复
      const replyContent = await getBotReply(input);

      // 移除思考中的消息
      const thinkingIndex = messages.value.findIndex((msg) => msg.id === thinkingMessage.id);
      if (thinkingIndex !== -1) {
        messages.value.splice(thinkingIndex, 1);
      }

      // 添加机器人回复（带延迟效果）
      setTimeout(() => {
        addMessage(replyContent.trim(), false);
        isLoading.value = false;
      }, replyDelay.value);
    } catch (error) {
      console.error('机器人回复失败:', error);

      // 移除思考中的消息
      const thinkingIndex = messages.value.findIndex((msg) => msg.id === thinkingMessage.id);
      if (thinkingIndex !== -1) {
        messages.value.splice(thinkingIndex, 1);
      }

      setTimeout(() => {
        addMessage('抱歉，我现在无法回答您的问题'.trim(), false);
        isLoading.value = false;
      }, replyDelay.value);
    }
  };

  // 获取机器人回复（核心逻辑）
  const getBotReply = async (input: string): Promise<string> => {
    // 优先使用用户自定义回复逻辑
    if (props.options.defaultReply) {
      return props.options.defaultReply(input);
    }

    // 使用工具函数中的回复逻辑
    return getSimpleReply(input, botName.value);
  };

  // 滚动到底部
  const scrollToBottom = () => {
    nextTick(() => {
      if (messagesContainer.value) {
        messagesContainer.value.scrollTop = messagesContainer.value.scrollHeight;
      }
    });
  };

  // 自动调整文本框高度
  const adjustTextareaHeight = () => {
    const textarea = document.querySelector('.input-container textarea') as HTMLTextAreaElement;
    if (textarea) {
      textarea.style.height = 'auto';
      textarea.style.height = Math.min(textarea.scrollHeight, 100) + 'px';
    }
  };

  // 监听输入变化
  watch(inputText, () => {
    nextTick(() => {
      adjustTextareaHeight();
    });
  });

  // 暴露方法给父组件
  defineExpose({
    reopenBot,
    addMessage,
    isClosed: readonly(isClosed),
  });
</script>

<style scoped>
  @keyframes slide-in {
    from {
      transform: translateY(20px) scale(0.9);
      opacity: 0;
    }

    to {
      transform: translateY(0) scale(1);
      opacity: 1;
    }
  }

  @keyframes loading {
    0%,
    100% {
      opacity: 1;
    }

    50% {
      opacity: 0.3;
    }
  }

  @keyframes spin {
    0% {
      transform: rotate(0deg);
    }

    100% {
      transform: rotate(360deg);
    }
  }

  @keyframes pulse-glow {
    0%,
    100% {
      box-shadow:
        0 4px 12px rgb(0 0 0 / 15%),
        0 0 0 0 rgb(239 68 68 / 70%);
    }

    50% {
      box-shadow:
        0 4px 12px rgb(0 0 0 / 15%),
        0 0 0 10px rgb(239 68 68 / 0%);
    }
  }

  @keyframes blink {
    0%,
    90%,
    100% {
      opacity: 1;
    }

    95% {
      opacity: 0;
    }
  }

  @keyframes pulse {
    0% {
      transform: scale(1);
    }

    50% {
      transform: scale(1.2);
    }

    100% {
      transform: scale(1);
    }
  }

  /* 响应式设计 */
  @media (max-width: 480px) {
    .chat-bot-widget {
      right: 10px;
      bottom: 10px;
    }

    .chat-window {
      right: -10px;
      width: calc(100vw - 20px);
    }

    .chat-bot-widget.collapsed {
      transform: translateX(calc(100vw - 70px));
    }

    .chat-bot-widget.collapsed.zoomed {
      transform: translateX(calc(100vw - 93px)) scale(1.33);
    }

    .chat-bot-widget.zoomed {
      bottom: 20%;
      left: 50%;
      transform: translateX(-50%);
    }

    .chat-bot-widget.zoomed .chat-window {
      width: calc(100vw - 20px);
      max-width: none;
      height: calc(100vh - 100px);
      max-height: none;
    }

    /* 移动端关闭按钮样式调整 */
    .bot-close-btn {
      top: -4px;
      right: -4px;
      width: 24px;
      height: 24px;
      border-width: 2px;
    }

    .bot-close-btn svg {
      width: 12px;
      height: 12px;
    }
  }

  /* 全屏模式 */
  @media (min-width: 1200px) {
    .chat-bot-widget.zoomed .chat-window {
      width: 60vw;
      max-width: 1000px;
      height: 80vh;
      max-height: 800px;
    }
  }

  .chat-bot-widget {
    position: fixed;
    z-index: 1000;
    right: 20px;
    bottom: 20px;
    transition: all 0.15s ease; /* 缩短过渡时间 */
    user-select: none; /* 防止拖动时选中文本 */
  }

  .chat-bot-widget.dragging {
    transition: none; /* 拖动时禁用过渡动画 */
  }

  .chat-bot-widget.collapsed {
    transform: translateX(calc(100% - 60px));
  }

  .chat-bot-widget.collapsed.zoomed {
    transform: translateX(calc(100% - 80px)) scale(1.33);
  }

  .chat-bot-widget.zoomed {
    position: fixed;
    z-index: 1000;
    bottom: 20%;
    left: 50%;
    transform: translateX(-50%);
  }

  /* 聊天窗口样式 */
  .chat-window {
    display: flex;
    position: absolute;
    right: 0;
    bottom: 70px;
    flex-direction: column;
    width: 320px;
    height: 450px;
    overflow: hidden;
    animation: slide-in 0.15s ease; /* 缩短动画时间 */
    border-radius: 16px;
    background: white;
    box-shadow: 0 8px 32px rgb(0 0 0 / 12%);
  }

  .chat-bot-widget.zoomed .chat-window {
    position: relative;
    right: auto;
    bottom: auto;
    width: 60vw;
    max-width: 1000px;
    height: 80vh;
    max-height: 800px;
  }

  /* 窗口头部 */
  .window-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 16px;
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    color: white;
  }

  .bot-info {
    display: flex;
    align-items: center;
    gap: 12px;
  }

  .bot-avatar {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 32px;
    height: 32px;
    border-radius: 50%;
    background: rgb(255 255 255 / 20%);
  }

  .bot-face-small {
    position: relative;
    width: 20px;
    height: 20px;
  }

  .bot-eyes-small {
    display: flex;
    justify-content: space-between;
    margin-bottom: 4px;
  }

  .eye-small {
    width: 4px;
    height: 4px;
    border-radius: 50%;
    background: white;
  }

  .bot-details {
    flex: 1;
  }

  .bot-name {
    margin-bottom: 2px;
    font-size: 14px;
    font-weight: 600;
  }

  .bot-status {
    opacity: 0.8;
    font-size: 12px;
  }

  .window-controls {
    display: flex;
    gap: 8px;
  }

  .control-btn {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 24px;
    height: 24px;
    transition: background 0.1s; /* 缩短按钮过渡时间 */
    border: none;
    border-radius: 50%;
    background: rgb(255 255 255 / 20%);
    color: white;
    cursor: pointer;
  }

  .control-btn:hover {
    background: rgb(255 255 255 / 30%);
  }

  .close-btn:hover {
    background: rgb(239 68 68 / 80%) !important;
  }

  /* 消息区域 */
  .messages-container {
    display: flex;
    flex: 1;
    flex-direction: column;
    padding: 16px;
    overflow-y: auto;
    background: #f8fafc;
    gap: 12px;
  }

  .message-item {
    display: flex;
    max-width: 85%;
  }

  .user-message {
    align-self: flex-end;
  }

  .message-content {
    padding: 8px 12px;
    border-radius: 12px;
    font-size: 14px;
    line-height: 1.4;
  }

  .user-message .message-content {
    border-bottom-right-radius: 4px;
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    color: white;
  }

  .message-item:not(.user-message) .message-content {
    border: 1px solid #e5e7eb;
    border-bottom-left-radius: 4px;
    background: white;
    color: #374151;
  }

  /* 可选择的文本样式 */
  .selectable-text {
    position: relative;
    cursor: text;
    user-select: text;
  }

  .selectable-text pre {
    margin: 0;
    padding: 0;
    font-family: inherit;
    font-size: inherit;
    line-height: inherit;
    word-wrap: break-word;
    white-space: pre-wrap;
    overflow-wrap: break-word;
  }

  /* 文本内容样式 */
  .text-content {
    margin: 0;
    padding: 0;
    line-height: 1.5;
    word-wrap: break-word;
    overflow-wrap: break-word;
  }

  .text-content p {
    margin: 8px 0;
  }

  .text-content p:first-child {
    margin-top: 0;
  }

  .text-content p:last-child {
    margin-bottom: 0;
  }

  /* 确保没有p标签的内容也能正确显示 */

  /* .text-content {

  } */

  .selectable-text::selection {
    background: rgb(102 126 234 / 30%);
    color: #374151;
  }

  /* 文本选择时的样式 */
  .selectable-text:focus {
    outline: none;
  }

  /* Markdown内容样式 */
  .markdown-content {
    line-height: 1.6;
    word-wrap: break-word;
    overflow-wrap: break-word;
  }

  .markdown-content h1,
  .markdown-content h2,
  .markdown-content h3,
  .markdown-content h4,
  .markdown-content h5,
  .markdown-content h6 {
    margin: 12px 0 8px;
    font-weight: 600;
    line-height: 1.4;
  }

  .markdown-content h1 {
    padding-bottom: 4px;
    border-bottom: 1px solid #e5e7eb;
    font-size: 1.5em;
  }

  .markdown-content h2 {
    padding-bottom: 3px;
    border-bottom: 1px solid #e5e7eb;
    font-size: 1.3em;
  }

  .markdown-content h3 {
    font-size: 1.1em;
  }

  .markdown-content p {
    margin: 8px 0;
  }

  .markdown-content ul,
  .markdown-content ol {
    margin: 8px 0;
    padding-left: 20px;
  }

  .markdown-content li {
    margin: 4px 0;
  }

  .markdown-content blockquote {
    margin: 8px 0;
    padding: 8px 12px;
    border-left: 4px solid #667eea;
    background: #f8fafc;
    color: #6b7280;
  }

  .markdown-content code {
    padding: 2px 6px;
    border-radius: 4px;
    background: #f3f4f6;
    font-family: Monaco, Menlo, 'Ubuntu Mono', monospace;
    font-size: 0.9em;
  }

  .markdown-content pre {
    margin: 8px 0;
    padding: 12px;
    overflow-x: auto;
    border-radius: 6px;
    background: #1f2937;
    color: #f9fafb;
    font-family: Monaco, Menlo, 'Ubuntu Mono', monospace;
    font-size: 0.9em;
    line-height: 1.5;
  }

  .markdown-content pre code {
    padding: 0;
    background: none;
    color: inherit;
  }

  .markdown-content table {
    width: 100%;
    margin: 8px 0;
    border-collapse: collapse;
    border: 1px solid #e5e7eb;
  }

  .markdown-content th,
  .markdown-content td {
    padding: 8px 12px;
    border: 1px solid #e5e7eb;
    text-align: left;
  }

  .markdown-content th {
    background: #f9fafb;
    font-weight: 600;
  }

  .markdown-content a {
    color: #667eea;
    text-decoration: none;
  }

  .markdown-content a:hover {
    text-decoration: underline;
  }

  .markdown-content img {
    max-width: 100%;
    height: auto;
    border-radius: 4px;
  }

  .markdown-content strong {
    font-weight: 600;
  }

  .markdown-content em {
    font-style: italic;
  }

  .markdown-content del {
    color: #6b7280;
    text-decoration: line-through;
  }

  .markdown-content hr {
    margin: 16px 0;
    border: none;
    border-top: 1px solid #e5e7eb;
  }

  .loading-content {
    display: flex;
    flex-direction: column;
    align-items: flex-start;
    gap: 8px;
    padding: 8px 12px;
  }

  .thinking-text {
    color: #6b7280;
    font-size: 14px;
    font-style: italic;
  }

  .thinking-dots {
    display: flex;
    gap: 3px;
  }

  .thinking-dots span {
    width: 6px;
    height: 6px;
    animation: loading 1s infinite;
    border-radius: 50%;
    background: #9ca3af;
  }

  .thinking-dots span:nth-child(2) {
    animation-delay: 0.2s;
  }

  .thinking-dots span:nth-child(3) {
    animation-delay: 0.4s;
  }

  /* 输入区域 */
  .input-container {
    display: flex;
    align-items: flex-end;
    padding: 16px;
    border-top: 1px solid #e5e7eb;
    background: white;
    gap: 8px;
  }

  textarea {
    flex: 1;
    min-height: 36px;
    max-height: 100px;
    padding: 8px 12px;
    transition: border-color 0.1s; /* 缩短输入框过渡时间 */
    border: 1px solid #e5e7eb;
    border-radius: 20px;
    background: #f8fafc;
    font-size: 14px;
    resize: none;
  }

  textarea:focus {
    border-color: #667eea;
    outline: none;
    box-shadow: 0 0 0 2px rgb(102 126 234 / 10%);
  }

  .send-btn {
    display: flex;
    flex-shrink: 0;
    align-items: center;
    justify-content: center;
    width: 36px;
    height: 36px;
    transition: all 0.1s; /* 缩短发送按钮过渡时间 */
    border: none;
    border-radius: 50%;
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    color: white;
    cursor: pointer;
  }

  .send-btn:hover:not(:disabled) {
    transform: scale(1.05);
    box-shadow: 0 2px 8px rgb(102 126 234 / 30%);
  }

  .send-btn:disabled {
    opacity: 0.6;
    cursor: not-allowed;
  }

  .loading-spinner {
    width: 16px;
    height: 16px;
    animation: spin 1s linear infinite;
    border: 2px solid rgb(255 255 255 / 30%);
    border-top: 2px solid white;
    border-radius: 50%;
  }

  /* 机器人头部样式（固定在右下角） */
  .bot-head {
    display: flex;
    position: absolute;
    z-index: 1001;
    right: 0;
    bottom: 0;
    align-items: center;
    justify-content: center;
    width: 60px;
    height: 60px;
    transition: all 0.1s ease; /* 大幅缩短过渡时间，提高响应速度 */
    border-radius: 50%;
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    box-shadow: 0 4px 12px rgb(0 0 0 / 15%);
    cursor: pointer;
  }

  .bot-head:hover {
    transform: scale(1.05);
    box-shadow: 0 6px 16px rgb(0 0 0 / 20%);
  }

  .bot-head.has-notification {
    animation: pulse-glow 1.5s infinite; /* 缩短通知动画时间 */
  }

  .bot-head.dragging {
    transform: scale(1.05);
    transition: none; /* 拖动时禁用所有过渡动画 */
    box-shadow: 0 6px 16px rgb(0 0 0 / 20%);
    cursor: grabbing;
  }

  .bot-head.dragging .eye {
    animation: none; /* 拖动时禁用眨眼动画 */
  }

  .bot-head.dragging .notification-dot {
    animation: none; /* 拖动时禁用通知点动画 */
  }

  .bot-face {
    position: relative;
    width: 40px;
    height: 40px;
  }

  .bot-eyes {
    display: flex;
    justify-content: space-between;
    margin-bottom: 8px;
  }

  .eye {
    width: 8px;
    height: 8px;
    animation: blink 2s infinite; /* 缩短眨眼动画时间 */
    border-radius: 50%;
    background: white;
  }

  .eye.left-eye {
    animation-delay: 0s;
  }

  .eye.right-eye {
    animation-delay: 0.5s;
  }

  .bot-mouth {
    width: 16px;
    height: 4px;
    margin: 0 auto;
    border-radius: 2px;
    background: white;
  }

  .notification-dot {
    position: absolute;
    top: 8px;
    right: 8px;
    width: 12px;
    height: 12px;
    animation: pulse 1.5s infinite; /* 缩短通知点动画时间 */
    border-radius: 50%;
    background: #ef4444;
  }

  /* 头像上的关闭按钮 */
  .bot-close-btn {
    display: flex;
    position: absolute;
    z-index: 1002;
    top: -6px;
    right: -6px;
    align-items: center;
    justify-content: center;
    width: 28px;
    height: 28px;
    transform: scale(0.8);
    transition: all 0.1s ease; /* 缩短关闭按钮过渡时间 */
    border: 3px solid white;
    border-radius: 50%;
    opacity: 0;
    background: #ef4444;
    box-shadow: 0 2px 6px rgb(0 0 0 / 20%);
    color: white;
    cursor: pointer;
  }

  .bot-head:hover .bot-close-btn {
    transform: scale(1);
    opacity: 1;
  }

  .bot-close-btn:hover {
    transform: scale(1.1);
    background: #dc2626;
    box-shadow: 0 4px 12px rgb(239 68 68 / 50%);
  }

  .zoom-hint {
    position: absolute;
    bottom: -25px;
    left: 50%;
    padding: 4px 8px;
    transform: translateX(-50%);
    transition: opacity 0.15s ease; /* 缩短提示框过渡时间 */
    border-radius: 4px;
    opacity: 0;
    background: rgb(0 0 0 / 80%);
    color: white;
    font-size: 10px;
    white-space: nowrap;
  }

  .chat-bot-widget.collapsed .bot-head:hover .zoom-hint {
    opacity: 1;
  }
</style>

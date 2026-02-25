<template>
  <ConfigProvider :locale="getAntdLocale" :theme="themeConfig">
    <AppProvider>
      <RouterView />
      <template v-if="disableChatBot !== 'true'">
        <!-- 浮动聊天机器人 -->
        <ChatBot :options="chatBotOptions" />
      </template>
    </AppProvider>
  </ConfigProvider>
</template>

<script lang="ts" setup>
  import { AppProvider } from '@/components/Application';
  import { useTitle } from '@/hooks/web/useTitle';
  import { useLocale } from '@/locales/useLocale';
  import { ConfigProvider } from 'ant-design-vue';
  import ChatBot from '@/chat-bot/ChatBot.vue';
  import type { ChatBotOptions } from '@/chat-bot/types';
  import { GetDifiApi } from '@/chat-bot/openapi';
  import { useDarkModeTheme } from '@/hooks/setting/useDarkModeTheme';
  import 'dayjs/locale/zh-cn';
  import { computed } from 'vue';
  import { useGlobSetting } from '@/hooks/setting';

  const { disableChatBot } = useGlobSetting();

  // support Multi-language
  const { getAntdLocale } = useLocale();

  const { isDark, darkTheme } = useDarkModeTheme();

  const themeConfig = computed(() =>
    Object.assign(
      {
        token: {
          colorPrimary: '#0960bd',
          colorSuccess: '#55D187',
          colorWarning: '#EFBD47',
          colorError: '#ED6F6F',
          colorInfo: '#0960bd',
        },
      },
      isDark.value ? darkTheme : {},
    ),
  );

  // 聊天机器人配置
  const chatBotOptions: ChatBotOptions = {
    botName: '智能问答-知你所问',
    placeholder: '请输入您的问题...',
    replyDelay: 800,

    // 自定义回复逻辑
    defaultReply: async (input: string) => {
      // 模拟API调用
      // await new Promise(resolve => setTimeout(resolve, 500));

      // 简单的关键词回复逻辑
      const lowerInput = input.toLowerCase();

      if (lowerInput.startsWith('帮助') || lowerInput.startsWith('help')) {
        return `
- **简单对话交流** - 和我聊天
- **HTML/Markdown支持** - 支持富文本内容
有什么问题都可以问我哦！`.trim();
      }
      try {
        const res = await GetDifiApi({
          content: lowerInput,
        });
        console.log(res);
        // 这里可以根据实际返回内容自定义回复
        if (res && res.msg) {
          console.log(res);
          return `${res.msg}`;
        }
      } catch (error: any) {
        console.error(error);
        return '抱歉，我无法回答您的问题，请您稍后再试！';
      }
      return `已收到您的消息："${input}"，后续会陆续开放更多功能！`;
    },

    // 消息回调
    onMessage: (message) => {
      console.log('新消息:', message);
    },
  };

  // Listening to page changes and dynamically changing site titles
  useTitle();
</script>

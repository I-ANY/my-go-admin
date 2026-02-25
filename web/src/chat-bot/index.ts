import { App, Plugin } from 'vue';
import ChatBot from './ChatBot.vue';
import type { ChatBotPluginOptions, BotMessageType } from './types';
import { renderMessageContent, detectContentType } from './utils';

// 扩展 Vue 全局属性类型
declare module 'vue' {
  interface ComponentCustomProperties {
    $chatBot: {
      addMessage: (content: string, isUser: boolean, type?: BotMessageType) => void;
      renderMessageContent: (message: any) => string;
      detectContentType: (content: string) => 'text' | 'html' | 'markdown';
    };
  }
}

// 定义插件
const ChatBotPlugin: Plugin = {
  install(app: App, options: ChatBotPluginOptions = {}) {
    // 注册全局组件
    app.component('ChatBot', ChatBot);

    // 挂载全局属性
    app.config.globalProperties.$chatBot = {
      addMessage: (content: string, isUser: boolean, type?: BotMessageType) => {
        // 这里可以通过 provide/inject 或其他方式与组件实例通信
        console.log('全局添加消息:', { content, isUser, type });
      },
      renderMessageContent: (message: any) => {
        return renderMessageContent(message);
      },
      detectContentType: (content: string) => {
        return detectContentType(content);
      },
    };

    // 提供全局配置
    app.provide('chatBotOptions', options);
  },
};

export default ChatBotPlugin;
export type { ChatBotOptions, ChatBotPluginOptions, BotMessageType } from './types';
export { renderMessageContent, detectContentType } from './utils';

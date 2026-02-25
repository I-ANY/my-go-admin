import 'uno.css';
import '@/design/index.less';
import '@/components/VxeTable/src/css/index.scss';
import 'ant-design-vue/dist/reset.css';
// Register icon sprite
import 'virtual:svg-icons-register';

import { createApp } from 'vue';

import { registerGlobComp } from '@/components/registerGlobComp';
import { setupGlobDirectives } from '@/directives';
import { setupI18n } from '@/locales/setupI18n';
import { setupErrorHandle } from '@/logics/error-handle';
import { initAppConfigStore } from '@/logics/initAppConfig';
import { router, setupRouter } from '@/router';
import { setupRouterGuard } from '@/router/guard';
import { setupStore } from '@/store';
import ChatBotPlugin from '@/chat-bot';
import App from './App.vue';
import Antd from 'ant-design-vue';
import './utils/checkVersion/checkVersion';

async function bootstrap() {
  const app = createApp(App);

  // Configure store
  // 配置 store
  setupStore(app);

  // Initialize internal system configuration
  // 初始化内部系统配置
  initAppConfigStore();

  // Register global components
  // 注册全局组件
  registerGlobComp(app);

  // Multilingual configuration
  // 多语言配置
  // Asynchronous case: language files may be obtained from the server side
  // 异步案例：语言文件可能从服务器端获取
  await setupI18n(app);

  // Configure routing
  // 配置路由
  setupRouter(app);

  // router-guard
  // 路由守卫
  setupRouterGuard(router);

  // Register global directive
  // 注册全局指令
  setupGlobDirectives(app);

  // Configure global error handling
  // 配置全局错误处理
  setupErrorHandle(app);

  // https://next.router.vuejs.org/api/#isready
  // await router.isReady();
  app.use(Antd);
  app.use(ChatBotPlugin, {
    botName: '小助手',
    botAvatar: ' https://example.com/bot-avatar.png ',
    userAvatar: ' https://example.com/user-avatar.png ',
    placeholder: '输入您的问题...',
    replyDelay: 1500,
    defaultReply: async (input: string) => {
      // 这里可以接入真实 AI 接口（如对话 API）
      // 示例：模拟调用 API
      await new Promise((resolve) => setTimeout(resolve, 500));
      return `这是自定义回复逻辑处理的内容："${input}"`;
    },
  });
  app.mount('#app');
}

bootstrap();

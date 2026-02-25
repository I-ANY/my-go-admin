// 机器人消息类型
export type BotMessageType = 'text' | 'markdown' | 'image' | 'loading';

// 消息项接口
export interface MessageItem {
  id: string;
  content: string;
  type: BotMessageType;
  isUser: boolean; // 是否是用户发送的消息
  timestamp: number;
  // 新增：消息渲染选项
  renderOptions?: {
    allowHtml?: boolean; // 是否允许HTML标签
    markdownOptions?: {
      tables?: boolean; // 是否支持表格
      strikethrough?: boolean; // 是否支持删除线
      tasklists?: boolean; // 是否支持任务列表
      ghCodeBlocks?: boolean; // 是否支持GitHub风格代码块
      ghMentions?: boolean; // 是否支持GitHub风格提及
      ghMentionsLink?: string; // GitHub提及链接
      emoji?: boolean; // 是否支持emoji
      underline?: boolean; // 是否支持下划线
      completeHTMLDocument?: boolean; // 是否生成完整HTML文档
      metadata?: boolean; // 是否支持元数据
      splitAdjacentBlockquotes?: boolean; // 是否分割相邻的引用块
      simpleLineBreaks?: boolean; // 是否支持简单换行
      requireSpaceBeforeHeadingText?: boolean; // 标题文本前是否需要空格
      ghCompatibleHeaderId?: boolean; // 是否使用GitHub兼容的标题ID
      prefixHeaderId?: string; // 标题ID前缀
      rawPrefixHeaderId?: boolean; // 是否使用原始前缀
      parseImgDimensions?: boolean; // 是否解析图片尺寸
      simplifiedAutoLink?: boolean; // 是否简化自动链接
      excludeTrailingPunctuationFromURLs?: boolean; // 是否从URL中排除尾随标点
      literalMidWordUnderscores?: boolean; // 是否字面处理单词中的下划线
      tablesHeaderId?: boolean; // 表格标题是否使用ID
      smoothLivePreview?: boolean; // 是否平滑实时预览
      smartIndentationFix?: boolean; // 是否智能缩进修复
      disableForced4SpacesIndentedSublists?: boolean; // 是否禁用强制4空格缩进子列表
    };
  };
}

// 机器人配置选项
export interface ChatBotOptions {
  botName?: string; // 机器人名称
  botAvatar?: string; // 机器人头像 URL
  userAvatar?: string; // 用户头像 URL
  placeholder?: string; // 输入框占位符
  onMessage?: (message: MessageItem) => void; // 消息回调
  replyDelay?: number; // 回复延迟（毫秒）
  defaultReply?: (input: string) => string | Promise<string>; // 默认回复逻辑
  requireLogin?: boolean; // 是否要求用户登录才能使用，默认为true
  showCloseButton?: boolean; // 是否显示关闭按钮，默认为true
  onClose?: () => void; // 关闭回调
  onReopen?: () => void; // 重新打开回调
}

// 插件安装选项
export interface ChatBotPluginOptions extends ChatBotOptions {
  // 可扩展其他全局配置
}

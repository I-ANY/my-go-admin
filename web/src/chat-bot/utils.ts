import showdown from 'showdown';
import type { MessageItem } from './types';

// 创建Markdown转换器实例
const createMarkdownConverter = (
  options?: NonNullable<MessageItem['renderOptions']>['markdownOptions'],
) => {
  const converter = new showdown.Converter();

  if (options) {
    // 设置Markdown选项
    if (options.tables !== undefined) converter.setOption('tables', options.tables);
    if (options.strikethrough !== undefined)
      converter.setOption('strikethrough', options.strikethrough);
    if (options.tasklists !== undefined) converter.setOption('tasklists', options.tasklists);
    if (options.ghCodeBlocks !== undefined)
      converter.setOption('ghCodeBlocks', options.ghCodeBlocks);
    if (options.ghMentions !== undefined) converter.setOption('ghMentions', options.ghMentions);
    if (options.ghMentionsLink !== undefined)
      converter.setOption('ghMentionsLink', options.ghMentionsLink);
    if (options.emoji !== undefined) converter.setOption('emoji', options.emoji);
    if (options.underline !== undefined) converter.setOption('underline', options.underline);
    if (options.completeHTMLDocument !== undefined)
      converter.setOption('completeHTMLDocument', options.completeHTMLDocument);
    if (options.metadata !== undefined) converter.setOption('metadata', options.metadata);
    if (options.splitAdjacentBlockquotes !== undefined)
      converter.setOption('splitAdjacentBlockquotes', options.splitAdjacentBlockquotes);
    if (options.simpleLineBreaks !== undefined)
      converter.setOption('simpleLineBreaks', options.simpleLineBreaks);
    if (options.requireSpaceBeforeHeadingText !== undefined)
      converter.setOption('requireSpaceBeforeHeadingText', options.requireSpaceBeforeHeadingText);
    if (options.ghCompatibleHeaderId !== undefined)
      converter.setOption('ghCompatibleHeaderId', options.ghCompatibleHeaderId);
    if (options.prefixHeaderId !== undefined)
      converter.setOption('prefixHeaderId', options.prefixHeaderId);
    if (options.rawPrefixHeaderId !== undefined)
      converter.setOption('rawPrefixHeaderId', options.rawPrefixHeaderId);
    if (options.parseImgDimensions !== undefined)
      converter.setOption('parseImgDimensions', options.parseImgDimensions);
    if (options.simplifiedAutoLink !== undefined)
      converter.setOption('simplifiedAutoLink', options.simplifiedAutoLink);
    if (options.excludeTrailingPunctuationFromURLs !== undefined)
      converter.setOption(
        'excludeTrailingPunctuationFromURLs',
        options.excludeTrailingPunctuationFromURLs,
      );
    if (options.literalMidWordUnderscores !== undefined)
      converter.setOption('literalMidWordUnderscores', options.literalMidWordUnderscores);
    if (options.tablesHeaderId !== undefined)
      converter.setOption('tablesHeaderId', options.tablesHeaderId);
    if (options.smoothLivePreview !== undefined)
      converter.setOption('smoothLivePreview', options.smoothLivePreview);
    if (options.smartIndentationFix !== undefined)
      converter.setOption('smartIndentationFix', options.smartIndentationFix);
    if (options.disableForced4SpacesIndentedSublists !== undefined)
      converter.setOption(
        'disableForced4SpacesIndentedSublists',
        options.disableForced4SpacesIndentedSublists,
      );
  }

  return converter;
};

// 处理换行符的函数
const processNewlines = (content: string): string => {
  // 去除首尾空白字符
  const trimmedContent = content.trim();

  if (!trimmedContent) {
    return '';
  }

  // 检查是否包含连续换行符（段落分隔）
  if (trimmedContent.includes('\n\n')) {
    // 有段落分隔，使用<p>标签
    let processed = trimmedContent
      .replace(/\n\n+/g, '</p><p>') // 连续换行符转换为段落分隔
      .replace(/\n/g, '<br>'); // 单个换行符转换为<br>

    // 确保内容被段落标签包裹
    if (!processed.startsWith('<p>')) {
      processed = '<p>' + processed;
    }
    if (!processed.endsWith('</p>')) {
      processed = processed + '</p>';
    }

    return processed;
  } else {
    // 没有段落分隔，只处理单个换行符，不使用<p>标签
    return trimmedContent.replace(/\n/g, '<br>');
  }
};

// 渲染消息内容
export const renderMessageContent = (message: MessageItem): string => {
  const { content, type, renderOptions } = message;

  switch (type) {
    case 'text':
      // 处理纯文本中的换行符，优化段落显示
      return processNewlines(content);

    case 'markdown':
      const converter = createMarkdownConverter(renderOptions?.markdownOptions);
      const markdownHtmlContent = converter.makeHtml(content);
      // 对于Markdown内容，直接使用转换后的HTML，不进行额外的换行符处理
      return markdownHtmlContent.trim();

    case 'image':
      return `<img src="${content}" alt="图片" style="max-width: 100%; height: auto;" />`;

    default:
      return content;
  }
};

// HTML转义函数
export const escapeHtml = (text: string): string => {
  const div = document.createElement('div');
  div.textContent = text;
  return div.innerHTML;
};

// 检测内容类型（优先级：Markdown > 换行符）
export const detectContentType = (content: string): 'text' | 'markdown' => {
  // 优先级1：检测Markdown语法（最高优先级）
  const markdownPatterns = [
    /^#{1,6}\s/, // 标题
    /\*\*.*?\*\*/, // 粗体
    /\*.*?\*/, // 斜体
    /`.*?`/, // 行内代码
    /```[\s\S]*?```/, // 代码块
    /\[.*?\]\(.*?\)/, // 链接
    /!\[.*?\]\(.*?\)/, // 图片
    /^\s*[-*+]\s/, // 无序列表
    /^\s*\d+\.\s/, // 有序列表
    /^\s*>\s/, // 引用
    /^\|.*\|$/, // 表格
    /~~.*?~~/, // 删除线
    /^\s*- \[ \]/, // 任务列表
  ];

  for (const pattern of markdownPatterns) {
    if (pattern.test(content)) {
      return 'markdown';
    }
  }

  // 优先级2：如果包含换行符，处理为文本类型
  if (content.includes('\n')) {
    return 'text';
  }

  // 默认返回文本类型
  return 'text';
};

// 安全地设置innerHTML
export const setInnerHTML = (element: HTMLElement, html: string): void => {
  // 简单的XSS防护
  const sanitizedHtml = html.replace(/<script\b[^<]*(?:(?!<\/script>)<[^<]*)*<\/script>/gi, '');
  element.innerHTML = sanitizedHtml;
};

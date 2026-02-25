import { Rule } from 'ant-design-vue/es/form';
import { message } from 'ant-design-vue';

export function splitByLineAndTrim(str: string): string[] | null {
  const ss: string[] = [];
  if (str && str.length > 0) {
    const hs = str.split('\n');
    hs.forEach((item) => {
      if (item.trim()) {
        ss.push(item.trim());
      }
    });
  }
  return ss;
}

export function commonCustomHeaderCell() {
  return () => {
    return {
      style: {
        whiteSpace: 'pre-wrap',
        fontSize: '12px',
        lineHeight: '1',
        // maxHeight: '20px',
      },
    };
  };
}

export const getCheckVariableFn = function (len: number) {
  return async (_rule: Rule, value: string) => {
    // 检查字符串是否只包含字母、数字和下划线
    const regex = /^[a-zA-Z0-9_]+$/;
    // 检查字符串是否为纯数字或纯下划线
    const isPureNumber = /^\d+$/.test(value);
    const isPureUnderscore = /^_+$/.test(value);
    if (value?.length < len) {
      return Promise.reject('长度不能小于' + len);
    }
    // 校验条件
    if (regex.test(value)) {
      if (isPureNumber || isPureUnderscore) {
        return Promise.reject('不能是纯数字或者下划线');
      } else {
        return Promise.resolve();
      }
    } else {
      return Promise.reject('只能包含字母、数字和下划线');
    }
  };
};

export function customCeilDivide(num, divisor = 1000 * 1000, decimalPlaces = 3) {
  if (typeof num !== 'number' || isNaN(num)) {
    return;
  }
  if (divisor === 0) {
    return;
  }
  const factor = Math.pow(10, decimalPlaces);
  const scaled = (num * factor) / divisor;
  // 使用 EPSILON 处理浮点数精度
  const corrected = Math.round(scaled / Number.EPSILON) * Number.EPSILON;
  const ceiled = Math.ceil(corrected);
  const result = (ceiled / factor).toFixed(decimalPlaces);
  return result;
}

export function customFloorDivide(num, divisor = 1000 * 1000, decimalPlaces = 3) {
  if (typeof num !== 'number' || isNaN(num)) {
    return;
  }
  if (divisor === 0) {
    return;
  }
  const factor = Math.pow(10, decimalPlaces);
  const scaled = (num * factor) / divisor;
  // 使用 EPSILON 处理浮点数精度
  const corrected = Math.round(scaled / Number.EPSILON) * Number.EPSILON;
  const ceiled = Math.floor(corrected);
  const result = (ceiled / factor).toFixed(decimalPlaces);
  return result;
}

// 计算单元格合并的行数
export function calTableCellRowSpan(
  dataSource: any[],
  currentIndex: number,
  needMergeFn: (current: Recordable, target: Recordable) => boolean,
): number {
  if (
    !dataSource ||
    dataSource.length === 0 ||
    currentIndex < 0 ||
    currentIndex >= dataSource.length ||
    !needMergeFn
  ) {
    return 1;
  }
  const startIndex = currentIndex;
  let endIndex = currentIndex;

  // 向上查找一行，如果向上能够找得到说明这一行不是起始行，直接返回跨越0行
  if (startIndex > 0 && needMergeFn(dataSource[currentIndex], dataSource[startIndex - 1])) {
    return 0;
  }

  //  向下查找可以连续合并的行，直到找到不能合并的行退出
  while (
    endIndex < dataSource.length - 1 &&
    needMergeFn(dataSource[currentIndex], dataSource[endIndex + 1])
  ) {
    endIndex++;
  }
  return endIndex - startIndex + 1;
}

export function cutStringLength(result: string, length: number = 50): string {
  if (!result) return '';
  return result.length > length ? result.substring(0, length) + '...' : result;
}

// 复制文本内容
export async function handleCopy(text: string) {
  if (!text) {
    message.warning('内容为空');
    return;
  }
  try {
    if (navigator.clipboard) {
      navigator.clipboard.writeText(text);
    } else {
      // 兼容旧浏览器
      const textarea = document.createElement('textarea');
      textarea.value = text;
      document.body.appendChild(textarea);
      textarea.select();
      document.execCommand('copy');
      document.body.removeChild(textarea);
    }
    message.success('复制成功');
  } catch (error) {
    message.error('复制失败');
  }
}

// 复制表格行数据
export async function handleCopyTableRowData<T>(
  dataSource: T[],
  fieldExtractor: (item: T) => string,
  fieldName: string,
  successMessage?: string,
) {
  if (!dataSource || dataSource.length === 0) {
    message.warning('暂无数据可复制');
    return;
  }

  // 提取指定字段，用换行符连接
  const extractedData = dataSource
    .map(fieldExtractor)
    .filter((value) => value) // 过滤掉空值
    .join('\n');

  if (!extractedData) {
    message.warning(`未找到${fieldName}数据`);
    return;
  }

  try {
    if (navigator.clipboard) {
      navigator.clipboard.writeText(extractedData);
    } else {
      // 兼容旧浏览器
      const textarea = document.createElement('textarea');
      textarea.value = extractedData;
      document.body.appendChild(textarea);
      textarea.select();
      document.execCommand('copy');
      document.body.removeChild(textarea);
    }
    if (successMessage) {
      message.success(successMessage);
    } else {
      message.success(`已复制${dataSource.length} 个${fieldName}`);
    }
  } catch (error) {
    message.error('复制失败');
  }
}

/**
 * 高亮特定关键字
 */
export function highlightKeywords(text: string): string {
  if (!text) return '';

  // 1. 转义 HTML (防止 XSS)
  let html = text.replace(
    /[&<>"']/g,
    (m) => ({ '&': '&amp;', '<': '&lt;', '>': '&gt;', '"': '&quot;', "'": '&#39;' })[m] || m,
  );

  // 1.5 处理 ANSI 颜色代码 (简单支持)
  // 颜色映射表
  const ansiColors: Record<string, string> = {
    '30': 'black',
    '31': '#ff4d4f', // red
    '32': '#52c41a', // green
    '33': '#faad14', // yellow
    '34': '#1890ff', // blue
    '35': '#722ed1', // magenta
    '36': '#13c2c2', // cyan
    '37': '#ffffff', // white
    '90': '#8c8c8c', // gray
  };

  // 替换 ANSI 颜色码 \x1b[33m -> <span style="color:...">
  // 支持格式: \x1b[n m 或 \x1b[n;n m (简单处理)
  // eslint-disable-next-line no-control-regex
  html = html.replace(/\x1b\[(\d+)(?:;(\d+))?m/g, (_match, p1) => {
    const code = p1;
    if (code === '0') {
      return '</span>';
    }
    const color = ansiColors[code];
    if (color) {
      return `<span style="color: ${color}">`;
    }
    return ''; // 移除不支持的颜色码
  });

  // 移除其他常见 ANSI 控制序列 (防止乱码)
  // eslint-disable-next-line no-control-regex
  html = html.replace(/\x1b\[[\d;?]*[a-zA-Z]/g, '');

  // 2. 定义关键字和颜色
  const map: { [key: string]: string } = {
    // Red - 错误/失败/危险
    error: '#ff4d4f',
    failed: '#ff4d4f',
    fail: '#ff4d4f',
    failure: '#ff4d4f',
    false: '#ff4d4f',
    fatal: '#ff4d4f',
    critical: '#ff4d4f',
    exception: '#ff4d4f',
    crash: '#ff4d4f',
    denied: '#ff4d4f',
    refused: '#ff4d4f',
    unreachable: '#ff4d4f',
    timeout: '#ff4d4f',
    down: '#ff4d4f',
    offline: '#ff4d4f',
    stopped: '#ff4d4f',
    err: '#ff4d4f',
    closed: '#ff4d4f',
    异常: '#ff4d4f',
    失败: '#ff4d4f',

    // Green - 成功/正常/运行中
    正常: '#52c41a',
    成功: '#52c41a',
    ok: '#52c41a',

    success: '#52c41a',
    true: '#52c41a',
    completed: '#52c41a',
    passed: '#52c41a',
    connected: '#52c41a',
    running: '#52c41a',
    up: '#52c41a',
    online: '#52c41a',
    active: '#52c41a',
    enabled: '#52c41a',
    done: '#52c41a',
    open: '#52c41a',
    started: '#52c41a',

    // Orange - 警告/等待/未知
    warning: '#faad14',
    warn: '#faad14',
    pending: '#faad14',
    waiting: '#faad14',
    disabled: '#faad14',
    unknown: '#faad14',
    paused: '#faad14',
    blocked: '#faad14',

    // Blue - 信息/跳过
    info: '#1890ff',
    notice: '#1890ff',
    skipped: '#1890ff',
    debug: '#1890ff',
  };

  // 3. 构建正则
  // 按长度降序排序，确保长词优先匹配
  const keys = Object.keys(map).sort((a, b) => b.length - a.length);

  // 区分处理：英文单词加边界，中文/符号直接匹配
  const patterns = keys.map((key) => {
    // 检查是否全由单词字符组成 (a-z, 0-9, _)
    if (/^\w+$/.test(key)) {
      return `\\b${key}\\b`;
    } else {
      // 中文或其他字符，需要转义正则特殊字符
      return key.replace(/[.*+?^${}()|[\]\\]/g, '\\$&');
    }
  });

  const pattern = new RegExp(`(${patterns.join('|')})`, 'gi');

  // 4. 替换
  return html.replace(pattern, (match) => {
    const lower = match.toLowerCase();
    const color = map[lower];
    return `<span style="color: ${color}; font-weight: bold;">${match}</span>`;
  });
}

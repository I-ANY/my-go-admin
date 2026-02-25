export const roomTypeMap = {
  1: 'IDC',
  2: 'ACDN',
  3: 'MCDN',
};
export const stateMap = {
  0: '未知',
  1: '达标',
  2: '不达标',
};
export const reportTypeMap = {
  2: '保底',
  3: '削峰',
};
export const billingTypeOptions = [
  { label: '95', value: '95' },
  { label: '日95', value: '日95' },
  { label: '月95', value: '月95' },
  { label: '单机95', value: '单机95' },
  { label: '买断', value: '买断' },
];

export const isMainMap = {
  true: '是',
  false: '否',
};

export const trTypeOptions = [
  { label: '正常', value: '正常' },
  { label: '跨省超出', value: '跨省超出' },
  { label: '跨省浪费', value: '跨省浪费' },
  { label: '未开启', value: '未开启' },
];
export const profitMap = {
  0: { label: '未知', color: '#999999' }, // 灰色
  1: { label: '盈利S', color: '#52c41a' }, // 绿色
  2: { label: '盈利A', color: '#73d13d' }, // 浅绿色
  3: { label: '盈利B', color: '#95de64' }, // 更浅绿色
  4: { label: '盈利C', color: '#b7eb8f' }, // 最浅绿色
  5: { label: '收支平衡', color: '#1890ff' }, // 蓝色
  6: { label: '亏损', color: '#f5222d' }, // 红色
};

export const isps = ['电信', '移动', '联通', '教育网'];
// 获取同一节点编号的行范围
export const getNodeRowRange = (dataSource: any[], index: number) => {
  if (!dataSource || dataSource.length === 0) return { start: index, end: index };

  const currentRoomNo = dataSource[index].roomNo;
  let start = index;
  let end = index;

  // 向上查找相同节点编号的起始位置
  while (start > 0 && dataSource[start - 1].roomNo === currentRoomNo) {
    start--;
  }

  // 向下查找相同节点编号的结束位置
  while (end < dataSource.length - 1 && dataSource[end + 1].roomNo === currentRoomNo) {
    end++;
  }

  return { start, end };
};

// 计算基于节点编号的行合并函数
export const calculateRowSpanByNode = (
  dataSource: any[],
  index: number,
  key: string,
  compareFn?: (current: any, next: any) => boolean,
) => {
  try {
    if (!dataSource || dataSource.length === 0 || index < 0 || index >= dataSource.length) {
      return 1;
    }

    const { start, end } = getNodeRowRange(dataSource, index);
    const currentRecord = dataSource[index];

    // 在同一节点编号范围内，查找连续相同值的范围
    let spanStart = index;
    let spanEnd = index;

    // 向上查找相同值（但不超出节点范围）
    if (compareFn) {
      // 使用自定义比较函数
      while (spanStart > start && compareFn(dataSource[spanStart - 1], currentRecord)) {
        spanStart--;
      }

      // 向下查找相同值（但不超出节点范围）
      while (spanEnd < end && compareFn(dataSource[spanEnd + 1], currentRecord)) {
        spanEnd++;
      }
    } else {
      // 使用默认比较逻辑
      const currentValue = currentRecord[key];

      // 向上查找相同值（但不超出节点范围）
      while (spanStart > start && dataSource[spanStart - 1][key] === currentValue) {
        spanStart--;
      }

      // 向下查找相同值（但不超出节点范围）
      while (spanEnd < end && dataSource[spanEnd + 1][key] === currentValue) {
        spanEnd++;
      }
    }

    const rowSpan = spanEnd - spanStart + 1;

    // 如果当前行是第一个相同值的行，返回总的行数；否则返回0（不显示）
    if (spanStart === index) {
      return rowSpan;
    }
    return 0;
  } catch (error) {
    console.warn('calculateRowSpanByNode error:', error);
    return 1;
  }
};

export const filterOptionBySpace = (inputValue, option: any) => {
  if (!inputValue) {
    return true;
  }
  let flag = true;
  inputValue.split(' ')?.forEach((char: string) => {
    if (char === '') {
      return;
    }
    // 其中任何一个没有匹配到，则返回false
    if (option.label.toLowerCase().indexOf(char.toLowerCase()) < 0) {
      flag = false;
      return;
    }
  });
  return flag;
};

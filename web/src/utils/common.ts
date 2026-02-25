import dayjs from 'dayjs';

export const RangePickPresetsExact = () => [
  {
    label: '今天',
    value: [
      dayjs(dayjs().format('YYYY-MM-DD 00:00:00'), 'YYYY-MM-DD HH:mm:ss'),
      dayjs(dayjs().format('YYYY-MM-DD 23:59:59'), 'YYYY-MM-DD HH:mm:ss'),
    ],
  },
  {
    label: '昨天',
    value: [
      dayjs(dayjs().subtract(1, 'day').format('YYYY-MM-DD 00:00:00'), 'YYYY-MM-DD HH:mm:ss'),
      dayjs(dayjs().subtract(1, 'day').format('YYYY-MM-DD 23:59:59'), 'YYYY-MM-DD HH:mm:ss'),
    ],
  },
  {
    label: '前天',
    value: [
      dayjs(dayjs().subtract(2, 'day').format('YYYY-MM-DD 00:00:00'), 'YYYY-MM-DD HH:mm:ss'),
      dayjs(dayjs().subtract(2, 'day').format('YYYY-MM-DD 23:59:59'), 'YYYY-MM-DD HH:mm:ss'),
    ],
  },
  {
    label: '本周',
    value: [
      dayjs(
        dayjs().startOf('week').add(1, 'day').format('YYYY-MM-DD 00:00:00'),
        'YYYY-MM-DD HH:mm:ss',
      ),
      dayjs(dayjs().endOf('week').format('YYYY-MM-DD 23:59:59'), 'YYYY-MM-DD HH:mm:ss'),
    ],
  },
  {
    label: '上周',
    value: [
      dayjs(
        dayjs().add(-1, 'week').startOf('week').add(1, 'day').format('YYYY-MM-DD 00:00:00'),
        'YYYY-MM-DD HH:mm:ss',
      ),
      dayjs(
        dayjs().add(-1, 'week').endOf('week').add(1, 'day').format('YYYY-MM-DD 23:59:59'),
        'YYYY-MM-DD HH:mm:ss',
      ),
    ],
  },
  {
    label: '本月',
    value: [
      dayjs(dayjs().startOf('month').format('YYYY-MM-DD 00:00:00'), 'YYYY-MM-DD HH:mm:ss'),
      dayjs(dayjs().endOf('month').format('YYYY-MM-DD 23:59:59'), 'YYYY-MM-DD HH:mm:ss'),
    ],
  },
  {
    label: '上月',
    value: [
      dayjs(
        dayjs().add(-1, 'month').startOf('month').format('YYYY-MM-DD 00:00:00'),
        'YYYY-MM-DD HH:mm:ss',
      ),
      dayjs(
        dayjs().add(-1, 'month').endOf('month').add(0, 'day').format('YYYY-MM-DD 23:59:59'),
        'YYYY-MM-DD HH:mm:ss',
      ),
    ],
  },
  { label: '最近3小时', value: [dayjs().add(-3, 'h'), dayjs()] },
  { label: '最近6小时', value: [dayjs().add(-6, 'h'), dayjs()] },
  { label: '最近8小时', value: [dayjs().add(-8, 'h'), dayjs()] },
  { label: '最近12小时', value: [dayjs().add(-12, 'h'), dayjs()] },
  { label: '最近1天', value: [dayjs().add(-1, 'd'), dayjs()] },
  { label: '最近3天', value: [dayjs().add(-3, 'd'), dayjs()] },
  { label: '最近7天', value: [dayjs().add(-7, 'd'), dayjs()] },
  { label: '最近30天', value: [dayjs().add(-30, 'd'), dayjs()] },
];

export const RangeDataPickPresetsExact = () => [
  {
    label: '本周',
    value: [
      dayjs(
        dayjs().startOf('week').add(1, 'day').format('YYYY-MM-DD 00:00:00'),
        'YYYY-MM-DD HH:mm:ss',
      ),
      dayjs(dayjs().endOf('week').format('YYYY-MM-DD 23:59:59'), 'YYYY-MM-DD HH:mm:ss'),
    ],
  },
  {
    label: '上周',
    value: [
      dayjs(
        dayjs().add(-1, 'week').startOf('week').add(1, 'day').format('YYYY-MM-DD 00:00:00'),
        'YYYY-MM-DD HH:mm:ss',
      ),
      dayjs(
        dayjs().add(-1, 'week').endOf('week').add(1, 'day').format('YYYY-MM-DD 23:59:59'),
        'YYYY-MM-DD HH:mm:ss',
      ),
    ],
  },
  {
    label: '本月',
    value: [
      dayjs(dayjs().startOf('month').format('YYYY-MM-DD 00:00:00'), 'YYYY-MM-DD HH:mm:ss'),
      dayjs(dayjs().endOf('month').format('YYYY-MM-DD 23:59:59'), 'YYYY-MM-DD HH:mm:ss'),
    ],
  },
  {
    label: '上月',
    value: [
      dayjs(
        dayjs().add(-1, 'month').startOf('month').format('YYYY-MM-DD 00:00:00'),
        'YYYY-MM-DD HH:mm:ss',
      ),
      dayjs(
        dayjs().add(-1, 'month').endOf('month').format('YYYY-MM-DD 23:59:59'),
        'YYYY-MM-DD HH:mm:ss',
      ),
    ],
  },
  {
    label: '最近3天',
    value: [
      dayjs(dayjs().add(-3, 'd').format('YYYY-MM-DD 00:00:00'), 'YYYY-MM-DD HH:mm:ss'),
      dayjs(dayjs().format('YYYY-MM-DD 23:59:59'), 'YYYY-MM-DD HH:mm:ss'),
    ],
  },
  {
    label: '最近7天',
    value: [
      dayjs(dayjs().add(-7, 'd').format('YYYY-MM-DD 00:00:00'), 'YYYY-MM-DD HH:mm:ss'),
      dayjs(dayjs().format('YYYY-MM-DD 23:59:59'), 'YYYY-MM-DD HH:mm:ss'),
    ],
  },
  {
    label: '最近15天',
    value: [
      dayjs(dayjs().add(-15, 'd').format('YYYY-MM-DD 00:00:00'), 'YYYY-MM-DD HH:mm:ss'),
      dayjs(dayjs().format('YYYY-MM-DD 23:59:59'), 'YYYY-MM-DD HH:mm:ss'),
    ],
  },
  {
    label: '最近30天',
    value: [
      dayjs(dayjs().add(-30, 'd').format('YYYY-MM-DD 00:00:00'), 'YYYY-MM-DD HH:mm:ss'),
      dayjs(dayjs().format('YYYY-MM-DD 23:59:59'), 'YYYY-MM-DD HH:mm:ss'),
    ],
  },
  {
    label: '最近60天',
    value: [
      dayjs(dayjs().add(-60, 'd').format('YYYY-MM-DD 00:00:00'), 'YYYY-MM-DD HH:mm:ss'),
      dayjs(dayjs().format('YYYY-MM-DD 23:59:59'), 'YYYY-MM-DD HH:mm:ss'),
    ],
  },
];

export const RangeDatafor1dayPickPresetsExact = () => [
  {
    label: '最近1天',
    value: [
      dayjs(dayjs().add(-1, 'd').format('YYYY-MM-DD 00:00:00'), 'YYYY-MM-DD HH:mm:ss'),
      dayjs(dayjs().format('YYYY-MM-DD 23:59:59'), 'YYYY-MM-DD HH:mm:ss'),
    ],
  },
];

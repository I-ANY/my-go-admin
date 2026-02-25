import { getDictDataMapFromDict, getSelectOptionsFromDict } from '@/utils/dict';

import { BasicColumn } from '@/components/Table';
import type { EChartsOption } from 'echarts';
import { FormSchema } from '@/components/Form';
import { KEnum } from '@/enums/dictTypeCode';
import dayjs from 'dayjs';
import echarts from '@/utils/lib/echarts';
import { h } from 'vue';

export const ispOptions = getSelectOptionsFromDict(KEnum.ISP)?.filter((item) => {
  return item.value != '-1';
});
export const ispMap = getDictDataMapFromDict(KEnum.ISP);

export const devTypeOptions = getSelectOptionsFromDict(KEnum.DEVICE_TYPE);
export const devTypeMap = getDictDataMapFromDict(KEnum.DEVICE_TYPE);

export const providerTypeMap = getDictDataMapFromDict(KEnum.PROVIDER_TYPE);
export const providerTypeOptions = getSelectOptionsFromDict(KEnum.PROVIDER_TYPE);

// 获取前一天
const yesterday = dayjs().subtract(1, 'day').format('YYYY-MM-DD');

// 搜索表单schema
export const searchFormSchema: FormSchema[] = [
  {
    field: 'date',
    label: '日期',
    component: 'DatePicker',
    componentProps: {
      format: 'YYYY-MM-DD',
      valueFormat: 'YYYY-MM-DD',
      defaultValue: yesterday,
    },
    colProps: { span: 4 },
  },
  {
    field: 'provider_id',
    label: '厂商',
    component: 'Select',
    colProps: { span: 4 },
    componentProps: {
      options: providerTypeOptions,
      // 默认值
      defaultValue: providerTypeMap[55].value,
    },
  },
  {
    field: 'isp',
    label: '运营商',
    component: 'Select',
    colProps: { span: 4 },
    componentProps: {
      options: ispOptions,
      mode: 'multiple',
    },
  },
  {
    field: 'dev_type',
    label: '设备类型',
    component: 'Select',
    componentProps: {
      options: devTypeOptions,
      allowClear: true,
      mode: 'multiple',
    },
    colProps: { span: 4 },
  },
  {
    field: 'area',
    label: '区域',
    component: 'ApiSelect',
    colProps: {
      span: 4,
    },
    componentProps: {
      placeholder: '请选择',
      // mode: 'multiple',
    },
  },
  {
    field: 'province',
    label: '省份',
    component: 'ApiSelect',
    colProps: {
      span: 4,
    },
    componentProps: {
      placeholder: '请选择',
      // mode: 'multiple',
    },
  },
];

// 图表配置生成函数
export const getChartOptions = (data: any[]): EChartsOption => {
  // 空数据
  const safeData = data || [];

  return {
    legend: {
      data: ['腾讯后台流量', 'autoops流量', 'ECDN流量'],
      left: 'center',
    },
    tooltip: {
      trigger: 'axis',
      axisPointer: {
        lineStyle: {
          width: 1,
          color: '#019680',
        },
      },
      formatter: (params: any) => {
        const time = dayjs(params[0].axisValue).format('YYYY-MM-DD HH:mm:ss');
        let result = `${time}<br/>`;
        params.forEach((param: any) => {
          // 在名称前添加颜色块
          result += `<span style="display:inline-block;margin-right:5px;border-radius:50%;width:10px;height:10px;background-color:${param.color};"></span>${param.seriesName}: ${param.value[1]} Mbps<br/>`;
        });
        return result;
      },
    },
    grid: {
      left: '1%',
      right: '1%',
      top: 80,
      bottom: 0,
      containLabel: true,
    },
    xAxis: {
      type: 'time',
      name: '时间',
      boundaryGap: [0, 0],
      splitLine: {
        show: true,
        lineStyle: {
          width: 1,
          type: 'solid',
          color: 'rgba(226,226,226,0.5)',
        },
      },
      axisTick: {
        show: false,
      },
    },
    yAxis: {
      type: 'value',
      name: '带宽(Mbps)',
      max: (value) => Math.max(value.max, 100),
      splitNumber: 4,
      splitArea: {
        show: true,
        areaStyle: {
          color: ['rgba(255,255,255,0.2)', 'rgba(226,226,226,0.2)'],
        },
      },
    },
    series: [
      {
        name: '腾讯后台流量',
        type: 'line',
        smooth: true,
        symbol: 'none',
        itemStyle: {
          color: '#0057ff', // 蓝色
        },
        areaStyle: {
          color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
            { offset: 0, color: 'rgba(0, 123, 255, 0.18)' },
            { offset: 1, color: 'rgba(0, 123, 255, 0)' },
          ]),
        },
        data: safeData.map((item) => [
          new Date(item.time_seq),
          (((item.tencent_flow / 300) * 8) / 1000 / 1000).toFixed(3),
        ]),
      },
      {
        name: 'autoops流量',
        type: 'line',
        smooth: true,
        symbol: 'none',
        itemStyle: {
          color: '#ff7d00', // 橙色
        },
        data: safeData.map((item) => [
          new Date(item.time_seq),
          (((item.autoops_flow / 300) * 8) / 1000 / 1000).toFixed(3), // 单位bytes
        ]),
      },
      {
        name: 'ECDN流量',
        type: 'line',
        smooth: true,
        symbol: 'none',
        itemStyle: {
          color: '#52c41a', // 绿色
        },
        data: safeData.map((item) => [
          new Date(item.time_seq),
          (item.ecdn_flow / 1000 / 1000).toFixed(3), // ecdn流量单位是bps，1000进位
        ]),
      },
    ],
  };
};

export const dailyPeekSearchFormSchema: FormSchema[] = [
  {
    field: 'date',
    label: '日期',
    component: 'DatePicker',
    componentProps: {
      format: 'YYYY-MM-DD',
      valueFormat: 'YYYY-MM-DD',
      defaultValue: yesterday,
    },
    colProps: { span: 4 },
  },
  {
    field: 'provider_id',
    label: '厂商',
    component: 'Select',
    componentProps: {
      options: providerTypeOptions,
      allowClear: true,
      // defaultValue: providerTypeMap[55].dictLabel,
    },
    colProps: { span: 4 },
  },
  {
    field: 'isp',
    label: '运营商',
    component: 'Select',
    componentProps: {
      options: ispOptions,
      allowClear: true,
      mode: 'multiple',
    },
    colProps: { span: 4 },
  },
  {
    field: 'dev_type',
    label: '设备类型',
    component: 'Select',
    componentProps: {
      options: devTypeOptions,
      allowClear: true,
      mode: 'multiple',
    },
    colProps: { span: 4 },
  },
  {
    field: 'area',
    label: '区域',
    component: 'ApiSelect',
    componentProps: {
      options: [],
      allowClear: true,
      // mode: 'multiple',
    },
    colProps: { span: 4 },
  },
  {
    field: 'province',
    label: '省份',
    component: 'ApiSelect',
    componentProps: {
      options: [],
      allowClear: true,
      // mode: 'multiple',
    },
    colProps: { span: 4 },
  },
  {
    field: 'ecdn_diff_rate_range',
    slot: 'ecdn_diff_rate_range',
    label: 'ecdn差值率',
    colProps: { span: 8 },
  },
  {
    field: 'autoops_diff_rate_range',
    slot: 'autoops_diff_rate_range',
    label: 'autoops差值率',
    colProps: { span: 8 },
  },
];

export const dailyPeekColumns: BasicColumn[] = [
  {
    title: '日期',
    dataIndex: 'date',
    width: 100,
    resizable: true,
  },
  {
    title: '厂商',
    dataIndex: 'provider_id',
    width: 100,
    resizable: true,
  },
  {
    title: '运营商',
    dataIndex: 'isp',
    width: 100,
    resizable: true,
  },
  {
    title: '设备类型',
    dataIndex: 'dev_type',
    width: 100,
    resizable: true,
  },
  {
    title: '区域',
    dataIndex: 'area',
    width: 100,
    resizable: true,
  },
  {
    title: '省份',
    dataIndex: 'province',
    width: 100,
    resizable: true,
  },
  {
    title: 'Ecdn峰值(Mb)',
    dataIndex: 'ecdn_peak',
    width: 120,
    resizable: true,
    sorter: true,
    customRender: ({ record }) => {
      if (record.ecdn_peak) {
        return (((record.ecdn_peak / 300) * 8) / 1000 / 1000).toFixed(3);
      }
      return '-';
    },
  },
  {
    title: 'Autoops峰值(Mb)',
    dataIndex: 'autoops_peak',
    width: 120,
    resizable: true,
    sorter: true,
    customRender: ({ record }) => {
      if (record.autoops_peak) {
        return (((record.autoops_peak / 300) * 8) / 1000 / 1000).toFixed(3);
      }
      return '-';
    },
  },
  {
    title: '腾讯峰值(Mb)',
    dataIndex: 'tencent_peak',
    width: 120,
    resizable: true,
    sorter: true,
    customRender: ({ record }) => {
      if (record.tencent_peak) {
        return (((record.tencent_peak / 300) * 8) / 1000 / 1000).toFixed(3);
      }
      return '-';
    },
  },
  {
    title: 'Ecdn峰值时间',
    dataIndex: 'ecdn_peak_time',
    width: 100,
    resizable: true,
    customRender: ({ record }) => {
      if (record.ecdn_peak_time) {
        return dayjs(record.ecdn_peak_time).format('HH:mm:ss');
      }
      return '-';
    },
  },
  {
    title: 'Autoops峰值时间',
    dataIndex: 'autoops_peak_time',
    width: 100,
    resizable: true,
    customRender: ({ record }) => {
      if (record.autoops_peak_time) {
        return dayjs(record.autoops_peak_time).format('HH:mm:ss');
      }
      return '-';
    },
  },
  {
    title: '腾讯峰值时间',
    dataIndex: 'tencent_peak_time',
    width: 100,
    resizable: true,
    format: (value) => {
      if (value) {
        return dayjs(value).format('HH:mm:ss');
      }
      return '-';
    },
  },
  {
    title: 'ecdn差值率',
    dataIndex: 'ecdn_diff_rate',
    width: 120,
    resizable: true,
    sorter: true,
    customRender: ({ record }) => {
      if (record.ecdn_diff_rate !== -1001) {
        const percentage = record.ecdn_diff_rate;
        const percentageText = `${percentage.toFixed(2)}%`;

        // 差值率大于 ±1% 时显示红色
        const color = Math.abs(percentage) > 1 ? '#ff4d4f' : '#262626';

        return h('span', { style: { color } }, percentageText);
      }
      return '-';
    },
  },
  {
    title: 'autoops差值率',
    dataIndex: 'autoops_diff_rate',
    width: 120,
    resizable: true,
    sorter: true,
    customRender: ({ record }) => {
      if (record.autoops_diff_rate !== -1001) {
        const percentage = record.autoops_diff_rate;
        const percentageText = `${percentage.toFixed(2)}%`;

        // 差值率大于 ±1% 时显示红色
        const color = Math.abs(percentage) > 1 ? '#ff4d4f' : '#262626';

        return h('span', { style: { color } }, percentageText);
      }
      return '-';
    },
  },
];

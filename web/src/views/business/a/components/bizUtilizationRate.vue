<template>
  <BasicTable
    @register="registerTable"
    @change="handleTableChange"
    @edit-cancel="
      (data) => {
        // 手动回滚到原始值
        if (data && data.record && data.key) {
          data.record[data.key] = data.value;
        }
      }
    "
  >
    <template #toolbar>
      <a-button
        type="primary"
        v-auth="APermissionCodeEnum.BUSINESS_A_UTILIZATION_RATE_BIZ_EXPORT"
        @click="handleExportData"
        :loading="data.exporting"
        >{{ data.exportButTitle }}</a-button
      >
    </template>
    <template #bodyCell="{ column, record }">
      <template v-if="column.key == 'roomType'">
        <Tag style="margin: 1px">{{ roomTypeMap[record.roomType] }}</Tag>
      </template>
      <template v-if="column.key == 'state'">
        <Tag :color="record.state === 1 ? 'green' : 'red'" style="margin: 1px">
          {{ stateMap[record.state] }}
        </Tag>
      </template>
      <template v-if="column.key == 'reportType'">
        <Tag style="margin: 1px">{{ reportTypeMap[record.reportType] }}</Tag>
      </template>
      <template v-if="column.key == 'bwUsageRateDay'">
        {{ formatPercentage(record.bwUsageRateDay) }}
      </template>
      <template v-if="column.key == 'bwUsageRateNight'">
        {{ formatPercentage(record.bwUsageRateNight) }}
      </template>
      <template v-if="column.key == 'bizs'">
        {{ record.bizs ? record.bizs.join(', ') : '-' }}
      </template>
      <!-- 前一天 -->
      <template v-if="column.key == 'lastOneDayBwUsageRateNight'">
        <Tooltip v-if="record.bwUsageRateNightExtend?.lastOneDay?.opsDescribe">
          <Badge>
            <span
              :style="{
                display: 'block',
                paddingRight: '8px',
              }"
              >{{ formatPercentage(record.bwUsageRateNightExtend?.lastOneDay?.bwUsageRateNight) }}
            </span>
            <template #count>
              <!-- 橙色 -->
              <ExclamationCircleFilled style="color: #ff8c00" />
            </template>
          </Badge>
          <template #title>
            <span>{{ record.bwUsageRateNightExtend?.lastOneDay?.opsDescribe }}</span>
          </template>
        </Tooltip>
        <span v-else>{{
          formatPercentage(record.bwUsageRateNightExtend?.lastOneDay?.bwUsageRateNight)
        }}</span>
      </template>
      <!-- 前两天 -->
      <template v-if="column.key == 'lastTwoDaysBwUsageRateNight'">
        <Tooltip v-if="record.bwUsageRateNightExtend?.lastTwoDay?.opsDescribe">
          <Badge>
            <span
              :style="{
                display: 'block',
                paddingRight: '8px',
              }"
              >{{ formatPercentage(record.bwUsageRateNightExtend?.lastTwoDay?.bwUsageRateNight) }}
            </span>
            <template #count>
              <!-- 橙色 -->
              <ExclamationCircleFilled style="color: #ff8c00" />
            </template>
          </Badge>
          <template #title>
            <span>{{ record.bwUsageRateNightExtend?.lastTwoDay?.opsDescribe }}</span>
          </template>
        </Tooltip>
        <span v-else>{{
          formatPercentage(record.bwUsageRateNightExtend?.lastTwoDay?.bwUsageRateNight)
        }}</span>
      </template>
    </template>
    <template #tableTop>
      <div style="margin-bottom: 8px; padding: 8px; background-color: white">
        <div
          >1、<span style="color: #4ebef1">节点当前利用率</span
          >：当前利用率为最近一次打点记录的利用率数据，如果节点类型为保底，则当前利用率为保底业务整体跑量/机房整体规划带宽；否则为该机房整体跑量/机房规划带宽；</div
        >
        <div
          >2、<span style="color: #4ebef1">跨省占比</span
          >：跨省占比=(跨省跑量带宽-业务规划带宽)/业务规划带宽，超出实际规划带宽为业务超出，低于实际规划带宽为业务浪费；</div
        >
      </div>
    </template>
  </BasicTable>
</template>
<script setup lang="ts">
  import { Modal, notification, Tag, Tooltip, Badge } from 'ant-design-vue';
  import { BasicTable, useTable, BasicColumn, FormSchema } from '@/components/Table';
  import { h, nextTick, onMounted, reactive } from 'vue';
  import dayjs from 'dayjs';
  import { getOwners, getBizUtilRateList, EditBizUtilRate, getBizs, Api } from '@/api/business/a';
  import type { Key } from 'ant-design-vue/lib/table/interface';
  import { calTableCellRowSpan, commonCustomHeaderCell, customCeilDivide } from '@/utils/util';
  import { APermissionCodeEnum } from '@/enums/permissionCodeEnum';
  import { ExclamationCircleOutlined, ExclamationCircleFilled } from '@ant-design/icons-vue';
  import { downloadFileByUrl } from '@/utils/download';
  import { usePermission } from '@/hooks/web/usePermission';
  import {
    billingTypeOptions,
    roomTypeMap,
    stateMap,
    reportTypeMap,
    isMainMap,
    isps,
    calculateRowSpanByNode,
    trTypeOptions,
    filterOptionBySpace,
  } from './data';

  defineOptions({
    name: 'BizUtilizationRate',
  });
  const { hasPermission } = usePermission();

  const data = reactive({
    exporting: false,
    exportButTitle: '导出数据',
    filters: {},
    sorter: {} as Recordable,
  });
  // 获取前一天
  const yesterday = dayjs().subtract(1, 'day').format('YYYY-MM-DD');
  // 格式化百分比函数，保留两位小数
  const formatPercentage = (value: number | null | undefined): string => {
    if (value == null || value === undefined) return '-';
    return `${(value * 100).toFixed(2)}%`;
  };

  const columns: BasicColumn[] = [
    {
      title: '节点编号',
      dataIndex: 'roomNo',
      width: 120,
      resizable: true,
      fixed: 'left',
      customCell: (_record: any, index: number | undefined) => {
        if (index === undefined) return {};
        const dataSource = safeGetDataSource();
        const rowSpan = calculateRowSpanByNode(dataSource, index, 'roomNo');
        return {
          rowSpan: rowSpan,
        };
      },
    },
    {
      title: '计算日期',
      dataIndex: 'date',
      width: 100,
      resizable: true,
      customCell: (_record: any, index: number | undefined) => {
        if (index === undefined) return {};
        const dataSource = safeGetDataSource();
        const rowSpan = calculateRowSpanByNode(dataSource, index, 'date');
        return {
          rowSpan: rowSpan,
        };
      },
    },
    {
      title: '机房类型',
      dataIndex: 'roomType',
      width: 80,
      resizable: true,
      customCell: (_record: any, index: number | undefined) => {
        if (index === undefined) return {};
        const dataSource = safeGetDataSource();
        const rowSpan = calculateRowSpanByNode(dataSource, index, 'roomType');
        return {
          rowSpan: rowSpan,
        };
      },
    },

    {
      title: '运营商',
      dataIndex: 'localIsp',
      width: 100,
      resizable: true,
      customCell: (_record: any, index: number | undefined) => {
        if (index === undefined) return {};
        const dataSource = safeGetDataSource();
        const rowSpan = calculateRowSpanByNode(dataSource, index, 'localIsp');
        return {
          rowSpan: rowSpan,
        };
      },
    },
    {
      title: '总出口带宽\n(Gbps)',
      dataIndex: 'planBw',
      customHeaderCell: commonCustomHeaderCell(),
      width: 100,
      resizable: true,
      sorter: true,
      customCell: (_record: any, index: number | undefined) => {
        if (index === undefined) return {};
        const dataSource = safeGetDataSource();
        const rowSpan = calTableCellRowSpan(dataSource, index, (current, next) => {
          return (
            current.roomNo === next.roomNo &&
            current.localIsp === next.localIsp &&
            current.location === next.location &&
            current.planBw === next.planBw
          );
        });
        return {
          rowSpan: rowSpan,
        };
      },
      customRender: ({ record }) => {
        if (record.planBw <= 0) return '-';
        return customCeilDivide(record.planBw, 1000 * 1000 * 1000, 2);
      },
    },
    {
      title: '保底带宽\n(Gbps)',
      dataIndex: 'minBw',
      customHeaderCell: commonCustomHeaderCell(),
      width: 100,
      resizable: true,
      sorter: true,
      customCell: (_record: any, index: number | undefined) => {
        if (index === undefined) return {};
        const dataSource = safeGetDataSource();
        const rowSpan = calTableCellRowSpan(dataSource, index, (current, next) => {
          return (
            current.roomNo === next.roomNo &&
            current.localIsp === next.localIsp &&
            current.location === next.location &&
            current.minBw === next.minBw
          );
        });
        return {
          rowSpan: rowSpan,
        };
      },
      customRender: ({ record }) => {
        if (record.minBw <= 0) return '-';
        return customCeilDivide(record.minBw, 1000 * 1000 * 1000, 2);
      },
    },
    {
      title: '计费方式',
      dataIndex: 'billingType',
      width: 100,
      resizable: true,
      customRender: ({ record }) => {
        if (!record.billingType) return '-';
        return record.billingType;
      },
      customCell: (_record: any, index: number | undefined) => {
        if (index === undefined) return {};
        const dataSource = safeGetDataSource();
        const rowSpan = calTableCellRowSpan(dataSource, index, (current, next) => {
          return (
            current.roomNo === next.roomNo &&
            current.localIsp === next.localIsp &&
            current.location === next.location &&
            current.billingType === next.billingType
          );
        });
        return {
          rowSpan: rowSpan,
        };
      },
    },
    {
      title: '节点当前利用率',
      dataIndex: 'curRatio',
      width: 135,
      resizable: true,
      sorter: true,
      customRender: ({ record }) => {
        if (record.curRatio === -1) return '-';
        return `${(record.curRatio * 100).toFixed(2)}%`;
      },
      customCell: (_record: any, index: number | undefined) => {
        if (index === undefined) return {};
        const dataSource = safeGetDataSource();

        // 使用自定义比较函数，比较渲染值是否相同
        const rowSpan = calculateRowSpanByNode(dataSource, index, 'curRatio', (current, next) => {
          // 确保在同一节点编号、同一运营商、同一所在地下比较
          if (
            current.roomNo !== next.roomNo ||
            current.localIsp !== next.localIsp ||
            current.location !== next.location
          ) {
            return false;
          }

          // 获取渲染值
          const getRenderValue = (record) => {
            if (record.curRatio === -1) return '-';
            return `${(record.curRatio * 100).toFixed(2)}%`;
          };

          // 比较两条记录的渲染值是否相同
          return getRenderValue(current) === getRenderValue(next);
        });

        return {
          rowSpan: rowSpan,
        };
      },
    },
    {
      title: '节点晚高峰95利用率',
      dataIndex: 'bwUsageRateNight',
      children: [
        {
          title: '当天',
          dataIndex: 'bwUsageRateNight',
          width: 100,
          sorter: true,
          resizable: true,
          customCell: (_record: any, index: number | undefined) => {
            if (index === undefined) return {};
            const dataSource = safeGetDataSource();

            // 使用自定义比较函数，比较渲染值是否相同
            const rowSpan = calculateRowSpanByNode(
              dataSource,
              index,
              'bwUsageRateNight',
              (current, next) => {
                // 确保在同一节点编号、同一运营商、同一所在地下比较
                if (
                  current.roomNo !== next.roomNo ||
                  current.localIsp !== next.localIsp ||
                  current.location !== next.location
                ) {
                  return false;
                }

                // 获取渲染值
                const getRenderValue = (record) => {
                  return record.bwUsageRateNight;
                };

                // 比较两条记录的渲染值是否相同
                return getRenderValue(current) === getRenderValue(next);
              },
            );

            return {
              rowSpan: rowSpan,
            };
          },
        },
        {
          title: '前一天',
          dataIndex: 'lastOneDayBwUsageRateNight',
          width: 100,
          // sorter: true,
          resizable: true,
          customCell: (_record: any, index: number | undefined) => {
            if (index === undefined) return {};
            const dataSource = safeGetDataSource();

            // 使用自定义比较函数，比较渲染值是否相同
            const rowSpan = calculateRowSpanByNode(
              dataSource,
              index,
              'bwUsageRateNight',
              (current, next) => {
                // 确保在同一节点编号、同一运营商、同一所在地下比较
                if (
                  current.roomNo !== next.roomNo ||
                  current.localIsp !== next.localIsp ||
                  current.location !== next.location
                ) {
                  return false;
                }

                // 获取渲染值
                const getRenderValue = (record) => {
                  return record.bwUsageRateNight;
                };

                // 比较两条记录的渲染值是否相同
                return getRenderValue(current) === getRenderValue(next);
              },
            );

            return {
              rowSpan: rowSpan,
            };
          },
        },
        {
          title: '前两天',
          dataIndex: 'lastTwoDaysBwUsageRateNight',
          width: 100,
          // sorter: true,
          resizable: true,
          customCell: (_record: any, index: number | undefined) => {
            if (index === undefined) return {};
            const dataSource = safeGetDataSource();

            // 使用自定义比较函数，比较渲染值是否相同
            const rowSpan = calculateRowSpanByNode(
              dataSource,
              index,
              'bwUsageRateNight',
              (current, next) => {
                // 确保在同一节点编号、同一运营商、同一所在地下比较
                if (
                  current.roomNo !== next.roomNo ||
                  current.localIsp !== next.localIsp ||
                  current.location !== next.location
                ) {
                  return false;
                }

                // 获取渲染值
                const getRenderValue = (record) => {
                  return record.bwUsageRateNight;
                };

                // 比较两条记录的渲染值是否相同
                return getRenderValue(current) === getRenderValue(next);
              },
            );

            return {
              rowSpan: rowSpan,
            };
          },
        },
      ],
    },
    {
      title: '所在地',
      dataIndex: 'location',
      width: 100,
      resizable: true,
    },
    {
      title: '跨省类型',
      dataIndex: 'trType',
      width: 100,
      resizable: true,
      customRender: ({ record }) => {
        return record.trType || '-';
      },
      customCell: (_record: any, index: number | undefined) => {
        if (index === undefined) return {};
        const dataSource = safeGetDataSource();
        const rowSpan = calculateRowSpanByNode(dataSource, index, 'trType');
        return {
          rowSpan: rowSpan,
        };
      },
    },
    {
      title: '跨省占比',
      dataIndex: 'trZoneRatio',
      width: 100,
      resizable: true,
      sorter: true,
      customRender: ({ record }) => {
        if (record.trZoneRatio === -1) return '-';
        return `${(record.trZoneRatio * 100).toFixed(2)}%`;
      },
      customCell: (_record: any, index: number | undefined) => {
        if (index === undefined) return {};
        const dataSource = safeGetDataSource();
        const rowSpan = calculateRowSpanByNode(dataSource, index, 'trZoneRatio');
        return {
          rowSpan: rowSpan,
        };
      },
    },
    {
      title: '业务大类',
      dataIndex: 'bizCategory', // 修改为实际的数据字段名
      width: 110,
      // sorter: true,
      resizable: true,
    },
    {
      title: '业务',
      dataIndex: 'biz', // 修改为实际的数据字段名
      width: 110,
      // sorter: true,
      resizable: true,
      // filters: [
      //   { text: '主线业务', value: 1 as any },
      //   { text: '其他业务', value: 0 as any },
      // ],
      customRender: ({ record }) => {
        const obj = {
          1: '机房总览',
          2: '保底',
          3: '削峰',
        };
        if (record.reportType === 1) {
          return record.biz;
        } else {
          return record.biz + '-' + record.localIsp + '-' + obj[record.reportType];
        }
      },
    },
    {
      title: '是否主线业务',
      dataIndex: 'isMain', // 修改为实际的数据字段名
      width: 120,
      resizable: true,
      customRender: ({ record }) => {
        return isMainMap[record.isMain];
      },
    },
    {
      title: '带宽占比',
      dataIndex: 'bwRatio',
      width: 120,
      sorter: true,
      resizable: true,
      helpMessage: '较前一日上升5%时，显示为绿色 ，下降5%时，显示为红色',
      customRender: ({ record }) => {
        if (record.bwRatio === -1) return '-';
        if (record.calExtendField && record.calExtendField.bwRatioCompareYesterdayDiff) {
          if (record.calExtendField.bwRatioCompareYesterdayDiff > 0.05) {
            return h(
              'span',
              { style: { color: 'green' } },
              `${(record.bwRatio * 100).toFixed(2)}%`,
            );
          } else if (record.calExtendField.bwRatioCompareYesterdayDiff < -0.05) {
            return h('span', { style: { color: 'red' } }, `${(record.bwRatio * 100).toFixed(2)}%`);
          } else {
            return `${(record.bwRatio * 100).toFixed(2)}%`;
          }
        }
        return `${(record.bwRatio * 100).toFixed(2)}%`;
      },
    },
    // {
    //   title: '节点备注',
    //   dataIndex: 'opsNodeDescribe',
    //   width: 200,
    //   resizable: true,
    //   edit: hasPermission(APermissionCodeEnum.BUSINESS_A_UTILIZATION_RATE_BIZ_NODE_DESCRIBE_EDIT),
    //   editComponent: 'InputTextArea',
    //   editComponentProps: {
    //     style: {
    //       textAlign: 'left',
    //       height: '100%',
    //     },
    //     autoSize: {
    //       minRows: 1,
    //       maxRows: 4,
    //     },
    //     placeholder: '请输入节点备注',
    //     submitOnEnter: false,
    //   },
    //   customRender: () => {},
    //   editRender: ({ record }) => {
    //     const text = record.opsNodeDescribe || '';
    //     const isLongText = text.length > 60;
    //     const displayText = isLongText ? text.substring(0, 60) + '...' : text;
    //     return h(
    //       Tooltip,
    //       {
    //         title: text,
    //         placement: 'topLeft',
    //       },
    //       {
    //         default: () =>
    //           h(
    //             'div',
    //             {
    //               style: {
    //                 wordWrap: 'break-word',
    //                 whiteSpace: 'pre-wrap',
    //                 overflowWrap: 'break-word',
    //                 textAlign: 'left',
    //               },
    //             },
    //             displayText,
    //           ),
    //       },
    //     );
    //   },
    //   customCell: (_record: any, index: number | undefined) => {
    //     if (index === undefined) return {};
    //     const dataSource = safeGetDataSource();
    //     const rowSpan = calculateRowSpanByNode(dataSource, index, 'roomNo');
    //     return {
    //       rowSpan: rowSpan,
    //     };
    //   },
    //   editRule: async (text) => {
    //     if (text.length > 200) {
    //       return '不能超过200字符';
    //     }
    //     return '';
    //   },
    // },
    {
      title: '日备注',
      dataIndex: 'opsDescribe',
      width: 200,
      resizable: true,
      edit: hasPermission(APermissionCodeEnum.BUSINESS_A_UTILIZATION_RATE_BIZ_DAY_DESCRIBE_EDIT),
      editComponent: 'InputTextArea',
      editComponentProps: {
        style: {
          textAlign: 'left',
          height: '100%',
        },
        autoSize: {
          minRows: 1,
          maxRows: 4,
        },
        submitOnEnter: false,
        placeholder: '请输入日备注',
      },
      customRender: () => {},
      editRender: ({ record }) => {
        const text = record.opsDescribe || '';
        const isLongText = text.length > 60;
        const displayText = isLongText ? text.substring(0, 60) + '...' : text;
        return h(
          Tooltip,
          {
            title: text,
            placement: 'topLeft',
          },
          {
            default: () =>
              h(
                'div',
                {
                  style: {
                    wordWrap: 'break-word',
                    whiteSpace: 'pre-wrap',
                    overflowWrap: 'break-word',
                    textAlign: 'left',
                  },
                },
                displayText,
              ),
          },
        );
      },
      // customCell: (_record: any, index: number | undefined) => {
      //   if (index === undefined) return {};
      //   const dataSource = safeGetDataSource();
      //   const rowSpan = calculateRowSpanByNode(dataSource, index, 'roomNo');
      //   return {
      //     rowSpan: rowSpan,
      //   };
      // },
      editRule: async (text) => {
        if (text.length > 200) {
          return '不能超过200字符';
        }
        return '';
      },
    },
    {
      title: '操作人',
      dataIndex: 'modifyUser',
      width: 120,
      resizable: true,
      customRender: ({ record }) => {
        if (record.modifyUser && record.modifyUser.nickName) {
          return record.modifyUser.nickName;
        }
        return '-';
      },
      // customCell: (_record: any, index: number | undefined) => {
      //   if (index === undefined) return {};
      //   const dataSource = safeGetDataSource();
      //   // const rowSpan = calculateRowSpanByNode(dataSource, index, 'roomNo');
      //   const rowSpan = calTableCellRowSpan(dataSource, index, (current, next) => {
      //     return (
      //       current.roomNo != null && current.roomNo != undefined && current.roomNo === next.roomNo
      //     );
      //   });
      //   return {
      //     rowSpan: rowSpan,
      //   };
      // },
    },
    {
      title: '操作时间',
      dataIndex: 'opsUpdateAt',
      width: 160,
      resizable: true,
      // customCell: (_record: any, index: number | undefined) => {
      //   if (index === undefined) return {};
      //   const dataSource = safeGetDataSource();
      //   // const rowSpan = calculateRowSpanByNode(dataSource, index, 'roomNo');
      //   const rowSpan = calTableCellRowSpan(dataSource, index, (current, next) => {
      //     return (
      //       current.roomNo != null && current.roomNo != undefined && current.roomNo === next.roomNo
      //     );
      //   });
      //   return {
      //     rowSpan: rowSpan,
      //   };
      // },
      customRender: ({ record }) => {
        if (record.opsUpdateAt) {
          return record.opsUpdateAt;
        }
        return '-';
      },
    },
  ];
  const searchFormSchema: FormSchema[] = [
    {
      field: 'date',
      label: '计算日期',
      component: 'DatePicker',
      componentProps: {
        placeholder: '请选择日期',
        format: 'YYYY-MM-DD',
        valueFormat: 'YYYY-MM-DD',
        disabledDate: (current: any) => {
          // 不能选择未来的日期
          return current && current >= new Date().setHours(0, 0, 0, 0);
        },
        onChange: onDateChange,
      },
      colProps: { span: 6 },
      required: true,
    },
    {
      field: 'roomType',
      label: '机房类型',
      component: 'Select',
      componentProps: {
        placeholder: '请选择机房类型',
        options: [
          { label: 'IDC', value: 1 },
          { label: 'ACDN', value: 2 },
          { label: 'MCDN', value: 3 },
        ],
        mode: 'multiple',
      },
      colProps: { span: 6 },
    },
    {
      field: 'roomNo',
      label: '节点编号',
      component: 'ApiSelect',
      colProps: { span: 6 },
      componentProps: {
        options: [],
        mode: 'multiple',
        showSearch: true,
        placeholder: '请选择节点编号',
        autoClearSearchValue: false,
        api: async () => {
          const data = await getOwners();
          return data.map((item: any) => ({
            label: item.name,
            value: item.id,
          }));
        },
      },
    },
    {
      field: 'localIsp',
      label: '运营商',
      component: 'Select',
      componentProps: {
        placeholder: '请选择运营商',
        mode: 'multiple',
        options: isps.map((item) => ({ label: item, value: item })),
      },
      colProps: { span: 6 },
    },
    {
      field: 'billingType',
      label: '计费方式',
      component: 'Select',
      componentProps: {
        placeholder: '请选择计费方式',
        options: billingTypeOptions,
        mode: 'multiple',
      },
      colProps: { span: 6 },
    },
    {
      field: 'bizCategory',
      label: '业务大类',
      component: 'Select',
      componentProps: {
        placeholder: '请选择业务大类',
        mode: 'multiple',
        options: [],
        maxTagCount: 3,
        showSearch: true,
        autoClearSearchValue: false,
        filterOption: filterOptionBySpace,
        onChange: () => {
          resetBizs();
        },
      },
      colProps: { span: 6 },
    },
    {
      field: 'bizs',
      label: '业务',
      component: 'Select',
      componentProps: {
        placeholder: '请选择业务',
        mode: 'multiple',
        options: [],
        maxTagCount: 3,
        showSearch: true,
        autoClearSearchValue: false,
        filterOption: filterOptionBySpace,
      },
      colProps: { span: 6 },
    },
    {
      field: 'trType',
      label: '跨省类型',
      component: 'Select',
      componentProps: {
        placeholder: '请选择跨省类型',
        options: trTypeOptions,
        mode: 'multiple',
      },
      colProps: { span: 6 },
    },
    {
      field: 'isMain',
      label: '是否主线业务',
      component: 'Select',
      componentProps: {
        placeholder: '请选择是否主线业务',
        options: Object.entries(isMainMap).map(([key, value]) => ({
          label: value,
          value: key === 'true' ? 1 : 0,
        })),
      },
      colProps: { span: 6 },
    },
    // {
    //   field: 'opsNodeDescribe',
    //   label: '节点备注',
    //   component: 'Input',
    //   componentProps: {
    //     placeholder: '请输入节点备注',
    //   },
    //   colProps: { span: 6 },
    // },
    {
      field: 'opsDescribe',
      label: '日备注',
      component: 'Select',
      componentProps: {
        placeholder: '输入 - 代表无日备注',
        options: [],
        mode: 'tags',
        showSearch: false,
        open: false,
        maxTagCount: 3,
      },
      colProps: { span: 6 },
    },

    // {
    //   field: 'isMain',
    //   label: '业务类型',
    //   component: 'Select',
    //   componentProps: {
    //     options: [
    //       { label: '主线业务', value: 1 },
    //       { label: '其他业务', value: 0 },
    //     ],
    //     mode: 'multiple',
    //   },
    //   colProps: { span: 6 },
    // },
  ];
  function handleTableChange(pagination, filters, sorter) {
    data.filters = filters;
    data.sorter = sorter;
  }

  const beforeEditSubmit = async (data: { record: any; index: number; key: Key; value: any }) => {
    const { record, value, key } = data;
    const params = { id: record.id } as Recordable;
    if (key === 'opsDescribe') {
      params.opsDescribe = value;
    } else if (key === 'opsNodeDescribe') {
      params.opsNodeDescribe = value;
    }
    await EditBizUtilRate(params);
    reload();
  };

  // 安全获取数据源的函数
  const safeGetDataSource = () => {
    try {
      return getDataSource() || [];
    } catch (error) {
      console.warn('getDataSource error:', error);
      return [];
    }
  };
  const [registerTable, { getDataSource, reload, getForm }] = useTable({
    title: '利用率列表',
    beforeEditSubmit,
    api: getBizUtilRateList,
    columns,
    formConfig: {
      labelWidth: 120,
      schemas: searchFormSchema,
      autoSubmitOnEnter: true,
      submitOnReset: false,
      resetFunc() {
        nextTick(() => {
          resetReportTime();
        });
        return Promise.resolve();
      },
      autoAdvancedLine: 2,
    },
    useSearchForm: true,
    bordered: true,
    showIndexColumn: true,
    showTableSetting: true,
  });
  function resetReportTime() {
    return getForm().setFieldsValue({
      date: yesterday,
    });
  }
  async function resetBizs() {
    const values = getForm().getFieldsValue();
    let bizOptions: any[] = [];
    const { bizs, bizCategory } = await getBizs({
      date: values.date,
      bizCategory: values.bizCategory === undefined ? null : values.bizCategory,
    });
    bizs?.forEach((biz) => {
      bizOptions.push({
        label: biz,
        value: biz,
      });
    });
    let bizCategoryOptions: any[] = [];
    bizCategory?.forEach((bizCategory) => {
      bizCategoryOptions.push({
        label: bizCategory,
        value: bizCategory,
      });
    });
    await getForm().updateSchema([
      {
        field: 'bizs',
        componentProps: {
          options: bizOptions,
        },
      },
      {
        field: 'bizCategory',
        componentProps: {
          options: bizCategoryOptions,
        },
      },
    ]);
    await getForm().setFieldsValue({
      bizs: (values.bizs || []).filter((item) => bizs?.includes(item)), // 修改现有选项的值
    });
  }
  onMounted(async () => {
    await resetReportTime();
  });
  function onDateChange() {
    resetBizs();
  }

  function handleExportData() {
    Modal.confirm({
      title: '是否确认导出筛选的数据?',
      icon: h(ExclamationCircleOutlined),
      content: h(
        'div',
        { style: 'color:red;' },
        '导出将花费一些时间，请耐心等待！如果数据过大，可能会导出失败，请缩小导出条件后重试',
      ),
      async onOk() {
        await getForm().validate();
        let value = await getForm().getFieldsValue();
        // 追加过滤条件和排序条件
        value = { ...value, ...data.filters };
        if (value.biz) {
          value.isMain = value.biz;
          delete value.biz;
        }
        if (data?.sorter) {
          value.field = data.sorter.field || null;
          value.order = data.sorter.order || null;
        }
        nextTick(() => {
          data.exporting = true;
          data.exportButTitle = '导出中';
        });
        await exportData(value);
      },
    });
  }

  async function exportData(value: Recordable) {
    try {
      let filename = await downloadFileByUrl(Api.ExportBizUtilRate, 'POST', 5 * 60, value, null);
      notification.success({
        message: '导出成功',
        description: '文件名：' + filename,
        duration: null,
      });
    } catch (error) {
      notification.error({
        message: '导出失败',
        description: error.message,
        duration: null,
      });
    } finally {
      nextTick(() => {
        data.exporting = false;
        data.exportButTitle = '导出数据';
      });
    }
  }
</script>
<style lang="less" scoped>
  :deep(.ant-tabs-nav-list) {
    margin-left: 20px !important;
  }

  .vben-basic-table-form-container {
    padding-top: 0 !important;
  }

  :deep(.vben-basic-table-form-container .ant-form) {
    margin-bottom: 8px !important;
  }

  :deep(.ant-picker) {
    width: 100% !important;
  }
</style>

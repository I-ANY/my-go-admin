<template>
  <div>
    <BasicTable @register="registerTable">
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'action'">
          <TableAction
            :actions="[
              {
                tooltip: '查看请求体',
                // label: '请求体',
                icon: 'ooui:log-in-ltr',
                onClick: handleViewOperaLogReqBody.bind(null, record),
              },
              {
                tooltip: '查看响应体',
                // label: '响应体',
                icon: 'ooui:log-in-rtl',
                onClick: handleViewOperaLogResBody.bind(null, record),
              },
            ]"
          />
        </template>
        <template v-if="column.dataIndex === 'requestMethod'">
          <Tag
            v-if="requestMethodMap[record.requestMethod]"
            :color="requestMethodMap[record.requestMethod].color || 'default'"
            >{{ requestMethodMap[record.requestMethod]?.dictLabel }}</Tag
          >
        </template>
        <template v-if="column.dataIndex === 'handleSource'">
          <Tag
            v-if="handleSourceMap[record.handleSource]"
            :color="handleSourceMap[record.handleSource].color || 'default'"
            >{{ handleSourceMap[record.handleSource]?.dictLabel }}</Tag
          >
          <span v-else></span>
        </template>
        <template v-if="column.dataIndex === 'httpCode'">
          <Tag v-if="record.httpCode < 300" color="green">{{ record.httpCode }}</Tag>
          <Tag v-else-if="record.httpCode < 400" color="orange">{{ record.httpCode }}</Tag>
          <Tag v-else color="red">{{ record.httpCode }}</Tag>
        </template>
        <template v-if="column.dataIndex === 'bizCode'">
          <span v-if="record.bizCode == ''"></span>
          <Tag v-else-if="record.bizCode < 300" color="green">{{ record.bizCode }}</Tag>
          <Tag v-else-if="record.bizCode < 400" color="orange">{{ record.bizCode }}</Tag>
          <Tag v-else color="red">{{ record.bizCode }}</Tag>
        </template>
      </template>
    </BasicTable>
    <JsonDataModal @register="registerJsonModal" />
  </div>
</template>
<script lang="ts">
  import { defineComponent, nextTick, onMounted } from 'vue';
  import { BasicTable, useTable, TableAction } from '@/components/Table';
  import { getAccountList, getOperaLogList } from '@/api/demo/system';
  import JsonDataModal from '../account/JsonDataModal.vue';
  import { useModal } from '@/components/Modal';
  import { Tag } from 'ant-design-vue';
  import {
    operaLogColumns,
    getOperaLogSearchFormSchema,
    requestMethodMap,
    handleSourceMap,
  } from './operalog.data';
  import { useDebounceFn } from '@vueuse/core';
  import dayjs from 'dayjs';
  import { RangePickPresetsExact } from '@/utils/common';

  export default defineComponent({
    name: 'OperaLogManagement',
    components: { BasicTable, TableAction, JsonDataModal, Tag },
    setup() {
      const [registerTable, { getForm }] = useTable({
        title: '操作日志列表',
        api: getOperaLogList,
        // size: 'small',
        fetchSetting: {
          totalField: 'total',
        },
        // pagination: {
        //   pageSizeOptions: ['10', '30', '50', '80', '100', '200'],
        // },
        columns: operaLogColumns,
        formConfig: {
          labelWidth: 90,
          schemas: getOperaLogSearchFormSchema(
            useDebounceFn(queryAccountData, 100),
            onTimePikerOpen,
          ),
          autoSubmitOnEnter: true,
          autoAdvancedLine: 3,
          submitOnReset: false,
          resetFunc() {
            nextTick(() => {
              resetRequestTime();
            });
            return Promise.resolve();
          },
        },
        useSearchForm: true,
        showTableSetting: true,
        bordered: true,
        showIndexColumn: false,
        actionColumn: {
          width: 100,
          title: '操作',
          dataIndex: 'action',
          fixed: 'right',
        },
      });
      const [registerJsonModal, { openModal: openJsonModal, setModalProps }] = useModal();

      async function handleViewOperaLogResBody(record: Recordable) {
        setModalProps({ showOkBtn: false, cancelText: '关闭' });
        openJsonModal(true, { content: record.jsonRes, title: '响应消息体' });
      }
      async function handleViewOperaLogReqBody(record: Recordable) {
        setModalProps({ showOkBtn: false, cancelText: '关闭' });
        openJsonModal(true, { content: record.reqBody, title: '请求消息体' });
      }

      async function queryAccountData(value: string) {
        const searchForm = getForm();
        await searchForm.updateSchema({
          field: 'userId',
          componentProps: {
            loading: true,
          },
        });
        const params = {
          search: value,
          pageSize: 10,
          pageIndex: 1,
        };
        let { items } = await getAccountList(params);
        let options: Array<any> = [];
        items.forEach((e) => {
          let option = {
            ...e,
            label: e.nickName,
            value: e.id,
          };
          options.push(option);
        });
        await searchForm.updateSchema({
          field: 'userId',
          componentProps: {
            loading: false,
            options: options,
          },
        });
      }

      onMounted(async () => {
        queryAccountData('');
        await resetRequestTime();
      });
      async function resetRequestTime(): Promise<void> {
        return getForm().setFieldsValue({
          requestTimeRangeStart: dayjs(
            dayjs().add(-2, 'day').format('YYYY-MM-DD 00:00:00'),
            'YYYY-MM-DD HH:mm:ss',
          ),
          requestTimeRangeEnd: dayjs(dayjs().format('YYYY-MM-DD 23:59:59'), 'YYYY-MM-DD HH:mm:ss'),
        });
      }

      // let timeout: any;

      // // 延迟搜索
      // async function onSearchUser(value: string) {
      //   if (timeout) {
      //     clearTimeout(timeout);
      //     timeout = null;
      //   }
      //   async function queryAccountData() {
      //     const searchForm = getForm();
      //     await searchForm.updateSchema({
      //       field: 'userId',
      //       componentProps: {
      //         loading: true,
      //       },
      //     });
      //     const params = {
      //       search: value,
      //       pageSize: 10,
      //       pageIndex: 1,
      //     };
      //     let { items } = await getAccountList(params);
      //     let options: Array<any> = [];
      //     items.forEach((e) => {
      //       let option = {
      //         ...e,
      //         label: e.nickName,
      //         value: e.id,
      //       };
      //       options.push(option);
      //     });
      //     await searchForm.updateSchema({
      //       field: 'userId',
      //       componentProps: {
      //         loading: false,
      //         options: options,
      //       },
      //     });
      //   }
      //   timeout = setTimeout(queryAccountData, 400);
      // }
      function onTimePikerOpen() {
        getForm().updateSchema({
          field: '[requestTimeRangeStart, requestTimeRangeEnd]',
          componentProps: {
            presets: RangePickPresetsExact(),
          },
        });
      }

      return {
        registerTable,
        registerJsonModal,
        handleViewOperaLogResBody,
        handleViewOperaLogReqBody,
        requestMethodMap,
        handleSourceMap,
      };
    },
  });
</script>

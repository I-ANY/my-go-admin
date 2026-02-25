<template>
  <div>
    <BasicTable @register="registerTable" virtual>
      <template #toolbar>
        <a-button
          type="primary"
          v-auth="PermissionCodeEnum.BUSINESS_INSPECT_SERVER"
          @click="handleServerInspect"
          >设备巡检</a-button
        ></template
      >
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'action'">
          <TableAction
            :actions="[
              {
                // icon: 'el:ok',
                // color: 'error',
                auth: PermissionCodeEnum.BUSINESS_SERVER_INSPECT_RECORD,
                label: '巡检记录',
                onClick: handInspectHistoryShow.bind(null, record),
                tooltip: {
                  title: '查看设备巡检记录',
                  // placement: 'top',
                },
              },
            ]"
          />
        </template>
      </template>
    </BasicTable>
  </div>
</template>
<script lang="ts" setup>
  import { BasicTable, useTable, TableAction } from '@/components/Table';
  import { columns, searchFormSchema, excludeServerTags } from './data';
  import { GetInspectServerList, InspectServer } from '@/api/business/inspect';
  import { GetSubcategoryList } from '@/api/business/biz';
  import { h, nextTick, onMounted } from 'vue';
  import { splitByLineAndTrim } from '@/utils/util';
  import { usePermission } from '@/hooks/web/usePermission';
  import { useGo } from '@/hooks/web/usePage';
  import { message, Modal } from 'ant-design-vue';
  import { ExclamationCircleOutlined } from '@ant-design/icons-vue';
  import { PermissionCodeEnum } from '@/enums/permissionCodeEnum';

  defineOptions({ name: 'InspectServer' });

  const { hasPermission } = usePermission();
  const go = useGo();
  const [registerTable, { getForm, getSelectRowKeys, reload, setSelectedRowKeys }] = useTable({
    title: '设备巡检',
    api: GetInspectServerList,
    beforeFetch: (params) => {
      params.hostnames = splitByLineAndTrim(params.hostnames);
      return params;
    },
    columns,
    formConfig: {
      labelWidth: 120,
      schemas: searchFormSchema(onIdcTypeChange, onCategoryIdsChange),
      autoSubmitOnEnter: true,
      showAdvancedButton: true,
      autoAdvancedLine: 2,
      submitOnReset: false,
      alwaysShowLines: 1,
      resetFunc() {
        nextTick(async () => {
          reloadBusinessOptions(false);
          getForm().setFieldsValue({
            excludeServerTags: excludeServerTags,
          });
        });
        return Promise.resolve();
      },
    },
    rowSelection: {
      type: 'checkbox',
      // getCheckboxProps: () => {
      //   return {
      //     disabled: false,
      //   };
      // },
    },
    rowKey: 'id',
    showSelectionBar: true, // 显示多选状态栏
    clickToRowSelect: false,
    useSearchForm: true,
    showTableSetting: true,
    bordered: true,
    showIndexColumn: true,
    pagination: {
      pageSizeOptions: ['10', '50', '100', '200'],
    },
    actionColumn: {
      width: 200,
      title: '操作',
      dataIndex: 'action',
      // slots: { customRender: 'action' },
      fixed: 'right',
      ifShow: function (): boolean {
        return hasPermission(PermissionCodeEnum.BUSINESS_SERVER_INSPECT_RECORD);
      },
    },
  });

  onMounted(() => {
    reloadBusinessOptions();
  });
  function onIdcTypeChange() {
    reloadBusinessOptions();
  }

  function onCategoryIdsChange() {
    reloadBusinessOptions();
  }

  async function reloadBusinessOptions(includeParmas = true) {
    let selectedBusiness: any[] = [];
    let params = {
      idcType: null,
      categoryIds: null,
    };

    if (includeParmas) {
      const values = await getForm().getFieldsValue();
      params.idcType = values.idcType;
      params.categoryIds = values.categoryIds;
      selectedBusiness = values.business;
    }

    let subcategory = await GetSubcategoryList(params);
    let options: any[] = [];
    subcategory?.items?.forEach((item) => {
      options.push({
        label: item.name,
        value: item.name,
      });
    });
    await getForm().updateSchema([
      {
        field: 'business',
        componentProps: {
          options: options,
        },
      },
    ]);
    // 获取已选择的业务，过滤掉Options中不存在的选项
    let keepBusiness: any[] = [];
    selectedBusiness?.forEach((item) => {
      options.forEach((option) => {
        if (option.value === item) {
          keepBusiness.push(item);
        }
      });
    });
    getForm().setFieldsValue({
      business: keepBusiness,
    });
  }
  function handInspectHistoryShow(record: Recordable) {
    go('/business/inspect/record?hostname=' + record.hostname);
  }
  function handleServerInspect() {
    const selectKeys = getSelectRowKeys();
    if (selectKeys.length === 0) {
      message.warning({ content: '请选择要巡检的设备' });
      return;
    }
    const m = Modal.confirm({
      title: '是否确认下发巡检任务?',
      icon: h(ExclamationCircleOutlined),
      content: h('div', { style: 'color:red;' }, '共' + selectKeys.length + '个设备'),
      async onOk() {
        m.update({
          okButtonProps: {
            loading: true,
          },
        });
        const data = {
          ids: selectKeys,
        };
        try {
          await InspectServer(data);
          message.success({ content: '下发成功' });
          reload();
          setSelectedRowKeys([]);
        } finally {
          m.update({
            okButtonProps: {
              loading: false,
            },
          });
        }
      },
      class: 'serverInspect',
    });
  }
</script>

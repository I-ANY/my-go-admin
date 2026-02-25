<template>
  <BasicModal
    v-bind="$attrs"
    @register="registerModal"
    :title="modalTitle"
    width="1550px"
    :footer="null"
    :destroyOnClose="true"
    :minHeight="100"
    :useWrapper="false"
    style="top: 10px; max-height: calc(100vh - 20px); overflow: auto"
    wrapClassName="history-modal-wrap"
  >
    <BasicTable @register="registerTable" />
  </BasicModal>
</template>

<script lang="ts" setup>
  import { ref, h } from 'vue';
  import { Tag } from 'ant-design-vue';
  import { BasicModal, useModalInner } from '@/components/Modal';
  import { BasicTable, useTable } from '@/components/Table';
  import { GetServerHistoryInfo } from '@/api/business/biz';
  import { formatToDateTime } from '@/utils/dateUtil';

  const modalTitle = ref('主机历史信息');
  const currentHostname = ref('');
  let location = ref('');
  let owner = ref('');

  const historyColumns = [
    {
      title: '机房',
      dataIndex: 'owner',
      width: 100,
      resizable: true,
    },
    {
      title: 'SN',
      dataIndex: 'sn',
      width: 180,
      align: 'left',
    },
    {
      title: '主机名',
      dataIndex: 'hostname',
      width: 180,
      resizable: true,
      customRender: ({ record }) => {
        const isMatch = record.hostname === currentHostname.value;
        return h(
          'span',
          { style: isMatch ? { color: '#589dd5', fontWeight: 'bold' } : {} },
          record.hostname,
        );
      },
    },
    {
      title: '服务器ID',
      dataIndex: 'serverID',
      width: 300,
      align: 'left',
    },
    {
      title: '业务',
      dataIndex: 'business',
      width: 120,
    },
    {
      title: 'FrankID',
      dataIndex: 'frankID',
      width: 380,
      align: 'left',
      customRender: ({ record }) => {
        return record.frankID || '-';
      },
    },
    {
      title: '地区',
      dataIndex: 'location',
      width: 100,
      resizable: true,
    },
    {
      title: '审批状态',
      dataIndex: 'status',
      width: 90,
      resizable: true,
      customRender: ({ record }) => {
        const status = record.status;
        if (status === 0) {
          return h(Tag, { color: 'default' }, () => '待审核');
        } else if (status === 1) {
          return h(Tag, { color: 'green' }, () => '已上线');
        } else if (status === 4) {
          return h(Tag, { color: 'orange' }, () => '已审核');
        } else if (status === 5) {
          return h(Tag, { color: 'red' }, () => '未通过');
        } else if (status === 6) {
          return h(Tag, { color: 'red' }, () => '休眠');
        }
        return status;
      },
    },
    {
      title: '在线状态',
      dataIndex: 'online',
      width: 90,
      resizable: true,
      customRender: ({ record }) => {
        const online = record.online;
        if (online === 1) {
          return h(Tag, { color: 'green' }, () => '在线');
        }
        return h(Tag, { color: 'red' }, () => '离线');
      },
    },
    {
      title: '创建时间',
      dataIndex: 'created_at',
      width: 160,
      customRender: ({ record }) => {
        return record.created_at ? formatToDateTime(record.created_at) : '-';
      },
    },
  ];

  const [registerModal] = useModalInner(async (data) => {
    const sn = data.sn || '';
    currentHostname.value = data.hostname || '';
    location.value = data.location || '';
    owner.value = data.owner || '';
    modalTitle.value = `${data.hostname} 历史业务追溯详情`;
    // 设置查询参数并刷新表格
    setProps({
      searchInfo: {
        sn: sn,
        location: location.value,
      },
    });
    await reload();
  });

  // 执行查询
  const handleOnlyLocalChange = async (e: any) => {
    const checked = e.target.checked;
    if (checked) {
      await getForm().setFieldsValue({ location: location.value, owner: owner.value });
    } else {
      await getForm().setFieldsValue({ location: '', owner: '' });
    }
    // 使用 submit 触发表单查询，确保表单值被正确传递
    await getForm().submit();
  };

  // 当输入框值与原始值不同时，取消勾选
  const createFieldChangeHandler = (field: 'owner' | 'location') => () => {
    const currentValue = getForm().getFieldsValue()[field] || '';
    const originalValue = field === 'owner' ? owner.value : location.value;
    if (currentValue !== originalValue) {
      getForm().setFieldsValue({ onlyLocal: false });
    }
  };

  const [registerTable, { reload, setProps, getForm }] = useTable({
    api: GetServerHistoryInfo,
    columns: historyColumns,
    formConfig: {
      labelWidth: 10,
      schemas: [
        {
          field: 'owner',
          label: ' ',
          component: 'Input',
          colProps: { span: 6 },
          componentProps: {
            onChange: createFieldChangeHandler('owner'),
            allowClear: true,
            placeholder: '输入机房',
          },
        },
        {
          field: 'location',
          label: ' ',
          component: 'Input',
          colProps: { span: 6 },
          componentProps: {
            onChange: createFieldChangeHandler('location'),
            allowClear: true,
            placeholder: '输入地区',
          },
        },
        {
          field: 'onlyLocal',
          label: ' ',
          component: 'Checkbox',
          colProps: { span: 4 },
          renderComponentContent: '仅查询同地区同机房',
          componentProps: {
            onChange: handleOnlyLocalChange,
          },
        },
      ],
      showAdvancedButton: false,
    },
    useSearchForm: true,
    showSelectionBar: false,
    showTableSetting: false,
    bordered: true,
    showIndexColumn: false,
    immediate: false,
    canResize: false,
    pagination: {
      pageSizeOptions: ['10', '30', '50'],
    },
  });
</script>

<style lang="less">
  .history-modal-wrap {
    .ant-modal {
      top: 50px !important;
      max-height: calc(100vh - 70px);
      padding-bottom: 20px;
    }

    .ant-modal-content {
      display: flex;
      flex-direction: column;
      max-height: calc(100vh - 70px);
      overflow: hidden;
    }

    .ant-modal-body {
      flex: 1;
      overflow: auto;
    }
  }
</style>

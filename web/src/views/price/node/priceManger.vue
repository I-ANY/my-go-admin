<template>
  <div>
    <BasicTable @register="registerTable">
      <template #toolbar>
        <!-- <Excel :excelType="excelType" @success-upload="uploadSuccess" /> -->
        <a-button type="primary" @click="handleUpload"> 导入 </a-button>
        <Dropdown>
          <template #overlay>
            <Menu>
              <MenuItem @click="exportData">
                <span> <ExportOutlined /> 导出数据 </span>
              </MenuItem>
              <MenuItem
                :disabled="getSelectRowKeys().length === 0"
                @click="batchOpenModal('price')"
              >
                <span> <EnvironmentOutlined /> 批量编辑采购单价 </span>
              </MenuItem>
              <MenuItem
                :disabled="getSelectRowKeys().length === 0"
                @click="batchOpenModal('priceType')"
              >
                <span> <EnvironmentOutlined /> 批量编辑计费方式 </span>
              </MenuItem>
            </Menu>
          </template>
          <a-button type="primary">更多...</a-button>
        </Dropdown>
      </template>
      <template #bodyCell="{ column, record }">
        <template v-if="column.key == 'priceType'">
          <span>{{ PriceTypeMap[record.priceType] }}</span>
        </template>
        <template v-if="column.key === 'action'">
          <TableAction
            :actions="[
              {
                icon: 'clarity:note-edit-line',
                onClick: (event) => {
                  event.stopPropagation();
                  handleEdit(record);
                },
                tooltip: '编辑',
              },
            ]"
          />
        </template>
      </template>
    </BasicTable>
    <ManagerDrawer @register="registerDrawer" @success="handleSuccess" />
    <Modal
      :bodyStyle="{ height: '300px', overflowY: 'auto' }"
      v-model:open="batchOpen"
      :title="title"
      @ok="batchEdit"
    >
      <BasicForm @register="registerForm" />
    </Modal>
    <ExcelUpload :visible="uploadModalVisible" :close="closeExcelUploadModal" />
  </div>
</template>
<script lang="ts" setup>
  import { Modal, message, Dropdown, Menu, MenuItem } from 'ant-design-vue';
  import { ExportOutlined } from '@ant-design/icons-vue';
  import { ref, reactive } from 'vue';
  import type { Key } from 'ant-design-vue/es/table/interface';
  import dayjs from 'dayjs';
  import { BasicTable, useTable, TableAction, BasicColumn, FormSchema } from '@/components/Table';
  import { getOwners, getNodePrice, updateNodes, exportNodeRecord } from '@/api/price/price';
  import { getAccountList } from '@/api/demo/system';
  import { useDrawer } from '@/components/Drawer';
  import { BasicForm, useForm } from '@/components/Form';
  import ManagerDrawer from './mangerDrawer.vue';
  // import Excel from '../excel.vue';
  import { jsonToSheetXlsx, ExportModalResult } from '@/components/Excel';
  import ExcelUpload from './ExcelUpload.vue';

  defineOptions({ name: 'PriceNodeManger' });

  const batchOpen = ref<boolean>(false);
  const title = ref<string>('');
  const currentEditField = ref<string>('');
  // const excelType = 'nodeRecord';
  const bitchFieldShow = reactive({});
  const serache = reactive({});
  const showSelection = ref(false);
  // 导出弹窗
  const uploadModalVisible = ref(false);

  const PriceTypeMap = {
    0: '未知',
    1: '日95(集群日95)',
    2: '单机日95',
    3: '买断',
    4: '月95',
    5: '单口月95',
  };

  const columns: BasicColumn[] = [
    {
      title: '节点',
      dataIndex: 'name',
      width: 250,
      fixed: 'left',
    },
    {
      title: '所在地',
      dataIndex: 'location',
      width: 200,
      fixed: 'left',
    },
    {
      title: '机房归属',
      dataIndex: 'origin',
      width: 100,
      fixed: 'left',
    },
    {
      title: '计费方式',
      dataIndex: 'priceType',
      fixed: 'left',
    },
    {
      title: '采购单价(元/Gbps)',
      dataIndex: 'price',
      width: 200,
      sorter: true,
    },
    {
      title: '运营商',
      dataIndex: 'localIsp',
      width: 100,
    },
    {
      title: '创建时间',
      dataIndex: 'createdAt',
      customRender: ({ text }) => {
        const date = new Date(text);
        return date.toLocaleString();
      },
    },
    {
      title: '更新时间',
      dataIndex: 'updatedAt',
      customRender: ({ text }) => {
        const date = new Date(text);
        return date.toLocaleString();
      },
    },
    {
      title: '最后更新人',
      dataIndex: 'userName',
      width: 100,
    },
  ];
  const searchFormSchema: FormSchema[] = [
    {
      label: '节点',
      field: 'name',
      component: 'ApiSelect',
      componentProps: {
        placeholder: '请选择节点',
        api: getOwners,
        resultField: 'items',
        labelField: 'name',
        valueField: 'id',
        showSearch: true,
        filterOption: (input: string, option: any) => {
          return (
            option.label?.toLowerCase().indexOf(input.toLowerCase()) >= 0 ||
            option.identify?.toLowerCase().indexOf(input.toLowerCase()) >= 0
          );
        },
      },
      colProps: { span: 6 },
    },
    {
      field: 'location',
      label: '所在地',
      component: 'Select',
      componentProps: {
        placeholder: '请选择省份',
        showSearch: true,
        options: [
          { label: '北京', value: '北京' },
          { label: '天津', value: '天津' },
          { label: '上海', value: '上海' },
          { label: '重庆', value: '重庆' },
          { label: '河北', value: '河北' },
          { label: '山西', value: '山西' },
          { label: '辽宁', value: '辽宁' },
          { label: '吉林', value: '吉林' },
          { label: '黑龙江', value: '黑龙江' },
          { label: '江苏', value: '江苏' },
          { label: '浙江', value: '浙江' },
          { label: '安徽', value: '安徽' },
          { label: '福建', value: '福建' },
          { label: '江西', value: '江西' },
          { label: '山东', value: '山东' },
          { label: '河南', value: '河南' },
          { label: '湖北', value: '湖北' },
          { label: '湖南', value: '湖南' },
          { label: '广东', value: '广东' },
          { label: '海南', value: '海南' },
          { label: '四川', value: '四川' },
          { label: '贵州', value: '贵州' },
          { label: '云南', value: '云南' },
          { label: '陕西', value: '陕西' },
          { label: '甘肃', value: '甘肃' },
          { label: '青海', value: '青海' },
          { label: '内蒙古', value: '内蒙古' },
          { label: '广西', value: '广西' },
          { label: '西藏', value: '西藏' },
          { label: '宁夏', value: '宁夏' },
          { label: '新疆', value: '新疆' },
        ],
      },
      colProps: { span: 6 },
    },
    {
      field: 'origin',
      label: '机房归属',
      component: 'Select',
      componentProps: {
        placeholder: '请选择机房归属',
        options: [
          { label: '招募', value: '招募' },
          { label: '自建', value: '自建' },
        ],
      },
      colProps: { span: 6 },
    },
    {
      field: 'localIsp',
      label: '运营商',
      component: 'Select',
      componentProps: {
        placeholder: '请选择运营商',
        mode: 'multiple',
        options: [
          { label: '电信', value: '电信' },
          { label: '移动', value: '移动' },
          { label: '联通', value: '联通' },
        ],
      },
      colProps: { span: 6 },
    },
    {
      label: '采购单价',
      field: 'startPrice',
      component: 'InputNumber',
      componentProps: {
        placeholder: '最低价',
        precision: 2,
      },
      colProps: { span: 4 },
    },
    {
      label: '',
      field: 'endPrice',
      component: 'InputNumber',
      componentProps: {
        placeholder: '最高价',
        precision: 2,
      },
      colProps: { span: 2 },
    },
    {
      field: 'priceType',
      label: '计费方式',
      component: 'Select',
      componentProps: {
        placeholder: '请选择计费方式',
        options: [
          { label: '未知', value: 0 },
          { label: '日95(集群日95)', value: 1 },
          { label: '单机日95', value: 2 },
          { label: '买断', value: 3 },
          { label: '月95', value: 4 },
          { label: '单口月95', value: 5 },
        ],
      },
      colProps: { span: 6 },
    },
    {
      field: 'updateAt',
      label: '更新时间',
      component: 'RangePicker',
      componentProps: {
        format: 'YYYY-MM-DD HH:mm:ss',
        placeholder: ['开始日期、时间', '结束日期、时间'],
        showTime: { format: 'HH:mm:ss' },
      },
      colProps: { span: 6 },
    },
    {
      field: 'updateBy',
      label: '最后更新人',
      component: 'ApiSelect',
      componentProps: {
        placeholder: '请选择更新人',
        api: getAccountList,
        params: { pageSize: 5000, pageIndex: 1 },
        resultField: 'items',
        labelField: 'nickName',
        valueField: 'id',
        showSearch: true,
        filterOption: (input, option) => {
          return option.label.toLowerCase().includes(input.toLowerCase());
        },
      },
      colProps: { span: 6 },
    },
  ];

  const formSchemaBatch: FormSchema[] = [
    {
      field: 'price',
      label: '采购单价(元/Gbps)',
      required: true,
      component: 'InputNumber',
      componentProps: {
        precision: 2,
      },
      colProps: { span: 16 },
      ifShow: () => bitchFieldShow['price'],
    },
    {
      field: 'priceType',
      label: '计费方式',
      required: true,
      component: 'Select',
      componentProps: {
        placeholder: '请选择计费方式',
        options: [
          { label: '未知', value: 0 },
          { label: '日95(集群日95)', value: 1 },
          { label: '单机日95', value: 2 },
          { label: '买断', value: 3 },
          { label: '月95', value: 4 },
          { label: '单口月95', value: 5 },
        ],
      },
      colProps: { span: 16 },
      ifShow: () => bitchFieldShow['priceType'],
    },
  ];

  const [registerForm, { resetFields, getFieldsValue }] = useForm({
    schemas: formSchemaBatch,
    showActionButtonGroup: false,
  });

  const [registerDrawer, { openDrawer }] = useDrawer();
  const [registerTable, { reload, clearSelectedRowKeys, getSelectRowKeys }] = useTable({
    title: '节点单价列表',
    api: getNodePrice,
    columns,
    formConfig: {
      labelWidth: 120,
      schemas: searchFormSchema,
      autoSubmitOnEnter: true,
      autoAdvancedLine: 5,
      showAdvancedButton: false,
      actionColOptions: {
        span: 18,
        style: { textAlign: 'right', paddingTop: '4px' },
      },
    },
    useSearchForm: true,
    bordered: true,
    showIndexColumn: true,
    showTableSetting: true,
    actionColumn: {
      width: 80,
      title: '操作',
      dataIndex: 'action',
      fixed: 'right',
    },
    rowKey: 'id',
    rowSelection: {
      type: 'checkbox',
      onChange: (selectedRowKeys: Key[]) => {
        showSelection.value = selectedRowKeys.length > 0;
      },
    },
    showSelectionBar: showSelection,
    clickToRowSelect: false,
  });

  function handleEdit(record: Recordable) {
    openDrawer(true, {
      record,
      isUpdate: true,
    });
  }

  function handleSuccess() {
    reload();
    let msg = '编辑成功';
    message.success(msg);
  }

  // function uploadSuccess(resp: string) {
  //   reload();
  //   message.success(resp, 10);
  // }

  function clearSelect() {
    clearSelectedRowKeys();
    showSelection.value = false;
  }

  function batchOpenModal(field: string) {
    batchOpen.value = true;
    currentEditField.value = field;
    resetFields();
    if (field === 'price') {
      title.value = '编辑单价';
      bitchFieldShow['price'] = true;
      bitchFieldShow['priceType'] = false;
      return;
    }
    if (field === 'priceType') {
      title.value = '编辑计费类型';
      bitchFieldShow['price'] = false;
      bitchFieldShow['priceType'] = true;
      return;
    }
  }

  function batchEdit() {
    let data = { id: getSelectRowKeys() };
    let values = getFieldsValue();

    if (currentEditField.value === 'price') {
      if (values.price <= 0) {
        message.error('单价禁止设置为0');
        return;
      }
      data['field'] = 'price';
      data['valueFloat'] = values.price;
    } else if (currentEditField.value === 'priceType') {
      if (values.priceType === undefined || values.priceType === null) {
        message.error('请选择计费方式');
        return;
      }
      data['field'] = 'priceType';
      data['value'] = values.priceType;
    }

    updateNodes(data).finally(() => {
      resetFields();
      clearSelect();
      reload();
      batchOpen.value = false;
    });
  }

  function exportData() {
    // openModal();
    const timestamp = dayjs().format('YYYYMMDDHHmm');
    defaultHeader({
      filename: `单价导出数据_${timestamp}.xlsx`,
      bookType: 'xlsx',
    });
  }

  function defaultHeader({ filename, bookType }: ExportModalResult) {
    console.log(filename);
    exportNodeRecord({
      ...serache,
      id: getSelectRowKeys(),
    })
      .then((resp) => {
        let data = resp.items;
        let result: { [key: string]: any }[] = [];
        for (let i = 0; i < data.length; i++) {
          let map = {};
          map['节点'] = data[i]['节点'];
          map['所在地'] = data[i]['所在地'];
          map['机房归属'] = data[i]['机房归属'];
          map['计费方式'] = data[i]['计费方式'];
          map['采购单价(元/Gbps)'] = data[i]['采购单价'];
          map['运营商'] = data[i]['运营商'];
          map['创建时间'] = data[i]['创建时间'];
          map['更新时间'] = data[i]['更新时间'];
          map['最后更新人'] = data[i]['更新人'];
          result.push(map);
        }
        jsonToSheetXlsx({
          data: result,
          filename: filename,
          write2excelOpts: {
            bookType,
          },
        });
      })
      .finally(() => {
        clearSelect();
      });
  }
  function handleUpload() {
    uploadModalVisible.value = true;
  }
  function closeExcelUploadModal() {
    uploadModalVisible.value = false;
    reload();
  }
</script>

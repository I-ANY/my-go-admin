<template>
  <div>
    <BasicTable @register="registerTable">
      <template #toolbar>
        <Excel :excelType="excelType" @success-upload="uploadSuccess" />
        <a-button type="primary" @click="handleCreate"> 新增 </a-button>
        <Dropdown>
          <template #overlay>
            <Menu>
              <MenuItem @click="exportData">
                <span> <ExportOutlined /> 导出数据 </span>
              </MenuItem>
              <MenuItem :disabled="getSelectRowKeys().length === 0">
                <Popconfirm
                  title="确认删除?"
                  ok-text="确认"
                  cancel-text="取消"
                  @confirm="batchDelete"
                >
                  <span> <DeleteOutlined /> 批量删除 </span>
                </Popconfirm>
              </MenuItem>
              <MenuItem
                :disabled="getSelectRowKeys().length === 0"
                @click="batchOpenModal('region_id')"
              >
                <span> <EnvironmentOutlined /> 批量调整大区 </span>
              </MenuItem>
              <MenuItem
                :disabled="getSelectRowKeys().length === 0"
                @click="batchOpenModal('zone_id')"
              >
                <span> <EnvironmentOutlined /> 批量调整省份 </span>
              </MenuItem>
              <MenuItem
                :disabled="getSelectRowKeys().length === 0"
                @click="batchOpenModal('price')"
              >
                <span> <DollarOutlined /> 批量调整单价 </span>
              </MenuItem>
              <MenuItem
                :disabled="getSelectRowKeys().length === 0"
                @click="batchOpenModal('mode_id')"
              >
                <span> <FunctionOutlined /> 批量调整计费方式 </span>
              </MenuItem>
            </Menu>
          </template>
          <a-button type="primary">更多...</a-button>
        </Dropdown>
      </template>
      <template #headerCell="{ column }">
        <template v-if="column.key === 'price'">
          <span v-if="showPrice"> <EyeOutlined @click="togglePriceColumn" /> 单价(元)</span>
          <span v-else> <EyeInvisibleOutlined @click="togglePriceColumn" /> 单价(元)</span>
        </template>
      </template>
      <template #bodyCell="{ column, record }">
        <template v-if="column.key == 'category'">
          <span>{{ record.category ? record.category.name : '/' }}</span>
        </template>
        <template v-if="column.key == 'bizs'">
          <a-button type="primary" @click.stop="handleQueryBiz(record.categoryID)">详情</a-button>
        </template>
        <template v-if="column.key == 'region'">
          {{ record.region ? record.region.name : '/' }}
        </template>
        <template v-if="column.key == 'zone'">
          {{ record.zone ? record.zone.name : '/' }}
        </template>
        <template v-if="column.key == 'isNotlocalIsp'">
          {{ record.bizIspMode || record.bizIsp ? '是' : '否' }}
        </template>
        <template v-if="column.key == 'bizIspMode'">
          {{ record.bizIspMode == 2 ? '异网计费' : record.bizIspMode == 1 ? '本网计费' : '/' }}
        </template>
        <template v-if="column.key == 'bizIsp'">
          {{ record.bizIsp ? record.bizIsp : '/' }}
        </template>
        <template v-if="column.key == 'mode'">
          {{ record.mode ? record.mode.name : '/' }}
        </template>
        <template v-if="column.key === 'price'">
          <span v-if="showPrice">{{ record.price }}</span>
          <span v-else>***</span>
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
              {
                icon: 'ant-design:delete-outlined',
                tooltip: '删除',
                popConfirm: {
                  title: '确认删除?',
                  placement: 'left',
                  confirm: handleDelete.bind(null, record),
                },
              },
            ]"
          />
        </template>
      </template>
    </BasicTable>
    <ManagerDrawer @register="registerDrawer" @success="handleSuccess" />
    <Modal
      :bodyStyle="{ height: '300px', overflowY: 'auto' }"
      v-model:open="open"
      title="业务详情"
      @ok="open = false"
    >
      <template v-for="item in bizData.items" :key="item.id">
        <div
          style="
            display: flex;
            align-items: center;
            padding: 8px 16px;
            border-bottom: 1px solid #f0f0f0;
          "
        >
          <span style="flex: 3; text-align: left">{{ item.name || '无业务' }}</span>
        </div>
      </template>
    </Modal>
    <Modal
      :bodyStyle="{ height: '300px', overflowY: 'auto' }"
      v-model:open="batchOpen"
      :title="title"
      @ok="batchEdit"
    >
      <BasicForm @register="registerForm" />
    </Modal>
    <!-- <ExpExcelModal @register="register" @success="defaultHeader" /> -->
  </div>
</template>
<script lang="ts" setup>
  import { Modal, message, Dropdown, Menu, MenuItem, Popconfirm } from 'ant-design-vue';
  import {
    EyeOutlined,
    EyeInvisibleOutlined,
    DollarOutlined,
    FunctionOutlined,
    EnvironmentOutlined,
    DeleteOutlined,
    ExportOutlined,
  } from '@ant-design/icons-vue';
  import { ref, reactive } from 'vue';
  import type { Key } from 'ant-design-vue/es/table/interface';
  import dayjs from 'dayjs';
  import { BasicTable, useTable, TableAction, BasicColumn, FormSchema } from '@/components/Table';
  import {
    getRecord,
    deleteRecord,
    deleteRecords,
    getBiz,
    getRegion,
    getMode,
    getZone,
    updateBatchNumPointRecord,
    updateBatchNumRecord,
    exportRecord,
    getCategory,
  } from '@/api/price/price';
  import { useDrawer } from '@/components/Drawer';
  import { BasicForm, useForm } from '@/components/Form';
  import ManagerDrawer from './mangerDrawer.vue';
  import Excel from '../excel.vue';
  // import { useModal } from '@/components/Modal';
  // import { ExpExcelModal } from '@/components/Excel';
  import { jsonToSheetXlsx, ExportModalResult } from '@/components/Excel';

  defineOptions({ name: 'PriceManger' });

  // const [register, { openModal }] = useModal();

  const open = ref<boolean>(false);
  const batchOpen = ref<boolean>(false);
  const title = ref<string>('');
  const showPrice = ref(false);
  const bizData = reactive({ items: [] as { id: number; name: string }[] });
  const excelType = 'record';
  const bitchFieldShow = reactive({});
  const serache = reactive({});
  const showSelection = ref(false);

  const columns: BasicColumn[] = [
    {
      title: '本网运营商',
      dataIndex: 'localIsp',
      width: 100,
      fixed: 'left',
    },
    {
      title: '业务组',
      dataIndex: 'category',
      width: 200,
      fixed: 'left',
    },
    {
      title: '业务详情',
      dataIndex: 'bizs',
      width: 100,
      fixed: 'left',
    },
    {
      title: '溜缝业务',
      dataIndex: 'low',
      width: 100,
      customRender: ({ text }) => (text ? '是' : '否'),
    },
    {
      title: '大区',
      dataIndex: 'region',
      width: 100,
    },
    {
      title: '省份',
      dataIndex: 'zone',
      width: 100,
    },
    {
      title: '是否异网',
      dataIndex: 'isNotlocalIsp',
      width: 100,
    },
    {
      title: '异网运营商',
      dataIndex: 'bizIsp',
      width: 100,
    },
    {
      title: '跨网计费方式',
      dataIndex: 'bizIspMode',
      width: 100,
    },
    {
      title: '计费方式',
      dataIndex: 'mode',
      width: 300,
    },
    {
      title: '单价(元)',
      dataIndex: 'price',
      width: 100,
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
      title: '备注',
      dataIndex: 'describe',
    },
  ];
  const searchFormSchema: FormSchema[] = [
    {
      label: '业务组',
      field: 'name',
      component: 'ApiSelect',
      componentProps: {
        api: getCategory,
        params: { pageSize: 50000, pageIndex: 1 },
        resultField: 'items',
        labelField: 'name',
        valueField: 'name',
        showSearch: true,
        mode: 'multiple',
        filterOption: (input: string, option: any) => {
          return (
            option.label?.toLowerCase().indexOf(input.toLowerCase()) >= 0 ||
            option.identify?.toLowerCase().indexOf(input.toLowerCase()) >= 0
          );
        },
        onChange: (v) => {
          serache['name'] = v;
          clearSelect();
        },
      },
      colProps: { span: 6 },
    },
    {
      field: 'localIsp',
      label: '本网运营商',
      component: 'Select',
      componentProps: {
        options: [
          { label: '电信', value: '电信' },
          { label: '移动', value: '移动' },
          { label: '联通', value: '联通' },
        ],
        onChange: (v) => {
          serache['localIsp'] = v;
          clearSelect();
        },
      },
      colProps: { span: 6 },
    },
    {
      field: 'diffNet',
      label: '是否异网',
      component: 'Select',
      componentProps: {
        options: [
          { label: '否', value: 0 },
          { label: '是', value: 1 },
        ],
        onChange: (v) => {
          serache['diffNet'] = v;
          clearSelect();
        },
      },
      colProps: { span: 6 },
    },
    {
      label: '大区',
      field: 'regionID',
      component: 'ApiSelect',
      componentProps: {
        api: getRegion,
        params: { pageSize: 5000, pageIndex: 1 },
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
      label: '省份',
      field: 'zoneID',
      component: 'ApiSelect',
      componentProps: {
        api: getZone,
        params: { pageSize: 5000, pageIndex: 1 },
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
  ];

  const formSchema: FormSchema[] = [
    {
      label: '大区',
      field: 'regionID',
      component: 'ApiSelect',
      componentProps: {
        api: getRegion,
        params: { pageSize: 5000, pageIndex: 1 },
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
      ifShow: () => {
        return bitchFieldShow['regionID'];
      },
    },
    {
      label: '省份',
      field: 'zoneID',
      component: 'ApiSelect',
      componentProps: {
        api: getZone,
        params: { pageSize: 5000, pageIndex: 1 },
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
      ifShow: () => {
        return bitchFieldShow['zoneID'];
      },
    },
    {
      label: '计费方式',
      field: 'modeID',
      component: 'ApiSelect',
      colProps: { span: 20 },
      componentProps: {
        api: getMode,
        params: { pageSize: 50000, pageIndex: 1 },
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
      ifShow: () => {
        return bitchFieldShow['modeID'];
      },
      required: true,
    },
    {
      field: 'price',
      label: '单价(元)',
      component: 'InputNumber',
      ifShow: () => {
        return bitchFieldShow['price'];
      },
      required: true,
    },
  ];

  const [registerForm, { resetFields, getFieldsValue }] = useForm({
    labelWidth: 150,
    // baseColProps: { span: 24 },
    // baseColProps: { lg: 12, md: 24 },
    schemas: formSchema,
    showActionButtonGroup: false,
  });

  const [registerDrawer, { openDrawer }] = useDrawer();
  const [registerTable, { reload, clearSelectedRowKeys, getSelectRowKeys }] = useTable({
    title: '业务单价列表',
    api: getRecord,
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

  function handleCreate() {
    openDrawer(true, {
      isUpdate: false,
    });
  }

  function handleEdit(record: Recordable) {
    openDrawer(true, {
      record,
      isUpdate: true,
    });
  }

  function handleQueryBiz(id: Recordable) {
    bizData.items = [];
    getBiz({ id: id })
      .then(function (resp) {
        for (var i = 0; i < resp.length; i++) {
          bizData.items.push(resp[i]);
        }
        open.value = true;
      })
      .catch(function (err) {
        message.error(err);
      });
  }

  async function handleDelete(record: Recordable) {
    await deleteRecord(record.id);
    message.success('删除成功');
    reload();
  }

  function handleSuccess(isAdd: boolean) {
    reload();
    let msg = '新增成功';
    if (!isAdd) {
      msg = '编辑成功';
    }
    message.success(msg);
  }

  function uploadSuccess(resp: string) {
    reload();
    message.success(resp, 10);
  }

  function togglePriceColumn() {
    showPrice.value = !showPrice.value;
  }

  function clearSelect() {
    clearSelectedRowKeys();
    showSelection.value = false;
  }

  async function batchDelete() {
    deleteRecords({ id: getSelectRowKeys() }).then(() => {
      clearSelect();
      reload();
    });
  }

  function batchOpenModal(field: string) {
    batchOpen.value = true;
    if (field === 'region_id') {
      title.value = '编辑大区';
      bitchFieldShow['regionID'] = true;
      bitchFieldShow['zoneID'] = false;
      bitchFieldShow['modeID'] = false;
      bitchFieldShow['price'] = false;
      return;
    }
    if (field === 'zone_id') {
      title.value = '编辑省份';
      bitchFieldShow['regionID'] = false;
      bitchFieldShow['zoneID'] = true;
      bitchFieldShow['modeID'] = false;
      bitchFieldShow['price'] = false;
      return;
    }
    if (field === 'mode_id') {
      title.value = '编辑计费方式';
      bitchFieldShow['regionID'] = false;
      bitchFieldShow['zoneID'] = false;
      bitchFieldShow['modeID'] = true;
      bitchFieldShow['price'] = false;
      return;
    }
    if (field === 'price') {
      title.value = '编辑单价';
      bitchFieldShow['regionID'] = false;
      bitchFieldShow['zoneID'] = false;
      bitchFieldShow['modeID'] = false;
      bitchFieldShow['price'] = true;
      return;
    }
  }

  function batchEdit() {
    let field = '';
    let data = { id: getSelectRowKeys() };
    let values = getFieldsValue();
    for (const key in bitchFieldShow) {
      if (bitchFieldShow[key]) {
        field = key;
        break;
      }
    }
    if (field === 'regionID') {
      data['field'] = 'region_id';
      if (values['regionID'] !== undefined) {
        data['value'] = values['regionID'];
      }
      updateBatchNumPointRecord(data).finally(() => {
        resetFields();
        clearSelect();
        reload();
        batchOpen.value = false;
      });
    }
    if (field === 'zoneID') {
      data['field'] = 'zone_id';
      if (values['zoneID'] !== undefined) {
        data['value'] = values['zoneID'];
      }
      updateBatchNumPointRecord(data).finally(() => {
        resetFields();
        clearSelect();
        reload();
        batchOpen.value = false;
      });
    }
    if (field === 'modeID') {
      data['field'] = 'mode_id';
      if (values['modeID'] !== undefined) {
        data['value'] = values['modeID'];
      } else {
        data['value'] = 0;
        return;
      }
      updateBatchNumRecord(data).finally(() => {
        resetFields();
        clearSelect();
        reload();
        batchOpen.value = false;
      });
    }
    if (field === 'price') {
      data['field'] = 'price';
      if (values['price'] !== undefined) {
        data['value'] = values['price'];
      } else {
        data['value'] = 0;
        return;
      }
      updateBatchNumRecord(data).finally(() => {
        resetFields();
        clearSelect();
        reload();
        batchOpen.value = false;
      });
    }
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
    exportRecord({
      ...serache,
      id: getSelectRowKeys(),
    })
      .then((resp) => {
        let data = resp.items;
        let result: { [key: string]: any }[] = [];
        for (let i = 0; i < data.length; i++) {
          let map = {};
          map['本网运营商'] = data[i]['本网运营商'];
          map['业务组'] = data[i]['业务组'];
          map['溜缝业务'] = data[i]['溜缝业务'];
          map['大区'] = data[i]['大区'];
          map['省份'] = data[i]['省份'];
          map['是否异网'] = data[i]['是否异网'];
          map['异网运营商'] = data[i]['异网运营商'];
          map['跨网计费方式'] = data[i]['跨网计费方式'];
          map['计费方式'] = data[i]['计费方式'];
          map['更新时间'] = data[i]['更新时间'];
          map['单价'] = data[i]['单价'];
          map['备注'] = data[i]['备注'];
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
</script>

<template>
  <div>
    <BasicDrawer
      title="资源权限"
      @register="registerModal"
      :width="1200"
      @ok="handleOk"
      okText="保存"
      cancelText="取消"
      @close="
        () => {
          emit('reload');
        }
      "
    >
      <Spin :spinning="allData.modalSpinning">
        <div class="resource-checkbox-container">
          <BasicForm @register="registerForm" @field-value-change="onValuesChange" />
        </div>
        <div>
          <Divider style="margin: 0">授权所有资源</Divider>
          <FormItem
            label="权限类型"
            :wrapper-col="{ span: 18 }"
            labelAlign="right"
            :label-col="{ span: 2 }"
            style="margin-left: 13px"
          >
            <CheckboxGroup
              v-model:value="allData.authedAllResourcePermissionTypes"
              @change="onAuthedAllResourcePermissionTypesChange()"
              :options="
                allData?.allPermissionTypes?.map((item) => ({
                  label: item.name,
                  value: item.code,
                })) || []
              "
            />
          </FormItem>
          <Divider style="margin: 0">授权部分资源</Divider>
          <BasicTable @register="registerTable">
            <template #toolbar>
              <a-button type="default" @click="handleResetChange">重置修改</a-button>
              <Dropdown>
                <template #overlay>
                  <Menu>
                    <MenuItem
                      v-for="(item, i) in allData.allPermissionTypes"
                      :key="'permissionType' + i"
                      :disabled="allData.authedAllResourcePermissionTypes.includes(item.code)"
                      @click="handleSelectPermissionType(item.code)"
                    >
                      {{ `全选本页-${item.name}` }}
                    </MenuItem>
                    <MenuDivider />
                    <MenuItem key="1" @click="handleSelectPage(true)"> 全选本页-所有权限 </MenuItem>
                    <MenuItem key="2" @click="handleSelectPage(false)">
                      取消本页-所有权限
                    </MenuItem>
                  </Menu>
                </template>
                <a-button>
                  更多操作
                  <DownOutlined />
                </a-button>
              </Dropdown>
            </template>
            <template #bodyCell="{ column, record }">
              <!-- 使用枚举值展示 -->
              <template v-for="(columnName, i) in Object.keys(allData.showEnumFields)" :key="i">
                <template v-if="column.key == columnName">
                  <span v-if="allData.showEnumFields[columnName][record[columnName]]">
                    {{ allData.showEnumFields[columnName][record[columnName]].dictLabel }}
                  </span>
                  <span v-else>{{ record[columnName] }}</span>
                </template>
              </template>

              <!-- 使用枚举值tag展示 -->
              <template v-for="(columnName, i) in Object.keys(allData.showTagFields)" :key="i">
                <template v-if="column.key == columnName">
                  <Tag
                    style="font-weight: bold"
                    v-if="allData.showTagFields[columnName][record[columnName]]"
                    :color="
                      allData.showTagFields[columnName][record[columnName]].color || 'default'
                    "
                    >{{ allData.showTagFields[columnName][record[columnName]].dictLabel }}</Tag
                  >
                  <span v-else>{{ record[columnName] }}</span>
                </template>
              </template>
              <!-- 权限类型 -->
              <template v-if="column.key == 'permissionTypes'">
                <div style="display: flex; align-items: center">
                  <Checkbox
                    v-model:checked="allData.permissionTypeCheckStatus[record.id].checkAll"
                    :indeterminate="allData.permissionTypeCheckStatus[record.id].indeterminate"
                    @change="onPermissionTypeCheckAllChange(record.id)"
                    >全选</Checkbox
                  >
                  <Divider type="vertical" :style="{ height: '15px' }" />
                  <CheckboxGroup
                    v-model:value="allData.permissionTypeCheckStatus[record.id].permissionTypes"
                    @change="onPermissionTypeCheckChange(record.id)"
                    :options="
                      allData.allPermissionTypes.map((item) => ({
                        label: item.name,
                        value: item.code,
                        disabled: allData.authedAllResourcePermissionTypes.includes(item.code),
                      })) || []
                    "
                  />
                </div>
              </template>
            </template>
          </BasicTable>
        </div>
      </Spin>
    </BasicDrawer>
  </div>
</template>

<script lang="ts" setup>
  import { BasicDrawer, useDrawerInner } from '@/components/Drawer';
  import { BasicTable, useTable } from '@/components/Table';
  import { BasicForm, useForm } from '@/components/Form';
  import {
    getRoleResourceDetailList,
    getResourceViewColumns,
    getResourceViewSearchFormSchemas,
    getRoleResourceInfo,
    updateRoleResource,
  } from '@/api/sys/resource';
  import { DownOutlined } from '@ant-design/icons-vue';
  import { reactive } from 'vue';
  import {
    Spin,
    Tag,
    Divider,
    CheckboxGroup,
    Checkbox,
    Modal,
    message,
    Dropdown,
    MenuItem,
    Menu,
    FormItem,
    MenuDivider,
  } from 'ant-design-vue';
  import { getDictDataMapFromDict } from '@/utils/dict';
  import { getResourcePermissionFormSchema } from './role.data';

  defineOptions({ name: 'ResourcePermission' });
  const emit = defineEmits(['reload', 'register']);
  const allData = reactive({
    record: {} as Recordable,
    currentResourceTypeId: 0,
    showEnumFields: {},
    showTagFields: {},
    modalSpinning: true,
    permissionTypeCheckStatus: {}, // 权限类型全选状态
    isResetChange: false, // 是否重置授权部分资源权限修改
    allPermissionTypes: [] as Recordable[], // 所有权限类型
    authedAllResourcePermissionTypes: [] as string[], // 已授权的所有资源的权限类型
  });
  function cleanAllData() {
    allData.record = {};
    allData.currentResourceTypeId = 0;
    allData.showEnumFields = {};
    allData.showTagFields = {};
    allData.modalSpinning = true;
    allData.permissionTypeCheckStatus = {};
    allData.isResetChange = false;
    allData.allPermissionTypes = [];
    allData.authedAllResourcePermissionTypes = [];
  }
  const [registerForm, { setFieldsValue, updateSchema, getFieldsValue, validate }] = useForm({
    schemas: getResourcePermissionFormSchema(onResourceTypeChange, getFuncs),
    showActionButtonGroup: false,
    labelWidth: 110,
  });
  const [registerTable, { reload, setColumns, setProps, getDataSource, setTableData }] = useTable({
    title: '资源权限',
    api: (params) => {
      // 直接调用 API，传递正确的参数
      return getRoleResourceDetailList(
        allData.record.id, // roleId
        {
          ...params,
          resourceTypeId: allData.currentResourceTypeId,
        },
      );
    },
    afterFetch: (data) => {
      // 需要重置所有已选择的权限类型
      if (allData.isResetChange) {
        allData.permissionTypeCheckStatus = {};
      } else {
        // 每次加载数据都需要将没有修改过的ID删除，仅保留修改过的ID，防止分页次数过多导致页面数据量过大
        Object.keys(allData.permissionTypeCheckStatus).forEach((key) => {
          const id = parseInt(key);
          // 这个记录没有改变，删除
          if (!allData.permissionTypeCheckStatus[id].changed) {
            allData.permissionTypeCheckStatus[id] = {
              checkAll: false,
              indeterminate: false,
              permissionTypes: [],
              changed: false,
            };
          }
        });
      }
      data.forEach((item) => {
        // 数据中不存在这个ID，或者这个ID没有被修改过，则每次加载的时候需要重新渲染
        if (
          !allData.permissionTypeCheckStatus[item.id] ||
          !allData.permissionTypeCheckStatus[item.id].changed
        ) {
          allData.permissionTypeCheckStatus[item.id] = {
            checkAll: false,
            indeterminate: false,
            permissionTypes: item.permissionTypes || [],
            changed: false, // 当前这个记录的值是否发生了变化，后续通过这个字段来标识是否需要提交到后端
          };
          // 重新计算全选框状态
          refreshPermissionTypeCheckAllStatus(item.id);
        }
        // 将已授权的所有资源的权限类型添加到当前记录的权限类型中
        allData.permissionTypeCheckStatus[item.id].permissionTypes.push(
          ...(allData.authedAllResourcePermissionTypes || []),
        );
      });
      allData.isResetChange = false;
      return data;
    },
    formConfig: {},
    useSearchForm: false,
    showTableSetting: true,
    bordered: true,
    showIndexColumn: false,
    size: 'small',
    rowKey: 'id',
    pagination: {
      pageSizeOptions: ['10', '30', '50', '100', '200'],
    },
    immediate: false,
    canResize: true,
    scroll: {
      y: '400px',
    },
  });
  const [registerModal, { setDrawerProps, closeDrawer }] = useDrawerInner(async (data) => {
    try {
      cleanAllData();
      // 加载初始数据
      allData.modalSpinning = true;
      const { record } = data;
      allData.record = record;
      setDrawerProps({
        confirmLoading: false,
        destroyOnClose: true,
        width: 1200,
        title: '编辑【' + record.name + '】资源权限',
        showFooter: true,
      });
    } finally {
      allData.modalSpinning = false;
    }
  });

  async function initSearchFormSchema() {
    let { schemas } = await getResourceViewSearchFormSchemas({
      id: allData.currentResourceTypeId,
    });
    schemas = schemas || [];
    setProps({
      formConfig: {
        schemas: schemas,
        labelWidth: 120,
        // schemas: getResourceViewSearchFormSchema(),
        autoSubmitOnEnter: true,
        autoAdvancedLine: 3,
        showAdvancedButton: true,
        compact: true,
      },
      useSearchForm: schemas.length != 0,
    });
  }
  async function initTableColumns() {
    let { columns, showEnumFields, showTagFields } = await getResourceViewColumns({
      id: allData.currentResourceTypeId,
    });
    // 添加权限类型列
    columns = columns || [];
    columns.push({
      title: '权限类型',
      dataIndex: 'permissionTypes',
      width: 300,
      resizable: true,
      fixed: 'right',
    });
    // 设置表格列
    setColumns(columns);

    // 设置需要展示为枚举的字段
    let showEnumFieldsMap = {};
    Object.keys(showEnumFields)?.forEach((key) => {
      const dictCode = showEnumFields[key];
      showEnumFieldsMap[key] = getDictDataMapFromDict(dictCode);
    });
    allData.showEnumFields = showEnumFieldsMap;

    // 获取需要展示为Tag的字段
    let showTagFieldsMap = {};
    Object.keys(showTagFields)?.forEach((key) => {
      const dictCode = showTagFields[key];
      showTagFieldsMap[key] = getDictDataMapFromDict(dictCode);
    });
    allData.showTagFields = showTagFieldsMap;

    return;
  }

  // 资源类型 发生改变
  async function onResourceTypeChange() {
    const values = getFieldsValue();
    // TODO 加载资源相关基础信息
    const res = await getRoleResourceInfo({
      roleId: allData.record.id,
      resourceTypeId: values.resourceTypeId,
    });
    // 初始化这个角色对应的资源信息
    allData.allPermissionTypes = res.permissionTypes || [];
    allData.authedAllResourcePermissionTypes = res.authedAllResourcePermissionTypes || [];
    onValuesChange();
    // 需要重置授权部分资源权限修改
    allData.isResetChange = true;
    reRenderTable();
  }

  function onValuesChange() {
    allData.currentResourceTypeId = getFieldsValue().resourceTypeId || 0;
  }

  // 重新加载表格和搜索框
  async function reRenderTable() {
    try {
      setTableData([]);
      allData.modalSpinning = true;
      await initSearchFormSchema();
      await initTableColumns();
      reload();
    } finally {
      allData.modalSpinning = false;
    }
  }
  function getFuncs() {
    return {
      setFieldsValue,
      updateSchema,
    };
  }
  // 获取可选择的权限类型=所有的权限类型-已授权的所有资源的权限类型
  function getCanSelectPermissionTypes() {
    return allData.allPermissionTypes
      .filter((item) => !allData.authedAllResourcePermissionTypes.includes(item.code))
      .map((item) => item.code);
  }
  // 全选选项发生改变
  function onPermissionTypeCheckAllChange(id: number) {
    if (!allData.permissionTypeCheckStatus[id]) {
      return;
    }
    allData.permissionTypeCheckStatus[id].indeterminate = false;
    allData.permissionTypeCheckStatus[id].changed = true; //值发生了变化，需要提交到后端

    // 获取可选择的权限类型
    const allPermissionTypes = getCanSelectPermissionTypes();

    // 点击全选框不影响已授权的所有资源的权限类型
    // 这里不能直接赋值，因为allData.authedAllResourcePermissionTypes是引用类型，如果直接赋值，会导致allData.authedAllResourcePermissionTypes的值发生变化
    // 所以需要使用展开运算符来创建一个新的数组，否则会导致修改permissionTypes时allData.authedAllResourcePermissionTypes的值发生变化
    let selectedPermissionTypes = [...(allData.authedAllResourcePermissionTypes || [])];

    // 全选=授权所有的权限类型+可选择的权限类型
    if (allData.permissionTypeCheckStatus[id].checkAll === true) {
      allData.permissionTypeCheckStatus[id].permissionTypes = [
        ...selectedPermissionTypes,
        ...allPermissionTypes,
      ];
    } else {
      // 取消全选=授权所有的权限类型
      allData.permissionTypeCheckStatus[id].permissionTypes = [...selectedPermissionTypes];
    }
  }

  // 权限类型选项发生改变
  function onPermissionTypeCheckChange(id: number) {
    if (!allData.permissionTypeCheckStatus[id]) {
      return;
    }
    allData.permissionTypeCheckStatus[id].changed = true; //值发生了变化，需要提交到后端
    refreshPermissionTypeCheckAllStatus(id);
  }

  // 重新计算全选状态
  function refreshPermissionTypeCheckAllStatus(id: number) {
    // 获取可选择的权限类型
    const canSelectPermissionTypes = getCanSelectPermissionTypes();

    // 当前记录选中的权限类型，去除已授权的所有资源的权限类型
    const currentSelectedPermissionTypes = (
      allData.permissionTypeCheckStatus[id].permissionTypes || []
    ).filter((item) => !allData.authedAllResourcePermissionTypes.includes(item));

    // 未选中有效的记录
    if (
      !currentSelectedPermissionTypes || // 未选择任何记录
      currentSelectedPermissionTypes.length === 0 || // 未选择任何记录
      !currentSelectedPermissionTypes.some((item) => canSelectPermissionTypes.includes(item)) // 选择的记录中不存在所有权限类型中，所有的权限类型都不存在
    ) {
      allData.permissionTypeCheckStatus[id].checkAll = false;
      allData.permissionTypeCheckStatus[id].indeterminate = false;
      return;
    }

    let currentIncludeAllPermissionTypes = true;
    // 判断是否全部选中
    canSelectPermissionTypes.forEach((item) => {
      if (!currentSelectedPermissionTypes.includes(item)) {
        currentIncludeAllPermissionTypes = false;
      }
    });
    // 全部选中
    if (currentIncludeAllPermissionTypes) {
      allData.permissionTypeCheckStatus[id].checkAll = true;
      allData.permissionTypeCheckStatus[id].indeterminate = false;
    } else {
      // 部分选中
      allData.permissionTypeCheckStatus[id].checkAll = false;
      allData.permissionTypeCheckStatus[id].indeterminate = true;
    }
  }
  // 重置授权部分资源权限修改
  function handleResetChange() {
    Modal.confirm({
      title: '是否重置“授权部分资源”相关修改？',
      onOk: () => {
        allData.isResetChange = true;
        reload();
      },
    });
  }
  async function handleOk() {
    const headerValues = await validate();
    // 授权所有资源
    let requestData = {
      resourceTypeId: headerValues.resourceTypeId, // 资源类型
      authedAllResourcePermissionTypes: allData.authedAllResourcePermissionTypes, // 已授权的所有资源的权限类型
      changedPermissionTypes: [] as Recordable[], // 修改过的权限类型
    };

    // 授权部分资源
    // 获取修改过的权限类型
    Object.keys(allData.permissionTypeCheckStatus).forEach((key) => {
      const id = parseInt(key);
      if (allData.permissionTypeCheckStatus[id].changed) {
        requestData.changedPermissionTypes.push({
          resourceId: id,
          permissionTypes:
            allData.permissionTypeCheckStatus[id].permissionTypes?.filter(
              (item) => !allData.authedAllResourcePermissionTypes.includes(item),
            ) || [],
        });
      }
    });
    // 提交请求
    try {
      setDrawerProps({ confirmLoading: true });
      await updateRoleResource(allData.record.id, requestData);
      message.success('操作成功');
    } finally {
      setDrawerProps({ confirmLoading: false });
    }
    closeDrawer();
    emit('reload');
    onResourceTypeChange();
  }
  function handleSelectPage(status: boolean) {
    getDataSource().forEach((item) => {
      allData.permissionTypeCheckStatus[item.id].checkAll = status;
      onPermissionTypeCheckAllChange(item.id);
    });
  }
  function onAuthedAllResourcePermissionTypesChange() {
    allData.isResetChange = true;
    reload();
  }
  function handleSelectPermissionType(permissionCode: string) {
    // 如果已授权的所有资源的权限类型中包含这个权限类型，则不进行操作
    if (
      allData.authedAllResourcePermissionTypes?.length > 0 &&
      allData.authedAllResourcePermissionTypes?.includes(permissionCode)
    ) {
      return;
    }
    getDataSource().forEach((item) => {
      // 如果当前记录的权限类型中不包含这个权限类型，则添加这个权限类型
      if (allData.permissionTypeCheckStatus[item.id]) {
        if (!allData.permissionTypeCheckStatus[item.id].permissionTypes.includes(permissionCode)) {
          allData.permissionTypeCheckStatus[item.id].permissionTypes.push(permissionCode);
          allData.permissionTypeCheckStatus[item.id].changed = true;
        }
      } else {
        allData.permissionTypeCheckStatus[item.id] = {
          checkAll: false,
          indeterminate: false,
          permissionTypes: [permissionCode],
          changed: true,
        };
      }
      onPermissionTypeCheckChange(item.id);
    });
  }
</script>

<style scoped>
  :deep(.vben-basic-table.vben-basic-table-form-container) {
    margin-top: 5px;
    padding-top: 0;
  }

  :deep(.ant-form.ant-form-horizontal.ant-form-default.vben-basic-form) {
    margin-top: 0;
    margin-bottom: 0;
    padding-top: 5px;
    padding-bottom: 0;
  }

  .resource-checkbox-container {
    /* margin: 0 20px; */

    padding-top: 5px;

    /* border: 1px solid #f0f0f0; */

    /* border-radius: 5px; */

    /* background-color: #fff; */
  }

  :deep(.ant-form-item.css-dev-only-do-not-override-7oeufo) {
    margin-top: 10px;
    margin-bottom: 10px;
  }
</style>

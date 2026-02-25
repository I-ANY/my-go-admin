<template>
  <div>
    <BasicModal
      v-bind="$attrs"
      @register="registerModal"
      title="设备列表"
      width="1200px"
      :maxHeight="600"
      :destroyOnClose="true"
      @ok="handleConfirm"
      style="top: 10px"
    >
      <BasicTable
        @register="registerTable"
        style="margin-top: 0; padding: 0"
        body-style="padding: 1px"
      >
        <template #form-categoryId="{ model, field }">
          <Select
            v-model:value="model[field]"
            placeholder="选择业务大类"
            allowClear
            showSearch
            :filterOption="filterOption"
            @change="reloadBusinessOptions"
          >
            <Select.Option
              v-for="item in categoryOptions"
              :key="item.value"
              :value="item.value"
              :label="item.label"
            >
              <div style="display: flex; align-items: center; justify-content: space-between">
                <span>{{ item.label }}</span>
                <Tag
                  v-if="item.isVirtual == 1"
                  color="green"
                  style=" transform: scale(0.8);font-size: 12px"
                >
                  虚拟
                </Tag>
              </div>
            </Select.Option>
          </Select>
        </template>
      </BasicTable>
    </BasicModal>
  </div>
</template>
<script setup lang="ts">
  import { BasicModal, useModalInner } from '@/components/Modal';
  import { BasicTable, useTable } from '@/components/Table';
  import { deviceColumns, deviceSearchFormSchema } from '@/views/ops/batchExecute/data';
  import { GetAuthedBiz, GetAuthedServer } from '@/api/business/biz';
  import { splitByLineAndTrim } from '@/utils/util';
  import { ref } from 'vue';
  import { Tag, Select } from 'ant-design-vue';

  const emit = defineEmits(['confirmDevices', 'register']);

  interface Category {
    value: number;
    label: string;
    text: string;
    isVirtual?: number;
    subcategories: SubCategory[];
  }

  interface SubCategory {
    value: number;
    label: string;
  }

  const categoryOptions = ref<Category[]>([]);
  const subCategoryOptions = ref<SubCategory[]>([]);

  const [registerModal, { closeModal }] = useModalInner(async () => {
    // 获取数据
    await initBusinessOptions();
  });

  const [registerTable, { getForm, getSelectRows, setTableData, setSelectedRows }] = useTable({
    title: '选择设备',
    api: GetAuthedServer,
    columns: deviceColumns,
    formConfig: {
      labelWidth: 0,
      schemas: deviceSearchFormSchema(reloadBusinessOptions),
      showActionButtonGroup: true,
      compact: true,
      autoSubmitOnEnter: true,
      submitOnChange: true,
      showAdvancedButton: true,
      showResetButton: true,
      showSubmitButton: false,
      baseRowStyle: {
        marginTop: '0',
        paddingTop: '0',
      },
    },
    useSearchForm: true,
    showTableSetting: false,
    bordered: true,
    showIndexColumn: false,
    showSelectionBar: true, // 显示多选状态栏
    rowSelection: {
      type: 'checkbox',
    },
    immediate: false, // 设置 immediate 为 false，阻止初始化时自动加载数据
    size: 'small',
    scroll: {
      y: 320,
    },
    pagination: {
      pageSizeOptions: ['10', '30', '50', '100', '200', '500'],
    },
    beforeFetch: async (params) => {
      params.hostnames = splitByLineAndTrim(params.hostnames);
      params.categoryIds = [params.categoryId];
      params.excludeStatus = [6];
      delete params.categoryId;
      return params;
    },
  });

  async function initBusinessOptions() {
    categoryOptions.value = [];
    subCategoryOptions.value = [];
    await getForm().setFieldsValue({
      categoryId: [],
      businessIds: [],
      hostnames: '',
      // Add other fields as needed
    });
    setTableData([]);
    const res = await GetAuthedBiz({});
    if (res.categories) {
      res.categories.forEach((item) => {
        let children: any[] = [];
        if (item.subcategories) {
          item.subcategories.forEach((subItem) => {
            children.push({
              value: subItem.id,
              label: subItem.name,
            });
          });
        }
        if (item.virtualSubcategories) {
          item.virtualSubcategories.forEach((subItem) => {
            children.push({
              value: subItem.id,
              label: subItem.name,
            });
          });
        }
        subCategoryOptions.value.push(...children);
        categoryOptions.value.push({
          value: item.id,
          label: item.name,
          text: item.name,
          isVirtual: item.isVirtual,
          subcategories: children,
        });
      });

      await getForm().updateSchema([
        {
          field: 'businessIds',
          componentProps: {
            options: subCategoryOptions.value,
          },
        },
      ]);
    }
  }

  async function reloadBusinessOptions() {
    const values = getForm().getFieldsValue();
    let selectCategoryId = values.categoryId;
    // 当清除选择时（selectCategoryId 为空），重置 businessIds 选项
    if (!selectCategoryId) {
      await getForm().updateSchema([
        {
          field: 'businessIds',
          componentProps: {
            options: subCategoryOptions.value, // 恢复所有子分类选项
          },
        },
      ]);
      return;
    }

    // 获取已选择的业务，过滤掉Options中不存在的选项
    categoryOptions.value.forEach((option) => {
      if (option.value === selectCategoryId) {
        if (option.subcategories) {
          getForm().updateSchema([
            {
              field: 'businessIds',
              componentProps: {
                options: option.subcategories,
              },
            },
          ]);
        }
      }
    });
  }

  const filterOption = (input: string, option: any) => {
    return option.label.toLowerCase().indexOf(input.toLowerCase()) >= 0;
  };

  async function handleConfirm() {
    const selectedRows = getSelectRows();
    if (selectedRows && selectedRows.length > 0) {
      // 将选中的设备数据传递给父组件
      emit('confirmDevices', selectedRows);
      await setSelectedRows([]);
      closeModal();
    }
  }
</script>

<style scoped lang="less">
  .vben-basic-table-form-container .ant-form {
    width: 100%;
    margin-bottom: 2px;
    padding: 5px 10px 6px;
    border-radius: 2px;
  }
</style>

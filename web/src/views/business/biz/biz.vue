<template>
  <div class="h-full flex p-2 gap-x-0 -space-x-4">
    <div class="flex w-1/3">
      <BasicTable @register="registerSubcategoryTable">
        <template #bodyCell="{ record, column }">
          <template v-if="column.key == 'status'">
            <Switch
              :checked="record.status === 1"
              :checked-children="'启用'"
              :un-checked-children="'禁用'"
              :class="record.status === 1 ? 'custom-switch-on' : 'custom-switch-off'"
              :loading="subcategoryTag[record.id] || false"
              @change="(checked) => handleSubcategoryStatus(record, checked)"
            />
          </template>
        </template>
      </BasicTable>
    </div>
    <div class="flex-1 flex flex-col h-full">
      <div class="h-1/3 mb-4">
        <div class="flex-1 h-0">
          <BasicTable @register="registerCategoryTable">
            <template #toolbar>
              <!--        <Excel :excelType="excelType" @success-upload="uploadSuccess" />-->
              <a-button type="primary" @click="handleCategoryCreate"> 新增</a-button>
            </template>
            <template #bodyCell="{ column, record }">
              <template v-if="column.key == 'status'">
                <Switch
                  :checked="record.status === 1"
                  :checked-children="'启用'"
                  :un-checked-children="'禁用'"
                  :class="record.status === 1 ? 'custom-switch-on' : 'custom-switch-off'"
                  :loading="categoryTag[record.id] || false"
                  @change="(checked) => handleStatus(record, checked)"
                />
              </template>
              <template v-if="column.key == 'subcategories'">
                <template v-if="record.subcategories && record.subcategories.length > 0">
                  <div class="subcategory-tags">
                    <Tag
                      v-for="(item, index) in record.subcategories"
                      :key="index"
                      :color="item.status === 1 ? '#108ee9' : '#FF1100'"
                      style="margin-right: 4px"
                    >
                      <Tooltip :title="item.status === 1 ? '启用' : '禁用'">
                        <span>{{ item.name }}</span>
                      </Tooltip>
                    </Tag>
                  </div>
                </template>
                <template v-else>
                  <span>-</span>
                </template>
              </template>
              <template v-if="column.key == 'action'">
                <TableAction
                  :actions="[
                    {
                      icon: 'clarity:note-edit-line',
                      onClick: (event) => {
                        event.stopPropagation();
                        handleCategoryEdit(record);
                      },
                      tooltip: '编辑',
                    },
                  ]"
                />
              </template>
            </template>
          </BasicTable>
          <CategoryDrawer @register="registerCategoryDrawer" @success="handleSuccess" />
        </div>
      </div>
    </div>
  </div>
</template>
<script setup lang="ts">
  import { BasicTable, TableAction, useTable } from '@/components/Table';
  import {
    GetCategoryListAll,
    GetSubcategoryListAll,
    UpdateCategory,
    UpdateSubcategory,
  } from '@/api/business/biz';
  import {
    categoryColumns,
    categorySchemas,
    subcategoryColumns,
    subcategorySchemas,
  } from '@/views/business/biz/data';
  // import { splitByLineAndTrim } from '@/utils/util';
  // import { CopyOutlined } from '@ant-design/icons-vue';
  import { message, Switch, Tag, Tooltip } from 'ant-design-vue';
  import { useDrawer } from '@/components/Drawer';
  import CategoryDrawer from './categoryDrawer.vue';
  import { reactive, unref } from 'vue';

  const categoryTag = reactive({});
  const subcategoryTag = reactive({});

  const [registerCategoryTable, { reload: reloadCategoryTable }] = useTable({
    title: '业务大类',
    api: GetCategoryListAll,
    columns: categoryColumns,
    formConfig: {
      labelWidth: 120,
      schemas: categorySchemas,
      autoSubmitOnEnter: true,
      showAdvancedButton: true,
      autoAdvancedLine: 4,
      // alwaysShowLines: 2,
    },
    useSearchForm: true,
    showTableSetting: true,
    bordered: true,
    showIndexColumn: false,
    pagination: {
      // pageSizeOptions: ['10', '30', '50', '100', '500', '1000'],
    },
    actionColumn: {
      width: 50,
      title: '操作',
      dataIndex: 'action',
      // slots: { customRender: 'action' },
      fixed: 'right',
    },
  });

  const [registerSubcategoryTable, { reload: reloadSubcategoryTable }] = useTable({
    title: '子业务列表',
    api: GetSubcategoryListAll,
    columns: subcategoryColumns,
    formConfig: {
      labelWidth: 120,
      schemas: subcategorySchemas,
      autoSubmitOnEnter: true,
      showAdvancedButton: true,
      autoAdvancedLine: 4,
      // alwaysShowLines: 2,
    },
    useSearchForm: true,
    showTableSetting: true,
    bordered: true,
    showIndexColumn: false,
    pagination: {
      // pageSizeOptions: ['10', '30', '50', '100', '500', '1000'],
    },
  });

  // 侧边框
  const [registerCategoryDrawer, { openDrawer: openCategoryDrawer }] = useDrawer();

  // 新增业务大类
  function handleCategoryCreate() {
    openCategoryDrawer(true, {
      record: {},
      isUpdate: false,
    });
  }

  // 业务大类启用/禁用
  async function handleStatus(record: Recordable, checked) {
    categoryTag[record.id] = true;
    try {
      let data = { ...record };
      data.status = 0;
      if (checked) {
        data.status = 1;
      }
      if (record.subcategories) {
        const ids = record.subcategories.map((item) => item.id);
        data.subcategories = ids;
      }
      if (record.virtualSubcategories) {
        const ids = record.virtualSubcategories.map((item) => item.id);
        data.virtualSubcategories = ids;
      }
      await UpdateCategory(record.id, unref(data));
      message.success('操作成功');
      await reloadCategoryTable();
    } finally {
      categoryTag[record.id] = false;
    }
  }

  // 业务小类启用/禁用
  async function handleSubcategoryStatus(record: Recordable, checked) {
    subcategoryTag[record.id] = true;
    try {
      let data = { ...record };
      data.status = 0;
      if (checked) {
        data.status = 1;
      }
      delete data.subcategories;
      await UpdateSubcategory(unref(data));
      message.success('操作成功');
      await reloadSubcategoryTable();
      await reloadCategoryTable();
    } finally {
      subcategoryTag[record.id] = false;
    }
  }

  // 编辑业务大类
  function handleCategoryEdit(record: Recordable) {
    openCategoryDrawer(true, {
      record: record,
      isUpdate: true,
    });
  }

  async function handleSuccess() {
    message.success('操作成功');
    await reloadCategoryTable();
  }
</script>

<style scoped lang="less">
  .custom-switch-on .ant-switch-checked {
    background-color: #108ee9 !important;
  }

  .custom-switch-off .ant-switch {
    background-color: #f10 !important;
  }

  .subcategory-tags {
    display: flex;
    flex-wrap: wrap;
    gap: 4px;
  }

  .subcategory-tag {
    margin-right: 0 !important; // 覆盖默认的margin-right
  }
</style>

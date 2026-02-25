<template>
  <div>
    <BasicModal
      title="资源数据"
      @register="registerModal"
      :width="1600"
      @cancel="
        () => {
          emit('reload');
        }
      "
    >
      <Spin :spinning="allData.modalSpinning">
        <BasicTable @register="registerTable">
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
                  :color="allData.showTagFields[columnName][record[columnName]].color || 'default'"
                  >{{ allData.showTagFields[columnName][record[columnName]].dictLabel }}</Tag
                >
                <span v-else>{{ record[columnName] }}</span>
              </template>
            </template>
          </template>
        </BasicTable>
      </Spin>
    </BasicModal>
  </div>
</template>

<script lang="ts" setup>
  import { BasicModal, useModalInner } from '@/components/Modal';
  import { BasicTable, useTable } from '@/components/Table';
  import {
    getResourceDetailList,
    getResourceViewColumns,
    getResourceViewSearchFormSchemas,
  } from '@/api/sys/resource';
  import { reactive } from 'vue';
  import { Spin, Tag } from 'ant-design-vue';
  import { getDictDataMapFromDict } from '@/utils/dict';

  defineOptions({ name: 'ResourceViewModal' });
  const emit = defineEmits(['reload', 'register']);
  const allData = reactive({
    record: {} as Recordable,
    showEnumFields: {},
    showTagFields: {},
    modalSpinning: true,
  });
  const [registerTable, { reload, setColumns, setProps }] = useTable({
    title: '资源数据',
    api: getResourceDetailList,
    // columns: getResourceViewColumns(),
    beforeFetch: (params) => {
      params.id = allData.record.id;
      return params;
    },
    formConfig: {},
    useSearchForm: true,
    showTableSetting: true,
    bordered: true,
    showIndexColumn: false,
    size: 'small',
    pagination: {
      pageSizeOptions: ['10', '30', '50', '100', '200'],
    },
    immediate: false,
    canResize: true,
    scroll: {
      y: '350px',
    },
  });
  const [registerModal, { setModalProps }] = useModalInner(async (data) => {
    resetAllData();
    try {
      allData.modalSpinning = true;
      const { record } = data;
      allData.record = record;
      setModalProps({
        confirmLoading: false,
        destroyOnClose: true,
        width: 1200,
        showCancelBtn: false,
        showOkBtn: false,
        title: '资源数据：' + record.name,
      });
      initSearchFormSchema();
      await initTableColumns();
      reload();
    } finally {
      allData.modalSpinning = false;
    }
  });

  async function initSearchFormSchema() {
    let { schemas } = await getResourceViewSearchFormSchemas({
      id: allData.record.id,
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
    const { columns, showEnumFields, showTagFields } = await getResourceViewColumns({
      id: allData.record.id,
    });
    // 设置表格列
    setColumns(columns || []);

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

  function resetAllData() {
    allData.record = {};
    allData.showEnumFields = {};
    allData.showTagFields = {};
    allData.modalSpinning = true;
  }
</script>

<style scoped>
  :deep(.vben-basic-table.vben-basic-table-form-container) {
    margin-top: 5px;
    padding-top: 0;
  }

  :deep(.ant-form.ant-form-horizontal.ant-form-default.vben-basic-form) {
    margin-top: 5px;
    margin-bottom: 5px;
    padding-top: 5px;
    padding-bottom: 5px;
  }
</style>

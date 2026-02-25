<template>
  <BasicDrawer
    v-bind="$attrs"
    @register="registerDrawer"
    @ok="handleSubmit"
    destroyOnClose
    :width="1200"
    :maskClosable="true"
    :showFooter="true"
    :showCancelBtn="true"
    :showOkBtn="true"
    @cancel="
      () => {
        emit('reload');
      }
    "
  >
    <BasicForm @register="registerForm">
      <template #pt="{ field }">
        <a-button v-if="field === 'pt_1'" @click="() => addPermissionCode(field)">+</a-button>
        <a-button v-else @click="() => delPermissionCode(field)">-</a-button>
      </template>
      <template #rf="{ field }">
        <a-button v-if="field === 'rf_1'" @click="() => addResourceField(field)">+</a-button>
        <a-button v-else @click="() => delResourceField(field)">-</a-button>
      </template>
      <template #rf_empty="{}"> </template> </BasicForm
  ></BasicDrawer>
</template>
<script lang="ts" setup>
  import { BasicDrawer, useDrawerInner } from '@/components/Drawer';
  import { BasicForm, FormSchema, useForm } from '@/components/Form/index';
  import { nextTick, reactive } from 'vue';
  import { getResourceFormSchema } from './data';
  import { getResourceTableField, addResource, updateResource } from '@/api/sys/resource';
  import { getCheckVariableFn } from '@/utils/util';

  defineOptions({
    name: 'ResourceDrawer',
  });
  const emit = defineEmits(['success', 'reload', 'register']);
  const allData = reactive({
    isUpdate: false,
    record: {} as Recordable,
    permissionTypeNumber: 0,
    resourceFieldNumber: 0,
    tableFieldOPtions: [] as any[],
  });

  const [
    registerForm,
    {
      validate,
      setFieldsValue,
      updateSchema,
      getFieldsValue,
      removeSchemaByField,
      appendSchemaByField,
    },
  ] = useForm({
    labelWidth: 100,
    schemas: getResourceFormSchema(onTabelChage),
    actionColOptions: {
      span: 24,
    },
    showActionButtonGroup: false,
  });
  const [registerDrawer, { setDrawerProps, closeDrawer }] = useDrawerInner((data) => {
    const { isUpdate, record } = data;
    allData.isUpdate = isUpdate;
    allData.record = record;
    allData.permissionTypeNumber = 0;
    allData.resourceFieldNumber = 0;
    if (isUpdate === true) {
      setDrawerProps({
        title: `编辑资源：${record.name}`,
      });
      updateEditSchema(record);
    } else {
      setDrawerProps({
        title: '新增资源',
      });
      updateAddSchema();
    }
  });

  // 选择的表发生变化
  async function onTabelChage() {
    let options: any[] = [];
    const values = getFieldsValue();
    const { fields } = await getResourceTableField({ tableName: values.table });
    fields?.forEach((field) => {
      options.push({
        label: field,
        value: field,
      });
    });
    allData.tableFieldOPtions = options;

    // 修改表格字段的选项
    let schemas: any[] = [];
    let fieldNameValues: Recordable = {};
    for (let i = 1; i <= allData.resourceFieldNumber; i++) {
      const fieldName = `rf_${i}_field_name`;
      schemas.push({
        field: fieldName,
        componentProps: {
          options: options,
        },
      });

      // 如果Options里面的选项没有fieldName则置空
      if (!options.some((option) => option.value === values[fieldName])) {
        fieldNameValues[fieldName] = null;
      }
    }
    if (schemas.length > 0) {
      updateSchema(schemas);
      setFieldsValue(fieldNameValues);
    }
  }
  async function updateAddSchema() {
    appendPermissionTypeSchema(1);
    appendresourceFieldSchema(1);
  }
  async function updateEditSchema(record: Recordable) {
    const { permissionTypes, fields } = record;
    // 回显权限类型表单
    let permissionTypeValues = {};
    // let permissionTypeSchema: any[] = [];
    if (permissionTypes && permissionTypes.length > 0) {
      let permissionTypeNumber = allData.permissionTypeNumber;
      // 先追加指定数量的 权限类型 表单
      await appendPermissionTypeSchema(permissionTypes.length);
      // 初始化前面添加的表单的值
      for (let i = 0; i < permissionTypes.length; i++) {
        permissionTypeNumber++;
        permissionTypeValues[`pt_${permissionTypeNumber}_name`] = permissionTypes[i].name;
        permissionTypeValues[`pt_${permissionTypeNumber}_code`] = permissionTypes[i].code;
        // code不支持修改
        // permissionTypeSchema.push({
        //   field: `pt_${allData.permissionTypeNumber}_code`,
        //   componentProps: {
        //     disabled: true,
        //     readonly: true,
        //   },
        // });
      }
    } else {
      // 没有就添加一个
      await appendPermissionTypeSchema(1);
    }

    // 回显字段数据
    let resourceFieldValues = {};
    if (fields && fields.length > 0) {
      let resourceFieldNumber = allData.resourceFieldNumber;
      // 先追加指定数量的 资源字段 表单
      await appendresourceFieldSchema(fields.length);
      // 初始化前面添加的表单的值
      for (let i = 0; i < fields.length; i++) {
        resourceFieldNumber++;
        resourceFieldValues[`rf_${resourceFieldNumber}_field_name`] = fields[i].fieldName;
        resourceFieldValues[`rf_${resourceFieldNumber}_column_name`] = fields[i].columnName;
        resourceFieldValues[`rf_${resourceFieldNumber}_support_filter`] =
          fields[i].supportFilter || 0;
        resourceFieldValues[`rf_${resourceFieldNumber}_is_dict`] = fields[i].isDict || 0;
        resourceFieldValues[`rf_${resourceFieldNumber}_show_with_tag`] = fields[i].showWithTag || 0;
        resourceFieldValues[`rf_${resourceFieldNumber}_dict_key`] = fields[i].dictKey;
        resourceFieldValues[`rf_${resourceFieldNumber}_sort`] = fields[i].sort;
      }
    } else {
      // 没有就添加一个
      await appendresourceFieldSchema(1);
    }

    // 更新表单结构
    await updateSchema([
      {
        field: 'identify',
        componentProps: {
          disabled: true,
          readonly: true,
        },
      },
      // ...permissionTypeSchema,
    ]);
    // 更新表单值
    await setFieldsValue({
      name: record.name || null,
      identify: record.identify || null,
      // table: record.table || null,// 这里先不更新表名，因为更新表名会导致“表字段”丢失，因为更新表名时“表字段”还没有值所以会重置“表字段”值
      sort: record.sort || 10,
      filter: record.filter || null,
      ...permissionTypeValues,
      ...resourceFieldValues,
    });
    // 延迟更新表名
    nextTick(() => {
      setFieldsValue({
        table: record.table || null,
      });
    });
  }

  function addPermissionCode(_field: string) {
    appendPermissionTypeSchema(1);
  }
  async function delPermissionCode(field: string) {
    await removeSchemaByField([`${field}_name`, `${field}_code`, field]);
    allData.permissionTypeNumber--;
  }
  // 添加权限类型表单
  async function appendPermissionTypeSchema(num: number) {
    let prefixField = 'permissionTypeDivider';
    if (allData.permissionTypeNumber > 0) {
      prefixField = `pt_${allData.permissionTypeNumber}`;
    }
    let schemas: FormSchema[] = [];
    for (let i = 0; i < num; i++) {
      allData.permissionTypeNumber++;
      let permissionTypeNumber = allData.permissionTypeNumber;
      const schema = getPermissionTypeSchema(permissionTypeNumber);
      schemas.push(...schema);
    }
    await appendSchemaByField(schemas, prefixField);
  }
  function getPermissionTypeSchema(permissionTypeNumber: number): FormSchema[] {
    return [
      {
        field: `pt_${permissionTypeNumber}_name`,
        label: '名称',
        component: 'Input',
        required: true,
        defaultValue: null,
        componentProps: {
          placeholder: '请输入名称',
          allowClear: false,
        },
        rules: [{ min: 2, trigger: 'change' }],
        colProps: {
          span: 11,
        },
      },
      {
        field: `pt_${permissionTypeNumber}_code`,
        label: '编码',
        component: 'Input',
        required: true,
        defaultValue: null,
        componentProps: {
          placeholder: '请输入编码',
          allowClear: false,
        },
        rules: [{ validator: getCheckVariableFn(4), trigger: 'change' }],
        colProps: {
          span: 11,
        },
      },
      {
        field: `pt_${permissionTypeNumber}`,
        // component: 'Input',
        label: ' ',
        slot: 'pt',
        colProps: {
          span: 1,
        },
      },
    ];
  }
  async function appendresourceFieldSchema(num: number) {
    let prefixField = 'resourceField';
    if (allData.resourceFieldNumber > 0) {
      prefixField = `rf_${allData.resourceFieldNumber}_empty`;
    }
    let schemas: FormSchema[] = [];
    for (let i = 0; i < num; i++) {
      allData.resourceFieldNumber++;
      let resourceFieldNumber = allData.resourceFieldNumber;
      const schema = getResourceFieldSchema(resourceFieldNumber);
      schemas.push(...schema);
    }
    await appendSchemaByField(schemas, prefixField);
  }
  function getResourceFieldSchema(resourceFieldNumber: number): FormSchema[] {
    return [
      {
        field: `rf_${resourceFieldNumber}_field_name`,
        component: 'Select',
        label: '表字段',
        colProps: {
          span: 6,
        },
        componentProps: {
          options: allData.tableFieldOPtions,
          showSearch: true,
        },
        required: true,
      },
      {
        field: `rf_${resourceFieldNumber}_column_name`,
        component: 'Input',
        label: '列名',
        colProps: {
          span: 6,
        },
        required: true,
        rules: [{ min: 2, trigger: 'change' }],
      },
      {
        field: `rf_${resourceFieldNumber}_support_filter`,
        component: 'RadioButtonGroup',
        label: '支持筛选',
        componentProps: {
          options: [
            {
              label: '否',
              value: 0,
            },
            {
              label: '是',
              value: 1,
            },
          ],
        },
        defaultValue: 0,
        colProps: {
          span: 5,
        },
        required: true,
      },
      {
        field: `rf_${resourceFieldNumber}_sort`,
        label: '字段排序',
        component: 'InputNumber',
        helpMessage: '字段展示的顺序，越小越靠前',
        defaultValue: 1,
        componentProps: {
          min: -999999,
          max: 999999,
          precision: 0,
          step: 1,
        },
        colProps: {
          span: 5,
        },
        required: true,
        show: true,
      },
      {
        field: `rf_${resourceFieldNumber}`,
        label: ' ',
        slot: 'rf',
        colProps: {
          span: 1,
        },
      },
      {
        field: `rf_${resourceFieldNumber}_is_dict`,
        component: 'RadioButtonGroup',
        label: '为字典值',
        componentProps: {
          options: [
            {
              label: '否',
              value: 0,
            },
            {
              label: '是',
              value: 1,
            },
          ],
          onChange: (value: any) => {
            // 为字典值 修改时需要修数据
            if (value === 1) {
              updateSchema([
                {
                  field: `rf_${resourceFieldNumber}_dict_key`,
                  show: true,
                  required: true,
                },
                {
                  field: `rf_${resourceFieldNumber}_show_with_tag`,
                  show: true,
                  required: true,
                },
                {
                  field: `rf_${resourceFieldNumber}_empty`,
                  show: true,
                  colProps: {
                    span: 7,
                  },
                },
              ]);
            } else {
              updateSchema([
                {
                  field: `rf_${resourceFieldNumber}_dict_key`,
                  show: false,
                  required: false,
                },
                {
                  field: `rf_${resourceFieldNumber}_show_with_tag`,
                  show: false,
                  required: false,
                },
                {
                  field: `rf_${resourceFieldNumber}_empty`,
                  show: true,
                  colProps: {
                    span: 18,
                  },
                },
              ]);
            }
            // let v: Recordable = {};
            // v[`rf_${resourceFieldNumber}_show_with_tag`] = 0;
            // v[`rf_${allData.resourceFieldNumber}_dict_key`] = null;
            // setFieldsValue(v);
          },
        },
        defaultValue: 0,
        colProps: {
          span: 6,
        },
        required: true,
      },
      {
        field: `rf_${resourceFieldNumber}_show_with_tag`,
        component: 'RadioButtonGroup',
        label: '展示为Tag',
        componentProps: {
          options: [
            {
              label: '否',
              value: 0,
            },
            {
              label: '是',
              value: 1,
            },
          ],
        },
        defaultValue: 0,
        colProps: {
          span: 6,
        },
        required: false,
        show: false,
      },
      {
        field: `rf_${resourceFieldNumber}_dict_key`,
        label: '字典编码',
        component: 'Input',
        defaultValue: null,
        colProps: {
          span: 5,
        },
        required: false,
        show: false,
      },
      {
        field: `rf_${resourceFieldNumber}_empty`,
        label: ' ',
        slot: 'rf_empty',
        show: true,
        colProps: {
          span: 18,
        },
      },
    ];
  }
  function addResourceField(_field: string) {
    appendresourceFieldSchema(1);
  }
  async function delResourceField(field: string) {
    await removeSchemaByField([
      `${field}_field_name`,
      `${field}_column_name`,
      `${field}_support_filter`,
      `${field}_sort`,
      `${field}_is_dict`,
      `${field}_show_with_tag`,
      `${field}_dict_key`,
      `${field}_empty`,
      field,
    ]);
    allData.resourceFieldNumber--;
  }

  // 将表单数据转换成请求数据
  function convert2RequestData(values: Recordable): Recordable {
    let requestData: Recordable = {};

    // 基础字段
    const baseKeys = ['name', 'identify', 'table', 'filter', 'sort'];
    baseKeys.forEach((key) => {
      requestData[key] = values[key] === undefined ? null : values[key];
    });
    let permissionTypes: Recordable[] = [];
    let fields: Recordable[] = [];
    // 遍历表单内所有的字段，分别获取出pt_和rf_字段的值，组成数组
    const keys = Object.keys(values);
    keys?.forEach((key) => {
      if (!key) {
        return;
      }
      // permissionTypes
      if (/^pt_\d+$/.test(key)) {
        let permissionType = {};
        const match = key.match(/\d+/);
        if (!match) return;
        const showSort = parseInt(match[0], 10);
        permissionType['showSort'] = showSort;
        permissionType['code'] = values[key + '_code'] === undefined ? null : values[key + '_code'];
        permissionType['name'] = values[key + '_name'] === undefined ? null : values[key + '_name'];
        permissionTypes.push(permissionType);
      }
      // 资源字段
      if (/^rf_\d+$/.test(key)) {
        let field = {};
        const match = key.match(/\d+/);
        if (!match) return;
        const showSort = parseInt(match[0], 10);
        field['showSort'] = showSort;
        field['fieldName'] =
          values[key + '_field_name'] === undefined ? null : values[key + '_field_name'];
        field['columnName'] =
          values[key + '_column_name'] === undefined ? null : values[key + '_column_name'];
        field['supportFilter'] =
          values[key + '_support_filter'] === undefined ? null : values[key + '_support_filter'];
        field['isDict'] = values[key + '_is_dict'] === undefined ? null : values[key + '_is_dict'];
        field['showWithTag'] =
          values[key + '_show_with_tag'] === undefined ? null : values[key + '_show_with_tag'];
        field['dictKey'] =
          values[key + '_dict_key'] === undefined ? null : values[key + '_dict_key'];
        field['sort'] = values[key + '_sort'] === undefined ? null : values[key + '_sort'];
        fields.push(field);
      }
    });

    // permissionTypes 排序
    permissionTypes.sort((a, b) => {
      return a.showSort - b.showSort;
    });
    // fields 排序
    fields.sort((a, b) => {
      return a.showSort - b.showSort;
    });

    requestData['permissionTypes'] = permissionTypes;
    requestData['fields'] = fields;

    return requestData;
  }
  async function handleSubmit() {
    try {
      const values = await validate();
      setDrawerProps({ confirmLoading: true });
      const data = convert2RequestData(values);
      if (allData.isUpdate === true) {
        await updateResource(allData.record.id, data);
      } else {
        await addResource(data);
      }
      emit('success');
      closeDrawer();
    } finally {
      setDrawerProps({ confirmLoading: false });
    }
  }
</script>
<style lang="less" scoped></style>

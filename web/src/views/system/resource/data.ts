import { getResourceTableList } from '@/api/sys/resource';
import { BasicColumn, FormSchema } from '@/components/Table';
import { getCheckVariableFn } from '@/utils/util';

export const getResourceColumns = function (): BasicColumn[] {
  return [
    {
      title: 'id',
      dataIndex: 'id',
      width: 80,
      // resizable: true,
    },
    {
      title: '资源名称',
      dataIndex: 'name',
      width: 150,
      resizable: true,
    },
    {
      title: '资源编码',
      dataIndex: 'identify',
      width: 150,
      resizable: true,
    },
    {
      title: '排序',
      dataIndex: 'sort',
      width: 80,
      resizable: true,
    },
    {
      title: '权限类型',
      dataIndex: 'permissionTypes',
      width: 260,
      resizable: true,
    },
    {
      title: '表名',
      dataIndex: 'table',
      width: 160,
      resizable: true,
    },
    {
      title: '过滤条件',
      dataIndex: 'filter',
      width: 200,
      resizable: true,
    },
    {
      title: '修改时间',
      dataIndex: 'updatedAt',
      width: 150,
      resizable: true,
    },
  ];
};

export const getSearchResourceColumnsFormSchema = function (): FormSchema[] {
  return [
    {
      label: '资源名称',
      field: 'name',
      component: 'Input',
      colProps: { span: 6 },
    },
    {
      label: '资源编码',
      field: 'identify',
      component: 'Input',
      colProps: { span: 6 },
    },
    {
      label: '表名',
      field: 'table',
      component: 'Input',
      colProps: { span: 6 },
    },
  ];
};

export const getResourceFormSchema = function (onTabelChage): FormSchema[] {
  return [
    {
      label: '资源名称',
      field: 'name',
      component: 'Input',
      colProps: { span: 8 },
      required: true,
      rules: [{ min: 2, trigger: 'change' }],
    },
    {
      label: '资源编码',
      field: 'identify',
      component: 'Input',
      colProps: { span: 8 },
      required: true,
      helpMessage: '资源的唯一标识符，不能重复',
      rules: [{ validator: getCheckVariableFn(5), trigger: 'change' }],
    },
    {
      label: '排序',
      field: 'sort',
      component: 'InputNumber',
      colProps: { span: 8 },
      required: true,
      helpMessage: '资源的排序，越小越靠前展示',
      componentProps: {
        step: 1,
        // defaultValue: 10,
      },
      defaultValue: 10,
    },
    {
      label: '表名',
      field: 'table',
      component: 'ApiSelect',
      required: true,
      colProps: { span: 8 },
      componentProps: {
        api: getResourceTableList,
        params: {
          pageIndex: 1,
          pageSize: 9999,
        },
        immediate: true,
        showSearch: true,
        resultField: 'items',
        labelField: 'name',
        valueField: 'name',
        onChange: onTabelChage,
        allowClear: false,
      },
      helpMessage: '资源对应的数据库表名称',
    },
    {
      label: '过滤条件',
      field: 'filter',
      component: 'Input',
      colProps: { span: 8 },
      helpMessage: '输入SQL过滤条件，输入后仅展示过滤后的数据',
      required: false,
    },
    {
      label: '资源权限类型',
      field: 'permissionTypeDivider',
      component: 'Divider',
      componentProps: {
        orientation: 'center',
      },
      colProps: {
        span: 24,
      },
    },
    {
      label: '资源字段',
      field: 'resourceField',
      component: 'Divider',
      componentProps: {
        orientation: 'center',
      },
      colProps: {
        span: 24,
      },
    },
  ];
};

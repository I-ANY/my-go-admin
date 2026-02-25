<template>
  <BasicDrawer
    v-bind="$attrs"
    @register="registerDrawer"
    showFooter
    :title="getTitle"
    width="500px"
    @ok="handleSubmit"
    destroyOnClose
  >
    <BasicForm @register="registerForm">
      <template #businessPermission="{ model, field }">
        <BasicTree
          v-model:value="model[field]"
          :treeData="allData.treeData"
          :fieldNames="{ title: 'name', key: 'key', children: 'subcategories' }"
          :checkStrictly="false"
          checkable
          :toolbar="false"
          title="业务权限配置"
          :disabled="allData.disableSelect"
        />
      </template>
    </BasicForm>
  </BasicDrawer>
</template>
<script lang="ts" setup>
  import { computed, reactive } from 'vue';
  import { BasicForm, useForm } from '@/components/Form';
  import {
    AuthAllResource,
    BusinessPermissionCode,
    getBusinessPermissionFormSchema,
    ResourceIdentify,
  } from './role.data';
  import { BasicDrawer, useDrawerInner } from '@/components/Drawer';
  import { BasicTree, TreeItem } from '@/components/Tree';
  import { getRoleAuthedResource, getBusinessResource, roleResourceAuth } from '@/api/sys/resource';

  defineOptions({ name: 'BusinessPermissionDrawer' });

  const emit = defineEmits(['success', 'register']);
  const allData = reactive({
    treeData: [] as TreeItem[],
    recordId: 0,
    disableSelect: false,
    businessIds: [] as string[],
    authAllBusiness: AuthAllResource.NO,
    allBusinessIds: [] as string[],
  });

  const [registerForm, { setFieldsValue, validate }] = useForm({
    labelWidth: 100,
    baseColProps: { span: 24 },
    schemas: getBusinessPermissionFormSchema(onAuthAllBusinessChange),
    showActionButtonGroup: false,
  });

  const [registerDrawer, { setDrawerProps, closeDrawer }] = useDrawerInner(async (data) => {
    const { record } = data;
    allData.recordId = record.id;
    const { resourceIds, authAllResource } = await getRoleAuthedResource({
      roleId: record.id,
      identify: ResourceIdentify.SUBCATEGORY,
      permissionCode: BusinessPermissionCode.GENERAL_PERMISSION,
    });
    // 赋值权限变量和表单值
    allData.businessIds = resourceIds?.map((item: any) => `subcategory_${item}`) || [];
    allData.authAllBusiness =
      authAllResource === AuthAllResource.YES ? AuthAllResource.YES : AuthAllResource.NO;
    if (allData.authAllBusiness === AuthAllResource.YES) {
      setFieldsValue({
        authAllBusiness: AuthAllResource.YES,
        businessIds: allData.businessIds,
      });
    } else {
      setFieldsValue({
        authAllBusiness: AuthAllResource.NO,
        businessIds: null,
      });
    }
  });

  const getTitle = computed(() => '业务权限配置');

  async function handleSubmit() {
    try {
      const body = {} as any;
      const values = await validate();
      if (values.authAllBusiness === AuthAllResource.YES) {
        body.authAllResource = AuthAllResource.YES;
        body.identify = ResourceIdentify.SUBCATEGORY;
        body.permissionCode = BusinessPermissionCode.GENERAL_PERMISSION;
      } else {
        body.authAllResource = AuthAllResource.NO;
        body.identify = ResourceIdentify.SUBCATEGORY;
        body.permissionCode = BusinessPermissionCode.GENERAL_PERMISSION;
        body.resourceIds =
          values.businessIds
            ?.filter((item: any) => item.startsWith('subcategory_'))
            .map((item: any) => item.split('_')[1]) || [];
      }
      await roleResourceAuth(allData.recordId, body);
      emit('success', true);
      closeDrawer();
    } finally {
      setDrawerProps({ confirmLoading: false });
    }
  }
  // 授权所有业务 发生改变，授权所有业务时，禁用选择业务
  // 授权部分业务时，启用选择业务
  async function onAuthAllBusinessChange(value: number) {
    await setTreeData();
    if (value === AuthAllResource.NO) {
      allData.disableSelect = false;
      setFieldsValue({
        businessIds: allData.businessIds,
      });
    } else {
      allData.disableSelect = true;
      setFieldsValue({
        businessIds: allData.allBusinessIds,
      });
    }
  }
  async function setTreeData() {
    const { data } = await getBusinessResource({});
    data?.forEach((category: any) => {
      category.key = `category_${category.id}`;
      category?.subcategories?.forEach((subcategory: any) => {
        const subcategoryId = `subcategory_${subcategory.id}`;
        allData.allBusinessIds.push(subcategoryId);
        subcategory.key = subcategoryId;
      });
    });
    allData.treeData = data || [];
  }
</script>

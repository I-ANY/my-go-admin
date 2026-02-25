<template>
  <!-- <div> -->
    <BasicTable @register="registerTable">
      <template #toolbar>
        <a-button type="primary" @click="handleCreate">
          新增配置
        </a-button>
      </template>
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'action'">
          <TableAction :actions="[
            {
              icon: 'clarity:note-edit-line',
              title: '编辑',
              onClick: handleEdit.bind(null, record),
            },
            {
              icon: 'ic:outline-delete-outline',
              title: '删除',
              popConfirm: {
                title: '是否确认删除',
                placement: 'left',
                confirm: handleDelete.bind(null, record),
              },
            },
          ]" />
        </template>
      </template>
    </BasicTable>
    <EncryptionConfigurationModal @register="registerModal" @success="handleSuccess" />
  <!-- </div> -->
</template>
<script lang="ts" setup>
import { BasicTable, useTable, TableAction } from '@/components/Table';
import EncryptionConfigurationModal from './EncryptionConfigurationModal.vue';
import { getAllBasenameList, getEncryptionConfigurationList, getAllOsTypeList, getAllTarVersionList } from '@/linuxApi/config';
import { columns, searchFormSchema } from './data';
import { useModal } from '@/components/Modal';
import { message } from 'ant-design-vue';
import { deleteEncryptionConfiguration } from '@/linuxApi/config';
import { onBeforeMount, ref } from 'vue';

const basenameList = ref<any[]>([]);
const osTypeList = ref<any[]>([]);
const tarVersionList = ref<any[]>([]);

defineOptions({ name: 'EncryptionConfiguration' });
const [registerTable, { reload, getForm }] = useTable({
  title: '业务包配置列表',
  api: getEncryptionConfigurationList,
  columns,
  pagination: true,
  formConfig: {
    labelWidth: '100px',
    schemas: searchFormSchema(basenameList.value),
    autoSubmitOnEnter: true,
  },
  useSearchForm: true,
  showTableSetting: true,
  bordered: true,
  showIndexColumn: false,
  actionColumn: {
    width: 120,
    title: '操作',
    dataIndex: 'action',
  },
});

async function getBasenameList() {
  const res = await getAllBasenameList();
  basenameList.value = res.map((item: any) => ({ label: item, value: item }));
  // 动态更新表单 schema
  const form = getForm();
  await form.updateSchema({
    field: 'basename',
    componentProps: {
      options: basenameList.value,
      onChange: async (value: string) => {
        // 清空操作系统和版本号
        form.setFieldsValue({ os_type: undefined, tar_version: undefined });
        osTypeList.value = [];
        tarVersionList.value = [];
        
        if (value) {
          // 获取操作系统列表
          const osTypes = await getAllOsTypeList(value);
          osTypeList.value = osTypes.map((item: any) => ({ label: item, value: item }));
          await form.updateSchema({
            field: 'os_type',
            componentProps: {
              options: osTypeList.value,
            },
          });
        } else {
          // 清空操作系统选项
          await form.updateSchema({
            field: 'os_type',
            componentProps: {
              options: [],
            },
          });
          await form.updateSchema({
            field: 'tar_version',
            componentProps: {
              options: [],
            },
          });
        }
      },
    },
  });
  
  // 为操作系统添加 onChange 事件
  await form.updateSchema({
    field: 'os_type',
    componentProps: {
      options: osTypeList.value,
      onChange: async (value: string) => {
        // 清空版本号
        form.setFieldsValue({ tar_version: undefined });
        tarVersionList.value = [];
        
        const basename = form.getFieldsValue().basename;
        if (basename && value) {
          // 获取版本号列表
          const versions = await getAllTarVersionList(basename, value);
          tarVersionList.value = versions.map((item: any) => ({ label: item, value: item }));
          await form.updateSchema({
            field: 'tar_version',
            componentProps: {
              options: tarVersionList.value,
            },
          });
        } else {
          // 清空版本号选项
          await form.updateSchema({
            field: 'tar_version',
            componentProps: {
              options: [],
            },
          });
        }
      },
    },
  });
}

const [registerModal, { openModal }] = useModal();
// 新增配置
function handleCreate() {
  openModal(true, {
    record: null,
    isUpdate: false,
  });
}
// 编辑配置
function handleEdit (record: any) {
  openModal(true, {
    record: {...record},
    isUpdate: true,
  });
}
// 删除配置
async function handleDelete(record: any) {
  await deleteEncryptionConfiguration(record);
  message.success('删除成功');
  reload();
}
/**
 * 操作成功回调
 */
function handleSuccess() {
  reload();
}

onBeforeMount(() => {
  getBasenameList();
});
</script>
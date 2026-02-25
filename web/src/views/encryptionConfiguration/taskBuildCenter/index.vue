<template>
  <div>
    <BasicTable @register="registerTable">
      <template #toolbar>
        <a-button type="primary" @click="handleCreate">
          新增任务
        </a-button>
      </template>
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'action'">
          <TableAction :actions="[
            {
              label: '详情',
              tooltip: '查看任务详情',
              onClick: handleView.bind(null, record),
            },
          ]" />
        </template>
      </template>
    </BasicTable>
    <AddTaskModal @register="registerModal" @success="handleSuccess" />
    <TaskLogModal @register="registerLogModal" />
  </div>
</template>
<script lang="ts" setup>
import { BasicTable, useTable, TableAction } from '@/components/Table';
import AddTaskModal from './AddTaskModal.vue';
import TaskLogModal from './TaskLogModal.vue';
import { getAllOsTypeList, getAllOsTypeListByIsoSignAndBasename, getAllTarVersionList, getBasenameListByIsoSign, getTaskBuildCenterList } from '@/linuxApi/config';
import { columns, searchFormSchema } from './data';
import { useModal } from '@/components/Modal';
import { onBeforeMount, ref } from 'vue';


const basenameList = ref<any[]>([]);
const osTypeList = ref<any[]>([]);
const tarVersionList = ref<any[]>([]);


defineOptions({ name: 'TaskBuildCenter' });
const [registerTable, { reload, getForm }] = useTable({
  title: '任务构建中心列表',
  api: getTaskBuildCenterList,
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
  const res = await getBasenameListByIsoSign('');
  basenameList.value = res?.basenames?.map((item: any) => ({ label: item, value: item }));
  // 动态更新表单 schema
  const form = getForm();
  await form.updateSchema({
    field: 'iso_sign',
    componentProps: {
      onChange: async (e) => {
        const res = await getBasenameListByIsoSign(e.target.value);
        const isoSignList = res?.basenames?.map((item: any) => ({ label: item, value: item }));
        await form.updateSchema({
          field: 'basename',
          componentProps: {
            options: isoSignList,
          },
        });
        const result = await getAllOsTypeListByIsoSignAndBasename(e.target.value, form.getFieldsValue().basename);
        const osTypeList = result.os_types.map((item: any) => ({ label: item, value: item }));
        await form.updateSchema({
          field: 'os_type',
          componentProps: {
            options: osTypeList,
          },
        });
      },
    },
  },
  );
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
          const osTypes = await getAllOsTypeListByIsoSignAndBasename(form.getFieldsValue().iso_sign, value);
          osTypeList.value = osTypes.os_types.map((item: any) => ({ label: item, value: item }));
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
const [registerLogModal, { openModal: openLogModal }] = useModal();
// 新增配置
function handleCreate() {
  openModal(true, {
    record: null,
    isUpdate: false,
  });
}
// 操作成功回调
async function handleSuccess(isUpdate: boolean, record?: any) {
  reload();
  if (record) {
    try {
      // const logText = await getTaskBuildCenterDetail(taskRecord);
      openLogModal(true, {
        record,
      });
    } catch (error) {
      console.error('获取任务日志失败:', error);
    }
  }
}

// 查看任务详情
async function handleView(record: any) {
  // const res = await getTaskBuildCenterDetail(record);
  openLogModal(true, {
    record
  });
}

onBeforeMount(() => {
  getBasenameList();
});
</script>
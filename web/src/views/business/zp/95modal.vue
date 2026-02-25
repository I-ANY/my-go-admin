<template>
  <BasicModal
    v-bind="$attrs"
    destroyOnClose
    @register="register"
    title="95值计算"
    width="1000px"
    @fullscreen="onFullscreen"
    @open-change="onOpenChange"
    okText="发布"
    @ok="handleSubmit"
  >
    <BasicTable @register="registerTable" />
  </BasicModal>
</template>
<script lang="ts" setup>
  import { ref, nextTick } from 'vue';
  import { BasicModal, useModalInner } from '@/components/Modal';
  import { BasicTable, useTable } from '@/components/Table';
  import { traffic95Columns } from './data';
  import { GetZPTraffic95Gather, CreateZPTraffic95 } from '@/api/business/zp';
  import { message } from 'ant-design-vue';
  import { formatToDateTime } from '@/utils/dateUtil';

  const formData = ref({});

  const [register, { closeModal }] = useModalInner((data) => {
    formData.value = data.props.formData;
  });

  const [registerTable, { reload, getDataSource }] = useTable({
    title: '95值详情(单位:Gbps)',
    api: GetZPTraffic95Gather,
    columns: traffic95Columns.map((column) => {
      if (column.dataIndex === 'is_publish') {
        return {
          ...column,
          ifShow: false,
        };
      }
      return column;
    }),
    showTableSetting: false,
    showIndexColumn: false,
    pagination: false,
    maxHeight: 150,
    beforeFetch: (params) => {
      // 在这里可以使用传递过来的参数
      if (formData.value) {
        // 将表单参数合并到请求参数中
        Object.assign(params, formData.value);
      }
    },
  });

  const selectTable = ref<InstanceType<typeof BasicTable> | undefined>();

  const onFullscreen = async () => {
    await nextTick();
    selectTable.value?.redoHeight();
  };
  const onOpenChange = (open: boolean) => {
    if (open) {
      reload();
    }
  };

  async function handleSubmit() {
    // 获取当前表格数据
    const tableData = getDataSource();
    const listData: any[] = [];
    let date = '';
    if (tableData.length > 0) {
      tableData.forEach((item) => {
        date = formatToDateTime(tableData[0].ecdn_timestmp).split(' ')[0];
        listData.push({ ...item });
      });
    }
    try {
      await CreateZPTraffic95({ date: date, items: listData }).then(() => {
        closeModal();
        message.success('95值发布成功');
      });
    } catch (e) {
      message.error('95值发布失败');
    }
  }
</script>

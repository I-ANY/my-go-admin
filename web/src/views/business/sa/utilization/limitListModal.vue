<template>
  <BasicModal
    @register="registerModal"
    v-bind="$attrs"
    :title="title"
    width="1300px"
    :footer="null"
  >
    <BasicTable @register="registerTable" ref="selectTable">
      <template #toolbar>
        <div style=" width: 95%;text-align: left">
          <span>（总共拉黑{{ lineCounts }}条线路，{{ bw_limit_total }}G带宽）</span>
        </div>
      </template>
    </BasicTable>
  </BasicModal>
</template>
<script setup lang="ts">
  import { ref, nextTick } from 'vue';
  import { BasicModal, useModalInner } from '@/components/Modal';
  import { BasicTable, useTable } from '@/components/Table';
  import { LimitListColumns } from './data';
  import { GetLimitList } from '@/api/business/sa';

  let title = ref('');
  let lineCounts = ref(0);
  let bw_limit_total = ref(0);
  const modalData = ref({
    date: '',
    guid: '',
  });

  // 弹窗
  const [registerModal] = useModalInner(async (data) => {
    // 接收并存储传入的参数
    if (data) {
      modalData.value = {
        date: data.date,
        guid: data.guid,
      };
      title.value = data.date + '拉黑线路详情';
    }
    await nextTick();
    // 重置表格数据
    await reloadTable(); // 重新加载表格数据
  });

  const [registerTable, { reload: reloadTable }] = useTable({
    canResize: true,
    title: '详情信息',
    api: GetLimitList,
    columns: LimitListColumns,
    rowKey: 'id',
    beforeFetch: (params) => {
      params.date = modalData.value.date;
      params.guid = modalData.value.guid;
    },
    afterFetch: (data) => {
      lineCounts.value = data.length;
      bw_limit_total.value = 0;
      data.forEach((item) => {
        bw_limit_total.value += item.bw_limit;
      });
      console.log(lineCounts, bw_limit_total);
    },
  });
</script>
<style scoped lang="less"></style>

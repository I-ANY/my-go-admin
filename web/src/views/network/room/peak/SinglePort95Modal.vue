<template>
  <BasicModal v-bind="$attrs" @register="registerModal" @cancel="onModalCancel"
    ><BasicTable
      @register="registerTable"
      :dataSource="data.items"
      :loading="data.loading"
      @change="onTableChange"
    >
      <template #toolbar>
        <Button type="primary" @click="reloadTable">刷新</Button>
      </template>
      <template #bodyCell="{ column, record, index }">
        <template v-if="column.key === 'action'">
          <template v-if="index < data.items.length - 1">
            <TableAction
              :actions="[
                {
                  icon: 'material-symbols:monitoring',
                  tooltip: '95值拉升预测',
                  onClick: handleImprove95Prediction.bind(null, record),
                  label: '预测',
                },
              ]"
            />
          </template>
          <template v-else>-</template>
        </template>
      </template>
    </BasicTable>
    <Improve95PredictionModal
      @register="registerImprove95PredictionModal"
      @success="onImprove95PredictionSuccess"
    />
  </BasicModal>
</template>
<script lang="ts" setup>
  import { BasicTable, useTable, TableAction } from '@/components/Table';
  import { BasicModal, useModalInner, useModal } from '@/components/Modal';
  import { getNumberResult, SinglePort95Columns } from './data';
  import { reactive } from 'vue';
  import { getSinglePort95 } from '@/api/network/room_peak';
  import { Button } from 'ant-design-vue';
  import Improve95PredictionModal from './Improve95PredictionModal.vue';
  import { customCeilDivide } from '@/utils/util';

  defineOptions({ name: 'SinglePort95Modal' });
  const emit = defineEmits(['register', 'success']);
  const [registerImprove95PredictionModal, { openModal: openImprove95PredictionModal }] =
    useModal();
  let roomRecord: any = undefined;
  let data = reactive({
    items: [] as any[],
    summary: {} as any,
    loading: false,
    sorter: {} as any,
  });

  const [registerModal, { setModalProps, closeModal }] = useModalInner(async (d) => {
    data.items = [];
    data.summary = {};
    data.sorter = {};
    data.loading = false;
    roomRecord = d.record;
    setModalProps({
      title: `${roomRecord.name}（${roomRecord.month}）端口95详情`,
      width: 1000,
      height: 550,
      destroyOnClose: true,
      showCancelBtn: false,
      showOkBtn: false,
      canFullscreen: false,
    });
    reloadTable();
  });

  const [registerTable, { getDataSource }] = useTable({
    // title: '95值拉高预测',
    columns: SinglePort95Columns(getTabelFunctions),
    size: 'small',
    // beforeFetch: (params) => {
    //   params.id = record.id;
    //   return params;
    // },
    canResize: true,
    scroll: { y: 450 },
    useSearchForm: false,
    showTableSetting: false,
    bordered: true,
    showIndexColumn: false,
    rowKey: 'id',
    pagination: false,
    actionColumn: {
      width: 80,
      title: '操作',
      dataIndex: 'action',
    },
  });
  function onModalCancel() {
    closeModal();
    emit('success');
  }
  async function reloadTable() {
    data.loading = true;
    try {
      let res = await getSinglePort95(roomRecord.id, {
        ...data.sorter,
      });
      data.items = res.items || [];
      const summary = res.summary || {};
      if (data.items.length > 0 && res.summary) {
        data.items.push({
          roomName: `总计(${data.items.length}个端口)`,
          monthTotalPointCount: summary.monthTotalPointCount,
          maxPeakShavingPointCount: summary.maxPeakShavingPointCount,
          used95PointCount: summary.used95PointCount,
          rest95PointCount: summary.rest95PointCount,
          restPeakShavingMinute: summary.restPeakShavingMinute,
          now95Bw: summary.now95Bw,
          month95Bw: summary.month95Bw,
        });
      }
    } finally {
      data.loading = false;
    }
  }
  function onTableChange(_pagination: any, _filters: any, sorter: any, _extra: any) {
    if (sorter && sorter.field && sorter.order) {
      data.sorter = {
        field: sorter.field,
        order: sorter.order,
      };
    } else {
      data.sorter = {};
    }
    reloadTable();
  }
  function getTabelFunctions(): any {
    return {
      getDataSource,
    };
  }
  function handleImprove95Prediction(record: any) {
    let preset95Bwbps = getNumberResult(
      customCeilDivide(record.preset95Bw, 1000 * 1000 * 1000, 2) as any,
    );
    const title = `95值拉升预测-${record.roomName}-${record.name}-预95值：${preset95Bwbps || ' - '}Gbps`;
    const improve95Prediction = record.improve95Prediction || [];
    openImprove95PredictionModal(true, {
      title,
      improve95Prediction,
    });
  }
  function onImprove95PredictionSuccess() {
    reloadTable();
  }
</script>

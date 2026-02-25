<template>
  <div class="demand-gap-container">
    <div class="header">
      <h1>腾讯业务需求缺口地区统计</h1>
      <p class="date">{{ data.date }}</p>
      <a-button type="primary" @click="handleRefresh" :loading="data.loading">
        <template #icon><ReloadOutlined /></template>
        刷新
      </a-button>
    </div>

    <div class="tables-wrapper">
      <!-- KP2业务 -->
      <div class="table-section">
        <div class="table-title">腾讯KP2业务需求缺口地区</div>
        <a-table
          :columns="columns"
          :data-source="data.kp2"
          :loading="data.loading"
          :pagination="false"
          :rowKey="(_record, index) => `row-${index}`"
          :scroll="{ y: 'calc(100vh - 280px)' }"
          size="small"
          bordered
        />
      </div>

      <!-- KPC业务 -->
      <div class="table-section">
        <div class="table-title">腾讯KPC业务需求缺口地区</div>
        <a-table
          :columns="columns"
          :data-source="data.kpc"
          :loading="data.loading"
          :pagination="false"
          :rowKey="(_record, index) => `row-${index}`"
          :scroll="{ y: 'calc(100vh - 280px)' }"
          size="small"
          bordered
        />
      </div>

      <!-- 专线业务 -->
      <div class="table-section">
        <div class="table-title">腾讯非80专线需求缺口地区</div>
        <a-table
          :columns="columns"
          :data-source="data.specialLine"
          :loading="data.loading"
          :pagination="false"
          :rowKey="(_record, index) => `row-${index}`"
          :scroll="{ y: 'calc(100vh - 280px)' }"
          size="small"
          bordered
        />
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
  import { onMounted, reactive } from 'vue';
  import { Table } from 'ant-design-vue';
  import { ReloadOutlined } from '@ant-design/icons-vue';
  import { GetDemandGapList } from '@/api/business/k';
  import { useMessage } from '@/hooks/web/useMessage';

  const ATable = Table;
  const { notification } = useMessage();

  const columns = [
    { title: '大区', dataIndex: 'area', width: 100, align: 'center' as const },
    { title: '省份', dataIndex: 'province', width: 120, align: 'center' as const },
    { title: '电信', dataIndex: 'dx', width: 80, align: 'center' as const },
    { title: '联通', dataIndex: 'lt', width: 80, align: 'center' as const },
    { title: '移动', dataIndex: 'yd', width: 80, align: 'center' as const },
  ];

  const data = reactive({
    loading: false,
    date: '',
    kp2: [] as any[],
    kpc: [] as any[],
    specialLine: [] as any[],
  });

  const fetchData = async () => {
    data.loading = true;
    try {
      const res: any = await GetDemandGapList();
      if (res) {
        data.date = res.date || '';
        data.kp2 = res.kp2 || [];
        data.kpc = res.kpc || [];
        data.specialLine = res.special_line || [];
      }
    } catch (error: any) {
      notification.error({
        message: '获取数据失败',
        description: error.message || '未知错误',
      });
    } finally {
      data.loading = false;
    }
  };

  const handleRefresh = () => {
    fetchData();
  };

  onMounted(() => {
    fetchData();
  });

  defineOptions({ name: 'DemandGapList' });
</script>

<style lang="less" scoped>
  .demand-gap-container {
    display: flex;
    flex-direction: column;
    height: 100vh;
    padding: 24px;
    overflow: hidden;
    background: #f0f2f5;
  }

  .header {
    display: flex;
    flex-shrink: 0;
    align-items: center;
    justify-content: space-between;
    margin-bottom: 16px;
    padding: 20px 24px;
    border-radius: 4px;
    background: #fff;
    box-shadow: 0 2px 8px rgb(0 0 0 / 10%);

    h1 {
      margin: 0;
      color: #333;
      font-size: 20px;
      font-weight: 600;
    }

    .date {
      margin: 0;
      color: #666;
      font-size: 14px;
    }
  }

  .tables-wrapper {
    display: flex;
    flex: 1;
    flex-direction: row;
    min-height: 0;
    overflow: auto hidden;
    gap: 16px;
  }

  .table-section {
    display: flex;
    flex: 1;
    flex-direction: column;
    min-width: 360px;
    overflow: hidden;
    border-radius: 4px;
    background: #fff;
    box-shadow: 0 2px 8px rgb(0 0 0 / 10%);
  }

  .table-title {
    padding: 12px 16px;
    border-bottom: 1px solid #f0f0f0;
    background: #fafafa;
    color: #333;
    font-size: 16px;
    font-weight: 600;
  }

  :deep(.ant-table) {
    .ant-table-thead > tr > th {
      background: #fafafa;
      font-weight: 600;
    }

    .ant-table-tbody > tr > td {
      padding: 8px 16px;
    }
  }
</style>

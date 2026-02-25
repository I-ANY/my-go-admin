<template>
  <div>
    <BasicTable @register="registerTable">
      <template #toolbar>
        <a-button
          style="background-color: red; color: white"
          @click="handleMockOpen"
          v-if="getDataSource()?.length > 0"
        >
          营收模拟
        </a-button>
      </template>
      <template #bodyCell="{ column, record }">
        <template v-if="column.key == 'category'">
          <span>{{ record.category ? record.category.name : '/' }}</span>
        </template>
        <template v-if="column.key == 'region'">
          {{ record.region ? record.region.name : '/' }}
        </template>
        <template v-if="column.key == 'zone'">
          {{ record.zone ? record.zone.name : '/' }}
        </template>
        <template v-if="column.key == 'isNotlocalIsp'">
          {{ record.bizIspMode || record.bizIsp ? '是' : '否' }}
        </template>
        <template v-if="column.key == 'bizIsp'">
          {{ record.bizIsp ? record.bizIsp : '/' }}
        </template>
        <template v-if="column.key == 'bizIspMode'">
          {{ record.bizIspMode == 2 ? '异网计费' : record.bizIspMode == 1 ? '本网计费' : '/' }}
        </template>
        <template v-if="column.key == 'mode'">
          {{ record.mode ? record.mode.name : '/' }}
        </template>
        <template v-if="column.key == 'hosts'">
          <a-button
            type="primary"
            @click="
              handleQueryHost(
                record.categoryID,
                record.owner,
                record.bizIsp,
                record.region?.name,
                record.zone?.name,
                record.localIsp,
                record.location,
              )
            "
            >详情</a-button
          >
        </template>
      </template>
    </BasicTable>
    <Modal
      :bodyStyle="{ height: '300px', overflowY: 'auto' }"
      v-model:open="open"
      title="设备详情"
      @ok="open = false"
    >
      <template v-for="item in hosts.items" :key="item.id">
        <div
          style="
            display: flex;
            align-items: center;
            padding: 8px 16px;
            border-bottom: 1px solid #f0f0f0;
          "
        >
          <span style="flex: 1; text-align: left">{{ item.category || '无分类' }}</span>
          <span style="flex: 2; text-align: left">{{ item.biz || '无业务' }}</span>
          <span style="flex: 3; text-align: left">{{ item.hostname || '无主机名' }}</span>
        </div>
      </template>
    </Modal>
    <Modal
      width="800px"
      :bodyStyle="{ height: '400px', overflow: 'auto' }"
      v-model:open="mockOpen"
      :title="selectedNodeName ? `${selectedNodeName}-营收模拟` : '营收模拟'"
      okText="我已知晓"
      @ok="mockOpen = false"
    >
      <template #title>
        <div style="display: flex; align-items: center; gap: 8px">
          <span>{{ selectedNodeName ? `${selectedNodeName}-营收模拟` : '营收模拟' }}</span>
        </div>
      </template>
      <div style="margin-left: 16px">
        <div>
          <b style="font-size: 16px">业务选择</b>
          <Tooltip placement="top" :overlayStyle="{ maxWidth: '400px' }">
            <template #title>
              <div
                >业务单价受 "业务类型、本网运营商、大区、省份、异网运营商、跨网计费方式"
                共同决定，任一参数不同，都可能引发单价变化</div
              >
            </template>
            <QuestionCircleOutlined style="font-size: 16px; cursor: pointer" />
          </Tooltip>
        </div>
        <div>
          <ApiSelect
            v-model:value="selectedValue"
            style="width: 99%"
            :api="fetchOptions"
            :params="searchParams"
            @search="handleSearch"
            @change="handleBusinessChange"
            :filterOption="false"
            showSearch
            allowClear
            placeholder="请搜索选择"
          />
        </div>
        <div style="margin-top: 16px">
          <b style="font-size: 16px">模拟结果</b>
          <Tooltip placement="top" :overlayStyle="{ maxWidth: '400px' }">
            <template #title>
              <div>
                <div style="margin-bottom: 8px"
                  >计算逻辑：投资回报率 = (业务单价 - 节点单价) / 节点单价</div
                >
                <div style="margin-bottom: 4px">
                  <span style="color: #52c41a">盈利S ></span>
                  <span style="color: #73d13d">盈利A ></span>
                  <span style="color: #95de64">盈利B ></span>
                  <span style="color: #b7eb8f">盈利C ></span>
                  <span style="color: #1890ff">收支平衡 ></span>
                  <span style="color: #f5222d">亏损</span>
                </div>
                <div style="margin-bottom: 4px">
                  <span style="color: #1890ff">收支平衡： </span>
                  <span> 投资回报率 = 0</span>
                </div>
                <div style="white-space: nowrap">
                  <span style="color: #f5222d">亏损：</span>
                  <span> 投资回报率 &lt; 0</span>
                </div>
              </div>
            </template>
            <QuestionCircleOutlined style="font-size: 16px; cursor: pointer" />
          </Tooltip>
        </div>
        <div v-if="loadingMock" style="padding: 20px; text-align: center">
          <a-spin />
        </div>
        <div v-else-if="mockResults.length > 0" style="margin-top: 16px">
          <table style="width: 99%; border-collapse: collapse">
            <thead>
              <tr style="background-color: #f5f5f5">
                <th style="padding: 8px; border: 1px solid #d9d9d9; text-align: left">盈亏状态</th>
                <th style="padding: 8px; border: 1px solid #d9d9d9; text-align: left">节点名称</th>
                <th style="padding: 8px; border: 1px solid #d9d9d9; text-align: left">所在地</th>
                <th style="padding: 8px; border: 1px solid #d9d9d9; text-align: left">运营商</th>
              </tr>
            </thead>
            <tbody>
              <tr
                v-for="item in mockResults"
                :key="item.name"
                style="border-bottom: 1px solid #f0f0f0"
              >
                <td style="padding: 8px; border: 1px solid #d9d9d9">
                  <span :style="{ color: getProfitColor(item.profit) }">
                    {{ getProfitText(item.profit) }}
                  </span>
                </td>
                <td style="padding: 8px; border: 1px solid #d9d9d9">{{ item.name }}</td>
                <td style="padding: 8px; border: 1px solid #d9d9d9">{{ item.location }}</td>
                <td style="padding: 8px; border: 1px solid #d9d9d9">{{ item.localIsp }}</td>
              </tr>
            </tbody>
          </table>
        </div>
        <div v-else-if="selectedValue" style="padding: 20px; color: #999; text-align: center">
          暂无模拟数据
        </div>
      </div>
    </Modal>
  </div>
</template>
<script lang="ts" setup>
  import { ref, reactive, h, computed } from 'vue';
  import { Modal, message, Tooltip } from 'ant-design-vue';
  import { BasicTable, useTable, BasicColumn, FormSchema } from '@/components/Table';
  import {
    getOwnerRecord,
    getHosts,
    getOwners,
    getRecordOption,
    mockProfit,
    getBusiness,
  } from '@/api/price/price';
  import { ApiSelect } from '@/components/Form';
  import { useDebounceFn } from '@vueuse/core';
  import { QuestionCircleOutlined } from '@ant-design/icons-vue';

  defineOptions({ name: 'PriceOwnerQuery' });

  const open = ref<boolean>(false);
  const mockOpen = ref<boolean>(false);
  const businessGroupList = ref([]);
  const hosts = reactive({
    items: [] as { id: string; biz: string; category: string; hostname: string }[],
  });

  const selectedValue = ref<string>('');
  const selectedNodeName = ref<string>('');
  const keyword = ref<string>('');
  const mockResults = ref<any[]>([]);
  const loadingMock = ref<boolean>(false);

  // 使用computed来创建searchParams，这样当keyword变化时会自动更新
  const searchParams = computed(() => {
    return { query: keyword.value };
  });

  const cellContent = () => ({
    colSpan: 1,
  });

  function getMergeHeaderColumns(): BasicColumn[] {
    return [
      {
        title: '本网运营商',
        dataIndex: 'localIsp',
        width: 100,
        customCell: cellContent,
      },
      {
        title: '节点',
        dataIndex: 'owner',
        minWidth: 100,
        maxWidth: 300,
        customCell: cellContent,
      },
      {
        title: '所在地',
        dataIndex: 'location',
        minWidth: 100,
        maxWidth: 300,
        customCell: cellContent,
      },
      {
        title: '业务组',
        dataIndex: 'category',
        width: 150,
        customCell: cellContent,
      },
      {
        title: '业务大区',
        dataIndex: 'region',
        width: 100,
        customCell: cellContent,
      },
      {
        title: '业务省份',
        dataIndex: 'zone',
        width: 100,
        customCell: cellContent,
      },
      {
        title: '是否异网',
        dataIndex: 'isNotlocalIsp',
        width: 100,
        customCell: cellContent,
      },
      {
        title: '异网运营商',
        dataIndex: 'bizIsp',
        width: 100,
        customCell: cellContent,
      },
      {
        title: '跨网计费方式',
        dataIndex: 'bizIspMode',
        width: 100,
        customCell: cellContent,
      },
      {
        title: '计费方式',
        dataIndex: 'mode',
        width: 300,
        customCell: cellContent,
      },
      {
        title: '盈亏状态',
        dataIndex: 'profit',
        width: 100,
        sorter: true,
        customCell: cellContent,
        customRender: ({ text }) => {
          const profitMap = {
            0: { label: '未知', color: '#999999' }, // 灰色
            1: { label: '盈利S', color: '#52c41a' }, // 绿色
            2: { label: '盈利A', color: '#73d13d' }, // 浅绿色
            3: { label: '盈利B', color: '#95de64' }, // 更浅绿色
            4: { label: '盈利C', color: '#b7eb8f' }, // 最浅绿色
            5: { label: '收支平衡', color: '#1890ff' }, // 蓝色
            6: { label: '亏损', color: '#f5222d' }, // 红色
          };
          // 确保text是数字类型
          const profitValue = Number(text);
          const config = profitMap[profitValue] || { label: '/', color: '#000000' };

          return h('span', { style: { color: config.color } }, config.label);
        },
      },
      {
        title: '设备详情',
        dataIndex: 'hosts',
        width: 100,
        customCell: cellContent,
      },
      {
        title: '备注',
        dataIndex: 'describe',
        customCell: cellContent,
      },
    ];
  }
  const searchFormSchema: FormSchema[] = [
    {
      label: '节点',
      field: 'owner',
      component: 'ApiSelect',
      componentProps: {
        api: getOwners,
        resultField: 'items',
        labelField: 'name',
        valueField: 'id',
        showSearch: true,
        filterOption: (input: string, option: any) => {
          return (
            option.label?.toLowerCase().indexOf(input.toLowerCase()) >= 0 ||
            option.identify?.toLowerCase().indexOf(input.toLowerCase()) >= 0
          );
        },
        onChange: async (item) => {
          const result = await getBusiness({ owner: item });
          businessGroupList.value = result.map((item: { name: string; id: string }) => ({
            label: item.name,
            value: item.id,
          }));
          // 清空业务组的选择
          const form = getForm();
          if (form) {
            form.setFieldsValue({ category: undefined });
          }
        },
      },
      required: true,
      colProps: { span: 5 },
    },
    {
      label: '业务组',
      field: 'category',
      component: 'Select',
      componentProps: () => ({
        options: businessGroupList.value,
        placeholder: '请选择业务组',
      }),
      colProps: { span: 5 },
    },
    {
      field: 'profit',
      label: '盈亏状态',
      component: 'Select',
      componentProps: {
        options: [
          { label: '未知', value: 0 },
          { label: '盈利S', value: 1 },
          { label: '盈利A', value: 2 },
          { label: '盈利B', value: 3 },
          { label: '盈利C', value: 4 },
          { label: '收支平衡', value: 5 },
          { label: '亏损', value: 6 },
        ],
      },
      defaultValue: null,
      colProps: { span: 4 },
    },
  ];

  const [registerTable, { getDataSource, getForm }] = useTable({
    title: '业务单价列表',
    api: getOwnerRecord,
    columns: getMergeHeaderColumns(),
    formConfig: {
      labelWidth: 120,
      schemas: searchFormSchema,
      autoSubmitOnEnter: true,
    },
    useSearchForm: true,
    bordered: true,
    pagination: false,
    showIndexColumn: false,
    showTableSetting: true,
  });

  function handleQueryHost(
    id: Recordable,
    owner: Recordable,
    bizIsp: Recordable,
    region: Recordable,
    zone: Recordable,
    localIsp: Recordable,
    location: Recordable,
  ) {
    getHosts({
      id: id,
      owner: owner,
      bizIsp: bizIsp,
      region: region,
      zone: zone,
      localIsp: localIsp,
      location: location,
    })
      .then(function (resp) {
        hosts.items = [];
        for (var i = 0; i < resp.items.length; i++) {
          hosts.items.push({
            id: resp.items[i].id,
            biz: resp.items[i].biz,
            category: resp.items[i].category,
            hostname: resp.items[i].hostname,
          });
        }
        open.value = true;
      })
      .catch(function (err) {
        message.error(err);
      });
  }

  function handleMockOpen() {
    selectedValue.value = '';
    keyword.value = '';
    mockResults.value = [];

    // 获取当前选中的群号作为节点名称
    const currentOwner = getDataSource()?.[0]?.owner || '';
    selectedNodeName.value = currentOwner;

    mockOpen.value = true;
  }

  // 处理业务选择变化
  async function handleBusinessChange(value: string) {
    if (!value) {
      mockResults.value = [];
      // 当清除选择时，重新调用getRecordOption刷新业务选项
      keyword.value = '';
      return;
    }

    // 获取当前选中的群号
    const currentOwner = getDataSource()?.[0]?.owner || '';

    loadingMock.value = true;
    try {
      const res = await mockProfit({
        id: value,
        owner: currentOwner,
      });
      mockResults.value = res?.items || [];
    } catch (error) {
      message.error('模拟请求失败');
      mockResults.value = [];
    } finally {
      loadingMock.value = false;
    }
  }

  const fetchOptions = async () => {
    // 使用getRecordOption API，但需要确保参数正确传递
    const res = await getRecordOption(searchParams.value);
    return res.items.map((item) => ({
      label: item.value,
      value: item.id,
    }));
  };

  // 使用防抖函数来处理搜索，增加延迟时间避免频繁请求
  const handleSearch = useDebounceFn((value: string) => {
    keyword.value = value;
  }, 800); // 增加到800ms延迟

  // 获取盈亏状态文本
  function getProfitText(profit: number): string {
    const profitMap = {
      0: '未知',
      1: '盈利S',
      2: '盈利A',
      3: '盈利B',
      4: '盈利C',
      5: '收支平衡',
      6: '亏损',
    };
    return profitMap[profit] || '/';
  }

  // 获取盈亏状态颜色
  function getProfitColor(profit: number): string {
    const colorMap = {
      0: '#999999', // 灰色
      1: '#52c41a', // 绿色
      2: '#73d13d', // 浅绿色
      3: '#95de64', // 更浅绿色
      4: '#b7eb8f', // 最浅绿色
      5: '#1890ff', // 蓝色
      6: '#f5222d', // 红色
    };
    return colorMap[profit] || '#000000';
  }
</script>

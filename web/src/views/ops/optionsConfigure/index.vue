<template>
  <div>
    <div class="h-full flex p-3 -space-x-3">
      <div class="flex w-1/6">
        <Col :span="24">
          <BasicTree
            title="  "
            search
            :treeData="navTreeData"
            style="min-width: 200px"
            :actionList="actionList"
            :renderIcon="createIcon"
            @select="handleNavSelect"
          />
        </Col>
      </div>
      <div class="flex-1 flex flex-col w-5/6" style="margin-top: 1px">
        <div class="h-1/3 mb-2">
          <BasicTable @register="registerTable" style="padding-top: 0">
            <template #toolbar>
              <a-button type="primary" @click="handleCreateScript" style="margin-right: auto">
                新 增
              </a-button>
            </template>
            <template #bodyCell="{ column, record }">
              <template v-if="column.key === 'params'">
                <Tooltip
                  placement="topLeft"
                  :overlayStyle="{ maxWidth: '1000px', maxHeight: '400px', overflowY: 'auto' }"
                >
                  <template #title>
                    <div class="content">{{ record.params }}</div>
                  </template>
                  <span> {{ record.params }}</span>
                </Tooltip>
              </template>
              <template v-if="column.key === 'action'">
                <TableAction
                  :actions="[
                    {
                      icon: 'clarity:note-edit-line',
                      tooltip: '编辑',
                      onClick: handleEditScript.bind(null, record),
                    },
                    {
                      icon: 'ant-design:copy-outlined',
                      tooltip: '拷贝',
                      onClick: handleCopyScript.bind(null, record),
                    },
                    {
                      icon: 'ant-design:delete-outlined',
                      color: 'error',
                      tooltip: '删除',
                      popConfirm: {
                        title: '是否确认删除',
                        placement: 'left',
                        confirm: handleDeleteScript.bind(null, record),
                      },
                    },
                  ]"
                />
              </template>
            </template>
          </BasicTable>
        </div>
      </div>
    </div>
    <OptionDrawer @register="registerDrawer" @success="handleDrawerEditSuccess" />
    <OptionsGroupModal @register="registerModal" @success="getTreeData" />
    <OptionsCopyModal @register="registerCopyModal" @success="handleCopySuccess" />
  </div>
</template>

<script lang="ts" setup>
  import { Col, message, Tooltip, Popconfirm } from 'ant-design-vue';
  import { BasicTree, TreeActionItem } from '@/components/Tree';
  import { BasicTable, TableAction, useTable } from '@/components/Table';
  import { onMounted, h, reactive } from 'vue';
  import { PlusOutlined, DeleteOutlined, EditOutlined, CopyOutlined } from '@ant-design/icons-vue';
  import { taskConfigColumns, taskConfigSearchSchema } from './data';
  import { useDrawer } from '@/components/Drawer';
  import { useModal } from '@/components/Modal';
  import {
    GetScriptTasks,
    // EditScriptTask,
    DelOperation,
    DelScriptTask,
    // AddScriptTask,
  } from '@/api/ops/execute';
  import OptionDrawer from './optionsDrawer.vue';
  import OptionsGroupModal from './optionsGroupModal.vue';
  import OptionsCopyModal from './optionsCopyModal.vue';
  import { GetAuthedBiz } from '@/api/business/biz';

  const [registerDrawer, { openDrawer }] = useDrawer();
  const [registerModal, { openModal }] = useModal();
  const [registerCopyModal, { openModal: openCopyModal }] = useModal();

  // 导航树数据
  let navTreeData = reactive<any[]>([
    {
      title: '公共',
      id: 0,
      key: 0,
      children: [],
    },
  ]);
  let businessId = 0;
  let optionGroupId = null;

  const [registerTable, { reload }] = useTable({
    api: GetScriptTasks,
    rowKey: 'id',
    columns: taskConfigColumns,
    formConfig: {
      labelWidth: 120,
      schemas: taskConfigSearchSchema,
      autoSubmitOnEnter: true,
    },
    useSearchForm: true,
    showTableSetting: false,
    bordered: true,
    showIndexColumn: false,
    immediate: false,
    actionColumn: {
      width: 100,
      title: '操作',
      dataIndex: 'action',
      fixed: 'right',
    },
    beforeFetch: async (params) => {
      params.parentId = optionGroupId;
      params.businessId = businessId;
      params.type = 'function';
      return params;
    },
  });

  // 处理导航选择
  async function handleNavSelect(selectedKeys, { selectedNodes }) {
    if (selectedKeys.length > 0) {
      businessId = selectedNodes[0].id;
      optionGroupId = null;
      console.log('selectedNodes', selectedNodes);
      if (selectedNodes[0].parentId || selectedNodes[0].parentId === 0) {
        businessId = selectedNodes[0].parentId;
        optionGroupId = selectedNodes[0].id;
      }
      // 根据选中的节点加载对应的表格数据
      // parentId = selectedNodes[0].parentId;
      // const resTasks = await GetScriptTasks({
      //   parentId: selectedNodes[0].parentId,
      //   businessId: selectedNodes[0].id,
      // });
      // // 手动设置表格数据
      // setTableData(resTasks.items);
      await reload();
    }
  }

  const actionList: TreeActionItem[] = [
    {
      render: (node) => {
        let nodeData = { ...node };
        // 只有第一层(level=0)显示新增按钮
        // 如果无法获取 level 信息，可以通过 node 的其他属性判断
        // 例如通过 node.children 或其他标识字段
        if (node.children) {
          // 假设只有父节点有 children 属性
          return h(PlusOutlined, {
            class: 'ml-4',
            onClick: (e) => {
              e.stopPropagation(); // 阻止事件冒泡
              handleCreateOperationGroup(nodeData);
            },
          });
        }
        if (!node.children) {
          // 或其他判断条件
          return [
            // 编辑按钮
            h(EditOutlined, {
              class: 'mr-3',
              onClick: (e) => {
                e.stopPropagation();
                handleEditOperationGroup(nodeData);
              },
            }),
            // 拷贝按钮
            h(CopyOutlined, {
              class: 'mr-3',
              onClick: (e) => {
                e.stopPropagation();
                handleCopyGroup(nodeData);
              },
            }),
            // 在 actionList 中修改删除按钮
            h(
              Popconfirm,
              {
                title: '确定要删除此操作组吗？',
                okText: '确认',
                cancelText: '取消',
                onConfirm: (e) => {
                  e.stopPropagation();
                  handleDeleteOperationGroup(nodeData);
                },
                onCancel: (e) => {
                  e.stopPropagation();
                },
              },
              {
                default: () =>
                  h(DeleteOutlined, {
                    onClick: (e) => {
                      e.stopPropagation();
                    },
                  }),
              },
            ),
          ];
        }
        return null; // 其他层级不显示
      },
    },
  ];

  function createIcon({ level }) {
    if (level === 1) {
      return 'ion:git-compare-outline';
    }
    if (level === 2) {
      return 'ion:home';
    }
    if (level === 3) {
      return 'ion:airplane';
    }
    return '';
  }

  // 新增操作组
  function handleCreateOperationGroup(node: any) {
    openModal(true, {
      isUpdate: false,
      businessId: node.id,
    });
  }

  // 编辑操作组
  function handleEditOperationGroup(node: any) {
    openModal(true, {
      isUpdate: true,
      title: node.title,
      optionGroupId: node.id,
      businessId: node.parentId,
    });
  }

  // 删除操作组
  async function handleDeleteOperationGroup(node: any) {
    try {
      await DelOperation(node.id);
      message.success('删除成功');
    } finally {
      await getTreeData();
    }
  }

  // 新增脚本
  function handleCreateScript() {
    openDrawer(true, {
      isUpdate: false,
      businessId: businessId,
      optionGroupId: optionGroupId,
    });
  }

  // 编辑脚本
  function handleEditScript(record: any) {
    openDrawer(true, {
      isUpdate: true,
      record,
    });
  }

  // 删除脚本
  async function handleDeleteScript(record: any) {
    // 实现删除逻辑
    try {
      await DelScriptTask(record.id);
      message.success('操作成功');
    } finally {
      await reload();
    }
  }

  // 拷贝单个脚本
  function handleCopyScript(record: any) {
    openCopyModal(true, {
      sourceId: record.id,
      sourceType: 'function',
    });
  }

  // 拷贝操作组（组下所有功能）
  function handleCopyGroup(node: any) {
    openCopyModal(true, {
      sourceId: node.id,
      sourceType: 'group',
    });
  }

  // 拷贝成功后刷新树和表格
  async function handleCopySuccess() {
    await getTreeData();
    await reload();
  }

  async function handleDrawerEditSuccess() {
    await reload();
  }

  async function getTreeData() {
    navTreeData.length = 1;
    navTreeData[0].children = [];
    const resTasks = await GetScriptTasks({
      type: 'group',
      parentId: 0,
      pageIndex: 1,
      pageSize: 1000,
    });
    // 处理公共的内容
    if (resTasks.items) {
      resTasks.items.forEach((option: any) => {
        if (option.businessId == 0) {
          navTreeData[0].children.push({
            title: option.name,
            key: '_' + option.id,
            id: option.id,
            parentId: 0,
          });
        }
      });
    }
    const resBiz = await GetAuthedBiz({});
    if (resBiz.categories) {
      resBiz.categories.forEach((item: any) => {
        const children: any[] = [];
        if (resTasks.items) {
          resTasks.items.forEach((option: any) => {
            if (option.businessId == item.id) {
              children.push({
                title: option.name,
                key: '_' + option.id,
                id: option.id,
                parentId: item.id,
              });
            }
          });
        }
        navTreeData.push({
          title: item.name,
          key: item.id,
          id: item.id,
          children: children,
        });
      });
    }
  }

  onMounted(async () => {
    await getTreeData();
  });
</script>
<style scoped>
  .content {
    word-break: break-word;
    white-space: pre-wrap;
  }
</style>

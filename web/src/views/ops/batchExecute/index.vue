<template>
  <div>
    <!-- 确保有根div -->
    <PageWrapper contentBackground contentClass="p-4" style="text-align: center">
      <div style="display: flex; justify-content: center">
        <BasicForm @register="registerForm" style="width: 100%">
          <template #submitBefore>
            <div style="display: flex; justify-content: center">
              <AButton
                type="primary"
                size="large"
                style="width: 250px"
                :loading="executeLoading"
                @click="handelExecute"
              >
                执行
              </AButton>
            </div>
          </template>
          <template #colSlot_selectDevice="{ field }">
            <FormItem :name="field" label="主机选择" style="margin-left: 90px">
              <a-button type="primary" block style="width: 150px" @click="handleSelectDevice">
                选择设备
              </a-button>
            </FormItem>
          </template>
        </BasicForm>
      </div>
    </PageWrapper>
    <div style="display: flex; flex-direction: column">
      <BasicTable @register="registerTable">
        <template #tableTitle>
          <span class="custom-table-title">{{ title }}</span>
        </template>
        <template #headerCell="{ column }">
          <template v-if="column.key == 'hostname'">
            主机名
            <Tooltip title="主机名复制">
              <CopyOutlined @click="handleCopyHostnames" />
            </Tooltip>
          </template>
          <template v-if="column.key == 'result'">
            结果详情
            <Tooltip title="结果复制">
              <CopyOutlined @click="handleCopyResult" />
            </Tooltip>
          </template>
        </template>
        <template #toolbar>
          <div style="display: flex; align-items: center; justify-content: flex-start"></div>
          <a-form>
            <a-form-item style="margin-top: 5px">
              <a-button
                type="primary"
                size="small"
                @click="handleExportData"
                style="margin-left: 5px"
              >
                {{ exportState.title }}
              </a-button>
              <a-button type="primary" size="small" @click="handleRefresh" style="margin-left: 5px">
                刷新
              </a-button>
              <a-radio-group
                size="small"
                v-model:value="formState.status"
                @change="handleSelectDetail"
                style="margin-left: 5px"
              >
                <a-radio-button value="-1">全部({{ statusCount.all }})</a-radio-button>
                <a-radio-button value="1" style="color: green"
                  >成功({{ statusCount.success }})
                </a-radio-button>
                <a-radio-button value="2" style="color: red"
                  >失败({{ statusCount.failed }})
                </a-radio-button>
                <a-radio-button value="0" style="color: orange"
                  >执行中({{ statusCount.running }})
                </a-radio-button>
              </a-radio-group>
            </a-form-item>
          </a-form>
        </template>
        <template #bodyCell="{ record, column }">
          <template v-if="column.key == 'result'">
            <Popover
              trigger="hover"
              placement="top"
              :mouseLeaveDelay="0.2"
              :overlayStyle="{ maxWidth: '1200px' }"
              :overlayInnerStyle="{ padding: '0px' }"
            >
              <template #content>
                <div style="position: relative">
                  <div class="code-copy-btn" @click="handleCopy(record.result)">
                    <CopyOutlined />
                    <span class="copy-text">复制</span>
                  </div>
                  <div class="result-code-box" v-html="formatResult(record.result)"></div>
                </div>
              </template>
              <span>{{ cutStringLength(record.result, 80) }}</span>
            </Popover>
            <span v-if="record.result && record.result.length > 80">{{
              cutStringLength(record.result, 500).slice(80)
            }}</span>
          </template>
          <template v-else-if="column.key == 'status'">
            <Tag color="orange" v-if="record.status === 0">
              执行中
              <Icon
                icon="svg-spinners:12-dots-scale-rotate"
                width="24"
                height="24"
                style="color: orange"
              />
            </Tag>
            <Tag color="green" v-else-if="record.status === 1">成功</Tag>
            <Tag color="red" v-else>失败</Tag>
          </template>
          <template v-else-if="column.key == 'action'">
            <a @click="handleExecuteDetail(record.hostname)"> 查看历史</a>
          </template>
        </template>
      </BasicTable>
    </div>
    <DeviceSelectModal @register="registerSelectModal" @confirm-devices="handleDeviceConfirm" />
    <ExecuteDetailModal @register="registerDetailModal" />
  </div>
</template>

<script lang="ts" setup>
  import { PageWrapper } from '@/components/Page';
  import { BasicTable, useTable } from '@/components/Table';
  import { BasicForm, useForm } from '@/components/Form';
  import { executeSearchFormSchema } from './data';
  import { executeResultColumns } from '../history/data';
  import { FormItem, Tooltip, message, Tag, notification, Popover } from 'ant-design-vue';
  import { useModal } from '@/components/Modal';
  import DeviceSelectModal from './deviceSelectModal.vue';
  import ExecuteDetailModal from './executeDetailModal.vue';
  import { h, ref, UnwrapRef, reactive, onMounted, onUnmounted, nextTick } from 'vue';
  import Icon from '@/components/Icon/Icon.vue';
  import {
    ScriptExecute,
    GetOperationsTreeData,
    GetExecuteTaskRecordUser,
    GetExecuteRecordDetail,
  } from '@/api/ops/execute';
  import {
    cutStringLength,
    splitByLineAndTrim,
    highlightKeywords,
    handleCopy,
    handleCopyTableRowData,
  } from '@/utils/util';
  import { CopyOutlined } from '@ant-design/icons-vue';
  import { defHttp } from '@/utils/http/axios';

  const [registerSelectModal, { openModal: openSelectModal }] = useModal();
  const [registerDetailModal, { openModal: openDetailModal }] = useModal();

  interface FormState {
    status: '-1' | 1 | 2 | 0;
  }

  const formState: UnwrapRef<FormState> = reactive({
    status: '-1',
  });

  let statusCount = reactive({
    all: 0,
    success: 0,
    failed: 0,
    running: 0,
  });

  // 添加选中的设备数据
  // const selectedDevices = ref([]);
  let currentExecuteSearchFormSchema = ref<string[]>([]);
  let taskId = null;
  let title = ref('');
  let executeLoading = ref(false);
  let namePath = '';
  // 刷新间隔配置（秒）
  const reload_intervals = [1, 3, 5, 6, 6, 6, 10, 10];
  // 在 script setup 中添加
  const selectFormLabels = reactive(new Map<string, string>()); // 存储字段名到标签的映射
  const orderedSelectForms = ref<string[]>([]); // 存储字段名的顺序
  let hostnameCount = ref(0);
  const exportState = reactive({
    exporting: false,
    title: '导出',
  });

  // 轮询控制
  const isPolling = ref(false);
  let pollingTimer: ReturnType<typeof setTimeout> | null = null;
  let pollingIndex = 0;

  function startPolling() {
    if (isPolling.value) return;
    isPolling.value = true;
    pollingIndex = 0;
    scheduleNextPoll();
  }

  function scheduleNextPoll() {
    if (!isPolling.value) return;
    if (pollingIndex >= reload_intervals.length) {
      stopPolling();
      return;
    }

    const delay = reload_intervals[pollingIndex] * 1000; // 转为毫秒
    pollingTimer = setTimeout(async () => {
      await executePoll();
    }, delay);
  }

  async function executePoll() {
    if (!isPolling.value) return;

    // 执行刷新
    await getLastTask();
    // 保持当前分页刷新
    await reload();

    // 准备下一次
    pollingIndex++;
    scheduleNextPoll();
  }

  function stopPolling() {
    isPolling.value = false;
    pollingIndex = 0;
    if (pollingTimer) {
      clearTimeout(pollingTimer);
      pollingTimer = null;
    }
  }

  onUnmounted(() => {
    stopPolling();
  });

  const [
    registerForm,
    {
      setFieldsValue,
      getFieldsValue,
      updateSchema,
      appendSchemaByField,
      removeSchemaByField,
      validateFields,
    },
  ] = useForm({
    labelWidth: 150,
    wrapperCol: {
      span: 18,
    },
    schemas: executeSearchFormSchema(onHostnamesChange),
    actionColOptions: {
      offset: 5,
      span: 16,
    },
    showSubmitButton: false,
    showResetButton: false,
  });

  const [registerTable, { reload, getDataSource }] = useTable({
    // title: '执行记录 ' + stepName || '',
    api: async (params) => {
      if (taskId != null) {
        params.taskId = taskId;
      }
      if (formState.status !== '-1') {
        params.status = formState.status;
      }
      return await GetExecuteRecordDetail(params);
    },
    afterFetch: (data) => {
      // 检查数据中是否有执行中(status === 0)的记录
      if (data && Array.isArray(data)) {
        const hasRunning = data.some((item: any) => item.status === 0);
        if (!hasRunning && isPolling.value) {
          // 如果当前页没有执行中的任务，且正在轮询，可以考虑停止轮询
          // 但为了保险起见（可能其他页还有），可以检查是否是第一次加载完成
          // 这里简化策略：如果当前页所有任务都完成了，就停止轮询当前页面的刷新
          // 或者可以计数，如果连续几次都没有执行中的，就停止
          stopPolling();
        } else if (hasRunning && !isPolling.value) {
          // 如果发现有执行中的，自动开启轮询（针对手动刷新或翻页场景）
          startPolling();
        }
      }
      return data;
    },
    columns: executeResultColumns,
    useSearchForm: false,
    showTableSetting: false,
    bordered: true,
    showIndexColumn: false,
    immediate: false,
    canResize: false,
    actionColumn: {
      width: 70,
      title: '操作',
      dataIndex: 'action',
      fixed: 'right',
    },
    pagination: {
      pageSizeOptions: ['10', '30', '50', '100', '200', '500'],
    },
  });

  async function handleSelectDevice() {
    openSelectModal(true, {});
  }

  // -------- 新增状态管理 --------
  // 记录字段间的依赖关系: ParentField -> ChildFields[]，用于自动清理
  const fieldDependencies = reactive(new Map<string, string[]>());
  // 缓存配置数据的 Map: ConfigName -> ConfigObject
  let configMap = new Map<string, any>();
  // 记录所有顶级动态字段名
  const rootFields = ref<string[]>([]);

  // 获取对应主机业务的功能树
  async function onHostnamesChange() {
    const hostnames = splitByLineAndTrim(getFieldsValue().hostnames);
    const res = await GetOperationsTreeData({ hostnames });
    if (res.code != 200) {
      const errMsg = res.msg;
      notification.error({
        message: '',
        description: h('div', { style: 'white-space: pre-wrap' }, errMsg),
        duration: 5,
      });
    }

    if (hostnames && hostnames.length > 0) {
      await addFieldHint('hostnames', '已选择' + hostnames.length + '台设备');
    }
    const ops = processTreeData(res.data);

    // 检查当前选中的功能是否存在于新的 ops 中
    const currentFunctionOptions = getFieldsValue().functionOptions;
    if (currentFunctionOptions && currentFunctionOptions.length > 0) {
      const isValid = checkSelectionExists(ops, currentFunctionOptions);
      if (!isValid) {
        // 清空选择
        await setFieldsValue({ functionOptions: [] });
        // 清空动态字段
        if (rootFields.value.length > 0) {
          for (const field of rootFields.value) {
            await clearChildFields(field);
            await removeSchemaByField(field);
          }
          rootFields.value = [];
        }
        // 清空提示
        await removeFieldHint('functionOptions');
      }
    }

    if (hostnames?.length == 0) {
      await removeFieldHint('hostnames');
      await removeFieldHint('functionOptions');
      // 重置functionOptions字段
      await setFieldsValue({
        functionOptions: [],
      });
      return;
    }
    await updateSchema({
      field: 'functionOptions',
      componentProps: {
        options: ops,
        fieldNames: {
          label: 'label',
          value: 'value',
          id: 'id',
          params: 'params',
          popupStyle: {
            maxHeight: '1000px', // 设置下拉菜单最大高度
            overflowY: 'auto', // 启用垂直滚动条
          },
        },
        onChange: async (values, selectedOptions) => {
          // 清理所有之前的动态字段（重置表单）
          if (rootFields.value.length > 0) {
            for (const field of rootFields.value) {
              await clearChildFields(field);
              await removeSchemaByField(field);
            }
            rootFields.value = [];
          }
          // 重置累积数据，避免多次选择时数据重复
          selectFormLabels.clear();
          orderedSelectForms.value = [];

          if (selectedOptions[1]?.scriptPath && selectedOptions[1]?.scriptName) {
            const msg =
              '脚本路径: /etc/zabbix/scripts/' +
              selectedOptions[1]?.scriptPath +
              '/agent/' +
              selectedOptions[1]?.scriptName;
            await addFieldHint('functionOptions', msg);
          }
          // 重置 namePath，避免多次选择时累加重复
          namePath = '';
          selectedOptions?.map((item) => {
            namePath += item.name ? item.name + '/' : item.label ? item.label + '/' : '';
          });
          // 获取到参数配置
          // const params = selectedOptions.params || [];

          // selectedOptions 是一个数组，包含从顶层到底层的所有选中项
          // if (params && params.length > 0) {
          // 移除之前的动态字段
          if (currentExecuteSearchFormSchema.value.length > 0) {
            await removeSchemaByField(currentExecuteSearchFormSchema.value);
          }
          // 获取最内层（最后一个）选中项的所有字段数据
          const deepestLevelOption = selectedOptions[selectedOptions.length - 1];
          // 先移除之前的动态字段，保留基础字段
          await removeSchemaByField('options');
          const params = deepestLevelOption.params;
          if (params && params.length > 0) {
            // 建立索引 Map，方便后续查找
            configMap.clear();
            const referencedFields = new Set<string>();

            params.forEach((item) => {
              if (item && item.name) {
                configMap.set(item.name.trim(), item);
                // 收集被引用的字段
                if (item.options) {
                  item.options.forEach((opt) => {
                    if (opt.params && Array.isArray(opt.params)) {
                      opt.params.forEach((p) => referencedFields.add(p.trim()));
                    }
                  });
                }
              }
            });

            // 找出顶级字段（未被任何选项引用的字段）
            // 保持 params 中的原始顺序
            const topLevelConfigs = params.filter(
              (item) => !referencedFields.has(item.name.trim()),
            );

            if (topLevelConfigs.length > 0) {
              let previousField = 'functionOptions';
              const topLevelFieldNames: string[] = [];

              // 1. 先依次添加所有顶级字段
              for (const config of topLevelConfigs) {
                const schema = createDynamicSchema(config);
                try {
                  await appendSchemaByField(schema, previousField);
                  previousField = config.name;
                  topLevelFieldNames.push(config.name);
                } catch (error) {
                  console.error(`Failed to append top-level schema for ${config.name}`, error);
                }
              }

              // 记录根字段（用于清理），这里记录第一个添加的顶级字段作为根的标记
              rootFields.value = topLevelFieldNames;

              // 2. 再依次处理默认值展开（仅对有子参数的选项）
              for (const config of topLevelConfigs) {
                const isSelect = config.type && config.type.toLowerCase() === 'select';
                if (isSelect && config.default) {
                  const selectedOption = config.options?.find(
                    (opt: any) => String(opt.value) === String(config.default),
                  );
                  if (selectedOption) {
                    selectFormLabels.set(config.name, selectedOption.label);
                    if (!orderedSelectForms.value.includes(config.name)) {
                      orderedSelectForms.value.push(config.name);
                    }
                    await setFieldsValue({ [config.name]: config.default });

                    // 只有当选项带有子参数时才触发联动展开
                    if (
                      selectedOption.params &&
                      Array.isArray(selectedOption.params) &&
                      selectedOption.params.length > 0
                    ) {
                      await handleDynamicFieldChange(config.name, config.default, selectedOption);
                    }
                  }
                }
              }
            }
          }
          // }
        },
      },
    });
  }

  // 2. 辅助函数：创建动态 Schema
  function createDynamicSchema(config: any) {
    const isSelect = config.type.toLowerCase() === 'select';
    let type = '';
    switch (config.type) {
      case 'select':
        type = 'Select';
        break;
      case 'input':
        type = 'Input';
        break;
      case 'inputText':
        type = 'InputTextArea';
        break;
    }
    return {
      field: config.name, // 使用 name 作为字段名
      label: config.name,
      component: type,
      required: config.required,
      defaultValue: config.default || '',
      helpMessage: config.help
        ? config.help.split(/\\n|\n/).map((text) => h('div', text))
        : undefined,
      helpComponentProps: {
        placement: 'topLeft',
      },
      colProps: { span: 24 },
      componentProps: () => {
        return {
          placeholder: config.help || `请选择或者输入${config.name}`,
          options: config.options,
          allowClear: true,
          // 只有 Select 类型需要绑定 onChange 来触发下一级
          onChange: isSelect
            ? async (val, opt) => {
                selectFormLabels.set(config.name, opt.label);
                if (!orderedSelectForms.value.includes(config.name)) {
                  orderedSelectForms.value.push(config.name);
                }
                await handleDynamicFieldChange(config.name, val, opt);
              }
            : undefined,
        };
      },
    };
  }

  // 新增辅助函数：处理默认值的递归展开
  async function processDefaultValueExpansion(config: any) {
    const isSelect = config.type && config.type.toLowerCase() === 'select';
    if (isSelect && config.default) {
      // 增加类型兼容性检查 (转为字符串比较)
      const selectedOption = config.options?.find(
        (opt: any) => String(opt.value) === String(config.default),
      );
      if (selectedOption) {
        selectFormLabels.set(config.name, selectedOption.label);
        if (!orderedSelectForms.value.includes(config.name)) {
          orderedSelectForms.value.push(config.name);
        }
        // 关键修复：显式设置表单值，确保 Form Model 状态同步
        await setFieldsValue({ [config.name]: config.default });

        return await handleDynamicFieldChange(config.name, config.default, selectedOption);
      }
    }
    return config.name; // 如果没有展开，返回自身作为最后一个字段
  }

  // 3. 辅助函数：处理字段变更联动
  async function handleDynamicFieldChange(fieldName: string, value: any, option: any) {
    // A. 清理当前字段产生的所有子字段（递归清理）
    await clearChildFields(fieldName);

    let lastAddedField = fieldName;

    // B. 如果选中了有 params 的选项，添加新的子字段
    // 只要 option 存在且有 params 就应该展开，不应该强制要求 value 为 truthy (因为 value 可能是空字符串)
    if (
      (value !== undefined && value !== null) ||
      (option && option.params && Array.isArray(option.params))
    ) {
      const newChildren: string[] = [];
      const childConfigs: any[] = []; // 暂存配置，用于后续处理默认值

      // 1. 第一阶段：同步添加所有直接子字段的 Schema
      // 这样可以保证子字段本身的顺序严格按照 params 定义顺序排列
      let previousField = fieldName;

      for (const childName of option.params) {
        const safeChildName = childName ? childName.trim() : '';
        const childConfig = configMap.get(safeChildName);

        if (childConfig) {
          const schema = createDynamicSchema(childConfig);
          try {
            await appendSchemaByField(schema, previousField);
            // 更新锚点为当前添加的字段
            previousField = childConfig.name;

            newChildren.push(childConfig.name);
            childConfigs.push(childConfig);
          } catch (error) {
            console.error(
              `Failed to append schema for ${childConfig.name} after ${previousField}`,
              error,
            );
            // 如果添加失败，不更新 previousField，尝试继续在当前锚点后添加下一个
          }
        }
      }

      // 更新最后一个添加的字段，作为返回值
      lastAddedField = previousField;

      // 记录依赖关系
      if (newChildren.length > 0) {
        fieldDependencies.set(fieldName, newChildren);
      }

      // 2. 第二阶段：处理默认值递归展开
      // 在所有子字段都占位完成后，再进行递归展开，避免异步操作打乱字段顺序
      for (const childConfig of childConfigs) {
        // 这里不需要 await 返回值来更新 previousField，因为位置已经固定好了
        // 递归展开会插入到 childConfig 后面，不会影响后续兄弟节点的位置
        await processDefaultValueExpansion(childConfig);
      }
    }
    return lastAddedField;
  }

  // 4. 辅助函数：递归清理子字段
  async function clearChildFields(parentField: string) {
    const children = fieldDependencies.get(parentField);
    if (children) {
      for (const childField of children) {
        // 递归清理子字段的子字段
        await clearChildFields(childField);
        // 移除 Schema
        try {
          await removeSchemaByField(childField);
        } catch (e) {
          // 忽略移除不存在字段的错误
        }
        // 同步清理 orderedSelectForms 和 selectFormLabels 中的记录
        const index = orderedSelectForms.value.indexOf(childField);
        if (index > -1) {
          orderedSelectForms.value.splice(index, 1);
        }
        selectFormLabels.delete(childField);
      }
      fieldDependencies.delete(parentField);
    }
  }

  // 递归处理树形结构数据，为每个节点添加 value 和 label 字段
  function processTreeData(data) {
    return data.map((item) => {
      const processedItem = {
        ...item,
        value: item.id,
        label: item.name,
        children: [],
      };
      if (item.children && item.children.length > 0) {
        item.children.forEach((child) => {
          processedItem.children.push({
            ...child,
            value: child.id,
            label: child.name,
          });
        });
      }
      return processedItem;
    });
  }

  // 检查选中的路径是否存在于新的 options 中
  function checkSelectionExists(options: any[], valuePath: any[]): boolean {
    if (!valuePath || valuePath.length === 0) return true;
    let currentOptions = options;

    for (const val of valuePath) {
      if (!currentOptions || currentOptions.length === 0) return false;
      const found = currentOptions.find((opt) => opt.value === val);
      if (!found) return false;
      currentOptions = found.children;
    }
    return true;
  }

  // 添加处理设备确认选择的函数
  function handleDeviceConfirm(selectedRows: any[]) {
    // 可以在这里更新表单数据，比如将选中的主机名设置到表单中
    const hostnames = selectedRows.map((device) => device.hostname).join('\n');
    setFieldsValue({
      hostnames: hostnames,
    });
  }

  // 辅助函数：按顺序获取动态字段的值
  function getOrderedDynamicValues() {
    const values = getFieldsValue();
    const orderedValues: any[] = [];
    // 从所有根字段开始递归查找
    if (rootFields.value.length > 0) {
      for (const field of rootFields.value) {
        traverseFields(field, values, orderedValues);
      }
    }
    return orderedValues;
  }

  // 递归遍历字段依赖
  function traverseFields(currentField: string, allValues: any, result: any[]) {
    // 1. 获取当前字段的值
    const val = allValues[currentField];

    // 2. 将值添加到结果列表 (如果值存在)
    if (val !== undefined && val !== null && val !== '') {
      result.push(val);
    }

    // 3. 查找当前字段的子字段
    const children = fieldDependencies.get(currentField);
    if (children && children.length > 0) {
      // 按照子字段生成的顺序遍历
      for (const childField of children) {
        traverseFields(childField, allValues, result);
      }
    }
  }

  // 执行
  async function handelExecute() {
    // 参数验证
    await validateFields();
    const hostnames = splitByLineAndTrim(getFieldsValue().hostnames);
    if (hostnames?.length == 0) {
      message.error('请选择主机');
      return;
    }
    hostnameCount.value = hostnames?.length || 0;
    // 获取按顺序排列的动态参数值
    // 假设第一个动态字段是 scriptConfigId (业务操作/操作)，后续是 params
    let stepName = '';
    orderedSelectForms.value.forEach((item) => {
      const label = selectFormLabels.get(item);
      stepName += `${label}/`;
    });
    const orderedParams = getOrderedDynamicValues();
    const body = {
      name: namePath + stepName,
      hostnames: hostnames,
      // businessId: getFieldsValue().functionOptions[0],
      scriptConfigId: getFieldsValue().functionOptions[1],
      params: orderedParams || [],
    };
    try {
      executeLoading.value = true;
      const res = await ScriptExecute(body);
      taskId = res.taskId;
      if (taskId) {
        // 获取执行结果
        await reload();
        startPolling();
      }
      message.success('执行成功');
      executeLoading.value = false;
      await getLastTask();
    } finally {
      executeLoading.value = false;
    }
  }

  async function handleSelectDetail() {
    await reload();
  }

  async function handleExecuteDetail(hostname: string) {
    openDetailModal(true, {
      hostname: hostname,
    });
  }

  async function getLastTask() {
    // 获取最新的任务id
    const res = await GetExecuteTaskRecordUser({ pageSize: 1 });
    if (res.items) {
      taskId = res.items[0].id;
      title.value = '任务' + `：${res.items[0].name}` || '';
      statusCount.all = res.items[0].totalCount || 0;
      statusCount.success = res.items[0].successCount || 0;
      statusCount.failed = res.items[0].failedCount || 0;
      statusCount.running = statusCount.all - statusCount.success - statusCount.failed;
    }
  }

  // 在特定字段后添加提示信息
  async function addFieldHint(fieldName: string, text: string) {
    await updateSchema({
      field: fieldName,
      itemProps: {
        extra: h(
          'span',
          {
            style: { color: 'green', fontSize: '14px', margin: '0 0 0 0' },
          },
          text,
        ),
      },
    });
  }

  async function removeFieldHint(fieldName: string) {
    await updateSchema({
      field: fieldName,
      itemProps: {
        extra: null,
      },
    });
  }

  async function handleRefresh() {
    await getLastTask();
  }

  // 复制主机名到剪贴板
  async function handleCopyHostnames() {
    const dataSource = getDataSource();
    await handleCopyTableRowData(dataSource, (item) => item.hostname, '主机名');
  }

  // 复制结果到剪贴板
  async function handleCopyResult() {
    const dataSource = getDataSource();
    await handleCopyTableRowData(dataSource, (item) => item.result, '结果');
  }

  // 格式化结果文本，针对特定格式进行换行处理
  function formatResult(text: string): string {
    let html = highlightKeywords(text);
    // 针对特定格式 ", result:" 后面添加换行
    // 匹配 pattern: , result: 后面允许有空格
    html = html.replace(/(, result:)\s*/gi, '$1\n');
    return html;
  }

  // 数据导出
  function handleExportData() {
    (async function () {
      exportState.exporting = true;
      exportState.title = '导出中...';
      try {
        await DoExportExeResultDetail({ taskId: taskId });
      } catch (error) {
        notification.error({
          message: '导出失败',
          description: error.message,
        });
        exportState.exporting = false;
        exportState.title = '导出';
      }
    })();
  }

  async function DoExportExeResultDetail(value: Recordable) {
    const res = await defHttp.post(
      {
        url: '/v1/ops/execute/tasks/result/export',
        responseType: 'blob',
        data: value,
        timeout: 10 * 60 * 1000,
      },
      { isReturnNativeResponse: true },
    );
    try {
      if (!res.headers['content-type'].includes('application/octet-stream')) {
        // 将 Blob 转换为 JSON
        const reader = new FileReader();
        reader.onload = () => {
          const jsonResponse = JSON.parse(reader.result as any);
          notification.error({
            message: '导出失败',
            description: jsonResponse.msg || '未知错误',
            duration: null,
          });
        };
        reader.readAsText(res.data);
        return;
      }
      const blob = new Blob([res.data], { type: res.headers['content-type'] });
      // 创建新的URL并指向File对象或者Blob对象的地址
      const blobURL = window.URL.createObjectURL(blob);
      // 创建a标签，用于跳转至下载链接
      const tempLink = document.createElement('a');
      tempLink.style.display = 'none';
      tempLink.href = blobURL;
      const contentDisposition =
        res.headers['content-disposition'] || `attachment;filename=hdd_device_info.csv`;
      const filename = contentDisposition.split(';')[1].split('=')[1].split("''")[1];
      tempLink.setAttribute('download', filename);
      // 兼容：某些浏览器不支持HTML5的download属性
      if (typeof tempLink.download === 'undefined') {
        tempLink.setAttribute('target', '_blank');
      }
      // 挂载a标签
      document.body.appendChild(tempLink);
      tempLink.click();
      document.body.removeChild(tempLink);
      // 释放blob URL地址
      window.URL.revokeObjectURL(blobURL);
      notification.success({
        message: '导出成功',
        description: '文件名：' + filename,
        duration: null,
      });
    } finally {
      nextTick(() => {
        exportState.exporting = false;
        exportState.title = '导出';
      });
    }
  }

  onMounted(async () => {
    await getLastTask();
    await reload();
  });
</script>
<style scoped lang="less">
  .tooltip-content {
    max-height: 490px;
    overflow: auto;
    word-break: break-all;
    white-space: pre-wrap;
  }

  .result-code-box {
    max-width: 1150px;
    max-height: 490px;
    padding: 20px;
    padding-top: 32px; // 增加顶部内边距，防止第一行文字被复制按钮遮挡
    overflow: auto;
    border-radius: 6px;
    //background-color: #282c34; // 深色背景，像 IDE
    //color: #abb2bf; // 浅色文字
    background-color: #0b0b0b; // 深色背景，像 IDE
    color: #c6c1c1; // 浅色文字
    font-family: Consolas, Monaco, 'Andale Mono', 'Ubuntu Mono', monospace;
    font-size: 13px;
    line-height: 1.5;
    word-break: break-all;
    white-space: pre-wrap;
    user-select: text;
  }

  .result-trigger {
    cursor: text;
  }

  .code-copy-btn {
    display: flex;
    position: absolute;
    z-index: 10;
    top: 5px;
    right: 10px;
    align-items: center;
    padding: 4px 8px;
    transition: all 0.3s;
    border-radius: 4px;
    opacity: 0.6;
    background-color: rgb(255 255 255 / 10%);
    color: #d1d2da;
    cursor: pointer;

    &:hover {
      opacity: 1;
      background-color: rgb(255 255 255 / 20%);
      color: #fff;
    }

    .anticon {
      font-size: 14px;
    }

    .copy-text {
      margin-left: 4px;
      font-size: 12px;
    }
  }

  // 自定义表格标题样式 - 美观的灰色设计
  .custom-table-title {
    display: inline-flex;
    position: relative;
    align-items: center;
    padding-left: 8px;
    color: #666;
    font-size: 15px;
    font-weight: 500;
    line-height: 24px;

    // ops-batch-exec
    // 添加左侧装饰条
    &::before {
      content: '';
      position: absolute;
      top: 50%;
      left: 0;
      width: 3px;
      height: 16px;
      transform: translateY(-50%);
      border-radius: 2px;
      background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    }
  }
</style>
<style lang="less">
  .function-options-cascader-popup .ant-cascader-menu {
    height: auto !important;
    max-height: 400px !important;
  }
</style>

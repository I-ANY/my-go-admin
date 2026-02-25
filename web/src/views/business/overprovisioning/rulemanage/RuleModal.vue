<template>
  <BasicModal
    v-bind="$attrs"
    @register="registerModal"
    :title="isUpdate ? '编辑规则' : '新建规则'"
    @ok="handleSubmit"
    :width="1000"
  >
    <BasicForm @register="registerForm">
      <template #storageBwRatioGroup>
        <div class="flex items-center">
          <InputNumber
            v-model:value="storageRatio"
            placeholder="存储(TB)"
            :min="0.01"
            :precision="2"
            :step="0.01"
            style="width: 100px"
          />
          <span class="mx-2">:</span>
          <InputNumber
            v-model:value="BwRatio"
            placeholder="带宽(Gbps)"
            :min="0.01"
            :precision="2"
            :step="0.01"
            style="width: 100px"
          />
        </div>
      </template>
      <template #planconfGroup>
        <div class="flex flex-col gap-3">
          <div
            v-for="(item, index) in planconfList"
            :key="index"
            class="flex flex-col gap-2 p-3 border border-gray-200 rounded"
          >
            <div class="flex items-center justify-between mb-2">
              <span class="font-medium">规则配置 {{ index + 1 }}</span>
              <Button
                type="text"
                danger
                size="small"
                @click="removePlanConf(index)"
                :disabled="planconfList.length === 1"
              >
                <DeleteOutlined />
                删除
              </Button>
            </div>
            <div class="flex gap-2 items-center">
              <span style="width: 80px">表达式类型:</span>
              <FormItemRest>
                <Select
                  v-model:value="item.cond_type"
                  placeholder="请选择"
                  style="flex: 1"
                  :options="[
                    { label: '提交带宽(Gbps)', value: 1 },
                    { label: '规划带宽(Gbps)', value: 2 },
                    { label: '磁盘/带宽比(TB:Gbps)', value: 3 },
                    { label: '磁盘类型(SSD/HDD)', value: 4 },
                    { label: '系统盘大小(GB)', value: 5 },
                    { label: 'SSD大小(TB)', value: 6 },
                    { label: 'HDD大小(TB)', value: 7 },
                    { label: '业务盘总大小(TB)', value: 8 },
                    { label: '内存总大小(GB)', value: 9 },
                  ]"
                />
              </FormItemRest>
            </div>
            <div class="flex gap-2 items-center">
              <span style="width: 80px">条件运算符:</span>
              <FormItemRest>
                <Select
                  v-model:value="item.cond.cond"
                  placeholder="请选择运算符"
                  style="width: 150px"
                  :options="[
                    { label: '= (等于)', value: '=' },
                    { label: '!= (不等于)', value: '!=' },
                    { label: '> (大于)', value: '>' },
                    { label: '>= (大于等于)', value: '>=' },
                    { label: '< (小于)', value: '<' },
                    { label: '<= (小于等于)', value: '<=' },
                  ]"
                />
              </FormItemRest>
              <span style="width: 30px; text-align: center">值:</span>
              <!-- 默认值为 0，当表达式类型为标准配置时，值为0 -->
              <FormItemRest>
                <Input
                  v-model:value="item.cond.value"
                  placeholder="请输入值 (数字/字符串)"
                  style="flex: 1"
                />
              </FormItemRest>
            </div>
            <div class="flex gap-2 items-center" v-if="index !== 0">
              <span style="width: 80px">逻辑关系:</span>
              <FormItemRest>
                <Select
                  v-model:value="item.cond.iftype"
                  placeholder="请选择"
                  style="flex: 1"
                  :options="[
                    { label: 'AND (所有条件都满足)', value: 'and' },
                    { label: 'OR (至少一个条件满足)', value: 'or' },
                  ]"
                />
              </FormItemRest>
            </div>
            <div class="flex gap-2 items-center" v-else>
              <span style="width: 80px">逻辑关系:</span>
              <FormItemRest>
                <Input value="AND (默认)" disabled style="flex: 1" />
              </FormItemRest>
            </div>
          </div>
          <Button type="dashed" block @click="addPlanConf">
            <PlusOutlined />
            添加规则配置
          </Button>
        </div>
      </template>
    </BasicForm>
  </BasicModal>
</template>

<script lang="ts" setup>
  import { message, InputNumber, Select, Button, Input, Form } from 'ant-design-vue';
  import { PlusOutlined, DeleteOutlined } from '@ant-design/icons-vue';
  import { BasicForm, useForm } from '@/components/Form';
  import { BasicModal, useModalInner } from '@/components/Modal';
  import { ruleFormSchema } from './data';
  import { createRule, updateRule } from '@/api/business/overprovisioning';
  import { ref } from 'vue';

  defineOptions({ name: 'RuleModal' });

  const FormItemRest = Form.ItemRest;

  interface CondDetail {
    cond: string; // 条件运算符
    value: string; // 值 (可以是 int、string 或 array)
    iftype: string; // 逻辑关系
  }

  interface PlanConfItem {
    cond_type: number;
    cond: CondDetail;
  }

  const emit = defineEmits(['success', 'register']);
  const storageRatio = ref<any>(null);
  const BwRatio = ref<any>(null);
  const planconfList = ref<PlanConfItem[]>([]);
  const isUpdate = ref(false);
  const recordId = ref<number | null>(null);
  const [registerForm, { resetFields, validate, setFieldsValue, getFieldsValue }] = useForm({
    labelWidth: 120,
    baseColProps: { span: 24 },
    schemas: ruleFormSchema,
    showActionButtonGroup: false,
  });

  const addPlanConf = () => {
    planconfList.value.push({
      cond_type: 0,
      cond: {
        cond: '=',
        value: '0',
        iftype: 'and',
      },
    });
  };

  const removePlanConf = (index: number) => {
    if (planconfList.value.length > 1) {
      planconfList.value.splice(index, 1);
    }
  };

  const [registerModal, { setModalProps, closeModal }] = useModalInner(async (data) => {
    await resetFields();
    storageRatio.value = null;
    BwRatio.value = null;
    // 新建时不再默认添加一条规则配置，交由用户手动添加
    planconfList.value = [];
    isUpdate.value = !!data?.isUpdate;
    recordId.value = null;

    setModalProps({ confirmLoading: false });

    // 如果是编辑模式，填充表单数据
    if (data?.isUpdate && data?.record) {
      recordId.value = data.record.id;
      console.log('编辑记录数据:', data.record);

      await setFieldsValue({
        ruleName: data.record.ruleName,
        mem: data.record.mem,
        ssdSize: data.record.ssdSize,
        hddSize: data.record.hddSize,
        sysDiskSize: data.record.sysdisksize || data.record.sysDiskSize || data.record.SysDiskSize,
        businessIds:
          data.record.businesses && Array.isArray(data.record.businesses)
            ? data.record.businesses.map((b) => b.businessId)
            : [],
        ruleStatus: data.record.ruleStatus,
        ruleDescription: data.record.ruleDescription,
        isProvinceScheduling: data.record.isProvinceScheduling ?? 0,
        enableMemoryCheck: data.record.enableMemoryCheck ?? true,
        enableSysDiskCheck: data.record.enableSysDiskCheck ?? true,
        enableSsdCheck: data.record.enableSsdCheck ?? true,
        enableHddCheck: data.record.enableHddCheck ?? true,
        enableDataDiskCheck: data.record.enableDataDiskCheck ?? true,
        enableStorageBwRatioCheck: data.record.enableStorageBwRatioCheck ?? true,
        storageBwRatioDiskType: data.record.storageBwRatioDiskType ?? 0,
      });

      // 解析 storageBwRatio 字符串，如 "1:2" -> storageRatio=1, BwRatio=2
      if (data.record.storageBwRatio) {
        const [storage, bw] = data.record.storageBwRatio.split(':');
        storageRatio.value = parseFloat(storage);
        BwRatio.value = parseFloat(bw);
      }

      // 处理 planconf 字段，支持多种可能的字段名
      const planconfData = data.record.planconf || data.record.planConf || data.record.PlanConf;
      console.log('planconf 数据:', planconfData);

      if (planconfData && Array.isArray(planconfData)) {
        planconfList.value = planconfData.map((item: any) => ({
          cond_type: item.cond_type || item.condType || item.CondType || 0,
          cond: {
            cond: item.cond?.cond || item.cond?.Cond || '=',
            value: item.cond?.value || item.cond?.Value || '0',
            iftype: item.cond?.iftype || item.cond?.IfType || item.cond?.ifType || 'and',
          },
        }));
      }
    }
  });

  async function handleSubmit() {
    try {
      await validate();
      const data = getFieldsValue();
      // 验证存储和带宽必填且大于0
      // if (!storageRatio.value || storageRatio.value <= 0 || !BwRatio.value || BwRatio.value <= 0) {
      //   message.error('存储和带宽必须大于 0');
      //   return;
      // }

      // 验证规划配置至少有一条且必填字段已填写
      // if (planconfList.value.length === 0) {
      //   message.error('请至少添加一条超配规则');
      //   return;
      // }

      for (let i = 0; i < planconfList.value.length; i++) {
        const item = planconfList.value[i];
        if (!item.cond.cond) {
          message.error(`规则配置 ${i + 1} 的条件运算符不能为空`);
          return;
        }
        // 支持值为 0：仅当为空、null、undefined 或纯空白时才认为未填写
        const condValue = item.cond.value;
        if (condValue === null || condValue === undefined || String(condValue).trim() === '') {
          message.error(`规则配置 ${i + 1} 的值不能为空`);
          return;
        }
        if (!item.cond.iftype) {
          message.error(`规则配置 ${i + 1} 的逻辑关系不能为空`);
          return;
        }
      }

      // 处理业务归属：当选择“所有业务(0)”时，只保留 0
      if (Array.isArray(data.businessIds) && data.businessIds.includes(0)) {
        data.businessIds = [0];
      }

      // 组合存储带宽比字符串，如 "1:2" 或 "2:2.5"
      // 允许为空：当任一为空时，不传 storageBwRatio 字段
      if (
        storageRatio.value !== null &&
        storageRatio.value !== undefined &&
        storageRatio.value !== '' &&
        BwRatio.value !== null &&
        BwRatio.value !== undefined &&
        BwRatio.value !== ''
      ) {
        data.storageBwRatio = `${storageRatio.value}:${BwRatio.value}`;
      } else {
        delete data.storageBwRatio;
      }
      // 转换 ruleStatus 为 int 类型 ('1'/'0' -> 1/0)
      data.ruleStatus = Number(data.ruleStatus);
      // 使用 API 文档中的字段名 planConf
      data.planConf = planconfList.value;

      setModalProps({ confirmLoading: true });

      // 根据是否是更新模式调用不同的接口
      if (isUpdate.value) {
        // Remove ID from data if it exists to avoid redundancy in body if not needed,
        // but often it's ignored or useful. The key change is passing ID as first arg.
        await updateRule(recordId.value!, data);
        console.log(data);
        message.success('规则更新成功');
      } else {
        await createRule(data);
        message.success('规则创建成功');
      }

      emit('success');
      closeModal();
    } catch (error) {
      message.error(`规则${isUpdate.value ? '更新' : '创建'}失败: ${error.message || '未知错误'}`);
    } finally {
      setModalProps({ confirmLoading: false });
    }
  }
</script>

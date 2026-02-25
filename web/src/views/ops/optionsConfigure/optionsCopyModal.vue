<template>
  <BasicModal v-bind="$attrs" @register="registerModal" :title="getTitle" @ok="handleSubmit">
    <div v-if="sourceType === 'group'" style="margin-bottom: 12px; color: #faad14">
      将拷贝该操作组下的所有功能到目标业务与操作组中，请确认。
    </div>
    <BasicForm @register="registerForm" />
  </BasicModal>
</template>

<script setup lang="ts">
  import { computed, defineOptions, reactive, ref, unref } from 'vue';
  import { BasicModal, useModalInner } from '@/components/Modal';
  import { BasicForm, useForm } from '@/components/Form';
  import { copyFormSchema } from './data';
  import { GetScriptTasks, CopyScriptConfig } from '@/api/ops/execute';
  import { GetAuthedBiz } from '@/api/business/biz';
  import { message } from 'ant-design-vue';

  defineOptions({
    name: 'OptionsCopyModal',
  });

  const emit = defineEmits(['success', 'register']);

  const sourceId = ref<number | null>(null);
  const sourceType = ref<'function' | 'group' | ''>('');

  const getTitle = computed(() =>
    unref(sourceType) === 'group' ? '拷贝操作组下所有功能' : '拷贝功能',
  );

  const [registerForm, { resetFields, getFieldsValue, setFieldsValue, validate, updateSchema }] =
    useForm({
      labelWidth: 80,
      baseColProps: { span: 24 },
      schemas: copyFormSchema,
      showActionButtonGroup: false,
    });

  async function handleBizChange(value: number) {
    const res = await GetScriptTasks({ businessId: value, type: 'group' });
    if (res.items) {
      await updateSchema([
        {
          field: 'targetParentId',
          componentProps: {
            options: res.items,
            fieldNames: {
              label: 'name',
              value: 'id',
            },
            placeholder: '请选择目标操作组',
          },
        },
      ]);
    }

    // 目标业务为公共时，若未填写脚本目录，则默认填充 general
    if (value === 0) {
      const current = getFieldsValue();
      if (!current.scriptPath) {
        await setFieldsValue({ scriptPath: 'general' });
      }
    }
  }

  const [registerModal, { setModalProps, closeModal }] = useModalInner(async (data) => {
    setModalProps({ width: 500, minHeight: 60 });
    await resetFields();

    sourceId.value = data.sourceId;
    sourceType.value = data.sourceType || 'function';

    // 根据来源类型动态控制 targetParentId 是否必填
    if (sourceType.value === 'group') {
      await updateSchema([
        {
          field: 'targetParentId',
          required: false,
        },
      ]);
    } else {
      await updateSchema([
        {
          field: 'targetParentId',
          required: true,
        },
      ]);
    }

    const categoryOptions = reactive<any[]>([
      {
        name: '公共',
        id: 0,
      },
    ]);
    const bizData = await GetAuthedBiz({});
    if (bizData.categories) {
      categoryOptions.push(...bizData.categories);
    }

    await updateSchema([
      {
        field: 'targetBizId',
        componentProps: {
          options: categoryOptions,
          fieldNames: {
            label: 'name',
            value: 'id',
          },
          placeholder: '请选择目标业务',
          onChange: (value: number) => {
            handleBizChange(value);
          },
        },
      },
    ]);
  });

  async function handleSubmit() {
    await validate();
    const values = getFieldsValue() as Recordable;

    const params: Recordable = {
      sourceId: sourceId.value,
      targetBizId: values.targetBizId,
      targetParentId: values.targetParentId,
    };
    if (values.scriptPath) {
      params.scriptPath = values.scriptPath;
    }

    await CopyScriptConfig(params);
    message.success('拷贝成功');
    emit('success');
    closeModal();
  }
</script>

<style scoped lang="less"></style>

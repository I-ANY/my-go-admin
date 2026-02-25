<template>
  <BasicModal
    v-bind="$attrs"
    @register="registerModal"
    :title="getTitle"
    @ok="handleSubmit"
    width="50%"
  >
    <BasicForm @register="registerForm" class="assessment-form" />
  </BasicModal>
</template>

<script lang="ts" setup>
  import { ref, computed, unref } from 'vue';
  import { BasicModal, useModalInner } from '@/components/Modal';
  import { BasicForm, useForm } from '@/components/Form';
  import { formSchema } from './data';
  import { updateAssessmentRule } from '@/api/business/a';
  import { message } from 'ant-design-vue';

  defineOptions({ name: 'AssessmentModal' });

  const emit = defineEmits(['success', 'register']);

  const isUpdate = ref(true);
  const rowId = ref('');

  const [registerForm, { setFieldsValue, clearValidate, validate, resetFields }] = useForm({
    labelWidth: 120,
    schemas: formSchema,
    showActionButtonGroup: false,
    baseColProps: { lg: 24, md: 24 },
    fieldMapToTime: [],
  });

  const [registerModal, { setModalProps, closeModal }] = useModalInner(async (data) => {
    resetFields();
    clearValidate();
    setModalProps({ confirmLoading: false });
    isUpdate.value = !!data?.isUpdate;

    if (unref(isUpdate)) {
      rowId.value = data.record.id;
      // 处理考核标准JSON
      let formData = { ...data.record };
      if (data.record.assessmentStandard) {
        // try {
        const standard = data.record.assessmentStandard;
        formData = {
          ...formData,
          utilizationRateThreshold: standard.utilizationRateThreshold * 100, // 转换为百分比显示
          nightPeakPointsThreshold: standard.nightPeakPointsThreshold,
          customRules: standard.customRules,
        };
        // } catch {
        //   // JSON解析失败时使用默认值
        //   formData = {
        //     ...formData,
        //     utilizationRateThreshold: 90,
        //     nightPeakPointsThreshold: 36,
        //   };
        // }
      }
      setFieldsValue(formData);
    } else {
      // 新增时设置默认值
      setFieldsValue({
        utilizationRateThreshold: 90,
        nightPeakPointsThreshold: 36,
      });
    }
  });

  const getTitle = computed(() => (!unref(isUpdate) ? '新增考核规则' : '编辑考核规则'));

  async function handleSubmit() {
    try {
      const values = await validate();
      setModalProps({ confirmLoading: true });

      // 处理考核标准
      const assessmentStandard = {
        utilizationRateThreshold: values.utilizationRateThreshold / 100, // 转换为小数
        nightPeakPointsThreshold: values.nightPeakPointsThreshold,
      };

      // 验证考核标准格式
      const submitData = {
        assessmentStandard: assessmentStandard,
      };

      if (unref(isUpdate)) {
        await updateAssessmentRule(Number(rowId.value), submitData);
        message.success('编辑成功');
      } else {
        // 新增功能暂未实现
        message.info('新增功能暂未实现');
        return;
      }

      closeModal();
      emit('success');
    } catch (error) {
      console.error('操作失败:', error);
      message.error('操作失败');
    } finally {
      setModalProps({ confirmLoading: false });
    }
  }
</script>

<style lang="less" scoped>
  .assessment-form {
    /deep/ .ant-form-item-label {
      align-items: flex-start !important;
      min-width: 120px !important;
      padding-top: 2px !important;
      line-height: 1.4 !important;
      white-space: normal !important;
    }
  }
</style>

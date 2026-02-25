<template>
    <BasicModal v-bind="$attrs" @register="register" @ok="handleSubmit">
        <BasicForm @register="registerForm" />
    </BasicModal>
</template>
<script lang="ts" setup>
import { BasicModal, useModalInner } from '@/components/Modal';
import { BasicForm, useForm } from '@/components/Form';
import { getFormSchema } from './data';
import { reactive } from 'vue';
import { createEncryptionConfiguration, updateEncryptionConfiguration } from '@/linuxApi/config';

const emits = defineEmits(['success', 'register']);
const allData = reactive({
    isUpdate: false,
    record: {} as any,
});
const [registerForm, { resetFields, setFieldsValue, validate, updateSchema }] = useForm({
    labelWidth: 120,
    baseColProps: { span: 24 },
    schemas: getFormSchema(),
    showActionButtonGroup: false,
});
const [register, { closeModal, setModalProps }] = useModalInner(async (data) => {
    resetFields();
    allData.isUpdate = !!data.isUpdate;
    setModalProps({
        confirmLoading: false,
        destroyOnClose: true,
        title: allData.isUpdate ? '编辑加密脚本配置' : '新增加密脚本配置',
        width: 800,
        height: 500,
    });
    if (allData.isUpdate) {
        await updateSchema([
            {
                field: 'scheduleNode',
                componentProps: {
                    disabled: true,
                },
            },
            {
                field: 'basename',
                componentProps: {
                    disabled: true,
                },
            },
            {
                field: 'os_type',
                componentProps: {
                    disabled: true,
                },
            },
        ]);
        allData.record = data.record;
        setFieldsValue({
            ...data.record,
        });
    } else {
        await updateSchema([
            {
                field: 'scheduleNode',
                componentProps: {
                    disabled: false,
                },
            },
            {
                field: 'basename',
                componentProps: {
                    disabled: false,
                },
            },
            {
                field: 'os_type',
                componentProps: {
                    disabled: false,
                },
            },
        ]);
    }
});
async function handleSubmit() {
    let data = await validate();
    if (allData.isUpdate) {
        await updateEncryptionConfiguration({ basename: allData.record.basename, os_type: allData.record.os_type, version: allData.record.tar_version }, { new_version: data.tar_version, tar_url: data.tar_url, tar_name: data.tar_name, operator: JSON.parse(localStorage.getItem('userInfo') || '{}').nickname });
    } else {
        await createEncryptionConfiguration({ ...data, operator: JSON.parse(localStorage.getItem('userInfo') || '{}').nickname });
    }
    emits("success", !allData.isUpdate);

    closeModal();
}
</script>
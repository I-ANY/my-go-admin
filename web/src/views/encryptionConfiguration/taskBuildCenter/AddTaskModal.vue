<template>
    <BasicModal v-bind="$attrs" @register="register" @ok="handleSubmit" ok-text="开始构建" cancel-text="取消">
        <BasicForm @register="registerForm" />
    </BasicModal>
</template>
<script lang="ts" setup>
import { BasicModal, useModalInner } from '@/components/Modal';
import { BasicForm, useForm } from '@/components/Form';
import { getFormSchema } from './data';
import { reactive, ref } from 'vue';
import { createTaskBuildCenter, getAllBasenameList, getAllOsTypeList, getAllTarVersionList } from '@/linuxApi/config';
import { message } from 'ant-design-vue';

const emits = defineEmits(['success', 'register']);
const allData = reactive({
    isUpdate: false,
    record: {} as any,
});

const basenameList = ref<any[]>([]);
const osTypeList = ref<any[]>([]);
const tarVersionList = ref<any[]>([]);

const [registerForm, { resetFields, setFieldsValue, validate, updateSchema, getFieldsValue }] = useForm({
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
        title: allData.isUpdate ? '编辑任务' : '新增任务',
        width: 800,
        height: 500,
    });

    // 初始化基础配置名称列表
    await initBasenameList();
});

// 初始化基础配置名称列表
async function initBasenameList() {
    try {
        const res = await getAllBasenameList();
        basenameList.value = res.map((item: any) => ({ label: item, value: item }));

        // 更新 basename 字段配置并添加 onChange 事件
        await updateSchema({
            field: 'basename',
            componentProps: {
                options: basenameList.value,
                onChange: async (value: string) => {
                    // 清空操作系统和版本号
                    setFieldsValue({ os_type: undefined, tar_version: undefined });
                    osTypeList.value = [];
                    tarVersionList.value = [];

                    if (value) {
                        // 获取操作系统列表
                        const osTypes = await getAllOsTypeList(value);
                        osTypeList.value = osTypes.map((item: any) => ({ label: item, value: item }));
                        await updateSchema({
                            field: 'os_type',
                            componentProps: {
                                options: osTypeList.value,
                            },
                        });
                    } else {
                        // 清空操作系统和版本号选项
                        await updateSchema([
                            {
                                field: 'os_type',
                                componentProps: {
                                    options: [],
                                },
                            },
                            {
                                field: 'tar_version',
                                componentProps: {
                                    options: [],
                                },
                            },
                        ]);
                    }
                },
            },
        });

        // 更新 os_type 字段配置并添加 onChange 事件
        await updateSchema({
            field: 'os_type',
            componentProps: {
                options: osTypeList.value,
                onChange: async (value: string) => {
                    // 清空版本号
                    setFieldsValue({ tar_version: undefined });
                    tarVersionList.value = [];

                    // 使用 getFieldsValue 而不是 validate，避免验证错误
                    const formValues = getFieldsValue();
                    const basename = formValues.basename;

                    if (basename && value) {
                        try {
                            // 获取版本号列表
                            const versions = await getAllTarVersionList(basename, value);
                            tarVersionList.value = versions.map((item: any) => ({ label: item, value: item }));
                            await updateSchema({
                                field: 'tar_version',
                                componentProps: {
                                    options: tarVersionList.value,
                                },
                            });
                        } catch (error) {
                            console.error('获取版本号列表失败:', error);
                            message.error('获取版本号列表失败');
                        }
                    } else {
                        // 清空版本号选项
                        await updateSchema({
                            field: 'tar_version',
                            componentProps: {
                                options: [],
                            },
                        });
                    }
                },
            },
        });
    } catch (error) {
        console.error('初始化基础配置名称列表失败:', error);
    }
}

async function handleSubmit() {
    try {
        const formData = await validate();
        const res = await createTaskBuildCenter({ ...formData, creator: JSON.parse(localStorage.getItem('userInfo') || '{}').nickname });
        if (res.task_id) {
            message.success('任务创建成功');
            emits("success", !allData.isUpdate, { task_id: res.task_id,...formData, creator: JSON.parse(localStorage.getItem('userInfo') || '{}').nickname });
            closeModal();
        } else {
            message.error('任务创建失败');
        }
        // emits("success", !allData.isUpdate, { task_id: res.task_id,...formData, creator: JSON.parse(localStorage.getItem('userInfo') || '{}').nickname });
        // closeModal();
    } catch (error) {
        console.error('提交失败:', error);
    }
}
</script>
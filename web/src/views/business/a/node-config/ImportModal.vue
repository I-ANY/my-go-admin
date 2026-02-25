<template>
  <BasicModal
    v-bind="$attrs"
    @register="registerModal"
    title="批量导入节点配置"
    @ok="handleImport"
    width="800px"
  >
    <div class="import-container">
      <Alert
        type="info"
        message="导入说明"
        description="请下载模板文件，按照模板格式填写数据后再上传。支持Excel格式文件(.xlsx, .xls)。包含节点编号(必填)、计费方式(选填)、是否仅异网节点(选填)、是否考核(选填)、保底带宽(选填)。模板只需要填写节点编号+需要更新的字段，其他不需要更新的字段不填写。"
        show-icon
        class="mb-4"
      />

      <div class="template-download mb-4">
        <a-button @click="downloadTemplate" type="link">
          <template #icon>
            <DownloadOutlined />
          </template>
          下载导入模板
        </a-button>
      </div>

      <BasicForm @register="registerForm" />
    </div>
  </BasicModal>
</template>

<script lang="ts" setup>
  import { ref, computed } from 'vue';
  import { BasicModal, useModalInner } from '@/components/Modal';
  import { BasicForm, useForm } from '@/components/Form';
  import { Alert, message } from 'ant-design-vue';
  import { DownloadOutlined } from '@ant-design/icons-vue';
  import { importNodeConfigs } from '@/api/business/a';
  import * as XLSX from 'xlsx';

  defineOptions({ name: 'ImportModal' });

  const emit = defineEmits(['success', 'register']);

  const billingTypes = ref<string[]>([]);
  const fileData = ref<any[]>([]);

  const [registerForm, { resetFields }] = useForm({
    labelWidth: 100,
    baseColProps: { span: 24 },
    schemas: [
      {
        field: 'file',
        label: '选择文件',
        component: 'Upload',
        required: true,
        componentProps: {
          accept: ['.xlsx', '.xls'],
          api: handleFileUpload, // 直接处理文件的函数
          multiple: false, // 只允许选择一个文件
        } as any,
        itemProps: {
          rules: [
            {
              required: true,
              message: '请选择要上传的文件',
            },
          ],
        },
        helpMessage: '支持Excel格式文件(.xlsx, .xls)',
        colProps: { span: 24 },
      },
    ],
    showActionButtonGroup: false,
  });

  const [registerModal, { setModalProps, closeModal }] = useModalInner(async () => {
    resetFields();
    fileData.value = [];
    setModalProps({ confirmLoading: false });

    // 设置固定的计费方式列表
    billingTypes.value = ['日95', '月95', '买断'];
    console.log('计费方式列表:', billingTypes.value);
  });

  // 处理文件上传的api函数
  async function handleFileUpload({ file }: { file: File }) {
    try {
      if (file instanceof File) {
        await handleFileSelect(file);
        // 返回模拟的上传成功结果
        return {
          data: {
            url: `temp://${file.name}`,
            name: file.name,
          },
        };
      } else {
        console.log('文件类型不正确:', typeof file);
        message.error('文件类型不正确，请选择Excel文件');
        throw new Error('文件类型不正确');
      }
    } catch (error) {
      console.error('文件处理失败:', error);
      throw error;
    }
  }

  // 处理文件选择
  async function handleFileSelect(file: File) {
    try {
      const data = await readExcelFile(file);
      parseExcelData(data);
      console.log('fileData:', fileData.value);
    } catch (error) {
      console.error('读取文件失败:', error);
      message.error('文件读取失败，请检查文件格式');
    }
  }

  // 有效和无效数据
  const validData = computed(() => fileData.value.filter((item) => item.isValid));
  // const invalidData = computed(() => fileData.value.filter((item) => !item.isValid));

  // 读取Excel文件
  function readExcelFile(file: File): Promise<any[]> {
    return new Promise((resolve, reject) => {
      const reader = new FileReader();
      reader.onload = (e) => {
        try {
          const data = new Uint8Array(e.target?.result as ArrayBuffer);
          const workbook = XLSX.read(data, { type: 'array' });
          const worksheet = workbook.Sheets[workbook.SheetNames[0]];
          const jsonData = XLSX.utils.sheet_to_json(worksheet, { header: 1 });
          resolve(jsonData);
        } catch (error) {
          reject(error);
        }
      };
      reader.onerror = reject;
      reader.readAsArrayBuffer(file);
    });
  }

  // 解析Excel数据
  function parseExcelData(data: any[]) {
    if (data.length < 2) {
      message.error('文件数据不足，至少需要包含表头和一行数据');
      return;
    }

    const headers = data[0] as string[];
    const ownerIndex = headers.findIndex((h) => h && h.includes('节点编号'));
    const billingTypeIndex = headers.findIndex((h) => h && h.includes('计费方式'));
    const isExternalOnlyIndex = headers.findIndex((h) => h && h.includes('是否仅异网节点'));
    const isAssessmentIndex = headers.findIndex((h) => h && h.includes('是否考核'));
    // const minBwIndex = headers.findIndex((h) => h && h.includes('保底带宽'));
    // const assessmentRemarkIndex = headers.findIndex((h) => h && h.includes('考核备注'));
    // const planningRemarkIndex = headers.findIndex((h) => h && h.includes('规划备注'));
    // const opsRemarkIndex = headers.findIndex((h) => h && h.includes('运维备注'));

    if (ownerIndex === -1) {
      message.error('文件格式不正确，请确保包含"节点编号"列');
      return;
    }

    const parsedData = data
      .slice(1)
      .map((row: any[]) => {
        const owner = row[ownerIndex]?.toString().trim();
        const billingType = billingTypeIndex >= 0 ? row[billingTypeIndex]?.toString().trim() : '';
        const isExternalOnlyStr =
          isExternalOnlyIndex >= 0 ? row[isExternalOnlyIndex]?.toString().trim() : '';
        const isAssessmentStr =
          isAssessmentIndex >= 0 ? row[isAssessmentIndex]?.toString().trim() : '';
        // const minBwStr = minBwIndex >= 0 ? row[minBwIndex]?.toString().trim() : '';
        // const assessmentRemark =
        //   assessmentRemarkIndex >= 0 ? row[assessmentRemarkIndex]?.toString().trim() : '';
        // const planningRemark =
        //   planningRemarkIndex >= 0 ? row[planningRemarkIndex]?.toString().trim() : '';
        // const opsRemark = opsRemarkIndex >= 0 ? row[opsRemarkIndex]?.toString().trim() : '';

        let isValid = true;
        let errorMessage = '';

        // 验证节点编号（必填）
        if (!owner) {
          isValid = false;
          errorMessage += '节点编号不能为空; ';
        }

        // 验证计费方式（选填）
        let billingTypeValue = undefined;
        if (billingType) {
          const validBillingTypes = ['日95', '月95', '买断'];
          if (!validBillingTypes.includes(billingType)) {
            isValid = false;
            errorMessage += `无效的计费方式"${billingType}"，仅支持：${validBillingTypes.join('、')}; `;
          } else {
            billingTypeValue = billingType;
          }
        }

        // 验证是否仅异网节点（选填）
        let isExternalOnly: boolean | undefined = undefined;
        if (isExternalOnlyStr) {
          if (isExternalOnlyStr === '是') {
            isExternalOnly = true;
          } else if (isExternalOnlyStr === '否') {
            isExternalOnly = false;
          } else {
            isValid = false;
            errorMessage += `无效的是否仅异网节点值"${isExternalOnlyStr}"，应为"是"或"否"; `;
          }
        }

        // 验证是否考核（选填）
        let isAssessment: boolean | undefined = undefined;
        if (isAssessmentStr) {
          if (isAssessmentStr === '是') {
            isAssessment = true;
          } else if (isAssessmentStr === '否') {
            isAssessment = false;
          } else {
            isValid = false;
            errorMessage += `无效的是否考核值"${isAssessmentStr}"，应为"是"或"否"; `;
          }
        }

        // // 验证保底带宽（选填，单位：G）
        // let minBw: number | undefined = undefined;
        // if (minBwStr) {
        //   try {
        //     const value = parseFloat(minBwStr);
        //     if (!isNaN(value) && value >= 0) {
        //       minBw = value * 1000 * 1000 * 1000; // 转换为bps
        //     } else {
        //       isValid = false;
        //       errorMessage += '保底带宽必须是非负数字; ';
        //     }
        //   } catch (e) {
        //     isValid = false;
        //     errorMessage += '保底带宽格式不正确; ';
        //   }
        // }

        return {
          owner,
          billingType: billingTypeValue,
          isExternalOnly,
          isAssessment,
          // minBw,
          // assessmentRemark: assessmentRemark || undefined,
          // planningRemark: planningRemark || undefined,
          // opsRemark: opsRemark || undefined,
          isValid,
          errorMessage: errorMessage.trim(),
        };
      })
      .filter((item) => item.owner); // 只过滤节点编号为空的行

    fileData.value = parsedData;

    // 详细输出验证结果
    console.log('=== 文件解析详细结果 ===');
    console.log('总数据条数:', parsedData.length);

    const validItems = parsedData.filter((item) => item.isValid);

    if (parsedData.length === 0) {
      message.warning('没有读取到有效数据');
      return;
    }

    const validCount = validItems.length;
    const invalidCount = parsedData.filter((item) => !item.isValid).length;

    if (invalidCount > 0) {
      message.warning(
        `文件解析完成：共${parsedData.length}条数据，有效${validCount}条，无效${invalidCount}条`,
      );
    } else {
      message.success(`文件解析完成：共${validCount}条有效数据`);
    }
  }

  // 下载模板文件
  function downloadTemplate() {
    const templateData = [
      [
        '节点编号(必填)',
        '计费方式(选填)',
        '是否仅异网节点(选填，值：是、否)',
        '是否考核(选填，值：是、否)',
        // '保底带宽(选填，单位：G)',
      ],
      ['示例节点001', '日95', '否', '是'],
      ['示例节点002', '月95', '否', '是'],
      ['示例节点003', '', '否', '否'],
    ];

    const worksheet = XLSX.utils.aoa_to_sheet(templateData);
    const workbook = XLSX.utils.book_new();
    XLSX.utils.book_append_sheet(workbook, worksheet, '节点配置导入模板');

    XLSX.writeFile(workbook, '节点配置导入模板.xlsx');
  }

  // 执行导入
  async function handleImport() {
    if (validData.value.length === 0) {
      message.error('请先选择文件并确保有有效数据');
      return;
    }

    try {
      setModalProps({ confirmLoading: true });

      const importItems = validData.value.map((item) => ({
        owner: item.owner,
        billingType: item.billingType || undefined,
        minBw: item.minBw || undefined,
        isExternalOnly: item.isExternalOnly,
        isAssessment: item.isAssessment,
        assessmentRemark: item.assessmentRemark || undefined,
        planningRemark: item.planningRemark || undefined,
        opsRemark: item.opsRemark || undefined,
      }));

      await importNodeConfigs({ items: importItems });
      message.success(`成功导入 ${importItems.length} 条节点配置`);

      closeModal();
      emit('success');
    } catch (error) {
      console.error('导入失败:', error);
      message.error('导入失败，请检查网络连接或联系管理员');
    } finally {
      setModalProps({ confirmLoading: false });
    }
  }
</script>

<style lang="less" scoped>
  .import-container {
    .template-download {
      text-align: center;
    }
  }
</style>

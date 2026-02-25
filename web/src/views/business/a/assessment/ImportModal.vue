<template>
  <BasicModal
    v-bind="$attrs"
    @register="registerModal"
    title="批量导入考核规则"
    @ok="handleImport"
    width="800px"
  >
    <div class="import-container">
      <Alert
        type="info"
        message="导入说明"
        description="请下载模板文件，按照模板格式填写数据后再上传。支持Excel格式文件(.xlsx, .xls)。包含节点编号、运营商、所在地、主线业务(选填)、利用率、达标点数等字段。汇聚节点的主线业务必须填写。"
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
  import { importAssessmentRules } from '@/api/business/a';
  import * as XLSX from 'xlsx';

  defineOptions({ name: 'AssessmentImportModal' });

  const emit = defineEmits(['success', 'register']);

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
    const ispIndex = headers.findIndex((h) => h && h.includes('运营商'));
    const locationIndex = headers.findIndex((h) => h && h.includes('所在地'));
    const mainBizIndex = headers.findIndex((h) => h && h.includes('主线业务'));
    const reportTypeIndex = headers.findIndex((h) => h && h.includes('统计类型'));
    const utilizationRateIndex = headers.findIndex((h) => h && h.includes('利用率'));
    const nightPeakPointsIndex = headers.findIndex((h) => h && h.includes('达标点数'));

    if (ownerIndex === -1) {
      message.error('文件格式不正确，请确保包含"节点编号"列');
      return;
    }

    if (ispIndex === -1) {
      message.error('文件格式不正确，请确保包含"运营商"列');
      return;
    }

    if (locationIndex === -1) {
      message.error('文件格式不正确，请确保包含"所在地"列');
      return;
    }

    // 主线业务现在是选填，所以不强制检查
    // if (mainBizIndex === -1) {
    //   message.error('文件格式不正确，请确保包含"主线业务"列');
    //   return;
    // }

    if (reportTypeIndex === -1) {
      message.error('文件格式不正确，请确保包含"统计类型"列');
      return;
    }

    if (utilizationRateIndex === -1) {
      message.error('文件格式不正确，请确保包含"利用率"列');
      return;
    }

    if (nightPeakPointsIndex === -1) {
      message.error('文件格式不正确，请确保包含"达标点数"列');
      return;
    }

    const parsedData = data
      .slice(1)
      .map((row: any[]) => {
        const owner = row[ownerIndex]?.toString().trim();
        const isp = row[ispIndex]?.toString().trim();
        const location = row[locationIndex]?.toString().trim();
        const mainBiz = mainBizIndex >= 0 ? row[mainBizIndex]?.toString().trim() : '';
        const reportType = row[reportTypeIndex]?.toString().trim();
        const utilizationRate = row[utilizationRateIndex]?.toString().trim();
        const nightPeakPoints = row[nightPeakPointsIndex]?.toString().trim();

        let isValid = true;
        let errorMessage = '';

        // 验证所有字段是否为空（必填字段）
        if (!owner) {
          isValid = false;
          errorMessage += '节点编号不能为空; ';
        }

        if (!isp) {
          isValid = false;
          errorMessage += '运营商不能为空; ';
        }

        if (!location) {
          isValid = false;
          errorMessage += '所在地不能为空; ';
        }

        // 主线业务改为选填，但汇聚节点(机房总览)必须填写主线业务
        const isAggregationNode = reportType === '机房总览';
        if (isAggregationNode && !mainBiz) {
          isValid = false;
          errorMessage += '汇聚节点的主线业务不能为空; ';
        }

        if (!reportType) {
          isValid = false;
          errorMessage += '统计类型不能为空; ';
        } else if (!['机房总览', '保底业务', '削峰业务'].includes(reportType)) {
          isValid = false;
          errorMessage += '统计类型必须是"节点"、"保底"或"削峰"; ';
        }

        if (!utilizationRate) {
          isValid = false;
          errorMessage += '利用率不能为空; ';
        }

        if (!nightPeakPoints) {
          isValid = false;
          errorMessage += '达标点数不能为空; ';
        }

        return {
          owner,
          isp,
          location,
          mainBiz,
          reportType,
          utilizationRate,
          nightPeakPoints,
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
        '运营商(必填)',
        '所在地(必填)',
        '主线业务(选填，汇聚节点必填)',
        '统计类型(必填，值：机房总览、保底业务、削峰业务)',
        '利用率(必填)',
        '达标点数(必填)',
      ],
      ['示例节点001', '移动', '北京北京', 'KP2', '机房总览', '95.5', '36'],
      ['示例节点002', '电信', '上海上海', 'LE', '保底业务', '90.0', '30'],
      ['示例节点003', '联通', '广东广州', 'KP2', '削峰业务', '85.5', '24'],
    ];

    const worksheet = XLSX.utils.aoa_to_sheet(templateData);
    const workbook = XLSX.utils.book_new();
    XLSX.utils.book_append_sheet(workbook, worksheet, '考核规则导入模板');

    XLSX.writeFile(workbook, '考核规则导入模板.xlsx');
  }

  // 执行导入
  async function handleImport() {
    if (validData.value.length === 0) {
      message.error('请先选择文件并确保有有效数据');
      return;
    }

    try {
      setModalProps({ confirmLoading: true });

      // 统计类型映射: {"机房总览": 1, "保底业务": 2, "削峰业务":3}
      const reportTypeMap = {
        机房总览: 1,
        保底业务: 2,
        削峰业务: 3,
      };

      const importItems = validData.value.map((item) => ({
        owner: item.owner,
        isp: item.isp,
        location: item.location,
        mainBiz: item.mainBiz,
        // 统计类型转换为对应的数字
        reportType: reportTypeMap[item.reportType],
        // 利用率转换为0-1，保留3位小数
        utilizationRateThreshold: parseFloat(item.utilizationRate) / 100,
        // 达标点数转换为整数
        nightPeakPointsThreshold: parseInt(item.nightPeakPoints, 10),
      }));

      await importAssessmentRules({ items: importItems });
      message.success(`成功导入${importItems.length}条考核规则"`);

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

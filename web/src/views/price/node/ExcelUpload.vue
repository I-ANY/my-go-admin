<template>
  <Modal :open="visible" title="导入" @cancel="handleClose" :footer="null">
    <UploadDragger
      :fileList="fileList"
      name="file"
      :multiple="false"
      :before-upload="beforeUpload"
      class="upload-modal-dragger"
      v-if="uploadStatus === EUploadStatus.INIT"
    >
      <p class="ant-upload-drag-icon">
        <InboxOutlined />
      </p>
      <p class="ant-upload-text">点击选择文件</p>
    </UploadDragger>
    <div class="upload-status" v-if="uploadStatus === EUploadStatus.UPLOADING">
      <LoadingOutlined style="font-size: 48px" />
      <span>导入中</span>
    </div>
    <div class="upload-status" v-if="uploadStatus === EUploadStatus.SUCCESS">
      <CheckCircleFilled style="font-size: 48px" twoToneColor="#0960bd" />
      <span>导入成功</span>
    </div>
    <div class="upload-status" v-if="uploadStatus === EUploadStatus.ERROR">
      <CloseCircleFilled style="font-size: 48px" twoToneColor="#0960bd" />
      <span>导入失败</span>
    </div>
    <div class="upload-modal-tips">
      <div class="upload-modal-tips-item" @click="handleDownload"
        >1、导入数据文件请使用系统提供的模版<span class="upload-modal-tips-item-download"
          >点击下载模版文件</span
        ></div
      >
      <div class="upload-modal-tips-item">2、文件仅支持xls、xlsx格式,大小不得超过20M</div>
    </div>
  </Modal>
</template>

<script lang="ts" setup>
  import { Modal, UploadDragger, message } from 'ant-design-vue';
  import * as XLSX from 'xlsx';
  import {
    InboxOutlined,
    LoadingOutlined,
    CheckCircleFilled,
    CloseCircleFilled,
  } from '@ant-design/icons-vue';
  import { ref } from 'vue';
  import { exportNodeRecord, importNodeRecord } from '@/api/price/price';
  import { jsonToSheetXlsx } from '@/components/Excel';
  import { dateUtil } from '@/utils/dateUtil';

  enum EUploadStatus {
    INIT = 'init',
    UPLOADING = 'uploading',
    SUCCESS = 'success',
    ERROR = 'error',
  }
  const props = defineProps({
    visible: {
      type: Boolean,
      required: true,
    },
    close: {
      type: Function,
      required: true,
    },
  });
  const fileList = ref([]);
  const uploadExcelData = ref();
  const uploadStatus = ref(EUploadStatus.INIT);
  function shapeWorkSheel(sheet: XLSX.WorkSheet, range: XLSX.Range) {
    let str = ' ',
      char = 65,
      customWorkSheet = {
        t: 's',
        v: str,
        r: '<t> </t><phoneticPr fontId="1" type="noConversion"/>',
        h: str,
        w: str,
      };
    if (!sheet || !sheet['!ref']) return [];
    let c = 0,
      r = 1;
    while (c < range.e.c + 1) {
      while (r < range.e.r + 1) {
        if (!sheet[String.fromCharCode(char) + r]) {
          sheet[String.fromCharCode(char) + r] = customWorkSheet;
        }
        r++;
      }
      r = 1;
      str += ' ';
      customWorkSheet = {
        t: 's',
        v: str,
        r: '<t> </t><phoneticPr fontId="1" type="noConversion"/>',
        h: str,
        w: str,
      };
      c++;
      char++;
    }
  }

  /**
   * @description: 第一行作为头部
   */
  function getHeaderRow(sheet) {
    if (!sheet || !sheet['!ref']) return [];
    const headers: string[] = [];
    // A3:B7=>{s:{c:0, r:2}, e:{c:1, r:6}}
    const range: XLSX.Range = XLSX.utils.decode_range(sheet['!ref']);
    shapeWorkSheel(sheet, range);
    const R = range.s.r;
    /* start in the first row */
    for (let C = range.s.c; C <= range.e.c; ++C) {
      /* walk every column in the range */
      const cell = sheet[XLSX.utils.encode_cell({ c: C, r: R })];
      /* find the cell in the first row */
      let hdr = 'UNKNOWN ' + C; // <-- replace with your desired default
      if (cell && cell.t) hdr = XLSX.utils.format_cell(cell);
      headers.push(hdr);
    }
    return headers;
  }
  /**
   * @description: 获得excel数据
   */
  function getExcelData(workbook) {
    const excelData = [];
    for (const sheetName of workbook.SheetNames) {
      const worksheet = workbook.Sheets[sheetName];
      const header: string[] = getHeaderRow(worksheet);
      let results = XLSX.utils.sheet_to_json(worksheet, {
        raw: true,
        dateNF: 'YYYY-MM-DD',
      }) as object[];
      results = results.map((row: object) => {
        for (let field in row) {
          if (row[field] instanceof Date) {
            row[field] = dateUtil(row[field]).format('YYYY-MM-DD');
          }
        }
        return row;
      });

      excelData.push({
        header,
        results,
        meta: {
          sheetName,
        },
      });
    }
    return excelData;
  }
  function handleClose() {
    uploadStatus.value = EUploadStatus.INIT;
    props.close();
  }
  function handleOk() {
    const data: Array<Record<string, string>> = [];
    const excelDataList = uploadExcelData.value[0];
    const headers = excelDataList.header.map((item) => {
      if (item.trim().includes('采购单价')) {
        return '采购单价';
      }
      return item.trim();
    });
    for (const item of excelDataList.results) {
      data.push({
        name: String(item['节点'] || '').trim(),
        location: String(item['所在地'] || '').trim(),
        localIsp: String(item['运营商'] || '').trim(),
        priceType: String(item['计费方式'] || '').trim(),
        price: String(item['采购单价(元/Gbps)'] || '').trim(),
      });
    }
    importNodeRecord({ rows: data, header: headers })
      .then(function (resp) {
        uploadStatus.value = EUploadStatus.SUCCESS;
        message.success(resp.msg);
        fileList.value = [];
        uploadStatus.value = EUploadStatus.INIT;
        props.close();
      })
      .catch(() => {
        uploadStatus.value = EUploadStatus.ERROR;
        message.error('导入失败，请稍后重试');
      });
  }
  function upload(rawFile) {
    return new Promise((resolve, reject) => {
      const reader = new FileReader();
      reader.onload = async (e) => {
        try {
          const data = e.target && e.target.result;
          const workbook = XLSX.read(data, { type: 'array', cellDates: true });
          const excelData = getExcelData(workbook);
          uploadExcelData.value = excelData;
          resolve('');
        } catch (error) {
          console.error(e);
          message.error('导入失败，请稍后重试');
          reject();
        }
      };
      reader.readAsArrayBuffer(rawFile);
    });
  }
  async function beforeUpload(uploadFile) {
    uploadStatus.value = EUploadStatus.UPLOADING;
    try {
      await upload(uploadFile);
      handleOk();
    } catch (e) {
      console.error(e);
      uploadStatus.value = EUploadStatus.ERROR;
    }

    return false;
  }
  function handleDownload() {
    exportNodeRecord({})
      .then((resp) => {
        let data = resp.items;
        let result: { [key: string]: any }[] = [];
        for (let i = 0; i < data.length; i++) {
          let map = {};
          map['节点'] = data[i]['节点'];
          map['所在地'] = data[i]['所在地'];
          map['运营商'] = data[i]['运营商'];
          map['计费方式'] = data[i]['计费方式'];
          map['采购单价(元/Gbps)'] = data[i]['采购单价'];
          result.push(map);
        }
        jsonToSheetXlsx({
          data: result,
          filename: '单价导入模板.xlsx',
          write2excelOpts: {
            bookType: 'xlsx',
          },
        });
      })
      .then(() => {
        message.success('下载成功');
      })
      .catch((e) => {
        console.error(e);
        message.error('下载失败，请稍后重试');
      });
  }
</script>

<style lang="less" scoped>
  .upload-modal-dragger {
    display: flex;
    width: 480px;
    margin: 10px auto;
  }

  .upload-modal-tips {
    display: flex;
    flex-direction: column;
    padding: 20px;

    &-item {
      display: flex;
      flex-direction: row;
      align-items: center;

      &-download {
        margin-left: 10px;
        color: #0960bd;
        cursor: pointer;
      }
    }
  }

  .upload-status {
    display: flex;
    position: relative;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    width: 100%;
    width: 480px;
    height: 100%;
    height: 127.14px;
    margin: 10px auto;
    transition: border-color 0.3s;
    border: 1px dashed #d9d9d9;
    border-radius: 8px;
    background: rgb(0 0 0 / 2%);
    text-align: center;
    cursor: pointer;

    & > span {
      margin-top: 16px;
      font-size: 16px;
    }
  }
</style>

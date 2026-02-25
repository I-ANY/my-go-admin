<template>
  <ImpExcel @success="upload">
    <a-button type="primary" :disabled="disabled" :loading="disabled">导入</a-button>
  </ImpExcel>
</template>

<script lang="ts" setup>
  import { ref } from 'vue';
  import { ImpExcel, ExcelData } from '@/components/Excel';
  import { importCategory, importRecord, importNodeRecord } from '@/api/price/price';

  const emit = defineEmits(['successUpload']);
  const disabled = ref(false);

  const props = defineProps({
    excelType: {
      type: String,
      required: true,
    },
  });

  function upload(excelDataList: ExcelData[]) {
    disabled.value = true;
    const data: Array<Record<string, string>> = [];
    if (props.excelType === 'category') {
      for (const item of excelDataList[0].results) {
        data.push({
          category: String(item['业务组'] || '').trim(),
          bizs: String(item['业务详情'] || '').trim(),
          outName: String(item['对外业务名'] || '').trim(),
          describe: String(item['备注'] || '').trim(),
        });
      }
      importCategory({ rows: data, header: excelDataList[0].header })
        .then(function (resp) {
          emit('successUpload', resp.msg);
        })
        .finally(() => {
          disabled.value = false;
        });
      return;
    }

    if (props.excelType === 'nodeRecord') {
      for (const item of excelDataList[0].results) {
        data.push({
          name: String(item['节点'] || '').trim(),
          location: String(item['所在地'] || '').trim(),
          localIsp: String(item['运营商'] || '').trim(),
          price: String(item['采购单价'] || '').trim(),
        });
      }
      importNodeRecord({ rows: data, header: excelDataList[0].header })
        .then(function (resp) {
          emit('successUpload', resp.msg);
        })
        .finally(() => {
          disabled.value = false;
        });
      return;
    }

    for (const item of excelDataList[0].results) {
      data.push({
        category: String(item['业务组'] || '').trim(),
        region: item['大区'] && item['大区'].trim() !== '' ? item['大区'].trim() : null,
        zone: item['省份'] && item['省份'].trim() !== '' ? item['省份'].trim() : null,
        localIsp: String(item['本网运营商'] || '').trim(),
        bizIsp:
          item['异网运营商'] && item['异网运营商'].trim() !== '' ? item['异网运营商'].trim() : null,
        bizIspMode:
          item['跨网计费方式'] && item['跨网计费方式'].trim() !== ''
            ? item['跨网计费方式'].trim()
            : null,
        mode: String(item['计费方式'] || '').trim(),
        price: String(item['单价'] + '' || '').trim(),
        low: String(item['溜缝业务'] + '' || '否').trim(),
        describe: String(item['备注'] || '').trim(),
      });
    }
    importRecord({ rows: data, header: excelDataList[0].header })
      .then(function (resp) {
        emit('successUpload', resp.msg);
      })
      .finally(() => {
        disabled.value = false;
      });
  }
</script>

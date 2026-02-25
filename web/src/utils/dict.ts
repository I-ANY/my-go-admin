import { DictDataItem } from '@/api/sys/model/dict';
import { useDictStore } from '@/store/modules/dict';
import { DefaultOptionType } from 'ant-design-vue/es/select';

export const getSelectOptions = function (items: DictDataItem[]): DefaultOptionType[] {
  let result: DefaultOptionType[] = [];
  result = items.map((item) => {
    return { label: item.dictLabel, value: item.dictValue };
  });
  return result;
};

export const getOptionMap = function (items: DictDataItem[]) {
  const result: any = {};
  items.forEach((item) => {
    result[item.dictValue] = item;
  });
  return result;
};

export const getEnumsFromDict = function (typeCode: string): DictDataItem[] {
  const dictStore = useDictStore();
  let dictData: DictDataItem[] = [];
  dictStore.getAllDict.forEach((dictType) => {
    if (dictType.typeCode == typeCode) {
      dictData = dictType.dictData || [];
    }
  });
  return dictData;
};

// 将字典中指定 typeCode 的字典数据转换成 select 选项所需要的数据
export const getSelectOptionsFromDict = function (typeCode: string): DefaultOptionType[] {
  const dictData = getEnumsFromDict(typeCode);
  return getSelectOptions(dictData);
};

// 将字典中指定 typeCode 的字典数据转换成 map类型 key=字典值  value=dictData
export const getDictDataMapFromDict = function (typeCode: string): Map<string, DictDataItem> {
  const dictData = getEnumsFromDict(typeCode);
  return getOptionMap(dictData);
};

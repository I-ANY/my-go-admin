import type { FormSchema } from '@/components/Form';
import { GetAreaList } from '@/api/business/k';
import { ref } from 'vue';
import { DefaultOptionType } from 'ant-design-vue/es/select';

interface TreeData {
  areas: any[];
  provinces: Record<string, any[]>;
  cities: Record<string, any[]>;
}

interface AreaSelectOptions {
  hasCity?: boolean;
  form: any;
  fields: {
    area: string;
    province: string;
    city?: string;
  };
}

export function useAreaSelect(options: AreaSelectOptions) {
  const areaOptions = ref<any[]>([]);
  const provinceOptions = ref<any[]>([]);
  const cityOptions = ref<any[]>([]);
  const treeData = ref<TreeData>();
  const currentArea = ref('');
  const currentProvince = ref('');

  const initAreaData = async () => {
    try {
      const areaData = await GetAreaList({});
      treeData.value = convertToTree(areaData.items);
      areaOptions.value = treeData.value?.areas || [];
      updateFormSchemas();
    } catch (error) {
      console.error('初始化区域数据失败:', error);
    }
  };

  const convertToTree = (data: any[]): TreeData => {
    const areaSet = new Set<string>();
    const provinceMap = new Map<string, any[]>();
    const cityMap = new Map<string, any[]>();

    data.forEach((item) => {
      if (item.area_name) {
        areaSet.add(item.area_name);

        // 处理省份
        if (item.province_name) {
          const provinces = provinceMap.get(item.area_name) || [];
          if (!provinces.some((p) => p.value === item.province_name)) {
            provinces.push({ value: item.province_name, label: item.province_name });
          }
          provinceMap.set(item.area_name, provinces);
        }

        // 处理城市
        if (options.hasCity && item.province_name && item.city_name) {
          const provinceKey = `${item.area_name}-${item.province_name}`;
          const cities = cityMap.get(provinceKey) || [];
          if (!cities.some((c) => c.value === item.city_name)) {
            cities.push({ value: item.city_name, label: item.city_name });
          }
          cityMap.set(provinceKey, cities);
        }
      }
    });

    return {
      areas: Array.from(areaSet).map((area) => ({ value: area, label: area })),
      provinces: Object.fromEntries(provinceMap),
      cities: Object.fromEntries(cityMap),
    };
  };

  const handleAreaChange = (value: string) => {
    currentArea.value = value;
    provinceOptions.value = treeData.value?.provinces[value] || [];
    cityOptions.value = [];

    // 强制更新表单字段选项
    options.form.updateSchema([
      {
        field: options.fields.province,
        componentProps: {
          options: provinceOptions,
          onChange: options.hasCity ? handleProvinceChange : undefined,
        },
      },
      ...(options.hasCity
        ? [
            {
              field: options.fields.city,
              componentProps: {
                options: cityOptions,
              },
            },
          ]
        : []),
    ]);

    options.form.setFieldsValue({
      [options.fields.province || 'province_name']: '',
      ...(options.hasCity && { [options.fields.city || 'city_name']: '' }),
    });
  };

  const handleProvinceChange = (value: string) => {
    currentProvince.value = value;
    if (options.hasCity) {
      const provinceKey = `${currentArea.value}-${value}`;
      cityOptions.value = treeData.value?.cities[provinceKey] || [];
    }
    options.form.setFieldsValue({ [options.fields.city || 'city_name']: '' });
  };

  const updateFormSchemas = () => {
    const schemas: FormSchema[] = [];

    schemas.push({
      field: options.fields.area,
      component: 'Select',
      componentProps: {
        options: areaOptions,
        onChange: handleAreaChange,
      } as any,
    });

    schemas.push({
      field: options.fields.province,
      component: 'Select',
      componentProps: {
        options: provinceOptions,
        onChange: options.hasCity ? handleProvinceChange : undefined,
      } as any,
    });

    if (options.hasCity && options.fields.city) {
      schemas.push({
        field: options.fields.city,
        component: 'Select',
        componentProps: {
          options: cityOptions,
        } as any,
      });
    }

    options.form.updateSchema(schemas, true);
  };

  return {
    areaOptions,
    provinceOptions,
    cityOptions,
    initAreaData,
    currentArea,
    currentProvince,
  };
}

interface AreaInfo {
  areaProvinceTree: Record<string, DefaultOptionType[]>;
  allProvinceOptions: DefaultOptionType[];
  areaOptions: DefaultOptionType[];
}

export function getProvinceName(item: any): string {
  return item.province_name;
}

export function getProvince(item: any): string {
  return item.province;
}

export async function getAreaData(getProvinceName: (item: any) => string): Promise<AreaInfo> {
  const data: AreaInfo = {
    areaProvinceTree: {},
    allProvinceOptions: [],
    areaOptions: [],
  };
  const res = await GetAreaList({});
  const areaInfo = res.items;
  const areaProvinceTree: Record<string, DefaultOptionType[]> = {};

  for (let i = 0; i < areaInfo?.length; i++) {
    const item = areaInfo[i];
    const areaName = item.area_name;
    if (!areaProvinceTree[areaName]) {
      areaProvinceTree[areaName] = [];
    }

    // 判断当前区域中是否存在该省份
    let isExist = false;
    areaProvinceTree[areaName].forEach((provinceOption) => {
      if (provinceOption.label === getProvinceName(item)) {
        isExist = true;
      }
    });

    if (!isExist) {
      areaProvinceTree[areaName].push({
        label: getProvinceName(item),
        value: getProvinceName(item),
      });
    }
  }

  const allProvinceOptions: DefaultOptionType[] = [];
  const areaOptions: DefaultOptionType[] = [];
  Object.keys(areaProvinceTree).forEach((areaName: string) => {
    areaOptions.push({
      label: areaName,
      value: areaName,
    });
    areaProvinceTree[areaName].forEach((provinceOption) => {
      allProvinceOptions.push({
        label: provinceOption.label,
        value: provinceOption.value,
      });
    });
  });
  data.allProvinceOptions = allProvinceOptions;
  data.areaProvinceTree = areaProvinceTree;
  data.areaOptions = areaOptions;

  return data;
}

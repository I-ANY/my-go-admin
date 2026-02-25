export interface DictTypeItem {
  Id?: number;
  typeName: string;
  typeCode: string;
  sort?: number;
  status?: number;
  remark?: string;
  dictData?: DictDataItem[];
}

export interface DictDataItem {
  Id?: number;
  dictTypeId: number;
  dictLabel: string;
  dictValue: string;
  sort?: number;
  status?: number;
  remark: string;
  color: string;
}

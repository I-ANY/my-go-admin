export interface ApiItem {
  id: number;
  path: string;
  method: string;
  handler: string;
}

/**
 * @description: Get menu return value
 */
export type getApiListModel = ApiItem[];

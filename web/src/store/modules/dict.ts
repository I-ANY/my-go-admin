import { DictTypeItem } from '@/api/sys/model/dict';
import { defineStore } from 'pinia';
import { store } from '@/store';
import { GetAllDict } from '@/api/sys/dict';

interface DictState {
  allDict: DictTypeItem[];
}

export const useDictStore = defineStore({
  id: 'app-dict',
  state: (): DictState => ({
    allDict: [],
  }),
  getters: {
    getAllDict(state): DictTypeItem[] {
      return state.allDict || [];
    },
  },
  actions: {
    async setAllDict() {
      const data = await GetAllDict();
      this.allDict = data;
    },
    resetState() {
      this.allDict = [];
    },
  },
});

// Need to be used outside the setup
export function useDictStoreWithOut() {
  return useDictStore(store);
}

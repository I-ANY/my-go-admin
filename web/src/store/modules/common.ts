import { defineStore } from 'pinia';
import { store } from '@/store';

interface DictState {
  inspectHostname: Nullable<string>;
}

export const useCommonStore = defineStore({
  id: 'app-common',
  state: (): DictState => ({
    inspectHostname: null,
  }),
  getters: {
    getInspectHostname(state): Nullable<string> {
      return state.inspectHostname;
    },
  },
  actions: {
    setInspectHostname(inspectHostname: string) {
      this.inspectHostname = inspectHostname;
    },
    clearInspectHostname() {
      this.inspectHostname = null;
    },
    resetState() {
      this.inspectHostname = null;
    },
  },
});

// Need to be used outside the setup
export function useCommonStoreWithOut() {
  return useCommonStore(store);
}

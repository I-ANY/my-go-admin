<template>
  <div class="m-4 mr-0 overflow-hidden bg-white">
    <BasicTree
      title="部门列表"
      toolbar
      treeWrapperClassName="h-[calc(100%-35px)] overflow-auto"
      :clickRowToExpand="false"
      :treeData="treeData"
      :fieldNames="{ key: 'id', title: 'name' }"
      @select="handleSelect"
      ref="treeRef"
    />
  </div>
</template>
<script lang="ts" setup>
  import { nextTick, onMounted, ref } from 'vue';

  import { BasicTree } from '@/components/Tree';
  import { getDeptList } from '@/api/demo/system';

  defineOptions({ name: 'DeptTree' });

  const emit = defineEmits(['select']);
  const treeRef = ref<any>(null);
  const treeData = ref<any[]>([]);

  async function fetch() {
    let res = await getDeptList();
    treeData.value = res.items;
    nextTick(() => {
      treeRef.value.filterByLevel(1);
    });
  }

  function handleSelect(keys) {
    emit('select', keys[0]);
  }

  onMounted(() => {
    fetch();
  });
</script>

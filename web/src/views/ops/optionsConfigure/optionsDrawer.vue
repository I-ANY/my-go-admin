<template>
  <BasicDrawer
    v-bind="$attrs"
    @register="registerDrawer"
    :title="getTitle"
    width="1200px"
    destroyOnClose
    :showFooter="true"
    ok-text="保存"
    cancel-text="取消"
    @ok="handleSubmit"
    @close="handleClose"
  >
    <BasicForm @register="registerForm">
      <template #params="{ field }">
        <div
          v-if="field === 'params'"
          class="monaco"
          id="monaco"
          ref="monacoRef"
          style="height: 400px"
        ></div>
      </template>
    </BasicForm>
  </BasicDrawer>
</template>
<script lang="ts" setup>
  import { ref, computed, unref, nextTick, shallowRef, reactive } from 'vue';
  import { BasicForm, useForm } from '@/components/Form';
  import { configFormSchema } from '@/views/ops/optionsConfigure/data';
  import { BasicDrawer, useDrawerInner } from '@/components/Drawer';
  import * as monaco from 'monaco-editor';
  import { AddOperation, EditOperation, GetScriptTasks } from '@/api/ops/execute';
  import { message } from 'ant-design-vue';
  import { GetAuthedBiz } from '@/api/business/biz';

  const isUpdate = ref(false);
  const getTitle = computed(() => (!unref(isUpdate) ? '新增操作' : '编辑操作'));
  const emit = defineEmits(['success', 'register']);

  const [registerDrawer, { setDrawerProps, closeDrawer }] = useDrawerInner(async (data) => {
    let categoryOptions = reactive([
      {
        name: '公共',
        id: 0,
      },
    ]);
    await resetFields();
    isUpdate.value = data.isUpdate;
    // 获取业务大类选项
    const bizData = await GetAuthedBiz({});
    if (bizData.categories) {
      categoryOptions.push(...bizData.categories);
      await updateSchema([
        {
          field: 'businessId',
          componentProps: {
            options: categoryOptions,
            fieldNames: {
              label: 'name',
              value: 'id',
            },
            onChange: async (value) => {
              // 获取业务下的子分类作为 parentID 的选项
              if (value === 0) {
                await setFieldsValue({
                  scriptPath: 'general',
                });
              }
              // 根据选择的 businessID 找到对应的业务数据
              const res = await GetScriptTasks({ businessId: value, type: 'group' });
              if (res.items) {
                await updateSchema([
                  {
                    field: 'parentId',
                    componentProps: {
                      options: res.items,
                      fieldNames: {
                        label: 'name',
                        value: 'id',
                      },
                    },
                  },
                ]);
              }
            },
          },
        },
      ]);
    }

    setDrawerProps({ confirmLoading: false });
    // 重置编辑器状态
    if (editor.value) {
      editor.value.dispose();
      editor.value = null;
    }
    // 模态框打开后初始化编辑器
    await nextTick(() => {
      initEditor();
      // 在这里给编辑器赋值
      if (data.isUpdate && data.record?.params) {
        editor.value?.setValue(data.record.params);
      }
      if (data.isUpdate) {
        setFieldsValue(data.record);
      } else {
        setFieldsValue({
          businessId: data.businessId,
          parentId: data.optionGroupId,
          status: 1,
        });
      }
    });
  });

  const [registerForm, { setFieldsValue, resetFields, validate, getFieldsValue, updateSchema }] =
    useForm({
      labelWidth: 80,
      baseColProps: { span: 24 },
      schemas: configFormSchema,
      showActionButtonGroup: false,
    });

  const editor = shallowRef<monaco.editor.IStandaloneCodeEditor | null>(null);
  const monacoRef = ref<HTMLElement | null>(null);

  const initEditor = () => {
    if (!monacoRef.value || editor.value) return;
    editor.value = monaco.editor.create(monacoRef.value, {
      value: '',
      language: 'yaml',
      theme: 'vs-dark',
      automaticLayout: true,
      lineHeight: 24,
      tabSize: 4,
      fontSize: 16,
      placeholder: '请输入操作详情',
      minimap: { enabled: false },
      readOnly: false,
      domReadOnly: true,
      quickSuggestions: {
        other: true,
        comments: true,
        strings: true,
      },
      suggestOnTriggerCharacters: true, // 输入任意字符都可触发
    });
    // 注册YAML补全项提供器
    monaco.languages.registerCompletionItemProvider('yaml', {
      triggerCharacters: [],
      provideCompletionItems: (model, position) => {
        const word = model.getWordUntilPosition(position);
        const range = {
          startLineNumber: position.lineNumber,
          endLineNumber: position.lineNumber,
          startColumn: word.startColumn,
          endColumn: word.endColumn,
        };

        return {
          suggestions: [
            {
              label: 'name',
              kind: monaco.languages.CompletionItemKind.Field,
              insertText: 'name: ',
              detail: '操作名称',
              range: range,
            },
            {
              label: 'value',
              kind: monaco.languages.CompletionItemKind.Field,
              documentation: '脚本名称',
              insertText: 'value: ',
              detail: '脚本名称',
              range: range,
            },
            {
              label: 'params',
              kind: monaco.languages.CompletionItemKind.Field,
              documentation: '脚本参数列表',
              insertText: 'params: ',
              detail: '脚本参数列表',
              range: range,
            },
            {
              label: 'params.options',
              kind: monaco.languages.CompletionItemKind.Field,
              documentation: '选项',
              insertText: 'options: \n',
              detail: '选项',
              range: range,
            },
            {
              label: 'params.option',
              kind: monaco.languages.CompletionItemKind.Field,
              documentation: '选项',
              insertText: `- label:
      value:
      children: [] # 有 则级联选项
      params: [] # 这个选项需要哪些参数。只对此选项的后面参数有效。用法参考: 长A/业务操作/修改dcache硬盘配置`,
              detail: '选项',
              range: range,
            },
            {
              label: 'params.default',
              kind: monaco.languages.CompletionItemKind.Field,
              documentation: '参数的默认值',
              insertText: 'default:',
              detail: '参数的默认值',
              range: range,
            },
            {
              label: 'params.required',
              kind: monaco.languages.CompletionItemKind.Field,
              documentation: '参数的默认值',
              insertText: 'required:',
              detail: '是否必须',
              range: range,
            },
            {
              label: 'params',
              kind: monaco.languages.CompletionItemKind.Field,
              documentation: '参数',
              insertText: `- name:            # 参数名称
      default:         # 默认值
      required: false  # 是否必须
      csv: |-          # 使用csv生成options
        a,b
      options: []      # 选项
      type: input      # 类型,给了options select
      remote: null
    `,
              detail: '参数',
              range: range,
            },
            {
              label: 'params.remote',
              kind: monaco.languages.CompletionItemKind.Field,
              documentation: '从api获取options',
              insertText: `url:   # api地址 只能是get方法
    query:  # 如果是级联, 参考 AACDN/调整参数和备份/版本更新配置
    path:   # json的数据路径
    label:  # 可不填
    value:  #  可不填
    children: null # 如果配置,则是级联操作 参考 AACDN/调整参数和备份/版本更新配置`,
              detail: '从api获取options',
              range: range,
            },
          ],
        };
      },
    });

    // 监听内容变化
    // editor.value.onDidChangeModelContent(() => {
    //   console.log('当前内容：', editor.value!.getValue());
    // });
  };

  function handleClose() {
    closeDrawer();
    // 其他取消逻辑
  }

  async function handleSubmit() {
    await validate();
    // 获取参数
    const params = getFieldsValue();
    params.type = 'function';
    params.params = editor.value!.getValue();
    if (isUpdate.value) {
      await EditOperation(params.id, params);
      emit('success');
      message.success('操作成功');
      closeDrawer();
    } else {
      await AddOperation(params);
      emit('success');
      message.success('操作成功');
      closeDrawer();
    }
  }
</script>
<style lang="less" scoped>
  .monaco {
    min-height: 200px;
    max-height: 800px;
    overflow: auto;
    resize: vertical;
  }
</style>

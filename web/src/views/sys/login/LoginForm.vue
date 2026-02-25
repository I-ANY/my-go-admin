<template>
  <div>
    <LoginFormTitle v-show="getShow" class="enter-x" />
    <!-- <div class="mt-80px">
    <ARow class="w-600px" :align="'middle'" :justify="'center'" :gutter="24">
      <ACol :span="18">
        <Button type="primary" block>门户网站登录</Button>
      </ACol>
    </ARow>
  </div> -->
    <Form
      class="p-4 enter-xl"
      :model="formData"
      :rules="getFormRules"
      ref="formRef"
      v-show="getShow"
      @keypress.enter="handleLogin"
    >
      <FormItem name="account" class="enter-x">
        <Input
          size="large"
          v-model:value="formData.account"
          :placeholder="t('sys.login.userName')"
          class="fix-auto-fill"
        />
      </FormItem>
      <FormItem name="password" class="enter-x">
        <InputPassword
          size="large"
          visibilityToggle
          v-model:value="formData.password"
          :placeholder="t('sys.login.password')"
        />
      </FormItem>

      <!-- <ARow class="enter-x">
      <ACol :span="12">
        <FormItem>
          <Checkbox v-model:checked="rememberMe" size="small">
            {{ t('sys.login.rememberMe') }}
          </Checkbox>
        </FormItem>
      </ACol>
      <ACol :span="12">
        <FormItem :style="{ 'text-align': 'right' }">
          <Button type="link" size="small" @click="setLoginState(LoginStateEnum.RESET_PASSWORD)">
            {{ t('sys.login.forgetPassword') }}
          </Button>
        </FormItem>
      </ACol>
    </ARow> -->

      <FormItem class="enter-x">
        <Button type="primary" size="large" block @click="handleLogin" :loading="loading">
          {{ t('sys.login.loginButton') }}
        </Button>
        <!-- <Button size="large" class="mt-4 enter-x" block @click="handleRegister">
        {{ t('sys.login.registerButton') }}
      </Button> -->
      </FormItem>
      <!-- <ARow class="enter-x" :gutter="[16, 16]">
      <ACol :md="8" :xs="24">
        <Button block @click="setLoginState(LoginStateEnum.MOBILE)">
          {{ t('sys.login.mobileSignInFormTitle') }}
        </Button>
      </ACol>
      <ACol :md="8" :xs="24">
        <Button block @click="setLoginState(LoginStateEnum.QR_CODE)">
          {{ t('sys.login.qrSignInFormTitle') }}
        </Button>
      </ACol>
      <ACol :md="8" :xs="24">
        <Button block @click="setLoginState(LoginStateEnum.REGISTER)">
          {{ t('sys.login.registerButton') }}
        </Button>
      </ACol>
    </ARow> -->

      <Divider class="enter-x">{{ t('sys.login.otherSignIn') }}</Divider>
      <div
        class="flex justify-evenly enter-x min-h-14 min-w-14"
        :class="`${prefixCls}-sign-in-way`"
      >
        <!-- <GithubFilled />
      <WechatFilled />
      <AlipayCircleFilled />
      <GoogleCircleFilled />
      <TwitterCircleFilled /> -->
        <Tooltip title="星云门户平台登录">
          <HomeOutlined
            :style="{ fontSize: '40px', color: '#e72c1c' }"
            @click="handleGetStarPortalLoginUrl"
          />
        </Tooltip>
      </div>
    </Form>
  </div>
</template>
<script lang="ts" setup>
  import { reactive, ref, unref, computed, onBeforeMount } from 'vue';

  import { Form, Input, Button, Divider, Tooltip } from 'ant-design-vue';
  import {
    // GithubFilled,
    // WechatFilled,
    // AlipayCircleFilled,
    // GoogleCircleFilled,
    // TwitterCircleFilled,
    HomeOutlined,
  } from '@ant-design/icons-vue';
  import LoginFormTitle from './LoginFormTitle.vue';

  import { useI18n } from '@/hooks/web/useI18n';
  import { useMessage } from '@/hooks/web/useMessage';

  import { useUserStore } from '@/store/modules/user';
  import { LoginStateEnum, useLoginState, useFormRules, useFormValid } from './useLogin';
  import { useDesign } from '@/hooks/web/useDesign';
  import { getStarPortalLoginUrl } from '@/api/sys/user';
  import { useRoute, useRouter } from 'vue-router';
  import { PageEnum } from '@/enums/pageEnum';
  //import { onKeyStroke } from '@vueuse/core';

  // const ACol = Col;
  // const ARow = Row;
  const route = useRoute();
  const router = useRouter();
  const FormItem = Form.Item;
  const InputPassword = Input.Password;
  const { t } = useI18n();
  const { notification } = useMessage();
  const { prefixCls } = useDesign('login');
  const userStore = useUserStore();

  const { getLoginState } = useLoginState();
  const { getFormRules } = useFormRules();

  const formRef = ref();
  const loading = ref(false);
  const emits = defineEmits(['pageLoading']);

  const formData = reactive({
    account: '',
    password: '',
  });

  const { validForm } = useFormValid(formRef);

  //onKeyStroke('Enter', handleLogin);

  const getShow = computed(() => unref(getLoginState) === LoginStateEnum.LOGIN);

  async function handleLogin() {
    const data = await validForm();
    if (!data) return;
    try {
      loading.value = true;
      const userInfo = await userStore.login({
        password: data.password,
        username: data.account,
        mode: 'none', //不要默认的错误提示
      });
      if (userInfo) {
        notification.success({
          message: t('sys.login.loginSuccessTitle'),
          description: `${t('sys.login.loginSuccessDesc')}: ${userInfo.nickname}`,
          duration: 3,
        });
      }
    } catch (error) {
      // notification.error({
      //   message: t('登录失败'),
      //   description: (error as unknown as Error).message || t('sys.api.networkExceptionMsg'),
      //   duration: 3,
      // });
      // createErrorModal({
      //   title: t('sys.api.errorTip'),
      //   content: (error as unknown as Error).message || t('sys.api.networkExceptionMsg'),
      //   getContainer: () => document.body.querySelector(`.${prefixCls}`) || document.body,
      // });
    } finally {
      loading.value = false;
    }
  }
  async function handleGetStarPortalLoginUrl() {
    let redirectUri =
      location.origin + route.path + '?redirect=' + (route.query.redirect || PageEnum.BASE_HOME);
    redirectUri = encodeURIComponent(redirectUri);
    let params = {
      redirectUri,
    };
    const res = await getStarPortalLoginUrl(params);
    if (res && res.loginUrl) {
      window.location.href = res.loginUrl;
    }
  }
  async function handleStarPortalLogin() {
    const code = route.query.code as string;
    const state = route.query.state as string;
    let redirectUri = location.origin + route.path;
    if (route.query.redirect) {
      redirectUri += '?redirect=' + route.query.redirect;
    }
    if (!code || !state || !redirectUri) {
      return;
    }
    try {
      emits('pageLoading', true);
      const userInfo = await userStore.starPortallogin({
        code: code,
        state: state,
        redirectUri: redirectUri,
        mode: 'none', //不要默认的错误提示
      });
      if (userInfo) {
        notification.success({
          message: t('sys.login.loginSuccessTitle'),
          description: `${t('sys.login.loginSuccessDesc')}: ${userInfo.nickname}`,
          duration: 3,
        });
      }
    } catch (error) {
      // notification.error({
      //   message: t('登录失败'),
      //   description: (error as unknown as Error).message || t('sys.api.networkExceptionMsg'),
      //   duration: 3,
      // });
      // createErrorModal({
      //   title: t('sys.api.errorTip'),
      //   content: (error as unknown as Error).message || t('sys.api.networkExceptionMsg'),
      //   getContainer: () => document.body.querySelector(`.${prefixCls}`) || document.body,
      // });
      let querys = { ...route.query };
      // 登录失败时，获取删除当前query参数中的code、state参数，防止用户刷新页面是重复报错，因为code只有一次有效期
      delete querys.code;
      delete querys.state;
      await router.push({
        path: route.path,
        query: querys,
      });
    } finally {
      emits('pageLoading', false);
    }
  }
  onBeforeMount(handleStarPortalLogin);
  onBeforeMount(checkAutoAuth);
  function checkAutoAuth() {
    if (route.query.autoAuth) {
      handleGetStarPortalLoginUrl();
    }
  }
</script>

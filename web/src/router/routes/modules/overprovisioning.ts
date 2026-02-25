import type { AppRouteModule } from '@/router/types';
import { LAYOUT } from '@/router/constant';

const overprovisioning: AppRouteModule = {
  path: '/business/overprovisioning',
  name: 'OverProvisioning',
  component: LAYOUT,
  redirect: '/business/overprovisioning/rulemanage',
  meta: {
    orderNo: 2000,
    icon: 'ant-design:appstore-outlined',
    title: '超配管理',
  },
  children: [
    {
      path: 'rulemanage',
      name: 'OverProvisioningRuleManage',
      meta: {
        title: '规则管理',
        ignoreKeepAlive: false,
      },
      component: () => import('@/views/business/overprovisioning/rulemanage/rulemanage.vue'),
    },
    {
      path: 'records',
      name: 'OverProvisioningRecords',
      meta: {
        title: '检测记录',
        ignoreKeepAlive: false,
      },
      component: () => import('@/views/business/overprovisioning/records/index.vue'),
    },
  ],
};

export default overprovisioning;

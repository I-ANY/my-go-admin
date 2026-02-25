import axios from 'axios';
import { Modal } from 'ant-design-vue';
import { useGlobSetting } from '@/hooks/setting';

const { disableVersionCheck } = useGlobSetting();

let localVersion: string = '';
let modalInstance: any = null;
let interval: any = null;
const DURATION = 1000 * 15;
// 检查服务端是否已经更新，如果更新刷新页面
export async function checkAppNewVersion() {
  // disableVersionCheck
  if (disableVersionCheck === 'true') {
    return;
  }
  const url = `/version.json?t=${Date.now()}`;
  let res: any = null;
  try {
    res = await axios.get(url);
  } catch (err) {
    console.error('checkAppNewVersion error: ', err);
  }
  if (!res || res.status !== 200 || !res.data || !res.data.version) return;
  const remoteVersion = res?.data?.version || '';
  if (localVersion !== remoteVersion && remoteVersion) {
    // 检查到新的版本，需要刷新页面
    if (!modalInstance && localVersion != '') {
      modalInstance = Modal.confirm({
        title: '发现新版本',
        content: '检测到系统已更新，建议立即刷新页面以获得最新功能。不刷新可能导致部分功能异常！',
        okText: '立即刷新',
        cancelText: '稍后刷新',
        closable: false,
        onOk: () => {
          window.location.reload();
        },
        onCancel: () => {
          modalInstance = null;
        },
      });
    }
    localVersion = remoteVersion;
  }
}
// 监听页面打开显示
document.addEventListener('visibilitychange', function () {
  if (document.hidden) {
    // 页面不可见时，停止检查新版本
    if (interval) {
      clearInterval(interval);
      interval = null;
    }
  } else {
    // 页面可见时，检查新版本
    try {
      checkAppNewVersion();
    } catch (error) {
      console.error('checkAppNewVersion error: ', error);
    }
    if (!interval) {
      interval = setInterval(() => {
        checkAppNewVersion();
      }, DURATION);
    }
  }
});
interval = setInterval(() => {
  checkAppNewVersion();
}, DURATION);

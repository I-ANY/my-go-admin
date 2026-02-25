import type { GlobConfig } from '#/config';

import { getAppEnvConfig } from '@/utils/env';

export const useGlobSetting = (): Readonly<GlobConfig> => {
  const {
    VITE_GLOB_APP_TITLE,
    VITE_GLOB_API_URL,
    VITE_GLOB_API_URL_PREFIX,
    VITE_GLOB_UPLOAD_URL,
    VITE_GLOB_LINUX_API_URL,
    VITE_GLOB_DISABLE_CHAT_BOT,
    VITE_GLOB_DISABLE_VERSION_CHECK,
  } = getAppEnvConfig();
  getAppEnvConfig();

  // Take global configuration
  const glob: Readonly<GlobConfig> = {
    title: VITE_GLOB_APP_TITLE,
    apiUrl: VITE_GLOB_API_URL,
    shortName: VITE_GLOB_APP_TITLE.replace(/\s/g, '_').replace(/-/g, '_'),
    urlPrefix: VITE_GLOB_API_URL_PREFIX,
    uploadUrl: VITE_GLOB_UPLOAD_URL,
    linuxApiUrl: VITE_GLOB_LINUX_API_URL,
    disableChatBot: VITE_GLOB_DISABLE_CHAT_BOT,
    disableVersionCheck: VITE_GLOB_DISABLE_VERSION_CHECK,
  };
  return glob as Readonly<GlobConfig>;
};

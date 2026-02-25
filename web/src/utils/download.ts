import { defHttp } from './http/axios';

export const downloadWithRes = function (res: any): string {
  const blob = new Blob([res.data], { type: res.headers['content-type'] });
  // 创建新的URL并指向File对象或者Blob对象的地址
  const blobURL = window.URL.createObjectURL(blob);
  // 创建a标签，用于跳转至下载链接
  const tempLink = document.createElement('a');
  tempLink.style.display = 'none';
  tempLink.href = blobURL;
  const contentDisposition =
    res.headers['content-disposition'] || `attachment;filename=unknown_filename`;

  // 解析文件名，优先使用 filename*=UTF-8'' 格式
  let filename = 'unknown_filename';
  const dispositionParts = contentDisposition.split(';');

  for (const part of dispositionParts) {
    const trimmedPart = part.trim();
    if (trimmedPart.startsWith("filename*=UTF-8''")) {
      // 处理 filename*=UTF-8'' 格式
      const encodedFilename = trimmedPart.substring(17); // 去掉 'filename*=UTF-8'''
      try {
        filename = decodeURIComponent(encodedFilename);
        break;
      } catch (e) {
        console.warn('Failed to decode filename:', encodedFilename);
      }
    } else if (trimmedPart.startsWith('filename=')) {
      // 处理普通 filename= 格式
      const value = trimmedPart.substring(9);
      if (value.startsWith('"') && value.endsWith('"')) {
        filename = value.substring(1, value.length - 1);
      } else {
        filename = value;
      }
    }
  }

  tempLink.setAttribute('download', filename);
  // 兼容：某些浏览器不支持HTML5的download属性
  if (typeof tempLink.download === 'undefined') {
    tempLink.setAttribute('target', '_blank');
  }
  // 挂载a标签
  document.body.appendChild(tempLink);
  tempLink.click();
  document.body.removeChild(tempLink);
  // 释放blob URL地址
  window.URL.revokeObjectURL(blobURL);
  return filename;
};

export async function downloadFileByUrl(
  url: string,
  method: string,
  timeoutSeconds: number,
  data: any,
  params: any,
) {
  const res = await defHttp.request(
    {
      url: url,
      responseType: 'blob',
      data,
      params,
      timeout: timeoutSeconds * 1000,
      method: method,
    },
    { isReturnNativeResponse: true },
  );
  if (!res.headers['content-type'].includes('application/octet-stream')) {
    // 将 Blob 转换为 JSON
    const text = await new Promise<string>((resolve, reject) => {
      const reader = new FileReader();
      reader.onload = () => {
        resolve(reader.result as string);
      };
      reader.onerror = () => {
        reject(new Error('Blob读取失败'));
      };
      reader.readAsText(res.data);
    });

    let jsonResponse;
    try {
      jsonResponse = JSON.parse(text);
    } catch (e) {
      throw new Error('JSON解析失败');
    }
    // 抛出错误信息
    throw new Error(jsonResponse.msg || '未知错误');
  }

  const filename = downloadWithRes(res);
  return Promise.resolve(filename);
}

// 文件流导出方法
export async function downloadFileStream(
  apiUrl: string,
  params: Recordable,
  method: 'POST' | 'GET' = 'POST',
  timeoutSeconds = 600,
) {
  try {
    const res = await defHttp.request(
      {
        url: apiUrl,
        method,
        data: method === 'POST' ? params : undefined,
        params: method === 'GET' ? params : undefined,
        responseType: 'json',
        timeout: timeoutSeconds * 1000,
      },
      { isReturnNativeResponse: true },
    );

    const { code, data, msg } = res.data;
    if (code === 200 && data?.file_content) {
      const byteArray = base64ToUint8Array(data.file_content);
      downloadBlob(byteArray, data.file_type, data.file_name);
      return { success: true, fileName: data.file_name };
    }
    throw new Error(msg || '导出失败');
  } catch (e) {
    return { success: false, error: e.message };
  }
}

// 公共转换方法
const base64ToUint8Array = (base64: string) => {
  const byteCharacters = atob(base64);
  return new Uint8Array(Array.from(byteCharacters, (char) => char.charCodeAt(0)));
};

// 公共下载方法
const downloadBlob = (data: Uint8Array, type: string, filename: string) => {
  const blob = new Blob([data], { type });
  const url = URL.createObjectURL(blob);
  const a = document.createElement('a');
  a.href = url;
  a.download = filename;
  document.body.appendChild(a);
  a.click();
  URL.revokeObjectURL(url);
  document.body.removeChild(a);
};

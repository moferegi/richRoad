import configData from '../../common/config.js'
export const env = {
  ...configData,
  timeout: 87600,
  maxSize: 128 * 1024 * 1024 //上传文件的大小限制,128mb
};

// import configData from '../../common/config.js'
// var config = {
//   ...configData,
//   timeout: 87600,
//   maxSize: 128 * 1024 * 1024 //上传文件的大小限制,128mb
// };
// module.exports = config

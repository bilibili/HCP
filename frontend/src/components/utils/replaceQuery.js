export default function replaceQuery(str,username,formData,that) {
  return str.replace(/\$([a-zA-Z_]\w*)/g, (match, key) => {
    if (key === 'USER') {
      return username;
    }
    const value = findValueByKey(formData, key);
    return value !== null ? value : "ISvalue_null";
  });
}

function findValueByKey(data, key) {
  // 如果在当前层级找到了键，直接返回其值
  if (key in data) {
    return data[key];
  }
  // 否则，递归查找嵌套对象中的键
  for (const prop in data) {
    if (typeof data[prop] === 'object' && data[prop] !== null) {
      // 如果是数组，则遍历数组中的每个元素
      if (Array.isArray(data[prop])) {
        for (const item of data[prop]) {
          const value = findValueByKey(item, key);
          if (value !== undefined) {
            return value;
          }
        }
      } else {
        // 如果是对象，则递归查找
        const value = findValueByKey(data[prop], key);
        if (value !== undefined) {
          return value;
        }
      }
    }
  }
  // 如果键不存在于任何层级，返回 null
  return null;
}
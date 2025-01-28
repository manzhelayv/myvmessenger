import * as SecureStore from "expo-secure-store";

export function SaveSecureStore(key, value) {
  try {
    const data = SecureStore.setItem(key, value);
    return data;
  } catch (err) {
    return false;
  }
}

export function GetSecureStore(key) {
  let result = SecureStore.getItem(key);
  if (result) {
    return result;
  } else {
    return false;
  }
}

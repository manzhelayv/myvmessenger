import AsyncStorage from '@react-native-async-storage/async-storage';

export async function SetAsyncStorage(key, value) {
     try {
         await AsyncStorage.setItem(key, value);
     } catch (e) {
         console.log("SetAsyncStorage error:", e);
     }
}

export async function SetAsyncStorageObject(key, value) {
     try {
         const jsonValue = JSON.stringify(value);
         await AsyncStorage.setItem(key, jsonValue);
     } catch (e) {
         console.log("SetAsyncStorageObject error:", e);
     }
}

export async function GetAsyncStorageObject(key) {
    try {
        let promise = await AsyncStorage.getItem(key);
        return promise != null ? JSON.parse(promise) : null;
    } catch (e) {
        console.log("GetAsyncStorageObject", e)
    }
}

export async function GetAsyncStorage(key) {
    let promise = await AsyncStorage.getItem(key);
    return promise
}

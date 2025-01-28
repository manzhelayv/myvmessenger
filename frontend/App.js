import React, { useState, createContext } from 'react';
import { StatusBar } from 'expo-status-bar';
import { StyleSheet, View } from 'react-native';
import { GestureHandlerRootView } from "react-native-gesture-handler";
import TopMenu from './components/menu/TopMenu';
import Cookies from "universal-cookie";
import "./css/styles.css";
import BackgroundNotification from "./components/notification/BackgroundNotification";
import {SetAsyncStorage} from './components/storage/AsyncStorage';
import { useFonts } from 'expo-font';
import {styles} from "./css/App";

export const cookies =new Cookies(null, {path: '/'})


export const HTTP = "http://"

// export const HOST = "localhost"
export const HOST = "192.168.1.239"
export const HOST_SERVER = HTTP + HOST + ":30096";
export const HOST_CHAT = HTTP + HOST + ":30095";
export const WS_HOST = 'ws://'+ HOST + ":30012"

export const UserContext = createContext();

export default function App() {
  const [notification, setNotification] = useState('');
  let [userStatus, setUserStatus] = useState("");

  BackgroundNotification(notification, setNotification)

  SetAsyncStorage('pushMessage', 'false')
  SetAsyncStorage('statusBackground', 'foreground')
  SetAsyncStorage('notification', '')

  cookies.set('oldAmailActiveContact', '')

    const [loaded, error] = useFonts({
        'Wittgenstein-Italic-Regular': require('./assets/fonts/Wittgenstein-Italic-Regular.otf'),
        'Wittgenstein-Regular': require('./assets/fonts/Wittgenstein-Regular.otf'),
    });

 return (
   <>
     <GestureHandlerRootView style={styles.container}>
       <View style={styles.header}></View>
       <View>
           <UserContext.Provider value={{userStatus, setUserStatus, cookies, HOST_SERVER, HOST_CHAT, WS_HOST, notification, setNotification}}>
             <TopMenu />
           </UserContext.Provider>
       </View>
       <StatusBar style="light" />
     </GestureHandlerRootView>
   </>
 );
}



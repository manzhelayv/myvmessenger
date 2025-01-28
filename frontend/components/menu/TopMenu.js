import * as React from 'react';
import { useNavigation, createNavigationContainerRef, NavigationContainer } from '@react-navigation/native';
import { createMaterialTopTabNavigator } from '@react-navigation/material-top-tabs';
import { useState, useContext} from 'react';
import Chat from './chat/Chat';
import Chats from './chats/Chats';
import Contacts from './contacts/Contacts';
import Settings from './settings/Settings';
import Login from '../user/Login';
import Registration from '../user/Registration';
import {UserContext, WS_HOST} from '../../App';
import { GetSecureStore, SaveSecureStore } from '../storage/SecureStore';
import {SetAsyncStorage, GetAsyncStorage, SetAsyncStorageObject, GetAsyncStorageObject} from '../storage/AsyncStorage.js';
import * as FileSystem from "expo-file-system";
import Avatar from "./settings/Avatar";
import Profile from "./settings/Profile";
import Password from "./settings/Password";
import MyTabBar from  "./MyTabBar"
import {styles} from "./css/TopMenu";

const navigationRef = createNavigationContainerRef();
const Tab = createMaterialTopTabNavigator();

export let globalOnPress = ""
export let globalLatsMessageForChats = ""
export let globalOnPressNotification = ""
export let globalNotification = ""
export let globalOnPressWebsocket

export default function TopMenu() {
  const [restartWebsocket, setRestartWebsocket] = useState("");
  const [clientWS, setClientWS] = useState("");
  const [latsMessageForChats, setLatsMessageForChats] = useState([]);
  const {cookies, notification} = useContext(UserContext);
  const [key, setKey] = useState(1);
  const [messageForChats, setMessageForChats] = useState([]);
  const navigation = useNavigation();
  const [userStatuses, setUserStatuses] = useState("");
  SetAsyncStorage('activeTabChat', '')

  const waitForNavigationToBeReady = (navigationRef, avatar) => {
    const checkIfReady = (resolve) => {
      if (navigationRef.isReady()) {
        resolve(true);
        
        navigationRef.navigate('chat', {
          img: avatar
        })

        return
      }

      setTimeout(() => {
        checkIfReady(resolve);
      }, 25);
    };

    return new Promise(resolve => {
      checkIfReady(resolve);
    });
  };

  const setNewKey = async (email, name, userTo, userFrom, date, dateDay, dateTime, message, avatar, menu) => {
    SaveSecureStore('userTo', userTo)
    SaveSecureStore('nameActiveContact', name)
    SetAsyncStorage('nameActiveContact', name)

    cookies.set('nameActiveContact', name)
    cookies.set('emailActiveContact', email)

    let mess = []
    mess['email'] = email
    mess['name'] = name
    mess['user_to'] = userTo

    if (date !== "") {
      mess['date'] = date
      mess['date_day'] = dateDay
      mess['date_time'] = dateTime
      mess['message'] = message
      mess['user_from'] = userFrom
    } else {
      mess['date'] = latsMessageForChats['date']
      mess['date_day'] = latsMessageForChats['date_day']
      mess['date_time'] = latsMessageForChats['date_time']
      mess['message'] = ""
      mess['user_from'] = latsMessageForChats['user_from']
    }

    GetAsyncStorageObject('userTo' + userTo).then((user) => {
      if(user){
        avatar = user.avatar
        SetAsyncStorage('selectedImageGlobal', avatar)
        FileSystem.readAsStringAsync(user.avatarFile, { encoding: FileSystem.EncodingType.Base64 })
            .then(async (b64avatar) => {
              if (b64avatar) {
                mess['avatar'] = b64avatar
                if (b64avatar !== "") {
                  SetAsyncStorage('selectedImageGlobal', mess['avatar']);

                  if (navigationRef.isReady()) {
                    navigationRef.navigate('chat', {
                      img: mess['avatar']
                    })
                  } else {
                    waitForNavigationToBeReady(navigationRef, mess['avatar'])
                  }
                }

                const objUserTo = {email: email, name: name, userTo: userTo, userFrom: userFrom, message: message, date: date, avatar:mess['avatar'], avatarFile: user.avatarFile, countNotReadMessage: user.countNotReadMessage};
                SetAsyncStorageObject('userTo' + userTo, objUserTo)

                setLatsMessageForChats(mess)
              }
            }).catch(e => console.log('err', e));
      }
    });
  }

  const setNewMessage = async (dateTimestamp, date, dateDay, dateTime, email, message, name, userFrom, userTo, avatar, menu, countNotReadMessage) => {
    if (menu !== 'SendWSUserTo') {
      cookies.set('emailActiveContact', email)
    }

    let mess = []
    mess['date_timestamp'] = dateTimestamp
    mess['date'] = date
    mess['date_day'] = dateDay
    mess['date_time'] = dateTime
    mess['email'] = email
    mess['message'] = message
    mess['name'] = name
    mess['user_from'] = userFrom
    mess['user_to'] = userTo

    if (countNotReadMessage !== undefined) {
      mess['countNotReadMessage'] = countNotReadMessage
      await GetAsyncStorage('statusBackground').then(async (status) => {
        await GetAsyncStorageObject('userChatInfo' + userTo).then(async (userChatInfo) => {
          const userToChat = GetSecureStore('userTo')
          if (userChatInfo) {
            if ((status === 'foreground' || status === 'active') &&
                userTo === userChatInfo.userTo &&
                userChatInfo.userTo === userToChat &&
                cookies.get('activeTabChat') === 'chat'
            ) {
              let objUser = {
                userTo: userChatInfo.userTo,
                lastAppearance: userChatInfo.lastAppearance,
                countNotReadMessage: 0
              };

              mess['countNotReadMessage'] = 0
              await SetAsyncStorageObject('userChatInfo' + userChatInfo.userTo, objUser)
            }
          }
        })
      })
    }

    GetAsyncStorageObject('userTo' + userTo).then((user) => {
      if(user){
        avatar = user.avatar

        SetAsyncStorage('selectedImageGlobal', avatar)

        FileSystem.readAsStringAsync(user.avatarFile, { encoding: FileSystem.EncodingType.Base64 })
            .then(async (b64avatar) => {
              mess['avatar'] = b64avatar

              if (avatar !== "") {
                SetAsyncStorage('selectedImageGlobal', mess['avatar'])
              }
              const objUserTo = {email: email, name: name, userTo: userTo, userFrom: userFrom, message: message, date: date, avatar:mess['avatar'], avatarFile: user.avatarFile, countNotReadMessage: user.countNotReadMessage};
              SetAsyncStorageObject('userTo' + userTo, objUserTo)

              setLatsMessageForChats(mess)
              setMessageForChats(mess)
            }).catch(e => console.log('err', e));
      }
    });
  }

  const setConnectionWebsocket = async () => {
    const getUserFrom = GetSecureStore('user')
    let userFrom = ""
    if (getUserFrom !== false) {
      const userFromParse = JSON.parse(getUserFrom)
      userFrom = userFromParse[0].tdid
    }

    let URL =WS_HOST+'?user_from='+userFrom//+"&user_to="+userTo
    let client = new WebSocket(URL)

    setClientWS(client)

    setRestartWebsocket("false")
  }

  if (restartWebsocket === "") {
    setConnectionWebsocket()
  }

  globalOnPress = setNewKey
  globalLatsMessageForChats = latsMessageForChats
  globalOnPressNotification = setNewMessage
  globalNotification = notification
  globalOnPressWebsocket = setConnectionWebsocket

  let token = GetSecureStore('token')

  return(
      <NavigationContainer independent={true} ref={navigationRef}>
        <Tab.Navigator
            initialRouteName="chats"
            screenOptions={{
              headerShown: false,
              scrollEnabled: false,
              tabBarScrollEnabled: false,
              onPageScroll:false,
              swipeEnabled:false,
              animationEnabled:false,
            }}
            overScrollMode='never'

            style={[styles.menuContainer, {}]}
            tabBar={(props) => <MyTabBar {...props} UserContext={UserContext} userStatuses={userStatuses}/>}
        >
          {token === false ? (
                  <>
                    <Tab.Screen
                        name="login"
                        children={() =>  <Login setKey={setKey} UserContext={UserContext} setRestartWebsocket={setRestartWebsocket}/>}
                        options={{
                          tabBarLabel: 'Логин',
                          swipeEnabled:true,
                        }}
                    />
                    <Tab.Screen
                        name="registration"
                        children={() =>  <Registration UserContext={UserContext} setRestartWebsocket={setRestartWebsocket}/>}
                        options={{
                          tabBarLabel: 'Регистрация',
                          swipeEnabled:true,
                        }}
                    />
                  </>
              ) :
              <>
                <Tab.Screen
                    name="chats"
                    children={() => <Chats onPress={setNewKey} latsMessageForChats={latsMessageForChats}/>}
                    options={{
                      tabBarLabel: 'Чаты',
                      headerLeft: ()=> null,
                      gesturesEnabled: false,
                      headerShown: false,
                      scrollEnabled: false,
                      tabBarScrollEnabled: false,
                    }}
                />
                <Tab.Screen
                    name="chat"
                    children={() => <Chat onPress={setNewMessage} onPressNotification={setNewKey} setUserStatuses={setUserStatuses} clientWS={clientWS} latsMessageForChats={latsMessageForChats} />}
                    options={({ route }) => ({
                      tabBarLabel: 'Чат',
                      headerLeft: ()=> null,
                      gesturesEnabled: false,
                      headerShown: false,
                      scrollEnabled: false,
                      tabBarScrollEnabled: false,
                    })}
                />
                <Tab.Screen
                    name="contacts"
                    children={() => <Contacts />}
                    options={{
                      tabBarLabel: 'Контакты',
                      headerLeft: ()=> null,
                      gesturesEnabled: false,
                      headerShown: false,
                      scrollEnabled: false,
                      tabBarScrollEnabled: false,
                    }}
                />
                <Tab.Screen
                    name="settings"
                    children={() => <Settings />}
                    options={{
                      tabBarLabel: 'Настройки',
                      headerLeft: ()=> null,
                      gesturesEnabled: false,
                      headerShown: false,
                      scrollEnabled: false,
                      tabBarScrollEnabled: false,
                    }}
                />
                <Tab.Screen
                    name="avatar"
                    children={() => <Avatar />}
                    options={{
                      tabBarLabel: 'Аватар',
                      headerLeft: ()=> null,
                      gesturesEnabled: false,
                      headerShown: false,
                      scrollEnabled: false,
                      tabBarScrollEnabled: false,
                    }}
                />
                <Tab.Screen
                    name="profile"
                    children={() => <Profile />}
                    options={{
                      tabBarLabel: 'Профайл',
                      headerLeft: ()=> null,
                      gesturesEnabled: false,
                      headerShown: false,
                      scrollEnabled: false,
                      tabBarScrollEnabled: false,
                    }}
                />
                <Tab.Screen
                    name="password"
                    children={() => <Password />}
                    options={{
                      tabBarLabel: 'Пароль',
                      headerLeft: ()=> null,
                      gesturesEnabled: false,
                      headerShown: false,
                      scrollEnabled: false,
                      tabBarScrollEnabled: false,
                    }}
                />
              </>
          }
        </Tab.Navigator>
      </NavigationContainer>
  )
}

function newTopMenu(){
  // TopMenu(false, 2)
  // const root = createRoot(document.getElementById('root'));
  // root.render(<App />);
}


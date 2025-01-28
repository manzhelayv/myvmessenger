import React, {useEffect, useRef, useState} from 'react';
import { Text, View, ScrollView, AppState } from 'react-native';
import axios from 'axios';
import AddNewMessages from './AddNewMessages';
import WebSockets from './WebSockets';
import {GetSecureStore, SaveSecureStore} from '../../storage/SecureStore.js';
import BackgroundTimer from "react-native-background-timer";
import {GetAsyncStorage, GetAsyncStorageObject, SetAsyncStorageObject,} from '../../storage/AsyncStorage.js';
import * as FileSystem from "expo-file-system";
import {addSubscription} from "./notification/Notification";
import {reconnectingSocket } from "./notification/StartWS";
import {styles} from "./css/MessagesCss";
import ImageViewerSettings from "../../form/ImageViewerSettings";
import {getFileExtensionData} from "./Function";
import {cookies, WS_HOST, HOST_CHAT, HOST_SERVER} from '../../../App';

let mess = []
let blockMessage = ''

export default function Messages({ onPress, setUserStatuses, clientWS, latsMessageForChats}) {
    const [messages, setMessages] = useState([]);
    const [message, setMessage] = useState([]);
    const [lenMessages, setLenMessages] = useState(0);
    const [wsConnected, setWsConnected] = useState(true);
    const appState = useRef(AppState.currentState);
    const [appStateVisible, setAppStateVisible] = useState(appState.current);

    let token = GetSecureStore('token')
    const config = {
        headers: { Authorization: `Bearer ${token}` }
    };

    let emailActiveContact = cookies.get('emailActiveContact')
    let oldAmailActiveContact = cookies.get('oldAmailActiveContact')
    const userTo = GetSecureStore('userTo')

    async function loadData() {
        axios.get(HOST_CHAT + '/chat?user_to=' + userTo, config)
            .then(async function (response) {
                cookies.set('oldAmailActiveContact', emailActiveContact)

                mess = response.data.messages

                async function Messages(mess) {
                    for (const value of mess) {
                        if (value.file !== '') {

                            let urlFile = FileSystem.documentDirectory + value.file

                            let file = await FileSystem.getInfoAsync(urlFile);

                            if (file.exists === true) {
                                let base64 =  await FileSystem.readAsStringAsync(urlFile, {encoding: FileSystem.EncodingType.Base64})
                                value.message = base64
                            } else {
                                value.message = 'Файл удален'
                                value.file = ''
                            }
                        }
                    }

                    return mess
                }

                let newMess = await Messages(mess)

                setMessages(newMess)
                setLenMessages(newMess.length)
                setMessage([])
            })
            .catch(function (error) {
                setMessages([])
                setLenMessages(0)
                setMessage([])
                setIsset(false)
                if (error.request !== undefined) {
                    const err = JSON.parse(error.request.response);
                    if (err.error.code === 404) {
                        cookies.set('oldAmailActiveContact', emailActiveContact)
                    }
                }
            });
    }

    const getUserFrom = GetSecureStore('user')
    let userFrom = ""
    if (getUserFrom !== false) {
        const userFromParse = JSON.parse(getUserFrom)
        userFrom = userFromParse[0].tdid
    }


    GetAsyncStorage('statusBackground').then((status) => {
        GetAsyncStorageObject('userChatInfo' + userTo).then((userChatInfo) => {
            if (userChatInfo) {
                if ((status === 'foreground' || status === 'active') &&
                    userTo === userChatInfo.userTo &&
                    cookies.get('activeTabChat') === 'chat' &&
                    userChatInfo.countNotReadMessage !== 0 &&
                    (latsMessageForChats['user_to'] === userTo)
                ) {
                    let objUser = {
                        userTo: userChatInfo.userTo,
                        lastAppearance: userChatInfo.lastAppearance,
                        countNotReadMessage: 0
                    };

                    cookies.set('countNotReadMessage' + userChatInfo.userTo, 0)

                    SetAsyncStorageObject('userChatInfo' + userChatInfo.userTo, objUser)

                    GetAsyncStorageObject("userTo" + userChatInfo.userTo).then((userToObj) => {
                        userToObj.countNotReadMessage = 0
                        SetAsyncStorageObject("userTo", userToObj)

                        onPress('', userToObj.date, "", userToObj.date, userToObj.email, userToObj.message, userToObj.name, userToObj.userFrom, userToObj.userTo, userToObj.avatar, 'SendWSUserTo', 0);
                    })
                }
            }
        })
    })

    let URL ='ws://'+WS_HOST+'?user_from='+userFrom//+"&user_to="+userTo
    let client = ''

    client = reconnectingSocket(URL, message, setMessage, onPress, cookies, userTo, setUserStatuses, clientWS)

    addSubscription(appState, appStateVisible, setAppStateVisible, cookies)

    if (cookies.get('wsConnected') !== true) {
        cookies.set('wsConnected', true)
        BackgroundTimer.setInterval(() => {
            client.getClient().submitMessagePing()
            addSubscription(appState, appStateVisible, setAppStateVisible, cookies)

        }, 5000);
    }

    if ((emailActiveContact !== oldAmailActiveContact || oldAmailActiveContact === '' || emailActiveContact === '')
         && emailActiveContact !== undefined && emailActiveContact !== "undefined") { // && client.isConnected() === true
        loadData();
        return
    }

    return (
        <View style={{ flex: 1, width: '100%' }}>
            <View style={{ height: '85%'}}>
                    <ScrollView
                        indStyle={styles.scrollView}
                        className="classN"
                        ref={ref => {this.scrollView = ref}}
                        onContentSizeChange={() => this.scrollView.scrollToEnd({animated: false})}
                    >
                        {messages.map((value, index) => {
                            let paddingTop = 0
                            if (index == 0) {
                                paddingTop = 10
                            }

                            if (value.file !== '') {
                                let fileData = getFileExtensionData(value.file)

                                blockMessage = <>
                                                 <Text style={styles.buttonTextDate}>{value.date}</Text>
                                                 <View style={styles.leftImage}>
                                                    <ImageViewerSettings placeholderImageSource="" selectedImage={fileData + value.message} />
                                                 </View>
                                                 <View style={styles.rightTextDate}>
                                                    <Text style={styles.buttonTextDate}>{value.date_time}</Text>
                                                 </View>
                                                </>
                            } else {
                                blockMessage = <>
                                                 <Text style={styles.buttonTextDate}>{value.date}</Text>
                                                 <View style={styles.leftText}>
                                                    <Text style={styles.buttonText}>{value.message}</Text>
                                                 </View>
                                                 <View style={styles.rightTextDate}>
                                                    <Text style={styles.buttonTextDate}>{value.date_time}</Text>
                                                 </View>
                                               </>
                            }

                            if(value.user_from === userFrom){
                                return (
                                    <View style={{display: 'block', paddingTop: paddingTop}} key={index}>
                                        <View style={styles.rightChatText}>
                                            {blockMessage}
                                        </View>
                                    </View>
                                )
                            } else {
                                return (
                                    <View style={{display: 'block', paddingTop: paddingTop}} key={index}>
                                        <View style={styles.leftChatText}>
                                            {blockMessage}
                                        </View>
                                    </View>
                                )
                            }
                        })}
                            <AddNewMessages message={message} styles={styles}/>
                    </ScrollView>
            </View>
            <WebSockets
                message={message}
                setMessage={setMessage}
                onPress={onPress}
                wsConnected={wsConnected}
                client={client}
                userTo={userTo}
            />
        </View>
    );
}


import {GetAsyncStorage, GetAsyncStorageObject, SetAsyncStorage, SetAsyncStorageObject} from "../../../storage/AsyncStorage";
import * as Notifications from "expo-notifications";
import BackgroundTimer from "react-native-background-timer";
import * as FileSystem from "expo-file-system";
import {GetSecureStore} from "../../../storage/SecureStore";
import {cookies} from "../../../../App";
import {globalOnPressWebsocket} from "../../TopMenu";

const sleep = (time) => new Promise((resolve) => setTimeout(() => resolve(), time));

export function reconnectingSocket(url, message, setMessage, onPress, cookies, userToTo, setUserStatuses, client) {
    let onPressWS = globalOnPressWebsocket
    let isConnected = false;
    let reconnectOnClose = true;
    let messageListeners = [];
    let stateChangeListeners = [];

    function on(fn) {
        messageListeners.push(fn);
    }

    function off(fn) {
        messageListeners = messageListeners.filter(l => l !== fn);
    }

    function onStateChange(fn) {
        stateChangeListeners.push(fn);
        return () => {
            stateChangeListeners = stateChangeListeners.filter(l => l !== fn);
        };
    }

    function start() {
        client.submitMessagePing = () => {
            console.log('ping')
            try {
                GetAsyncStorage('statusBackground').then((status) => {
                    if (status === 'foreground') {
                        let messageObj = {"message": "pingWS123443321", "file": '', "status": "online"}
                        client.send(JSON.stringify(messageObj));
                    } else {
                        let status = "" + cookies.get('backgroundDate') +  cookies.get('backgroundDateTime') + "";
                        let messageObj = {"message": "pingWS123443321", "file": '', "status": status}
                        client.send(JSON.stringify(messageObj));
                    }
                }).catch(e => console.log('err', e));
            } catch (error) {
                console.log("sendWS", error);
            }
        };

        client.submitMessage = (message) => {
            try {
                client.send(JSON.stringify(message));
            } catch (error) {
                console.log("sendWS", error);
            }
        };

        client.onopen = () => {
            isConnected = true;
            stateChangeListeners.forEach(fn => fn(true));
        }

        const close = client.close;

        client.close = () => {
            reconnectOnClose = false;
            close.call(client);
        }

        client.onmessage = (e) => {
            let data = JSON.parse(e.data);
            data = JSON.parse(data);

            if (data.Message === 'pongWS123443321') {
                if (data.UserTo === userToTo) {
                    SetAsyncStorage('statusUser' + data.UserTo, data.Status)
                    cookies.set('statusUser' + data.UserTo, data.Status)
                }

                setUserStatuses(cookies.get('statusUser' + userToTo))

                client.onopen()
                return
            }

            if (data.File !== '') {
                FileSystem.writeAsStringAsync(FileSystem.documentDirectory + data.File, data.Message, { encoding: FileSystem.EncodingType.Base64 })
                    .then(async (uri) => {

                    })
                    .catch(e => console.log('err', value.user_to, e));

                data.Message = "file"
            }

            GetAsyncStorageObject('userChatInfo' + data.UserTo).then((userChatInfo) => {
                if (!userChatInfo) {
                    let objUser = {userTo: data.UserTo, lastAppearance: 'online', countNotReadMessage: 0};

                    SetAsyncStorageObject('userChatInfo' + data.UserTo, objUser)
                }
            })

            let userTo = GetSecureStore('userTo')
            GetAsyncStorageObject('userTo' + data.UserTo).then((userToObj) => {
                SetAsyncStorage('nameActiveContact', userToObj.name)

                let notificationMessages = ""
                GetAsyncStorage('statusBackground').then((status) => {
                    GetAsyncStorageObject('userChatInfo' + data.UserTo).then((userChatInfo) => {
                        if (userChatInfo) {
                            let objUser = {userTo: userChatInfo.userTo, lastAppearance: userChatInfo.lastAppearance, countNotReadMessage: userChatInfo.countNotReadMessage + 1};

                            if ((status === 'foreground') && userTo === userChatInfo.userTo && cookies.get('activeTabChat') !== 'chat') {
                                if (userToObj.countNotReadMessage !== undefined) {
                                    userToObj.countNotReadMessage += 1
                                } else {
                                    userToObj.countNotReadMessage = 1
                                }
                                SetAsyncStorageObject('userChatInfo' + userChatInfo.userTo, objUser)
                                SetAsyncStorageObject('userTo' + data.UserTo, userToObj)
                            } else {
                                if (userToObj.countNotReadMessage !== undefined) {
                                    userToObj.countNotReadMessage += 1
                                } else {
                                    userToObj.countNotReadMessage = 1
                                }
                                SetAsyncStorageObject('userChatInfo' + userChatInfo.userTo, objUser)
                                SetAsyncStorageObject('userTo' + data.UserTo, userToObj)
                            }

                        }
                    }).catch(function (error) {})

                    if (status !== 'foreground') {
                        GetAsyncStorage('notification').then((messages) => {
                            if (messages === null) {
                                notificationMessages = data.Message
                            } else {
                                notificationMessages = messages + "\n" + data.Message
                            }

                            if (notificationMessages !== '') {
                                    GetAsyncStorage('notificationIdentifier' + userToObj.userTo).then((identifier) => {
                                        if (identifier) {
                                            Notifications.dismissNotificationAsync(identifier)
                                        }
                                    }).catch(e => console.log('notificationIdentifier err', e));

                                    let d = new Date()
                                    let minutes = (d.getMinutes() < 10 ? '0' : '') + d.getMinutes()
                                    let hours = d.getHours() + ":" + minutes;

                                    let date = hours
                                    const notificationIdentifier = Notifications.scheduleNotificationAsync({
                                        content: {
                                            title: userToObj.name,
                                            body: notificationMessages,
                                            data: {
                                                userTo: userToObj.userTo,
                                                date: userToObj.date,
                                                dateTime: date,
                                                message: data.Message,
                                                url: 'chat'
                                            },
                                            attachments: {}
                                        },
                                        trigger: {seconds: 1, repeats: false}
                                    })

                                    notificationIdentifier.then((identifier) => {
                                        SetAsyncStorage('notificationIdentifier' + userToObj.userTo, identifier).then((res) => {

                                        }).catch(e => console.log('notificationIdentifier err', e));
                                    }).catch(e => console.log('notificationIdentifier2 err', e));
                            }

                            SetAsyncStorage('notification', notificationMessages)
                        })
                    }
                })

                let msg = []
                if (message.length > 0) {
                    message.map((value, index) => {
                        msg.push(value)
                    })
                }

                let messageObj = {}
                if (data.Message !== '' && data.Message !== 'file') {
                    messageObj = {"message": data.Message, "file": '', "userTo": data.UserTo}
                } else {
                    messageObj = {"message": "", "file": data.File, "userTo": data.UserTo}
                }

                let left = []
                left['left'] = messageObj
                msg.push(left)

                if (data.UserTo === userTo) {
                    setMessage(msg)
                }

                let d = new Date()
                let minutes = (d.getMinutes() < 10 ? '0' : '') + d.getMinutes()
                let hours = d.getHours() + ":" + minutes;

                let date = hours
                let dateTimestamp = Date.now();

                GetAsyncStorageObject('userChatInfo' + data.UserTo).then((userChatInfo) => {
                    cookies.set('countNotReadMessage' + data.UserTo, userChatInfo.countNotReadMessage + 1)
                    onPress(dateTimestamp, userToObj.date, "", date, userToObj.email, data.Message, userToObj.name, userToObj.userFrom, userToObj.userTo, userToObj.avatar, 'SendWSUserTo', userChatInfo.countNotReadMessage + 1);
                })
            })
        }

        client.onerror = (e) => console.error(e);

        client.onclose = () => {
            isConnected = false;
            stateChangeListeners.forEach(fn => fn(false));

            if (!reconnectOnClose) {
                return;
            }

            let intervalId = BackgroundTimer.setInterval(() => {
                onPressWS()

                BackgroundTimer.clearInterval(intervalId);
            }, 5000);
        }
    }

    start();

    return {
        on,
        off,
        onStateChange,
        close: () => client.close(),
        getClient: () => client,
        isConnected: () =>isConnected,
    };
}
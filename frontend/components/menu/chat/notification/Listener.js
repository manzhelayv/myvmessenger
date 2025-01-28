import * as Notifications from "expo-notifications";
import {useEffect, useRef} from "react";
import {GetAsyncStorageObject, SetAsyncStorage, SetAsyncStorageObject} from "../../../storage/AsyncStorage";
import * as FileSystem from "expo-file-system";
import Cookies from "universal-cookie";
import {navigate} from "expo-router/build/global-state/routing";

const cookies = new Cookies(null, { path: '/' });

Notifications.setNotificationHandler({
    handleNotification: async () => ({
        shouldShowAlert: true,
        shouldPlaySound: false,
        shouldSetBadge: false,
    }),
});

const registerForPushNotificationsAsync = async () => {
    const { status: existingStatus } = await Notifications.getPermissionsAsync();
    let finalStatus = existingStatus;

    if (existingStatus !== 'granted') {
        // if we dontt have access to it, we ask for it
        const { status } = await Notifications.requestPermissionsAsync();
        finalStatus = status;
    }
    if (finalStatus !== 'granted') {
        return;
    }

    const projectId = '0f541b61-a9ad-4cfc-a5e0-639d27ec2119'
    let pushTokenString = ""
    try {
        pushTokenString = (
            await Notifications.getExpoPushTokenAsync({
                projectId,
            })
        ).data;
    } catch (e) {
        console.log("expo-notifications err", e);
    }

    cookies.set('pushTokenString', pushTokenString);

    Notifications.setNotificationChannelAsync('default', {
        name: 'default',
        importance: Notifications.AndroidImportance.MAX,
        vibrationPattern: [0, 250, 250, 250],
        lightColor: '#FF231F7C',
    });

    return pushTokenString;
}

export function StatefullApp (onPressNotification) {
    const notificationListener = useRef();
    const responseListener = useRef();

    useEffect(() => {
        // Register for push notification
        const token = registerForPushNotificationsAsync();
        // This listener is fired whenever a notification is received while the app is foregrounded
        notificationListener.current = Notifications.addNotificationReceivedListener(notification => {
            notificationCommonHandler(notification);
        });

        // This listener is fired whenever a user taps on or interacts with a notification
        // (works when app is foregrounded, backgrounded, or killed)
        responseListener.current = Notifications.addNotificationResponseReceivedListener(response => {
            notificationCommonHandler(response.notification);
            notificationNavigationHandler(response.notification.request.content, onPressNotification);
        });
    }, []);

    const notificationCommonHandler = (notification) => {
        // save the notification to reac-redux store
        console.log('A notification has been received', notification)
    }

    const notificationNavigationHandler = (data, onPressNotification) => {
        // navigate to app screen
        if (data.data.userTo === undefined) {
            return
        }

        GetAsyncStorageObject('userTo' + data.data.userTo).then((userTo) => {
            SetAsyncStorage('nameActiveContact', userTo.name)

            const objUserTo = {
                email: userTo.email,
                name: userTo.name,
                userTo: userTo.userTo,
                userFrom: userTo.userTo,
                message: data.data.message,
                date: data.data.dateTime,
                avatar: FileSystem.documentDirectory + userTo.userTo + '.jpeg',
                avatarFile: userTo.avatarFile,
                countNotReadMessage: userTo.countNotReadMessage
            };
            SetAsyncStorageObject('userTo' + userTo.user_to, objUserTo)
            SetAsyncStorage('activeContactNotification', data.data.userTo)
            onPressNotification(userTo.email, userTo.name, userTo.userTo, userTo.userTo, data.data.date, "", data.data.dateTime, data.data.message, "")
        })
    }
}

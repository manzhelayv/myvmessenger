import {AppState} from 'react-native';
import { SetAsyncStorage } from '../../../storage/AsyncStorage.js';
import {GetStringMonth} from "../Function";

export function addSubscription (appState, appStateVisible, setAppStateVisible, cookies) {
    let statusBackground = ""

    // useEffect(() => {
    const subscription = AppState.addEventListener('change', nextAppState => {
        if (
            appState.current.match(/inactive|background/) &&
            nextAppState === 'active'
        ) {
            setAppStateVisible('background');
        } else {
            let d = new Date()
            let minutes = (d.getMinutes() < 10 ? '0' : '') + d.getMinutes()
            let hours = d.getHours() + ":" + minutes;

            let dateTime = hours

            let date = d.getDate() + " " + GetStringMonth(d.getMonth()) + " " + d.getFullYear() + " г."

            cookies.set('backgroundDate', date);
            cookies.set('backgroundDateTime', dateTime);

            setAppStateVisible('foreground');
        }

        cookies.set('statusBackground', appStateVisible);

        SetAsyncStorage('statusBackground', appStateVisible);

        appState.current = nextAppState;
        statusBackground  = appStateVisible
    });

    // subscription.remove();

    //     return () => {
    //         subscription.remove();
    //     };
    // }, [appStateVisible, setAppStateVisible]);

    return statusBackground
}
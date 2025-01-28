import * as React from 'react';
import {View, ImageBackground} from 'react-native';
import Messages from './Messages';
import {StatefullApp} from './notification/Listener';
import {styles} from "./css/ChatCss";
import {cookies, WS_HOST} from '../../../App';
import {GetSecureStore} from "../../storage/SecureStore";

const image =  require("../../../images/126.jpg");

export default function Chat({onPress, onPressNotification, setUserStatuses, clientWS, latsMessageForChats}) {
    StatefullApp(onPressNotification)
    cookies.set('active_menu', 'Chat')
    setUserStatuses(cookies.get('statusUser' + GetSecureStore('userTo')))

    return (
      <View style={styles.containerTop}>
          <ImageBackground source={image} style={styles.image}>
            <View style={styles.container}>
              <Messages onPress={onPress} setUserStatuses={setUserStatuses} clientWS={clientWS} latsMessageForChats={latsMessageForChats}/>
            </View>
          </ImageBackground>
      </View>
    );
}

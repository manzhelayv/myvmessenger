import React, { useState } from 'react';
import { View, TextInput } from 'react-native';
import Button from '../../form/Button';
import {GetAsyncStorageObject, SetAsyncStorageObject,} from '../../storage/AsyncStorage.js';
import {GetSecureStore} from "../../storage/SecureStore";
import {styles} from "./css/WebSocketsCss";
import * as ImagePicker from "expo-image-picker";
import * as FileSystem from "expo-file-system";
import {generateFileName} from "./Function";

export default function WebSockets({message, setMessage, onPress, wsConnected, client, userTo}) {
  const [messageText, setMessageText] = React.useState('');

  let topButton = -25
  let topTextInput = 22
  const [topView, setTopView] = useState(0);

  const handleOnFocus = () => {
    setTopView(-25)
  };

  const handleOnBlur = () => {
    setTopView(0)
  };

  if (wsConnected === true) {
    const sendMessage = (image, filename) => {
      if (messageText === '' && typeof image === 'object') {
        return
      }

      let messageObj = {}
      if (messageText !== '') {
        messageObj = {"userTo": userTo, "message": messageText, "file": '', "status": "online"}
      } else {
        messageObj = {"userto": userTo, "message": image, "file": filename, "status": "online"}
      }

      client.getClient().submitMessage(messageObj);

      setMessageText('')

      let msg = []
      if (message.length > 0) {
        message.map((value, index) => {
          msg.push(value)
        })
      }

      let right = []
      right['right'] = messageObj
      msg.push(right)

      let d = new Date()
      let minutes = (d.getMinutes() < 10 ? '0' : '') + d.getMinutes()
      let hours = d.getHours() + ":" + minutes;

      let date = hours
      let dateTimestamp = Date.now();

      let userToId =  GetSecureStore('userTo')
      GetAsyncStorageObject('userTo' + userToId).then((userTo) => {
        onPress(dateTimestamp, userTo.date, "", date, userTo.email, messageText, userTo.name, userTo.userFrom, userTo.userTo, userTo.avatar);
      })

      setMessage(msg)
    }

    const pickImageAsync = async () => {
      let result = await ImagePicker.launchImageLibraryAsync({
        quality: 1,
      });

      let filename = ""

      await FileSystem.readAsStringAsync(result.assets[0].uri, { encoding: FileSystem.EncodingType.Base64 })
          .then(async (b64avatar) => {
            let fileName = result.assets[0].uri
            filename = generateFileName(fileName, 12)
            await FileSystem.writeAsStringAsync(FileSystem.documentDirectory + filename, b64avatar, { encoding: FileSystem.EncodingType.Base64 })
                .then(async (uri) => {
                  sendMessage(b64avatar, filename)
                })
                .catch(e => console.log('err', value.user_to, e));
          })
          .catch(e => console.log('err', value.user_to, e));
    };

    return (
        <View style={[styles.containerButton, {top: topView}]}>
          <View style={styles.addContact} key={1}>
            <TextInput
                style={[styles.textArea, {top: topTextInput}]}
                underlineColorAndroid="transparent"
                placeholderTextColor="grey"
                numberOfLines={10}
                multiline={true}
                onChangeText={(text) => setMessageText(text)}
                value={messageText}
                onFocus={() => handleOnFocus()}
                onBlur={() => handleOnBlur()}
            />
          </View>
          <View style={[styles.addButton, {top: topButton}]} key={2}>
            <Button theme="send" label="Отправить" onPress={sendMessage}/>
            <View>
              <Button theme="primary" label="Фото" onPress={pickImageAsync} />
            </View>
          </View>
        </View>
    );
  }
}

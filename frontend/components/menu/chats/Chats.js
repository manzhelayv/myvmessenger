import * as React from 'react';
import {Pressable, Text, View} from 'react-native';
import axios from 'axios';
import {useState} from 'react';
import {ScrollViewIndicator} from '@fanchenbao/react-native-scroll-indicator';
import ImageViewer from '../../form/ImageViewer';
import * as FileSystem from "expo-file-system";
import {GetSecureStore} from '../../storage/SecureStore.js';
import { SetAsyncStorageObject, GetAsyncStorageObject } from "../../storage/AsyncStorage";
import {styles} from "./css/ChatsCss";
import {cookies, HOST_CHAT} from '../../../App';
import Animated from "react-native-reanimated";
import Entypo from "@expo/vector-icons/Entypo";
import {OpportunityAttachmentView} from "../PopupMenu";

let userChats = []

export default function Chats({onPress, latsMessageForChats}) {
    const [chatsLenght, setChatsLenght] = useState(0);
    const [countNotReadMessageLenght, setCountNotReadMessageLenght] = useState(0);
    const [chats, setChats] = useState([]);

    cookies.set('active_menu', 'Chats')

    const token = GetSecureStore('token')
    const config = {
      headers: { Authorization: `Bearer ${token}` }
    };

    async function loadData() {
        axios.get(HOST_CHAT + '/chat/chats', config)
        .then(async function (response) {
            userChats = response.data

            setChatsLenght(1)
            setChats(userChats)
        })
        .catch(function (error) {
            console.log(error)
        });
    }

    if (chatsLenght === 0) {
        loadData();
    }

    let newChats = true

    chats.map((value, index) => {
        GetAsyncStorageObject('userChatInfo' + value.user_to).then( async (userChatInfo) => {
            if (userChatInfo) {
                cookies.set('countNotReadMessage' + value.user_to, userChatInfo.countNotReadMessage)
            }
        })
    })

    chats.map((value, index) => {
        if (latsMessageForChats['user_to'] !== undefined && latsMessageForChats['user_to'] === value.user_to && latsMessageForChats['message'] !== ''){
            chats[index] = latsMessageForChats

            newChats = false
        }

        chats[index].countNotReadMessage = cookies.get('countNotReadMessage' + value.user_to)
    })

    if (newChats === true && latsMessageForChats['user_to'] !== undefined && latsMessageForChats['message'] !== '') {
        chats[chats.length] = latsMessageForChats
    }

    chats.sort(function (a, b) {
        if (a.date_timestamp < b.date_timestamp) {
            return 1;
          }
          if (a.date_timestamp > b.date_timestamp) {
            return -1;
          }
          return 0;
    })

    return (
        <View style={styles.container}>
            <View style={{width:'100%'}}>
                <OpportunityAttachmentView />
            </View>
            <View
                accessibilityRole="button"
                accessibilityState={{ selected: true }}
                tabBarIndicatorStyle= {{backgroundColor: "#03fc5e", top:50, zIndex:-9999}}
                style={styles.viewButton}
            >
                <Entypo
                    name="chat"
                    size={18}
                    color="#25292e"
                    style={styles.buttonIconChat}
                />
                <Animated.Text style={styles.tabBarLabelStyle}>Чаты</Animated.Text>
            </View>

            <ScrollViewIndicator indStyle={styles.scrollView}>
                {chats.map((value, index) => {
                    let day = ""
                    if (value.date_time === "0") {
                        day = <Text style={styles.buttonTextRight} >{value.date_day}</Text>
                    } else if (value.date_day === "") {
                        day = <Text style={styles.buttonTextRightDate} >{value.date_time}</Text>
                    } else {
                        day = <Text style={styles.buttonTextRightDate} >{value.date_day}</Text>
                    }

                   let top = index === 0 ? 15 : 0

                    FileSystem.writeAsStringAsync(FileSystem.documentDirectory + value.user_to + '.jpeg', value.avatar, { encoding: FileSystem.EncodingType.Base64 })
                        .then(async (uri) => {

                        })
                        .catch(e => console.log('err', value.user_to, e));

                    FileSystem.readAsStringAsync(FileSystem.documentDirectory + value.user_to + '.jpeg', { encoding: FileSystem.EncodingType.Base64 })
                        .then(async (b64avatar) => {

                    }).catch(e => console.log('err', e));

                    if (value.message === '' && value.file !== '' && value.file !== undefined) {
                        value.message = 'Файл'
                    }

                    const objUserTo = {email: value.email, name: value.name, userTo: value.user_to, userFrom: value.user_to, message: value.message, date: value.date_time, avatar: FileSystem.documentDirectory + value.user_to + '.jpeg', avatarFile: FileSystem.documentDirectory + value.user_to + '.jpeg', countNotReadMessage: value.countNotReadMessage};

                    SetAsyncStorageObject('userTo' + value.user_to, objUserTo)

                    let img = 'data:image/jpeg;base64,' + value.avatar
                    let image = ""

                    if (img !== "") {
                        image = <ImageViewer placeholderImageSource="" selectedImage={img} size="25"/>
                    }
                   return <View style={[styles.leftChatText, {top:top}]} key={value.user_to} >
                        <Pressable
                            onPress={()=> onPress(value.email, value.name, value.user_to, value.user_from, value.date, value.date_day, value.date_time, value.message, value.avatar, 'Chats')}
                        >
                            {image}
                            <View style={styles.rightText}>
                                <View>
                                    <Text style={styles.buttonNameText}>{value.name}</Text>
                                </View>

                                <View style={styles.buttonTextView}>
                                     <Text style={styles.buttonText}>{value.message}</Text>
                                </View>

                                {value.countNotReadMessage === 0 || value.countNotReadMessage === undefined? (
                                    <View style={styles.dateRight}>
                                        {day}
                                    </View>
                                ) :
                                    <></>
                                }

                                {value.countNotReadMessage !== undefined && value.countNotReadMessage !== 0 ? (
                                    <View style={styles.dateRightNotReadMessageView}>
                                        <Text style={{color: '#21912b'}}>
                                            {day}
                                        </Text>
                                        <View style={styles.countNotReadMessageView}>
                                            <Text style={styles.countNotReadMessageText}>
                                                {value.countNotReadMessage}
                                            </Text>
                                        </View>
                                    </View>
                                    ) :
                                    <></>
                                }
                            </View>
                        </Pressable>
                    </View>
                })}
            </ScrollViewIndicator>
            <View style={{height:50, paddingTop: 10}}></View>
        </View>
    );
  }










  
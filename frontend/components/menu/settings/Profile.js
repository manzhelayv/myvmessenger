import { StatusBar } from 'expo-status-bar';
import {Text, View} from 'react-native';
import { GestureHandlerRootView } from "react-native-gesture-handler";
import Button from '../../form/Button';
import {useState} from 'react';
import * as MediaLibrary from 'expo-media-library';
import axios from 'axios';
import {GetSecureStore, SaveSecureStore} from "../../storage/SecureStore";
import { HOST_SERVER} from '../../../App';
import * as React from "react";
import {Input} from "react-native-elements";
import {styles} from "./css/Profile";
import {Result} from './Settings';
import MaskInput from "react-native-mask-input";

export default function Profile() {
    const [status, requestPermission] = MediaLibrary.usePermissions();
    const userInfo = GetSecureStore('user')
    const [subject, setSubject] = useState("");
    const [color, setColor] = useState("green");
    const [visibility, setVisibility] = useState(true);

    let emailUser = ''
    let nameUser = ''
    let loginUser = ''
    let phoneUser = ''
    let tokenUser = ''
    let tdidUser = ''
    if (userInfo !== false) {
        const userInfoParse = JSON.parse(userInfo)
        emailUser = userInfoParse[0].email
        nameUser = userInfoParse[0].name
        loginUser = userInfoParse[0].login
        phoneUser = userInfoParse[0].phone
        tokenUser = userInfoParse[0].access_token
        tdidUser = userInfoParse[0].tdid
    }

    const [email, setEmail] = useState(emailUser);
    const [name, setName] = useState(nameUser);
    const [login, setLogin] = useState(loginUser);
    const [phone, setPhone] = useState(phoneUser);

    const params = {
        email: email,
        name: name,
        login: login,
        phone: phone,
        tdid: tdidUser,
    }

    let token = GetSecureStore('token')
    const config = {
        headers: { Authorization: `Bearer ${token}` },
    };

    if (status === null) {
        requestPermission();
    }

    const onSaveProfile = async () => {
        axios.put(HOST_SERVER + '/user', params, config)
            .then(function (response) {

                const user = [
                    {
                        'token': tokenUser,
                        'login': login,
                        'name': name,
                        'email': email,
                        'phone': phone,
                        'tdid': tdidUser,
                    },
                ];

                SaveSecureStore('user', JSON.stringify(user))

                setColor("green")
                setSubject("Данные успешно обновлены")

                setVisibility(true)
                setTimeout(() => {
                    setVisibility(false)
                }, 5000);

            }).catch(function (error) {
                setColor("red")
                const err = JSON.parse(error.request.response);
                setSubject(err.error.message)

                setVisibility(true)
                setTimeout(() => {
                    setVisibility(false)
                }, 5000)
        });
    }

    return (
        <View style={styles.menuContainer}>
            <View style={styles.menuContainer}>
                <GestureHandlerRootView style={styles.container}>
                    { visibility &&
                        <Result subject={subject} color={color} top={5} bottom={10} />
                    }
                    <View style={styles.viewText}>
                        <Text>
                            Имя
                        </Text>
                    </View>
                    <Input
                        placeholder="Имя"
                        onChangeText={ (nameText) => setName(nameText)}
                        value={name}
                        containerStyle={styles.inputStyle}
                    />
                    <View style={styles.viewText}>
                        <Text>
                            Логин
                        </Text>
                    </View>
                    <Input
                        placeholder="Логин"
                        onChangeText={ (loginText) => setLogin(loginText)}
                        value={login}
                        containerStyle={styles.inputStyle}
                    />
                    <View style={styles.viewText}>
                        <Text>
                            Email
                        </Text>
                    </View>
                    <Input
                        placeholder="Email"
                        onChangeText={ (emailText) => setEmail(emailText)}
                        value={email}
                        containerStyle={styles.inputStyle}
                    />
                    <View style={styles.viewText}>
                        <Text>
                            Телефон
                        </Text>
                    </View>
                    <MaskInput
                        placeholder="Телефон"
                        placeholderTextColor={'#86939E'}
                        value={phone}
                        onChangeText={ (phoneText) => setPhone(phoneText)}
                        mask={['+',/\d/, ' ', /\d/, /\d/, /\d/, ' ', /\d/, /\d/, /\d/, ' ',  /\d/, /\d/, /\d/, /\d/]}
                        style={styles.inputMaskStyle}
                    />
                    <View style={styles.footerContainer}>
                        <Button theme="save" label="Сохранить" onPress={onSaveProfile} />
                    </View>
                    <StatusBar style="light" />
                </GestureHandlerRootView>
            </View>
        </View>
    );
}

import { StatusBar } from 'expo-status-bar';
import { View, Text} from 'react-native';
import { GestureHandlerRootView } from "react-native-gesture-handler";
import { useState } from 'react';
import * as MediaLibrary from 'expo-media-library';
import axios from 'axios';
import {GetSecureStore} from "../../storage/SecureStore";
import {HOST_SERVER} from '../../../App';
import * as React from "react";
import {Input} from "react-native-elements";
import Button from "../../form/Button";
import {styles} from "./css/Password";
import {Result} from './Settings';

export default function Password() {
    const [status, requestPermission] = MediaLibrary.usePermissions();
    const [oldPassword, setOldPassword] = useState("");
    const [newPassword, setNewPassword] = useState("");
    const [examinationNewPassword, setExaminatioNewPassword] = useState("");
    const [subject, setSubject] = useState("");
    const [color, setColor] = useState("red");
    const [visibility, setVisibility] = useState(true);

    const userInfo = GetSecureStore('user')
    let loginUser = ''
    if (userInfo !== false) {
        const userInfoParse = JSON.parse(userInfo)
        loginUser = userInfoParse[0].login
    }

    const params = {
        new_password: newPassword,
        password: oldPassword,
        email_or_login: loginUser,
    }

    let token = GetSecureStore('token')
    const config = {
        headers: { Authorization: `Bearer ${token}` },
    };

    if (status === null) {
        requestPermission();
    }

    const onSavePassword = async () => {
        if (newPassword !== examinationNewPassword) {
            setVisibility(true)
            setColor("red")
            setSubject("Пароли не совпадают")
        } else {
            axios.put(HOST_SERVER + '/updatepassword', params, config)
                .then(function (response) {
                    setColor("green")
                    setSubject("Пароль изменен")
                    setOldPassword("")
                    setNewPassword("")
                    setExaminatioNewPassword("")

                    setVisibility(true)
                    setTimeout(() => {
                        setVisibility(false)
                    }, 5000)
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
    }

    return (
        <View style={styles.menuContainer}>
            <View style={styles.menuContainer}>
                <GestureHandlerRootView style={styles.container}>
                    { visibility &&
                        <Result subject={subject} color={color} top={0} bottom={8}/>
                    }
                    <View style={styles.viewText}>
                        <Text>
                            Старый пароль
                        </Text>
                    </View>
                    <Input
                        placeholder="Старый пароль"
                        onChangeText={ (oldPassword) => setOldPassword(oldPassword)}
                        value={oldPassword}
                        secureTextEntry={true}
                        containerStyle={styles.inputStyle}
                    />
                    <View style={styles.viewText}>
                        <Text>
                            Новый пароль
                        </Text>
                    </View>
                    <Input
                        placeholder="Новый пароль"
                        onChangeText={ (newPassword) => setNewPassword(newPassword)}
                        value={newPassword}
                        secureTextEntry={true}
                        containerStyle={styles.inputStyle}
                    />
                    <View style={styles.viewText}>
                        <Text>
                            Повторите новый пароль
                        </Text>
                    </View>
                    <Input
                        placeholder="Повторите новый пароль"
                        onChangeText={ (examinationNewPassword) => setExaminatioNewPassword(examinationNewPassword)}
                        value={examinationNewPassword}
                        secureTextEntry={true}
                        containerStyle={styles.inputStyle}
                    />
                    <View style={styles.footerContainer}>
                        <Button theme="save" label="Сохранить" onPress={onSavePassword} />
                    </View>
                    <StatusBar style="light" />
                </GestureHandlerRootView>
            </View>
        </View>
    );
}


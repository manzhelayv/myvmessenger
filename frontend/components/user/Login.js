import * as React from 'react';
import { Text, View } from 'react-native';
import { Input } from 'react-native-elements';
import axios from 'axios';
import Button from '../form/Button';
import { useState, useContext } from 'react';
import {SaveSecureStore} from '../storage/SecureStore';
import {styles} from "./css/Login";

const current = new Date();
const nextYear = new Date();

nextYear.setFullYear(current.getFullYear() + 1);

function LoginResult({subject, color}){
    return (
      <View style={styles.viewLoginResult}>
        <Text style={{color: color}}>
        {subject}
        </Text>
      </View>
    )
}
  
export default function Login({setKey, UserContext, setRestartWebsocket}) {
    const {HOST_SERVER} = useContext(UserContext);
    const [emailOrLogin, setEmailOrLogin] = useState("");
    const [password, setPassword] = useState("");
    const [subject, setSubject] = useState("Введите логин или email и пароль");
    const [color, setColor] = useState("green");
  
    const params = {
      email_or_login: emailOrLogin, 
      password: password,
    }

    const config = {
      headers: { Accept: 'application/json' }
    };

    const loginUser = async () => {
      axios.post(HOST_SERVER + '/login', params, config)
      .then(function (response) {
        let tokenUser = response.data.access_token
        let loginUser = response.data.login
        let nameUser = response.data.name
        let emailUser = response.data.email
        let tdidUser = response.data.tdid
        let phoneUser = response.data.phone

        const user = [
          {
            'token': tokenUser,
            'login': loginUser,
            'name': nameUser,
            'email': emailUser,
            'phone': phoneUser,
            'tdid': tdidUser,
          },
        ];

        SaveSecureStore('user', JSON.stringify(user))
        SaveSecureStore('token', tokenUser)
        setKey(2)
        setEmailOrLogin("")
        setPassword("")
        setRestartWebsocket("")
      })
      .catch(function (error) {
        setColor("red")
        const err = JSON.parse(error.request.response);
        setSubject(err.error.message)
      });
    };
  
    return (
      <View style={styles.containerLogin}>
        <LoginResult subject={subject} color={color} />
        <View style={styles.container}>
            <View style={styles.viewText}>
                <Text>
                    Email или логин
                </Text>
            </View>
            <Input 
              placeholder="Email или логин"
              onChangeText={ (emailOrLoginText) => setEmailOrLogin(emailOrLoginText)}
              value={emailOrLogin}
              containerStyle={styles.inputStyle}
            />
            <View style={styles.viewText}>
                <Text style={styles.textPasswordStyle}>
                    Пароль
                </Text>
            </View>
            <Input 
              placeholder="Пароль"
              onChangeText={ (passwordText) => setPassword(passwordText)}
              value={password}
              containerStyle={styles.inputStyle}
              secureTextEntry={true}
            />
          <View style={styles.footerContainer}>
            <Button theme="primary" label="Отправить" onPress={loginUser}  />
          </View>
        </View>
      </View>
    );
  }


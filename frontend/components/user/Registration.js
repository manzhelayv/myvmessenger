import * as React from 'react';
import { Text, View } from 'react-native';
import { Input } from 'react-native-elements';
import axios from 'axios';
import Button from '../form/Button';
import { useState, useContext } from 'react';
import {SaveSecureStore} from '../storage/SecureStore';
import {styles} from "./css/Registration";
import MaskInput from 'react-native-mask-input';

function LoginResult(subject){
    return (
      <View style={styles.viewLoginResult}>
        <Text style={{color: subject.color,  fontSize: 16}}>
        {subject.subject}
        </Text>
      </View>
    )
}

export default function Registration({UserContext, setRestartWebsocket}) {
    const {cookies, setUserStatus, HOST_SERVER} = useContext(UserContext);
    const [email, setEmail] = useState("");
    const [name, setName] = useState("");
    const [login, setLogin] = useState("");
    const [phone, setPhone] = useState("");
    const [password, setPassword] = useState("");
    const [subject, setSubject] = useState("Заполните все поля");
    const [color, setColor] = useState("green");
  
    const params = {
      email: email,
      name: name,
      login: login,
      password: password,
      phone: phone,
    }
  
    const userRegistration = async () => {
      axios.post(HOST_SERVER + '/user', params)
      .then(function (response) {
        let tokenUser = response.data.access_token
        let loginUser = response.data.login
        let nameUser = response.data.name
        let emailUser = response.data.email
        let tdidUser = response.data.tdid
        let phoneUser = response.data.phone
  
        setEmail("")
        setName("")
        setLogin("")
        setPassword("")
        setPhone("")

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

        setUserStatus(tokenUser)
        setRestartWebsocket("")
      })
      .catch(function (error) {
        setColor("red")
        const err = JSON.parse(error.request.response);
        setSubject(err.error.message)
      });
    };
  
    return (
      <View style={styles.containerRegistration}>
          <LoginResult subject={subject} color={color} />
          <View style={styles.container}>
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
            <Button theme="primary" label="Отправить" onPress={userRegistration}  />
          </View>
        </View>
      </View>
    );
  }


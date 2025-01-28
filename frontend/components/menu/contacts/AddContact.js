import * as React from 'react';
import { Text, View, TextInput } from 'react-native';
import axios from 'axios';
import { useContext, useState } from 'react';
import {keyAddContactsContext} from './Contacts';
import Button from '../../form/Button';
import { GetSecureStore } from '../../storage/SecureStore.js';
import { SetAsyncStorage } from '../../storage/AsyncStorage.js';
import {styles} from "./css/AddContact";

let keyThis = 1

function LoginResult(subject){
    return (
      <Text style={{color: subject.color, fontSize: 14}}>
       {subject.subject}
      </Text>
    )
}

export default function AddContact({HOST_SERVER, setContacts}) {
    const {key, setKey} = useContext(keyAddContactsContext);
    const [contact, setContact] = useState([]);
    const [color, setColor] = useState("green");
    const [subject, setSubject] = useState("");

    let token = GetSecureStore('token')
    const config = {
      headers: { Authorization: `Bearer ${token}` }
    };

    const params = {
        email_or_login: contact,
    }

    async function loadData() {
        axios.post(HOST_SERVER + '/contacts', params, config)
        .then(function (response) {
            setContact("")

            keyThis++
            setKey(keyThis)
            SetAsyncStorage('mess', true)

            setSubject("")

            let userContacts = response.data
            setContacts(userContacts)
        })
        .catch(function (error) {
            setColor("red")
            const err = JSON.parse(error.request.response);
            setSubject(err.error.message)
        });
    }

    let topButton = 23
    let topTextInput = 1
    const [topView, setTopView] = useState(0);

    const handleOnFocus = () => {
        setTopView(-10)
    };

    const handleOnBlur = () => {
        setTopView(0)
    };

    return (
        <View style={[styles.addContacts, {top: topView}]}>
          <View style={styles.viewLoginResult}>
            <LoginResult subject={subject} color={color} />
          </View>
          <View style={[styles.addContact, {top: topTextInput}]} key={1}>
            <TextInput
                style={styles.textArea}
                underlineColorAndroid="transparent"
                placeholderTextColor="grey"
                numberOfLines={10}
                multiline={true}
                onChangeText={ (contact) => setContact(contact)}
                value={contact}
                onFocus={() => handleOnFocus()}
                onBlur={() => handleOnBlur()}
                />
          </View>
          <View style={[styles.addButton, {top: topButton}]} key={2}>
              <Button theme="primary" label="Добавить" onPress={loadData}   />
          </View>
        </View>
    );
}


  
import * as React from 'react';
import {Text, View, Pressable} from 'react-native';
import axios from 'axios';
import { useState, createContext } from 'react';
import AddContact from './AddContact'
import {ScrollViewIndicator} from '@fanchenbao/react-native-scroll-indicator';
import ImageViewer from '../../form/ImageViewer';
import * as FileSystem from "expo-file-system";
import { GetSecureStore } from '../../storage/SecureStore.js';
import {SetAsyncStorage, GetAsyncStorage, SetAsyncStorageObject} from '../../storage/AsyncStorage.js';
import {cookies, HOST_SERVER, UserContext} from '../../../App';
import {globalOnPress} from "../TopMenu";
import * as Contact from "expo-contacts";
import Button from "../../form/Button";
import {styles} from "./css/Contacts";

let userContacts = []
let keyThis = 1

export const keyAddContactsContext = createContext();

const phoneFormat = (s = true) => {
    if (s.startsWith('8')) {
        let phone = s.substr(1);

        return '+7 ' + phone
    }

    return s.replace(/(\d{1})(\d{3})(\d{3})(\d{4})/g, `$1 $2 $3 $4`)
}

export default function Contacts() {
    let onPress = globalOnPress

    const [contacts, setContacts] = useState([]);
    const [key, setKey] = useState(1);
    const [error, setError] = useState("");
    const [color, setColor] = useState("");
    const [top, setTop] = useState(0);
    const [visibility, setVisibility] = useState(true);

    let token = GetSecureStore('token')
    const config = {
      headers: { Authorization: `Bearer ${token}` }
    };

    let mess = GetAsyncStorage('mess')
    mess.then((mess) => {
        if (mess === true) {
            userContacts = []
            SetAsyncStorage('mess', false)
        }
    })

    async function showFirstContactAsync() {
        const { status } = await Contact.requestPermissionsAsync();
        if (status === 'granted') {
            const { data } = await Contact.getContactsAsync({
                fields: [Contact.Fields.Emails, Contact.Fields.PhoneNumbers],
            });

            let contacts = []

            if (data.length > 0) {
                {data.map((contact, index) => {
                    let phone = phoneFormat(contact.phoneNumbers[0].number)
                    let objContact = {phone: phone, name: contact.name}
                    contacts.push(objContact)
                })}
            }

            let token = GetSecureStore('token')
            const config = {
                headers: { Authorization: `Bearer ${token}` }
            };

            const params = {
                contacts: contacts,
            }

            await axios.put(HOST_SERVER + '/contacts', params, config)
                .then(function (response) {
                    if (response.data !== null) {
                        userContacts = response.data

                        setContacts(userContacts)
                        keyThis++

                        setError("Контакты обновлены")
                        setColor('green')
                        setTop(15)

                        setVisibility(true)
                        setTimeout(() => {
                            setVisibility(false)
                        }, 5000);
                    }
                })
                .catch(function (error) {
                    const err = JSON.parse(error.request.response);
                    setError(err.error.message)
                    setColor('red')
                    setTop(15)

                    setVisibility(true)
                    setTimeout(() => {
                        setVisibility(false)
                    }, 5000);
                });
        }
    }

    async function loadData() {
        axios.get(HOST_SERVER + '/contacts', config)
        .then(function (response) {
            if (response.data !== null) {
              userContacts = response.data

              setContacts(userContacts)
              keyThis++
            }
        })
        .catch(function (error) {
            console.log(error.request.response);
        });
    } 

    if (userContacts.length === 0) {
        loadData();
    } else if (contacts.length === 0) {
        setContacts(userContacts)
    }

    return (
      <View style={styles.container}>
        { visibility &&
          <View style={styles.updateText}>
              <Text style={{color: color}}>{error}</Text>
          </View>
        }
        <View style={[styles.containerContacts, {top: top}]}>
            <View style={styles.addButton}>
                <Button theme="primary" label="Обновить" onPress={showFirstContactAsync} />
            </View>
          <ScrollViewIndicator indStyle={styles.scrollView}>
              {contacts.map((value, index) => {
                  FileSystem.writeAsStringAsync(FileSystem.documentDirectory + value.user_to + '.jpeg', value.avatar, { encoding: FileSystem.EncodingType.Base64 })
                      .then(async (uri) => {

                      })
                      .catch(e => console.log('err', value.user_to, e));

                  FileSystem.readAsStringAsync(FileSystem.documentDirectory + value.user_to + '.jpeg', { encoding: FileSystem.EncodingType.Base64 })
                      .then(async (b64avatar) => {

                      }).catch(e => console.log('err', e));

                  if (cookies.get('activeTabChat') === 'contacts'){
                      const objUserTo = {email: value.email, name: value.user_to_name, userTo: value.user_to, userFrom: '', message: '', date: '', avatar: FileSystem.documentDirectory + value.user_to + '.jpeg', avatarFile: FileSystem.documentDirectory + value.user_to + '.jpeg', countNotReadMessage: value.countNotReadMessage};

                      SetAsyncStorageObject('userTo' + value.user_to, objUserTo)
                  }

                  let img = 'data:image/jpeg;base64,' + value.avatar
                  let image = ""

                  if (img !== ''){
                   image = <ImageViewer placeholderImageSource="" selectedImage={img} size="25"/>
                  } 

                  return (
                    <View key={index} style={[styles.textView]} >
                      <Pressable
                        onPress={()=> onPress(value.email, value.user_to_name, value.user_to, "", "", "", "", "", value.avatar)}
                        key={index}
                      >
                        {image}
                        <View style={styles.rightTextView}>
                          <Text style={styles.rightText}>{value.user_to_name}</Text> 
                        </View>
                      </Pressable>
                    </View>
                  )
              })}
          </ScrollViewIndicator>
      </View>
      <keyAddContactsContext.Provider value={{key, setKey}} style={styles.addContacts}> 
        <AddContact HOST_SERVER={HOST_SERVER} setContacts={setContacts} />
      </keyAddContactsContext.Provider>
      </View>
    );
  }



  
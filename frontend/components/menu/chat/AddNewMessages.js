import React from 'react';
import { Text, View} from 'react-native';
import ImageViewerSettings from "../../form/ImageViewerSettings";
import * as FileSystem from "expo-file-system";
import {getFileExtensionData, GetStringMonth} from "./Function";
import {GetSecureStore} from "../../storage/SecureStore";
import {styles} from "./css/MessagesCss";

let blockMessage = ""

export default function AddNewMessages({message, styles}) {
    let newMesss = []
    for (const value of message) {
        let d = new Date()
        let minutes = (d.getMinutes() < 10 ? '0' : '') + d.getMinutes()
        let hours = d.getHours() + ":" + minutes;

        let dateTime = hours

        let date = d.getDate() + " " + GetStringMonth(d.getMonth()) + " " + d.getFullYear() + " г."

        if (value['right'] !== undefined) {
            if (value['right'].date === undefined) {
                value['right'].date = date
            }

            if (value['right'].dateTime === undefined) {
                value['right'].dateTime = dateTime
            }

            if (value['right'].file !== '') {
                FileSystem.readAsStringAsync(FileSystem.documentDirectory + value['right'].file, {encoding: FileSystem.EncodingType.Base64})
                    .then((base64) => {
                        value['right'].message = base64
                    })
            }
        }

        if (value['left'] !== undefined) {
            if (value['left'].date === undefined) {
                value['left'].date = date
            }

            if (value['left'].dateTime === undefined) {
                value['left'].dateTime = dateTime
            }

            if (value['left'].file !== '') {
                FileSystem.readAsStringAsync(FileSystem.documentDirectory + value['left'].file, {encoding: FileSystem.EncodingType.Base64})
                    .then((base64) => {
                        value['left'].message = base64
                    })
            }
        }

        newMesss.push(value)
    }

    const userTo = GetSecureStore('userTo')

    return (
            newMesss.map((value, index) => {
                let paddingTop = 0
                if (index === 0) {
                    paddingTop = 10
                }

                if(value['right'] !== undefined) {
                    if (value['right'].file !== '') {
                        let fileData = getFileExtensionData(value['right'].file)

                        blockMessage = <>
                                        <Text style={styles.buttonTextDate}>{value['right'].date}</Text>
                                        <View style={styles.leftImage}>
                                            <ImageViewerSettings placeholderImageSource="" selectedImage={fileData + value['right'].message} />
                                        </View>
                                        <View style={styles.rightTextDate}>
                                            <Text style={styles.buttonTextDate}>{value['right'].dateTime}</Text>
                                        </View>
                                       </>
                    } else {
                        blockMessage = <>
                                        <Text style={styles.buttonTextDate}>{value['right'].date}</Text>
                                        <View style={styles.leftText}>
                                            <Text style={styles.buttonText}>{value['right'].message}</Text>
                                        </View>
                                        <View style={styles.rightTextDate}>
                                            <Text style={styles.buttonTextDate}>{value['right'].dateTime}</Text>
                                        </View>
                                       </>
                    }

                    return (
                        <View style={{display: 'block', paddingTop: paddingTop}} key={index}>
                            <View style={styles.rightChatText}>
                                {blockMessage}
                            </View>
                        </View>
                    )
                } else if (value['left'].userTo !== undefined){
                    if (value['left'].file !== '') {
                        let fileData = getFileExtensionData(value['left'].file)
                        blockMessage = <>
                            <Text style={styles.buttonTextDate}>{value['left'].date}</Text>
                            <View style={styles.leftImage}>
                                <ImageViewerSettings placeholderImageSource="" selectedImage={fileData + value['left'].message} />
                            </View>
                            <View style={styles.rightTextDate}>
                                <Text style={styles.buttonTextDate}>{value['left'].dateTime}</Text>
                            </View>
                        </>
                    } else {
                        blockMessage = <>
                            <Text style={styles.buttonTextDate}>{value['left'].date}</Text>
                            <View style={styles.leftText}>
                                <Text style={styles.buttonText}>{value['left'].message}</Text>
                            </View>
                            <View style={styles.rightTextDate}>
                                <Text style={styles.buttonTextDate}>{value['left'].dateTime}</Text>
                            </View>
                        </>
                    }

                    return (
                        <View style={{display: 'block', paddingTop: paddingTop}} key={index}>
                            <View style={styles.leftChatText}>
                                {blockMessage}
                            </View>
                        </View>
                    )
                }
            })
    );
}


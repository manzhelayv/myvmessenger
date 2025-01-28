import { StatusBar } from 'expo-status-bar';
import {View, Platform, Text} from 'react-native';
import { GestureHandlerRootView } from "react-native-gesture-handler";
import * as FileSystem from 'expo-file-system';
import Button from '../../form/Button';
import ImageViewerSettings from '../../form/ImageViewerSettings';
import * as ImagePicker from 'expo-image-picker';
import { useState, useRef, useEffect } from 'react';
import * as MediaLibrary from 'expo-media-library';
import domtoimage from 'dom-to-image';
import axios from 'axios';
import {GetSecureStore} from "../../storage/SecureStore";
import { HOST_SERVER} from '../../../App';
import * as React from "react";
import {styles} from "./css/Avatar";
import {Result} from './Settings';

const PlaceholderImage = require('../../../assets/images/background-image.png');

let avatar = ""

export default function Avatar() {
    const [showAppOptions, setShowAppOptions] = useState(false);
    const [selectedImage, setSelectedImage] = useState(null);
    const [status, requestPermission] = MediaLibrary.usePermissions();
    const [statusImage, setStatusImage] = useState(false);
    const imageRef = useRef();
    const [subject, setSubject] = useState("");
    const [color, setColor] = useState("green");
    const [visibility, setVisibility] = useState(true);

    const pickImageAsync = async () => {
        let result = await ImagePicker.launchImageLibraryAsync({
            allowsEditing: true,
            quality: 1,
        });

        const base64 = await FileSystem.readAsStringAsync(result.assets[0].uri, { encoding: 'base64' });
        if (!result.canceled) {
            setSelectedImage('data:image/jpeg;base64,' + base64)
            setShowAppOptions(true);
        } else {
            alert('You did not select any image.');
        }
    };

    const onReset = () => {
        setShowAppOptions(false);
    };

    let token = GetSecureStore('token')
    const config = {
        headers: { Authorization: `Bearer ${token}` },
    };

    useEffect(() => {
        if (statusImage === false) {
            axios.get(HOST_SERVER + '/profile', config)
                .then(function (response) {
                    avatar = 'data:image/jpeg;base64,' + response.data.avatar

                    setSelectedImage(avatar)
                    setStatusImage(true)
                })
                .catch(function (error) {
                    const err = JSON.parse(error.request.response);
                    console.log(err.error.message)
                });
        }
    });

    const onSaveImageAsync = async () => {
        if (Platform.OS !== 'web') {
            try {
                const params = {
                    image: selectedImage,
                }

                axios.post(HOST_SERVER + '/profile', params, config)
                    .then(function (response) {
                        setColor("green")
                        setSubject("Аватар успешно обновлен")

                        avatar = selectedImage

                        setVisibility(true)
                        setTimeout(() => {
                            setVisibility(false)
                        }, 5000);

                    }).catch(async function (error) {
                        setColor("red")
                        const err = JSON.parse(error.request.response);
                        setSubject(err.error.message)

                        setSelectedImage(avatar)

                        setVisibility(true)
                        setTimeout(() => {
                            setVisibility(false)
                        }, 5000);
                });
                setShowAppOptions(false);
            } catch (e) {
                console.log(e);
            }
        } else {
            try {
                const dataUrl = await domtoimage.toJpeg(imageRef.current, {
                    quality: 1, // 0.95
                    width: 250,
                    height: 250,
                });

                const params = {
                    image: dataUrl,
                }

                axios.post(HOST_SERVER + '/profile', params, config)
                    .then(function (response) {

                    }).catch(function (error) {
                        const err = JSON.parse(error.request.response);
                        console.log(err.error.message)
                });
                setShowAppOptions(false);
            } catch (e) {
                console.log(e);
            }
        }
    };

    if (status === null) {
        requestPermission();
    }

    return (
        <View style={styles.menuContainer}>
            <View style={styles.menuContainer}>
                <GestureHandlerRootView style={styles.container}>
                    { visibility &&
                        <Result subject={subject} color={color} top={5} bottom={15} />
                    }
                    <View style={styles.imageContainer}>
                        <View ref={imageRef} collapsable={false}>
                            <ImageViewerSettings placeholderImageSource={PlaceholderImage} selectedImage={selectedImage} />
                        </View>
                    </View>
                    {showAppOptions ? (
                        <View>
                            <View style={styles.optionsRow}>
                                <Button theme="reset" label="Сбросить" onPress={onReset} />
                                <Button theme="save" label="Сохранить" onPress={onSaveImageAsync} />
                            </View>
                        </View>
                    ) : (
                        <View style={styles.footerContainer}>
                            <Button theme="primary" label="Выберите" onPress={pickImageAsync} />
                        </View>
                    )}
                    <StatusBar style="light" />
                </GestureHandlerRootView>
            </View>
        </View>
    );
}

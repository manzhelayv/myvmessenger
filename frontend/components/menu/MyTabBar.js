import {useContext, useRef} from "react";
import {SetAsyncStorage} from "../storage/AsyncStorage";
import {Image, Pressable, TouchableOpacity, View} from "react-native";
import MaterialCommunityIcons from "@expo/vector-icons/MaterialCommunityIcons";
import FontAwesome from "@expo/vector-icons/FontAwesome";
import ImageViewer from "../form/ImageViewer";
import PlaceholderImage from "../../assets/images/background-image.png";
import goBack from "../../assets/images/go_back.png";
import * as React from "react";
import Entypo from "@expo/vector-icons/Entypo";
import MaterialIcons from '@expo/vector-icons/MaterialIcons';
import Animated from 'react-native-reanimated';
import {styles} from "./css/TopMenu";

let tabStylesWidth = []
let tabBarIcon = ""

export default function MyTabBar({ state, descriptors, navigation, position, UserContext, userStatuses }) {
    const {cookies} = useContext(UserContext);
    const imageRef = useRef();

    SetAsyncStorage('activeTabChat', '')

    {state.routes.map((route, index) => {
        const isFocused = state.index === index;

        if (isFocused) {
            SetAsyncStorage('activeTabChat', route.name)
            cookies.set('activeTabChat', route.name)
        }
    })}

    let imgFromState = ""
    if (state.routes[1] !== undefined) {
        let stateRoutes = state.routes[1].params

        if (stateRoutes !== undefined) {
            imgFromState = 'data:image/jpeg;base64,' + stateRoutes.img
        }
    }

    let activeTabChat = cookies.get('activeTabChat')

    return (
        <View style={{ flexDirection: 'row', paddingTop: 20 }}>
            {state.routes.map((route, index) => {
                const {options} = descriptors[route.key];
                const label =
                    options.tabBarLabel !== undefined
                        ? options.tabBarLabel
                        : options.title !== undefined
                            ? options.title
                            : route.name;

                const isFocused = state.index === index;

                const onPress = () => {
                    const event = navigation.emit({
                        type: 'tabPress',
                        target: route.key,
                    });

                    if (!isFocused && !event.defaultPrevented) {
                        navigation.navigate(route.name);
                    }
                };

                if (route.name !== 'login' && route.name !== 'registration') {
                    tabStylesWidth[route.name] = 0
                    if (isFocused) {
                        tabStylesWidth[route.name] = '100%'
                    }
                } else {
                    tabStylesWidth[route.name] = "50%"
                }

                const onLongPress = () => {
                    navigation.emit({
                        type: 'tabLongPress',
                        target: route.key,
                    });
                };

                if (activeTabChat !== 'chat') {
                    if (route.name === 'registration') {
                        tabBarIcon = <MaterialIcons
                            name="app-registration"
                            size={18}
                            color="#25292e"
                            style={styles.buttonIconRegistration}
                        />
                    }

                    if (route.name === 'login') {
                        tabBarIcon = <MaterialCommunityIcons
                            name="login"
                            size={18}
                            color="#25292e"
                            style={styles.buttonIconLogin}
                        />
                    }

                    if (route.name === 'Contacts') {
                        tabBarIcon = <MaterialCommunityIcons
                            name="contacts"
                            size={18}
                            color="#25292e"
                            style={styles.buttonIconContacts}
                            accessibilityViewIsModal='true'
                        />
                    }

                    if (route.name === 'chats') {
                        tabBarIcon = <Entypo
                            name="chat"
                            size={18}
                            color="#25292e"
                            style={styles.buttonIconChat}
                        />
                    }

                    if (route.name === 'Settings') {
                        tabBarIcon = <FontAwesome
                            name="user"
                            size={18}
                            color="#25292e"
                            style={styles.buttonIconSettings}
                        />
                    }
                }

                let justifyContent = "flex-start"
                let alignItems = "start"

                if (route.name !== 'chat') {
                    alignItems = 'center'
                    justifyContent = 'center'
                }

                if (route.name === 'login' || route.name === 'registration') {
                    return <TouchableOpacity
                        accessibilityRole="button"
                        accessibilityState={isFocused ? { selected: true } : {}}
                        accessibilityLabel={options.tabBarAccessibilityLabel}
                        testID={options.tabBarTestID}
                        onPress={onPress}
                        onLongPress={onLongPress}
                        tabBarIndicatorStyle= {{backgroundColor: "#03fc5e", top:509}}
                        style={{
                            flex: 1,
                            backgroundColor: isFocused && (route.name === 'login' || route.name === 'registration')? '#03fc5e' : "white",// visibility: isFocused ? 'hidden' : 'visible',
                            maxWidth: tabStylesWidth[route.name],
                            height: 30,
                            margin:5,
                            marginLeft:2,
                            marginRight:1,
                            borderRadius:100,
                            top: route.name === 'chat' ? 10 : -10,
                            justifyContent: justifyContent,
                            alignItems: alignItems,
                            position: "relative",
                        }}
                    >
                        {tabBarIcon}

                        {route.name === 'login' ? (
                                <Animated.Text style={styles.tabBarLabelStyleL}>{label}</Animated.Text>
                            )  :
                            <></>
                        }

                        {route.name === 'registration' ? (
                                <Animated.Text style={styles.tabBarLabelStyleR}>{label}</Animated.Text>
                            )  :
                            <></>
                        }
                    </TouchableOpacity>
                }

                let styleChildSettingsRight = 4

                if (route.name === 'profile') {
                    styleChildSettingsRight = 7
                }

                if (route.name === 'password') {
                    styleChildSettingsRight = 10
                }

                return (
                    <View
                        style={{
                            flex: 1,
                            backgroundColor: isFocused && (route.name === 'login' || route.name === 'registration')? '#03fc5e' : "white",// visibility: isFocused ? 'hidden' : 'visible',
                            maxWidth: tabStylesWidth[route.name],
                            height: 30,
                            margin:5,
                            marginLeft:2,
                            marginRight:1,
                            borderRadius:100,
                            top: route.name === 'chat' ? 10 : -10,
                            justifyContent: justifyContent,
                            alignItems: alignItems,
                            position: "relative",
                        }}
                    >
                        {(route.name === 'chats' && isFocused) ? (
                                <View style={[styles.containerTopChats]}>
                                    <Animated.Text style={styles.tabBarLabelStyleChats}>MYVMessenger</Animated.Text>
                                </View>
                            ) :
                            <></>
                        }
                        {(route.name === 'chat' && isFocused) ? (
                                <View style={[styles.containerTopChat]}>
                                    <View style={styles.GoBack}>
                                        <Pressable
                                            onPress={() => {
                                                cookies.set('active_menu', 'Chats');
                                                navigation.navigate('chats');
                                            }}
                                        >
                                            <Image
                                                source={goBack}
                                                style={{
                                                    width: 20,
                                                    height: 20,
                                                }}
                                            />
                                        </Pressable>
                                    </View>
                                    { (imgFromState !== '') ? (
                                        <View ref={imageRef} collapsable={false} style={styles.ChatView}>
                                            <ImageViewer placeholderImageSource={PlaceholderImage} selectedImage={imgFromState} chat="" />
                                        </View>
                                    ) :
                                        <></>
                                    }
                                    <Animated.Text style={styles.tabBarLabelStyleChat}>{cookies.get('nameActiveContact')}</Animated.Text>
                                    { (userStatuses !== undefined && userStatuses !== '') ? (
                                        <Animated.Text style={styles.tabBarLabelStyleChatStatus}>{userStatuses}</Animated.Text>
                                    ) :
                                        <Animated.Text style={styles.tabBarLabelStyleChatStatus}>Не известен</Animated.Text>
                                    }
                                </View>
                            ) :
                            <></>
                        }

                        { (route.name === 'settings' && isFocused) ? (
                                <View style={[styles.containerTopSettingsAndContacts]}>
                                    <View style={styles.GoBackSettingsAndContacts}>
                                        <Pressable
                                            onPress={() => {
                                                cookies.set('active_menu', 'Chats');
                                                navigation.navigate('chats');
                                            }}
                                        >
                                            <Image
                                                source={goBack}
                                                style={{
                                                    width: 20,
                                                    height: 20,
                                                }}
                                            />
                                        </Pressable>
                                    </View>
                                    <Animated.Text style={styles.tabBarLabelStyle}>{label}</Animated.Text>
                                </View>
                            ) :
                            <></>
                        }

                        { (route.name === 'avatar' && isFocused) || (route.name === 'profile' && isFocused) || (route.name === 'password' && isFocused) ? (
                                <View style={[styles.containerTopSettingsAndContacts, {right: styleChildSettingsRight}]}>
                                    <View style={styles.GoBackSettingsAndContacts}>
                                        <Pressable
                                            onPress={() => {
                                                cookies.set('active_menu', 'Settings');
                                                navigation.navigate('settings');
                                            }}
                                        >
                                            <Image
                                                source={goBack}
                                                style={{
                                                    width: 20,
                                                    height: 20,
                                                }}
                                            />
                                        </Pressable>
                                    </View>
                                    <Animated.Text style={styles.tabBarLabelStyle}>{label}</Animated.Text>
                                </View>
                            ) :
                            <></>
                        }

                        { (route.name === 'contacts' && isFocused) ? (
                                <View style={[styles.containerTopSettingsAndContacts]}>
                                    <View style={styles.GoBackSettingsAndContacts}>
                                        <Pressable
                                            onPress={() => {
                                                cookies.set('active_menu', 'Chats');
                                                navigation.navigate('chats');
                                            }}
                                        >
                                            <Image
                                                source={goBack}
                                                style={{
                                                    width: 20,
                                                    height: 20,
                                                }}
                                            />
                                        </Pressable>
                                    </View>
                                    <Animated.Text style={styles.tabBarLabelStyle}>{label}</Animated.Text>
                                </View>
                            ) :
                            <></>
                        }
                    </View>
                );
            })}
        </View>
    );
}



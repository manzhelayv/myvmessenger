import React from "react";
import { View, Text, Pressable} from "react-native";
import {useNavigation} from "@react-navigation/native";
import {styles} from "./css/Settings";
import MaterialIcons from "@expo/vector-icons/MaterialIcons";
import {AntDesign} from "@expo/vector-icons";

export function Result({subject, color, top, bottom}){
    return (
        <View style={{top: top, paddingBottom: bottom}}>
            <Text style={{color: color}}>
                {subject}
            </Text>
        </View>
    )
}

export default function Settings() {
    const navigation = useNavigation();

    return (
        <View style={styles.menuContainer}>
            <View style={[styles.profileLink]} >
                <Pressable
                    style={[styles.button, { backgroundColor: "#fff" }]}
                    onPress={() => {
                        navigation.navigate('avatar');
                    }}
                >
                    <MaterialIcons
                        name="person"
                        size={25}
                        color="#25292e"
                        style={styles.buttonIcon}
                    />
                    <View style={styles.viewBlock}>
                        <Text style={styles.text}>
                            Аватар
                        </Text>
                        <Text style={styles.textDescription}>
                            Создание, изменение аватара пользователя
                        </Text>
                    </View>
                </Pressable>
            </View>
            <View style={styles.profileLink}>
                <Pressable
                    style={[styles.button, { backgroundColor: "#fff" }]}
                    onPress={() => {
                        navigation.navigate('profile');
                    }}
                >
                    <AntDesign
                        name="profile"
                        size={25}
                        color="#25292e"
                        style={styles.buttonIcon}
                    />
                    <View style={styles.viewBlock}>
                        <Text style={styles.text}>
                            Профайл
                        </Text>
                        <Text style={styles.textDescription}>
                            Изменение данных профайла пользователя
                        </Text>
                    </View>
                </Pressable>
            </View>
            <View style={styles.profileLink}>
                <Pressable
                    style={[styles.button, { backgroundColor: "#fff" }]}
                    onPress={() => {
                        navigation.navigate('password');
                    }}
                >
                    <MaterialIcons
                        name="password"
                        size={25}
                        color="#25292e"
                        style={styles.buttonIcon}
                    />
                    <View style={styles.viewBlock}>
                        <Text style={styles.text}>
                            Пароль
                        </Text>
                        <Text style={styles.textDescription}>
                            Изменение пароля пользователя
                        </Text>
                    </View>
                </Pressable>
            </View>
        </View>
    )
}


import React, { memo } from 'react';
import { View } from 'react-native';
import ItemOptions from './ItemOptions';
import {useNavigation} from "@react-navigation/native";

const Item = ({ id, idMenu, setIdMenu }) => {
    const navigation = useNavigation();

    return (
        <>
            <View id={id}>
                <ItemOptions
                    id={id}
                    idMenu={idMenu}
                    setIdMenu={setIdMenu}
                    options={[
                        {
                            title: 'Настройки',
                            onPress: () => {
                                navigation.navigate('settings');
                            },
                        },
                        {
                            title: 'Контакты',
                            onPress: () => {
                                navigation.navigate('contacts');
                            },
                        },
                    ]}
                />
            </View>
        </>
    );
};

export default memo(Item);

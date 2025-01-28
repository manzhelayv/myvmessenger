import React, { memo } from 'react';
import {Image, Platform, StyleSheet, Text, View} from 'react-native';
import { Menu, MenuOptions, MenuOption, MenuTrigger } from 'react-native-popup-menu';

const PlaceholderImage = require('../../images/menu.png');

const ItemOptions = ({ id, idMenu, options, setIdMenu }) => {
    const handleOnPress = (id) => {
        if (id === idMenu) {
            setIdMenu('');
            return;
        }
        setIdMenu(id);
    };

    return (
        <Menu onClose={() => setIdMenu('')} style={styles.container}>
            <MenuTrigger
                onPress={() => handleOnPress(id)}
                customStyles={{
                    triggerWrapper: {
                        top: 25,
                        left: '25%',
                        maxWidth: '100%',
                        width: 1000,
                        height: 500,
                        borderRadius: 30,
                        float: 'right',
                        marginLeft: "auto",
                        fontSize: 20,

                    },
                    triggerTouchable: {
                        underlayColor: 'green',
                        maxWidth: '100%',
                        width: 1000,
                        borderRadius: 30,
                        fontSize: 20,
                        left: 30
                    },
                }}
            >
                <Image
                    source={PlaceholderImage}
                    style={{
                        width: 50,
                        height: 50,
                        zIndex:9999999,
                        opacity: 111,
                        float: 'right',
                        marginLeft: "auto",
                        display:'inline-block',
                        top: 5,
                        left:5,
                        position: 'absolute',
                    }}Z
                />
            </MenuTrigger>

            <MenuOptions optionsContainerStyle={styles.OptionsStyle}
                         customStyles={{
                             optionsWrapper: {
                                 bottom: Platform.OS === 'ios' ? -60 : -72,
                                 height: 100,
                                 backgroundColor: 'white',
                                 borderRadius: 8,
                                 padding: 8,
                                 width: "70%",
                                 shadowOffset: { width: 0, height: 3 },
                                 shadowOpacity: 0.2,
                                 shadowRadius: 4,
                                 elevation: 1,
                                 zIndex: 9999,
                                 float: 'right',
                                 marginLeft: "auto",
                                 top: 20,
                                 left:-202,
                                 fontSize: 25,
                             },
                             optionsContainer: {
                                 zIndex: 9999,
                                 fontSize: 25,
                                 top: 50,
                                 left:260,
                                 height: 0,
                             }
                         }}
            >
                {options?.map((option, index) => (
                    <>
                        {/*<MenuOption key={option.title} onSelect={option.onPress} navigationRef={option.navigationRef} text={option.title} style={styles.items} />*/}
                        <MenuOption key={option.title} onSelect={option.onPress}>
                            <Text style={styles.items}>{option.title}</Text>
                        </MenuOption>
                    </>
                ))}
            </MenuOptions>
        </Menu>
    );
};

export default memo(ItemOptions);

const styles = StyleSheet.create({
    container: {
        opacity: 111,
        zIndex: 9999,
        height: 100,
        float: 'right',
        marginLeft: "auto",
        display:'inline-block'
    },
    OptionsStyle: {
        float: 'right',
        marginLeft: "auto",
    },
    items: {
        zIndex: 100,
        fontSize: 16,
    }
});
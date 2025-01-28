import React, { useState } from 'react';
import { FlatList, StyleSheet, View } from 'react-native';
import { MenuProvider } from 'react-native-popup-menu';
import Item from './ThreeTopMenu';
import Animated from "react-native-reanimated";

export const OpportunityAttachmentView = () => {
    const [idMenu, setIdMenu] = useState('');

    const Items = [{ id: 'id', name: 'name', person: 'person', idMenu: idMenu, setIdMenu: setIdMenu }]

    return (
        <View style={styles.uploadContainer}>
            <Animated.Text style={styles.tabBarLabelStyle}>GGGGGGGGGGGG</Animated.Text>
            <MenuProvider  style={styles.MenuProvider}>
                <FlatList
                    data={Items}
                    keyExtractor={(item) => item.id}
                    renderItem={({ item }) => (
                        <Item
                            id={item.id}
                            idMenu={idMenu}
                            setIdMenu={setIdMenu}
                        />
                    )}
                />

            </MenuProvider>
        </View>
    );
};


const styles = StyleSheet.create({
    uploadContainer: {
        height: 178, // Выпадающего списка
        top: -45,
        left:'60%',
        position: 'absolute',
        float: 'right',
        marginLeft: "auto",
        zIndex: 9999999,
        opacity: 111
    },
});
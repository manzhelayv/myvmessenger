import {StyleSheet} from "react-native";

export const styles = StyleSheet.create({
    menuContainer: {
        flex: 1,
        color: 'black',
        width: 1000,
        maxWidth: "100%",
        borderRadius:15,
        position: 'relative',
        fontSize: 20,
    },
    container: {
        flex: 1,
        backgroundColor: '#fff',
        alignItems: 'center',
        borderTopColor: "#edf2ee",
        borderTopWidth: 1,
    },
    image: {
        width: 320,
        height: 440,
        borderRadius: 18,
    },
    footerContainer: {
        top: 20,
        bottom: 0,
        flex: 1 / 3,
        alignItems: 'center',
    },
    optionsRow: {
        alignItems: 'center',
        flexDirection: 'row',
        top: 20,
    },
});
import {StyleSheet} from "react-native";

export const styles = StyleSheet.create({
    menuContainer: {
        flex: 1,
        color: 'black',
        width: 1000,
        maxWidth: "100%",
        borderRadius:15,
        position: 'absolute',
        fontSize: 25,
        zIndex:9999,
        top: 5,
        backgroundColor: "white",
        height: "100%"
    },
    profileLink: {
        zIndex:9999,
        paddingLeft: 20,
        // borderBottomWidth: 1,
        // borderColor: 'black',
        paddingBottom: 10,
        position:'relative'

    },
    button: {
        paddingTop: 10,
        float: 'right',
    },
    text: {
        fontSize: 18,
        fontFamily: 'Wittgenstein-Italic-Regular'
    },
    textDescription: {
        fontSize: 14,
        fontFamily: 'Wittgenstein-Italic-Regular',
    },
    buttonIcon: {
        float: 'left',
        position:'absolute',
        top: 18,
        left: 0,
    },
    viewBlock: {
        left: 45
    },
});
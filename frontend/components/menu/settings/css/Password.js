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
        paddingTop: 10,
    },
    footerContainer: {
        top: 10,
        bottom: 0,
        flex: 1 / 3,
        alignItems: 'center',
    },
    message: {
        fontSize: 16
    },
    inputStyle: {
        padding: 0,
        fontFamily : 'RightLato',
        height: 55,
    },
    viewText: {
        width: 180,
        alignItems: 'center',
    },
});
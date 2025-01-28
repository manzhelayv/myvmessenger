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
    footerContainer: {
        top: 25,
        bottom: 0,
        flex: 1 / 3,
        alignItems: 'center',
    },
    viewText: {
        // borderBottomWidth: 1,
        width: 70,
        alignItems: 'center',
    },
    inputMaskStyle: {
        paddingBottom: 17,
        fontFamily : 'RightLato',
        borderBottomWidth: 1,
        borderBottomColor: '#86939E',
        width: '95%',
        fontSize: 18,
    },
    inputStyle: {
        padding: 0,
        fontFamily : 'RightLato',
        height: 55,
    },
});

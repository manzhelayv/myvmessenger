import {StyleSheet} from "react-native";

export const styles = StyleSheet.create({
    containerRegistration: {
        flex: 1,
        backgroundColor:'white',
        fontFamily : 'RightLato',
        borderTopColor: "#edf2ee",
        borderTopWidth: 1,
        width: "100%",
    },
    container: {
        top: 20,
        alignItems: 'center',
    },
    inputStyle: {
        padding: 0,
        fontFamily : 'RightLato',
        height: 55,
    },
    inputMaskStyle: {
        paddingBottom: 15,
        fontFamily : 'RightLato',
        borderBottomWidth: 1,
        borderBottomColor: '#86939E',
        width: '95%',
        fontSize: 18,
    },
    viewLoginResult: {
        alignItems: 'center',
        fontFamily : 'RightLato',
        marginTop: 10
    },
    footerContainer: {
        marginTop: 25,
        alignItems: 'center',
    },
    viewText: {
        // borderBottomWidth: 1,
        // position: 'absolute',
        // top: -10,
        alignItems: 'center',
    },
});
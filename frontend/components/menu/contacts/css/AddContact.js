import {StyleSheet} from "react-native";

export const styles = StyleSheet.create({
    addContacts: {
        height:50,
        paddingTop: 20,
        borderTopColor: "#edf2ee",
        borderTopWidth: 1
    },
    textArea: {
        width: "52%",
        height: 31,
        color: 'black',
        borderColor: 'black',
        borderWidth: 3,
        borderRadius: 10,
        padding: 5,
        marginBottom: 10,
    },
    addContact: {
        float: 'left',
        left: 10,
        display: 'block',
    },
    addButton: {
        position: 'absolute',
        float: 'right',
        right:-10,
        top: 20,
    },
    viewLoginResult: {
        // backgroundColor:'white',
        justifyContent: 'center',
        alignItems: 'center',
        fontFamily : 'RightLato',
        fontSize: 10,
        width: 300,
        position: 'absolute',
        bottom: 60,
        left: -12
    },
});
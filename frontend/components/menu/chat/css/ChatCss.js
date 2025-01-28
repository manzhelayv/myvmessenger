import {StyleSheet} from "react-native";

export const styles = StyleSheet.create({
    containerTop: {
        flex: 1,
        // width: "100%",
        backgroundColor: '#fff',
        // alignItems: 'center',
        borderTopColor: "#edf2ee",
        borderTopWidth: 1
    },
    image: {
        flex: 1,
        resizeMode: 'cover',
        justifyContent: 'center',
    },
    container: {
        flex: 1,
        alignItems: 'center',
    },
    textArea: {
        width: 160,
        height: 31,
        color: 'black',
        borderColor: 'black',
        borderWidth: 3,
        borderRadius: 10,
        padding: 5,
        marginTop: 11
    },
    addContact: {
        float: 'left',
        left: 10,
        backgroundColor: 'white',
        height: 100,
    },
    addButton: {
        position: 'absolute',
        float: 'right',
        right:-10,
        marginTop: 11
    },
    tabBarLabelStyle: {
        fontFamily : 'RightLato',
        fontSize: 22,
        textTransform: 'none',
        left: 100,
        fontWeight: 600,
        top: -62,
        alignItems: 'center',
        justifyContent: 'center',
    },
    buttonIconChat: {
        position: 'absolute',
        left: 10,
        top: 18
    },
    ChatView: {
        top: -19,
        left: 40,
        marginRight: "auto",
        float: "left",
    },
    GoBack: {
        top: 33,
        left: 5,
        marginRight: "auto",
        float: "left",
        fontSize: 30,
    },
});
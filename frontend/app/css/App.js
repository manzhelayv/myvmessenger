import {StyleSheet} from "react-native";

const styles = StyleSheet.create({
    containerMenu: {
        top:7,
        float: 'left',
        marginRight: "auto",
        zIndex: 100,
        height:50,
        width: '100%',
        fontSize:30,
    },
    buttonIconChat: {
        position: 'absolute',
        left: 28,
    },
    buttonIconSettings: {
        position: 'absolute',
        left: 10,
    },
    buttonIconContacts: {
        position: 'absolute',
        left: 10,
    },
    buttonIconRegistration: {
        position: 'absolute',
        right: 60,
        top: 30
    },
    buttonIconLogin: {
        position: 'absolute',
        right: 40,
        top: 30
    },
    tabBarLabelStyle: {
        fontFamily : 'RightLato',
        fontSize: 18,
        textTransform: 'none',
        left: 12,
        fontWeight: 600,
        top: -2,
        alignItems: 'center',
        justifyContent: 'center',
    },
    tabBarLabelStyleChat: {
        fontFamily : 'RightLato',
        fontSize: 22,
        textTransform: 'none',
        left: 40,
        zIndex:1,
        top: -70,
        alignItems: 'center',
        justifyContent: 'center',
    },
    ChatView: {
        top: -30,
        right: 18,
        marginRight: "auto",
        float: "left",
    },
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
        flex: 2,
        alignItems: 'center',
        position: 'absolute',
        backgroundColor: 'white',
        height: "100%",
        width: "100%",
        borderRadius:15
    },
    containerCenter: {
        height: "100%",
    },
    header: {
        width: "100%",
        height:"5%",
        backgroundColor: "black"
    },
    footer: {
        position: 'absolute',
        width: "100%",
        height:"2%",
        top: "98%",
        backgroundColor: "black"
    }
});
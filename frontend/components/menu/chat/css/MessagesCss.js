import {StyleSheet} from "react-native";

export const styles = StyleSheet.create({
    rightChatText: {
        width: 'auto',
        maxWidth: 260,
        float: 'right',
        //   position: "absolute",
        zIndex: 10,
        right: 8,
        top: 0,
        borderColor: '#03fc8c',
        borderWidth: 7,
        borderRadius: 10,
        padding: 5,
        paddingTop: 0,
        //   marginTop:5,
        marginBottom:10,
        backgroundColor: '#03fc5e',
        display: 'block',
        //   left: 'auto',
        marginLeft: "auto",
        fontFamily: 'Wittgenstein-Italic-Regular'
    },
    leftChatText: {
        width: 'auto',
        maxWidth: 260,
        float: 'left',
        top: 0,
        borderColor: '#03fc8c',
        borderWidth: 7,
        borderRadius: 10,
        padding: 5,
        paddingTop: 0,
        //   marginTop:10,
        marginBottom:10,
        backgroundColor: '#03fc5e',
        display: 'block',
        fontFamily: 'Wittgenstein-Italic-Regular',
        left: 8,
        marginRight: "auto",
    },
    scrollView: {
        marginHorizontal: 20,
        marginRight: 0,
        backgroundColor: 'green',
    },
    buttonTextDate: {
        fontFamily : 'RightLato',
        fontSize:11,
    },
    buttonText: {
        fontFamily : 'RightLato',
        fontWeight: 'bold',
        fontSize:12,
    },
    leftText: {
        float: 'left',
        marginRight: "auto",
    },
    rightTextDate: {
        float: 'right',
        marginLeft: "auto",
    },
    leftImage: {
        marginLeft: -7,
    },
});

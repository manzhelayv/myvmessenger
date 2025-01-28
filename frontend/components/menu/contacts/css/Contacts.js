import {StyleSheet} from "react-native";

export const styles = StyleSheet.create({
    container: {
        flex: 1,
        backgroundColor:'white',
        fontFamily : 'RightLato',
        borderTopColor: "#edf2ee",
        borderTopWidth: 1,
        width: "100%",
        // height: "90%",
    },
    scrollView: {
        marginHorizontal: 20,
        marginRight: 0,
        backgroundColor: 'green',
    },
    textView: {
        left: 15,
        height: 67,
    },
    rightTextView: {
        left: 60,
        top:-35,
    },
    rightText: {
        fontSize: 16,
        fontFamily: 'RightLato',
    },
    userNameText: {
        fontSize: 14,
        fontFamily: 'DimkinBold',
        fontWeight: 400,
        top: 10,
    },
    containerContacts: {
        height: "82%",
        top: -2,
    },
    GoBack: {
        top: 33,
        left: 5,
        marginRight: "auto",
        float: "left",
        fontSize: 30,
    },
    tabBarLabelStyle: {
        fontFamily : 'RightLato',
        fontSize: 22,
        textTransform: 'none',
        left: 40,
        fontWeight: 600,
        top: -11,
        alignItems: 'center',
        justifyContent: 'center',
    },
    addButton: {
        position: 'absolute',
        float: 'left',
        marginRight: "auto",
        top: 10,
        left: "58%",
        zIndex: 9999,
    },
    updateText: {
        color: "#000",
        top: 7,
        alignItems: 'center',
        justifyContent: 'center',
    },
});
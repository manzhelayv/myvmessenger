import React from "react";
import {globalOnPressNotification, globalNotification, globalOnPress} from "../../components/menu/TopMenu";
import Chat from "../../components/menu/chat/Chat";

export default function TabLayout() {

    return <Chat onPress={globalOnPress} notification={globalNotification} onPressNotification={globalOnPressNotification} />
}

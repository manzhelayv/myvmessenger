import React from "react";
import {globalLatsMessageForChats, globalOnPress} from "../../components/menu/TopMenu";
import Chats from "../../components/menu/chats/Chats";

export default function TabLayout() {
    return <Chats onPress={globalOnPress} latsMessageForChats={globalLatsMessageForChats} navigationRef="" />
}

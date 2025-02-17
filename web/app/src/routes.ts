import EntryMappingList from "./components/EntryMappingList.vue";
import Login from "./components/Login.vue";

export const routes = [
    {path: '/', component: EntryMappingList, name: "Entry Mapping"},
    {path: '/login', component: Login, name: "Login"},
]

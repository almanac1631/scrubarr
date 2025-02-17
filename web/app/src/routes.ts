import EntryMappingList from "./components/EntryMappingList.vue";
import Login from "./components/Login.vue";
import {RouteRecordRaw} from "vue-router";

import 'vue-router';

declare module 'vue-router' {
    interface RouteMeta {
        displayedInNavigation?: boolean;
    }
}

export const routes: Array<RouteRecordRaw> = [
    {path: '/', component: EntryMappingList, name: "Entry Mappings", meta: {displayedInNavigation: true}},
    {path: '/login', component: Login, name: "Login", meta: {displayedInNavigation: false}},
]

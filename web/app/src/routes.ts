import EntryMappingList from "./components/EntryMappingList.vue";
import Login from "./components/Login.vue";
import {RouteRecordRaw} from "vue-router";

import 'vue-router';
import Main from "./components/Main.vue";

declare module 'vue-router' {
    interface RouteMeta {
        displayedInNavigation?: boolean;
    }
}

export const routes: Array<RouteRecordRaw> = [
    {path: '/', component: Main, name: "Overview", meta: {displayedInNavigation: true}},
    {path: '/entry-mappings', component: EntryMappingList, name: "Files", meta: {displayedInNavigation: true}},
    {path: '/login', component: Login, name: "Login", meta: {displayedInNavigation: false}},
]

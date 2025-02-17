import {createApp} from 'vue'
import './style.css'
import App from './App.vue'
import {createRouter, createWebHistory} from "vue-router";
import {routes} from "./routes.ts";
import {isAuthenticated} from "./auth/auth.ts";

const router = createRouter({
    history: createWebHistory(),
    routes,
});

router.beforeEach((to) => {
    if (to.path === "/login") {
        return;
    }
    if (!isAuthenticated()) {
        return {path: "/login"};
    }
})

createApp(App).use(router).mount('#app')

import {createApp} from 'vue'
import './style.css'
import App from './App.vue'
import {createRouter, createWebHistory} from "vue-router";
import {routes} from "./routes.ts";
import {checkAndInitAuthentication} from "./auth/auth.ts";

const router = createRouter({
    history: createWebHistory(),
    routes,
});

router.beforeEach((to) => {
    const alreadyAuthenticated = checkAndInitAuthentication();
    const isLoginPath = to.path === "/login"
    if (isLoginPath && alreadyAuthenticated) {
        return {path: "/"};
    } else if (!isLoginPath && !alreadyAuthenticated) {
        return {path: "/login"};
    }
})

createApp(App).use(router).mount('#app')

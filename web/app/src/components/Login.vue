<script setup lang="ts">
import {ref} from "vue";
import {Configuration, DefaultApiFactory} from "../api";
import {AxiosError} from "axios";
import {notify} from "../utils/notificationList.ts";
import {initializeAuthToken} from "../auth/auth.ts";
import {useRouter} from "vue-router";
import {basePath} from "../utils/api.ts";

const username = ref("");
const password = ref("");

const router = useRouter();

async function submitLoginForm(e: any) {
  e.preventDefault();
  const unauthenticatedApiClient = DefaultApiFactory(new Configuration({basePath}));
  try {
    const loginResp = await unauthenticatedApiClient.login({
      username: username.value,
      password: password.value,
    });
    initializeAuthToken(loginResp.data.token);
    await router.push({path: "/"});
  } catch (e) {
    if (e instanceof AxiosError && e.status === 401) {
      notify("Invalid login credentials.");
      return;
    }
    console.log("unknown error occurred while logging in");
    console.log(e);
  }
}
</script>

<template>
  <div class="bg-gray-800 h-full">
    <div class="flex min-h-full flex-col justify-center px-6 py-12 lg:px-8 overflow-y-auto">
      <div class="sm:mx-auto sm:w-full sm:max-w-sm">
        <img class="mx-auto h-20 w-auto" src="/logo-text.svg" alt="logo of scrubarr">
      </div>

      <div class="mt-10 sm:mx-auto sm:w-full sm:max-w-sm">
        <form class="space-y-6" action="#" method="POST" @submit="submitLoginForm">
          <div>
            <label for="username" class="block text-sm/6 font-medium text-gray-100">Username</label>
            <div class="mt-2">
              <input type="text" name="username" id="username" autocomplete="username" required v-model="username"
                     class="block w-full rounded-md bg-white px-3 py-1.5 text-base text-gray-900 outline-1 -outline-offset-1 outline-gray-300 placeholder:text-gray-400 focus:outline-2 focus:-outline-offset-2 focus:outline-red-500 sm:text-sm/6">
            </div>
          </div>

          <div>
            <div class="flex items-center justify-between">
              <label for="password" class="block text-sm/6 font-medium text-gray-100">Password</label>
            </div>
            <div class="mt-2">
              <input type="password" name="password" id="password" autocomplete="current-password" required v-model="password"
                     class="block w-full rounded-md bg-white px-3 py-1.5 text-base text-gray-900 outline-1 -outline-offset-1 outline-gray-300 placeholder:text-gray-400 focus:outline-2 focus:-outline-offset-2 focus:outline-red-500 sm:text-sm/6">
            </div>
          </div>

          <div>
            <button type="submit"
                    class="flex w-full justify-center rounded-md bg-red-500 px-3 py-1.5 text-sm/6 font-semibold text-gray-100 shadow-xs hover:bg-red-600 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-red-500">
              Sign in
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<style scoped>

</style>

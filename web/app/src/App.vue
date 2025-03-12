<script setup lang="ts">
import {Disclosure} from '@headlessui/vue'
import {RouterLink} from "vue-router";
import {routes} from "./routes.ts";
import Notification from "./components/Notification.vue";
import {notificationList} from "./utils/notificationList.ts";
import {logout} from "./auth/auth.ts";
import RefreshButton from "./components/RefreshButton.vue";
import InfoFooter from "./components/InfoFooter.vue";
</script>

<template>
  <div class="fixed right-0 z-10 pt-14 pr-2">
    <TransitionGroup tag="div" enter-active-class="transition-opacity">
      <Notification v-for="notification in notificationList" :notification="notification">
      </Notification>
    </TransitionGroup>
  </div>
  <div class="flex flex-col h-full">
    <Disclosure as="nav" class="bg-gray-800" v-if="$route.path !== '/login'">
      <div class="mx-auto max-w-7xl px-2 sm:px-6 lg:px-8">
        <div class="relative flex h-16 items-center justify-between">
          <div class="flex flex-1 items-center justify-center sm:items-stretch sm:justify-start">
            <div class="flex flex-shrink-0 items-center text-gray-100">
              scrubarr
            </div>
            <div class="hidden sm:ml-6 sm:block">
              <div class="flex space-x-4">
                <RouterLink :to="item.path" v-for="item in routes.filter(route => route.meta?.displayedInNavigation)"
                            activeClass='bg-gray-900 text-white'
                            class='text-gray-300 hover:bg-gray-700 hover:text-white rounded-md px-3 py-2 text-sm font-medium'
                            :aria-current="$route.path == item.path">
                  {{ item.name }}
                </RouterLink>
              </div>
            </div>
          </div>

          <div class="flex space-x-4 items-end justify-center sm:items-stretch sm:justify-start">
            <RefreshButton></RefreshButton>
            <div class="flex items-center">
              <a aria-current="true" href="#" @click="logout"
                 class="text-gray-300 hover:bg-gray-700 hover:text-white rounded-md px-3 py-2 text-sm font-medium">
                Logout
              </a>
            </div>
          </div>
        </div>
      </div>
    </Disclosure>
    <main class="bg-gray-100 flex-grow relative" v-if="$route.meta.displayedInNavigation">
      <div class="py-2">
        <div class="container mx-auto text-gray-900">
          <h1 class="text-3xl my-5">
            {{ $route.name }}
          </h1>
          <RouterView/>
        </div>
      </div>
      <InfoFooter></InfoFooter>
    </main>
    <RouterView v-else class="flex-grow"/>
  </div>
</template>

<style>
/*noinspection CssUnusedSymbol*/
html, body, #app {
  height: 100%;
}
</style>

<script setup lang="ts">
import {Disclosure} from '@headlessui/vue'
import {RouterLink} from "vue-router";
import {routes} from "./routes.ts";
</script>

<template>
  <Disclosure as="nav" class="bg-gray-800">
    <div class="mx-auto max-w-7xl px-2 sm:px-6 lg:px-8">
      <div class="relative flex h-16 items-center justify-between">
        <div class="flex flex-1 items-center justify-center sm:items-stretch sm:justify-start">
          <div class="flex flex-shrink-0 items-center text-gray-100">
            scrubarr
          </div>
          <div class="hidden sm:ml-6 sm:block">
            <div class="flex space-x-4">
              <RouterLink :to="item.path" v-for="item in routes" :key="item.meta.title"
                          :class="[item.meta.current ? 'bg-gray-900 text-white' : 'text-gray-300 hover:bg-gray-700 hover:text-white', 'rounded-md px-3 py-2 text-sm font-medium']"
                          :aria-current="item.meta.current ? 'page' : undefined">
                {{ item.meta.title }}
              </RouterLink>
            </div>
          </div>
        </div>
      </div>
    </div>
  </Disclosure>
  <main class="bg-gray-100 py-12">
    <div class="container mx-auto rounded-md bg-white p-10 shadow">
      <h1 class="text-2xl my-3">
        Entry Mappings
      </h1>
      <RouterView/>
    </div>
  </main>
</template>

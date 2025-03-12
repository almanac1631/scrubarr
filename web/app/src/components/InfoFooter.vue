<script setup lang="ts">
import {Configuration, DefaultApiFactory} from "../api";
import {basePath} from "../utils/api.ts";
import {onMounted, ref} from "vue";

const unauthenticatedApiClient = DefaultApiFactory(new Configuration({basePath}));

const apiVersion = ref<string | null>(null);

onMounted(() => {
  unauthenticatedApiClient.getInfo().then((response) => {
    apiVersion.value = response.data.info.version;
  });
});
</script>

<template>
  <div class="h-8"></div>
  <div class="absolute bottom-2 w-full text-center">
    Made with ❤️ by almanac1631 - scrubarr {{ apiVersion }}
  </div>
</template>

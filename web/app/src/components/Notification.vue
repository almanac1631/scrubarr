<script setup lang="ts">


import {PropType} from "vue";
import {NotificationEntry} from "../utils/notificationList.ts";

const props = defineProps({
  notification: Object as PropType<NotificationEntry>
});
</script>

<template>
  <div class="bg-red-200 px-6 py-4 mx-2 my-4 rounded-md text-lg flex items-center mx-auto max-w-lg">
    <svg class="text-red-600 w-5 h-5 sm:w-5 sm:h-5 mr-3" aria-hidden="true"
         xmlns="http://www.w3.org/2000/svg" width="24" height="24" fill="currentColor" viewBox="0 0 24 24">
      <path fill-rule="evenodd"
            d="M2 12C2 6.477 6.477 2 12 2s10 4.477 10 10-4.477 10-10 10S2 17.523 2 12Zm7.707-3.707a1 1 0 0 0-1.414 1.414L10.586 12l-2.293 2.293a1 1 0 1 0 1.414 1.414L12 13.414l2.293 2.293a1 1 0 0 0 1.414-1.414L13.414 12l2.293-2.293a1 1 0 0 0-1.414-1.414L12 10.586 9.707 8.293Z"
            clip-rule="evenodd"/>
    </svg>
    <span class="text-red-800">{{ props.notification?.message }}</span>
  </div>
</template>

<style scoped>

</style>
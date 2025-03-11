<script setup lang="ts">


import {computed, PropType} from "vue";
import {NotificationEntry, NotificationType} from "../utils/notificationList.ts";

const props = defineProps({
  notification: Object as PropType<NotificationEntry>
});

const bgClass = computed(() => {
  return {
    'bg-red-200': props.notification?.type === NotificationType.Error,
    'bg-green-200': props.notification?.type === NotificationType.Success,
    'bg-blue-200': props.notification?.type === NotificationType.Info,
  }
});
const iconClass = computed(() => {
  return {
    'text-red-600': props.notification?.type === NotificationType.Error,
    'text-green-600': props.notification?.type === NotificationType.Success,
    'text-blue-600': props.notification?.type === NotificationType.Info,
  }
});
const textClass = computed(() => {
  return {
    'text-red-800': props.notification?.type === NotificationType.Error,
    'text-green-800': props.notification?.type === NotificationType.Success,
    'text-blue-800': props.notification?.type === NotificationType.Info,
  }
});
</script>

<template>
  <div class="px-6 py-4 my-4 rounded-md text-lg flex items-center mx-auto max-w-lg" :class="bgClass">
    <svg class="w-5 h-5 sm:w-5 sm:h-5 mr-3" :class="iconClass" aria-hidden="true"
         xmlns="http://www.w3.org/2000/svg" width="24" height="24" fill="currentColor" viewBox="0 0 24 24">
      <path v-if="props.notification?.type == NotificationType.Error"
            fill-rule="evenodd"
            d="M2 12C2 6.477 6.477 2 12 2s10 4.477 10 10-4.477 10-10 10S2 17.523 2 12Zm7.707-3.707a1 1 0 0 0-1.414 1.414L10.586 12l-2.293 2.293a1 1 0 1 0 1.414 1.414L12 13.414l2.293 2.293a1 1 0 0 0 1.414-1.414L13.414 12l2.293-2.293a1 1 0 0 0-1.414-1.414L12 10.586 9.707 8.293Z"
            clip-rule="evenodd"/>
      <path v-if="props.notification?.type == NotificationType.Success"
            fill-rule="evenodd"
            d="M17 3.34a10 10 0 1 1 -14.995 8.984l-.005 -.324l.005 -.324a10 10 0 0 1 14.995 -8.336zm-1.293 5.953a1 1 0 0 0 -1.32 -.083l-.094 .083l-3.293 3.292l-1.293 -1.292l-.094 -.083a1 1 0 0 0 -1.403 1.403l.083 .094l2 2l.094 .083a1 1 0 0 0 1.226 0l.094 -.083l4 -4l.083 -.094a1 1 0 0 0 -.083 -1.32z"
            clip-rule="evenodd"/>
      <path v-if="props.notification?.type == NotificationType.Info"
            fill-rule="evenodd"
            d="M12 2c5.523 0 10 4.477 10 10a10 10 0 0 1 -19.995 .324l-.005 -.324l.004 -.28c.148 -5.393 4.566 -9.72 9.996 -9.72zm0 9h-1l-.117 .007a1 1 0 0 0 0 1.986l.117 .007v3l.007 .117a1 1 0 0 0 .876 .876l.117 .007h1l.117 -.007a1 1 0 0 0 .876 -.876l.007 -.117l-.007 -.117a1 1 0 0 0 -.764 -.857l-.112 -.02l-.117 -.006v-3l-.007 -.117a1 1 0 0 0 -.876 -.876l-.117 -.007zm.01 -3l-.127 .007a1 1 0 0 0 0 1.986l.117 .007l.127 -.007a1 1 0 0 0 0 -1.986l-.117 -.007z"
            clip-rule="evenodd"/>
    </svg>
    <span :class="textClass">{{ props.notification?.message }}</span>
  </div>
</template>

<style scoped>

</style>
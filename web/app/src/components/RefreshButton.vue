<script setup lang="ts">
import {getApiClient} from "../utils/api.ts";
import {ref} from "vue";
import {NotificationType, notify} from "../utils/notificationList.ts";

const refreshInProgress = ref(false);

async function refreshEntryMapping() {
  notify("Refreshing entry mapping...", NotificationType.Info);
  refreshInProgress.value = true;
  try {
    await getApiClient().refreshEntryMappings();
    notify("Entry mapping refreshed successfully.", NotificationType.Success);
  } catch (e) {
    notify("Failed to refresh entry mapping.", NotificationType.Error);
    console.error("unknown error occurred while requesting a refresh in");
    console.error(e);
  } finally {
    refreshInProgress.value = false;
  }
}
</script>

<template>
  <div class="flex items-center">
    <button @click="refreshEntryMapping" title="Refresh entry mapping"
            class="text-gray-300 px-3 py-2 disabled:text-gray-100"
            :disabled="refreshInProgress"
            :class="{ 'hover:bg-gray-700': !refreshInProgress, 'hover:text-white rounded-md': !refreshInProgress }">
      <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none"
           stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"
           :class="{ 'animate-spin': refreshInProgress }">
        <path stroke="none" d="M0 0h24v24H0z" fill="none"/>
        <path d="M20 11a8.1 8.1 0 0 0 -15.5 -2m-.5 -4v4h4"/>
        <path d="M4 13a8.1 8.1 0 0 0 15.5 2m.5 4v-4h-4"/>
      </svg>
    </button>
  </div>
</template>

<style scoped>

</style>
<script setup lang="ts">

import {computed, onMounted, ref, Ref} from "vue";
import {Stats} from "../api";
import {getApiClient} from "../utils/api.ts";
import {formatFileSize} from "../utils/fileSize.ts";

const stats: Ref<Stats | null> = ref(null);


const apiClient = getApiClient();

onMounted(async () => {
  stats.value = (await apiClient.getStats()).data.stats;
});

const percentage = computed(() => {
  if (stats.value === null) {
    return;
  }
  const diskSpace = stats.value.diskSpace;
  return (diskSpace.bytesUsed / diskSpace.bytesTotal) * 100;
});
</script>

<template>
  <div class="grid grid-cols-3 gap-4">
    <div class="rounded-md bg-white p-4 shadow">
      <p class="text-lg mb-1">Storage use</p>
      <div class="relative rounded-xl overflow-hidden">
        <div class="h-6 bg-gray-200">
        </div>
        <div :style="'width: ' + percentage?.toFixed(2) + '%'"
             class="h-6 bg-red-500 absolute left-0 top-0 z-10 px-2 flex items-center justify-end">
          <div class="text-sm">{{ percentage?.toFixed(2) }}%</div>
        </div>
      </div>
      <div class="text-xs p-1">{{ formatFileSize(stats?.diskSpace.bytesUsed) }}/{{
          formatFileSize(stats?.diskSpace.bytesTotal)
        }}
      </div>
    </div>
  </div>
</template>

<style scoped>

</style>
<script setup lang="ts">

import {onMounted, ref} from "vue";
import {getApiClient} from "../../utils/api.ts";

const props = defineProps({
  entryId: String
});

function highlightJSON(obj: unknown): string {
  const jsonString = JSON.stringify(obj, null, 2);

  return jsonString
      .replace(/(&)/g, "&amp;") // escape HTML special chars
      .replace(/</g, "&lt;")
      .replace(/>/g, "&gt;")
      .replace(
          /("(\\u[\da-fA-F]{4}|\\[^u]|[^\\"])*"(\s*:)?|\b(true|false|null)\b|-?\d+(?:\.\d*)?(?:[eE][+\-]?\d+)?)/g,
          match => {
            let cls = "text-gray-300"; // default color
            if (/^"/.test(match)) {
              cls = /:$/.test(match) ? "text-blue-400" : "text-green-400";
            } else if (/true|false/.test(match)) {
              cls = "text-purple-400";
            } else if (/null/.test(match)) {
              cls = "text-pink-400";
            } else if (/^-?\d+/.test(match)) {
              cls = "text-yellow-400";
            }
            return `<span class="${cls}">${match}</span>`;
          }
      );
}

type RetrieverApiResponse = {
  retrieverId: string;
  active: boolean;
  apiResponse: string;
};

const apiClient = getApiClient();

const highlightedRetrieverApiResponses = ref([] as RetrieverApiResponse[]);

onMounted(async () => {
  if (!props.entryId) {
    return [];
  }
  const details = await apiClient.getEntryMappingDetails(props.entryId);
  let firstActive = true;
  for (const retrieverDetail of details.data.entry.retrieverDetails) {
    highlightedRetrieverApiResponses.value.push({
      retrieverId: retrieverDetail.id,
      active: firstActive,
      apiResponse: highlightJSON(JSON.parse(retrieverDetail.apiResp))
    });
    firstActive = false;
  }
});

function activateEntry(entry: RetrieverApiResponse) {
  highlightedRetrieverApiResponses.value?.forEach(e => e.active = false);
  entry.active = true;
}
</script>

<template>
  <tr>
    <td colspan="3" class="px-3 py-3">
      <div class="flex gap-3 border-b-1 border-b-gray-200">
        <template v-for="entry in highlightedRetrieverApiResponses">
          <div
              :class="[ entry.active ? 'border-b-red-500 bg-gray-50' : 'border-b-transparent' ]"
              class="cursor-pointer p-2 text-gray-400 hover:bg-gray-50 rounded-t-md border-b-2 hover:border-b-red-500"
              v-on:click="activateEntry(entry)"
          >
            {{ entry.retrieverId }}
          </div>
        </template>
      </div>
      <div class="flex overflow-x-auto">
        <template v-for="entry in highlightedRetrieverApiResponses">
          <div v-if="entry.active" class="flex-1 p-2 bg-gray-50 rounded-b-md">
            <pre v-html="entry.apiResponse" class="text-sm"></pre>
          </div>
        </template>
      </div>
    </td>
  </tr>
</template>

<style scoped>

</style>
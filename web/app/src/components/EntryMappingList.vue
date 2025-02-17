<script setup lang="ts">
import {computed, onMounted, Ref, ref, watch} from "vue";
import {EntryMapping, GetEntryMappingsFilterEnum, Retriever} from "../api";
import PreloaderTableEntry from "./entry-mapping-list/PreloaderTableEntry.vue";
import Dropdown, {DropdownOption} from "./common/Dropdown.vue";
import Pagination from "./common/Pagination.vue";
import {getCategoriesFromRetrieverList, sortRetrieverList} from "../utils/retrievers.ts";
import TableRetrieverStateHeader from "./entry-mapping-list/TableRetrieverStateHeader.vue";
import TableRetrieverStateRowEntry from "./entry-mapping-list/TableRetrieverStateRowEntry.vue";
import {getApiClient} from "../utils/api.ts";

const contentLoaded = ref(false);
const entryMappingList: Ref<Array<EntryMapping> | null> = ref(null);
const entryMappingTotalAmount: Ref<number | null> = ref(null);
const retrieverList: Ref<Array<Retriever> | null> = ref(null);

function isEntryPresentInRetriever(entryMapping: EntryMapping, retrieverId: string): boolean {
  return entryMapping.retrieverFindings?.filter(finding => {
    return finding.id == retrieverId;
  }).length > 0;
}

function isEntryPresentInRetrieverCategory(entryMapping: EntryMapping, retrieverCategoryName: String): boolean {
  return entryMapping.retrieverFindings?.filter(finding => {
    const matchedRetriever = retrieverList.value?.filter((retriever) => retriever.id === finding.id)[0];
    return matchedRetriever?.category == retrieverCategoryName;
  }).length > 0;
}

function toggleAll() {
  if (areAllSelected()) {
    selectedItems.value = [];
  } else {
    selectedItems.value = entryMappingList.value ?? [];
  }
}

const selectedItems: Ref<Array<EntryMapping>> = ref([]);

const selectedPage: Ref<number> = ref(1);

function areAllSelected() {
  return entryMappingList.value?.length === selectedItems.value.length;
}

const apiClient = getApiClient();

async function fetchAndDisplayEntries() {
  if (retrieverList.value?.values === null) {
    throw Error("cannot fetch and display entries if retrievers are not resolved yet");
  }

  if (selectedFilter.value === null || selectedPageSize.value === null || selectedPageSize.value.value === undefined) {
    return;
  }
  try {
    const requestPage = selectedPage.value - 1;
    const entryMappingResponse = (await apiClient.getEntryMappings(requestPage, +selectedPageSize.value.value, selectedFilter.value.value as GetEntryMappingsFilterEnum)).data;
    entryMappingList.value = entryMappingResponse.entries;
    entryMappingTotalAmount.value = entryMappingResponse.totalAmount;
  } finally {
    contentLoaded.value = true;
  }
}

const retrieverGroupingEnabled = ref(false);
const retrieverCategoryList = computed(() => {
  if (retrieverList.value === null) {
    return;
  }
  return getCategoriesFromRetrieverList(retrieverList.value);
});

onMounted(async () => {
  const retrieverListResp = (await apiClient.getRetrievers()).data.retrievers;
  sortRetrieverList(retrieverListResp);
  retrieverList.value = retrieverListResp;
});

const filterElemList = [
  {displayName: "Filter by", value: undefined},
  {displayName: "No filter", value: undefined},
  {displayName: "Incomplete entries", value: GetEntryMappingsFilterEnum.IncompleteEntries},
  {displayName: "Complete entries", value: GetEntryMappingsFilterEnum.CompleteEntries}
];

const selectedFilter: Ref<DropdownOption | null> = ref(null);

const pageSizeElemList = [
  {displayName: "Page size", value: "5"},
  {displayName: "10", value: "10"},
  {displayName: "25", value: "25"},
  {displayName: "50", value: "50"},
  {displayName: "100", value: "100"}
];

const selectedPageSize: Ref<DropdownOption | null> = ref(null);

watch([selectedFilter, selectedPageSize, selectedPage], () => {
  fetchAndDisplayEntries();
});
</script>

<template>
  <div class="h-full bg-gray-100 py-12">
    <div class="container mx-auto rounded-md bg-white p-10 shadow">
      <h1 class="text-2xl my-3">
        Entry Mappings
      </h1>
      <div class="my-2 flex justify-between">
        <div class="flex items-center">
          <div>
            <label class="inline-flex items-center cursor-pointer">
              <span class="me-3 text-sm font-medium text-gray-900">Group by retriever category</span>
              <input type="checkbox" value="" class="sr-only peer" v-model="retrieverGroupingEnabled">
              <span
                  class="relative w-11 h-6 bg-gray-200 peer-focus:ring-2 rounded-full peer peer-checked:after:translate-x-full rtl:peer-checked:after:-translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:start-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-red-500"></span>
            </label>
          </div>
        </div>
        <div class="my-2 flex justify-end gap-2">
          <Dropdown :options="pageSizeElemList" :default-option="pageSizeElemList[0]" v-model="selectedPageSize"/>
          <Dropdown :options="filterElemList" :default-option="filterElemList[0]" v-model="selectedFilter"/>
        </div>
      </div>
      <table class="table-fixed w-full">
        <thead>
        <tr class="text-left border-b-2" v-if="contentLoaded">
          <th class="w-[50px]">
            <div class="flex justify-center">
              <div class="relative h-5 w-5">
                <input type="checkbox"
                       class="peer h-5 w-5 cursor-pointer transition-all appearance-none rounded border border-slate-300"
                       :checked="areAllSelected()" @change="toggleAll">
                <div
                    class="flex justify-center items-center absolute top-0 left-0 z-10 h-5 w-5 pointer-events-none peer-checked:opacity-100 opacity-0 transition-opacity">
                  <svg class="h-3 w-3 text-red-500" viewBox="0 0 24 24" fill="none" stroke="currentColor"
                       stroke-width="5"
                       stroke-linecap="round" stroke-linejoin="round">
                    <polyline points="20 6 9 17 4 12"/>
                  </svg>
                </div>
              </div>
            </div>
          </th>
          <th class="py-3 pr-3 font-medium">
            Name
          </th>
          <th v-if="retrieverGroupingEnabled" v-for="retrieverCategory in retrieverCategoryList"
              class="w-[120px] p-3 font-medium text-center truncate">
            <div class="h-6 flex justify-center">
              <div class="relative">
                {{ retrieverCategory.displayName }}
              </div>
            </div>
          </th>
          <th v-else v-for="retriever in retrieverList" class="w-[100px] p-3 font-medium text-center">
            <TableRetrieverStateHeader
                :name="retriever.name"
                :hover-text="retriever.softwareName"
                :logo-filename="`${retriever.softwareName}-128x128.png`"
                :logo-alt-text="`The logo of the ${retriever.softwareName} software project.`"
            />
          </th>
        </tr>
        <PreloaderTableEntry class="border-b-2" v-else/>
        </thead>
        <tbody>
        <tr v-for="entryMapping in entryMappingList" class="hover:bg-stone-100 border-t" v-if="contentLoaded">
          <td class="py-3 pl-3 pr-3">
            <div class="flex justify-center">
              <div class="relative h-5 w-5">
                <input type="checkbox"
                       class="peer h-5 w-5 cursor-pointer transition-all appearance-none rounded border border-slate-300"
                       v-bind:value="entryMapping" v-model="selectedItems">
                <div
                    class="flex justify-center items-center absolute top-0 left-0 z-10 h-5 w-5 pointer-events-none peer-checked:opacity-100 opacity-0 transition-opacity">
                  <svg class="h-3 w-3 text-red-500" viewBox="0 0 24 24" fill="none" stroke="currentColor"
                       stroke-width="5"
                       stroke-linecap="round" stroke-linejoin="round">
                    <polyline points="20 6 9 17 4 12"/>
                  </svg>
                </div>
              </div>
            </div>
          </td>
          <td class="py-3 pr-3 font-medium truncate" :title="entryMapping.name">
            {{ entryMapping.name }}
          </td>

          <TableRetrieverStateRowEntry
              v-if="retrieverGroupingEnabled" v-for="retrieverCategory in retrieverCategoryList"
              :present="isEntryPresentInRetrieverCategory(entryMapping, retrieverCategory.name)"
          />
          <TableRetrieverStateRowEntry
              v-else v-for="retriever in retrieverList"
              :present="isEntryPresentInRetriever(entryMapping, retriever.id)"
          />
        </tr>
        <PreloaderTableEntry v-for="_ in 10" v-else/>
        </tbody>
      </table>
      <Pagination
          v-if="contentLoaded && entryMappingTotalAmount !== null &&selectedPageSize !== null &&  selectedPageSize.value !== undefined"
          :page-size="+selectedPageSize.value"
          :selected-page="1"
          :total-amount-of-items="entryMappingTotalAmount" v-model="selectedPage"/>
    </div>
  </div>
</template>

<style scoped>

</style>
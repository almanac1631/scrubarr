<script setup lang="ts">
import {computed, onMounted, Ref, ref, watch} from "vue";
import {EntryMapping, GetEntryMappingsFilterEnum, GetEntryMappingsSortByEnum, Retriever} from "../api";
import PreloaderTableEntry from "./entry-mapping-list/PreloaderTableEntry.vue";
import Dropdown, {DropdownOption} from "./common/Dropdown.vue";
import Pagination from "./common/Pagination.vue";
import {getCategoriesFromRetrieverList, sortRetrieverList} from "../utils/retrievers.ts";
import TableRetrieverStateHeader from "./entry-mapping-list/TableRetrieverStateHeader.vue";
import TableRetrieverStateRowEntry from "./entry-mapping-list/TableRetrieverStateRowEntry.vue";
import {getApiClient} from "../utils/api.ts";
import {formatFileSize} from "../utils/fileSize.ts";
import {notify} from "../utils/notificationList.ts";
import EntryMappingDetails from "./entry-mapping-list/EntryMappingDetails.vue";

const contentLoaded = ref(false);

type EntryMappingWrapper = EntryMapping & {
  unfolded: boolean
};

const entryMappingList: Ref<Array<EntryMappingWrapper> | null> = ref(null);
const entryMappingTotalAmount: Ref<number | null> = ref(null);

interface RetrieverWrapper extends Retriever {
  hasMultipleInstances: boolean;
}

const retrieverList: Ref<Array<RetrieverWrapper> | null> = ref(null);

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

const selectedPage: Ref<number> = ref(1);

const apiClient = getApiClient();

async function fetchAndDisplayEntries() {
  if (retrieverList.value?.values === null) {
    throw Error("cannot fetch and display entries if retrievers are not resolved yet");
  }

  if (selectedFilter.value === null || selectedPageSize.value === null || selectedPageSize.value.value === undefined) {
    return;
  }
  try {
    const nameVal = name.value === null ? undefined : name.value;
    const entryMappingResponse = (await apiClient.getEntryMappings(selectedPage.value, +selectedPageSize.value.value, selectedFilter.value.value as GetEntryMappingsFilterEnum, selectedSortBy.value, nameVal)).data;
    entryMappingList.value = [];
    const entries = entryMappingResponse.entries ?? [];
    for (const response of entries) {
      const entryMappingWrapper: EntryMappingWrapper = {
        ...response,
        unfolded: false
      };
      entryMappingList.value.push(entryMappingWrapper);
    }
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

async function prepareEntryMappingDelete(entryMapping: EntryMapping, event: MouseEvent) {
  entrySelectedForDeletion.value = entryMapping;

  if (event.shiftKey) {
    await deleteEntryMapping(entryMapping);
  }
}

async function deleteEntryMapping(entryMapping: EntryMapping) {
  if (entrySelectedForDeletion.value !== entryMapping) {
    throw Error("cannot delete entry mapping if it is not selected for deletion");
  }
  try {
    await apiClient.deleteEntryMapping(entryMapping.id);
  } catch (e) {
    console.error("Error deleting entry mapping.");
    console.error(e);
    notify("Could not delete entry mapping: " + e, "error");
    return;
  } finally {
    entrySelectedForDeletion.value = null;
  }
  notify(`Entry mapping ${entryMapping.name} deleted.`, "success");
  await fetchAndDisplayEntries();
}

const entrySelectedForDeletion: Ref<EntryMapping | null> = ref(null);

onMounted(async () => {
  const retrieverListResp = (await apiClient.getRetrievers()).data.retrievers;
  sortRetrieverList(retrieverListResp);
  const wrappedRetrieverList = new Array<RetrieverWrapper>();
  for (const retriever of retrieverListResp) {
    const hasMultipleInstances = retrieverListResp.filter((filterRetriever) => filterRetriever.id !== retriever.id &&
        filterRetriever.category === retriever.category && filterRetriever.softwareName === retriever.softwareName).length > 0;
    wrappedRetrieverList.push({
      ...retriever,
      hasMultipleInstances: hasMultipleInstances
    });
  }
  retrieverList.value = wrappedRetrieverList;
});

const filterElemList = [
  {displayName: "Filter by", value: undefined},
  {displayName: "No filter", value: undefined},
  {displayName: "Incomplete entries", value: GetEntryMappingsFilterEnum.IncompleteEntries},
  {displayName: "Complete entries", value: GetEntryMappingsFilterEnum.CompleteEntries}
];

const selectedFilter: Ref<DropdownOption | null> = ref(null);

const selectedSortBy: Ref<GetEntryMappingsSortByEnum | undefined> = ref(undefined);

function toggleSortBy(availableOptions: Array<GetEntryMappingsSortByEnum>) {
  if (selectedSortBy.value === undefined) {
    selectedSortBy.value = availableOptions[0];
  } else {
    const currentIndex = availableOptions.indexOf(selectedSortBy.value);
    if (currentIndex === availableOptions.length - 1) {
      selectedSortBy.value = undefined;
    } else {
      selectedSortBy.value = availableOptions[currentIndex + 1];
    }
  }
}

const pageSizeElemList = [
  {displayName: "Page size", value: "25"},
  {displayName: "10", value: "10"},
  {displayName: "25", value: "25"},
  {displayName: "50", value: "50"},
  {displayName: "100", value: "100"}
];

const selectedPageSize: Ref<DropdownOption | null> = ref(null);

watch([selectedFilter, selectedPageSize, selectedPage, selectedSortBy], () => {
  fetchAndDisplayEntries();
});

const name: Ref<string | null> = ref(null);

watch([name], () => {
  const valueChangedTo = name.value;
  setTimeout(() => {
    if (name.value !== valueChangedTo) {
      return;
    }

    if (selectedPage.value !== 1) {
      selectedPage.value = 1; // Will trigger previous watch and fetch
    } else {
      fetchAndDisplayEntries();
    }
  }, 200);
});
</script>

<template>
  <div class="container mx-auto rounded-md bg-white px-8 py-6 shadow">
    <div class="my-2 flex">
      <div>
        <label class="inline-flex items-center cursor-pointer">
          <span class="me-3 text-sm font-medium text-gray-900">Group by retriever category</span>
          <input type="checkbox" value="" class="sr-only peer" v-model="retrieverGroupingEnabled">
          <span
              class="relative w-11 h-6 bg-gray-200 peer-focus:ring-2 rounded-full peer peer-checked:after:translate-x-full rtl:peer-checked:after:-translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:start-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-red-500"></span>
        </label>
      </div>
    </div>
    <div class="my-2 flex justify-between">
      <div class="relative my-2">
        <input type="text" id="search-bar"
               class="bg-gray-100 text-gray-400 font-medium rounded focus:ring-red-500 block pl-9 py-2 pr-2 w-80"
               placeholder="Search" v-model="name">
        <div class="absolute left-0 top-0 flex items-center p-2 text-gray-500">
          <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none"
               stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"
               class="icon icon-tabler icons-tabler-outline icon-tabler-search">
            <path stroke="none" d="M0 0h24v24H0z" fill="none"/>
            <path d="M10 10m-7 0a7 7 0 1 0 14 0a7 7 0 1 0 -14 0"/>
            <path d="M21 21l-6 -6"/>
          </svg>
        </div>
      </div>
      <div class="my-2 flex justify-end gap-2">
        <Dropdown :options="pageSizeElemList" :default-option="pageSizeElemList[0]" v-model="selectedPageSize"/>
        <Dropdown :options="filterElemList" :default-option="filterElemList[0]" v-model="selectedFilter"/>
      </div>
    </div>
    <table class="table-fixed w-full">
      <thead>
      <tr class="text-left border-b-2 border-gray-200" v-if="contentLoaded">
        <th class="py-3 pr-3 font-medium pl-3">
          <button class="flex items-center"
                  v-on:click="toggleSortBy([GetEntryMappingsSortByEnum.NameAsc, GetEntryMappingsSortByEnum.NameDesc])">
            Name
            <svg class="w-4 h-4 ms-1 inline" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" width="24"
                 height="24" fill="none" viewBox="0 0 24 24">
              <path :class="{'text-slate-300': selectedSortBy !== GetEntryMappingsSortByEnum.NameAsc}"
                    stroke="currentColor" stroke-linecap="round"
                    stroke-linejoin="round"
                    stroke-width="3"
                    d="m16 9-4-4-4 4"/>
              <path :class="{'text-slate-300': selectedSortBy !== GetEntryMappingsSortByEnum.NameDesc}"
                    stroke="currentColor" stroke-linecap="round" stroke-linejoin="round"
                    stroke-width="3"
                    d="m8 15 4 4 4-4"/>
            </svg>
          </button>
        </th>
        <th class="w-28 pr-3 font-medium">
          <button class="flex items-center"
                  v-on:click="toggleSortBy([GetEntryMappingsSortByEnum.SizeAsc, GetEntryMappingsSortByEnum.SizeDesc])">
            Size
            <svg class="w-4 h-4 ms-1 inline" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" width="24"
                 height="24" fill="none" viewBox="0 0 24 24">
              <path :class="{'text-slate-300': selectedSortBy !== GetEntryMappingsSortByEnum.SizeAsc}"
                    stroke="currentColor" stroke-linecap="round"
                    stroke-linejoin="round"
                    stroke-width="3"
                    d="m16 9-4-4-4 4"/>
              <path :class="{'text-slate-300': selectedSortBy !== GetEntryMappingsSortByEnum.SizeDesc}"
                    stroke="currentColor" stroke-linecap="round" stroke-linejoin="round"
                    stroke-width="3"
                    d="m8 15 4 4 4-4"/>
            </svg>
          </button>
        </th>
        <th class="w-52 pr-3 font-medium">
          <button class="flex items-center"
                  v-on:click="toggleSortBy([GetEntryMappingsSortByEnum.DateAddedAsc, GetEntryMappingsSortByEnum.DateAddedDesc])">
            Added
            <svg class="w-4 h-4 ms-1 inline" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" width="24"
                 height="24" fill="none" viewBox="0 0 24 24">
              <path :class="{'text-slate-300': selectedSortBy !== GetEntryMappingsSortByEnum.DateAddedAsc}"
                    stroke="currentColor" stroke-linecap="round"
                    stroke-linejoin="round"
                    stroke-width="3"
                    d="m16 9-4-4-4 4"/>
              <path :class="{'text-slate-300': selectedSortBy !== GetEntryMappingsSortByEnum.DateAddedDesc}"
                    stroke="currentColor" stroke-linecap="round" stroke-linejoin="round"
                    stroke-width="3"
                    d="m8 15 4 4 4-4"/>
            </svg>
          </button>
        </th>
        <th v-if="retrieverGroupingEnabled" v-for="retrieverCategory in retrieverCategoryList"
            class="w-[120px] p-3 font-medium text-center truncate">
          <div class="h-6 flex justify-center">
            <div class="relative">
              {{ retrieverCategory.displayName }}
            </div>
          </div>
        </th>
        <th v-else v-for="retriever in retrieverList" class="w-[80px] p-3 font-medium text-center">
          <TableRetrieverStateHeader
              :name="retriever.hasMultipleInstances ? retriever.name : null"
              :hover-text="retriever.softwareName"
              :logo-filename="`${retriever.softwareName}-128x128.png`"
              :logo-alt-text="`The logo of the ${retriever.softwareName} software project.`"
          />
        </th>
        <th class="w-[80px] p-3"></th>
      </tr>
      <PreloaderTableEntry class="border-b-2 border-gray-200" v-else/>
      </thead>
      <tbody>
      <template v-if="contentLoaded" v-for="entryMapping in entryMappingList">
        <tr class="hover:bg-stone-100 border-t border-gray-200">
          <td class="py-3 pr-3 pl-3 font-medium truncate cursor-pointer" :title="entryMapping.name"
              v-on:click="entryMapping.unfolded = !entryMapping.unfolded">
            <div class="flex overflow-hidden">
              <div>
                <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none"
                     stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"
                     class="icon icon-tabler icons-tabler-outline icon-tabler-chevron-right transform transition-transform"
                     :class="[entryMapping.unfolded ? 'rotate-90' : '']">
                  <path stroke="none" d="M0 0h24v24H0z" fill="none"/>
                  <path d="M9 6l6 6l-6 6"/>
                </svg>
              </div>
              <div>
                {{ entryMapping.name }}
              </div>
            </div>
          </td>

          <td class="py-3 pr-3 font-medium truncate" :title="formatFileSize(entryMapping.size)">
            {{ formatFileSize(entryMapping.size) }}
          </td>

          <td class="py-3 pr-3 font-medium truncate" :title="entryMapping.dateAdded">
            {{ new Date(entryMapping.dateAdded).toISOString().replace(".000", "") }}
          </td>

          <TableRetrieverStateRowEntry
              v-if="retrieverGroupingEnabled" v-for="retrieverCategory in retrieverCategoryList"
              :present="isEntryPresentInRetrieverCategory(entryMapping, retrieverCategory.name)"
          />
          <TableRetrieverStateRowEntry
              v-else v-for="retriever in retrieverList"
              :present="isEntryPresentInRetriever(entryMapping, retriever.id)"
          />

          <td class="p-3 flex justify-center">
            <button class="text-gray-500" v-if="entrySelectedForDeletion !== entryMapping"
                    v-on:click="prepareEntryMappingDelete(entryMapping, $event)">
              <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none"
                   stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"
                   class="icon icon-tabler icons-tabler-outline icon-tabler-trash">
                <path stroke="none" d="M0 0h24v24H0z" fill="none"/>
                <path d="M4 7l16 0"/>
                <path d="M10 11l0 6"/>
                <path d="M14 11l0 6"/>
                <path d="M5 7l1 12a2 2 0 0 0 2 2h8a2 2 0 0 0 2 -2l1 -12"/>
                <path d="M9 7v-3a1 1 0 0 1 1 -1h4a1 1 0 0 1 1 1v3"/>
              </svg>
            </button>
            <button class="text-green-500 hover:bg-green-100 rounded-md"
                    v-if="entrySelectedForDeletion === entryMapping"
                    v-on:click="deleteEntryMapping(entryMapping)">
              <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none"
                   stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"
                   class="icon icon-tabler icons-tabler-outline icon-tabler-check">
                <path stroke="none" d="M0 0h24v24H0z" fill="none"/>
                <path d="M5 12l5 5l10 -10"/>
              </svg>
            </button>
            <button class="text-red-500 hover:bg-red-100 rounded-md" v-if="entrySelectedForDeletion === entryMapping"
                    v-on:click="entrySelectedForDeletion = null">
              <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none"
                   stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"
                   class="icon icon-tabler icons-tabler-outline icon-tabler-x">
                <path stroke="none" d="M0 0h24v24H0z" fill="none"/>
                <path d="M18 6l-12 12"/>
                <path d="M6 6l12 12"/>
              </svg>
            </button>
          </td>
        </tr>
        <EntryMappingDetails v-if="entryMapping.unfolded" :entry-id="entryMapping.id"></EntryMappingDetails>
      </template>
      <PreloaderTableEntry v-for="_ in 10" v-else/>
      </tbody>
    </table>
    <Pagination
        v-if="contentLoaded && entryMappingTotalAmount !== null &&selectedPageSize !== null &&  selectedPageSize.value !== undefined"
        :page-size="+selectedPageSize.value"
        :selected-page="1"
        :total-amount-of-items="entryMappingTotalAmount" v-model="selectedPage"/>
  </div>
</template>

<style scoped>

</style>
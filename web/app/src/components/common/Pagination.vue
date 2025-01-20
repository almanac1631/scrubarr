<script setup lang="ts">
import {computed} from "vue";
import {getPageList} from "../../utils/pagination.ts";

const model = defineModel<number>({default: null});

const props = defineProps<{
  totalAmountOfItems: number
  pageSize: number
}>()

const minElem = computed(() => {
  return (model.value - 1) * props.pageSize + 1;
});

const maxElem = computed(() => {
  return Math.min(model.value * props.pageSize, props.totalAmountOfItems);
});

function updateSelectedPageNumber(pageNumber: number) {
  model.value = pageNumber;
}

function updateSelectedPageNumberRelative(delta: number) {
  if (model.value <= 1 && delta < 0) {
    return;
  } else {
    const previousPageNumber = pageList.value[pageList.value.length - 1].pageNumber
    if (previousPageNumber === null) {
      return;
    }
    if (model.value >= previousPageNumber && delta > 0) {
      return;
    }
  }
  model.value += delta;
}

const pageList = computed(() => getPageList(props.totalAmountOfItems, props.pageSize, model.value, 6));
</script>

<template>
  <div class="flex items-center justify-between border-t border-gray-200 bg-white py-6 px-3">
    <div class="flex flex-1 justify-between sm:hidden">
      <a href="#"
         class="relative inline-flex items-center rounded-md border border-gray-300 bg-white px-4 py-2 text-sm font-medium text-gray-700 hover:bg-gray-50">Previous</a>
      <a href="#"
         class="relative ml-3 inline-flex items-center rounded-md border border-gray-300 bg-white px-4 py-2 text-sm font-medium text-gray-700 hover:bg-gray-50">Next</a>
    </div>
    <div class="hidden sm:flex sm:flex-1 sm:items-center sm:justify-between">
      <div>
        <p class="text-sm text-gray-700">
          Showing
          <span class="font-medium">{{ minElem }}</span>
          to
          <span class="font-medium">{{ maxElem }}</span>
          of
          <span class="font-medium">{{ props.totalAmountOfItems }}</span>
          results
        </p>
      </div>
      <div>
        <nav class="isolate inline-flex -space-x-px rounded-md shadow-sm" aria-label="Pagination">
          <a href="#" v-on:click.prevent="updateSelectedPageNumberRelative(-1)"
             class="relative inline-flex items-center rounded-l-md px-2 py-2 text-gray-400 ring-1 ring-inset ring-gray-300 hover:bg-gray-50 focus:z-20 focus:outline-offset-0">
            <span class="sr-only">Previous</span>
            <svg class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor" aria-hidden="true" data-slot="icon">
              <path fill-rule="evenodd"
                    d="M11.78 5.22a.75.75 0 0 1 0 1.06L8.06 10l3.72 3.72a.75.75 0 1 1-1.06 1.06l-4.25-4.25a.75.75 0 0 1 0-1.06l4.25-4.25a.75.75 0 0 1 1.06 0Z"
                    clip-rule="evenodd"/>
            </svg>
          </a>

          <span v-for="page in pageList">
            <span v-if="page.pageNumber === null"
                  class="relative inline-flex items-center px-4 py-2 text-sm font-semibold text-gray-700 ring-1 ring-inset ring-gray-300 focus:outline-offset-0">
              ...
            </span>
            <a href="#" v-else-if="page.isSelected" v-on:click.prevent="updateSelectedPageNumber(page.pageNumber)"
               aria-current="page"
               class="relative z-10 inline-flex items-center bg-red-500 px-4 py-2 text-sm font-semibold text-white focus:z-20 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600">
              {{ page.pageNumber }}
            </a>
            <a href="#" v-else v-on:click.prevent="updateSelectedPageNumber(page.pageNumber)"
               class="relative inline-flex items-center px-4 py-2 text-sm font-semibold text-gray-900 ring-1 ring-inset ring-gray-300 hover:bg-gray-50 focus:z-20 focus:outline-offset-0">
              {{ page.pageNumber }}
            </a>
          </span>
          <a href="#" v-on:click.prevent="updateSelectedPageNumberRelative(1)"
             class="relative inline-flex items-center rounded-r-md px-2 py-2 text-gray-400 ring-1 ring-inset ring-gray-300 hover:bg-gray-50 focus:z-20 focus:outline-offset-0">
            <span class="sr-only">Next</span>
            <svg class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor" aria-hidden="true" data-slot="icon">
              <path fill-rule="evenodd"
                    d="M8.22 5.22a.75.75 0 0 1 1.06 0l4.25 4.25a.75.75 0 0 1 0 1.06l-4.25 4.25a.75.75 0 0 1-1.06-1.06L11.94 10 8.22 6.28a.75.75 0 0 1 0-1.06Z"
                    clip-rule="evenodd"/>
            </svg>
          </a>
        </nav>
      </div>
    </div>
  </div>
</template>

<style scoped>

</style>
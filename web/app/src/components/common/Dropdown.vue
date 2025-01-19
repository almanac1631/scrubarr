<script setup lang="ts">
import {onMounted, PropType} from "vue";

const model = defineModel();

export interface DropdownOption {
  displayName: string
  value: string | undefined
}

const props = defineProps({
  options: Object as PropType<Array<DropdownOption>>,
  defaultOption: Object as PropType<DropdownOption>,
});

onMounted(async () => {
  model.value = props.defaultOption;
});
</script>

<template>
  <div class="relative text-gray-400">
    <select name="filterSelection" class="p-2 bg-gray-100 rounded font-medium appearance-none w-[200px]"
            v-model="model">
      <option v-for="elem in props.options" :value="elem" :selected="elem === defaultOption"
              :disabled="elem === defaultOption" :hidden="elem === defaultOption">
        {{ elem.displayName }}
      </option>
    </select>
    <div class="h-full absolute top-0 right-0 flex items-center pointer-events-none">
      <svg class="h-8 w-8" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"
           stroke-linecap="round" stroke-linejoin="round">
        <polyline points="6 9 12 15 18 9"/>
      </svg>
    </div>
  </div>
</template>

<style scoped>

</style>
<script setup lang="ts">

import {computed} from "vue";

const props = defineProps({
  present: Boolean
});

const presentText = computed(() => props.present ? "✅" : "❌");
</script>

<template>
  <td class="p-3 text-center">
    {{ presentText }}
  </td>
</template>

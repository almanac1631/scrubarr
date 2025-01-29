<script setup lang="ts">

const props = defineProps({
  name: String,
  hoverText: String,
  logoFilename: String,
  logoAltText: String,
});

const logoFilepath = `/retriever/${props.logoFilename}`
</script>

<template>
  <div class="h-6 flex justify-center">
    <div class="relative" :title="props.hoverText">
      <img :src="logoFilepath"
           :alt="logoAltText"
           class="max-h-full w-full h-full">
      <div v-if="name !== undefined" class="bg-white absolute rounded -bottom-2 -right-5 z-10">
        <div
            class="bg-red-100 leading-5 h-5 uppercase text-[0.5rem] font-bold px-0.5 border-red-500 rounded">
          {{ props.name }}
        </div>
      </div>
    </div>
  </div>
</template>

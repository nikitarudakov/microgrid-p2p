<script setup lang="ts">
defineProps<{
  max_capacity: number
}>()

const inputValue = defineModel()
</script>

<template>
  <section>
    <label class="block font-[500]" for="energy">Select Capacity</label>
    <div class="flex items-center gap-x-[1em] my-4">
      <input v-model="inputValue" type="range" id="energy" name="energy" min="1" :max="max_capacity" step="0.01"/>
      <span class="font-medium text-gray-800 w-[50px] shrink-2">{{ inputValue }} <br> kWh</span>
    </div>
    <div class="flex justify-between">
      <p class="text-gray-600">Min: 1 kWh</p>
      <p class="text-gray-600">Max: {{ max_capacity.toFixed(0) }} kWh</p>
    </div>
  </section>
</template>


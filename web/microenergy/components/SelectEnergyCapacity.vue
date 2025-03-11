<script setup lang="ts">
defineProps<{
  max_capacity: number
}>()

const inputValue = defineModel()
</script>

<template>
  <section>
    <label class="block" for="energy">Select Capacity</label>
    <div class="flex items-center gap-x-[1em]">
      <input v-model="inputValue" type="range" id="energy" name="energy" min="1" :max="max_capacity" step="0.01"/>
      <span>{{ inputValue }} <br> kWh</span>
    </div>
    <div class="flex justify-between">
      <p>Min: 1 kWh</p>
      <p>Max: {{ max_capacity }} kWh</p>
    </div>
  </section>
</template>


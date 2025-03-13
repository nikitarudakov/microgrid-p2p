<script setup lang="ts">
import { computed } from 'vue'
import WindIcon from "~/assets/img/windIcon.svg"
import SolarPanelIcon from "~/assets/img/solarPanelIcon.svg"

const props = defineProps<{
  name: string
}>()

// Define icon based on energy resource type
const energyTypeIcon = computed(() => {
  if (props.name == "Solar") {
    return SolarPanelIcon
  } else {
    return WindIcon
  }
})
</script>

<template>
  <img class="icon" :src=energyTypeIcon alt="energy_type_icon">
  <p>{{ name }}</p>
</template>

<style scoped>
</style>

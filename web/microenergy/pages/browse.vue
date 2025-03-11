<template>
  <div v-if="error">Here is an error: {{ error.statusCode }}</div>
  <ul class="whitespace-nowrap">
    <li v-for="resource in data?.energyResources" :key="resource.name">
      <EnergyListing :energy-resource=resource />
    </li>
  </ul>
</template>

<style scoped>
ul {
  --n: 6; /* The maximum number of columns */
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(max(200px, 100%/var(--n)), 1fr));
  gap: 2rem;
}
</style>

<script setup lang="ts">
const { data, error } = await useAsyncGql({ operation: 'EnergyResources' });
</script>

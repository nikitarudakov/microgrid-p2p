<script setup lang="ts">
import { useQuery } from '@vue/apollo-composable'
import { computed, type ComputedRef } from 'vue'
import { ENERGY_RESOURCES } from '@/application/queries.ts'
import type { EnergyResource } from '@/interfaces/interfaces.ts'
import ListingElement from '@/components/ListingElement.vue'

const { result, loading, error } = useQuery(ENERGY_RESOURCES)

const listing: ComputedRef<EnergyResource[]> = computed(() => {
  if (!result.value) return []
  return result.value.energyResources.map(resource => ({
    id: resource.id as string,
    name: resource.name as string,
    producer: {
      firstName: resource.producer.first_name,
      lastName: resource.producer.last_name,
    },
    capacity: resource.capacity as number,
    price: resource.price as number,
  }))
})
</script>


<template>
  <div v-if="loading">Loading...</div>
  <div v-else-if="error">Error: {{ error.message }}</div>
  <ul class="whitespace-nowrap">
    <li v-for="resource in listing.values()" :key="resource.name">
      <ListingElement :energy-resource=resource />
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

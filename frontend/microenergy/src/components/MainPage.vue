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
  <ul class="flex gap-10 whitespace-nowrap">
    <li v-for="resource in listing.values()" :key="resource.name">
      <ListingElement :energy-resource=resource />
    </li>
  </ul>
</template>

<script setup lang="ts">
import { useQuery } from '@vue/apollo-composable'
import { computed, type ComputedRef } from 'vue'
import { ENERGY_RESOURCES } from '@/application/queries.ts'
import type { EnergyResource, User } from '@/interfaces/interfaces.ts'

const { result, loading, error } = useQuery(ENERGY_RESOURCES)

const listing: ComputedRef<EnergyResource[]> = computed(() => {
  if (!result.value) return []
  return result.value.energyResources.map(resource => ({
    name: resource.name as string,
    producer: resource.producer as User,
    capacity: resource.capacity as number,
    price: resource.price as number,
  }))
})
</script>


<template>
  <div v-if="loading">Loading...</div>
  <div v-else-if="error">Error: {{ error.message }}</div>
  <ul class="flex" v-for="resource in listing.values()" :key="resource.name">
    <li>{{ resource.name }} </li>
  </ul>
</template>

<script setup lang="ts">
import type { EnergyResource } from "@/interfaces/interfaces.ts"
import UserElement from '@/components/UserElement.vue'
import EnergyGrayIcon from '@/assets/energy-gray-icon.svg'
import SolarIcon from '@/assets/solar-icon.svg'
import WindIcon from '@/assets/wind-icon.svg'
import { computed, ref } from 'vue'
import PurchasePopUp from '@/components/PurchasePopUp.vue'

const props = defineProps<{
  energyResource: EnergyResource
}>()

const username = ref(
  `${props.energyResource.producer.firstName} ${props.energyResource.producer.lastName}`
)

// Used for showing pop-up
const showPurchasePopUp = ref(false)

// Define icon based on energy resource type
const energyTypeIcon = computed(() => {
  if (props.energyResource.name == "Solar Panel") {
    return SolarIcon
  } else {
    return WindIcon
  }
})
</script>


<template>
  <div class="px-10 py-8 drop-shadow-xl bg-white rounded-xl">
    <UserElement :user=energyResource.producer />

    <ul class="mt-3">
      <li class="flex">
        <img :src=energyTypeIcon alt="energy_type_icon">
        <p>{{ energyResource.name }}</p>
      </li>
      <li class="flex">
        <img :src=EnergyGrayIcon alt="energy_icon">
        <p>{{ `${energyResource.capacity.toFixed(0) } kWh Available` }} </p>
      </li>
    </ul>

    <p class="mt-3 font-bold text-emerald-600 text-[1.3rem]">{{ `\$${energyResource.price}/kWh` }}</p>

    <button @click="showPurchasePopUp = !showPurchasePopUp" class="purchase-btn w-[175px] mt-3 rounded-2xl">
      <img src="../assets/energy-icon.svg" alt="energy icon">
      <span>Purchase Now</span>
    </button>
  </div>

  <div v-if="showPurchasePopUp" class="w-full h-full t-0 absolute flex justify-center items-center">
    <PurchasePopUp
      :price=energyResource.price
      :producer-name=username
      :resource-name=energyResource.name
      :capacity=energyResource.capacity
    />
  </div>
</template>

<style scoped>
li {
  align-items: center;
  gap: 1em;
  max-height: 60px;

  img {
    max-width: 20px;
    width: 100%;
    max-height: 16px;
    object-fit: contain;
  }
}
</style>

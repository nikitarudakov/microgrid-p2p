<script setup lang="ts">
import { ref } from 'vue'
import UserIcon from "~/components/UserIcon.vue";
import type {EnergyResource} from "~/stores/energyResources";

const props = defineProps<{
  energyResource: EnergyResource
}>()

// Used for showing pop-up
const showPurchasePopUp = ref(false)
</script>

<template>
  <div class="flex flex-wrap gap-3 px-10 py-8 drop-shadow-xl bg-white rounded-xl">
    <UserIcon :user=energyResource.producer />

    <ul class="mt-3">
      <li class="flex">
        <EnergyTypeIcon :name=energyResource.name />
      </li>
      <li class="flex">
        <img src="~/assets/img/grayLightningIcon.svg" alt="energy_icon">
        <p>{{ `${ energyResource.capacity.toFixed(0) } kWh Available` }} </p>
      </li>
    </ul>

    <p class="mt-3 font-bold text-emerald-600 text-[1.3rem]">{{ `\$${energyResource.price}/kWh` }}</p>

    <button @click="showPurchasePopUp = true" class="w-full purchase-btn mt-3 rounded-2xl">
      <img src="~/assets/img/lightningIcon.svg" alt="energy icon">
      <span>Purchase Now</span>
    </button>
  </div>

  <div v-if="showPurchasePopUp" class="fixed w-full h-full top-0 left-0 z-99">
    <EnergyPurchasePopUp
        @close = "showPurchasePopUp = false"
        :id=energyResource.id
        :name=energyResource.name
        :price=energyResource.price
        :producer=energyResource.producer
        :capacity=energyResource.capacity
    />
  </div>
</template>

<style scoped>
li {
  align-items: center;
  gap: 1em;
  max-height: 60px;
}
</style>

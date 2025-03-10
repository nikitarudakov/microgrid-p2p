<script setup lang="ts">
import { ref } from "vue";
import type { User } from '@/interfaces/interfaces.ts'
import UserElement from '@/components/UserElement.vue'

// Events this component will emit
defineEmits(['close'])

defineProps<{
  resourceName: string,
  producer: User,
  capacity: number,
  price: number,
}>()

// Updated when user uses range input to select desired energy capacity for purchase
const selectedEnergyCapacity = ref(0);

const calculateTotalPrice = (): number => {
  const totalPrice = Number(selectedEnergyCapacity.value) * 1.1
  return Math.round((totalPrice + Number.EPSILON) * 100) / 100
}

async function submitPurchaseRequest() {
    console.log("request was sent!")
}
</script>

<template>
 <div class="h-full inset-0 flex items-center justify-center bg-gray-900/40">
  <div class="bg-gradient-to-b from-white to-emerald-50 rounded-2xl px-15 py-10 border-2 border-white">

    <div class="flex justify-between">
      <h1 class="text-2xl font-bold">{{ resourceName }}</h1>
      <img @click="$emit('close')" class="cursor-pointer" src="../assets/close.png" alt="close" />
    </div>

    <div class="mt-5 mb-10">
      <p class="text-base mb-4">Purchase sustainable energy from local <br> green producers: </p>
      <UserElement :user=producer />
    </div>

    <label for="energy" class="flex flex-col my-auto">
      <span class="mb-2 font-medium">Select Energy Amount:</span> <!-- Acts as a header -->

      <div class="flex items-center gap-2"> <!-- Input and value in a row -->
        <input
          type="range"
          v-model="selectedEnergyCapacity"
          id="energy" name="energy"
          min="1" :max="capacity" step="0.01"
          class="flex-grow"
        />
        <p class="text-xl leading-6 font-medium">{{ selectedEnergyCapacity }} <br> kWh</p>
      </div>
    </label>

    <div class="text-emerald-600 pt-10 mx-10">
      <div class="flex">
        <img src="../assets/plant-icon.svg" alt="plant">
        <p class="ml-2">Environmental Impact</p>
      </div>
      <p>Reduces CO₂ emissions by approximately 40kg</p>
    </div>

    <div class="mt-10">
      <button class="flex w-full justify-center cursor-pointer  bg-emerald-600 text-white p-4 rounded-xl drop-shadow-xl" type="submit" @click="submitPurchaseRequest">
        <img class="mr-2" src="../assets/energy-icon.svg" alt="energy icon">
        Purchase Energy
      </button>
      <p class="text-center mt-3 text-gray-600">{{ 'Total Price: ' + '$' + calculateTotalPrice() + ' USD' }}</p>
    </div>

  </div>
</div>
</template>

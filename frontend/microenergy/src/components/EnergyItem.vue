<script setup lang="ts">
import { ref, computed } from "vue";
import { api } from "@/application/api.ts";
import defaultPhoto from "@/assets/profile.png";

const props = defineProps<{
  sellerUsername: string,
  photoUrl?: string,
  availableEnergyCapacity: number,
  kWhPrice: number,
}>()


const withDefaultPhotoURL = computed(() => {
  return props.photoUrl || defaultPhoto
})

// Updated when user uses range input to select desired energy capacity for purchase
const selectedEnergyCapacity = ref(0);

const calculateTotalPrice = (): number => {
  const totalPrice = Number(selectedEnergyCapacity.value) * props.kWhPrice
  return Math.round((totalPrice + Number.EPSILON) * 100) / 100
}

async function submitPurchaseRequest() {
    try {
      const response = await api.post( "/", {
        seller: props.sellerUsername,
        capacity: Number(selectedEnergyCapacity.value),
      }, {
        headers: {
          'Content-Type': 'application/json',
        }
      });

      if (response.status != 200) {
        throw new Error(`HTTP error! Status: ${response.status}`);
      }

      console.log("Response:", response.data);
    } catch (error) {
      console.error("Error sending request:", error);
    }
}
</script>

<template>
<div class="fixed inset-0 flex items-center justify-center bg-gray-900/70">
  <div class="flex flex-col min-h-fit bg-gradient-to-b from-green-100 to-green-100 text-black rounded-2xl px-15 py-10">
    <h1 class="text-2xl font-bold">Eco Energy Trading</h1>

    <div class="mt-5 mb-10">
      <p class="text-base mb-4">Purchase sustainable energy from local <br> green producers: </p>
      <div class="flex items-center">
        <img class="max-w-10 max-h-auto" :src="withDefaultPhotoURL" alt="photo"/>
        <p class="font-bold ml-2">{{ sellerUsername }}</p>
      </div>
    </div>

    <label for="energy" class="flex flex-col my-auto">
      <span class="mb-2 font-medium">Select Energy Amount:</span> <!-- Acts as a header -->

      <div class="flex items-center gap-2"> <!-- Input and value in a row -->
        <input
          type="range"
          v-model="selectedEnergyCapacity"
          id="energy" name="energy"
          min="1" :max="availableEnergyCapacity"
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
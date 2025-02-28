<script setup lang="ts">
import {ref} from "vue";

const props = defineProps<{
  sellerUsername: string,
  photoUrl?: string,
  availableEnergyCapacity: number,
  kWhPrice: number,
}>()

const selectedEnergyCapacity = ref<number | string>(0);
async function submitPurchaseRequest() {
    try {
      const response = await fetch("http://localhost:8080/", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          seller: props.sellerUsername,
          capacity: Number(selectedEnergyCapacity.value),
        }),
      });

      if (!response.ok) {
        throw new Error(`HTTP error! Status: ${response.status}`);
      }

      const data = await response.json();
      console.log("Response:", data);
    } catch (error) {
      console.error("Error sending request:", error);
    }
}
</script>

<template>
<div>
  <h1>Eco Energy Trading</h1>

  <p>Purchase sustainable energy from local green producers</p>
  <img v-if="photoUrl" :src="photoUrl" alt="photo"/>
  <p>{{ sellerUsername }}</p>

  <label for="energy">
    Select Energy Amount
    <input
      type="range"
      v-model="selectedEnergyCapacity"
      id="energy" name="energy"
      min="1" :max="availableEnergyCapacity"
    />
  </label>
  <p>{{selectedEnergyCapacity}} kWh</p>
  <button type="submit" @click="submitPurchaseRequest">Purchase Energy</button>

</div>
</template>

<style scoped>

</style>

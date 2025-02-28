<script setup lang="ts">
import {ref} from "vue";
import axios from "axios";

const props = defineProps<{
  sellerUsername: string,
  photoUrl?: string,
  availableEnergyCapacity: number,
  kWhPrice: number,
}>()

const selectedEnergyCapacity = ref<number>(0);
const submitPurchaseRequest = () => {
  axios.post(`http://localhost:8080/`, {
    "seller": props.sellerUsername,
    "capacity": selectedEnergyCapacity.value,
  })
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

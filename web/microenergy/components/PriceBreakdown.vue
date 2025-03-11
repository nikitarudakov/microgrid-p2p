<script setup lang="ts">
const props = defineProps<{
  price: number
}>()

const inputEnergyCapacityValue = defineModel()

const totalPrice = computed(() => {
  let totalPrice = Number(inputEnergyCapacityValue.value) * props.price
  return Math.round((totalPrice + Number.EPSILON) * 100) / 100
})

</script>

<template>
  <section>
    <p>Price Breakdown</p>
    <ul>
      <li>
        <span>Price per kWh: </span>
        <span>{{ price }}</span>
      </li>
      <li>
        <span>Transaction Fee: </span>
        <span>$2.50</span>
      </li>
    </ul>
    <p>
      <span>Total Price: </span>
      <span>${{ totalPrice }}</span>
    </p>
  </section>
</template>

<style scoped>
</style>

<script setup lang="ts">
const props = defineProps<{
  price: number
  capacity: number
}>()

const totalPrice = computed(() => {
  let totalPrice = Number(props.capacity) * props.price
  return Math.round((totalPrice + Number.EPSILON) * 100) / 100
})

</script>

<template>
  <section>
    <p class="font-[500]">Price Breakdown</p>
    <ul class="my-2 border-b-1 border-gray-300 py-2">
      <li class="flex justify-between">
        <span class="text-gray-600">Price per kWh: </span>
        <span>${{ price }}</span>
      </li>
      <li class="flex justify-between">
        <span class="text-gray-600">Transaction Fee: </span>
        <span>$2.50</span>
      </li>
    </ul>
    <p class="flex justify-between font-medium">
      <span>Total Price: </span>
      <span>${{ totalPrice }}</span>
    </p>
  </section>
</template>

<style scoped>
</style>

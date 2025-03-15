<script setup lang="ts">
import {type EnergyResource, useEnergyResourcesStore} from "~/stores/energyResources";
import EnvironmentImpactIcon from "~/components/SelectEnergyCapacity.vue";
import SelectEnergyCapacity from "~/components/SelectEnergyCapacity.vue";
import PriceBreakdown from "~/components/PriceBreakdown.vue";

// Energy Resource prop
const props = defineProps<EnergyResource>()

// Close event controls closing this pop up
const emit = defineEmits(['close'])

// Handle energy capacity selection and pricing
const inputEnergyCapacityValue = useState<number>(() => { return 0; });

// Handling purchase request
const store = useEnergyResourcesStore()

async function submitPurchaseRequest() {
  try {
    const GqlInstance = useGql();
    const { purchaseEnergy } = await GqlInstance('PurchaseEnergy', {
      id: props.id,
      capacity: inputEnergyCapacityValue.value,
    });

    if (purchaseEnergy) {
      store.updateEnergyResource(purchaseEnergy);
      emit('close');
    }
  } catch (error) {
    console.error('Failed to purchase energy:', error);
  }
}
</script>

<template>
  <div class="h-full flex items-center justify-center bg-gray-900/40">
    <div id="purchase-info" class="bg-emerald-50 rounded-2xl p-10">
      <div class="col-span-2">
        <div class="flex justify-between items-center">
          <h1 class="font-bold text-2xl/10">{{ name }}</h1>
          <img @click="$emit('close')" id="purchase-info-close"  src="~/assets/img/closeIcon.png" alt="close" />
        </div>

        <p>Purchase clean energy directly from local producers</p>
      </div>

      <div id="purchase-info-user" class="box grid col-span-2">
        <div class="row-span-2">
          <UserIcon :user=producer />
        </div>
        <div class="flex col-start-2 gap-2 ml-2">
          <img class="icon" src="~/assets/img/starIcon.svg" alt="starIcon" />
          <div>4.9 <span class="text-gray-600">(234 trades)</span></div>
        </div>
        <div class="flex gap-2 items-center col-start-2 ml-2">
          <EnergyTypeIcon :name=name />
        </div>
        <p class="mt-2 text-gray-600">Available Energy: <span class="font-medium">{{ capacity.toFixed(0) }} kWh</span></p>
      </div>

      <SelectEnergyCapacity class="box" v-model="inputEnergyCapacityValue" :max_capacity=capacity />
      <PriceBreakdown class="box" :capacity=inputEnergyCapacityValue :price=price />
      <OrderInformation class="box col-span-2"/>

      <button
          :disabled="inputEnergyCapacityValue <= 0 || inputEnergyCapacityValue > props.capacity"
          @click="submitPurchaseRequest"
          class="box col-span-2 bg-emerald-600 text-white drop-shadow-xl cursor-pointer disabled:bg-gray-400">
        Buy Now
      </button>
    </div>
  </div>
</template>

<style scoped>
#purchase-info {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1.5rem;
}

#purchase-info-user {
  grid-template-columns: auto 1fr;
  grid-template-rows: auto;
}

#purchase-info-close {
  max-width: 20px;
  width: 100%;
  height: 20px;
  cursor: pointer;
}
</style>

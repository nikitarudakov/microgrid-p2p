<script setup lang="ts">
import type {User} from "~/stores/energyResources";

defineProps<{
  user: User
}>()
</script>

<template>
  <div class="flex items-center gap-3">
    <img class="max-w-[40px] w-full max-h-[40px] object-contain" src="~/assets/img/userLogo.svg" alt="user_photo">
    <div>
      <p class="font-bold">{{ `${user.first_name} ${user.last_name}` }}</p>
      <div class="flex gap-1">
        <img src="~/assets/img/verified.svg" alt="verified">
        <p class="text-emerald-800">Verified Producer</p>
      </div>
    </div>
  </div>
</template>

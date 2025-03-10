<script setup lang="ts">
import type { User } from '@/interfaces/interfaces';
defineProps<{
  user: User
}>()
</script>

<template>
  <div class="flex items-center gap-3">
    <img class="max-w-[40px] w-full max-h-[40px] object-contain" src="@/assets/user.svg" alt="user_photo">
    <div>
      <p class="font-bold">{{ `${user.firstName} ${user.lastName}` }}</p>
      <div class="flex gap-1">
        <img src="@/assets/verified-icon.svg" alt="verified">
        <p class="text-emerald-800">Verified</p>
      </div>
    </div>
  </div>
</template>


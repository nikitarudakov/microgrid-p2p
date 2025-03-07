<script setup lang="ts">
import {provide} from 'vue'
import MainPage from "@/components/MainPage.vue";
import { apolloClient } from '@/application/apollo.ts'
import { DefaultApolloClient } from "@vue/apollo-composable"
provide(DefaultApolloClient, apolloClient)
</script>

<template>
  <MainPage/>
</template>

<style scoped>
</style>

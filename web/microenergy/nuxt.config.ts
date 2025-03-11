import tailwindcss from "@tailwindcss/vite";

// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  modules: ['nuxt-graphql-client'],
  compatibilityDate: '2024-11-01',
  devtools: { enabled: true },
  css: ['~/assets/css/main.css'],
  vite: {
    plugins: [
      tailwindcss(),
    ],
  },
  runtimeConfig: {
    public: {
      GQL_HOST: 'http://localhost:8081/query' // overwritten by process.env.GQL_HOST
    }
  },
  'graphql-client': {
    clients: {
      default: {
        host: 'http://localhost:8081/query',
        schema: '../../graph/schema.graphqls',
      }
    },
  }
})

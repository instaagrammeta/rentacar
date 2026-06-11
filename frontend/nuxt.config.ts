// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  compatibilityDate: '2025-01-01',
  devtools: { enabled: true },

  // Single Page Application: all data is fetched client-side from the Go API,
  // which keeps deployment simple (no SSR proxy to the backend required).
  ssr: false,

  modules: ['@nuxtjs/tailwindcss', '@pinia/nuxt'],

  css: ['~/assets/css/main.css'],

  runtimeConfig: {
    public: {
      // Overridden at runtime via NUXT_PUBLIC_API_BASE.
      apiBase: 'http://localhost:5000/api',
    },
  },

  app: {
    head: {
      title: 'Wheelzie — Аренда автомобилей',
      htmlAttrs: { lang: 'ru' },
      meta: [
        { charset: 'utf-8' },
        { name: 'viewport', content: 'width=device-width, initial-scale=1' },
        { name: 'description', content: 'Wheelzie — система управления прокатом автомобилей' },
      ],
      link: [{ rel: 'icon', type: 'image/svg+xml', href: '/favicon.svg' }],
    },
  },

  typescript: {
    strict: true,
  },
})

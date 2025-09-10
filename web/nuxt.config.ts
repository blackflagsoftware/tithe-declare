// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  devtools: { enabled: true },
  postcss: {
    plugins: {
      tailwindcss: {},
      autoprefixer: {},
    },
  },
  modules: ["@nuxtjs/tailwindcss", "nuxt-security"],
  css: ["~/assets/css/main.css"],
  runtimeConfig: {
    public: {
      apiURL: process.env.AUTH_SERVER_API_URL,
      appName: process.env.APP_NAME,
      adminPwd: process.env.ADMIN_PWD
    }
  },
  security: {
    headers: {
      strictTransportSecurity: false,
      crossOriginResourcePolicy: false,
      crossOriginOpenerPolicy: false,
      crossOriginEmbedderPolicy: false,
    }
  }
})
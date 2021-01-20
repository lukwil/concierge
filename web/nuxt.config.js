import colors from 'vuetify/es5/util/colors'

export default {
  /*
   ** Nuxt rendering mode
   ** See https://nuxtjs.org/api/configuration-mode
   */
  mode: 'spa',
  /*
   ** Nuxt target
   ** See https://nuxtjs.org/api/configuration-target
   */
  target: 'static',
  /*
   ** Headers of the page
   ** See https://nuxtjs.org/api/configuration-head
   */
  server: {
    host: '0.0.0.0', // default: localhost
  },
  head: {
    titleTemplate: '%s - ' + process.env.npm_package_name,
    title: process.env.npm_package_name || '',
    meta: [
      { charset: 'utf-8' },
      { name: 'viewport', content: 'width=device-width, initial-scale=1' },
      {
        hid: 'description',
        name: 'description',
        content: process.env.npm_package_description || '',
      },
    ],
    link: [
      { rel: 'icon', type: 'image/x-icon', href: '/favicon.ico' },
      {
        rel: 'stylesheet',
        href:
          'https://fonts.googleapis.com/css?family=Roboto:300,400,500,700|Material+Icons',
      },

      {
        rel: 'stylesheet',
        href: 'https://fonts.googleapis.com/css?family=Raleway',
      },
    ],
  },
  /*
   ** Global CSS
   */
  css: [],
  /*
   ** Plugins to load before mounting the App
   ** https://nuxtjs.org/guide/plugins
   */
  plugins: [],
  /*
   ** Auto import components
   ** See https://nuxtjs.org/api/configuration-components
   */
  components: true,
  /*
   ** Nuxt.js dev-modules
   */
  buildModules: ['@nuxt/typescript-build', '@nuxtjs/vuetify'],
  /*
   ** Nuxt.js modules
   */
  modules: ['@nuxtjs/axios', '@nuxtjs/apollo', '@nuxtjs/auth-next'],
  router: {
    middleware: ['auth'],
  },
  axios: {
    baseURL: 'http://localhost:8888',
    proxyHeaders: false,
    credentials: false,
  },
  apollo: {
    clientConfigs: {
      default: {
        httpEndpoint: 'http://10.5.12.6:31189/v1/graphql',
        wsEndpoint: 'ws://10.5.12.6:31189/v1/graphql',
        // httpEndpoint: 'http://localhost:8082/v1/graphql',
        // wsEndpoint: 'ws://localhost:8082/v1/graphql',
      },
    },
  },
  auth: {
    strategies: {
      local: false,
      keycloak: {
        scheme: 'oauth2',
        endpoints: {
          authorization: `http://localhost:8091/auth/realms/concierge/protocol/openid-connect/auth`,
          token: `http://localhost:8091/auth/realms/concierge/protocol/openid-connect/token`,
          logout:
            `http://localhost:8091/auth/realms/concierge/protocol/openid-connect/logout?redirect_uri=` +
            encodeURIComponent('http://localhost:3000/'),
          userInfo:
            'http://localhost:8091/auth/realms/concierge/protocol/openid-connect/userinfo',
        },
        user: {
          property: false, // <--- Default "user"
          autoFetch: true,
        },
        token: {
          property: 'access_token',
          type: 'Bearer',
          name: 'Authorization',
          maxAge: 1800, // Can be dynamic ?
        },
        refreshToken: {
          property: 'refresh_token',
          maxAge: 60 * 60 * 24 * 30, // Can be dynamic ?
        },
        responseType: 'code',
        token_type: 'Bearer',
        token_key: 'access_token',
        grantType: 'authorization_code',
        clientId: 'concierge-vue',
        scope: ['openid', 'profile', 'email'],
        codeChallengeMethod: 'S256',
      },
    },
  },

  /*
   ** vuetify module configuration
   ** https://github.com/nuxt-community/vuetify-module
   */
  vuetify: {
    customVariables: ['~/assets/variables.scss'],
    theme: {
      dark: false,
      themes: {
        dark: {
          primary: colors.blue.darken2,
          accent: colors.grey.darken3,
          secondary: colors.amber.darken3,
          info: colors.teal.lighten1,
          warning: colors.amber.base,
          error: colors.deepOrange.accent4,
          success: colors.green.accent3,
        },
        light: {
          primary: colors.blue.darken2,
          accent: colors.grey.darken3,
          secondary: colors.amber.darken3,
          info: colors.teal.lighten1,
          warning: colors.amber.base,
          error: colors.deepOrange.accent4,
          success: colors.green.accent4,
        },
      },
    },
  },
  /*
   ** Build configuration
   ** See https://nuxtjs.org/api/configuration-build/
   */
  build: {},
}

<template>
  <v-app dark>
    <v-navigation-drawer
      v-model="drawer"
      :mini-variant="miniVariant"
      :clipped="clipped"
      expand-on-hover
      permanent
      app
      width="300px"
    >
      <v-list>
        <v-list-item
          v-for="(item, i) in items"
          :key="i"
          :to="item.to"
          router
          exact
        >
          <v-list-item-action>
            <v-icon>{{ item.icon }}</v-icon>
          </v-list-item-action>
          <v-list-item-content>
            <v-list-item-title v-text="item.title" />
          </v-list-item-content>
        </v-list-item>
      </v-list>
    </v-navigation-drawer>
    <v-app-bar :clipped-left="clipped" fixed app>
      <!-- <v-app-bar-nav-icon @click.stop="drawer = !drawer" /> -->

      <v-img
        class="mx-2"
        src="/logo_concierge.png"
        max-height="50"
        max-width="50"
        contain
      ></v-img>

      <v-toolbar-title id="toolbar-title" v-text="title" />
      <v-spacer />

      <v-menu v-if="$store.state.auth.loggedIn" offset-y>
        <template v-slot:activator="{ on }">
          <!-- <v-btn text slot="activator"> -->
          <v-avatar slot="activator" color="primary" v-on="on">
            <v-icon dark> mdi-account-circle </v-icon>
            <!-- <v-icon right>expand_more</v-icon> -->
          </v-avatar>
        </template>
        <v-card class="mx-auto" max-width="344" outlined>
          <v-list-item three-line>
            <v-list-item-content>
              <div class="overline mb-4">
                {{ this.$auth.user.preferred_username }}
                <!-- {{ $auth.strategy.token.get() }} -->
              </div>
              <v-list-item-title class="headline mb-1">
                {{ this.$auth.user.given_name }}
                {{ this.$auth.user.family_name }}
              </v-list-item-title>
            </v-list-item-content>
          </v-list-item>

          <v-card-actions class="justify-center">
            <v-btn color="primary" @click="logout()">Logout</v-btn>
          </v-card-actions>
        </v-card>
      </v-menu>
    </v-app-bar>
    <v-main>
      <v-container>
        <nuxt />
      </v-container>
    </v-main>
  </v-app>
</template>

<script>
export default {
  data() {
    return {
      clipped: true,
      drawer: true,
      fixed: false,
      items: [
        {
          icon: 'mdi-view-dashboard',
          title: 'Dashboard',
          to: '/',
        },
        {
          icon: 'mdi-database',
          title: 'MinIO buckets',
          to: '/minio',
        },
        {
          icon: 'mdi-docker',
          title: 'Single Deployments',
          to: '/singleContainer',
        },
        {
          icon: 'mdi-arrow-decision',
          title: 'Distributed Deployments',
          to: '/distributed',
        },
        {
          icon: 'mdi-help-circle',
          title: 'Help',
          to: '/help',
        },
      ],
      miniVariant: false,
      right: true,
      rightDrawer: false,
      title: 'concierge',
    }
  },
  methods: {
    logout() {
      this.$auth.logout('keycloak')
    },
  },
}
</script>

<style scoped>
#toolbar-title {
  font-family: 'Raleway', sans-serif;
  font-size: 1.8em;
}
</style>

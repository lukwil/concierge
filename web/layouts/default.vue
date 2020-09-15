<template>
  <v-app dark>
    <v-navigation-drawer
      v-model="drawer"
      :mini-variant="miniVariant"
      :clipped="clipped"
      fixed
      app
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
      <v-app-bar-nav-icon @click.stop="drawer = !drawer" />

      <v-img
        class="mx-2"
        src="https://vectr.com/lukwil/h4Hf1xmLBA.png?width=250&height=350&select=h4Hf1xmLBApage0"
        max-height="50"
        max-width="50"
        contain
      ></v-img>

      <v-toolbar-title id="toolbar-title" v-text="title" />
      <v-spacer />
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
      clipped: false,
      drawer: false,
      fixed: false,
      items: [
        {
          icon: 'mdi-view-dashboard',
          title: 'Dashboard',
          to: '/',
        },
        {
          icon: 'mdi-docker',
          title: 'Single Container',
          to: '/singleContainer',
        },
        {
          icon: 'mdi-arrow-decision',
          title: 'Distributed',
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
}
</script>

<style scoped>
#toolbar-title {
  font-family: 'Raleway', sans-serif;
  font-size: 1.8em;
}
</style>

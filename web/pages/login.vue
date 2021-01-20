<template>
  <div>
    <br />
    <div class="text-h3 font-weight-light" align="center">Login</div>
    <br />
    <div class="text-center">
      <v-btn color="primary" @click="login()">Login via Keycloak</v-btn>
    </div>
  </div>
</template>

<script lang="ts">
import Vue from 'vue'
export default Vue.extend({
  methods: {
    login(): void {
      this.$auth.loginWith('keycloak')
    },
  },
})
</script>

<style lang="scss"></style>

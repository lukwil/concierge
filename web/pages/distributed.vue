<template>
  <v-card>
    <v-tabs
      v-model="activeTab"
      background-color="primary"
      centered
      dark
      icons-and-text
    >
      <v-tabs-slider></v-tabs-slider>

      <v-tab href="#tab-1">
        Overview
        <v-icon>mdi-format-list-bulleted</v-icon>
      </v-tab>
      <v-tab href="#tab-2">
        New Deployment
        <v-icon>mdi-plus-circle-outline</v-icon>
      </v-tab>
    </v-tabs>

    <v-tabs-items v-model="activeTab">
      <v-tab-item :key="1" :value="'tab-1'">
        <distributed-overview />
      </v-tab-item>

      <v-tab-item :key="2" :value="'tab-2'">
        <distributed-flow-stepper />
      </v-tab-item>
    </v-tabs-items>
  </v-card>
</template>

<script lang="ts">
import Vue from 'vue'
import DistributedFlowStepper from '@/components/distributed/flowStepper.vue'
import DistributedOverview from '@/components/distributed/overview.vue'
export default Vue.extend({
  components: {
    DistributedFlowStepper,
    DistributedOverview,
  },
  data() {
    return {
      activeTab: null,
      text:
        'Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat.',
    }
  },
})
</script>

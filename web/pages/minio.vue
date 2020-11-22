<template>
  <v-card flat>
    <br />
    <div class="text-h3 font-weight-light" align="center">MinIO Buckets</div>
    <br />
    <v-tabs
      v-model="tab"
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
        New MinIO bucket
        <v-icon>mdi-plus-circle-outline</v-icon>
      </v-tab>
    </v-tabs>

    <v-tabs-items v-model="tab">
      <v-tab-item :key="1" :value="'tab-1'">
        <minio-overview />
      </v-tab-item>

      <v-tab-item :key="2" :value="'tab-2'">
        <minio-flow-stepper />
      </v-tab-item>
    </v-tabs-items>
  </v-card>
</template>

<script lang="ts">
import Vue from 'vue'
import MinioFlowStepper from '@/components/minio/flowStepper.vue'
import MinioOverview from '@/components/minio/overview.vue'
export default Vue.extend({
  components: {
    MinioFlowStepper,
    MinioOverview,
  },
  data() {
    return {
      e6: 1,
      tab: null,
      text:
        'Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat.',
    }
  },
})
</script>

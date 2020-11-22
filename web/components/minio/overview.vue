<template>
  <v-layout column justify-center align-center>
    <v-card width="100%">
      <v-card-title>
        <v-spacer></v-spacer>
        <v-text-field
          v-model="search"
          append-icon="mdi-magnify"
          label="Search"
          single-line
          hide-details
        ></v-text-field>
      </v-card-title>
      <v-data-table
        :search="search"
        :headers="headers"
        :items="buckets"
        class="elevation-1"
        :height="height"
        fixed-header
        multi-sort
        :show-select="false"
        hide-default-footer
        disable-pagination
        loading-text="Loading... Please wait"
      >
        <template v-slot:[`item.actions`]="{ item }">
          <template v-if="item.status">
            <v-tooltip bottom>
              <template v-slot:activator="{ on, attrs }">
                <v-progress-circular
                  v-show="item.status.name == 'stopping'"
                  :size="20"
                  :width="2"
                  indeterminate
                  v-bind="attrs"
                  color="warning"
                  v-on="on"
                ></v-progress-circular>
              </template>
              <span>Deleting...</span>
            </v-tooltip>

            <v-tooltip bottom>
              <template v-slot:activator="{ on, attrs }">
                <v-progress-circular
                  v-show="item.status.name == 'starting'"
                  :size="20"
                  :width="2"
                  indeterminate
                  color="success"
                  v-bind="attrs"
                  v-on="on"
                ></v-progress-circular>
              </template>
              <span>Creating...</span>
            </v-tooltip>
          </template>
          <v-icon
            v-else
            medium
            class="mr-2"
            color="error"
            @click="showItem(item)"
            >mdi-alert-circle</v-icon
          >

          <v-tooltip bottom>
            <template v-slot:activator="{ on, attrs }">
              <v-icon
                medium
                class="mr-2"
                color="primary"
                v-bind="attrs"
                @click="editItem(item)"
                v-on="on"
                >mdi-square-edit-outline</v-icon
              >
            </template>
            <span>Edit</span>
          </v-tooltip>

          <v-tooltip bottom>
            <template v-slot:activator="{ on, attrs }">
              <v-icon
                medium
                color="error"
                v-bind="attrs"
                v-on="on"
                @click="deleteBucket(item.id)"
                >mdi-delete</v-icon
              >
            </template>
            <span>Delete</span>
          </v-tooltip>
        </template>
      </v-data-table>
    </v-card>
    <snackbar ref="snackbar"></snackbar>
  </v-layout>
</template>

<script lang="ts">
import Vue from 'vue'
import Snackbar from '@/components/Snackbar.vue'
import gql from 'graphql-tag'
import { Bucket } from '@/plugins/conciergeApi'
export default Vue.extend({
  apollo: {
    $subscribe: {
      minioBuckets: {
        query: gql`
          subscription {
            minio_bucket {
              id
              name
              status {
                name
              }
            }
          }
        `,
        result(data: any) {
          this.buckets = data.data.minio_bucket
        },
      },
    },
  },
  name: 'ContainerOverview',
  components: {
    Snackbar,
  },
  data() {
    return {
      snackbar: new Snackbar(),
      buckets: [] as Bucket[],
      search: '',
      headers: [
        { text: 'Name', value: 'name' },
        { text: 'Users', value: '' },
        { text: 'Actions', value: 'actions', sortable: false },
      ],

      height: 'calc(80vh - 150px)',
      editedIndex: -1,
      tab: [],
    }
  },
  computed: {},
  watch: {},
  mounted() {
    this.snackbar = this.$refs.snackbar as any
    console.log(this.tab)
  },
  methods: {
    deleteBucket(id: number) {
      this.$apollo.mutate({
        mutation: gql`
          mutation deleteMinioBucketByPk($id: Int!) {
            delete_minio_bucket_by_pk(id: $id) {
              id
            }
          }
        `,
        variables: {
          id,
        },
      })
    },
  },
})
</script>

<style scoped>
.v-data-table__wrapper {
  height: calc(100vh - 150px) !important;
}
</style>

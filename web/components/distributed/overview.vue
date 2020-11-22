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
        :items="deployments"
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
                <v-icon
                  :disabled="item.status.name !== 'running'"
                  medium
                  class="mr-2"
                  v-bind="attrs"
                  v-on="on"
                  @click="openInNewTab(item.name_k8s)"
                  >mdi-open-in-new</v-icon
                >
              </template>
              <span>Open in new tab</span>
            </v-tooltip>

            <v-tooltip bottom>
              <template v-slot:activator="{ on, attrs }">
                <v-icon
                  medium
                  class="mr-2"
                  v-bind="attrs"
                  v-on="on"
                  @click="showWarnings(item.name_k8s)"
                  >mdi-text-box</v-icon
                >
              </template>
              <span>Show warnings</span>
            </v-tooltip>

            <v-dialog v-model="dialog" :retain-focus="false" max-width="90vw">
              <v-card>
                <v-card-title>Warnings</v-card-title>
                <v-divider></v-divider>

                <v-data-table
                  :items="warnings"
                  :headers="dialogHeaders"
                  class="elevation-1"
                  fixed-header
                  multi-sort
                  :show-select="false"
                  hide-default-footer
                  disable-pagination
                  loading-text="Loading... Please wait"
                ></v-data-table>
                <!-- <v-divider></v-divider> -->
                <v-card-actions>
                  <v-spacer></v-spacer>
                  <v-btn
                    color="secondary"
                    text
                    @click="getWarnings(currentNameK8s)"
                    >Reload</v-btn
                  >
                  <v-btn color="error" text @click="closeDialog">Close</v-btn>
                </v-card-actions>
              </v-card>
            </v-dialog>

            <v-tooltip bottom>
              <template v-slot:activator="{ on, attrs }">
                <v-progress-circular
                  v-show="item.status.name == 'created'"
                  :size="20"
                  :width="2"
                  v-bind="attrs"
                  indeterminate
                  color="success"
                  v-on="on"
                ></v-progress-circular>
              </template>
              <span>Configuring/Starting...</span>
            </v-tooltip>

            <v-tooltip bottom>
              <template v-slot:activator="{ on, attrs }">
                <v-icon
                  v-show="item.status.name == 'running'"
                  medium
                  color="success"
                  class="mr-2"
                  v-bind="attrs"
                  v-on="on"
                  @click="stopContainer(item.name_k8s)"
                  >mdi-checkbox-blank-circle</v-icon
                >
              </template>
              <span>running...</span>
            </v-tooltip>

            <v-tooltip bottom>
              <template v-slot:activator="{ on, attrs }">
                <v-icon
                  v-show="item.status.name == 'succeeded'"
                  medium
                  color="success"
                  class="mr-2"
                  v-bind="attrs"
                  v-on="on"
                  >mdi-check-circle</v-icon
                >
              </template>
              <span>succeeded!</span>
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
                @click="deleteContainer(item.id)"
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
import { DistributedDeployment, Warning } from '@/plugins/conciergeApi'
export default Vue.extend({
  apollo: {
    $subscribe: {
      distributedDeployments: {
        query: gql`
          subscription {
            distributed_deployment {
              id
              name
              name_k8s
              container_image
              worker_count
              launcher_cpu
              launcher_ram
              worker_cpu
              worker_ram
              worker_gpu
              status {
                name
              }
            }
          }
        `,
        result(data: any) {
          console.log(data.data.distributed_deployment)
          this.deployments = data.data.distributed_deployment
        },
      },
    },
  },
  name: 'DistributedOverview',
  components: {
    Snackbar,
  },
  data() {
    return {
      snackbar: new Snackbar(),
      deployments: [] as DistributedDeployment[],
      currentNameK8s: '',
      search: '',
      headers: [
        { text: 'Name', value: 'name' },
        { text: 'Container image', value: 'container_image' },
        { text: 'Worker count', value: 'worker_count' },
        { text: 'Launcher: CPU', value: 'launcher_cpu' },
        { text: 'Launcher: RAM (MB)', value: 'launcher_ram' },
        { text: 'Worker: CPU', value: 'worker_cpu' },
        { text: 'Worker: RAM (MB)', value: 'worker_ram' },
        { text: 'Worker: GPU', value: 'worker_gpu' },
        { text: 'Actions', value: 'actions', sortable: false },
        // { text: 'Status', value: 'Status' },
      ],
      dialog: false,
      dialogHeaders: [
        { text: 'Timestamp', value: 'Timestamp' },
        { text: 'Reason', value: 'Reason' },
        { text: 'Message', value: 'Message' },
      ],
      warnings: [] as Warning[],

      height: 'calc(80vh - 150px)',
      editedIndex: -1,
      tab: [],
    }
  },
  computed: {},
  watch: {},
  mounted() {
    this.snackbar = this.$refs.snackbar as any
    // const q = this.$apollo.queries.hello
    console.log(this.tab)
    // window.setInterval(() => {
    //   this.initialize()
    // }, 2000)
  },
  methods: {
    showItem(item: any) {
      console.log(item)
    },
    openInNewTab(name: string) {
      const url: string = `http://10.5.12.6:31380/${name}/`
      window.open(url, '_blank')
    },
    showWarnings(name: string) {
      this.currentNameK8s = name
      this.dialog = true
      this.getWarnings(name)
    },
    getWarnings(name: string) {
      this.$apollo
        .query({
          query: gql`
            query containerWarnings(
              $nameK8s: ContainerWarningsInput! = { nameK8s: "" }
            ) {
              containerWarnings(nameK8s: $nameK8s) {
                Timestamp
                Reason
                Message
              }
            }
          `,
          variables: {
            nameK8s: { nameK8s: name },
          },
        })
        .then((data) => {
          console.log(data.data.containerWarnings)
          this.warnings = data.data.containerWarnings
        })
    },
    closeDialog() {
      this.dialog = false
      this.warnings = []
    },
    deleteContainer(id: number) {
      this.$apollo.mutate({
        mutation: gql`
          mutation deleteDistributedDeploymentByPk($id: Int!) {
            delete_distributed_deployment_by_pk(id: $id) {
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

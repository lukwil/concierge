<template>
  <v-layout column justify-center align-center>
    <v-card width="100%">
      <v-card-title>
        <!-- Single Deployments -->
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
                  :disabled="item.status.name !== 'started'"
                  medium
                  v-bind="attrs"
                  class="mr-2"
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

            <v-dialog v-model="dialog" :retain-focus="false">
              <v-card>
                <v-card-title>
                  <span class="headline">Warnings</span>
                </v-card-title>

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
                <v-icon
                  v-show="item.status.name == 'stopped'"
                  medium
                  color="success"
                  class="mr-2"
                  v-bind="attrs"
                  v-on="on"
                  @click="startContainer(item.name_k8s)"
                  >mdi-play</v-icon
                >
              </template>
              <span>Start</span>
            </v-tooltip>

            <v-tooltip bottom>
              <template v-slot:activator="{ on, attrs }">
                <v-icon
                  v-show="item.status.name == 'started'"
                  medium
                  color="warning"
                  class="mr-2"
                  v-bind="attrs"
                  v-on="on"
                  @click="stopContainer(item.name_k8s)"
                  >mdi-pause</v-icon
                >
              </template>
              <span>Stop</span>
            </v-tooltip>

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
              <span>Stopping...</span>
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
              <span>Starting...</span>
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
import { Deployment, Warning } from '@/plugins/conciergeApi'
export default Vue.extend({
  apollo: {
    $subscribe: {
      singleDeployments: {
        query: gql`
          subscription {
            single_deployment {
              id
              cpu
              gpu
              name
              name_k8s
              container_image
              ram
              status {
                name
              }
              volume {
                size
              }
            }
          }
        `,
        result(data: any) {
          this.deployments = data.data.single_deployment
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
      deployments: [] as Deployment[],
      currentNameK8s: '',
      search: '',
      headers: [
        { text: 'Name', value: 'name' },
        { text: 'User', value: 'user' },
        { text: 'Container image', value: 'container_image' },
        { text: 'CPU', value: 'cpu' },
        { text: 'RAM (MB)', value: 'ram' },
        { text: 'GPU', value: 'gpu' },
        { text: 'Volume Size (MB)', value: 'volume.size' },
        { text: 'Actions', value: 'actions', sortable: false },
        // { text: 'Status', value: 'Status' },
      ],
      dialog: false,
      dialogHeaders: [
        { text: 'Timestamp', value: 'Timestamp' },
        { text: 'Message', value: 'Message' },
      ],
      warnings: [] as Warning[],

      height: 'calc(80vh - 150px)',
      editedIndex: -1,
      tab: [],
      singleDeployments: [] as Deployment[],
    }
  },
  computed: {},
  watch: {},
  mounted() {
    const depl = []
    const d1 = {
      name: 'Experiment 1',
      user: 'Lukas Willburger',
      cpu: 2000,
      ram: 4096,
      gpu: 1,
      status: {
        name: 'started',
      },
      volume: { size: 512 },
      container_image: 'lukwil/jupyter-gpu',
    }
    depl.push(d1)
    const d2 = {
      name: 'Wichtige Berechnung',
      user: 'Lukas Willburger',
      cpu: 4000,
      ram: 8192,
      gpu: 2,
      status: {
        name: 'starting',
      },
      volume: { size: 10240 },
      container_image: 'lukwil/jupyter-gpu',
    }
    depl.push(d2)
    const d3 = {
      name: 'CPU-Berechnung',
      user: 'Max Mustermann',
      cpu: 4000,
      ram: 2048,
      // gpu: 1,
      status: {
        name: 'stopping',
      },
      // volume: { size: 512 },
      container_image: 'lukwil/jupyter-cpu',
    }
    depl.push(d3)
    const d4 = {
      name: 'Experiment 2',
      user: 'Lukas Willburger',
      cpu: 1000,
      ram: 4096,
      gpu: 1,
      status: {
        name: 'stopped',
      },
      // volume: { size: 512 },
      container_image: 'lukwil/jupyter-gpu',
    }
    depl.push(d4)
    const d5 = {
      name: 'Fehlerhaftes Experiment',
      user: 'Max Mustermann',
      cpu: 1000,
      ram: 4096,
      gpu: 1,
      // status: {
      //   name: 'stopped',
      // },
      // volume: { size: 512 },
      container_image: 'lukwil/jupyter-gpu',
    }
    depl.push(d5)
    this.deployments = depl

    // this.snackbar = this.$refs.snackbar as any

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
                Message
              }
            }
          `,
          variables: {
            nameK8s: { nameK8s: name },
          },
        })
        .then((data) => {
          this.warnings = data.data.containerWarnings
        })
    },
    closeDialog() {
      this.dialog = false
      this.warnings = []
    },
    stopContainer(name: string) {
      this.$apollo
        .mutate({
          mutation: gql`
            mutation stopStatefulSet($nameK8s: Input! = { nameK8s: "" }) {
              stopStatefulSet(nameK8s: $nameK8s) {
                Replicas
              }
            }
          `,
          variables: {
            nameK8s: { nameK8s: name },
          },
        })
        .then((data) => {
          console.log(data)
        })
    },
    startContainer(name: string) {
      this.$apollo
        .mutate({
          mutation: gql`
            mutation startStatefulSet($nameK8s: Input! = { nameK8s: "" }) {
              startStatefulSet(nameK8s: $nameK8s) {
                Replicas
              }
            }
          `,
          variables: {
            nameK8s: { nameK8s: name },
          },
        })
        .then((data) => {
          console.log(data)
        })
    },
    deleteContainer(id: number) {
      this.$apollo.mutate({
        mutation: gql`
          mutation deleteSingleDeploymentByPk($id: Int!) {
            delete_single_deployment_by_pk(id: $id) {
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

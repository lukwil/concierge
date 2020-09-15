<template>
  <v-layout column justify-center align-center>
    <v-card>
      <v-card-title>
        Deployments
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
            <v-icon
              :disabled="item.status.name !== 'started'"
              medium
              class="mr-2"
              @click="openInNewTab(item.name_k8s)"
              >mdi-open-in-new</v-icon
            >

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
            <v-icon medium class="mr-2" @click="showWarnings(item.name_k8s)"
              >mdi-text-box</v-icon
            >
            <v-icon medium class="mr-2" @click="editItem(item)"
              >mdi-information</v-icon
            >

            <v-icon
              v-if="item.status.name == 'stopped'"
              medium
              color="success"
              class="mr-2"
              @click="startContainer(item.name_k8s)"
              >mdi-play</v-icon
            >
            <v-icon
              v-if="item.status.name == 'started'"
              medium
              color="warning"
              class="mr-2"
              @click="stopContainer(item.name_k8s)"
              >mdi-pause</v-icon
            >
            <v-progress-circular
              v-if="item.status.name == 'stopping'"
              :size="20"
              :width="2"
              indeterminate
              color="warning"
            ></v-progress-circular>
            <v-progress-circular
              v-if="item.status.name == 'starting'"
              :size="20"
              :width="2"
              indeterminate
              color="success"
            ></v-progress-circular>
          </template>
          <v-icon
            v-else
            medium
            class="mr-2"
            color="error"
            @click="showItem(item)"
            >mdi-alert-circle</v-icon
          >
          <v-icon medium color="error" @click="deleteContainer(item.id)"
            >mdi-delete</v-icon
          >
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
      const url: string = `http://localhost:8080/${name}/`
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

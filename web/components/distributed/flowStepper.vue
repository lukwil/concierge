<template>
  <v-stepper v-model="e6" vertical>
    <v-stepper-step :complete="e6 > 1" step="1">
      Deployment name <small>Select your desired deployment name</small>
    </v-stepper-step>
    <v-stepper-content step="1">
      <v-form v-model="deploymentNameValid">
        <v-card-text>
          <v-container>
            <v-row>
              <v-col cols="12" sm="12" md="12">
                <v-text-field
                  v-model="editedItem.deploymentName"
                  :rules="deploymentNameRules"
                  label="Name"
                  type="text"
                ></v-text-field>
              </v-col>
            </v-row>
          </v-container>
        </v-card-text>
      </v-form>
      <!-- Disabled muss !valid sein !!!-->
      <v-btn :disabled="!deploymentNameValid" color="primary" @click="e6 = 2"
        >Continue</v-btn
      >
    </v-stepper-content>

    <v-stepper-step :complete="e6 > 2" step="2"
      >Container image
      <small>Select your desired container image</small></v-stepper-step
    >
    <v-stepper-content step="2">
      <v-form v-model="containerImageValid">
        <v-card-text>
          <v-container>
            <v-row>
              <v-col cols="12" sm="12" md="12">
                <v-text-field
                  v-model="editedItem.containerImage"
                  :rules="containerImageRules"
                  label="Container image"
                  type="text"
                ></v-text-field>
              </v-col>
            </v-row>
          </v-container>
        </v-card-text>
      </v-form>
      <v-btn text @click="e6 = 1">Back</v-btn>
      <v-btn :disabled="!containerImageValid" color="primary" @click="e6 = 3"
        >Continue</v-btn
      >
    </v-stepper-content>

    <v-stepper-step :complete="e6 > 3" step="3"
      >Environment variables
      <small>Select your desired environment variables</small></v-stepper-step
    >
    <v-stepper-content step="3">
      <v-form v-model="environmentVariablesValid">
        <v-card-text>
          <v-container>
            <!-- <v-row>
              <v-col cols="12" sm="12" md="12">
                <v-switch
                  v-model="editedItem.useEnvironmentVariables"
                  class="ma-1"
                  label="Use environmentVariables"
                ></v-switch>
              </v-col>
            </v-row> -->
            <v-row>
              <v-col cols="12" sm="5" md="5">
                <v-text-field
                  v-model="environmentVariableNew.name"
                  label="Name"
                  type="text"
                ></v-text-field>
              </v-col>
              <v-col cols="12" sm="5" md="5">
                <v-text-field
                  v-model="environmentVariableNew.value"
                  label="Value"
                  type="text"
                ></v-text-field>
              </v-col>
              <v-col cols="12" sm="2" md="2">
                <v-btn
                  fab
                  small
                  color="primary"
                  :disabled="
                    !(
                      !!environmentVariableNew.name.trim() &&
                      !!environmentVariableNew.value.trim()
                    )
                  "
                >
                  <v-icon dark medium @click="addEnvironmentVariable()"
                    >mdi-playlist-plus</v-icon
                  >
                </v-btn>
              </v-col>
            </v-row>

            <v-divider></v-divider>
            <br />

            <v-row
              v-for="(variable, index) in environmentVariables"
              :key="index"
            >
              <v-col cols="12" sm="5" md="5">
                <v-text-field
                  v-model="variable.name"
                  label="Name"
                  type="text"
                ></v-text-field>
              </v-col>
              <v-col cols="12" sm="5" md="5">
                <v-text-field
                  v-model="variable.value"
                  label="Value"
                  type="text"
                ></v-text-field>
              </v-col>
              <v-col cols="12" sm="2" md="2">
                <v-btn fab small color="error">
                  <v-icon dark medium @click="deleteEnvironmentVariable(index)"
                    >mdi-delete</v-icon
                  >
                </v-btn>
              </v-col>
            </v-row>
          </v-container>
        </v-card-text>
      </v-form>
      <v-btn text @click="e6 = 2">Back</v-btn>
      <v-btn :disabled="!containerImageValid" color="primary" @click="e6 = 4"
        >Continue</v-btn
      >
    </v-stepper-content>

    <v-stepper-step :complete="e6 > 4" step="4">
      <v-tooltip right>
        <template v-slot:activator="{ on, attrs }">
          <div>
            URL prefix
            <v-icon
              v-bind="attrs"
              small
              :color="e6 == 4 ? 'primary' : 'grey'"
              v-on="on"
              >mdi-information</v-icon
            >
          </div>
        </template>
        <span
          >You can later access your application via a subpath in the URL.
          <br />
          To get this to work, some applications like e.g. Jupyter Notebooks
          need a given base URL. <br />
          Thus you have the opportunity to set a URL prefix according to the
          application's needs. <br />
          The final URL prefix will be given as an environment variable
          <span class="secondary--text">URL_PREFIX</span>

          to use with your application.</span
        >
      </v-tooltip>

      <small>Select URL prefix</small></v-stepper-step
    >
    <v-stepper-content step="4">
      <v-form v-model="urlPrefixValid">
        <v-card-text>
          <v-container>
            <v-row>
              <v-col cols="12" sm="12" md="12">
                <v-switch
                  v-model="editedItem.overrideURLPrefix"
                  class="ma-1"
                  label="Override URL prefix"
                  @change="urlPrefixChange"
                ></v-switch>
              </v-col>
            </v-row>
            <v-row>
              <v-col cols="12" sm="4" md="4">
                <v-switch
                  v-show="editedItem.overrideURLPrefix"
                  v-model="editedItem.useContainerNameAsURLPrefix"
                  class="ma-1"
                  label="Use name given by Kubernetes"
                ></v-switch>
              </v-col>
              <v-col cols="12" sm="8" md="8">
                <v-text-field
                  v-show="editedItem.overrideURLPrefix"
                  v-model="editedItem.urlPrefix"
                  :rules="urlPrefixRules"
                  :disabled="editedItem.useContainerNameAsURLPrefix"
                  label="URL prefix"
                  type="text"
                  hint="Key in URL prefix with leading /"
                  persistent-hint
                ></v-text-field>
              </v-col>
            </v-row>
          </v-container>
        </v-card-text>
      </v-form>
      <v-btn text @click="e6 = 3">Back</v-btn>
      <v-btn :disabled="!urlPrefixValid" color="primary" @click="e6 = 5"
        >Continue</v-btn
      >
    </v-stepper-content>

    <v-stepper-step :complete="e6 > 5" step="5"
      >Resources <small>Select your desired resources</small></v-stepper-step
    >
    <v-stepper-content step="5">
      <v-form>
        <v-card outlined>
          <v-card-title class="primary white--text">Launcher</v-card-title>
          <v-card-text class="black--text">
            <v-form v-model="launcherResourcesValid">
              <v-card-text>
                <v-container>
                  <v-row>
                    <v-col cols="12" sm="4" md="4">
                      <v-text-field
                        v-model="editedItem.launcherCpu"
                        :rules="cpuRules"
                        label="CPU"
                        type="number"
                        hint="Key in CPU shares in thousand denomiation; e.g. 1 CPU equals 1000, 0.3 CPU equals 300"
                        persistent-hint
                      ></v-text-field>
                    </v-col>
                    <v-col cols="12" sm="4" md="4">
                      <v-text-field
                        v-model="editedItem.launcherRam"
                        :rules="ramRules"
                        label="Memory"
                        type="number"
                        hint="Key in memory shares in Megabytes; e.g. 2GB equals 2048, 0.5GB equals 512"
                        persistent-hint
                      ></v-text-field>
                    </v-col>
                    <v-col cols="12" sm="4" md="4">
                      <v-text-field
                        v-model="editedItem.workerCount"
                        :rules="workerCountRules"
                        label="Worker count"
                        type="number"
                        hint="Key in the desired number of workers"
                        persistent-hint
                      ></v-text-field>
                    </v-col>
                  </v-row>
                </v-container>
              </v-card-text>
            </v-form>
          </v-card-text>
        </v-card>
        <br />
        <v-card outlined>
          <v-card-title class="primary white--text">Worker</v-card-title>
          <v-card-text class="black--text">
            <v-form v-model="workerResourcesValid">
              <v-card-text>
                <v-container>
                  <v-row>
                    <v-col cols="12" sm="4" md="4">
                      <v-text-field
                        v-model="editedItem.workerCpu"
                        :rules="cpuRules"
                        label="CPU"
                        type="number"
                        hint="Key in CPU shares in thousand denomiation; e.g. 1 CPU equals 1000, 0.3 CPU equals 300"
                        persistent-hint
                      ></v-text-field>
                    </v-col>
                    <v-col cols="12" sm="4" md="4">
                      <v-text-field
                        v-model="editedItem.workerRam"
                        :rules="ramRules"
                        label="Memory"
                        type="number"
                        hint="Key in memory shares in Megabytes; e.g. 2GB equals 2048, 0.5GB equals 512"
                        persistent-hint
                      ></v-text-field>
                    </v-col>
                    <v-col cols="12" sm="4" md="4">
                      <v-text-field
                        v-model="editedItem.workerGpu"
                        :rules="gpuRules"
                        label="GPU"
                        type="number"
                      ></v-text-field>
                    </v-col>
                  </v-row>
                </v-container>
              </v-card-text>
            </v-form>
          </v-card-text>
        </v-card>
      </v-form>
      <br />

      <v-btn text @click="e6 = 4">Back</v-btn>
      <v-btn
        :disabled="!(launcherResourcesValid && workerResourcesValid)"
        color="primary"
        @click="e6 = 6"
        >Continue</v-btn
      >
    </v-stepper-content>

    <v-stepper-step :complete="e6 > 6" step="6"
      >Object storage
      <small>Add MinIO buckets (optional)</small></v-stepper-step
    >
    <v-stepper-content step="6">
      <v-form>
        <v-card-text>
          <v-container>
            <v-row>
              <v-col cols="12" sm="8" md="8">
                <v-select
                  v-model="existingBucketsSelected"
                  return-object
                  chips
                  multiple
                  item-text="name"
                  :items="existingBuckets"
                  placeholder="Choose MinIO bucket(s)..."
                  label="Existing bucket(s)"
                ></v-select>
              </v-col>
            </v-row>
          </v-container>
        </v-card-text>
      </v-form>
      <v-btn text @click="e6 = 5">Back</v-btn>
      <v-btn color="primary" @click="e6 = 7">Continue</v-btn>
    </v-stepper-content>

    <v-stepper-step :complete="e6 > 7" step="7"
      >Overview <small>Verify your entered data</small></v-stepper-step
    >
    <v-stepper-content step="7">
      <v-form>
        <v-card outlined>
          <v-card-title class="primary white--text">{{
            editedItem.deploymentName
          }}</v-card-title>
          <v-card-text class="black--text">
            <v-row>
              <v-col cols="12" lg="4" md="4" sm="12">
                <div class="font-weight-bold">Container image:</div>
                {{ editedItem.containerImage }}
              </v-col>
              <v-col
                v-if="editedItem.overrideURLPrefix"
                cols="12"
                lg="4"
                md="4"
                sm="12"
              >
                <div class="font-weight-bold">URL prefix:</div>
                <div v-if="!editedItem.useContainerNameAsURLPrefix">
                  {{ editedItem.urlPrefix }}
                </div>
                <div v-else class="secondary--text">
                  Use name given by Kubernetes
                </div>
              </v-col>
            </v-row>

            <v-row v-if="environmentVariables.length">
              <v-col cols="12" lg="12" md="12" sm="12">
                <div class="font-weight-bold">Environment variables:</div>
                <v-chip
                  v-for="(variable, index) in environmentVariables"
                  :key="index"
                  class="ma-2"
                  color="primary"
                >
                  {{ variable.name }}: {{ variable.value }}
                </v-chip>
              </v-col>
            </v-row>

            <v-row v-if="existingBucketsSelected.length">
              <v-col cols="12" lg="4" md="4" sm="12">
                <div class="font-weight-bold">Use existing bucket(s):</div>
                <v-chip
                  v-for="bucket in existingBucketsSelected"
                  :key="bucket.id"
                  class="ma-2 secondary"
                >
                  {{ bucket.name }}
                </v-chip>
              </v-col>
            </v-row>

            <v-divider></v-divider>

            <v-row>
              <v-col cols="12" lg="4" md="4" sm="12">
                <div class="font-weight-bold">Launcher CPU:</div>
                {{ editedItem.launcherCpu }}m (equals
                {{ editedItem.launcherCpu / 1000 }} core(s))
              </v-col>
              <v-col cols="12" lg="4" md="4" sm="12">
                <div class="font-weight-bold">Launcher Memory:</div>
                {{ editedItem.launcherRam }} MB
              </v-col>
              <v-col cols="12" lg="4" md="4" sm="12">
                <div class="font-weight-bold">Worker count:</div>
                {{ editedItem.workerCount }}
              </v-col>
            </v-row>

            <v-divider></v-divider>

            <v-row>
              <v-col cols="12" lg="4" md="4" sm="12">
                <div class="font-weight-bold">Worker CPU:</div>
                {{ editedItem.workerCpu }}m (equals
                {{ editedItem.workerCpu / 1000 }} core(s))
              </v-col>
              <v-col cols="12" lg="4" md="4" sm="12">
                <div class="font-weight-bold">Worker Memory:</div>
                {{ editedItem.workerRam }} MB
              </v-col>
              <v-col cols="12" lg="4" md="4" sm="12">
                <div class="font-weight-bold">Worker GPU:</div>
                {{ editedItem.workerGpu }}
              </v-col>
            </v-row>
          </v-card-text>
        </v-card>
      </v-form>
      <br />
      <v-btn text @click="e6 = 6">Back</v-btn>
      <v-btn color="success" @click="createDistributedDeployment()"
        >Create new deployment</v-btn
      >
    </v-stepper-content>

    <snackbar ref="snackbar"></snackbar>
  </v-stepper>
</template>

<script lang="ts">
import Vue from 'vue'
import { SingleContainerDto } from '@/plugins/conciergeApi'
import Snackbar from '@/components/Snackbar.vue'
import gql from 'graphql-tag'
export default Vue.extend({
  apollo: {
    $subscribe: {
      minioBuckets: {
        query: gql`
          subscription {
            minio_bucket {
              id
              name
            }
          }
        `,
        result(data: any) {
          console.log(data.data.minio_bucket)
          // this.existingBuckets = data.data.minio_bucket
        },
      },
    },
  },
  name: 'DistributedFlowStepper',
  components: {
    Snackbar,
  },
  data() {
    return {
      e6: 1,
      snackbar: new Snackbar(),
      editedItem: {
        deploymentName: '',
        containerImage: '',
        useEnvironmentVariables: false,
        overrideURLPrefix: false,
        urlPrefix: '/',
        useContainerNameAsURLPrefix: true,
        launcherCpu: 0,
        launcherRam: 0,
        workerCount: 0,
        workerCpu: 0,
        workerRam: 0,
        workerGpu: 0,
        useObjectStorage: false,
        useExistingBucket: false,
      },
      environmentVariableNew: {
        name: '',
        value: '',
      },
      environmentVariables: [] as { name: string; value: string }[],
      existingBuckets: [
        { id: 3, name: 'balksadf' },
        { id: 4, name: 'askdfjlka' },
      ],
      existingBucketsSelected: [],
      deploymentNameValid: true,
      deploymentNameRules: [
        (v: string) => !!v || 'Deployment name is mandatory!',
      ],
      containerImageValid: true,
      containerImageRules: [
        (v: string) => !!v || 'Container image is mandatory!',
      ],
      urlPrefixValid: true,
      urlPrefixRules: [(v: string) => !!v || 'URL prefix is mandatory!'],
      launcherResourcesValid: true,
      workerResourcesValid: true,
      cpuRules: [
        (v: number) => !!v || 'CPU is mandatory!',
        (v: number) => Number.isInteger(+v) || 'CPU must be an integer!',
        (v: number) => v > 0 || 'CPU must be a positive integer (>0)!',
      ],
      ramRules: [
        (v: number) => !!v || 'Memory is mandatory!',
        (v: number) => Number.isInteger(+v) || 'Memory must be an integer!',
        (v: number) => v > 0 || 'Memory must be a positive integer (>0)!',
      ],
      gpuRules: [
        (v: number) => !!v || 'GPU is mandatory!',
        (v: number) => Number.isInteger(+v) || 'GPU must be an integer!',
        (v: number) =>
          v >= 0 || 'GPU must be a positive integer or zero (>=0)!',
      ],
      workerCountRules: [
        (v: number) => !!v || 'Worker count is mandatory!',
        (v: number) =>
          Number.isInteger(+v) || 'Worker count must be an integer!',
        (v: number) =>
          v > 0 || 'Worker count must be a positive integer (>=1)!',
      ],
      persistenceValid: true,
      volumeSizeRules: [
        (v: number) => !!v || 'Volume size is mandatory!',
        (v: number) =>
          Number.isInteger(+v) || 'Volume size must be an integer!',
        (v: number) => v > 0 || 'Volume size must be a positive integer (>0)!',
      ],
      volumeMountPathRules: [
        (v: string) => !!v || 'Volume mount path is mandatory!',
      ],
    }
  },
  computed: {
    computedPersistenceValid: {
      get(): boolean {
        console.log('GOT')
        return this.persistenceValid
      },
      set(v: boolean) {
        console.log('SET' + v)
        this.persistenceValid = v
      },
    },
  },
  mounted() {
    this.snackbar = this.$refs.snackbar as any
  },
  methods: {
    urlPrefixChange(event: boolean) {
      if (!event) {
        this.urlPrefixValid = true
      } else {
        this.editedItem.urlPrefix = '/'
      }
    },
    addEnvironmentVariable() {
      const variable: { name: string; value: string } = {
        name: this.environmentVariableNew.name,
        value: this.environmentVariableNew.value,
      }

      this.environmentVariables.push(variable)

      this.environmentVariableNew.name = ''
      this.environmentVariableNew.value = ''
    },
    deleteEnvironmentVariable(index: number) {
      this.environmentVariables.splice(index, 1)
    },
    createDistributedDeployment() {
      const deployment: any = {
        name: this.editedItem.deploymentName,
        container_image: this.editedItem.containerImage,
        worker_count: this.editedItem.workerCount,
        launcher_cpu: this.editedItem.launcherCpu,
        launcher_ram: this.editedItem.launcherRam,
        worker_cpu: this.editedItem.workerCpu,
        worker_ram: this.editedItem.workerRam,
        worker_gpu: this.editedItem.workerGpu,
      }

      if (this.editedItem.overrideURLPrefix) {
        deployment.url_prefix = this.editedItem.urlPrefix
        if (this.editedItem.useContainerNameAsURLPrefix) {
          deployment.url_prefix = 'name_k8s'
        }
      }

      if (this.environmentVariables.length) {
        const distributedEnvironmentVariables = {
          data: this.environmentVariables,
        }
        deployment.distributed_environment_variables = distributedEnvironmentVariables
      }

      if (this.existingBucketsSelected.length) {
        const bucketIds = this.existingBucketsSelected.map((e: any) => ({
          minio_bucket_id: e.id,
        }))
        const deploymentBuckets = {
          data: bucketIds,
        }
        deployment.distributed_deployment_minio_buckets = deploymentBuckets
      }

      this.$apollo
        .mutate({
          mutation: gql`
            mutation insertDistributedDeployment(
              $deployment: distributed_deployment_insert_input!
            ) {
              insert_distributed_deployment_one(object: $deployment) {
                id
              }
            }
          `,
          variables: {
            deployment,
          },
        })
        .then((data) => {
          console.log(data)
        })
    },
  },
})
</script>
<style scoped></style>

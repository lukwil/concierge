<template>
  <v-stepper v-model="e6" vertical>
    <snackbar ref="snackbar"></snackbar>
    <v-stepper-step :complete="e6 > 1" step="1">
      Bucket name <small>Select your desired bucket name</small>
    </v-stepper-step>
    <v-stepper-content step="1">
      <v-form v-model="bucketNameValid">
        <v-card-text>
          <v-container>
            <v-row>
              <v-col cols="12" sm="12" md="12">
                <v-text-field
                  v-model="editedItem.bucketName"
                  :rules="bucketNameRules"
                  label="Name"
                  type="text"
                ></v-text-field>
              </v-col>
            </v-row>
          </v-container>
        </v-card-text>
      </v-form>
      <!-- Disabled muss !valid sein !!!-->
      <v-btn :disabled="!bucketNameValid" color="primary" @click="e6 = 2"
        >Continue</v-btn
      >
    </v-stepper-content>

    <v-stepper-step step="2"
      >Users
      <small
        >Assign bucket to additional users (optional)</small
      ></v-stepper-step
    >
    <v-stepper-content step="2">
      <v-form>
        <v-card-text>
          <v-container>
            <v-row>
              <v-col cols="12" sm="8" md="8">
                <v-select
                  v-model="existingUsersSelected"
                  chips
                  multiple
                  :items="existingUsers"
                  placeholder="Choose users..."
                  label="Additional users"
                ></v-select>
              </v-col>
            </v-row>
          </v-container>
        </v-card-text>
      </v-form>
      <v-btn text @click="e6 = 1">Back</v-btn>
      <v-btn color="success" @click="createBucket()"
        >Create new MinIO bucket</v-btn
      >
    </v-stepper-content>
  </v-stepper>
</template>

<script lang="ts">
import Vue from 'vue'
import Snackbar from '@/components/Snackbar.vue'
import gql from 'graphql-tag'
export default Vue.extend({
  name: 'SingleContainerFlowStepper',
  components: {
    Snackbar,
  },
  data() {
    return {
      e6: 1,
      snackbar: new Snackbar(),
      editedItem: {
        bucketName: '',
        additionalUsers: false,
      },
      existingUsers: [] as string[],
      existingUsersSelected: '',
      bucketNameValid: true,
      bucketNameRules: [(v: string) => !!v || 'Bucket name is mandatory!'],
    }
  },
  created() {
    this.initialize()
  },
  mounted() {
    this.snackbar = this.$refs.snackbar as any
  },
  methods: {
    initialize() {
      this.existingUsers = ['user1', 'user2', 'user3']
    },
    createBucket() {
      const bucket: any = {
        name: this.editedItem.bucketName,
      }

      this.$apollo
        .mutate({
          mutation: gql`
            mutation insertMinioBucket($container: minio_bucket_insert_input!) {
              insert_minio_bucket_one(object: $bucket) {
                id
              }
            }
          `,
          variables: {
            bucket,
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

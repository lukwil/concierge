<template>
  <v-snackbar
    v-model="show"
    top
    right
    multi-line
    :color="color"
    :timeout="timeout"
  >
    <v-card :color="color" flat>
      <v-card-title class="white--text">{{ title }}</v-card-title>
      <v-card-text class="white--text">{{ text }}</v-card-text>
    </v-card>
    <v-btn dark text @click="show = false">Schließen</v-btn>
  </v-snackbar>
</template>

<script lang="ts">
import Vue from 'vue'
import { AxiosResponse } from 'axios'

export default Vue.extend({
  name: 'Snackbar',
  data() {
    return {
      show: false,
      color: 'success',
      title: '',
      text: '',
      timeout: 6000,
    }
  },
  methods: {
    showSnackbar(
      response: AxiosResponse<any>,
      title: string = '',
      message: string = ''
    ) {
      let status = ''
      switch (response.status) {
        case 200: {
          if (!message) message = 'Datensatz wurde erfolgreich aktualisiert!'
          status = 'success'
          break
        }
        case 201: {
          if (!message) message = 'Datensatz wurde erfolgreich angelegt!'
          status = 'success'
          break
        }
        case 204: {
          if (!message) message = 'Datensatz wurde erfolgreich gelöscht!'
          status = 'success'
          break
        }
        case 400: {
          if (!message)
            message =
              'Format des Datensatzes für die Anfrage wurde nicht eingehalten!'
          status = 'warning'
          break
        }
        case 404: {
          if (!message)
            message =
              'Gewünschter Datensatz ist nicht (mehr) in der Datenbank vorhanden! \nBitte Seite neu laden!'
          status = 'error'
          break
        }
        case 409: {
          if (!message)
            message =
              'Ein zugehöriger Datensatz mit den gleichen Werten ist bereits vorhanden!'
          status = 'warning'
          break
        }
        case 500: {
          if (!message)
            message = 'Ein unerwarteter Serverfehler ist aufgetreten!'
          status = 'error'
          break
        }
      }
      this.title = title
      this.text = message
      this.color = status
      this.show = true
    },
  },
})
</script>

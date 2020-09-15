import Vue from 'vue'
import { ConciergeApi } from './conciergeApi'
declare module 'vue/types/vue' {
  interface Vue {
    $conciergeApi(): ConciergeApi
  }
}

Vue.prototype.$conciergeApi = () =>
  new ConciergeApi('http://192.168.178.22:3000')

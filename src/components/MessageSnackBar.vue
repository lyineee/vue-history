<template>
  <v-snackbar v-model="snackbar">
    {{ msg }}

    <template v-slot:action="{ attrs }">
      <v-btn color="pink" text v-bind="attrs" @click="snackbar = false">
        Close
      </v-btn>
    </template>
  </v-snackbar>
</template>

<script lang="ts">
import { Vue, Component } from "vue-property-decorator";
@Component
export default class MessageSnackBar extends Vue {
  msg = "";
  snackbar = false;
  created() {
    this.$root.$on("message", (data: string) => {
      this.snackbar = true;
      this.msg = data;
    });
  }
}
</script>

<style>
</style>
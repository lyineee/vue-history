<template>
  <v-dialog v-model="dialog" width="500">
    <template v-slot:activator="{ on, attrs }">
      <v-btn
        id="newButton"
        class="new-history mx-2"
        fab
        dark
        :small="isMobile()"
        :autofocus="true"
        v-bind="attrs"
        v-on="on"
        color="indigo"
      >
        <v-icon dark> mdi-plus </v-icon>
      </v-btn>
    </template>

    <v-card>
      <v-card-title class="text-h5 grey"> 输入链接 </v-card-title>

      <v-text-field
        :loading="btnLoading"
        clearable
        label="链接"
        outlined
        class="upload-textbox ma-4 mb-0"
        v-model="updateUrl"
      ></v-text-field>

      <v-divider></v-divider>

      <v-card-actions>
        <v-spacer></v-spacer>
        <v-btn color="primary" text @click="addHistory"> 上传 </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script  lang="ts">
import { environment } from "@/environments/environment";
import axios, { AxiosError } from "axios";
import { Vue, Component } from "vue-property-decorator";
@Component
export default class AddHistoryBtn extends Vue {
  dialog = false;
  btnLoading = false;
  updateUrl = "";
  logObject(obj: any): any {
    console.log("testesste");
    console.log(obj);
    return obj;
  }
  addHistory() {
    this.btnLoading = true;
    axios
      .get(
        `${environment.apiUrl}/history/update?url=${encodeURIComponent(
          this.updateUrl
        )}`
      )
      .then(() => {
        this.$root.$emit("message", "添加成功");
        this.updateUrl = "";
        this.dialog = false;
        setTimeout(() => {
          this.$emit("new");
        }, 3000);
      })
      .catch((err: AxiosError) => {
        console.log(err);
        this.$root.$emit("message", err.response?.data.message);
      })
      .finally(() => {
        this.btnLoading = false;
      });
  }
  isMobile() {
    switch (this.$vuetify.breakpoint.name) {
      case "xs":
        return true;
      default:
        return false;
    }
  }
}
</script>

<style>
#newButton {
  position: fixed;
  bottom: 8%;
  right: 2%;
  z-index: 200;
}
</style>
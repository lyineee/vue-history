<template>
  <v-app-bar app>
    <v-app-bar-title>记录</v-app-bar-title>
    <v-container>
      <v-row>
        <v-spacer />
        <v-col cols="12"  sm="8" md="6">
          <v-combobox
            ref="searchBox"
            v-model="searchList"
            :items="allTags"
            :search-input.sync="search"
            :loading="loading"
            class="mx-4"
            flat
            solo-inverted
            hide-details
            dense
            multiple
            clearable
            hide-selected
            prepend-inner-icon="search"
          ></v-combobox>
        </v-col>
        <v-spacer />
      </v-row>
    </v-container>
    <v-menu left bottom v-if="isMobile()">
      <template v-slot:activator="{ on, attrs }">
        <v-btn icon v-bind="attrs" v-on="on">
          <v-icon>mdi-dots-vertical</v-icon>
        </v-btn>
      </template>

      <v-list>
        <v-list-item @click="logout()" v-if="!edit">
          <v-list-item-title>退出登录</v-list-item-title>
          <v-icon right>mdi-logout</v-icon>
        </v-list-item>
        <v-list-item @click="toggleEdit()" transition="fade-transition">
          <v-list-item-title>{{
            edit ? "退出编辑" : "编辑"
          }}</v-list-item-title>
          <v-icon right>edit</v-icon>
        </v-list-item>
      </v-list>
    </v-menu>
    <v-expand-x-transition>
      <v-btn icon v-if="!isMobile() && !edit" @click="logout()">
        <v-icon>mdi-logout</v-icon>
      </v-btn>
    </v-expand-x-transition>
    <v-btn icon v-if="!isMobile()" @click="toggleEdit()">
      <v-icon>edit</v-icon>
    </v-btn>
  </v-app-bar>
</template>

<script  lang="ts">
import { environment } from "@/environments/environment";
import axios from "axios";
import {
  Vue,
  Component,
  Watch,
  VModel,
  Emit,
  Ref,
} from "vue-property-decorator";
@Component
export default class AppBar extends Vue {
  @Ref() searchBox!: any;
  @Watch("search") onSearchChanged(val: string) {
    this.loading = true;
    axios
      .get(`${environment.apiUrl}/history/tag?search=${val}`)
      .then((resp) => {
        this.allTags = resp.data;
      })
      .catch((err) => {
        console.log(err);
      })
      .finally(() => {
        this.loading = false;
      });
  }
  @VModel({ type: Array }) searchList!: Array<string>;
  edit = false;
  @Emit("edit")
  toggleEdit() {
    this.edit = !this.edit;
    return this.edit;
  }
  isMobile() {
    switch (this.$vuetify.breakpoint.name) {
      case "xs":
        return true;
      case "sm":
        return true;
      default:
        return false;
    }
  }
  logout() {
    localStorage.removeItem("user");
    localStorage.removeItem("accessToken");
    this.$emit("logout");
  }
  allTags = [];
  search = null;
  loading = false;
  created() {
    window.addEventListener("keydown", this.escapeListener);
  }
  escapeListener(event: KeyboardEvent) {
    if (event.key == "/" && !this.searchBox.isFocused) {
      event.preventDefault();
      this.searchBox.focus();
    }
  }
}
</script>
<style>
.search-box {
  max-width: 200px;
}
</style>
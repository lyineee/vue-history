<template>
  <v-app>
    <AppBar
      v-model="searchList"
      @logout="loginDialog = true"
      @edit="appBar.edit = $event"
    />
    <v-content>
      <LoginDialog v-model="loginDialog" />
      <HistoryGroup :search="searchList" />
    </v-content>
  </v-app>
</template>

<script lang="ts">
import HistoryGroup from "./components/HistoryGroup.vue";
import LoginDialog from "./components/LoginDialog.vue";
import AppBar from "./components/AppBar.vue";
import { Vue, Component, Provide, Watch } from "vue-property-decorator";

@Component({
  components: { HistoryGroup, LoginDialog, AppBar },
})
export default class App extends Vue {
  searchList = [];
  loginDialog = false;
  // onEditChange(event:boolean){

  // }
  colorScheme!: string;
  @Watch("colorScheme")
  onColorSchemeChange(val: string, old: string) {
    console.log("in", this.colorScheme);
    switch (this.colorScheme) {
      case "dark":
        this.$vuetify.theme.dark = true;
        break;
      case "light":
        this.$vuetify.theme.dark = false;
        break;
      default:
        this.$vuetify.theme.dark = false;
        break;
    }
  }
  @Provide() appBar = { edit: false };
  mounted() {
    console.log("on mounted hook");
    if (
      window.matchMedia &&
      window.matchMedia("(prefers-color-scheme: dark)").matches
    ) {
      this.$vuetify.theme.dark = true;
    }
    window
      .matchMedia("(prefers-color-scheme: dark)")
      .addEventListener("change", (e) => {
        this.$vuetify.theme.dark = e.matches;
      });
  }
}
</script>

<style lang="scss">
@import "https://fonts.googleapis.com/icon?family=Material+Icons";
</style>

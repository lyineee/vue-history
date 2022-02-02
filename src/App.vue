<template>
  <v-app>
    <AppBar />
    <v-main>
      <LoginDialog />
      <HistoryGroup />
      <MessageSnackBar />
    </v-main>
    <Footer />
  </v-app>
</template>

<script lang="ts">
import HistoryGroup from "./components/HistoryGroup.vue";
import LoginDialog from "./components/LoginDialog.vue";
import AppBar from "./components/AppBar.vue";
import MessageSnackBar from "./components/MessageSnackBar.vue";
import Footer from "./components/Footer.vue";
import { Vue, Component, Watch } from "vue-property-decorator";

@Component({
  components: { HistoryGroup, LoginDialog, AppBar, MessageSnackBar, Footer },
})
export default class App extends Vue {
  colorScheme!: string;
  @Watch("colorScheme")
  onColorSchemeChange() {
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
  mounted() {
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

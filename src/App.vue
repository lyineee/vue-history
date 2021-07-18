<template>
  <v-app>
    <AppBar v-model="searchList"  @logout="loginDialog = true" @edit="appBar.edit = $event"/>
    <v-content>
      <LoginDialog v-model="loginDialog"/>
      <HistoryGroup :search="searchList" />
    </v-content>
  </v-app>
</template>

<script lang="ts">
import HistoryGroup from "./components/HistoryGroup.vue";
import LoginDialog from "./components/LoginDialog.vue";
import AppBar from "./components/AppBar.vue";
import { Vue, Component, Provide } from "vue-property-decorator";

@Component({
  components: { HistoryGroup, LoginDialog, AppBar },
})
export default class App extends Vue {
  searchList = [];
  loginDialog = false;
  // onEditChange(event:boolean){
    
  // }
  @Provide() appBar = {edit:false}
}
</script>

<style lang="scss">
@import "https://fonts.googleapis.com/icon?family=Material+Icons";
</style>

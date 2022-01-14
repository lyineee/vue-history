<template>
  <v-container>
    <v-scale-transition group leave-absolute>
      <v-row class="" v-for="item in items" :key="item.id">
        <v-col cols="12">
          <HistoryItem :item="item" />
        </v-col>
      </v-row>
    </v-scale-transition>
    <div v-intersect="onIntersect" class="text-h7 text-center mt-3">
      {{ loadingMsg == "" ? "加载中..." : loadingMsg }}
    </div>
  </v-container>
</template>

<script lang="ts">
import { environment } from "@/environments/environment";
import { Vue, Component, Watch, Prop } from "vue-property-decorator";
import axios from "axios";
import HistoryItem from "./HistoryItem.vue";
import * as _ from "lodash";

@Component({
  components: { HistoryItem },
})
export default class HistoryGroup extends Vue {
  loadingMsg = "";
  onIntersect(entries: Array<IntersectionObserverEntry>) {
    if (!entries[0].isIntersecting) {
      return;
    }
    if (this.historyGenerator) {
      this.historyGenerator();
    } else {
      console.log("no items");
    }
  }
  items = [];
  historyGenerator: any;
  debounceSearchHistory: any;
  onProcess = false; // prevent multirequest in same time
  getHistoryGenerator(limit = 20, searchString = "", skip = 0) {
    let page = Math.floor(skip / limit);
    let clear = true;
    let end = false;
    let debounceSearchHistory = _.debounce(() => {
      console.log(this.onProcess);
      if (end || this.onProcess) {
        return;
      }
      this.onProcess = true;
      this.loadingMsg = "";
      axios
        .get(
          `${environment.apiUrl}/history/user?limit=${limit}&skip=${
            limit * page
          }${searchString ? "&search=" + searchString : ""}`
        )
        .then((resp) => {
          if (resp.data.length == 0) {
            end = true;
            this.loadingMsg = "到底了";
          }
          this.items = clear ? resp.data : this.items.concat(resp.data);
          clear = false;
          page++;
        })
        .catch((err) => {
          console.log(err);
        })
        .finally(() => {
          this.onProcess = false;
        });
    }, 500);
    return debounceSearchHistory;
  }
  created() {
    this.$root.$on("login", () => {
      console.log("login");
      this.historyGenerator = this.getHistoryGenerator(10);
      this.historyGenerator();
    });

    this.$root.$on("logout", () => {
      this.historyGenerator = null;
      this.items = [];
    });

    this.$root.$on("onSearchListChange", (search: Array<string>) => {
      this.historyGenerator = this.getHistoryGenerator(5, search.join(" "));
      this.historyGenerator();
    });

    this.historyGenerator = this.getHistoryGenerator(20);
    this.historyGenerator();
  }
}
</script>

<style>
</style>
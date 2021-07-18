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
  // item = {
  //   title: "test title asdfsfasf[asfafsda]fasfdsafdsafsajflkjsaljkjljs",
  //   tags: ["tag1", "tag2"],
  //   url: "https://bbs.nga.cn/read.php?tid=24740201&page=222",
  //   lastModified: "2021-02-18T11:28:39.522+00:00",
  //   current_page: 222,
  // };
  loadingMsg = "";
  @Prop() search!: Array<string>;
  @Watch("search")
  onSearchChange(val: Array<string>) {
    this.historyGenerator = this.getHistoryGenerator(5, this.search.join(" "));
    this.historyGenerator();
  }
  onIntersect() {
    if (this.items) {
      this.historyGenerator();
    } else {
      console.log("no items");
    }
  }
  items = [];
  historyGenerator: any;
  debounceSearchHistory: any;
  getHistoryGenerator(limit = 20, searchString = "", skip = 0) {
    let page = Math.floor(skip / limit);
    let clear = true;
    let end = false;
    let debounceSearchHistory = _.debounce(() => {
      if (end) {
        return;
      }
      this.loadingMsg = ""
      axios
        .get(
          `${environment.apiUrl}/history/user?limit=${limit}&skip=${
            limit * page
          }${searchString ? "&search=" + searchString : ""}`
        )
        .then((resp) => {
          if (resp.data.length == 0) {
            end = true
            this.loadingMsg = "到底了";
          }
          this.items = clear ? resp.data : this.items.concat(resp.data);
          clear = false;
          page++;
        })
        .catch((err) => {
          console.log(err);
        });
    }, 500);
    return debounceSearchHistory;
  }
  created() {
    this.historyGenerator = this.getHistoryGenerator(20);
    this.historyGenerator();
  }
}
</script>

<style>
</style>
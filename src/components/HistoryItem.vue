<template>
  <v-card class="mx-auto" max-width="800px" elevation="5" >
    <v-card-title class="text-h5 his-title">
      {{ item.title }}
    </v-card-title>
    <v-card-subtitle>
      {{ modifiedTime }}
    </v-card-subtitle>
    <v-expand-transition>
    <v-combobox
      ref="tagBox"
      v-if="item.tag.length != 0 || edit"
      v-model="item.tag"
      :items="allTags"
      :search-input.sync="search"
      :loading="loading"
      label="没有tag"
      flat
      chips
      deletable-chips
      multiple
      solo
      persistent-hint
      hide-selected
      append-icon=""
      :readonly="!edit"
      @change="update()"
    >
      <template v-slot:no-data>
        <v-list-item>
          <v-list-item-content>
            <v-list-item-title>
              没有找到 "<strong>{{ search }}</strong
              >". 按下 <kbd>enter</kbd> 以创建新tag
            </v-list-item-title>
          </v-list-item-content>
        </v-list-item>
      </template>
    </v-combobox>
    </v-expand-transition>
    <v-card-actions>
      <v-btn block color="primary" @click="openUrl()"> 打开 </v-btn>
    </v-card-actions>
  </v-card>
</template>

<script lang="ts">
import { environment } from "@/environments/environment";
import axios from "axios";
import { Vue, Component, Watch, Prop, Inject, Ref } from "vue-property-decorator";

@Component
export default class HistoryItem extends Vue {
  // @Prop() edit = true;
  @Ref() tagBox!:any
  @Prop() item: any;
  @Inject() readonly appBar !: any
  edit = false;
  @Watch("appBar",{deep:true})
  onAppBarDataChange(){
    this.edit = this.appBar.edit
  }
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
  readonly = true;
  allTags = [];
  search = null;
  loading = false;
  get modifiedTime(): string {
    let date = new Date(this.item.lastModified);
    return date.toLocaleString()
  }
  openUrl() {
    window.open(this.item.url, "_blank");
  }

  update: any; // update tag function
  created() {
    (() => {
      let last = this.item.tag;
      this.update = () => {
        if (last != this.item.tag) {
          last = this.item.tag;
          let data = new FormData();
          data.append("tag", this.item.tag);
          axios
            .put(`${environment.apiUrl}/history/${this.item.id}/tag`, data)
            .then()
            .catch((err) => {
              console.log(err);
            });
        }
      };
    })();
    // window.addEventListener("keydown", (event)=>{
    //   if(event.key == "/" && this.tagBox.isFocused == true){
    //     event.stopPropagation()
    //   }
    // });
  }
}
</script>

<style>
.his-title {
  display: block;
}
</style>
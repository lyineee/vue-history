<template>
  <v-card class="mx-auto" max-width="800px" elevation="5">
    <v-card-title class="text-h5 his-title">
      {{ item.title }}
    </v-card-title>
    <v-card-subtitle>
      {{ modifiedTime }}
      <span class="float-right">上次看到: {{ item.current_page }}页</span>
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
    <v-card-actions class="card-action">
      <!-- <v-expand-x-transition> -->
      <!-- <div> -->
      <v-text-field
        :loading="btnLoading"
        clearable
        class="upload-textbox"
        v-model="updateUrl"
      ></v-text-field>
      <!-- </div> -->
      <!-- </v-expand-x-transition> -->
      <!-- <v-expand-x-transition> -->
      <v-btn
        :loading="btnLoading"
        :class="edit ? 'action-upload-btn' : 'action-open-btn'"
        class="action-btn"
        color="primary"
        @click="edit ? updatePage() : openUrl()"
      >
        <v-scroll-y-transition>
          <span v-if="!edit" style="position: absolute"> 打开</span>
        </v-scroll-y-transition>
        <v-scroll-y-transition>
          <span v-if="edit" style="position: absolute"> 确认</span>
        </v-scroll-y-transition>
      </v-btn>
      <!-- </v-expand-x-transition> -->
    </v-card-actions>
  </v-card>
</template>

<script lang="ts">
import { environment } from "@/environments/environment";
import axios from "axios";
import { Vue, Component, Watch, Prop, Ref } from "vue-property-decorator";

@Component
export default class HistoryItem extends Vue {
  @Ref() tagBox!: any;
  @Prop() item: any;
  edit = false;
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
  updateUrl = "";
  readonly = true;
  allTags = [];
  search = null;

  loading = false;
  btnLoading = false;
  get modifiedTime(): string {
    let date = new Date(this.item.lastModified);
    return date.toLocaleString();
  }
  openUrl() {
    window.open(this.item.url, "_blank");
  }

  updatePage() {
    this.btnLoading = true;
    axios
      .get(
        `${environment.apiUrl}/history/update?url=${encodeURIComponent(
          this.updateUrl
        )}`
      )
      .then(() => {
        this.$root.$emit("message", "上传成功");
      })
      .catch((err) => {
        console.log(err);
      })
      .finally(() => {
        this.btnLoading = false;
      });
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
    this.updateUrl = this.item.url;
    this.$root.$on("onEditStateChange", (editState: boolean) => {
      this.edit = editState;
    });
  }
}
</script>

<style>
.his-title {
  display: block;
}
.card-action {
  height: 52px;
  justify-content: flex-end;
  overflow: hidden;
}
.upload-textbox {
  margin: 0px 20px;
}
/* .upload-textbox-collapse {
  max-width: 0%;
} */
.action-open-btn {
  width: 100%;
}
.action-upload-btn {
  width: 20%;
}
button.action-btn {
  transition: width 0.3s;
}
</style>
<template>
  <v-card
    class="mx-auto"
    max-width="800px"
    elevation="5"
    @contextmenu="showEditMenu"
    v-intersect="onIntersect"
  >
    <v-menu
      v-model="showMenu"
      :position-x="x"
      :position-y="y"
      absolute
      offset-y
    >
      <v-list class="py-0">
        <v-list-item transition="fade-transition" @click="toggleEdit()">
          <v-list-item-title>{{
            edit ? "退出编辑" : "编辑"
          }}</v-list-item-title>
          <v-icon right>edit</v-icon>
        </v-list-item>
        <v-list-item transition="fade-transition" @click="refresh()">
          <v-list-item-title> 刷新 </v-list-item-title>
          <v-icon right>refresh</v-icon>
        </v-list-item>
        <v-list-item
          class="error-color"
          transition="fade-transition"
          @click="deleteDialog = true"
        >
          <v-list-item-title class="error-text"> 删除记录 </v-list-item-title>
          <v-icon right color="white">delete</v-icon>
        </v-list-item>
      </v-list>
    </v-menu>
    <v-card-title class="text-h5 his-title">
      {{ item.title }}
    </v-card-title>
    <v-card-subtitle>
      {{ modifiedTime }}
      <span class="float-right"
        >上次看到: {{ item.current_page }}/{{
          item.total_page ? item.total_page : "-"
        }}页</span
      >
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
      <v-text-field
        :loading="btnLoading"
        clearable
        class="upload-textbox"
        v-model="updateUrl"
      ></v-text-field>
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
    </v-card-actions>
    <DeleteConfirmDialog v-model="deleteDialog" @confirm="deleteHistory()" />
  </v-card>
</template>

<script lang="ts">
import { environment } from "@/environments/environment";
import axios from "axios";
import { Vue, Component, Watch, Prop, Ref } from "vue-property-decorator";
import DeleteConfirmDialog from "./DeleteConfirmDialog.vue";

@Component({
  components: { DeleteConfirmDialog },
})
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
        this.$root.$emit("message", err.response?.data.message);
      })
      .finally(() => {
        this.loading = false;
      });
  }
  updateUrl = "";
  readonly = true;
  allTags = [];
  search = null;

  showMenu = false;
  x = 0;
  y = 0;

  loading = false;
  btnLoading = false;
  get modifiedTime(): string {
    if (this.item.lastModified) {
      let date = new Date(this.item.lastModified);
      return date.toLocaleString();
    }
    return "";
  }

  deleteDialog = false;
  showEditMenu(e: MouseEvent) {
    e.preventDefault();
    this.showMenu = false;
    this.x = e.clientX;
    this.y = e.clientY;
    this.$nextTick(() => {
      this.showMenu = true;
    });
  }
  openUrl() {
    window.open(this.item.url, "_blank");
  }

  toggleEdit() {
    this.edit = !this.edit;
  }
  refresh() {
    this.$emit("update");
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
        this.edit = false;
        this.$emit("update");
      })
      .catch((err) => {
        console.log(err);
        this.$root.$emit("message", err.response?.data.message);
      })
      .finally(() => {
        this.btnLoading = false;
      });
  }

  deleteHistory() {
    this.loading = true;
    axios
      .delete(`${environment.apiUrl}/history/${this.item.id}`)
      .then(() => {
        this.$root.$emit("message", "删除成功");
        this.$emit("delete");
      })
      .catch((err) => {
        console.log(err);
        this.$root.$emit("message", err.response?.data.message);
      })
      .finally(() => {
        this.loading = false;
      });
  }

  isIntersecting = false;
  refreshIntervalId?: number;
  onIntersect(
    entries: Array<IntersectionObserverEntry>,
    observer: IntersectionObserver
  ) {
    if (
      this.isIntersecting != entries[0].isIntersecting &&
      entries[0].isIntersecting == true
    ) {
      //console.debug("start auto refresh: ", this.item.title);
      this.refreshIntervalId = setInterval(() => {
        this.$emit("update");
      }, 1000 * 10); // refresh every 10 second
    }
    if (!entries[0].isIntersecting && this.refreshIntervalId) {
      //console.debug("stop auto refresh: ", this.item.title);
      clearInterval(this.refreshIntervalId);
    }
    this.isIntersecting = entries[0].isIntersecting;
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
              this.$root.$emit("message", err.response?.data.message);
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
.action-open-btn {
  width: 100%;
}
.action-upload-btn {
  width: 20%;
}
button.action-btn {
  transition: width 0.3s;
}
.error-color {
  /* color: red; */
  background-color: #ff5252;
}
.error-text {
  color: white;
}
</style>
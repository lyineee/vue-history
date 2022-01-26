<template>
  <div class="text-center">
    <v-dialog v-model="deleteDialog" width="500">
      <v-card>
        <v-card-title class="text-h5"> 删除记录? </v-card-title>

        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn color="primary darken-1" text @click="deleteDialog = false">
            取消
          </v-btn>
          <v-btn color="error" text @click="emit('confirm')">
            删除
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </div>
</template>

<script lang="ts">
import { Vue, Component, VModel } from "vue-property-decorator";

@Component({})
export default class DeleteConfirmDialog extends Vue {
  @VModel() deleteDialog!: boolean;
}
</script>

<style>
</style>
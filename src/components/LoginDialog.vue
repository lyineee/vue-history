<template>
  <v-dialog v-model="dialog" persistent max-width="290">
    <v-card :loading="loading">
      <template slot="progress">
        <v-progress-linear height="3" indeterminate></v-progress-linear>
      </template>
      <v-card-title class="text-h5"> 登录 </v-card-title>
      <v-card-text>登录以使用功能</v-card-text>
      <v-container class="m-1">
        <v-row>
          <v-col cols="12" class="mx-auto">
            <v-form ref="form" v-model="valid">
              <v-text-field
                @keydown.enter="login"
                v-model="username"
                :rules="[rules.max, rules.requireUsername]"
                label="用户名"
                required
              >
              </v-text-field>
              <v-text-field
                @keydown.enter="login"
                v-model="password"
                :rules="[rules.max, rules.requirePassword]"
                label="密码"
                :append-icon="showPassword ? 'mdi-eye' : 'mdi-eye-off'"
                :type="showPassword ? 'text' : 'password'"
                @click:append="showPassword = !showPassword"
                required
                clearable
              >
              </v-text-field>
            </v-form>
          </v-col>
        </v-row>
      </v-container>
      <v-card-actions>
        <v-spacer></v-spacer>
        <v-btn
          color="green darken-1"
          :disabled="!valid"
          :loading="loading"
          block
          text
          @click="login()"
        >
          登录
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
  <!-- </v-row> -->
</template>

<script lang="ts">
import axios from "axios";
import { environment } from "@/environments/environment";
import { Vue, Component, VModel, Watch } from "vue-property-decorator";

@Component
export default class LoginDialog extends Vue {
  dialog = false;
  showPassword = false;
  loading = false;
  username = "";
  password = "";
  valid = false;

  rules = {
    max: (v: string): any => {
      let test = !!v;
      if (!test) {
        return false;
      } else {
        return v.length <= 25 || "最多25个字符";
      }
    },
    requirePassword: (v: any): any => {
      return !!v || "请输入密码";
    },
    requireUsername: (v: any): any => {
      return !!v || "请输入用户名";
    },
  };
  login() {
    if (!this.valid) {
      return;
    }
    this.loading = true;
    axios
      .post(`${environment.apiUrl}/auth/token`, {
        username: this.username,
        password: this.password,
      })
      .then((resp) => {
        localStorage.setItem("user", JSON.stringify(resp.data));
        localStorage.setItem("accessToken", resp.data["access_token"]);
        this.dialog = false;
        this.setToken(resp.data["access_token"]);
        this.$root.$emit("login");
      })
      .catch((err) => {
        console.error(err);
      })
      .finally(() => (this.loading = false));
  }
  setToken(token: string) {
    axios.interceptors.request.use((config) => {
      config.headers = {
        Authorization: `Bearer ${token}`,
      };
      return config;
    });
  }

  revokeToken() {
    axios.interceptors.request.use((config) => {
      config.headers = {
        Authorization: "",
      };
      return config;
    });
  }

  created() {
    this.$root.$on("logout", () => {
      localStorage.removeItem("user");
      localStorage.removeItem("accessToken");
      this.revokeToken();
      this.dialog = true;
    });

    console.log("Try to grab user token.");
    let token = localStorage.getItem("accessToken");
    if (token == null) {
      console.log("Dont get access token, login first.");
      this.dialog = true;
      return;
    } else {
      console.log("Successfully get access token.");
      this.setToken(token);
    }
    console.log("Registry login fallback operaion.");
    axios.interceptors.response.use((resp) => {
      if (resp.status == 401) {
        console.log("Auth fail, try to login.");
        this.dialog = true;
        return resp;
      }
      if (resp.status != 200) {
        console.log("Unexpected error occur");
      }
      return resp;
    });
    // fire
    axios
      .get(`${environment.apiUrl}/auth/user`)
      .then()
      .catch((err) => {
        console.log(err);
      });
  }
}
</script>

<style>
</style>
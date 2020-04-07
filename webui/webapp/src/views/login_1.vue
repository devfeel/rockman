<template>
  <div id="index-container">
      <Form ref="formInline">
        <FormItem>
            <a href="javascript:void()" class="demo-badge"><div class="txt"><span>任务调度系统</span></div></a>
        </FormItem>
        <FormItem>
            <i-input name="token" size="large"  type="password" password placeholder="输入口令" v-model="token" />
        </FormItem>
        <FormItem>
            <i-button type="primary" size="large" v-on:click="login" long>提交</i-button>
        </FormItem>
    </Form>
</div>
</template>
<script>
import { login } from '@/api/login.js';
export default {
  data() {
    return {
      token: ''
    };
  },
  methods: {
    login() {
      login({ token: this.token }).then((data) => {
        if (data.Code === 200) {
          this.$Message.info('登陆成功,正在跳转!');
          this.$store.commit('SET_TOKEN', this.token)
          window.sessionStorage.setItem('Token', this.token)
          this.$router.push({path: 'home'})
        } else {
          this.$Message.warning(data.Msg);
        }
      });
    }
  }
};
</script>
<style lang="less" scoped>

#index-container {
    transform: translateY(-50%);
    /* background: #ae9696; */
    top: 40%;
    position: absolute;
    margin: 0 auto;
    left: 0;
    width: 350px;
    line-height: 100px;
    right: 0;
    text-align: center;
}

#index-container form > div {
    padding: 5px;
    margin-bottom: 15px;
}
.txt {
    color: #999;
    height: 40px;
    font-size: 28px;
    font-weight: bold;
}

</style>

<template>
  <Layout class="content-main">
    <Sider hide-trigger :style="{background: '#f5f7f9',padding: '24px 0px 0px 0px'}">
        <Menu active-name="1" theme="light" width="auto" :open-names="['1']">
            <MenuItem name="1" to="/static/settings/nodes"><Icon type="ios-keypad"></Icon>NODE管理</MenuItem>
            <MenuItem name="2" to="/static/settings/users"><Icon type="ios-keypad"></Icon>用户管理</MenuItem>
        </Menu>
    </Sider>
    <Layout :style="{padding: '0 0'}">
        <Content :style="{padding: '24px', minHeight: '280px', background: '#fff'}">
            <keep-alive>
            <router-view />
          </keep-alive>
        </Content>
    </Layout>
</Layout>
</template>
<script>
export default {
  data() {
    return {
      token: ''
    };
  },
  methods: {
    loginOut() {
      this.$store.commit('SET_TOKEN', null)
      window.sessionStorage.setItem('Token', null)
      this.$router.push({ path: 'home' })
    }
  }
}
</script>
<style scoped>
.content-main{
    position: absolute;
    left: 0;
    right: 0;
    padding: 0 0px;
    height: calc(100%);
    /* height: calc(100% - 104px); */
}
.ivu-menu{
    background:transparent;
}
.ivu-menu-light.ivu-menu-vertical .ivu-menu-item-active:not(.ivu-menu-submenu){
    background: rgb(230,247,255);
}

</style>

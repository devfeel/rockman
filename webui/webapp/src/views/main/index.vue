<template>
  <el-container>
    <el-header  class="main-header">
      <!-- <div class="layout-logo"></div> -->
      <div class="layout-header-text"><el-link :underline="false" class="layout-header-text">Rockman</el-link></div>
      <div class="layout-nav">
        <el-menu :default-active="activeData" mode="horizontal" router @select="selectMenu">
          <el-menu-item  index="/static/home" >首页</el-menu-item>
          <el-menu-item  index="/static/node" >节点中心</el-menu-item>
          <el-menu-item  index="/static/task" >运行管理</el-menu-item>
        </el-menu>
      </div>
      <div class="layout-header-user">
        <el-dropdown @command="onDropDownItemClick">
        <span class="el-dropdown-link">
          管理员<i class="el-icon-arrow-down el-icon--right"></i>
        </span>
        <el-dropdown-menu slot="dropdown">
          <el-dropdown-item command="loginOut">退出登录</el-dropdown-item>
        </el-dropdown-menu>
      </el-dropdown>
      </div>
    </el-header>
    <el-main >
      <div class="main-content">
        <keep-alive>
            <router-view />
        </keep-alive>
      </div>
    </el-main>
  </el-container>
</template>
<script>
export default {
  data() {
    return {
      token: '',
      activeData: '1'
    }
  },
  mounted() {
      var nameVal = window.sessionStorage.getItem('selectMenu');
      if (nameVal) {
        this.activeData = nameVal;
      } else {
        this.activeData = '1';
      }
  },
  methods: {
    onDropDownItemClick(command) {
      console.log(this.$store.state)
      if (command === 'loginOut') {
        this.loginOut()
      }
    },
    selectMenu(name) {
      window.sessionStorage.setItem('selectMenu', name)
    },
    loginOut() {
      debugger;
      this.$store.commit('SET_TOKEN', null)
      window.sessionStorage.removeItem('Token')
      this.$store.commit('SET_INFO', null)
      window.sessionStorage.removeItem('UserInfo')
      window.sessionStorage.removeItem('selectMenu')
      this.$router.push({ path: '/static/login' })
    }
  }
}
</script>
<style scoped>
.main {
  height: 100%;
}
.main-header,
.main-header-menu,
.el-menu-item {
  height: 64px;

}
.el-menu-item{
  width: 130px;
  text-align: center;
}

.ivu-layout-header {
  /* background-color:#E7E8EC; */
}
.main-content {

  left: 0;
  right: 0;
  padding: 0 0px;
  margin: 0 auto;
  /* margin-top: 40px; */
  /* background: #EDF0F5; */
  min-height: calc(100% - 64px);
}
.main-content-con {
  padding: 15px;
  min-width: 280px;
  background-color: #fff;
  height: calc(100% - 40px);
}
.layout-logo {
  width: 55px;
  height: 55px;
  border-radius: 3px;
  float: left;
  position: relative;
  top: 2px;
}
.layout-header-text {
  width: 80px;
  height: 64px;
  padding-left: 10px;
  border-radius: 3px;
  color:rgb(24, 144, 255);
  font-size: 24px;
  float: left;
  font-family: "Helvetica Neue",Helvetica,"PingFang SC","Hiragino Sans GB","Microsoft YaHei","微软雅黑",Arial,sans-serif;
}
.layout-header-user {
  width: 150px;
  height: 30px;
  border-radius: 3px;
  /* color: #fff; */
  font-size: 18px;
  float: right;
}
.layout-header-user a {
  /* color: #fff; */
  font-size: 14px;
}

.layout-nav {
  /* width: 420px; */
  padding-left: 50px;
  float: left;
  /* margin: 0 auto; */
  /* margin-right: 20px; */
}
.ivu-menu-item-active {
  background: #33a6f7;
}
</style>

<template>
  <Layout>
    <Header class="main-header">
      <div class="layout-logo"><img src="~@/assets/rockman3@2x.png" style="height:55px;width:55px;"></div>
      <div class="layout-header-text"></div>
      <div class="layout-nav">
        <Menu mode="horizontal" theme="dark" :active-name="activeData" class="main-header-menu" @on-select="selectMenu">
          <MenuItem name="1" :to='{name:"home"}'>
          <!-- <Icon type="ios-navigate"></Icon> -->
          监控
          </MenuItem>
          <MenuItem name="2" :to='{name:"settings"}'>
          <!-- <Icon type="ios-keypad"></Icon> -->
          配置中心
          </MenuItem>
          <MenuItem name="3" :to='{name:"runtimes"}'>
          <!-- <Icon type="ios-analytics"></Icon> -->
          运行管理
          </MenuItem>
        </Menu>
      </div>
      <div class="layout-header-user">
        <!-- <a v-on:click="loginOut">安全退出</a> -->
        <Dropdown @on-click="onDropDownItemClick">
            <a href="javascript:void(0)">
                管理员
                <Icon type="ios-arrow-down"></Icon>
            </a>
            <DropdownMenu slot="list" >
              <DropdownItem >个人资料</DropdownItem>
              <DropdownItem name="loginOut">安全退出</DropdownItem>
            </DropdownMenu>
        </Dropdown>
      </div>
    </Header>
    <Layout>
      <Layout>
        <Content class="main-content">
          <keep-alive>
            <router-view />
          </keep-alive>
        </Content>
      </Layout>
    </Layout>
  </Layout>
</template>
<script>
export default {
  data() {
    return {
      token: '',
      activeData: '1'
    };
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
    onDropDownItemClick(name) {
      console.log(this.$store.state)
      if (name === 'loginOut') {
        this.loginOut()
      }
    },
    selectMenu(name) {
      window.sessionStorage.setItem('selectMenu', name)
    },
    loginOut() {
      this.$store.commit('SET_TOKEN', null)
      window.sessionStorage.removeItem('Token')
      this.$store.commit('SET_INFO', null)
      window.sessionStorage.removeItem('UserInfo')
      window.sessionStorage.removeItem('selectMenu')
      this.$router.push({ path: 'login' })
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
.ivu-menu-item {
  height: 64px;
}
.ivu-menu-item {
 width: 130px;
 text-align: center;
}

.ivu-layout-header {
  /* background-color:#E7E8EC; */
}
.main-content {
  position: absolute;
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
  width: 200px;
  height: 30px;
  position: relative;
  border-radius: 3px;
  color: #fff;
  font-size: 18px;
  float: left;
}
.layout-header-user {
  width: 150px;
  height: 30px;
  border-radius: 3px;
  color: #fff;
  font-size: 18px;
  float: right;
}
.layout-header-user a {
  color: #fff;
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

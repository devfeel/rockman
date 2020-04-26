<template>
  <div id="app" >
    <router-view/>
  </div>
</template>

<script>
  export default {
    name: 'App',
    computed: {},
    created() {
    }
  }
</script>
<style lang="less">
body {
  margin: 0px 0px;
  padding: 0px 0px;
  font-size: 14px;
  width: 100%;
  height: 100%;
  overflow: hidden;
}
#app {
  height: 1500px;
}
</style>

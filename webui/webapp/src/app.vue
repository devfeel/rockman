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
  font-size: 14px;
  width: 100%;
  height: 100%;
  overflow: hidden;
}
#app {
  height: 1500px;
}
</style>

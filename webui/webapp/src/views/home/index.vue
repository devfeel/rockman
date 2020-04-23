<template>
  <div class="content">
    <el-card class="cluster-info">

    </el-card>
    <el-card class="message-info">
      <div class="message-info-btn">
        <el-radio-group v-model="radioValue" size="medium">
          <el-radio-button label="节点消息" fill="#409EFF"></el-radio-button>
          <el-radio-button label="任务消息" fill="#409EFF"></el-radio-button>
        </el-radio-group>
      </div>
      <div class="message-info-content">
        <el-scrollbar style="height:100%">
          <el-timeline :reverse="reverse"
          v-infinite-scroll="load"
          infinite-scroll-disabled="disabled">
            <el-timeline-item
              v-for="(timeLine, index) in timeLineData"
              :key="index"
              :timestamp="timeLine.time" icon='el-icon-more' type='primary' placement='top'>
              <el-card>
                <el-tag type='warning'>{{timeLine.node}}</el-tag> {{timeLine.message}}
              </el-card>
            </el-timeline-item>
          </el-timeline>
          <div class="message-info-more" v-if="loading"><el-link :underline="false" @click="load">更多</el-link></div>
          <div class="message-info-more" v-if="timeLineData.length>10">没有更多了</div>
        </el-scrollbar>
      </div>
    </el-card>
  </div>
</template>
<script>
import { getClusterInfo } from '@/api/cluster.js';
export default {
  data() {
    return {
      radioValue: '节点消息',
      reverse: false,
      timeLineData: [],
      loading: false
    };
  },
  activated() {
    this.onInit();
  },
  methods: {
    onInit() {
      getClusterInfo().then(res => {
        if (res.RetCode === 0) {
        } else {
          this.$Message.warning(res.RetMsg);
        }
      })
      this.timeLineData.push({
        time: '2020-04-23 16:25:10',
        node: '10.139.160.174:40001',
        message: 'test-job任务加入'
      });
      this.timeLineData.push({
        time: '2020-04-23 14:25:10',
        node: '10.139.160.174:40001',
        message: 'test-job任务移出'
      });
      this.timeLineData.push({
        time: '2020-04-23 12:25:10',
        node: '10.139.160.174:40001',
        message: 'http-job任务加入'
      });
      this.loading = true;
    },
    load() {
      this.timeLineData.push({
        time: '2020-04-23 12:25:10',
        node: '10.139.160.174:40001',
        message: 'http-job任务加入'
      });
    },
    noMore () {
        return this.timeLineData.length >= 20
    },
    disabled () {
      return this.loading || this.noMore
    }
  }
};
</script>
<style lang="less" scoped>
.content {
  background-color:rgb(243, 243, 243);
  position: absolute;
  left: 0;
  right: 0;
  padding: 0 0px;
  margin: 0 auto;
  /* margin-top: 40px; */
  /* background: #EDF0F5; */
  height: calc(100%);
  overflow-y:scroll;
}
.cluster-info{
  margin: 10px 20px;
  background: white;
  height: 200px;
}
.message-info{
  // position: absolute;
  // left: 0;
  // right: 0;
  // padding: 0 0px;
  margin: 10px 20px;
  background: white;
  //height: calc(100%-30px);
}
.message-info-btn{
  margin-bottom:10px;
  height: 45px;
}
.message-info-content{
  padding: 10px;
  // height: 200px;
  //height: calc(100%-40px);
}
.el-scrollbar {
  height: 100%;
}
.el-scrollbar__wrap { overflow: scroll; width: 110%; height: 120%; }
.message-info-more{
  text-align: center;
}
</style>

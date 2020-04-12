<template>
    <div>
      <div style>
        <div data-v-542f4644 class="ivu-row" style="padding:15px;background: white;">
          <div class="headTitle">
              <Form label-position="left" :label-width="200" :model="clusterInfo">
                <Row>
                  <Col span="5">
                      <FormItem label="集群编码："><span v-text="clusterInfo.ClusterId"></span></FormItem>
                  </Col>
                  <Col span="5">
                      <FormItem label="注册服务地址："><span v-text="clusterInfo.RegistryServerUrl"></span></FormItem>
                  </Col>
                  <Col span="5">
                      <FormItem label="Leader服务器：">
                          <span v-text="clusterInfo.LeaderServer"></span>
                      </FormItem>
                  </Col>
                </Row>
              </Form>
          </div>
            <div
              v-for="item in topColor"
              :key="item.name"
              class="ivu-col ivu-col-span-6"
              style="padding-left: 8px; padding-right: 8px;"
            >
              <div data-v-542f4644 class="ivu-card" :style="{background:item.background}">
                <div class="icon-left">
                  <Icon :type="item.icon" />
                </div>
                <div class="ivu-card-body">
                  <div data-v-542f4644 class="demo-color-name">{{item.name}}</div>
                  <div data-v-542f4644 class="demo-color-desc">{{item.desc}}</div>
                </div>
              </div>
            </div>
        </div>
      </div>
    </div>
</template>
<script>
import { getClusterInfo } from '@/api/cluster.js';
export default {
    data() {
        return {
            clusterInfo: {
              RegistryServerUrl: '',
              LeaderServer: ''
            },
             topColor: [
                {
                name: 'Node数量',
                desc: '205',
                background: 'rgb(25, 190, 107)',
                icon: 'ios-home'
                },
                {
                name: '任务数量',
                desc: '412',
                background: 'rgb(45, 183, 245)',
                icon: 'ios-help-buoy'
                },
                {
                name: '累计执行次数',
                desc: '200',
                background: 'rgb(255, 153, 0)',
                icon: 'md-ionic'
                },
                {
                name: '累计异常次数',
                desc: '1020',
                background: 'rgb(237, 64, 20)',
                icon: 'ios-navigate'
                }
            ]
        }
    },
    mounted() {
      this.init();
    },
     methods: {
        init() {
          getClusterInfo().then(res => {
            if (res.RetCode === 0) {
                this.clusterInfo = res.Message;
            } else {
                this.$Message.error(res.RetMsg);
            }
          })
        }
     }
}
</script>
<style scoped>
.home-contianer {
  background: #efefef;
  width: 100%;
  height: 100%;
  /* padding: 20px; */
}
.headTitle {
  padding-left: 50px;
  font-size: 24px;
}
.home-app {
  display: inline-block;
  /* display: -webkit-flex;
  display: flex; */
  padding: 15px;
  padding-top: 5px;
}
.home-app > div {
  float: left;
  width: 33.33333%;
  padding-left: 8px;
  padding-right: 8px;
}
.ivu-card-body {
  text-align: center;
  padding: 25px 13px;
  padding-left: 80px;
}
.demo-color-name {
  color: #fff;
  font-size: 16px;
}
.demo-color-desc {
  color: #fff;
  opacity: 0.7;
}
.ivu-card {
  position: relative;
}
.ivu-card .icon-left {
  border-right: 1px solid;
  padding: 10px 24px;
  height: 100%;
  position: absolute;
  font-size: 50px;
  color: white;
}
.ivu-row {
  border-bottom: 2px dotted #eee;
  padding: 15px;
  margin-bottom: 15px;
}

.h5-desc {
  padding-top: 10px;
}
</style>
<style lang="less">
.charts {
  display: inline-block;
  width: 100%;
  margin-top: 20px;
  // padding: 0px 24px;
  .left {
    padding: 25px;
    background: white;
    height: 360px;
    width: 49%;
    float: left;
    margin-right: 1%;
    background: white;
  }
  .right {
    padding: 25px 45px;
    background: white;
    height: 360px;
    width: 49%;
    float: left;
    margin-left: 1%;
    .badge-count {
      padding: 3px 7px;
      position: relative;
      border: 1px solid #eee;
      border-radius: 50%;

      margin-right: 11px;
    }
    .badge {
      background: #e2e2e2;
      color: #3a3535;
    }
    .top3 {
      background: #2db7f5;
      color: white;
    }
    .cell {
      position: relative;
      display: flex;
      padding: 10px 0;
      border-bottom: 1px dotted #eee;
    }
    .primary {
      flex: 1;
    }
    .title {
      font-size: 16px;
      padding-bottom: 6px;
      border-bottom: 1px solid #eee;
      margin-bottom: 11px;
    }
    .name {
      font-size: 15px;
      position: relative;
      top: 5px;
      color: #303133;
      left: 12px;
    }
    .desc {
      margin-left: 27px;
      font-size: 12px;
      color: #b3b3b3;
      position: relative;
      top: 5px;
    }
  }
}
</style>

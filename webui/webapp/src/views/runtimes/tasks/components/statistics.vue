<template>
    <div>
        <div style="height:258px;" :id="bar"></div>
        <div style="height:258px;" :id="pie"></div>
    </div>
</template>
<script>
import options from './chartOptions';
let echarts = require('echarts');
export default {
    data() {
        return {
        bar: 'b-' + ~~(Math.random(10000, 100000) * 100000),
        pie: 'p-' + ~~(Math.random(10000, 100000) * 100000),
        options: options
        };
    },
    props: {
      data: {}
    },
    mounted() {
        this.options.bar.title = {
        subtext: '',
        x: 'center'
        };
        let $bar = echarts.init(document.getElementById(this.bar));
        $bar.setOption(this.options.bar);

        this.options.pie.legend.top = 50;
        this.options.pie.legend.right = 80;
        this.options.pie.legend.orient = 'vertical';
        let $pie = echarts.init(document.getElementById(this.pie));
        $pie.setOption(this.options.pie);
    }
}
</script>

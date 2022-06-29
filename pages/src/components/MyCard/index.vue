<template>
  <div class="MyCard">
    <el-card class="box-card mxCard">
      <div class="table">
        <span class="titleSpan">{{ title }}</span>
      </div>
      <slot />
    </el-card>
  </div>
</template>

<script>
export default {
  name: 'MyCard',
  props: {
    title: { type: String, default: '' }
  },
  data() {
    return {}
  },
  mounted() {
    // document.getElementsByClassName('mxCard')[0].style.minHeight = document.getElementById('app').clientHeight - 170 + 'px'
  },
  methods: {}
}
</script>

<style rel="stylesheet/scss" lang="scss" scoped>
  .mxCard {
    border-radius: 10px;
    height: 100%;
    min-height: 89vh;
    overflow: visible;
    .table {
      border-bottom: 2px solid rgba(0, 0, 0, 0.1);
      padding-bottom: 10px;
      margin-bottom: 15px;
      .titleSpan {
        font-size: 25px;
        font-weight: 700;
      }
    }
  }
</style>

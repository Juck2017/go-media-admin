<template>
  <div class="app-container">
    <my-card title="分屏播放">
      <video-player />
    </my-card>
  </div>
</template>
<script>
import MyCard from '@/components/MyCard'
import VideoPlayer from './videoPlayer'
export default {
  name: 'SplitScreen',
  components: {
    VideoPlayer,
    MyCard
  }
}
</script>
<style rel="stylesheet/scss" lang="scss" scoped>
</style>

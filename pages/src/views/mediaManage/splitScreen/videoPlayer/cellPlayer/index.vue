<template>
  <div class="player">
    player{{ title }}
    {{ currentIndex === title - 1 }}
  </div>
</template>
<script>
import '@/js/jsmpeg.min'
export default {
  name: 'CellPlayer',
  props: {
    title: {
      type: Number,
      default: 0
    },
    currentIndex: {
      type: Number,
      default: 0
    }
  },
  methods: {
    playVide() {
      const url = 'ws://127.0.0.1:5000'
      this.$nextTick(function() {
        const canvasDiv = document.getElementById('canvasDiv')
        const canvas = document.createElement('canvas')
        canvas.style.width = '640px'
        canvas.style.height = '350px'
        canvasDiv.appendChild(canvas)
            // eslint-disable-next-line
                this.player = new JSMpeg.Player(`${url}?code=${row.code}`,
          {
            canvas,
            onPlay: function() {
              console.log('play video')
            }
          })
      })
    }
  }
}
</script>
<style rel="stylesheet/scss" lang="scss" scoped>
.player {
    background-color: black;
    height: 100%;
    border: 1px solid white;
    color: white;
    text-align: center;
}
</style>

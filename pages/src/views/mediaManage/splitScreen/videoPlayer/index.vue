<template>
  <div class="cell">
    <el-row :gutter="20" style="margin:0px;width:100%">
      <el-col :span="6">
        <div class="grid-content bg-purple">
          <div class="cell-tool">
            <el-button size="small" type="primary" @click="cellCount = 1">1</el-button>
            <el-button size="small" type="primary" @click="cellCount = 4">4</el-button>
            <el-button size="small" type="primary" @click="cellCount = 9">9</el-button>
            <el-button size="small" type="primary" @click="cellCount = 16">16</el-button>
            <el-tree
              ref="tree"
              :data="treeData"
              node-key="code"
              :props="defaultProps"
              highlight-current
              check-on-click-node
              :current-node-key="currentNodeKey"
              @node-click="handleNodeClick"
            />
          </div>
        </div>
      </el-col>
      <el-col :span="18">
        <div class="grid-content bg-purple">
          <div class="cell-player">
            <div v-for="i in cellCount" :key="i" :class="cellClass(i)">
              <cell-player :title="i" :current-index="currentIndex" />
            </div>
          </div>
        </div>
      </el-col>
    </el-row>

  </div>
</template>
<script>
import CellPlayer from './cellPlayer'
export default {
  name: 'VideoPlayer',
  components: {
    CellPlayer
  },
  data() {
    return {
      cellCount: 4,
      currentNodeKey: '',
      currentIndex: 0,
      treeData: [
        {
          id: 1,
          label: '上海申铁在线',
          children: [{
            id: 3,
            label: '摄像头007',
            code: 'CHANNEL20220616007'
          },
          {
            id: 4,
            label: '摄像头008',
            code: 'CHANNEL20220616008'
          }]
        }, {
          id: 2,
          label: '坂银通道养护',
          children: [{
            id: 5,
            label: '摄像头003',
            code: 'CHANNEL20220616003'
          }, {
            id: 6,
            label: '摄像头004',
            code: 'CHANNEL20220616004'
          }]
        }],
      defaultProps: {
        children: 'children',
        label: 'label'
      }
    }
  },
  computed: {
    cellClass() {
      return function(index) {
        switch (this.cellCount) {
          case 1:
            return ['cell-player-1']
          case 4:
            return ['cell-player-4']
          case 9:
            return ['cell-player-9']
          case 16:
            return ['cell-player-16']
          default:
            break
        }
      }
    }
  },
  watch: {
    currentNodeKey(code) {
      if (code) {
        console.log('this.currentIndex', this.currentIndex)
        if (this.currentIndex > this.cellCount - 2) {
          this.currentIndex = 0
        } else {
          this.currentIndex++
        }
        console.log('current code:', code)
      }
    },
    cellCount() {
      this.currentIndex = 0
    }
  },
  methods: {
    handleNodeClick(data) {
      this.currentNodeKey = data.code
    }
  }
}
</script>
<style rel="stylesheet/scss" lang="scss" scoped>
.cell {
  height: 100%;

  .cell-tool {
    height: 736px;
    line-height: 30px;
    padding: 2.5px 7px;
  }

  .cell-player {
    display: flex;
    justify-content: space-between;
    height: 736px;
    flex-wrap: wrap;

    .cell-player-1 {
      width: 100%;
      box-sizing: border-box;
    }

    .cell-player-4 {
      width: 50%;
      height: 50% !important;
      box-sizing: border-box;
    }

    .cell-player-9 {
      width: 33.33%;
      height: 33.33% !important;
      box-sizing: border-box;
    }

    .cell-player-16 {
      width: 25%;
      height: 25% !important;
      box-sizing: border-box;
    }

  }

  .bg-purple {
    background: #d3dce6 !important;
  }

  .grid-content {
    border-radius: 4px;
    min-height: 36px;
  }
}
</style>

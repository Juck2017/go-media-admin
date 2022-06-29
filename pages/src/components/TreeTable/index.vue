<template>
  <el-table
    v-loading="listLoading"
    :data="formatData"
    :row-style="showRow"
    v-bind="$attrs"
    size="small"
    empty-text="暂无数据"
    element-loading-text="加载中..."
  >
    <!-- 列无标题数据 -->
    <el-table-column v-if="columns.length===0" width="150">
      <template slot-scope="scope">
        <span v-for="space in scope.row._level" :key="space" class="ms-tree-space" />
        <span v-if="iconShow(0,scope.row)" class="tree-ctrl" @click="toggleExpanded(scope.$index)">
          <i v-if="!scope.row._expanded" class="el-icon-plus" />
          <i v-else class="el-icon-minus" />
        </span>
        {{ scope.$index }}
      </template>
    </el-table-column>

    <!-- 不是按钮和图标操作 -->
    <template v-for="(column, index) in columns">
      <el-table-column
        v-if="column.type!='button'&&column.type!='iconButton'&&columns.length>0"
        :key="column.value"
        :label="column.text"
        :width="column.width"
      >
        <template slot-scope="scope">
          <!-- 这一行控制前后空格 -->
          <template v-for="space in scope.row.treeLevel">
            <span v-if="index === 0" :key="space" class="ms-tree-space" />
          </template>
          <!-- 展开与收缩图标 -->
          <span v-if="iconShow(index,scope.row)" class="tree-ctrl" @click="toggleExpanded(scope.$index)">
            <i v-if="!scope.row._expanded" class="el-icon-plus" />
            <i v-else class="el-icon-minus" />
          </span>
          <span v-html="formatter(scope.row[column.value],column.formatter,scope.row)" />
        </template>
      </el-table-column>
    </template>

    <template v-for="column in columns">
      <el-table-column
        v-if="column.type==='button'&&columns.length>0"
        :key="column.value"
        :label="column.text"
        :width="column.width"
        align="center"
      >
        <template v-for="(key, num) in column.list(scope.row)" slot-scope="scope">
          <el-button
            v-if="column.type==='button'"
            :key="num"
            :type="key.type"
            size="mini"
            @click="key.click(scope,$event)"
          >{{ key.value }}
          </el-button>
        </template>
      </el-table-column>
    </template>

    <!-- 操作按钮 -->
    <template v-for="(column, index) in columns">
      <el-table-column
        v-if="column.type === 'iconButton' && columns.length>0"
        :key="index"
        :label="column.text"
        :width="column.width"
        align="center"
      >
        <template slot-scope="scope">
          <div v-if="column.type === 'iconButton' && columns.length>0">
            <span v-for="(key, num) in column.list(scope.row)" :key="num">
              <!-- v-button-privilege="key.privilege" -->
              <el-tooltip :content="key.value" placement="top">
                <i style="font-size: 25px;" class="icon iconfont projectColor pointer" :class="key.class" @click="key.click(scope,$event)" />
              </el-tooltip>
              <span class="projectColor" style="position: absolute;margin-left: -4px;font-weight: 700;opacity: 0.5;">|</span>
            </span>
          </div>
        </template>
      </el-table-column>
    </template>
    <slot />
  </el-table>
</template>

<script>
import treeToArray from './eval'
// import buttonPrivilege from '@/directive/buttonPrivilege'
export default {
  name: 'TreeTable',
  directives: {
    // buttonPrivilege
  },
  props: {
    data: {
      type: [Array, Object],
      required: true
    },
    columns: {
      type: Array,
      default: () => []
    },
    evalFunc: {
      type: Function,
      default: () => {}
    },
    evalArgs: {
      type: Array,
      default: () => []
    },
    expandAll: {
      type: Boolean,
      default: false
    },
    listLoading: {
      type: Boolean,
      default: false
    }
  },
  computed: {
    // 格式化数据源
    formatData: function() {
      // console.log('执行formatData函数', this.data)
      let tmp
      if (!Array.isArray(this.data)) {
        tmp = [this.data]
      } else {
        if (this.data.length > 0) {
          // 组装参数
          tmp = this.treeListUtil(this.data)
        } else {
          tmp = this.data
        }
      }
      let func // 定义方法
      let args // 定义参数
      // const args = this.evalArgs ? Array.concat([tmp, this.expandAll], this.evalArgs) : [tmp, this.expandAll]
      if (this.evalFunc.toString() === 'function _default() {}') {
        // 是默认的函数,则设置系统函数
        func = treeToArray
        args = this.evalArgs.concat([tmp, this.expandAll])
      } else {
        // 传递了新的函数
        func = this.evalFunc
        args = [tmp, this.expandAll]
      }
      // apply为函数内部调用,第一个参数this,第二个参数为传递的参数(数组)
      return func.apply(null, args)
    }
  },
  methods: {
    // 设置行是否隐藏和展示, 返回对象
    showRow: function(row) {
      let returnStyle = {}
      const show = (row.row.parent ? (row.row.parent._expanded && row.row.parent._show) : true)
      row.row._show = show
      if (show) {
        returnStyle = {
          animation: 'treeTableShow 1s',
          '-webkit-animation': 'treeTableShow 1s'
        }
      } else {
        returnStyle = {
          display: 'none'
        }
      }
      return returnStyle
    },
    // 数据组织
    treeListUtil(data, parentId) {
      const itemArr = []
      for (let i = 0; i < data.length; i++) {
        const node = data[i]
        if (node.parentId === parentId) {
          node.children = this.treeListUtil(data, node.id)
          itemArr.push(node)
        }
      }
      return itemArr
    },
    // 切换下级是否展开
    toggleExpanded: function(trIndex) {
      const record = this.formatData[trIndex]
      record._expanded = !record._expanded
    },
    // 图标显示
    iconShow(index, record) {
      return (index === 0 && record.children && record.children.length > 0)
    },
    formatter(value, formatter, row) {
      if (!formatter) {
        return value
      } else {
        return formatter(row)
      }
    }
  }
}
</script>
<style rel="stylesheet/css">
  @keyframes treeTableShow {
    from {
      opacity: 0;
    }
    to {
      opacity: 1;
    }
  }

  @-webkit-keyframes treeTableShow {
    from {
      opacity: 0;
    }
    to {
      opacity: 1;
    }
  }
</style>

<style lang="scss" rel="stylesheet/scss" scoped>
  $color-blue: #2196F3;
  $space-width: 18px;
  .ms-tree-space {
    position: relative;
    top: 1px;
    display: inline-block;
    font-style: normal;
    font-weight: 400;
    line-height: 1;
    width: $space-width;
    height: 14px;
    &::before {
      content: ""
    }
  }

  .processContainer {
    width: 100%;
    height: 100%;
  }

  table td {
    line-height: 26px;
  }

  .tree-ctrl {
    position: relative;
    cursor: pointer;
    color: $color-blue;
    margin-left: -$space-width;
  }
</style>

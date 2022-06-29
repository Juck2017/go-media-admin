<template>
  <div class="chooseAuthDiv">
    <el-tree
      ref="tree"
      :data="formatData"
      show-checkbox
      node-key="id"
      :default-expand-all="true"
      :default-checked-keys="orgPrivilges"
      :props="defaultProps"
      @check-change="handleCheckChange"
    />
  </div>
</template>

<script>
export default {
  name: 'ChooseAuth',
  props: {
    data: {
      type: Array,
      default: () => []
    },
    orgPrivilges: {
      type: Array,
      default: () => []
    }
  },
  data() {
    return {
      defaultProps: {
        children: 'children',
        label: 'name'
      }
    }
  },
  computed: {
    formatData() {
      let tmp
      if (this.data.length > 0) {
        // 组装参数
        tmp = this.treeListUtil(this.data)
      } else {
        tmp = this.data
      }
      return tmp
    }
  },
  watch: {
    // 此处是切换不同的权限页面后进行初次权限设置
    orgPrivilges: function(val) {
      this.$refs.tree.setCheckedKeys(val)
      this.$store.dispatch('setCheckAuth', val)
    }
  },
  mounted() {
    // 初始化数据
    this.handleCheckChange()
  },
  methods: {
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
    handleCheckChange() {
      // 传入store
      this.$store.dispatch('setCheckAuth', [].concat(this.$refs.tree.getCheckedKeys(), this.$refs.tree.getHalfCheckedKeys()))
      // 触发给角色权限绑定的事件
      this.$emit('changeAuth')
    }
  }
}
</script>

<style rel="stylesheet/scss" lang="scss">

</style>

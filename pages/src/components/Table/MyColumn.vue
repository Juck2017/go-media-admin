<template>
  <el-table-column
    :prop="item.value"
    :label="item.text"
    align="center"
  >
    <template v-if="item.children">
      <my-column
        v-for="(item, index) in item.children"
        :key="index"
        :item="item"
      />
    </template>
    <template slot-scope="scope">
      <span
        v-if="!item.type"
        :style="item.style?item.style:''"
        :class="item.classFun?item.classFun(scope.row):''"
        @click="item.click?item.click(scope.row,$event):false"
        v-html="formatter(scope.row[item.value],item.formatter,scope.row,item.filter,item.filterParams)"
      />
      <img v-if="item.type === 'img'" :src="item.getSrc(scope.row)" :height="item.imgHeight">
      <div v-if="item.type === 'button'">
        <el-button
          v-for="(key, num) in item.list(scope.row)"
          :key="num"
          v-button-privilege="key.privilege"
          size="mini"
          :type="key.type"
          @click="key.click(scope,$event)"
        >{{ key.value }}
        </el-button>
      </div>
      <div v-if="item.type === 'iconButton'">
        <span v-if="item.value" style="margin-right:5px">{{ scope.row[item.value] }}</span>
        <span v-for="(key, num) in item.list(scope.row)" :key="num">
          <el-tooltip v-button-privilege="key.privilege" :content="key.value" placement="top">
            <i style="font-size: 25px; padding-right: 5px;" class="icon iconfont projectColor pointer" :class="key.class" @click="key.click(scope,$event)" />
          </el-tooltip>
          <span v-if="num < item.list(scope.row).length-1" class="projectColor" style="position: absolute;margin-left: -4px;font-weight: 700;">|</span>
        </span>
      </div>
      <div v-if="item.type === 'switch'">
        <el-switch
          v-model="scope.row[item.model]"
          :active-value="item.activeValue"
          :inactive-value="item.inactiveValue"
          :active-text="item.activeText"
          :inactive-text="item.inactiveText"
          @change="item.change(scope,scope.row[item.model])"
        />
      </div>
    </template>
  </el-table-column>
</template>

<script>
import buttonPrivilege from '@/directive/buttonPrivilege'

export default {
  name: 'MyColumn',
  directives: {
    buttonPrivilege
  },
  props: {
    item: {
      type: Object,
      default: () => {}
    }
  },
  methods: {
    formatter(value, formatter, row, filter, filterParams) {
      if (!formatter) {
        if (filter) {
          if (filterParams) {
            const tempArr = [value, ...filterParams]
            return filter(...tempArr)
          } else {
            return filter(value)
          }
        }
        return value
      } else {
        return formatter(row)
      }
    }
  }
}
</script>
<style scoped>
</style>


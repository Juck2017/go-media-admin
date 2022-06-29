<template>
  <div>
    <el-table
      ref="multipleTable"
      v-loading="listLoading"
      :header-cell-style="rowClass"
      :data="data"
      :height="tableHeight"
      :max-height="height"
      :highlight-current-row="true"
      :row-key="row_key"
      tooltip-effect="dark"
      :expand-row-keys="expandRowKeys"
      select-on-indeterminate
      border
      :stripe="stripe"
      size="small"
      element-loading-text="加载中..."
      :span-method="spanMethod"
      :row-class-name="rowClassName"
      empty-text="暂无数据"
      style="width: 100%"
      @expand-change="expandChange"
      @selection-change="selectChange"
      @row-click="rowClick"
      @row-dblclick="dblclick"
      @filter-change="filterChange"
      @sort-change="sortChange"
    >
      <!-- :span-method="spanMethod ? spanMethod: false" -->
      <!-- 显示多选框 -->
      <el-table-column
        v-if="select"
        type="selection"
        :selectable="selectable"
        width="55"
      />
      <!-- 显示行索引,从1开始 -->
      <el-table-column
        v-if="showIndex"
        type="index"
        align="center"
        label="序号"
        width="60"
      />
      <!-- 显示为一个可展开的按钮 -->
      <el-table-column
        v-if="expand"
        type="expand"
      >
        <template slot-scope="scope">
          <el-form label-position="left" inline class="demo-table-expand">
            <el-form-item
              v-for="(item, index) in expandList"
              :key="index"
              :style="{width: expandList[0].width}"
              :label="item.text"
            >
              <span v-html="formatter(scope.row[item.value],item.formatter,scope.row,item.filter,item.filterParams)" />
            </el-form-item>
          </el-form>
        </template>
      </el-table-column>
      <!-- filters: 数据过滤的选项 -->
      <el-table-column
        v-for="(item, index) in columns"
        :key="index"
        :label="item.text"
        :width="item.width"
        :show-overflow-tooltip="showOverflow"
        :filter-multiple="item.filterMultiple ? item.filterMultiple : false"
        :prop="item.value"
        :column-key="item.value"
        align="center"
        :sortable="item.sortable ? 'custom' : null"
        :filters="item.headFilters ? item.headFilters : null"
      >
        <!-- <my-column v-for="(child,i) in item.children" v-if="item.children" :key="i" :item="child" /> -->
        <template slot-scope="scope">
          <!-- <my-render v-if="item.render" :row="scope.row" :render="item.render" /> -->
          <span
            v-if="!item.type"
            :style="item.style?item.style:''"
            :class="item.classFun?item.classFun(scope.row):''"
            @click="item.click?item.click(scope.row,$event):false"
            v-html="formatter(scope.row[item.value],item.formatter,scope.row,item.filter,item.filterParams)"
          />
          <img
            v-if="item.type === 'img'"
            :src="item.getSrc(scope.row)"
            :height="item.imgHeight"
          >
          <div v-if="item.type === 'button'">
            <el-button
              v-for="(key, num) in item.list(scope.row)"
              :key="num"
              v-button-privilege="key.privilege"
              size="mini"
              :class="key.class"
              :type="key.type"
              @click="key.click(scope,$event)"
            >{{ key.value }}
            </el-button>
          </div>

          <div v-if="item.type=='iconPopover'" style="">
            <el-popover trigger="hover" placement="bottom-end">
              <!-- <div
                v-for="file in scope.row[item.value]"
                v-if="scope.row[item.value]&&scope.row[item.value].length!=0"
                :key="file.name"
                style="width:425px;font-size:16px;color:#0166DC;"
              >
                <div style="height:40px;line-height:40px">
                  <div>
                    <a
                      :download="file.name"
                      :href="file.fileName"
                      target="_blank"
                      style="display:inline-block;width:372px;height:30px;margin-right:30px;overflow: hidden;text-overflow:ellipsis;white-space: nowrap;text-decoration:underline"
                    >{{ file.name }}</a>
                  </div>
                </div>

              </div> -->
              <div v-if="!scope.row[item.value]||scope.row[item.value].length==0" style="text-align:center;">暂无数据</div>
              <div slot="reference" class="name-wrapper">
                <!-- <el-tag size="medium">aaaa</el-tag> -->
                <i v-if="item.iconfont&&item.iconfont!=''" class="iconfont" :class="item.iconfont" />
                <img v-if="!item.iconfont||item.iconfont==''" src="../../../static/img/contract/files.png" class="pointer">
              </div>
            </el-popover>
          </div>

          <div v-if="item.type=='reportView'" style="">
            <div
              v-for="file in scope.row[item.value]"
              :key="file.name"
            >
              <span style="font-size:14px;">{{ file.name }}</span>
              <!--                <a :href="file.viewUrl" title="预览" target="_blank" v-button-privilege="'B_zybyl'">-->
              <!--                  <i -->
              <!--                    class="iconfont icon-chakan1"-->
              <!--                    style="margin-left:20px;color:#FF973B;font-size:20px;"></i>-->
              <!--                </a>-->
              <a v-button-privilege="'B_zybxz'" :href="file.downloadUrl" title="下载" download target="_blank">
                <i
                  class="iconfont icon-download"
                  style="margin-left:20px;color:#0166DC;font-size:18px;"
                />
              </a>
              <i
                v-button-privilege="'B_zybsc'"
                class="iconfont icon-reduce-fill pointer"
                style="margin-left:20px;color:#FF3A56;font-size:20px;
                "
                @click="deleteFile(scope,file)"
              />
            </div>
          </div>

          <div v-if="item.type === 'textButton'" style="display:flex;justify-content: space-around;">
            <span
              v-for="(key, num) in item.list(scope.row)"
              :key="num"
              class="pointer"
              style="color:rgb(24,144,255);"
              @click="key.click(scope,$event)"
            >{{ key.value }}
            </span>
          </div>
          <!-- 展示图标按钮 -->
          <div v-if="item.type === 'iconButton'">
            <span
              v-if="item.value"
              style="margin-right:5px"
            >{{ scope.row[item.value] }}</span>
            <span v-for="(key, num) in item.list(scope.row)" :key="num">
              <el-tooltip
                v-button-privilege="key.privilege"
                :content="key.value"
                placement="top"
              >
                <i
                  style="font-size: 25px; padding-right: 5px;"
                  class="icon iconfont projectColor pointer"
                  :class="key.class"
                  @click="key.click(scope,$event)"
                />
              </el-tooltip>
              <!--<i v-button-privilege="key.privilege" style="font-size: 25px; padding-right: 5px;" class="icon iconfont projectColor pointer" :class="key.class" @click='key.click(scope,$event)'  ></i>-->
              <span
                v-show="num < item.list(scope.row).length-1"
                class="projectColor"
                style="position: absolute;margin-left: -4px;font-weight: 700;"
              >|</span>
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
          <div v-if="item.type === 'Tag'">
            <el-tag
              :key="item.text"
              :type="scope.row.status=='1'?'success':'info'"
              :effect="item.effect"
            >
              {{ item.getLabel(scope.row) }}
            </el-tag>
          </div>
        </template>
      </el-table-column>
    </el-table>
    <el-pagination
      v-if="showPagination && total>pageSize"
      ref="pagination"
      style="margin-top: 10px;"
      :current-page="pageNum"
      :page-size="pageSize"
      layout="total, prev, pager, next, jumper"
      :total="total"
      @current-change="currentChange"
    />
  </div>
</template>
<script>
import buttonPrivilege from '@/directive/buttonPrivilege'
// import myColumn from './MyColumn'
// import MyRender from './MyRender'

export default {
  name: 'TableList',
  directives: {
    buttonPrivilege
  },
  components: {
    // myColumn,
    // MyRender
  },
  props: {
    data: {
      type: Array,
      default: () => {
        return []
      }
    },
    height: {
      type: Number,
      default: () => {
        return
      }
    },
    columns: {
      type: Array,
      default: () => {
        return []
      }
    },
    styleObj: {
      type: Object,
      default: () => {
        return {}
      }
    },
    select: {
      type: Boolean,
      default: () => {
        return false
      }
    },
    isAccordion: {
      type: Boolean,
      default: () => {
        return false
      }
    },
    rowkey: {
      type: String,
      default: () => {
        return 'id'
      }
    },
    stripe: {
      type: Boolean,
      default: true
    },
    pageNum: { type: Number, default: 1 },
    pageSize: { type: Number, default: 10 },
    total: { type: Number, default: 0 },
    listLoading: { type: Boolean, default: false },
    showPagination: { type: Boolean, default: true },
    expand: { type: Boolean, default: false },
    showIndex: { type: Boolean, default: false },
    showOverflow: { type: Boolean, default: false },
    expandList: {
      type: Array,
      default: () => {
        return []
      }
    },
    headColor: { type: String, default: '' },
    spanMethod: {
      type: Function,
      default: () => function() {
      }
    },
    heightlightRows: {
      type: Array,
      default: () => {
        return []
      }
    },
    myCellStyle: { type: Function, default: null },
    tableHeight: {
      type: Number, default: () => {
        return
      }
    },
    rowClassName: {
      type: Function,
      default: () => function() {}
    }
  },
  data() {
    return {
      expandRowKeys: [],
      oldrowkey: null
    }
  },
  watch: {
    pageNum(val) {
      this.$refs.pagination.lastEmittedPage = val
    }
  },
  mounted() {
    // console.log(this.$refs.multipleTable)
  },
  methods: {
    selectable(row, index) {
      // row中无selectable默认复选框可用
      const { selectable = true } = row
      return selectable
    },
    filterChange(filters) {
      this.$emit('filterChange', filters)
    },
    sortChange(val) {
      // ascending 表示升序，descending 表示降序，null 表示还原为原始顺序
      this.$emit('sortChange', val)
    },
    rowClass({ row, rowIndex }) {
      // console.log(rowIndex)
      return 'background:' + this.headColor
    },
    selectChange(e) {
      this.$emit('selectChange', e)
    },
    rowClick(e, event, column) {
      this.$emit('rowClick', e, event, column)
    },
    clear() {
      this.$refs.multipleTable.clearSelection()
      return true
    },
    deleteFile(scope, file) {
      this.$emit('deleteFile', scope, file)
    },
    addSerch(rows) {
      rows.forEach(row => {
        this.$refs.multipleTable.toggleRowSelection(row)
      })
    },
    dblclick(e) { // 双击表格触发事件
      this.$emit('dbclick', e)
    },
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
    },
    highlight(row) {
      this.$refs.multipleTable.setCurrentRow(row)
    },
    currentChange(val) {
      this.$emit('currentChange', val)
    },
    // row 中有{row, column, rowIndex, columnIndex}
    // cellStyle(row) {
    //   if (this.myCellStyle) {
    //   }
    //     return this.myCellStyle(row)
    // },
    row_key(row) {
      return row[this.rowkey]
    },
    expandChange(e) {
      this.$emit('expandChange', e)
      if (this.isAccordion) {
        if (this.oldrowkey === e[this.rowkey] && this.expandRowKeys.length !== 0) {
          this.expandRowKeys = []
        } else {
          this.expandRowKeys = [e[this.rowkey]]
        }
      }
      this.oldrowkey = e[this.rowkey]
    }
  }
}
</script>
<!-- table测试组件 -->
<style rel="stylesheet/scss" lang="scss">
  .el-table--striped .el-table__body tr.el-table__row--striped td {
    background: #EAF3F8
  }

  .demo-table-expand {
    font-size: 0;
  }

  .demo-table-expand label {
    color: #99a9bf;
  }

  .demo-table-expand .el-form-item {
    margin-right: 0;
    margin-bottom: 0;
  }
  .el-table__row.alarm_column{
    background-color: #fffce1 !important;
    td{
      background-color: #fffce1 !important;
    }
  }
</style>

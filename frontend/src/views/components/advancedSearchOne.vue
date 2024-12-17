<template>
  <div class="advancedSearchOne">
    <div v-show="ownIndex != 0" style="border:1px solid #e1e1e1;position:relative;margin-bottom:20px;font-size:10px">
      <span style="position:absolute;top:-5px;left:-1px;background:#f4f5f9;font-weight: 600;">{{
        searchLevelData.category }}</span>
    </div>
    <div style="display:flex;height:auto">
      <div style="flex:1">
        <el-select v-model="searchLevelData.name" filterable placeholder="请选择" @change="searchFormchangeEdit"
          style="width:100%;">
          <el-option v-for="(item, indexs) in searchMoreArr" :key="indexs" :label="item.label" :value="item.name">
          </el-option>
        </el-select>
      </div>
      <div style="flex:1;margin-left:15px;">
        <el-select :disabled="!searchLevelData.name" size='mini' v-model="searchLevelData.type" placeholder="请选择"
          @change="searchFormEdit(searchLevelData.type)" style="width:100%;">
          <el-option v-for="(item, indexs) in searchLevelData.index" :key="indexs" :label="item.label"
            :value="item.value">
          </el-option>
        </el-select>
      </div>
      <div style="flex:1;display:flex;margin-top: -6px;margin-left:15px;width:220px">
        <el-form @submit.native.prevent v-if="Object.keys(searchLevelData.objContent).length != 0" :model="searchObjOne"
          ref="searchObjForm" label-width="0px" class="demo-ruleForm">
          <template v-if="searchLevelData.objContent.name != 'id'">
            <template v-if="searchLevelData.objContent.enum && searchLevelData.objContent.indexSearch != 'selectNull'">
              <el-form-item v-if="!(searchLevelData.objContent.type == 'object' && searchLevelData.objContent.list)">
                <el-select class="fanfan" :disabled="Object.keys(searchLevelData.objContent).length == 0" multiple
                  reserve-keyword placeholder="请选择" filterable clearable style="width:100%;"
                  v-model="searchLevelData.val">
                  <el-option v-for="(oitem, index) in searchLevelData.objContent.enum" :value="oitem.en"
                    :label="oitem.zh" :key="index"></el-option>
                </el-select>
              </el-form-item>
            </template>
            <template v-if="searchLevelData.objContent.indexSearch == 'selectNull'">
              <el-form-item>
                <el-select :disabled="Object.keys(searchLevelData.objContent).length == 0" reserve-keyword
                  placeholder="请选择" clearable style="width:100%;" v-model="searchLevelData.val">
                  <el-option :value="true" label="是"></el-option>
                  <el-option value="no" label="否"></el-option>
                </el-select>
              </el-form-item>
            </template>
            <template v-if="!searchLevelData.objContent.enum && searchLevelData.objContent.indexSearch != 'selectNull'">
              <template
                v-if="searchLevelData.objContent.DataList && searchLevelData.objContent.indexSearch == 'select'">
                <el-form-item v-if="!(searchLevelData.objContent.type == 'object' && searchLevelData.objContent.list)">
                  <el-select class="fanfan" :disabled="Object.keys(searchLevelData.objContent).length == 0" multiple
                    reserve-keyword placeholder="请选择" filterable clearable style="width:100%;"
                    :remote-method="getselectList" remote v-model="searchLevelData.val">
                    <el-option v-for="(oitem, index) in searchLevelData.objContent.DataList"
                      :value="oitem[searchLevelData.objContent.name.split('.')[searchLevelData.objContent.name.split('.').length - 1]]"
                      :label="oitem[searchLevelData.objContent.name.split('.')[searchLevelData.objContent.name.split('.').length - 1]]"
                      :key="index"></el-option>
                  </el-select>
                </el-form-item>
              </template>
              <template
                v-if="searchLevelData.objContent.TagsList && searchLevelData.objContent.indexSearch == 'select'">
                <el-form-item v-if="!(searchLevelData.objContent.type == 'object' && searchLevelData.objContent.list)">
                  <el-select :disabled="Object.keys(searchLevelData.objContent).length == 0" multiple reserve-keyword
                    placeholder="请选择" filterable clearable :remote-method="getselectList" remote style="width:100%;"
                    class="langClass" v-model="searchLevelData.val">
                    <el-option v-for="(oitem, index) in searchLevelData.objContent.TagsList" :value="oitem.value"
                      :label="oitem.value" :key="index"></el-option>
                  </el-select>
                </el-form-item>
              </template>
              <template v-if="searchLevelData.objContent.indexSearch == 'arr'">
                <el-form-item>
                  <el-input style="margin-top:6px;width:100%"
                    :disabled="Object.keys(searchLevelData.objContent).length == 0" type="textarea"
                    placeholder="换行输入多个值" :autosize="{ minRows: 3, maxRows: 20 }" resize="none"
                    v-model="searchLevelData.val"></el-input>
                </el-form-item>
              </template>
              <template
                v-if="searchLevelData.objContent.indexSearch == '_HGT' || searchLevelData.objContent.indexSearch == '_HLT'">
                <el-date-picker :disabled="Object.keys(searchLevelData.objContent).length == 0"
                  v-model="searchLevelData.val" type="datetime" placeholder="选择日期时间" align="right"
                  value-format="yyyy-MM-ddTHH:mm:ss+08:00" style="margin-top:6px;width:100%">
                </el-date-picker>
              </template>
              <template
                v-if="searchLevelData.objContent.indexSearch == '_DGT' || searchLevelData.objContent.indexSearch == '_DLT'">
                <el-date-picker v-model="searchLevelData.val" type="date" placeholder="选择日期" align="right"
                  value-format="yyyy-MM-dd" style="margin-top:6px;width:100%"
                  :disabled="Object.keys(searchLevelData.objContent).length == 0">
                </el-date-picker>
              </template>
              <template
                v-if="searchLevelData.objContent.indexSearch == '_GT' || searchLevelData.objContent.indexSearch == '_LT' || searchLevelData.objContent.indexSearch == '' || searchLevelData.objContent.indexSearch == '_REGEX'">
                <el-form-item style="width:100%;">
                  <el-input :disabled="Object.keys(searchLevelData.objContent).length == 0" placeholder="请输入"
                    style="width:100%;" clearable v-model.trim="searchLevelData.val"></el-input>
                </el-form-item>
              </template>

            </template>
          </template>

        </el-form>
        <el-select v-else disabled v-model="defaultData" placeholder="请选择" style="margin-top:6px;width:100%">
          <el-option :label="1" :value="1">
          </el-option>
        </el-select>
      </div>
      <div style="margin-top:10px;font-size:12px;color: #409EFF;cursor: pointer;">
        <i class="el-icon-delete" style="margin-left:20px;color:red" @click="delNewsearch"
          v-if="!(advSearch.length == 1 && (!advSearch.child || advSearch.child.length == 0))"></i>
       <!--  <span style="margin-left:10px" @click="addNewsearch('OR')" v-if="searchLevel != 3"><i
            class="el-icon-circle-plus-outline" style="margin-right:5px"></i>OR</span> -->
        <span style="margin-left:10px" @click="addNewsearch('AND')" v-if="searchLevel != 3"><i
            class="el-icon-circle-plus-outline" style="margin-right:5px"></i>AND</span>
        <span style="margin-left:10px" @click="addNewsearch()" v-if="searchLevel == 3"><i
            class="el-icon-circle-plus-outline"></i></span>
      </div>
    </div>
  </div>
</template>
<script>
import Http from "@/components/api/services";
export default {
  props: {
    searchMoreArr: {
      type: Array,
      default: []
    },
    tableHeaderArr: {
      type: Array,
      default: []
    },
    advSearch: {
      type: Array,
      default: []
    },
    searchLevelData: {
      type: Object,
      default: {}
    },
    searchLevel: {
      type: Number,
      default: null
    },
    ownIndex: {
      type: Number,
      default: null
    },
    parentIndex: {
      type: Number,
      default: null
    },
    oneIndex: {
      type: Number,
      default: null
    },
  },
  data() {
    return {
      searchObjOne: {},
      defaultData: '',
      searchMoreArrChange: [],
      searchLevelDataChange: {}
    }
  },
  watch: {

  },
  methods: {
    searchFormchangeEdit() {
      this.searchMoreArrChange = JSON.parse(JSON.stringify(this.searchMoreArr))
      this.searchLevelData.type = null
      this.searchLevelData.val = null
      this.searchLevelData.objContent = {}
      this.searchMoreArrChange.map((item, index) => {
        if (item.name == this.searchLevelData.name) {
          this.searchLevelData.index = item.children
          this.searchLevelData.type = item.children[0].value
          this.tableHeaderArr.map(key => {
            if (key.name == this.searchLevelData.name) {
              key.indexSearch = this.searchLevelData.type
              this.searchLevelData.objContent = JSON.parse(JSON.stringify(key))
            }
          })
        }
      })
      let firstStr = this.searchLevelData.name
      this.searchMoreArrChange.map(key => {
        if (key.name == firstStr) {
          var twoStr = key.children[0].value
        }
      })
    },
    searchFormEdit(item) {
      this.searchLevelData.objContent = {}
      this.searchLevelData.val = null
      this.tableHeaderArr.map(key => {
        if (key.name == this.searchLevelData.name) {
          key.indexSearch = item
          this.searchLevelData.objContent = JSON.parse(JSON.stringify(key))
        }
      })
    },
    addNewsearch(category) {
      this.$emit('addNewsearch', category, this.searchLevel, this.ownIndex, this.parentIndex, this.oneIndex, this.searchLevelData)
    },
    delNewsearch() {
      this.$emit('delNewsearch', this.searchLevel, this.ownIndex, this.parentIndex, this.oneIndex, this.searchLevelData)
    },
    getselectList(val) {
      if (val) {
        this.getSelList(this.searchLevelData.objContent, val)
      }
    },
    getSelList(obj, value) {
      let postData = {}
      let keyV = `${obj.name.split(".")[obj.name.split(".").length - 1]}`
      if (value && value != '') {
        postData = {
          schema: obj.schema,
          where: {
            [`${keyV}_REGEX:`]: value,
          },
        }
      } else {
        postData = {
          schema: obj.schema,
        }
      }
      Http.getQueryList(postData).then((response) => {
        if (response.data.data.data) {
          if (obj.DataList) {
            obj.DataList = response.data.data.data
          }
        }
      })
    },
  },
  mounted() { }
}
</script>
<style scoped>
.fanfan /deep/ .el-select__tags .el-tag.el-tag--info {
  max-width: 95% !important;
  overflow: hidden !important;
  text-overflow: ellipsis !important;
  position: relative !important;
  padding-right: 15px;
}

.fanfan /deep/ .el-tag.el-tag--info .el-tag__close {
  position: absolute;
  top: 5px;
  right: -2px;
}

.app-container
  .advancedSearchOne
  /deep/
  .el-select
  .el-input.is-disabled
  .el-input__inner {
  background: #ecf0f4 !important;
  color: #c0c4cc !important;
}

.advancedSearchOne /deep/ .el-form {
  width: 100% !important;
}
</style>

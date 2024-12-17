import { mapGetters } from "vuex";
import Http from "@/components/api/services";
import replaceQuery from "../utils/replaceQuery";
export const userMixin = {
  data() {
    return {
      users: [],
      checkIdcNameBQ: "",
      referenceShowArrSlist: [],
      referenceShowArrSlistObj: {},
      uploadObjs: {},
      ObjDataRowServer: {},
      searchObj: {},
      classWidth: "",
      poHide: false,
      searchMoreArr: [],
      cascaderMoreArr: [],
      searchObj1: {},
      searchFormData: {},
      searchObjTag: [],
      tableHeaderArr: [],
      referenceShowArrSlistArrStats: false,
      cascaderMorestr: "",
      searchMoreArrSecond: [],
      defaultData: "",
      cascaderMorestrSecond: null,
      advSearchIndex: null,
      enumAllArr: [],
      advSearch: [
        {
          name: "",
          type: "",
          val: "",
          category: "",
          index: [],
          child: [],
          objContent: {},
        },
      ],
      lastAdvSearch: [],
      batchBG: [],
    };
  },
  methods: {
    clearableRadioBtn(row, obj, result, cobj, iii, iiii, formValueT, key) {
      var value = row.paramList ? [] : "";
      formValueT[key] = value;
      this.tureOrFalseCoant(value, obj, result, cobj, iii, iiii);
    },
    // 搜索tag移除、
    handleCloseTag(item) {
      this.searchObj1[item.enName] = "";
      this.searchTableServer(true);
    },
    //高级搜索,点击菜单发生变化时
    searchForm(item) {
      this.searchFormData = {};
      this.tableHeaderArr.map((key) => {
        if (key.name == this.cascaderMorestr) {
          key.indexSearch = item;
          this.searchFormData = JSON.parse(JSON.stringify(key));
        }
      });
    },
    enterTagLast(index) {
      this.advSearchIndex = index;
      this.poHide = true;
      this.advSearch = this.lastAdvSearch[index];
    },
    closePopover() {
      this.classWidth = "smallwidth";
    },
    //高级搜索,默认选择第一个子菜单时
    searchFormchange() {
      this.cascaderMorestrSecond = null;
      this.searchFormData = {};
      this.searchMoreArr.map((item, index) => {
        if (item.name == this.cascaderMorestr) {
          this.searchMoreArrSecond = item.children;
        }
      });
      let firstStr = this.cascaderMorestr;
      this.searchMoreArr.map((key) => {
        if (key.name == firstStr) {
          twoStr = key.children[0].value;
        }
      });
    },
    searchShowStatus1() {
      this.advSearchIndex = null;
      this.poHide = true;
      this.advSearch = [
        {
          name: "",
          type: null,
          val: "",
          category: "",
          index: [],
          child: [],
          objContent: {},
        },
      ];
    },
    addNewsearch(type, level, ownIndex, parentIndex, oneIndex, obj) {
      if (level == 1) {
        if (this.advSearch[0].category == "") {
          this.advSearch[0].category = type;
          this.advSearch.push({
            name: "",
            type: "",
            val: "",
            category: type,
            index: [],
            child: [],
            objContent: {},
          });
        } else {
          if (this.advSearch.length == 1) {
            this.advSearch[0].category = type;
            this.advSearch.push({
              name: "",
              type: "",
              val: "",
              category: type,
              index: [],
              child: [],
              objContent: {},
            });
          } else {
            if (this.advSearch[ownIndex].category == type) {
              this.advSearch.push({
                name: "",
                type: "",
                val: "",
                category: type,
                index: [],
                child: [],
                objContent: {},
              });
            } else {
              let arrTemp = [];
              arrTemp.push({
                name: obj.name,
                type: obj.type,
                val: obj.val,
                category: type,
                index: obj.index,
                child: [],
                objContent: obj.objContent,
              });
              arrTemp.push({
                name: "",
                type: "",
                val: "",
                category: type,
                index: [],
                child: [],
                objContent: {},
              });
              this.advSearch[ownIndex].child = arrTemp;
            }
          }
        }
      }
      if (level == 2) {
        if (this.advSearch[parentIndex].child[ownIndex].category == type) {
          this.advSearch[parentIndex].child.push({
            name: "",
            type: "",
            val: "",
            category: type,
            index: [],
            child: [],
            objContent: {},
          });
        } else {
          let arrTemp = [];
          arrTemp.push({
            name: obj.name,
            type: obj.type,
            val: obj.val,
            category: type,
            index: obj.index,
            objContent: obj.objContent,
          });
          arrTemp.push({
            name: "",
            type: "",
            val: "",
            category: type,
            index: [],
            objContent: {},
          });
          this.advSearch[parentIndex].child[ownIndex].child = arrTemp;
        }
      }
      if (level == 3) {
        this.advSearch[oneIndex].child[parentIndex].child.push({
          name: "",
          type: "",
          val: "",
          category:
            this.advSearch[oneIndex].child[parentIndex].child[0].category,
          index: [],
          objContent: {},
        });
      }
      this.advSearch = JSON.parse(JSON.stringify(this.advSearch));
    },
    delNewsearch(level, ownIndex, parentIndex, oneIndex, obj) {
      let advSearchData = JSON.parse(JSON.stringify(this.advSearch));
      if (level == 1) {
        if (advSearchData.length == 1 && ownIndex == 0 && (!obj.child || !obj.child.length)) {
          return;
        }
        advSearchData.splice(ownIndex, 1);
        if (advSearchData.length == 1) {
          if (advSearchData[0].child && advSearchData[0].child.length) {
            advSearchData = advSearchData[0].child;
          }
        }
      }
      if (level == 2) {
        advSearchData[parentIndex].child.splice(ownIndex, 1);
        if (advSearchData[parentIndex].child.length == 1) {
          advSearchData[parentIndex] = {
            name: advSearchData[parentIndex].child[0].name,
            type: advSearchData[parentIndex].child[0].type,
            val: advSearchData[parentIndex].child[0].val,
            category: advSearchData[parentIndex].category,
            index: advSearchData[parentIndex].child[0].index,
            objContent: advSearchData[parentIndex].child[0].objContent,
            child: advSearchData[parentIndex].child[0].child
              ? advSearchData[parentIndex].child[0].child
              : [],
          };
        }
      }
      if (level == 3) {
        advSearchData[oneIndex].child[parentIndex].child.splice(ownIndex, 1);
        if (advSearchData[oneIndex].child[parentIndex].child.length == 1) {
          advSearchData[oneIndex].child[parentIndex] = {
            name: advSearchData[oneIndex].child[parentIndex].child[0].name,
            type: advSearchData[oneIndex].child[parentIndex].child[0].type,
            val: advSearchData[oneIndex].child[parentIndex].child[0].val,
            category: advSearchData[oneIndex].child[parentIndex].category,
            index: advSearchData[oneIndex].child[parentIndex].child[0].index,
            objContent: advSearchData[oneIndex].child[parentIndex].child[0].objContent,
            child: advSearchData[oneIndex].child[parentIndex].child[0].child
              ? advSearchData[oneIndex].child[parentIndex].child[0].child
              : [],
          };
        }
      }
      this.advSearch = JSON.parse(JSON.stringify(advSearchData));
      this.$forceUpdate();
    },
    searchTagChange(addvitem) {
      let eItemTxt = "";
      let searchObjTagss = [];
      if (addvitem.val) {
        this.tableHeaderArr.map((item) => {
          if (item.name == addvitem.name) {
            if (item.enum) {
              eItemTxt += "[";
              item.enum.map((eitem) => {
                if (Array.isArray(addvitem.val)) {
                  addvitem.val.map((kkey, kki) => {
                    if (eitem.en == kkey) {
                      eItemTxt += `${eitem.zh},`;
                    }
                  });
                } else {
                  if (eitem.en == addvitem.val) {
                    eItemTxt = eitem.zh;
                  }
                }
              });
              eItemTxt = eItemTxt.slice(0, eItemTxt.length - 1);
              eItemTxt += "]";
            } else {
              eItemTxt = addvitem.val;
            }
            searchObjTagss = {
              name: addvitem.objContent.label,
              value: eItemTxt,
              type: addvitem.type ? addvitem.type : "select",
              category: addvitem.category,
              child: addvitem.child,
            };
          }
        });
      }
      return searchObjTagss;
    },
    entrySreach(item, obj) {
      let oldname = item.enName.split("+")[0];
      let oldsecond = item.enName.split("+")[1] ? item.enName.split("+")[1] : "";
      this.cascaderMorestr = oldname;
      this.searchFormData = {};
      this.searchMoreArr.map((item, index) => {
        if (item.name == this.cascaderMorestr) {
          this.searchMoreArrSecond = item.children;
        }
      });
      let twoStr = "";
      let firstStr = this.cascaderMorestr;
      this.searchMoreArr.map((key) => {
        if (key.name == firstStr) {
          twoStr = key.children[0].value;
        }
      });
      this.cascaderMorestrSecond = oldsecond;
      this.tableHeaderArr.map((key) => {
        if (key.name == this.cascaderMorestr) {
          key.indexSearch = this.cascaderMorestrSecond;
          this.searchFormData = JSON.parse(JSON.stringify(key));
        }
      });
      this.poHide = true;
    },
    refAdvancedSearch(index) {
      this.resetSearch();
      this.searchObjTag.splice(index, 1);
      this.lastAdvSearch.splice(index, 1);
      this.searchTableServer(true);
    },
    resetSearch() {
      this.advSearch = [
        {
          name: "",
          type: null,
          val: "",
          category: "",
          index: [],
          child: [],
          objContent: {},
        },
      ];
    },
    refAdvancedSearchLast() {
      this.resetSearch();
      this.searchObjTag = [];
      this.lastAdvSearch = [];
      this.searchTableServer(true);
    },
    clickServerBatch(obj, row, infoReviewForm, infoReviewAttr) {
      this.listName = row.paramName;
      row.childAttrArr[0].map((mitem, mi) => {
        if (mitem.paramremainOperand == mitem.paramreference) {
          this.choseServerList(
            obj,
            mi,
            row.paramName,
            infoReviewForm,
            infoReviewAttr,
            mitem,
            mitem.paramName,
            0,
            mitem.paramTitle,
            true
          );
        }
      });
    },
    choseServerBtnBathjoin() {
      let serverArr = [];
      let paramremainOperandStr = this.infoReviewAttr[this.serverTypeParentIndex].childAttrArr[0][0].paramremainOperand; //Server
      let paramNameStr = "";
      if (!this.choseListServer.length) {
        this.$message({
          showClose: true,
          message: "请选择条目",
          type: "warning",
        });
        return;
      }
      this.infoReviewAttr[this.serverTypeParentIndex].childAttrArr[0].map((item, i) => {
        if (item.paramreference == paramremainOperandStr) {
          paramNameStr = item.paramName;
        }
      }
      );
      this.infoReviewForm[this.childRowTitle].map((item, i) => {
        serverArr.push(item[paramNameStr]);
      });
      this.choseListServer.map((citem, ci) => {
        if (serverArr.indexOf(citem.id) == -1) {
          let obj = {};
          let attrLineone = JSON.parse(
            JSON.stringify(
              this.infoReviewAttr[this.serverTypeParentIndex].childAttrArr[0]
            )
          );
          attrLineone.map((oitem, oi) => {
            obj[oitem.paramName] = oitem.paramList ? [] : "";
            if (oitem.paramName == "cloudAccount") {
              obj[oitem.paramName] = oitem.paramList ? [this.choseListServer[ci].account.id] : this.choseListServer[ci].account.id;
              oitem.paramDataList = JSON.parse(
                JSON.stringify([
                  {
                    name: this.choseListServer[ci].account.name,
                    value: this.choseListServer[ci].account.id,
                  }
                ])
              );
            }
            if (oitem.paramName == paramNameStr) {
              obj[oitem.paramName] = oitem.paramList ? [this.choseListServer[ci].id] : this.choseListServer[ci].id;
              oitem.paramDataList = JSON.parse(
                JSON.stringify([
                  {
                    name: this.choseListServer[ci].name,
                    value: this.choseListServer[ci].id,
                  }
                ])
              );
            }
          });
          if (this.infoReviewForm[this.childRowTitle][0][paramNameStr]) {
            this.infoReviewForm[this.childRowTitle].push(obj);
            this.infoReviewAttr[this.serverTypeParentIndex].childAttrArr.push(attrLineone);
          } else {
            this.infoReviewForm[this.childRowTitle] = [obj];
            this.infoReviewAttr[this.serverTypeParentIndex].childAttrArr = [attrLineone];
          }
        } else {
          this.$message({
            showClose: true,
            message: "已添加过该对象",
            type: "warning",
          });
        }
      });
      this.infoReviewAttr[this.serverTypeParentIndex].childAttrArr.map(
        (item, ii) => {
          item.map((jyitem, jykey) => {
            if (jyitem.paramRequired) {
              if (
                this.infoReviewForm[this.childRowTitle][ii][jyitem.paramName]
              ) {
                this.infoReviewAttr[this.serverTypeParentIndex].childAttrArr[
                  ii
                ][jykey].colorRed = "1";
              }
            }
          });
        }
      );
      this.choseServerLog = false;
    },
    emitChild(obj) {
      this.infoReviewForm[obj.d][obj.aiii.erindex][obj.e] = JSON.parse(JSON.stringify(obj.a));
      this.infoReviewAttr[Number(obj.waiIndex)].childAttrArr[obj.aiii.erindex][obj.aiii.sani].childAttrArrThird = JSON.parse(JSON.stringify(obj.b));
      this.infoReviewForm = JSON.parse(JSON.stringify(this.infoReviewForm));
      this.infoReviewAttr = JSON.parse(JSON.stringify(this.infoReviewAttr));
    },
    emitChildfocusChildGLArr(obj) {
      this.infoReviewAttr[Number(obj.waiIndex)].childAttrArr[obj.aiii.erindex][obj.aiii.sani].childAttrArrThird[obj.pcindex][obj.nrindex].paramDataList = obj.dataList;
      this.infoReviewAttr = JSON.parse(JSON.stringify(this.infoReviewAttr));
    },
    emitChildvalidating(obj) {
      this.infoReviewAttr[Number(obj.waiIndex)].childAttrArr[obj.aiii.erindex][obj.aiii.sani].childAttrArrThird = obj.attr;
      this.infoReviewAttr = JSON.parse(JSON.stringify(this.infoReviewAttr));
    },
    emitChildcopy(obj) {
      this.infoReviewForm[obj.d][obj.aiii.erindex][obj.e] = obj.a;
      this.infoReviewAttr[Number(obj.waiIndex)].childAttrArr[obj.aiii.erindex][obj.aiii.sani].childAttrArrThird = obj.b;
      this.infoReviewForm = JSON.parse(JSON.stringify(this.infoReviewForm));
      this.infoReviewAttr = JSON.parse(JSON.stringify(this.infoReviewAttr));
    },
    emitChilddelOne(obj) {
      this.infoReviewForm[obj.d][obj.aiii.erindex][obj.e] = obj.a;
      this.infoReviewForm = JSON.parse(JSON.stringify(this.infoReviewForm));
    },
    emitChildSureSer(obj) {
      this.infoReviewAttr[Number(obj.waiIndex)].childAttrArr[obj.aiii.erindex][obj.aiii.sani].childAttrArrThird = obj.b;
      this.infoReviewAttr = JSON.parse(JSON.stringify(this.infoReviewAttr));
    },
    checkShow(name, i) {
      this.infoReviewForm[name + "ChildObj"][i]["_isShowStatus"] = !this.infoReviewForm[name + "ChildObj"][i]["_isShowStatus"];
      this.infoReviewForm = JSON.parse(JSON.stringify(this.infoReviewForm));
      this.infoReviewAttr = JSON.parse(JSON.stringify(this.infoReviewAttr));
    },
    linkJumpTo(obj, row) {
      let link = obj.paramreferenceLink.replace("{$id}", row.id);
      window.open(`${link}`);
    },
    // 如果是服务器初始化则有批量导入按钮，点击导入按钮
    batchInBtn() {
      if (this.batchTextarea.replace(/\s/g, "")) {
        var batchInarr = [];
        var batchInNewarr = [];
        var batchInNewarrLast = [];
        var batchInNewarrObjLast = [];
        batchInarr = this.batchTextarea.split("\n");
        batchInarr.map((k, i) => {
          if (k.replace(/\s/g, "")) {
            batchInNewarr.push(k.replace(/\s/g, ""));
          }
        });
      }
      let arrArr = this.infoReviewAttr[this.batchIndex].childAttrArr[0];
      batchInNewarr.map((k, i) => {
        let hArr = k.split(",");
        batchInNewarrObjLast.push({});
        arrArr.map((k1, i1) => {
          let reg = new RegExp(k1.paramPatterns);
          if (k1.paramRequired && hArr[i1] == "") {
            k1.colorRed = "5";
          } else {
            k1.colorRed = "1";
          }
          if ((k1.paramMinLength != null && k1.paramMinLength != "undefind" && k1.paramMaxLength != null && k1.paramMaxLength != "undefind") || (k1.paramMinmum != null && k1.paramMinmum != "undefind" && k1.paramMaxmum != null && k1.paramMaxmum != "undefind")) {
            if (k1.paramType == "string" && (hArr[i1].length < k1.paramMinLength || hArr[i1].length > k1.paramMaxLength)) {
              k1.colorRed = "2";
            } else if (k1.paramType == "number" && (hArr[i1] < k1.paramMinmum || hArr[i1] > k1.paramMaxmum)) {
              k1.colorRed = "4";
            } else {
              k1.colorRed = "1";
            }
          }
          if (k1.paramPatterns && k1.colorRed == "1") {
            if (!reg.test(hArr[i1])) {
              k1.colorRed = "3";
            } else {
              k1.colorRed = "1";
            }
          }
          batchInNewarrObjLast[i][k1.paramName] = hArr[i1] ? hArr[i1] : "";
        });
        batchInNewarrLast.push(arrArr);
      });
      this.infoReviewForm[this.infoReviewAttr[this.batchIndex].paramName + "ChildObj"] = batchInNewarrObjLast;
      this.infoReviewAttr[this.batchIndex].childAttrArr = batchInNewarrLast;
      this.infoReviewForm = JSON.parse(JSON.stringify(this.infoReviewForm));
      this.infoReviewAttr = JSON.parse(JSON.stringify(this.infoReviewAttr));
      this.batchIn = false;
    },
    // 验证状态
    validating(value, row, minLen, maxLen, pattern, type, formDatathis) {
      let index = row.$index;
      let reg = new RegExp(pattern);
      if (value == "" && this.infoReviewAttr[index].paramRequired) {
        if (this.newTaskName != "资源模板新建" && this.newTaskName != "资源模板编辑" && this.newTaskName != "弹性资源模板新建" && this.newTaskName != "弹性资源模板编辑" && this.newTaskName != "弹性资源模板改配") {
          this.infoReviewAttr[index].colorRed = "5";
          this.infoReviewAttr[index].customDescription = "不可为空";
        }
      } else {
        if (value) {
          this.infoReviewAttr[index].colorRed = "1";
        }
        if (minLen != null && minLen != "undefind") {
          if (this.infoReviewAttr[index].colorRed == "1") {
            if (type == "string" && value.length < minLen) {
              this.infoReviewAttr[index].colorRed = "2";
              this.infoReviewAttr[index].customDescription =
                "最小长度" + minLen + "个字符";
            } else if (type == "number" && value < minLen) {
              this.infoReviewAttr[index].colorRed = "4";
              this.infoReviewAttr[index].customDescription =
                "数值不得小于" + minLen;
            } else {
              this.infoReviewAttr[index].colorRed = "1";
            }
          }
        }
        if (maxLen != null && maxLen != "undefind") {
          if (this.infoReviewAttr[index].colorRed == "1") {
            if (type == "string" && value.length > maxLen) {
              this.infoReviewAttr[index].colorRed = "2";
              this.infoReviewAttr[index].customDescription =
                "最大长度" + maxLen + "个字符";
            } else if (type == "number" && value > maxLen) {
              this.infoReviewAttr[index].colorRed = "4";
              this.infoReviewAttr[index].customDescription =
                "数值不得大于" + maxLen;
            } else {
              this.infoReviewAttr[index].colorRed = "1";
            }
          }
        }
        if (pattern && this.infoReviewAttr[index].colorRed == "1") {
          if (!reg.test(value)) {
            this.infoReviewAttr[index].colorRed = "3";
            this.infoReviewAttr[index].customDescription = "填写格式有误";
          } else {
            this.infoReviewAttr[index].colorRed = "1";
          }
        }
        if (!this.infoReviewAttr[index].paramRequired) {
          if (value == "") {
            this.infoReviewAttr[index].colorRed = "1";
          }
        }
      }
      if (formDatathis) {
        formDatathis.map((item, i) => {
          if (pattern && this.infoReviewAttr[index].colorRed == "1") {
            if (!reg.test(item)) {
              this.infoReviewAttr[index].colorRed = "3";
              this.infoReviewAttr[index].customDescription = "填写格式有误";
            } else {
              this.infoReviewAttr[index].colorRed = "1";
            }
          }
        });
      }
      this.infoReviewAttr = JSON.parse(JSON.stringify(this.infoReviewAttr));
    },
    // 验证状态子表单
    validatingchild(value, row, cindex, ccindex, minLen, maxLen, pattern, type, formDatathis) {
      let pIndex = row.$index;
      let reg = new RegExp(pattern);
      if (value == "" && this.infoReviewAttr[pIndex].childAttrArr[cindex][ccindex].paramRequired) {
        if (this.newTaskName != "资源模板新建" && this.newTaskName != "资源模板编辑" && this.newTaskName != "弹性资源模板新建" && this.newTaskName != "弹性资源模板编辑" && this.newTaskName != "弹性资源模板改配") {
          this.infoReviewAttr[pIndex].childAttrArr[cindex][ccindex].colorRed = "5";
          this.infoReviewAttr[pIndex].childAttrArr[cindex][ccindex].customDescription = "不可为空";
        }
      } else {
        if (value) {
          this.infoReviewAttr[pIndex].childAttrArr[cindex][ccindex].colorRed = "1";
        }
        if (minLen != null && minLen != "undefind") {
          if (type == "string" && value.length < minLen) {
            this.infoReviewAttr[pIndex].childAttrArr[cindex][ccindex].colorRed = "2";
            this.infoReviewAttr[pIndex].childAttrArr[cindex][ccindex].customDescription = "最小长度" + minLen + "个字符";
          } else if (type == "number" && value < minLen) {
            this.infoReviewAttr[pIndex].childAttrArr[cindex][ccindex].colorRed = "4";
            this.infoReviewAttr[pIndex].childAttrArr[cindex][ccindex].customDescription = "数值不得小于" + minLen;
          } else {
            this.infoReviewAttr[pIndex].childAttrArr[cindex][ccindex].colorRed = "1";
          }
        }
        if (maxLen != null && maxLen != "undefind" && this.infoReviewAttr[pIndex].childAttrArr[cindex][ccindex].colorRed == "1") {
          if (type == "string" && value.length > maxLen) {
            this.infoReviewAttr[pIndex].childAttrArr[cindex][ccindex].colorRed = "4";
            this.infoReviewAttr[pIndex].childAttrArr[cindex][ccindex].customDescription = "最大长度" + maxLen + "个字符";
          } else if (type == "number" && value > maxLen) {
            this.infoReviewAttr[pIndex].childAttrArr[cindex][ccindex].colorRed = "4";
            this.infoReviewAttr[pIndex].childAttrArr[cindex][ccindex].customDescription = "数值不得大于" + maxLen;
          } else {
            this.infoReviewAttr[pIndex].childAttrArr[cindex][ccindex].colorRed = "1";
          }
        }
        if (pattern && this.infoReviewAttr[pIndex].childAttrArr[cindex][ccindex].colorRed == "1") {
          if (!reg.test(value)) {
            this.infoReviewAttr[pIndex].childAttrArr[cindex][ccindex].colorRed = "3";
            this.infoReviewAttr[pIndex].childAttrArr[cindex][ccindex].customDescription = "填写格式有误";
          } else {
            this.infoReviewAttr[pIndex].childAttrArr[cindex][ccindex].colorRed = "1";
          }
        }
        if (!this.infoReviewAttr[pIndex].childAttrArr[cindex][ccindex].paramRequired) {
          if (value == "") {
            this.infoReviewAttr[pIndex].childAttrArr[cindex][ccindex].colorRed = "1";
          }
        }
      }
      if (formDatathis) {
        formDatathis.map((item, i) => {
          if (pattern && this.infoReviewAttr[pIndex].childAttrArr[cindex][ccindex].colorRed == "1") {
            if (!reg.test(item)) {
              this.infoReviewAttr[pIndex].childAttrArr[cindex][ccindex].colorRed = "3";
              this.infoReviewAttr[pIndex].childAttrArr[cindex][ccindex].customDescription = "填写格式有误";
            } else {
              this.infoReviewAttr[pIndex].childAttrArr[cindex][ccindex].colorRed = "1";
            }
          }
        });
      }
      this.infoReviewAttr = JSON.parse(JSON.stringify(this.infoReviewAttr));
    },
    // 下拉搜索
    filterMethod(val, obj, formData) {
      this.getstr = val;
      this.getMetaSchemasListval(obj.row.paramreference, obj.row, val, formData);
    },
    // 下拉菜单
    getMetaSchemasListvalFocus(val, obj, formData) {
      this.getstr = val;
      this.getMetaSchemasListval(obj.row.paramreference, obj.row, val, formData);
    },
    focusSelectDefaultValue(obj, name, formData, index) {
      var objData = obj.row;
      if (objData.paramreference) {
        if (!objData.paramreferenceQuery) {
          objData.paramreferenceQuery = "{}"
        }
      }else{
        return false
      }
      var postData = {
        schema: objData.paramreference,
        page_size:100,
        page_num:1,
        where:{}
      }
      if(objData.paramreferenceQuery){
        var queryStr = replaceQuery(objData.paramreferenceQuery,this.loginUserName,formData,this)
        console.log(queryStr)
        if (queryStr.indexOf("ISvalue_null") != -1) {
          return false;
        }
        queryStr = queryStr.replace(/\"\[/g, "[").replace(/\]\"/g, "]");
        var queryStrObj = JSON.parse(queryStr)
        postData.where = queryStrObj
      }

      Http.getQueryList(postData).then(
      (response) => {
        var arrList = []
        response.data.data.data.map((k,i)=>{
          arrList.push({
            name: k.name,
            value: k.id
          })
        })
        if(arrList.length&&!this.lookType){       
            if(objData.paramDefault){
                if(objData.paramDefault.indexOf("$")!=-1){
                    this.infoReviewForm[name] = arrList[Number(objData.paramDefault.split("$")[1])].value
                }
            }
        }
        let newArrD = JSON.parse(JSON.stringify(arrList));
        this.infoReviewAttr[index].paramDataList = newArrD;
        this.infoReviewAttr = JSON.parse(JSON.stringify(this.infoReviewAttr));
        this.infoReviewForm = JSON.parse(JSON.stringify(this.infoReviewForm));
      })
    },
    // 下拉聚焦
    focusSelectValue(obj, name, formData, index, df) {
      var objData = obj.row;
      if (objData.paramreference) {
        if (!objData.paramreferenceQuery) {
          objData.paramreferenceQuery = "{}"
        }
      }else{
        return false
      }
      var postData = {
        schema: objData.paramreference,
        page_size:100,
        page_num:1,
        where:{}
      }
      if(objData.paramreferenceQuery){
        var queryStr = replaceQuery(objData.paramreferenceQuery,this.loginUserName,formData,this)
        console.log(queryStr)
        if (queryStr.indexOf("ISvalue_null") != -1) {
          return false;
        }
        queryStr = queryStr.replace(/\"\[/g, "[").replace(/\]\"/g, "]");
        var queryStrObj = JSON.parse(queryStr)
        postData.where = queryStrObj
      }

      Http.getQueryList(postData).then(
      (response) => {
        var arrList = []
        response.data.data.data.map((k,i)=>{
          if (obj.paramreferenceJoinTitle) {
            this.idChangeNameObj[k.id] = k[obj.paramreferenceJoinTitle];
          }
          arrList.push({
            name: k.name,
            value: k.id
          })
        })
        let newArrD = JSON.parse(JSON.stringify(arrList));
        let lastArrD = newArrD;
        if (index) {
          this.infoReviewAttr[index].paramDataList = lastArrD;
          if (df) {
            if (objData.paramDefault && objData.paramDefault.indexOf("$") != -1) {
              formData[name] = lastArrD[Number(objData.paramDefault.split("$")[1])].value;
            }
          }
        } else {
          objData.paramDataList = lastArrD;
          if (objData.paramDefault && objData.paramDefault.indexOf("$") != -1) {
            formData[name] = lastArrD[Number(objData.paramDefault.split("$")[1])].value;
          }
        }
        this.infoReviewAttr = JSON.parse(JSON.stringify(this.infoReviewAttr));
        this.infoReviewAttr.map((k, i) => {
          if (k.paramreferenceQuery && k.paramreferenceQuery.indexOf(`$${name}`) != -1) {
            this.focusSelectValue({ row: this.infoReviewAttr[i] }, this.infoReviewAttr[i].paramName, this.infoReviewForm, i, "df");
          }
        });
      })
    },
    // 获取模型列表
    getMetaSchemasListval(name, obj, val, formData, pindex, index, cindex) {
      if (!name) {
        return;
      }
      if (name) {
        if (!obj.paramreferenceQuery) {
          obj.paramreferenceQuery = "{}"
        }
      }else{
        return false
      }
      var postData = {
        schema: name,
        page_size:100,
        page_num:1,
        where:{}
      }
      if(obj.paramreferenceQuery){
        var queryStr = replaceQuery(obj.paramreferenceQuery,this.loginUserName,formData,this)
        console.log(queryStr)
        if (queryStr.indexOf("ISvalue_null") != -1) {
          return false;
        }
        queryStr = queryStr.replace(/\"\[/g, "[").replace(/\]\"/g, "]");
        var queryStrObj = JSON.parse(queryStr)
        postData.where = queryStrObj
      }
      if (val) {
        if(!postData.where.and){
          postData.where.and = []
        }
        postData.where.and.push({name_REGEX : val})
      }
     

      Http.getQueryList(postData).then(
      (response) => {
        var arrList = []
        response.data.data.data.map((k,i)=>{
          arrList.push({
            name: k.name,
            value: k.id
          })
        })
        let newArrD = JSON.parse(JSON.stringify(arrList));
        if (pindex !== undefined) {
          this.infoReviewAttr[pindex].childAttrArr[index][cindex].paramDataList = newArrD;
          obj.paramDataList = newArrD;
        } else {
          obj.paramDataList = newArrD;
        }
      })
    },
    getMetaSchemasList(name, obj, i, i1, yiindex, erindex, twoIndex, pcindex, i33, diyiIn) {
      if (!name) {
        return;
      }
      if (name) {
        if (!obj.paramreferenceQuery) {
          obj.paramreferenceQuery = "{}"
        }
      }else{
        return false
      }
      if(i33) {
        var formData = this.infoReviewForm[this.infoReviewAttr[yiindex].paramName+"ChildObj"][erindex]
      }else{
        var formData = this.infoReviewForm
      }
      var postData = {
        schema: name,
        page_size:100,
        page_num:1,
        where:{}
      }
      if(obj.paramreferenceQuery){
        var queryStr = replaceQuery(obj.paramreferenceQuery,this.loginUserName,formData,this)
        console.log(queryStr)
        if (queryStr.indexOf("ISvalue_null") != -1) {
          return false;
        }
        queryStr = queryStr.replace(/\"\[/g, "[").replace(/\]\"/g, "]");
        var queryStrObj = JSON.parse(queryStr)
        postData.where = queryStrObj
      }
      obj.paramDataList = []
      Http.getQueryList(postData).then(
      (response) => {
        obj.paramDataList = []
        response.data.data.data.map((k,i)=>{
          if (obj.paramreferenceJoinTitle) {
            this.idChangeNameObj[k.id] = k[obj.paramreferenceJoinTitle];
          }
          obj.paramDataList.push({
            name: k.name,
            value: k.id
          })
        })
        if (i !== undefined && i !== null) {
          this.infoReviewAttr[yiindex].childAttrArr[i][i1] = obj;
        }
        if (erindex !== undefined && erindex !== null) {
          if (this.infoReviewAttr[yiindex]["childAttrArr"][erindex][twoIndex]["childAttrArrThird"][pcindex]) {
            this.infoReviewAttr[yiindex]["childAttrArr"][erindex][twoIndex]["childAttrArrThird"][pcindex][i33] = obj;
          }
        }
        if (diyiIn !== undefined && diyiIn !== null) {
          this.infoReviewAttr[diyiIn] = obj;
        }
        this.infoReviewAttr = JSON.parse(JSON.stringify(this.infoReviewAttr));
      })
    },
    resetForm(formName) {
      if (this.$refs[formName]) {
        this.$refs[formName].resetFields();
      }
    },
    // 是否关联必填 如果关联，会有一个选择  是就展示 不是隐藏
    tureOrFalseCoant(value, obj, result, cobj, iii, iiii, status) {
      if (cobj) {
        let ii = obj.$index;
        if (value == "") {
          if (this.newTaskName != "资源模板新建" && this.newTaskName != "资源模板编辑" && this.newTaskName != "弹性资源模板新建" && this.newTaskName != "弹性资源模板编辑" && this.newTaskName != "弹性资源模板改配") {
            if (this.infoReviewAttr[ii].childAttrArr[iii][iiii].paramRequired) {
              this.infoReviewAttr[ii].childAttrArr[iii][iiii].colorRed = "5";
              this.infoReviewAttr[ii].childAttrArr[iii][iiii].customDescription = "不可为空";
            } else {
              this.infoReviewAttr[ii].childAttrArr[iii][iiii].colorRed = "1";
            }
          }
        } else {
          this.infoReviewAttr[ii].childAttrArr[iii][iiii].colorRed = "1";
        }
        for (var Ai = iiii + 1; Ai < this.infoReviewAttr[ii].childAttrArr[iii].length; Ai++) {
          if (this.infoReviewAttr[ii].childAttrArr[iii][Ai].paramreferenceQuery) {
            if (this.infoReviewAttr[ii].childAttrArr[iii][Ai].paramreferenceQuery.indexOf(`$${this.infoReviewAttr[ii].childAttrArr[iii][iiii].paramName}`) != -1) {
              this.focusChildSelect(obj, Ai, obj.row.paramName, this.infoReviewForm, this.infoReviewAttr[ii].childAttrArr[iii][Ai], iii);
              if (this.infoReviewAttr[ii].childAttrArr[iii][Ai].paramList) {
                this.infoReviewForm[obj.row.paramName + "ChildObj"][iii][this.infoReviewAttr[ii].childAttrArr[iii][Ai].paramName] = [];
              } else {
                this.infoReviewForm[obj.row.paramName + "ChildObj"][iii][this.infoReviewAttr[ii].childAttrArr[iii][Ai].paramName] = "";
              }
              this.infoReviewAttr[ii].childAttrArr[iii][Ai].paramDataList = [];
            }
          }
        }
      } else {
        let index = obj.$index;
        if (value == "") {
          if (this.newTaskName != "资源模板新建" && this.newTaskName != "资源模板编辑" && this.newTaskName != "弹性资源模板新建" && this.newTaskName != "弹性资源模板编辑" && this.newTaskName != "弹性资源模板改配") {
            if (this.infoReviewAttr[index].paramRequired) {
              this.infoReviewAttr[index].colorRed = "5";
              this.infoReviewAttr[index].customDescription = "不可为空";
            } else {
              this.infoReviewAttr[index].colorRed = "1";
            }
          }
        } else {
          this.infoReviewAttr[index].colorRed = "1";
        }
        for (var Ai = index + 1; Ai < this.infoReviewAttr.length; Ai++) {
          if (this.infoReviewAttr[Ai].paramreferenceQuery) {
            if (this.infoReviewAttr[Ai].paramreferenceQuery.indexOf(`$${obj.row.paramName}`) != -1) {
              this.focusSelectValue({ row: this.infoReviewAttr[Ai] }, this.infoReviewAttr[Ai].paramName, this.infoReviewForm, Ai);
              this.infoReviewForm[this.infoReviewAttr[Ai].paramName] = this.infoReviewAttr[Ai].paramList || this.infoReviewAttr[Ai].paramStyle == "upload" ? [] : "";
              this.infoReviewAttr[Ai].paramDataList = [];
            }
          }
        }
      }
      this.infoReviewAttr.map((arItem) => {
        if (arItem.paramShemas) {
          let shemasObject = JSON.parse(arItem.paramShemas);
          let kgIndex = 0;
          for (let meItem in shemasObject) {
            if (Object.prototype.toString.call(shemasObject[meItem]).indexOf("Array") != -1) {
              if (Object.prototype.toString.call(this.infoReviewForm[meItem]).indexOf("Array") != -1) {
                let kg = false;
                this.infoReviewForm[meItem].map((k5, i5) => {
                  if (shemasObject[meItem].indexOf(k5) != -1) {
                    kg = true;
                  }
                });
                if (kg) {
                  kgIndex += 1;
                }
              } else {
                if (shemasObject[meItem].indexOf(this.infoReviewForm[meItem]) != -1) {
                  kgIndex += 1;
                }
              }
            } else {
              if (Object.prototype.toString.call(this.infoReviewForm[meItem]).indexOf("Array") != -1) {
                if (this.infoReviewForm[meItem].indexOf(shemasObject[meItem]) != -1) {
                  kgIndex += 1;
                }
              } else {
                if (this.infoReviewForm[meItem]) {
                  if (shemasObject[meItem] == this.infoReviewForm[meItem]) {
                    kgIndex += 1;
                  }
                }
              }
            }
          }

          if (kgIndex == Object.keys(shemasObject).length) {
            if (arItem.paramavailableCondition == "role:admin") {
              if (this.users.indexOf(this.usernameAll) != -1) {
                arItem.thisShowIf = "1";
              } else {
                arItem.thisShowIf = "2";
              }
            } else if (
              arItem.paramavailableCondition &&
              arItem.paramavailableCondition.indexOf("groupRole:") != -1
            ) {
              var arrGroup = arItem.paramavailableCondition.split("groupRole:")[1].split(",");
              arItem.thisShowIf = "2";
              arrGroup.map((gk1, gi1) => {
                if (this.manageUserGroup.indexOf(gk1) != -1) {
                  arItem.thisShowIf = "1";
                }
              });
            } else {
              arItem.thisShowIf = "1";
            }
          } else {
            arItem.thisShowIf = "2";
          }
          if (arItem.thisShowIf == "2") {
            if (arItem.paramList) {
              if (arItem.paramreference != "FormTemplate") {
                this.infoReviewForm[arItem.paramName] = [];
              }
            } else {
              if (arItem.paramDefault && arItem.paramDefault.indexOf("$") == -1) {
                this.infoReviewForm[arItem.paramName] = arItem.paramDefault;
              } else {
                if (arItem.paramreference != "FormTemplate") {
                  if (arItem.paramList || arItem.paramStyle == "upload") {
                    this.infoReviewForm[arItem.paramName] = [];
                  } else {
                    this.infoReviewForm[arItem.paramName] = "";
                  }
                }
              }
            }
          }
        }
        if (arItem.childAttrArr) {
          arItem.childAttrArr.map((achItem, achItemi) => {
            achItem.map((aaaItem, aaaItemi) => {
              if (achItemi == iii) {
                if (aaaItem.paramShemas) {
                  let shemasObjectch = JSON.parse(aaaItem.paramShemas);
                  var bgif = "2";
                  var a = 0;
                  for (let meItemch in shemasObjectch) {
                    if (cobj.paramName == meItemch) {
                      if (Object.prototype.toString.call(shemasObjectch[meItemch]).indexOf("Array") != -1) {
                        shemasObjectch[meItemch].map((k222, i222) => {
                          if (Object.prototype.toString.call(result).indexOf("Array") != -1) {
                            if (result.indexOf(k222) != -1) {
                              a += 1;
                            }
                          } else {
                            if (result == k222) {
                              a += 1;
                            }
                          }
                        });
                      } else {
                        if (shemasObjectch[meItemch] == result) {
                          a += 1;
                        }
                      }
                    } else {
                      if (this.infoReviewForm[arItem.paramName + "ChildObj"][achItemi][meItemch]) {
                        if (Object.prototype.toString.call(this.infoReviewForm[arItem.paramName + "ChildObj"][achItemi][meItemch]).indexOf("Array") != -1) {
                          this.infoReviewForm[arItem.paramName + "ChildObj"][achItemi][meItemch].map((k333, i333) => {
                            if (Object.prototype.toString.call(shemasObjectch[meItemch]).indexOf("Array") != -1) {
                              if (shemasObjectch[meItemch].indexOf(k333) != -1) {
                                a += 1;
                              }
                            } else {
                              if (shemasObjectch[meItemch] == k333) {
                                a += 1;
                              }
                            }
                          });
                        } else {
                          if (Object.prototype.toString.call(shemasObjectch[meItemch]).indexOf("Array") != -1) {
                            if (shemasObjectch[meItemch].indexOf(this.infoReviewForm[arItem.paramName + "ChildObj"][achItemi][meItemch]) != -1) {
                              a += 1;
                            }
                          } else {
                            if (shemasObjectch[meItemch] == this.infoReviewForm[arItem.paramName + "ChildObj"][achItemi][meItemch]) {
                              a += 1;
                            }
                          }
                        }
                      }
                    }
                  }
                  if (a == Object.keys(shemasObjectch).length) {
                    bgif = "1";
                  }
                  if (aaaItem.paramavailableCondition == "role:admin") {
                    if (this.users.indexOf(this.usernameAll) != -1) {
                      aaaItem.thisShowIf = bgif;
                    } else {
                      aaaItem.thisShowIf = "2";
                    }
                  } else if (aaaItem.paramavailableCondition && aaaItem.paramavailableCondition.indexOf("groupRole:") != -1) {
                    var arrGroup = aaaItem.paramavailableCondition.split("groupRole:")[1].split(",");
                    aaaItem.thisShowIf = "2";
                    arrGroup.map((gk1, gi1) => {
                      if (this.manageUserGroup.indexOf(gk1) != -1) {
                        aaaItem.thisShowIf = bgif;
                      }
                    });
                  } else {
                    aaaItem.thisShowIf = bgif;
                  }
                  if (bgif == 1) {
                    if (aaaItem.paramType == "number") {
                      this.infoReviewForm[arItem.paramName + "ChildObj"][achItemi][aaaItem.paramName] = this.infoReviewForm[arItem.paramName + "ChildObj"][achItemi][aaaItem.paramName] ? this.infoReviewForm[arItem.paramName + "ChildObj"][achItemi][aaaItem.paramName] : Number(aaaItem.paramDefault);
                    }
                  } else {
                    if (aaaItem.paramList) {
                      this.infoReviewForm[arItem.paramName + "ChildObj"][achItemi][aaaItem.paramName] = [];
                    } else {
                      this.infoReviewForm[arItem.paramName + "ChildObj"][achItemi][aaaItem.paramName] = "";
                    }
                  }
                }
              }
            });
          });
        }
      });
      this.infoReviewForm = JSON.parse(JSON.stringify(this.infoReviewForm));
      this.infoReviewAttr = JSON.parse(JSON.stringify(this.infoReviewAttr));
    },
    // 选择子表单弹框
    choseServerList(obj, index, name, formData, forAtt, ShowTodisplayObj, title, bodyIndex, serverTitle, type) {
      console.log(111)
      if (type) {
        this.serverLogStatus = true;
      } else {
        this.serverLogStatus = false;
      }
      var formData1 = this.infoReviewForm;
      this.searchMoreArr = [];
      this.cascaderMoreArr = [];
      this.searchObjTag = [];
      this.searchObj = {};
      this.searchObj1 = {};
      this.ObjDataRowServer = {};
      this.searchFormServe = "";
      this.serverTitle = serverTitle;
      this.searchFormServe = "";
      this.resetVal();
      this.serverTypeParentIndex = obj.$index;
      this.serverTypeIndex = bodyIndex;
      //  如果是弹框选择服务器，找到要定义数据位置
      this.childRowTitle = name + "ChildObj";
      this.childRowTitle1 = title;
      var objData = obj.row.childAttrArr[bodyIndex][index];
      this.mainOperandName = objData.paramremainOperand;
      this.mainOperandvsrefrename = objData.paramreference;
      if (!objData.paramreferenceQuery) {
        objData.paramreferenceQuery = "{}"
      }
      this.searchpageNum = 1;
      var queryobj  = {}
      // query替换 
      var queryTempStr = objData.paramreferenceQuery;
      if(queryTempStr){
        var queryStr = replaceQuery(queryTempStr,this.loginUserName,formData,this)
        console.log(queryStr)
        if (queryStr.indexOf("ISvalue_null") != -1) {
          return false;
        }
        queryStr = queryStr.replace(/\"\[/g, "[").replace(/\]\"/g, "]");
        queryobj = JSON.parse(queryStr);
      }
      this.choseServerLog = true;
      this.searchHostTable = [];
      this.searchtotalLen = 0;
      this.choseListServer = [];
      this.referenceShowArrHeader = [];
      let objTypeList = {};
      this.referenceShowArrSlistObj = {};
      
      if (ShowTodisplayObj.paramreferenceDisplay) {
        objData.paramDataList = [];
        let referenceShowArrHeader11 = [];
        let paramreferenceShowNew = [];
        this.paramreferenceShowNew1 = [];
        let referenceShowArr12 = [];
        let firstName = "";
        let newDispalyObj = {};
        let displayShow = ShowTodisplayObj.paramreferenceDisplay ? JSON.parse(ShowTodisplayObj.paramreferenceDisplay) : { properties: [] };
        displayShow.properties.map((ditem, dii) => {
          if (ditem.name != "oid") {
            paramreferenceShowNew.push(ditem.name);
          } else {
            paramreferenceShowNew.push("id");
          }
          if (ditem.index) {
            this.paramreferenceShowNew1.push(ditem.name);
          }
          newDispalyObj[ditem.name] = ditem.title;
          if (ditem.display) {
            if (firstName == "" || !firstName) {
              firstName = ditem.name;
            }
            if (ditem.name.indexOf(".") != -1) {
              let valTemp = "";
              let arrTemp = ditem.name.split(".");
              let valTemp1 = "";
              let valTempB = "";
              let valTempC = "";
              if (arrTemp.length == 2) {
                arrTemp.map((kk, ii) => {
                  valTemp1 += `{${kk}}`;
                  if (ii == 0) {
                    valTemp = kk;
                  } else {
                    valTemp1 = `${valTemp}{${kk}}`;
                  }
                });
              } else if (arrTemp.length == 3) {
                arrTemp.map((kk, ii) => {
                  valTemp1 += `{${kk}}`;
                  if (ii == 0) {
                    valTemp = kk;
                  } else if (ii == 1) {
                    valTempB = `${valTemp}{${kk}`;
                  } else {
                    valTemp1 = `${valTempB}{${kk}}`;
                  }
                });
              } else if (arrTemp.length == 4) {
                arrTemp.map((kk, ii) => {
                  valTemp1 += `{${kk}}`;
                  if (ii == 0) {
                    valTemp = kk;
                  } else if (ii == 1) {
                    valTempB = `${valTemp}{${kk}`;
                  } else if (i1 == 2) {
                    valTempC = `${valTempB}{${kk}`;
                  } else {
                    valTemp1 = `${valTempB}{{${kk}}}}`;
                  }
                });
              }
              referenceShowArr12.push(valTemp1);
              objTypeList[valTemp1] = ditem;
            }
            if (ditem.name.indexOf(".") == -1) {
              referenceShowArr12.push(ditem.name);
              objTypeList[ditem.name] = ditem;
            }
          }
        });
        this.referenceShowArrSlist = referenceShowArr12;
        this.referenceShowArrSlistObj = objTypeList;
        this.paramreferenceShowNew2 = paramreferenceShowNew;
        paramreferenceShowNew.map((k, i) => {
          referenceShowArrHeader11.push({
            title: k,
            value: k,
          });
        });
        for (let rekey in displayShow.properties) {
          referenceShowArrHeader11.map((reItem) => {
            if (
              displayShow.properties[rekey].name == reItem.value &&
              displayShow.properties[rekey].display
            ) {
              this.referenceShowArrHeader.push({
                title: displayShow.properties[rekey].title,
                value: displayShow.properties[rekey].name,
              });
            }
          });
        }
        this.referenceShowArrHeader = JSON.parse(
          JSON.stringify(this.referenceShowArrHeader)
        );
        this.theadArrShow[bodyIndex] = this.referenceShowArrHeader;
        this.tableHeaderArr = displayShow.properties;
        this.tableHeaderArr.map((k, i) => {
          k.title = k.title.replace(/\./g, "");
          this.tagAllArr = [];
          if (k.schema) {
            this.getSelList(k);
          }
          if (k.name != "id" && k.default) {
            this.checkboxTH.push(k.title);
          }
          if (k.name && k.name != "id" && !k.hide && k.display) {
            this.searchMoreArr.push(k);
          }
          if (k.name && k.display) {
            if (k.tagsArr) {
              let str = "";
              k.tagsArr.map((item) => {
                str = k.name + "." + item;
                this.searchObj[str] = "";
              });
            } else {
              this.searchObj[k.name] = "";
            }
          }
        });
        this.searchMoreArr.map((item) => {
          item.label = item.title;
          item.value = item.name;
          if (item.enum) {
            item.children = [
              {
                value: "select",
                label: "选择",
              },
            ];
          } else {
            if (item.type == "number" || item.type == "integer") {
              item.children = [
                {
                  value: "",
                  label: "精确",
                },
                {
                  value: "_GT",
                  label: "大于",
                },
                {
                  value: "_LT",
                  label: "小于",
                },
              ];
            } else if (item.type == "string") {
              let arrTemp = [];

              if (item.index) {
                if (item.DataList || item.TagsList) {
                  arrTemp.push({
                    value: "select",
                    label: "选择",
                  });
                }
                if (item.index.indexOf("hash") != -1) {
                  arrTemp.push(
                    {
                      value: "",
                      label: "精确",
                    },
                    {
                      value: "arrTemp",
                      label: "批量",
                    }
                  );
                }
                if (item.index.indexOf("regexp") != -1) {
                  arrTemp.push({
                    value: "_REGEX",
                    label: "模糊",
                  });
                }
                if (item.index.indexOf("hour") != -1) {
                  arrTemp.push(
                    {
                      value: "_HGT",
                      label: "大于",
                    },
                    {
                      value: "_HLT",
                      label: "小于",
                    }
                  );
                }
                if (item.index.indexOf("day") != -1) {
                  arrTemp.push(
                    {
                      value: "_DGT",
                      label: "大于",
                    },
                    {
                      value: "_DLT",
                      label: "小于",
                    }
                  );
                }
                item.children = arrTemp;
              } else {
                if (item.DataList || item.TagsList) {
                  item.children = [
                    {
                      value: "select",
                      label: "选择",
                    },
                  ];
                } else if (!item.DataList && !item.TagsList) {
                  item.children = [
                    {
                      value: "",
                      label: "精确",
                    },
                  ];
                }
              }
            } else {
              item.children = [
                {
                  value: "",
                  label: "精确",
                },
              ];
            }
          }
          item.children.push({
            value: "selectNull",
            label: "空值",
          });
        });
        this.searchObj = JSON.parse(JSON.stringify(this.searchObj));
      }
      if (ShowTodisplayObj.paramreferenceDisplay) {
        let displayShow = ShowTodisplayObj.paramreferenceDisplay ? JSON.parse(ShowTodisplayObj.paramreferenceDisplay) : { properties: [] };
        displayShow.properties.map((ditem, dii) => {
          if (ditem.name == "cid") {
            this.showCidSearch = true;
          }
        });
      }
      this.queryObjectSearch = queryobj;
      this.nameWhereSearchServe = `${objData.paramreference}Where`;
      this.nameWhereSearchServeOption = `${objData.paramreference}Option`;
      this.queryNameSearchServe = objData.paramreference;
      this.ObjDataRowServer = objData;
      //   获取服务器列表接口方法
      this.searchTableServer(true);
    },
    // 添加列表弹框
    choseList(obj, name, formData, forAtt, ShowTodisplayObj, title) {
      console.log(222)
      this.choseListServer = [];
      this.resetVal();
      //   r如果是弹框选择服务器，找到要定义数据位置
      var objData = obj.row;
      console.log(222,objData)
      this.searchMoreArr = [];
      this.cascaderMoreArr = [];
      this.searchObjTag = [];
      this.searchObj = {};
      this.searchObj1 = {};
      this.ObjDataRowServer = {};
      this.searchFormServe = "";
      this.listTitle = objData.paramTitle;
      this.listName = objData.paramName;
      if (!objData.paramreferenceQuery) {
        objData.paramreferenceQuery = "{}"
      }
      this.searchpageNum = 1;
      var queryTempStr = objData.paramreferenceQuery;
      var queryobj = {}
      if(queryTempStr){
        var queryStr = replaceQuery(queryTempStr,this.loginUserName,formData,this)
        console.log(queryStr)
        if (queryStr.indexOf("ISvalue_null") != -1) {
          return false;
        }
        queryStr = queryStr.replace(/\"\[/g, "[").replace(/\]\"/g, "]");
        queryobj = JSON.parse(queryStr);
      }
      if (!title) {
        this.choseListStatus = true;
      }
      this.referenceShowListArrHeader = [];
      this.searchHostTable = [];
      this.searchtotalLen = 0;
      let objTypeList = {};
      this.referenceShowArrSlistObj = {};
      if (ShowTodisplayObj.paramreferenceDisplay) {
        objData.paramDataList = [];
        let referenceShowArrHeader11 = [];
        let paramreferenceShowNew = [];
        this.paramreferenceShowNew1 = [];
        let referenceShowArr12 = [];
        let firstName = "";
        let newDispalyObj = {};
        let displayShow = ShowTodisplayObj.paramreferenceDisplay ? JSON.parse(ShowTodisplayObj.paramreferenceDisplay) : { properties: [] };
        displayShow.properties.map((ditem, dii) => {
          if (ditem.name != "oid") {
            paramreferenceShowNew.push(ditem.name);
          } else {
            paramreferenceShowNew.push("id");
          }
          if (ditem.index) {
            this.paramreferenceShowNew1.push(ditem.name);
          }
          newDispalyObj[ditem.name] = ditem.title;
          if (ditem.display) {
            if (firstName == "" || !firstName) {
              firstName = ditem.name;
            }
            if (ditem.name.indexOf(".") != -1) {
              let valTemp = "";
              let arrTemp = ditem.name.split(".");
              let valTemp1 = "";
              let valTempB = "";
              let valTempC = "";
              if (arrTemp.length == 2) {
                arrTemp.map((kk, ii) => {
                  valTemp1 += `{${kk}}`;
                  if (ii == 0) {
                    valTemp = kk;
                  } else {
                    valTemp1 = `${valTemp}{${kk}}`;
                  }
                });
              } else if (arrTemp.length == 3) {
                arrTemp.map((kk, ii) => {
                  valTemp1 += `{${kk}}`;
                  if (ii == 0) {
                    valTemp = kk;
                  } else if (ii == 1) {
                    valTempB = `${valTemp}{${kk}`;
                  } else {
                    valTemp1 = `${valTempB}{${kk}}`;
                  }
                });
              } else if (arrTemp.length == 4) {
                arrTemp.map((kk, ii) => {
                  valTemp1 += `{${kk}}`;
                  if (ii == 0) {
                    valTemp = kk;
                  } else if (ii == 1) {
                    valTempB = `${valTemp}{${kk}`;
                  } else if (i1 == 2) {
                    valTempC = `${valTempB}{${kk}`;
                  } else {
                    valTemp1 = `${valTempB}{{${kk}}}}`;
                  }
                });
              }
              referenceShowArr12.push(valTemp1);
              objTypeList[valTemp1] = ditem;
            }
            if (ditem.name.indexOf(".") == -1) {
              referenceShowArr12.push(ditem.name);
              objTypeList[ditem.name] = ditem;
            }
          }
        });
        this.referenceShowArrSlist = referenceShowArr12;
        this.referenceShowArrSlistObj = objTypeList;
        this.paramreferenceShowNew2 = paramreferenceShowNew;
        paramreferenceShowNew.map((k, i) => {
          referenceShowArrHeader11.push({
            title: k,
            value: k,
          });
        });
        for (let rekey in displayShow.properties) {
          referenceShowArrHeader11.map((reItem) => {
            if (displayShow.properties[rekey].name == reItem.value && displayShow.properties[rekey].display) {
              this.referenceShowListArrHeader.push({
                title: displayShow.properties[rekey].title,
                value: displayShow.properties[rekey].name,
              });
            }
          });
        }
        this.referenceShowListArrHeader = JSON.parse(
          JSON.stringify(this.referenceShowListArrHeader)
        );
        this.listObj[objData.paramName] = this.referenceShowListArrHeader;
        this.tableHeaderArr = displayShow.properties;
        this.tableHeaderArr.map((k, i) => {
          k.title = k.title.replace(/\./g, "");
          this.tagAllArr = [];
          if (k.schema) {
            this.getSelList(k);
          }
          if (k.name != "id" && k.default) {
            this.checkboxTH.push(k.title);
          }
          if (k.name && k.name != "id" && !k.hide) {
            this.searchMoreArr.push(k);
          }
          if (k.name) {
            if (k.tagsArr) {
              let str = "";
              k.tagsArr.map((item) => {
                str = k.name + "." + item;
                this.searchObj[str] = "";
              });
            } else {
              this.searchObj[k.name] = "";
            }
          }
          if (k.enum) {
            k.enum.map((euitem, eni) => {
              let valueDataForm = `${k.name}:::${euitem.zh}`;
              this.enumAllArr[valueDataForm] = euitem.en;
            });
          }
        });
        this.searchMoreArr.map((item) => {
          item.label = item.title;
          item.value = item.name;
          if (item.enum) {
            item.children = [
              {
                value: "select",
                label: "选择",
              },
            ];
          } else {
            if (item.type == "number" || item.type == "integer") {
              item.children = [
                {
                  value: "",
                  label: "精确",
                },
                {
                  value: "_GT",
                  label: "大于",
                },
                {
                  value: "_LT",
                  label: "小于",
                },
              ];
            } else if (item.type == "string") {
              let arrTemp = [];

              if (item.index) {
                if (item.DataList || item.TagsList) {
                  arrTemp.push({
                    value: "select",
                    label: "选择",
                  });
                }
                if (item.index.indexOf("hash") != -1) {
                  arrTemp.push(
                    {
                      value: "",
                      label: "精确",
                    },
                    {
                      value: "arrTemp",
                      label: "批量",
                    }
                  );
                }
                if (item.index.indexOf("regexp") != -1) {
                  arrTemp.push({
                    value: "_REGEX",
                    label: "模糊",
                  });
                }
                if (item.index.indexOf("hour") != -1) {
                  arrTemp.push(
                    {
                      value: "_HGT",
                      label: "大于",
                    },
                    {
                      value: "_HLT",
                      label: "小于",
                    }
                  );
                }
                if (item.index.indexOf("day") != -1) {
                  arrTemp.push(
                    {
                      value: "_DGT",
                      label: "大于",
                    },
                    {
                      value: "_DLT",
                      label: "小于",
                    }
                  );
                }
                item.children = arrTemp;
              } else {
                if (item.DataList || item.TagsList) {
                  item.children = [
                    {
                      value: "select",
                      label: "选择",
                    },
                  ];
                } else if (!item.DataList && !item.TagsList) {
                  item.children = [
                    {
                      value: "",
                      label: "精确",
                    },
                  ];
                }
              }
            } else {
              item.children = [
                {
                  value: "",
                  label: "精确",
                },
              ];
            }
          }
          item.children.push({
            value: "selectNull",
            label: "空值",
          });
        });
        this.searchObj = JSON.parse(JSON.stringify(this.searchObj));
      }
      if (ShowTodisplayObj.paramreferenceDisplay) {
        let displayShow = ShowTodisplayObj.paramreferenceDisplay ? JSON.parse(ShowTodisplayObj.paramreferenceDisplay) : { properties: [] };
        displayShow.properties.map((ditem, dii) => {
          if (ditem.name == "cid") {
            this.showCidSearch = true;
          }
        });
      } 
      this.queryObjectSearch = queryobj;
      this.nameWhereSearchServe = `${objData.paramreference}Where`;
      this.nameWhereSearchServeOption = `${objData.paramreference}Option`;
      this.queryNameSearchServe = objData.paramreference;
      this.ObjDataRowServer = objData;
      //   获取服务器列表接口方法
      this.searchTableServer(true);
    },
    getSelList(obj, value) {
      let postData = {};
      let keyV = `${obj.name.split(".")[obj.name.split(".").length - 1]}`;
      if (value && value != '') {
        postData = {
          schema: obj.schema,
          where: {
            [`${keyV}_REGEX:`]: value,
          }
        }
      } else {
        postData = {
          schema: obj.schema,
        }
      }
      Http.getQueryList(postData).then((response) => {
        if (response.data.data.data) {
          if (obj.DataList) {
            obj.DataList = response.data.data.data;
          }
        }
      }
      ).catch((err) => {
        console.log(err);
      })
    },
    // 添加服务器
    searchTableServer(type) {
      if(type){
        this.searchpageNum = 1
      }

      var postData = {
        schema: this.queryNameSearchServe,
        page_size:this.searchpagesize,
        page_num:this.searchpageNum,
        where:{}
      }
      if(this.queryObjectSearch){
        postData.where = this.queryObjectSearch
      }
      if(!postData.where.and){
        postData.where.and = []
      }
      if(!postData.where.or){
        postData.where.or = {}
      }

      if(this.searchFormServe.trim() != ''){
        var disTitle = this.ObjDataRowServer.paramreferenceDisplay
        if(disTitle){
          var or = {}
          var dataPro = JSON.parse(this.ObjDataRowServer.paramreferenceDisplay).properties
          dataPro.map(k=>{
            if(k.display && k.tpye!='number' && k.type !='integer'){
              var searKey1 = k.name.split('.').reverse()
              var str11 = ''
              searKey1.map((k1, i1) => {
                if (i1 == 0) {
                   str11 = `{"${k1}_REGEX":"${this.searchFormServe.trim()}"}`
                } else {
                  str11 = `{"${k1}":${str11}}`
                }
              })
              var getKet = Object.keys(JSON.parse(str11))[0]
              or[getKet] = JSON.parse(str11)[getKet]
            }
          })
          postData.where.or = or
        }else{
          postData.where = {
            "name_REGEX":this.searchFormServe.trim()
          }
        }
      }
      if (this.advSearch[0].val || this.advSearch[0].category) {
        if (this.advSearchIndex || this.advSearchIndex === 0) {
          this.lastAdvSearch[this.advSearchIndex] = JSON.parse(
            JSON.stringify(this.advSearch)
          );
        } else {
          this.lastAdvSearch.push(JSON.parse(JSON.stringify(this.advSearch)));
        }
        var searchObjTagLast = JSON.parse(JSON.stringify(this.lastAdvSearch));
        searchObjTagLast.map((lastItem, lastI) => {
          lastItem.map((advItem, advI) => {
            if (advItem.child && advItem.child.length) {
              advItem.child.map((advItem2, adv2) => {
                if (advItem2.child && advItem2.child.length) {
                  advItem2.child.map((advItem3, adv3) => {
                    if (!advItem3.child || !advItem3.child.length) {
                      lastItem[advI].child[adv2].child[adv3] =
                        this.searchTagChange(advItem3);
                    }
                  });
                } else {
                  lastItem[advI].child[adv2] = this.searchTagChange(advItem2);
                }
              });
            } else {
              lastItem[advI] = this.searchTagChange(advItem);
            }
          });
        });
        this.searchObjTag = searchObjTagLast;
      }
      var and = []
      if (this.lastAdvSearch.length) {
        this.lastAdvSearch.map((lastItem, lastI) => {
          if (lastItem[0].category) {
            let firstType = lastItem[0].category;
            var arr1 = {
              [firstType]: [],
            };
            lastItem.map((advItem, advI) => {
              if (advItem.child && advItem.child.length) {
                let secondType = advItem.child[0].category;
                var arr2 = {
                  [secondType]: [],
                };
                advItem.child.map((advItem2, adv2) => {
                  if (advItem2.child && advItem2.child.length) {
                    let thirdType = advItem2.child[0].category;
                    var arr3 = {
                      [thirdType]: [],
                    };
                    advItem2.child.map((advItem3, adv3) => {
                      if (!advItem3.child || !advItem3.child.length) {
                        if (!advItem3.val) {
                          return;
                        }
                        if (Object.prototype.toString.call(advItem3.val).indexOf("Array") != -1) {
                          if (advItem3.val.length == 0) {
                            return;
                          }
                        }
                        arr3[thirdType].push(
                          this.advSearchHandle(advItem3, name, this.tableHeaderArr)
                        );
                      }
                    });
                    arr2[secondType].push(arr3);
                  } else {
                    if (!advItem2.val) {
                      return;
                    }
                    if (Object.prototype.toString.call(advItem2.val).indexOf("Array") != -1) {
                      if (advItem2.val.length == 0) {
                        return;
                      }
                    }
                    arr2[secondType].push(
                      this.advSearchHandle(advItem2, name, this.tableHeaderArr)
                    );
                  }
                });
                arr1[firstType].push(arr2);
              } else {
                if (!advItem.val) {
                  return;
                }
                if (Object.prototype.toString.call(advItem.val).indexOf("Array") != -1) {
                  if (advItem.val.length == 0) {
                    return;
                  }
                }
                arr1[firstType].push(
                  this.advSearchHandle(advItem, name, this.tableHeaderArr)
                );
              }
            });
            and.push(arr1);
          } else if (!lastItem[0].category && lastItem[0].val) {
            if (Object.prototype.toString.call(lastItem[0].val).indexOf("Array") != -1) {
              if (lastItem[0].val.length == 0) {
                return;
              }
            }
            and.push(this.advSearchHandle(lastItem[0], name, this.tableHeaderArr));
          }
        });
      }
      if(Object.keys(and).length){
        var arrAnd = []
        and.map((k1, i1) => {
          if(k1.AND){
            arrAnd.push(...k1.AND)
          }else{
            arrAnd.push(k1)
          }
        })
        postData.where.and.push(...arrAnd)
      }
      
      this.searchHostTable = []
      Http.getQueryList(postData).then(
      (response) => {
        this.searchHostTable = response.data.data.data;
        this.searchtotalLen = response.data.data.total;
      }).catch((error) => {
        this.searchtotalLen = 0
      })

      this.poHide = false;
      this.advSearch = [
        {
          name: "",
          type: null,
          val: "",
          category: "",
          index: [],
          child: [],
          objContent: {},
        },
      ];
    },
    // 确定选择服务器按钮
    choseServerBtn() {
      if (!this.infoReviewForm[this.childRowTitle][this.serverTypeIndex][this.childRowTitle1]) {
        this.infoReviewForm[this.childRowTitle][this.serverTypeIndex][this.childRowTitle1] = [];
      }
      if (!this.choseListServer.length) {
        this.$message({
          showClose: true,
          message: "请选择条目",
          type: "warning",
        });
        return;
      }
      var arrTemp = [];
      this.infoReviewForm[this.childRowTitle].map((k, i) => {
        if (k[this.childRowTitle1] && Object.prototype.toString.call(k[this.childRowTitle1]).indexOf("Array") != -1) {
          k[this.childRowTitle1].map((k2, i2) => {
            if (k2.oid) {
              arrTemp.push(k2.oid);
            } else {
              arrTemp.push(k2.id);
            }
          });
        }
      });
      this.choseListServer.map((k, i) => {
        //服务器名称对比 一样加校验  不一样不加
        if (this.mainOperandName == this.mainOperandvsrefrename) {
          if (arrTemp.indexOf(k.id) == -1) {
            this.infoReviewForm[this.childRowTitle][this.serverTypeIndex][this.childRowTitle1].push(k);
          } else {
            this.$message({
              showClose: true,
              message: "已添加过该对象",
              type: "warning",
            });
          }
        } else {
          this.infoReviewForm[this.childRowTitle][this.serverTypeIndex][this.childRowTitle1].push(k);
        }
      });
      if (this.infoReviewForm[this.childRowTitle][this.serverTypeIndex][this.childRowTitle1].length) {
        this.infoReviewAttr[this.serverTypeParentIndex].childAttrArr[this.serverTypeIndex].map((k, i) => {
          if (k.paramName == this.childRowTitle1) {
            k.colorRed = "1";
          }
        });
      }
      this.infoReviewForm = JSON.parse(JSON.stringify(this.infoReviewForm));
      this.choseServerLog = false;
    },
    choseListBtn() {
      if (!this.infoReviewForm[this.listName]) {
        this.infoReviewForm[this.listName] = [];
      }
      if (!this.choseListServer.length) {
        this.$message({
          showClose: true,
          message: "请选择条目",
          type: "warning",
        });
        return;
      }
      var arrTemp = [];
      this.infoReviewForm[this.listName].map((k, i) => {
        arrTemp.push(k.id);
      });

      this.choseListServer.map((k, i) => {
        if (arrTemp.indexOf(k.id) == -1) {
          this.infoReviewForm[this.listName].push(k);
        } else {
          this.$message({
            showClose: true,
            message: "已添加过该对象",
            type: "warning",
          });
        }
      });
      if (this.infoReviewForm[this.listName].length) {
        this.infoReviewAttr.map((k, i) => {
          if (k.paramName == this.listName) {
            this.infoReviewAttr[i].colorRed = "1";
          }
        });
      }
      this.infoReviewForm = JSON.parse(JSON.stringify(this.infoReviewForm));
      this.choseListStatus = false;
    },
    // 选择服务器
    handleSelectionChangeSearth(val) {
      this.choseListServer = val;
    },
    // 子表单下拉搜索
    focusChildSelect(obj, index, name, formData, ShowTodisplayObj, bodyIndex, status, resetStatus) {
      var objData = obj.row.childAttrArr[bodyIndex][index];
      if (objData.paramreference) {
        objData.paramDataList = [];
        if (!objData.paramreferenceQuery) {
          objData.paramreferenceQuery = "{}"
        }
      }else{
        return false
      }
      if(objData.paramreference=='account'){
        var showName = "alias"
      }else{
        var showName = "name"
      }
      var postData = {
        schema: objData.paramreference,
        page_size:100,
        page_num:1,
        where:{}
      }
      if(objData.paramreferenceQuery){
        var queryStr = replaceQuery(objData.paramreferenceQuery,this.loginUserName,formData,this)
        console.log(queryStr)
        if (queryStr.indexOf("ISvalue_null") != -1) {
          return false;
        }
        queryStr = queryStr.replace(/\"\[/g, "[").replace(/\]\"/g, "]");
        var queryStrObj = JSON.parse(queryStr)
        postData.where = queryStrObj
      }
      objData.paramDataList = []
      Http.getQueryList(postData).then(
      (response) => {
        objData.paramDataList = []
        response.data.data.data.map((k,i)=>{
          if (objData.paramreferenceJoinTitle) {
            this.idChangeNameObj[k.id] = k[objData.paramreferenceJoinTitle];
          }
          objData.paramDataList.push({
            name: k[showName],
            value: k.id,
          });
        })
        if ((bodyIndex !== undefined || bodyIndex !== null) && !resetStatus) {
          this.infoReviewAttr.map((k, i) => {
            if (k.paramName == name) {
              k.childAttrArr[bodyIndex][index].paramDataList = objData.paramDataList;
              if (!this.lookType) {
                if (k.childAttrArr[bodyIndex][index].paramDefault) {
                  if (k.childAttrArr[bodyIndex][index].paramDefault.indexOf("$") != -1) {
                    if (!this.checkList.length) {
                      k.childAttrArr[bodyIndex][index].colorRed = "1";
                      this.infoReviewForm[name + "ChildObj"][bodyIndex][k.childAttrArr[bodyIndex][index].paramName] = objData.paramDataList[Number(k.childAttrArr[bodyIndex][index].paramDefault.split("$")[1])]?.value;
                    }
                    for (
                      var i22 = index + 1;
                      i22 < k.childAttrArr[bodyIndex].length;
                      i22++
                    ) {
                      if (k.childAttrArr[bodyIndex][i22].paramStyle == "radio" && k.childAttrArr[bodyIndex][i22].paramreference && k.childAttrArr[bodyIndex][i22].paramreference != "TempFormTemplate" && k.childAttrArr[bodyIndex][i22].paramreferenceQuery.indexOf("$" + objData.paramName) != -1) {
                        this.focusChildSelect({ row: k }, i22, k.paramName, this.infoReviewForm, k.childAttrArr[bodyIndex][i22], bodyIndex, k.childAttrArr[bodyIndex][i22].paramName);
                      }
                    }
                  }
                }
              }
            }
          });
          this.infoReviewForm = JSON.parse(JSON.stringify(this.infoReviewForm));
          this.infoReviewAttr = JSON.parse(JSON.stringify(this.infoReviewAttr));
        }
      })
    },
    // 子表单下拉确认
    changeChildSelect(value, obj, index, name, formData, formAttr, bodyIndex) {
      var ii = obj.$index;
      if (value == "") {
        if (this.newTaskName != "资源模板新建" && this.newTaskName != "资源模板编辑" && this.newTaskName != "弹性资源模板新建" && this.newTaskName != "弹性资源模板编辑" && this.newTaskName != "弹性资源模板改配") {
          this.infoReviewAttr[ii].childAttrArr[bodyIndex][index].colorRed = "5";
          this.infoReviewAttr[ii].childAttrArr[bodyIndex][index].customDescription = "不可为空";
        }
      } else {
        this.infoReviewAttr[ii].childAttrArr[bodyIndex][index].colorRed = "1";
      }
      formAttr[ii].childAttrArr[bodyIndex].map((k, i) => {
        if (i > index) {
          if (k.paramreferenceQuery) {
            if (
              k.paramreferenceQuery.indexOf(`$${formAttr[ii].childAttrArr[bodyIndex][index].paramName}`) != -1
            ) {
              if (k.paramList) {
                formData[name + "ChildObj"][bodyIndex][k.paramName] = [];
              } else {
                formData[name + "ChildObj"][bodyIndex][k.paramName] = "";
              }
              k.paramDataList = [];
              this.focusChildSelect(obj, i, name, formData, k, bodyIndex, k.paramName);
            }
          }
        }
      });
    },
    // 获取子表单
    showChildAttr(obj, name, val, index, pindex, moGroupstrWhereARR, firstType) {
      this.secondName = name;
      let nameORId = "";
      if (!index && index != 0) {
        this.batchIndex = obj.$index;
        this.infoReviewAttr[this.batchIndex].colorRed = "1";
        this.infoReviewAttr[obj.$index]["childAttrArr"] = this.infoReviewAttr[obj.$index]["childAttrArr"] ? this.infoReviewAttr[obj.$index]["childAttrArr"] : [];
        nameORId = "id";
      } else {
        this.batchIndex = index;
        this.infoReviewAttr[index]["childAttrArr"] = this.infoReviewAttr[index]["childAttrArr"] ? this.infoReviewAttr[index]["childAttrArr"] : [];
        //第一次进来如果默认只有一个批次 默认展开一个，由于参数不能同步所以多加一个参数判断
        if (firstType) {
          nameORId = "id";
        } else {
          nameORId = "name";
        }
      }
      let infoReviewAttrNewArr = [];
      let objArr = {};
      var postData = {"schema":"form_template","where":{[nameORId]:val}}
      Http.getFormTemplate(postData).then((response) => {
        if(response.data.data.data.length){
          var  objAttrData =  JSON.parse(response.data.data.data[0].data)
          if (objAttrData) {
            this.nodeCheckNameThird = objAttrData.name;
            this.historyparamreferenceShow = objAttrData;
            this.infoReviewForm[`${name}ChildObj`] = this.infoReviewForm[`${name}ChildObj`] ? this.infoReviewForm[`${name}ChildObj`] : [];
            let showTitleArr = [];
            objAttrData.parameters.map((k, i) => {
              var tempArr = [];
              if (k.templates) {
                k.templates.map((k, i) => {
                  tempArr.push(k.name);
                });
              }
              var enumQueryArr = [];
              if (k.enum) {
                enumQueryArr = [];
                k.enum.map((k, i) => {
                  enumQueryArr.push({ name: k.zh, value: k.en });
                });
              }
              if ((k.prerequisite || k.type != "string") && k.reference != "FormTemplate") {
                enumQueryArr = [];
              }
              if (k.type == "object" && k.reference && k.list) {
                this.historyparamreferenceShow1 = k.referenceShow ? k.referenceShow : [];
              }
              var bgif = "2";
              if (k.dependentSchema) {
                let objaa = JSON.parse(k.dependentSchema);
                if (moGroupstrWhereARR && Object.keys(moGroupstrWhereARR).length) {
                  var valueDataForm = 0;
                  for (var k222 in objaa) {
                    if (typeof objaa[k222] == "string") {
                      if (moGroupstrWhereARR[k222]) {
                        if (moGroupstrWhereARR[k222].indexOf(objaa[k222]) != -1) {
                          valueDataForm += 1;
                        }
                      }
                    } else {
                      if (moGroupstrWhereARR) {
                        objaa[k222].map((k333, i33) => {
                          if (moGroupstrWhereARR[k222]) {
                            if (JSON.stringify(moGroupstrWhereARR[k222]).indexOf(k333) != -1) {
                              valueDataForm += 1;
                            }
                          }
                        });
                      }
                    }
                  }
                  if (valueDataForm == Object.keys(objaa).length) {
                    bgif = "1";
                  }
                } else {
                  var valueDataForm = 0;
                  for (var k222 in objaa) {
                    if (typeof objaa[k222] == "string") {
                      if (objArr[k222]) {
                        if (objArr[k222].indexOf(objaa[k222]) != -1) {
                          valueDataForm += 1;
                        }
                      }
                    } else {
                      if (objArr[k222] && typeof objArr[k222] != "string") {
                        objArr[k222].map((k333, i333) => {
                          if (objaa[k222].indexOf(k333) != -1) {
                            valueDataForm += 1;
                          }
                        });
                      } else {
                        if (objaa[k222].indexOf(objArr[k222]) != -1) {
                          valueDataForm += 1;
                        }
                      }
                    }
                  }
                  if (valueDataForm == Object.keys(objaa).length) {
                    bgif = "1";
                  }
                }
              } else {
                bgif = "1";
              }
              if (bgif == "1") {
                if (k.availableCondition == "role:admin") {
                  if (this.users.indexOf(this.usernameAll) == -1) {
                    bgif = "2";
                  }
                } else if (
                  k.availableCondition &&
                  k.availableCondition.indexOf("groupRole:") != -1
                ) {
                  var arrGroup = k.availableCondition.split("groupRole:")[1].split(",");
                  var bgG = false;
                  arrGroup.map((gk1, gi1) => {
                    if (this.manageUserGroup.indexOf(gk1) != -1) {
                      bgG = true;
                    }
                  });
                  if (!bgG) {
                    bgif = "2";
                  }
                }
              }
              if (k.joinTitle) {
                showTitleArr.push(k.name);
              }
              infoReviewAttrNewArr.push({
                paramType: k.type,
                paramName: k.name,
                paramTitle: k.title,
                childAttrArrThird: k.reference == "FormTemplate" ? [] : "",
                paramList: k.list,
                paramDescription: k.description,
                paramShemas: k.dependentSchema,
                paramPatterns: k.pattern,
                paramConstraintDescription: k.constraintDescription,
                paramDefault: k.default,
                paramMaxmum: k.maximum,
                paramMinmum: k.minimum,
                paramMaxLength: k.maxLength,
                paramMinLength: k.minLength,
                paramStyle: k.style,
                paramRequired: k.required,
                paramEnum: k.enum ? true : false,
                paramDataList: enumQueryArr,
                paramreference: k.reference,
                paramreferenceZ: tempArr,
                paramreferenceQuery: k.referenceQuery,
                paramreferenceJoinTitle: k.joinTitle,
                paramreferenceLink: k.link,
                paramreferenceShow: k.referenceShow ? k.referenceShow : [],
                paramprerequisite: k.prerequisite ? k.prerequisite : [],
                paramreferenceMutation: k.referenceMutation ? k.referenceMutation : "",
                paramreferenceSelect: k.referenceSelect ? k.referenceSelect : "",
                parmreformat: k.format ? k.format : "",
                thisShowIf: bgif,
                colorRed: "1",
                customDescription: "",
                paramGroup: k.group,
                paramMutipleof: k.mutipleof,
                _id: i + 1,
                ID: k.id,
                paramFromTitle: objAttrData.title,
                paramFromName: objAttrData.name,
                paramTableTitle: [],
                paramreferenceDisplay: k.referenceDisplay ? k.referenceDisplay : "",
                paramavailableCondition: k.availableCondition ? k.availableCondition : "",
                paramremainOperand: objAttrData.mainOperand ? objAttrData.mainOperand.name : "",
                paramannotation: k.annotation,
                paramTemplates: k.templates ? k.templates : "",
              });
              if (bgif == 1) {
                if (k.type == "object") {
                  if (k.list) {
                    objArr[k.name] = [];
                  } else {
                    if (k.default) {
                      if (k.default.indexOf("$") == -1) {
                        objArr[k.name] = k.default ? k.default : "";
                      } else {
                        objArr[k.name] = "";
                      }
                    } else {
                      objArr[k.name] = "";
                    }
                  }
                } else {
                  if (k.list) {
                    objArr[k.name] = k.default ? [k.default] : [];
                  } else {
                    if (k.type == "number") {
                      objArr[k.name] = k.default ? Number(k.default) : "";
                    } else {
                      objArr[k.name] = k.default ? k.default : "";
                    }
                  }
                }
              } else {
                if (k.list) {
                  objArr[k.name] = [];
                } else {
                  objArr[k.name] = "";
                }
              }
              if (k.reference == "FormTemplate") {
                objArr[k.name + "ChildObj3"] = [];
              }
            }); //判断如果入口是预览 则正常
            //展示批次折叠title
            infoReviewAttrNewArr[0]["showTitleArr"] = showTitleArr;
            if (this.lookType && !index) {
              objArr["_isShowStatus"] = true;
            }
            infoReviewAttrNewArr.map((k, i) => {
              if (k.paramType == "object") {
                if (k.paramreferenceZ.length && (!k.paramreference || k.paramreference != "FormTemplate")) {
                  k.paramreferenceZ.map((k1, i1) => {
                    k.paramDataList.push({
                      name: k1,
                      value: k1,
                    });
                  });
                }
                if (!k.paramreferenceQuery) {
                  if (k.paramreference == "FormTemplate") {
                    var arr1 = [];
                    k.paramDataList = [];
                    k.paramTemplates.map((kss, i) => {
                      if (k.paramreferenceJoinTitle) {
                        this.idChangeNameObj[kss.id] = kss[k.paramreferenceJoinTitle];
                      }
                      k.paramDataList.push({
                        name: kss.title,
                        value: kss.id,
                      });
                      this.gaipeiArr3.push({
                        name: kss.name,
                        value: kss.id,
                      });
                      arr1.push({
                        name: kss.title,
                        value: kss.id,
                      });
                    });
                    k.paramDataList = arr1;
                    if (k.paramStyle == "table" && !this.lookType) {
                      if (!index && index != 0) {
                        let pcindex = this.infoReviewAttr[obj.$index]["childAttrArr"].length - 1;
                        this.showChildAttrThird(k, k.paramName, k.paramreferenceZ[0], null, i, {}, obj.$index, pcindex);
                      } else {
                        this.showChildAttrThird(k, k.paramName, k.paramreferenceZ[0], null, i, {}, index, pindex);
                      }
                    }
                  }
                }
              }
            });
            if (!index && index !== 0) {
              this.theadArrShow.push([]);
              this.typeCheckArr.push(val);
              this.infoReviewForm[name] = val;
              this.createChildARRName = name;
              this.infoReviewForm[`${name}ChildObj`].push(objArr);
              this.infoReviewAttr[obj.$index]["childAttrArr"].push(infoReviewAttrNewArr);
              var a = 0;
              infoReviewAttrNewArr.map((k, i) => {
                if (k.paramDefault) {
                  if (k.paramDefault.indexOf("$") != -1) {
                    if (a == 0) {
                      this.focusChildSelect(obj, i, name, this.infoReviewForm, k, this.infoReviewAttr[obj.$index]["childAttrArr"].length - 1);
                    }
                    a += 1;
                  } else {
                    if (this.infoReviewForm[k.paramName] === "") {
                      this.infoReviewForm[k.paramName] = k.paramDefault;
                    }
                  }
                } else {
                  if (k.paramStyle == "radio") {
                    if (a == 0) {
                      this.focusChildSelect(obj, i, name, this.infoReviewForm, k, this.infoReviewAttr[obj.$index]["childAttrArr"].length - 1);
                    }
                    a += 1;
                  }
                }
              });
            } else {
              console.log(this.checkList,333,moGroupstrWhereARR)
              if (this.checkList && this.checkList.length == 0) {
                objArr["formName"] = moGroupstrWhereARR["formName"];
                objArr["formName1"] = moGroupstrWhereARR["formName1"];
              }
              this.infoReviewForm[`${name}ChildObj`][pindex] = objArr;
              this.infoReviewAttr[index]["childAttrArr"][pindex] =
                infoReviewAttrNewArr;
              for (let moItem in moGroupstrWhereARR) {
                //判断如果是对象，则直接赋值name  如果是数组，则赋值到子表单的列表上也就是表格
                if (moGroupstrWhereARR[moItem]) {
                  if (Object.getPrototypeOf(moGroupstrWhereARR[moItem]) === Object.prototype) {
                    if (moGroupstrWhereARR[moItem] && moGroupstrWhereARR[moItem].name) {
                      if (this.infoReviewForm[`${name}ChildObj`][pindex]) {
                        this.infoReviewForm[`${name}ChildObj`][pindex][moItem] = moGroupstrWhereARR[moItem].oid;
                      }
                    } else {
                      this.infoReviewForm[`${name}ChildObj`][pindex][moItem] = moGroupstrWhereARR[moItem].oid;
                    }
                    infoReviewAttrNewArr.map((k222, i222) => {
                      if (k222.paramName == moItem) {
                        k222.paramDataList = [];
                        var nameTitle = "";
                        var referenceArr = []; // [name id path xxxxx]
                        var titleObj = {};
                        if (k222.paramreferenceDisplay) {
                          let mooNew = JSON.parse(JSON.stringify(moGroupstrWhereARR[moItem]));
                          let displayCopy = JSON.parse(k222.paramreferenceDisplay).properties;
                          let mooLast = {};
                          for (let ikey in mooNew) {
                            if (Object.prototype.toString.call(mooNew[ikey]).indexOf("String") != -1) {
                              mooLast[ikey] = mooNew[ikey];
                            } else if (Object.prototype.toString.call(mooNew[ikey]).indexOf("Object") != -1) {
                              for (let iikey in mooNew[ikey]) {
                                mooLast[`${ikey}.${iikey}`] = mooNew[ikey][iikey];
                              }
                            }
                          }
                          displayCopy.map((jitem, j1) => {
                            if (jitem.display) {
                              referenceArr.push(jitem.name);
                              titleObj[jitem.name] = jitem.title;
                            }
                          });
                          const copyArr = displayCopy.find(
                            (item) => item.display == true
                          );
                          nameTitle = copyArr.name;
                          referenceArr.splice(nameTitle, 1);
                        }
                        var str = "";
                        referenceArr.map((k3333, i333) => {
                          if (k3333.indexOf(".") != -1) {
                            let valTemp = JSON.parse(JSON.stringify(moGroupstrWhereARR[moItem]));
                            let arrTemp = k3333.split(".");
                            let valTemp1 = "";
                            arrTemp.map((kk, ii) => {
                              if (Object.prototype.toString.call(valTemp[kk]).indexOf("Object") != -1) {
                                valTemp = valTemp[kk];
                              } else {
                                valTemp1 = valTemp[kk];
                              }
                            });
                            if (valTemp1) {
                              str += titleObj[k3333] ? titleObj[k3333] + ":" + valTemp1 + " " : valTemp1 + " ";
                            } else {
                              str += titleObj[k3333] ? titleObj[k3333] + ":" + " " : " ";
                            }
                          } else {
                            if (moGroupstrWhereARR[moItem][k3333]) {
                              str += titleObj[k3333] ? titleObj[k3333] + ":" + moGroupstrWhereARR[moItem][k3333] + " " : moGroupstrWhereARR[moItem][k3333] + " ";
                            } else {
                              str += titleObj[k3333] ? titleObj[k3333] + ":" + " " : " ";
                            }
                          }
                        });
                        if (nameTitle.indexOf(".") != -1) {
                          let valTemp = JSON.parse(JSON.stringify(moGroupstrWhereARR[moItem]));
                          let arrTemp = nameTitle.split(".");
                          let valTemp1 = "";
                          arrTemp.map((kk, ii) => {
                            if (Object.prototype.toString.call(valTemp[kk]).indexOf("Object") != -1) {
                              valTemp = valTemp[kk];
                            } else {
                              valTemp1 = valTemp[kk];
                            }
                          });
                          if (k222.paramreferenceJoinTitle) {
                            this.idChangeNameObj[moGroupstrWhereARR[moItem].oid] = moGroupstrWhereARR[moItem][k222.paramreferenceJoinTitle];
                          }
                          k222.paramDataList.push({
                            name: valTemp1,
                            value: moGroupstrWhereARR[moItem].oid,
                            content: str,
                          });
                        } else {
                          if (k222.paramreferenceJoinTitle) {
                            this.idChangeNameObj[moGroupstrWhereARR[moItem].oid] = moGroupstrWhereARR[moItem][k222.paramreferenceJoinTitle];
                          }
                          let title = moGroupstrWhereARR[moItem][nameTitle];
                          k222.paramDataList.push({
                            name: title,
                            value: moGroupstrWhereARR[moItem].oid,
                            content: str,
                          });
                        }
                      }
                    });
                  } else if (Array.isArray(moGroupstrWhereARR[moItem])) {
                    let oneline = {};
                    this.infoReviewAttr[index]["childAttrArr"][pindex].map((nnitem, nni) => {
                      if (nnitem.paramName == moItem) {
                        oneline = nnitem;
                      }
                    }
                    );
                    if (Object.keys(oneline).length) {
                      if ((oneline.paramreferenceDisplay && oneline.paramreference != "FormTemplate" && oneline.paramList) || (oneline.paramStyle == "tableStyle" && oneline.paramreference != "FormTemplate")) {
                        this.infoReviewForm[`${name}ChildObj`][pindex][moItem] = moGroupstrWhereARR[moItem] ? moGroupstrWhereARR[moItem] : [];
                      } else {
                        let oidArr1 = [];
                        moGroupstrWhereARR[moItem].map((item, index) => {
                          oidArr1.push(item.oid);
                        });
                        this.infoReviewForm[`${name}ChildObj`][pindex][moItem] = oidArr1 ? oidArr1 : [];
                      }
                    }
                  } else {
                    this.infoReviewForm[`${name}ChildObj`][pindex][moItem] = moGroupstrWhereARR[moItem];
                  }
                } else {
                  if (moItem == "operateType") {
                    infoReviewAttrNewArr.map((kk, ii) => {
                      if (kk.paramName == "operateType") {
                        if (kk.paramList) {
                          this.infoReviewForm[`${name}ChildObj`][pindex][moItem] = [];
                        }
                      }
                    });
                  }
                }
              }
              let typeOfObjectStr = "";
              //判断如果入口是详情，则找到对应数据展示 pindex是Form的下标也就是子下标 index 是Attr下标也就是父下标
              infoReviewAttrNewArr.map((k, i) => {
                if (k.paramType == "object" && k.paramreference == "FormTemplate") {
                  let typeOfObjects = moGroupstrWhereARR[k.paramName];
                  if (typeOfObjects) {
                    typeOfObjects.map((tyItem, tyIndex) => {
                      for (let toItem in tyItem) {
                        if (tyItem[toItem] && toItem.indexOf("joint") != -1) {
                          //获取name请求子表单，相当于添加批次按钮
                          typeOfObjectStr = toItem.match(/joint(\S*)/)[1];
                          let objItem = tyItem[toItem];
                          objItem["formName1"] = tyItem["formName"];
                          this.showChildAttrThird(k, k.paramName, typeOfObjectStr, tyIndex, i, objItem, index, pindex);
                        }
                      }
                    });
                  }
                }
                if (
                  k.paramType == "string" &&
                  Object.keys(moGroupstrWhereARR).length
                ) {
                  objArr[k.paramName] = moGroupstrWhereARR[k.paramName];
                }
              });
              // console.log(infoReviewAttrNewArr)
              this.infoReviewAttr[index]["childAttrArr"][pindex] = infoReviewAttrNewArr;
              this.infoReviewAttr[index]["childAttrArr"][pindex].map((k, i) => {
                if (k.paramreference == "FormTemplate") {
                  this.infoReviewForm[`${name}ChildObj`][pindex][k.paramName + "ChildObj3"] = [];
                }
                if (k.paramList) {
                  this.infoReviewForm[`${name}ChildObj`][pindex][k.paramName] = this.infoReviewForm[`${name}ChildObj`][pindex][k.paramName] ? this.infoReviewForm[`${name}ChildObj`][pindex][k.paramName] : [];
                }
                if (k.paramStyle == "radio") {
                  this.focusChildSelect({ row: this.infoReviewAttr[index] }, i, name, this.infoReviewForm, k, pindex);
                }
              });
              this.theadArrShow.push([]);
              this.typeCheckArr[pindex] = val;
              this.infoReviewForm[name] = val;
              this.createChildARRName = name;
              // 资源列表编辑按钮回填
              if (this.checkList.length) {
                this.infoReviewAttr[index]["childAttrArr"][pindex].map(
                  (k22, i22) => {
                    if (this.checkList[0].qt_direction && k22.paramName == "direction") {
                      this.infoReviewForm[`${name}ChildObj`][pindex]["direction"] = this.checkList[0].qt_direction;
                    }
                    if (k22.paramName == "account") {
                      let valTemp = "";
                      if (this.checkList[0].account.provider.name && this.nodeCheckName != "FormDeleteCloudProduct" && this.nodeCheckName != "FormConfigCloudProduct") {
                        valTemp = "  [" + this.checkList[0].account.provider.name + "]";
                      }
                      if (k22.paramreferenceJoinTitle) {
                        this.idChangeNameObj[this.checkList[0].account.id] = this.checkList[0].account.alias;
                      }
                      k22.paramDataList = [
                        {
                          value: this.checkList[0].account.id,
                          name: this.checkList[0].account.alias,
                          newName: this.checkList[0].account.alias + valTemp,
                        },
                      ];
                      this.infoReviewForm[`${name}ChildObj`][pindex]["account"] = this.checkList[0].account.id;
                    }

                    if (k22.paramName == "cloudAccount") {
                      let valTemp = "";
                      if (this.checkList[0].account.provider.name && this.nodeCheckName != "FormDeleteCloudProduct" && this.nodeCheckName != "FormConfigCloudProduct") {
                        valTemp = "  [" + this.checkList[0].account.provider.name + "]";
                      }
                      if (k22.paramreferenceJoinTitle) {
                        this.idChangeNameObj[this.checkList[0].account.id] = this.checkList[0].account.name;
                      }
                        k22.paramDataList = [
                          {
                            value: this.checkList[0].account.id,
                            name: this.checkList[0].account.alias,
                            newName: this.checkList[0].account.alias + valTemp,
                          },
                        ];
                        this.infoReviewForm[`${name}ChildObj`][pindex]["cloudAccount"] = this.checkList[0].account.id;
                    } else {
                      if (k22.paramName == "operateType") {
                        this.infoReviewForm[`${name}ChildObj`][pindex]["operateType"] = this.checkList[0].operateType ? this.checkList[0].operateType : [];
                      } else {
                        if (k22.paramName == "formObjects" && (k22.paramreferenceDisplay) && k22.paramList) {
                          this.infoReviewForm[`${name}ChildObj`][pindex]["formObjects"] = this.checkList;
                        } else if ((k22.paramName == "formObjects"||k22.paramName == "formObject") && !k22.paramList) {
                          let attrA = [];
                          let attrI = 0;
                          let attrObj = {};
                          this.infoReviewAttr.map((k33, i33) => {
                            if (k33.childAttrArr) {
                              attrA = k33.childAttrArr[0];
                              attrI = i33;
                            }
                          });
                          attrA.map((k33, i33) => {
                            if (k33.paramList) {
                              attrObj[k33.paramName] = [];
                            } else {
                              attrObj[k33.paramName] = "";
                            }
                          });
                          this.checkList.map((k33, i33) => {
                            if (i33 != 0) {
                              this.infoReviewForm[`${name}ChildObj`][i33] = attrObj;
                              this.infoReviewAttr[attrI].childAttrArr.push(attrA);
                            }
                          });
                          this.infoReviewAttr = JSON.parse(JSON.stringify(this.infoReviewAttr));
                          this.infoReviewForm = JSON.parse(JSON.stringify(this.infoReviewForm));
                          this.infoReviewAttr[attrI].childAttrArr.map(
                            (k33, i33) => {
                              let inxx = i33;
                              let a1 = this.checkList[inxx].account.id;
                              let b1 = this.checkList[inxx].account.alias;
                              let a2 = this.checkList[inxx].id;
                              let b2 = this.checkList[inxx].name;
                              if (this.checkList[inxx].account.provider.name) {
                                b1 = this.checkList[inxx].account.alias;
                              }
                              this.infoReviewForm[`${name}ChildObj`][i33]["cloudAccount"] = a1;
                              this.infoReviewForm[`${name}ChildObj`][i33][k22.paramName] = a2;
                              k33.map((k44, i44) => {
                                if (k44.paramName == "cloudAccount") {
                                  if (k44.paramreferenceJoinTitle) {
                                    this.idChangeNameObj[a1] = b1;
                                  }
                                  k44.paramDataList = [
                                    { value: a1, name: b1, newName: b1 },
                                  ];
                                  this.infoReviewForm[`${name}ChildObj`][inxx]["cloudAccount"] = a1;
                                }
                                if (k44.paramName ==k22.paramName) {
                                  if (k44.paramreferenceJoinTitle) {
                                    this.idChangeNameObj[a2] = b2;
                                  }
                                  k44.paramDataList = [
                                    { value: a2, name: b2, newName: b2 },
                                  ];
                                  this.infoReviewForm[`${name}ChildObj`][inxx][k22.paramName] = a2;
                                }
                              });
                            }
                          );
                        }
                      }
                    }
                  }
                );
              }
              //得到atrr所属的那一行的paramreference，用来得到动态表格表头以及请求表头中文数据
              if (this.infoReviewAttr[index]["childAttrArr"]) {
                this.infoReviewAttr[index]["childAttrArr"].map(
                  (arrItem, arrIndex) => {
                    if (arrItem) {
                      arrItem.map((arritem2, arrii) => {
                        if (arritem2.paramList && arritem2.paramType == "object") {
                          if (arritem2.paramreferenceDisplay) {
                            let displayShow = JSON.parse(
                              arritem2.paramreferenceDisplay
                            );
                            let referenceShowArrHeader11 = [];
                            let paramreferenceShowNew = [];
                            let paramreferenceShowNewlast = [];
                            displayShow.properties.map((ditem, dii) => {
                              if (ditem.name != "oid") {
                                paramreferenceShowNew.push(ditem.name);
                              } else {
                                paramreferenceShowNew.push("id");
                              }
                            });
                            paramreferenceShowNew.map((k, i) => {
                              referenceShowArrHeader11.push({
                                title: k,
                                value: k,
                              });
                            });
                            for (let rekey in displayShow.properties) {
                              referenceShowArrHeader11.map((reItem) => {
                                if (displayShow.properties[rekey].name == reItem.value && displayShow.properties[rekey].display) {
                                  paramreferenceShowNewlast.push({
                                    title: displayShow.properties[rekey].title,
                                    value: displayShow.properties[rekey].name,
                                  });
                                }
                              });
                            }
                            arritem2.paramTableTitle = paramreferenceShowNewlast;
                            this.theadArrShow[arrIndex] = paramreferenceShowNewlast;
                          }
                        }
                      });
                    }
                  }
                );
              }
            }
          }
          this.infoReviewForm = JSON.parse(JSON.stringify(this.infoReviewForm));
          this.infoReviewAttr = JSON.parse(JSON.stringify(this.infoReviewAttr));
        }
      })
    },
    // 获取三层子表单
    showChildAttrThird(obj, name, vals, oneindex, twoIndex, moGroupstrWhereARR, yiindex, erindex) {
      this.childData = [];
      this.childObject = [];
      let infoReviewAttrNewArr = [];
      let nameORId = ''
      if(!moGroupstrWhereARR) {
          nameORId = 'id'
      }else{
          nameORId = 'name'
      }
      this.infoReviewAttr[yiindex].childAttrArr[erindex][twoIndex].colorRed = "1";
      let objArr = {};
      var postData = {"schema":"form_template","where":{[nameORId]:vals}}
      Http.getFormTemplate(postData).then((response) => {
        if(response.data.data.data.length){
          var  objAttrData =  JSON.parse(response.data.data.data[0].data)
          this.infoReviewForm[`${this.secondName}ChildObj`][erindex][name + "ChildObj3"] = this.infoReviewForm[`${this.secondName}ChildObj`][erindex][name + "ChildObj3"] ? this.infoReviewForm[`${this.secondName}ChildObj`][erindex][name + "ChildObj3"] : [];
          let showTitleArr = [];
          objAttrData.parameters.map((k, i) => {
            var tempArr = [];
            if (k.templates) {
              k.templates.map((k, i) => {
                tempArr.push(k.name);
              });
            }
            var enumQueryArr = [];
            if (k.enum) {
              enumQueryArr = [];
              k.enum.map((k, i) => {
                enumQueryArr.push({ name: k.zh, value: k.en });
              });
            }
            if (
              (k.prerequisite || k.type != "string") &&
              k.reference != "FormTemplate"
            ) {
              enumQueryArr = [];
            }
            if (k.type == "object" && k.reference && k.list) {
              this.historyparamreferenceShow1 = k.referenceShow
                ? k.referenceShow
                : [];
            }
            var bgif = "2";
            if (k.dependentSchema) {
              let objaa = JSON.parse(k.dependentSchema);
              if (this.lookType) {
                var valueDataForm = 0;
                for (var k222 in objaa) {
                  if (typeof objaa[k222] == "string") {
                    if (moGroupstrWhereARR[k222]) {
                      if (moGroupstrWhereARR[k222].indexOf(objaa[k222]) != -1) {
                        valueDataForm += 1;
                      }
                    }
                  } else {
                    if (moGroupstrWhereARR) {
                      objaa[k222].map((k333, i33) => {
                        if (moGroupstrWhereARR[k222]) {
                          if (
                            JSON.stringify(moGroupstrWhereARR[k222]).indexOf(
                              k333
                            ) != -1
                          ) {
                            valueDataForm += 1;
                          }
                        }
                      });
                    }
                  }
                }
                if (valueDataForm == Object.keys(objaa).length) {
                  bgif = "1";
                }
              } else {
                var valueDataForm = 0;
                for (var k222 in objaa) {
                  if (typeof objaa[k222] == "string") {
                    if (objArr[k222]) {
                      if (objArr[k222].indexOf(objaa[k222]) != -1) {
                        valueDataForm += 1;
                      }
                    }
                  } else {
                    if (objArr[k222] && typeof objArr[k222] != "string") {
                      objArr[k222].map((k333, i333) => {
                        if (objaa[k222].indexOf(k333) != -1) {
                          valueDataForm += 1;
                        }
                      });
                    } else {
                      if (objaa[k222].indexOf(objArr[k222]) != -1) {
                        valueDataForm += 1;
                      }
                    }
                  }
                }
                if (valueDataForm == Object.keys(objaa).length) {
                  bgif = "1";
                }
              }
            } else {
              bgif = "1";
            }
            if (bgif == "1") {
              if (k.availableCondition == "role:admin") {
                if (this.users.indexOf(this.usernameAll) == -1) {
                  bgif = "2";
                }
              } else if (
                k.availableCondition &&
                k.availableCondition.indexOf("groupRole:") != -1
              ) {
                var arrGroup = k.availableCondition.split("groupRole:")[1].split(",");
                var bgG = false;
                arrGroup.map((gk1, gi1) => {
                  if (this.manageUserGroup.indexOf(gk1) != -1) {
                    bgG = true;
                  }
                });
                if (!bgG) {
                  bgif = "2";
                }
              }
            }
            if (k.joinTitle) {
              showTitleArr.push(k.name);
            }
            infoReviewAttrNewArr.push({
              paramType: k.type,
              paramName: k.name,
              paramTitle: k.title,
              childAttrArrThird: k.reference == "FormTemplate" ? [] : "",
              paramList: k.list,
              paramDescription: k.description,
              paramShemas: k.dependentSchema,
              paramPatterns: k.pattern,
              paramConstraintDescription: k.constraintDescription,
              paramDefault: k.default,
              paramMaxmum: k.maximum,
              paramMinmum: k.minimum,
              paramMaxLength: k.maxLength,
              paramMinLength: k.minLength,
              paramStyle: k.style,
              paramRequired: k.required,
              paramEnum: k.enum ? true : false,
              paramDataList: enumQueryArr,
              paramreference: k.reference,
              paramreferenceZ: tempArr,
              paramreferenceQuery: k.referenceQuery,
              paramreferenceJoinTitle: k.joinTitle,
              paramreferenceLink: k.link,
              paramreferenceShow: k.referenceShow ? k.referenceShow : [],
              paramprerequisite: k.prerequisite ? k.prerequisite : [],
              paramreferenceMutation: k.referenceMutation ? k.referenceMutation : "",
              paramreferenceSelect: k.referenceSelect ? k.referenceSelect : "",
              parmreformat: k.format ? k.format : "",
              thisShowIf: bgif,
              colorRed: "1",
              customDescription: "",
              paramGroup: k.group,
              paramMutipleof: k.mutipleof,
              _id: i + 1,
              ID: k.id,
              paramFromTitle: objAttrData.title,
              paramFromName: objAttrData.name,
              paramTableTitle: [],
              paramreferenceDisplay: k.referenceDisplay ? k.referenceDisplay : "",
              paramavailableCondition: k.availableCondition ? k.availableCondition : "",
              paramremainOperand: objAttrData.mainOperand ? objAttrData.mainOperand.name : "",
              paramannotation: k.annotation,
            });
            if (bgif == 1) {
              if (k.type == "object") {
                if (k.list) {
                  objArr[k.name] = [];
                } else {
                  if (k.default) {
                    if (k.default.indexOf("$") == -1) {
                      objArr[k.name] = k.default ? k.default : "";
                    } else {
                      objArr[k.name] = "";
                    }
                  } else {
                    objArr[k.name] = "";
                  }
                }
              } else {
                if (k.list) {
                  objArr[k.name] = k.default ? [k.default] : [];
                } else {
                  if (k.type == "number") {
                    objArr[k.name] = k.default ? Number(k.default) : "";
                  } else {
                    objArr[k.name] = k.default ? k.default : "";
                  }
                }
              }
            } else {
              if (k.list) {
                objArr[k.name] = [];
              } else {
                objArr[k.name] = "";
              }
            }
          });
          infoReviewAttrNewArr[0]["showTitleArr"] = showTitleArr;
          if (!oneindex) {
            objArr["_isShowStatus"] = true;
          }
          infoReviewAttrNewArr.map((k, i) => {
            if (k.paramType == "object") {
              if (k.paramreferenceZ.length && (!k.paramreference || k.paramreference != "FormTemplate")) {
                k.paramreferenceZ.map((k1, i1) => {
                  k.paramDataList.push({
                    name: k1,
                    value: k1,
                  });
                });
              }
              if(k.paramreference){
                let pcindex = this.infoReviewAttr[yiindex]["childAttrArr"][erindex][twoIndex]["childAttrArrThird"].length;
                this.getMetaSchemasList(k.paramreference, k, null, null, yiindex, erindex, twoIndex, pcindex, i);
              }
            }
          });
          if (oneindex === null || oneindex === undefined) {
            this.infoReviewForm[`${this.secondName}ChildObj`][erindex][name] = vals;
            this.infoReviewForm[`${this.secondName}ChildObj`][erindex][name + "ChildObj3"].push(objArr);
            this.infoReviewAttr[yiindex]["childAttrArr"][erindex][twoIndex]["childAttrArrThird"].push(infoReviewAttrNewArr);
          } else {
            this.infoReviewForm[`${this.secondName}ChildObj`][erindex][name + "ChildObj3"][oneindex] = objArr;
            this.infoReviewAttr[yiindex]["childAttrArr"][erindex][twoIndex]["childAttrArrThird"][oneindex] = infoReviewAttrNewArr;
            if (moGroupstrWhereARR) {
              for (let moItem in moGroupstrWhereARR) {
                //判断如果是对象，则直接赋值name  如果是数组，则赋值到子表单的列表上也就是表格
                if (moGroupstrWhereARR[moItem] !== null && moGroupstrWhereARR[moItem] !== undefined) {
                  if (Object.getPrototypeOf(moGroupstrWhereARR[moItem]) === Object.prototype) {
                    infoReviewAttrNewArr.map((k222, i222) => {
                      if (k222.paramName == moItem) {
                        k222.paramDataList = [];
                        var nameTitle = "";
                        var referenceArr = []; // [name id path xxxxx]
                        var titleObj = {};
                        if (k222.paramreferenceDisplay) {
                          let mooNew = JSON.parse(JSON.stringify(moGroupstrWhereARR[moItem]));
                          let displayCopy = JSON.parse(k222.paramreferenceDisplay).properties;
                          let mooLast = {};
                          for (let ikey in mooNew) {
                            if (Object.prototype.toString.call(mooNew[ikey]).indexOf("String") != -1) {
                              mooLast[ikey] = mooNew[ikey];
                            } else if (Object.prototype.toString.call(mooNew[ikey]).indexOf("Object") != -1) {
                              for (let iikey in mooNew[ikey]) {
                                mooLast[`${ikey}.${iikey}`] = mooNew[ikey][iikey];
                              }
                            }
                          }
                          displayCopy.map((jitem, j1) => {
                            if (jitem.display) {
                              referenceArr.push(jitem.name);
                              titleObj[jitem.name] = jitem.title;
                            }
                          });
                          const copyArr = displayCopy.find(
                            (item) => item.display == true
                          );
                          nameTitle = copyArr.name;
                          referenceArr.splice(nameTitle, 1);
                        }
                        var str = "";
                        referenceArr.map((k3333, i333) => {
                          if (k3333.indexOf(".") != -1) {
                            let valTemp = JSON.parse(
                              JSON.stringify(moGroupstrWhereARR[moItem])
                            );
                            let arrTemp = k3333.split(".");
                            let valTemp1 = "";
                            arrTemp.map((kk, ii) => {
                              if (Object.prototype.toString.call(valTemp[kk]).indexOf("Object") != -1) {
                                valTemp = valTemp[kk];
                              } else {
                                valTemp1 = valTemp[kk];
                              }
                            });
                            if (valTemp1) {
                              str += titleObj[k3333] ? titleObj[k3333] + ":" + valTemp1 + " " : valTemp1 + " ";
                            } else {
                              str += titleObj[k3333] ? titleObj[k3333] + ":" + " " : " ";
                            }
                          } else {
                            if (moGroupstrWhereARR[moItem][k3333]) {
                              str += titleObj[k3333] ? titleObj[k3333] + ":" + moGroupstrWhereARR[moItem][k3333] + " " : moGroupstrWhereARR[moItem][k3333] + " ";
                            } else {
                              str += titleObj[k3333] ? titleObj[k3333] + ":" + " " : " ";
                            }
                          }
                        });

                        if (nameTitle.indexOf(".") != -1) {
                          let valTemp = JSON.parse(
                            JSON.stringify(moGroupstrWhereARR[moItem])
                          );
                          let arrTemp = nameTitle.split(".");
                          let valTemp1 = "";
                          arrTemp.map((kk, ii) => {
                            if (Object.prototype.toString.call(valTemp[kk]).indexOf("Object") != -1) {
                              valTemp = valTemp[kk];
                            } else {
                              valTemp1 = valTemp[kk];
                            }
                          });
                          if (k222.paramreferenceJoinTitle) {
                            this.idChangeNameObj[moGroupstrWhereARR[moItem].id] = moGroupstrWhereARR[moItem][k222.paramreferenceJoinTitle];
                          }
                          k222.paramDataList.push({
                            name: valTemp1,
                            value: moGroupstrWhereARR[moItem].oid,
                            content: str,
                          });
                        } else {
                          if (k222.paramreferenceJoinTitle) {
                            this.idChangeNameObj[moGroupstrWhereARR[moItem].id] = moGroupstrWhereARR[moItem][k222.paramreferenceJoinTitle];
                          }
                          let title = moGroupstrWhereARR[moItem][nameTitle];
                          k222.paramDataList.push({
                            name: title,
                            value: moGroupstrWhereARR[moItem].oid,
                            content: str,
                          });
                        }
                      }
                    });
                    if (moGroupstrWhereARR[moItem] && moGroupstrWhereARR[moItem].name) {
                      if (this.infoReviewForm[`${this.secondName}ChildObj`][erindex][`${name}ChildObj3`][oneindex]) {
                        this.infoReviewForm[`${this.secondName}ChildObj`][erindex][`${name}ChildObj3`][oneindex][moItem] = moGroupstrWhereARR[moItem].oid;
                      }
                    } else {
                      this.infoReviewForm[`${this.secondName}ChildObj`][erindex][`${name}ChildObj3`][oneindex][moItem] = moGroupstrWhereARR[moItem].oid;
                    }
                  } else if (Array.isArray(moGroupstrWhereARR[moItem])) {
                    if (moGroupstrWhereARR[moItem]) {
                      let arrorobj = false;
                      infoReviewAttrNewArr.map((k222, i222) => {
                        if (k222.paramList && (k222.paramreferenceDisplay) && k222.paramreference != "FormTemplate" && k222.paramStyle != "search") {
                          this.infoReviewForm[`${this.secondName}ChildObj`][erindex][`${name}ChildObj3`][oneindex][moItem] = moGroupstrWhereARR[moItem] ? moGroupstrWhereARR[moItem] : [];
                          arrorobj = true;
                        }
                      });
                      if (!arrorobj) {
                        moGroupstrWhereARR[moItem].map((sitem, sindex) => {
                          if (!this.infoReviewForm[`${this.secondName}ChildObj`][erindex][`${name}ChildObj3`][oneindex][moItem]) {
                            this.infoReviewForm[`${this.secondName}ChildObj`][erindex][`${name}ChildObj3`][oneindex][moItem] = [];
                          }
                          if (Object.prototype.toString.call(sitem).indexOf("Object") != -1) {
                            this.infoReviewForm[`${this.secondName}ChildObj`][erindex][`${name}ChildObj3`][oneindex][moItem].push(sitem.oid ? sitem.oid : sitem.id);
                            infoReviewAttrNewArr.map((k222, i222) => {
                              if (k222.paramName == moItem) {
                                k222.paramDataList.push({
                                  name: sitem.name,
                                  value: sitem.oid ? sitem.oid : sitem.id,
                                });
                              }
                            });
                          } else {
                            this.infoReviewForm[`${this.secondName}ChildObj`][erindex][`${name}ChildObj3`][oneindex][moItem].push(sitem);
                            infoReviewAttrNewArr.map((k222, i222) => {
                              if (k222.paramName == moItem) {
                                k222.paramDataList.push({
                                  name: sitem,
                                  value: sitem,
                                });
                              }
                            });
                          }
                        });
                      }
                    } else {
                      this.infoReviewForm[`${this.secondName}ChildObj`][erindex][`${name}ChildObj3`][oneindex][moItem] = [];
                    }
                  } else {
                    this.infoReviewForm[`${this.secondName}ChildObj`][erindex][`${name}ChildObj3`][oneindex][moItem] = moGroupstrWhereARR[moItem];
                  }
                } else {
                  if (moItem == "operateType") {
                    infoReviewAttrNewArr.map((kk, ii) => {
                      if (kk.paramName == "operateType") {
                        if (kk.paramList) {
                          this.infoReviewForm[`${this.secondName}ChildObj`][erindex][`${name}ChildObj3`][oneindex][moItem] = [];
                        }
                      }
                    });
                  }
                }
              }
            }
          }
          if (this.infoReviewAttr[yiindex]["childAttrArr"][erindex][twoIndex]["childAttrArrThird"] && Object.prototype.toString.call(this.infoReviewAttr[yiindex]["childAttrArr"][erindex][twoIndex]["childAttrArrThird"]).indexOf("Array") != -1) {
            this.infoReviewAttr[yiindex]["childAttrArr"][erindex][twoIndex]["childAttrArrThird"].map((arrItem, arrIndex) => {
              if (arrItem && Object.prototype.toString.call(arrItem).indexOf("Array") != -1)
                arrItem.map((arritem2, arrii) => {
                  if (arritem2.paramList && arritem2.paramType == "object") {
                    let displayShow = arritem2.paramreferenceDisplay ? JSON.parse(arritem2.paramreferenceDisplay) : { properties: [] };
                    let referenceShowArrHeader11 = [];
                    let paramreferenceShowNew = [];
                    let paramreferenceShowNewlast = [];
                    displayShow.properties.map((ditem, dii) => {
                      if (ditem.name != "oid") {
                        paramreferenceShowNew.push(ditem.name);
                      } else {
                        paramreferenceShowNew.push("id");
                      }
                    });
                    paramreferenceShowNew.map((k, i) => {
                      referenceShowArrHeader11.push({
                        title: k,
                        value: k,
                      });
                    });
                    for (let rekey in displayShow.properties) {
                      referenceShowArrHeader11.map((reItem) => {
                        if (displayShow.properties[rekey].name == reItem.value && displayShow.properties[rekey].display) {
                          paramreferenceShowNewlast.push({
                            title: displayShow.properties[rekey].title,
                            value: displayShow.properties[rekey].name,
                          });
                        }
                      });
                    }
                    arritem2.paramTableTitle = paramreferenceShowNewlast;
                  }
                });
            });
          }
          this.infoReviewForm[`${this.secondName}ChildObj`][erindex][name] = vals;
          this.createChildARRName3 = name;
          this.infoReviewAttr = JSON.parse(JSON.stringify(this.infoReviewAttr));
          this.infoReviewForm = JSON.parse(JSON.stringify(this.infoReviewForm));
          this.$forceUpdate();
          //最外层index
          this.waiIndex = yiindex + "";
        }
      })

    },
    // 复制子表单
    copyCardAttr(obj, name, childIndex) {
      let copyIndex = childIndex;
      let valTemp = this.typeCheckArr[childIndex];
      let bb = this.theadArrShow[childIndex];
      this.typeCheckArr.splice(childIndex + 1, 0, valTemp);
      this.theadArrShow.splice(childIndex + 1, 0, bb);
      let thisInfoForm = JSON.parse(JSON.stringify(this.infoReviewForm));
      let thisInfoAttr = JSON.parse(JSON.stringify(this.infoReviewAttr));
      let cc = thisInfoForm[`${name}ChildObj`][childIndex];
      let dd = thisInfoAttr[obj.$index]["childAttrArr"][childIndex];
      thisInfoForm[`${name}ChildObj`].splice(childIndex, 0, cc);
      thisInfoAttr[obj.$index]["childAttrArr"].splice(childIndex, 0, dd);
      thisInfoAttr[obj.$index]["childAttrArr"].map((k, i) => {
        if (i > copyIndex) {
          k.map((k2, i2) => {
            if (k2.paramprerequisite.length) {
              if (k2.paramList) {
                if (thisInfoForm[`${name}ChildObj`][i][k2.paramName]) {
                  let newArr = [];
                  thisInfoForm[`${name}ChildObj`][i][k2.paramName].map(
                    (k3, i3) => {
                      if (k3) {
                        let valTemp = k3.split(" - ");
                        if (valTemp[0].split("批次")[1] > childIndex + 1) {
                          let bb = `批次${Number(valTemp[0].split("批次")[1]) + 1} - ${valTemp[1]}`;
                          newArr.push(bb);
                        } else {
                          newArr.push(k3);
                        }
                      }
                    }
                  );
                  thisInfoForm[`${name}ChildObj`][i][k2.paramName] = newArr;
                }
              } else {
                if (thisInfoForm[`${name}ChildObj`][i][k2.paramName]) {
                  let valTemp =
                    thisInfoForm[`${name}ChildObj`][i][k2.paramName].split(" - ");
                  if (valTemp[0].split("批次")[1] > childIndex + 1) {
                    let bb = `批次${Number(valTemp[0].split("批次")[1]) + 1} - ${valTemp[1]}`;
                    thisInfoForm[`${name}ChildObj`][i][k2.paramName] = bb;
                  }
                }
              }
            }
            if (k2.childAttrArrThird) {
              k2.childAttrArrThird.map((k3, i3) => {
                k3.map((k4, i4) => {
                  if (k4.paramprerequisite.length) {
                    if (k4.paramList) {
                      if (
                        thisInfoForm[`${name}ChildObj`][i][
                          `${k2.paramName}ChildObj3`
                        ][i3][k4.paramName].length != 0
                      ) {
                        let newArr = [];
                        thisInfoForm[`${name}ChildObj`][i][
                          `${k2.paramName}ChildObj3`
                        ][i3][k4.paramName].map((k5, i5) => {
                          if (k5) {
                            let valTemp = k5.split(" - ");
                            if (valTemp[0].split("批次")[1] > childIndex + 1) {
                              let bb = `批次${Number(valTemp[0].split("批次")[1]) + 1} - ${valTemp[1]}`;
                              newArr.push(bb);
                            } else {
                              newArr.push(k5);
                            }
                          }
                        });
                        thisInfoForm[`${name}ChildObj`][i][`${k2.paramName}ChildObj3`][i3][k4.paramName] = newArr;
                      }
                    } else {
                      if (
                        thisInfoForm[`${name}ChildObj`][i][`${k2.paramName}ChildObj3`][i3][k4.paramName]
                      ) {
                        let valTemp =
                          thisInfoForm[`${name}ChildObj`][i][`${k2.paramName}ChildObj3`][i3][k4.paramName].split(" - ");
                        if (valTemp[0].split("批次")[1] > childIndex + 1) {
                          let bb = `批次${Number(valTemp[0].split("批次")[1]) + 1} - ${valTemp[1]}`;
                          thisInfoForm[`${name}ChildObj`][i][`${k2.paramName}ChildObj3`][i3][k4.paramName] = bb;
                        }
                      }
                    }
                  }
                });
              });
            }
          });
        }
      });
      this.infoReviewForm = JSON.parse(JSON.stringify(thisInfoForm));
      this.infoReviewAttr = JSON.parse(JSON.stringify(thisInfoAttr));
    },
    focusChildGLArr(obj, index, name, pindex) {
      let arrTemp = obj.paramprerequisite;
      let dataList = [];
      this.infoReviewAttr.map((k, i) => {
        if (k.paramName == name) {
          for (var v = 0; v < index; v++) {
            if (arrTemp.indexOf(k.childAttrArr[v][0].paramFromName) != -1) {
              dataList.push({
                name: `批次${v + 1} - ${k.childAttrArr[v][0].paramFromTitle}`,
                value: `批次${v + 1} - ${k.childAttrArr[v][0].paramFromTitle}`,
              });
            }
          }
          k.childAttrArr[index][pindex].paramDataList = dataList;
        }
      });
    },
    deletCardAttr(obj, name, childIndex) {
      let nameForm = obj.row.childAttrArr[childIndex][0].paramFromName;
      let title = `批次${childIndex + 1} - ${obj.row.childAttrArr[childIndex][0].paramFromTitle}`;
      this.infoReviewAttr[obj.$index].childAttrArr.map((k, i) => {
        if (i > childIndex) {
          k.map((k2, i2) => {
            if (k2.paramprerequisite.length && k2.paramprerequisite.indexOf(nameForm) != -1) {
              if (Object.prototype.toString.call(this.infoReviewForm[`${name}ChildObj`][i][k2.paramName]).indexOf("Array") != -1) {
                if (this.infoReviewForm[`${name}ChildObj`][i][k2.paramName].indexOf(title) != -1) {
                  this.infoReviewForm[`${name}ChildObj`][i][k2.paramName].splice(this.infoReviewForm[`${name}ChildObj`][i][k2.paramName].indexOf(title), 1);
                }
              } else {
                if (this.infoReviewForm[`${name}ChildObj`][i][k2.paramName] == title) {
                  this.infoReviewForm[`${name}ChildObj`][i][k2.paramName] = "";
                }
              }
            }
            if (k2.paramprerequisite.length) {
              if (Object.prototype.toString.call(this.infoReviewForm[`${name}ChildObj`][i][k2.paramName]).indexOf("Array") != -1) {
                let newArr = [];
                this.infoReviewForm[`${name}ChildObj`][i][k2.paramName].map(
                  (k3, i3) => {
                    if (k3) {
                      let valTemp = k3.split(" - ");
                      if (valTemp[0].split("批次")[1] > childIndex) {
                        let bb = `批次${Number(valTemp[0].split("批次")[1]) - 1} - ${valTemp[1]}`;
                        newArr.push(bb);
                      } else {
                        newArr.push(k3);
                      }
                    }
                  }
                );
                this.infoReviewForm[`${name}ChildObj`][i][k2.paramName] = newArr;
              } else {
                if (this.infoReviewForm[`${name}ChildObj`][i][k2.paramName] != "") {
                  if (this.infoReviewForm[`${name}ChildObj`][i][k2.paramName]) {
                    let valTemp = this.infoReviewForm[`${name}ChildObj`][i][k2.paramName].split(" - ");
                    if (valTemp[0].split("批次")[1] > childIndex) {
                      let bb = `批次${Number(valTemp[0].split("批次")[1]) - 1} - ${valTemp[1]}`;
                      this.infoReviewForm[`${name}ChildObj`][i][k2.paramName] = bb;
                    }
                  }
                }
              }
            }
            if (k2.childAttrArrThird) {
              k2.childAttrArrThird.map((k3, i3) => {
                k3.map((k4, i4) => {
                  if (k4.paramprerequisite.length && k4.paramprerequisite.indexOf(nameForm) != -1) {
                    if (Object.prototype.toString.call(this.infoReviewForm[`${name}ChildObj`][i][k2.paramName + "ChildObj3"][i3][k4.paramName]).indexOf("Array") != -1) {
                      if (this.infoReviewForm[`${name}ChildObj`][i][k2.paramName + "ChildObj3"][i3][k4.paramName].indexOf(title) != -1) {
                        this.infoReviewForm[`${name}ChildObj`][i][k2.paramName + "ChildObj3"][i3][k4.paramName].splice(this.infoReviewForm[`${name}ChildObj`][i][k2.paramName + "ChildObj3"][i3][k4.paramName].indexOf(title), 1);
                      }
                    } else {
                      if (this.infoReviewForm[`${name}ChildObj`][i][k2.paramName + "ChildObj3"][i3][k4.paramName] == title) {
                        this.infoReviewForm[`${name}ChildObj`][i][k2.paramName + "ChildObj3"][i3][k4.paramName] = "";
                      }
                    }
                  }
                  if (k4.paramprerequisite.length) {
                    if (Object.prototype.toString.call(this.infoReviewForm[`${name}ChildObj`][i][k2.paramName + "ChildObj3"][i3][k4.paramName]).indexOf("Array") != -1) {
                      let newArr = [];
                      this.infoReviewForm[`${name}ChildObj`][i][k2.paramName + "ChildObj3"][i3][k4.paramName].map((k5, i5) => {
                        if (k5) {
                          let valTemp = k5.split(" - ");
                          if (valTemp[0].split("批次")[1] > childIndex) {
                            let bb = `批次${Number(valTemp[0].split("批次")[1]) - 1} - ${valTemp[1]}`;
                            newArr.push(bb);
                          } else {
                            newArr.push(k5);
                          }
                        }
                      });
                      this.infoReviewForm[`${name}ChildObj`][i][k2.paramName + "ChildObj3"][i3][k4.paramName] = newArr;
                    } else {
                      if (this.infoReviewForm[`${name}ChildObj`][i][k2.paramName + "ChildObj3"][i3][k4.paramName] != "") {
                        let valTemp = this.infoReviewForm[`${name}ChildObj`][i][k2.paramName + "ChildObj3"][i3][k4.paramName].split(" - ");
                        if (valTemp[0].split("批次")[1] > childIndex) {
                          let bb = `批次${Number(valTemp[0].split("批次")[1]) - 1} - ${valTemp[1]}`;
                          this.infoReviewForm[`${name}ChildObj`][i][k2.paramName + "ChildObj3"][i3][k4.paramName] = bb;
                        }
                      }
                    }
                  }
                });
              });
            }
          });
        }
      });
      this.typeCheckArr.splice(childIndex, 1);
      this.theadArrShow.splice(childIndex, 1);
      let thisInfoForm = JSON.parse(JSON.stringify(this.infoReviewForm));
      let thisInfoAttr = JSON.parse(JSON.stringify(this.infoReviewAttr));
      thisInfoForm[`${name}ChildObj`].splice(childIndex, 1);
      thisInfoAttr[obj.$index]["childAttrArr"].splice(childIndex, 1);
      this.infoReviewForm = JSON.parse(JSON.stringify(thisInfoForm));
      this.infoReviewAttr = JSON.parse(JSON.stringify(thisInfoAttr));
      if (this.infoReviewAttr[obj.$index].childAttrArr.length == 0) {
        this.infoReviewAttr[obj.$index].colorRed = "5";
      }
    },
    changeSelectValue(value, obj, name, formData, formAttr) {
      let index = obj.$index;
      if (value == "") {
        if (this.newTaskName != "资源模板新建" && this.newTaskName != "资源模板编辑" && this.newTaskName != "弹性资源模板新建" && this.newTaskName != "弹性资源模板编辑" && this.newTaskName != "弹性资源模板改配") {
          if (this.infoReviewAttr[index].paramRequired) {
            this.infoReviewAttr[index].colorRed = "5";
            this.infoReviewAttr[index].customDescription = "不可为空";
          } else {
            this.infoReviewAttr[index].colorRed = "1";
          }
        }
      } else {
        this.infoReviewAttr[index].colorRed = "1";
      }
      for (var k in formAttr) {
        if (formAttr[k].paramreferenceQuery) {
          if (formAttr[k].paramreferenceQuery.indexOf(`$${name}`) != -1) {
            if (formAttr[k].paramList) {
              this.infoReviewForm[formAttr[k].paramName] = [];
            } else {
              this.infoReviewForm[formAttr[k].paramName] = "";
            }
          }
        }

        if (formAttr[k].childAttrArr) {
          if (JSON.stringify(formAttr[k].childAttrArr).indexOf(`$${name}`) != -1) {
            formAttr[k].childAttrArr.map((k11, i11) => {
              k11.map((k22, i22) => {
                if (k22.paramreferenceQuery) {
                  if (k22.paramreferenceQuery.indexOf("$") != -1) {
                    if (Object.prototype.toString.call(this.infoReviewForm[formAttr[k].paramName + "ChildObj"][i11][k22.paramName]).indexOf("Array") != -1) {
                      this.infoReviewForm[formAttr[k].paramName + "ChildObj"][i11][k22.paramName] = [];
                    } else {
                      this.infoReviewForm[formAttr[k].paramName + "ChildObj"][i11][k22.paramName] = "";
                    }
                    if (k22.paramName == "formObjects") {
                      this.infoReviewForm[formAttr[k].paramName + "ChildObj"][i11]["formObjects"] = k22.paramList ? [] : "";
                    }
                    if (k22.paramStyle == "radio") {
                      formAttr[k].childAttrArr[i11][i22].paramDataList = [];
                      this.focusChildSelect({ row: formAttr[k] }, i22, formAttr[k].paramName, this.infoReviewForm, k22, i11, k22.paramName);
                    }
                  }
                }
              });
            });
          }
        }
      }
      for (var i = index + 1; i < this.infoReviewAttr.length; i++) {
        if (this.infoReviewAttr[i].paramreferenceQuery) {
          if (this.infoReviewAttr[i].paramreferenceQuery.indexOf(`$${name}`) != -1) {
            this.focusSelectValue({ row: this.infoReviewAttr[i] }, this.infoReviewAttr[i].paramName, this.infoReviewForm);
          }
        }
      }
    },
    // 表单同意提交
    subAddModelDeny() {
      if (this.jobIdStatus && this.jobIdStatus != "审核中") {
        this.$message({
          showClose: true,
          message: "不能同意非审核中的任务！",
          type: "warning",
        });
        return false;
      }
      this.reviewload = true;
      let postDataG = {
        process_id: this.historyTableList._cometProcess,
        message_type: "disapprove",
        dealer: this.usernameAll,
      };
      
    },
    resetVal() {
      this.batchSearchName = "";
      this.batchSearchName1 = "";
      this.showCidSearch = false;
    },
    // 分页
    changehandleCurrent(val) {
      this.searchpageNum = val;
      this.searchTableServer();
    },
    handleSizeChange(val) {
      this.searchpagesize = val;
      this.searchpageNum = 1;
      this.searchTableServer(true);
    },
    // 上传文件
    successUpload(response, file, fileList, name, index) {
      if (response.code != 0) {
        this.$message({
          showClose: true,
          message: "文件上传失败：" + response.data,
          type: "error",
        });
      }
      var arrTemp = [];
      fileList.map((k, i) => {
        if (k.response.code == 0) {
          arrTemp.push({
            name: k.name,
            url: k.response.data,
            response: k.response,
          });
        }
      });
      this.infoReviewForm[name] = arrTemp;
      this.infoReviewAttr[index].colorRed = "1";
    },
    handleRemove(file, fileList, name) {
      var arrTemp = [];
      fileList.map((k, i) => {
        if (k.response.code == 0) {
          arrTemp.push({
            name: k.name,
            url: k.response.data,
            response: k.response,
          });
        }
      });
      this.infoReviewForm[name] = arrTemp;
    },
    successUpload1(response, file, fileList, name, index, chidlName, iii, iiii) {
      if (response.code != 0) {
        this.$message({
          showClose: true,
          message: "文件上传失败：" + response.data,
          type: "error",
        });
      }
      var arrTemp = [];
      fileList.map((k, i) => {
        if (k.response.code == 0) {
          arrTemp.push({
            name: k.name,
            url: k.response.data,
            response: k.response,
          });
        }
      });
      this.infoReviewForm[name + "ChildObj"][iii][chidlName] = arrTemp;
      this.infoReviewAttr[index].childAttrArr[iii][iiii].colorRed = "1";
    },
    handleRemove1(file, fileList, name, index, chidlName, iii, iiii) {
      var arrTemp = [];
      fileList.map((k, i) => {
        if (k.response.code == 0) {
          arrTemp.push({
            name: k.name,
            url: k.response.data,
            response: k.response,
          });
        }
      });
      this.infoReviewForm[name + "ChildObj"][iii][chidlName] = arrTemp;
      this.infoReviewAttr[index].childAttrArr[iii][iiii].colorRed = "1";
    },
    changebatchSearchInput(obj, name, str, attr, aiii) {
      var batcharr = [];
      var batchNewarr = [];
      batcharr = str.split("\n");
      batcharr.map((k, i) => {
        if (k.replace(/\s/g, "")) {
          batchNewarr.push(k.replace(/\s/g, ""));
        }
      });
      if (!batchNewarr.length) return;
      obj.map((k, i) => {
        k[name] = batchNewarr[i];
        if (attr[i]) {
          if (batchNewarr[i]) {
            attr[i][aiii].colorRed = "1";
          }
        }
      });
    },
  },
  computed: {
    ...mapGetters(["manageUser", "manageUserGroup"]),
  },
  created() {
    this.users = this.manageUser;
  },
};

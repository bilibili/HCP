import Http from "@/components/api/services";
export const productLIst = {
  //方法集合
  methods: {
    addNewsearch(type, level, ownIndex, parentIndex, oneIndex, obj) {
      if (level == 1) {
        if (this.advSearch[0].category == '') {
          this.advSearch[0].category = type
          this.advSearch.push({
            name: '',
            type: '',
            val: '',
            category: type,
            index: [],
            child: [],
            objContent: {},
          })
        } else {
          if (this.advSearch.length == 1) {
            this.advSearch[0].category = type
            this.advSearch.push({
              name: '',
              type: '',
              val: '',
              category: type,
              index: [],
              child: [],
              objContent: {},
            })
          } else {
            if (this.advSearch[ownIndex].category == type) {
              this.advSearch.push({
                name: '',
                type: '',
                val: '',
                category: type,
                index: [],
                child: [],
                objContent: {},
              })
            } else {
              let arr = []
              arr.push({
                name: obj.name,
                type: obj.type,
                val: obj.val,
                category: type,
                index: obj.index,
                child: [],
                objContent: obj.objContent,
              })
              arr.push({
                name: '',
                type: '',
                val: '',
                category: type,
                index: [],
                child: [],
                objContent: {},
              })
              this.advSearch[ownIndex].child = arr
            }
          }
        }
      }
      if (level == 2) {
        if (this.advSearch[parentIndex].child[ownIndex].category == type) {
          this.advSearch[parentIndex].child.push({
            name: '',
            type: '',
            val: '',
            category: type,
            index: [],
            child: [],
            objContent: {},
          })
        } else {
          let arr = []
          arr.push({
            name: obj.name,
            type: obj.type,
            val: obj.val,
            category: type,
            index: obj.index,
            objContent: obj.objContent,
          })
          arr.push({
            name: '',
            type: '',
            val: '',
            category: type,
            index: [],
            objContent: {},
          })
          this.advSearch[parentIndex].child[ownIndex].child = arr
        }
      }
      if (level == 3) {
        this.advSearch[oneIndex].child[parentIndex].child.push({
          name: '',
          type: '',
          val: '',
          category:
            this.advSearch[oneIndex].child[parentIndex].child[0].category,
          index: [],
          objContent: {},
        })
      }
      this.advSearch = JSON.parse(JSON.stringify(this.advSearch))
    },
    delNewsearch(level, ownIndex, parentIndex, oneIndex, obj) {
      let advSearchData = JSON.parse(JSON.stringify(this.advSearch))
      if (level == 1) {
        if (
          advSearchData.length == 1 &&
          ownIndex == 0 &&
          (!obj.child || !obj.child.length)
        ) {
          return
        }
        advSearchData.splice(ownIndex, 1)
        if (advSearchData.length == 1) {
          if (advSearchData[0].child && advSearchData[0].child.length) {
            advSearchData = advSearchData[0].child
          }
        }
      }
      if (level == 2) {
        advSearchData[parentIndex].child.splice(ownIndex, 1)
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
          }
        }
      }
      if (level == 3) {
        advSearchData[oneIndex].child[parentIndex].child.splice(ownIndex, 1)
        if (advSearchData[oneIndex].child[parentIndex].child.length == 1) {
          advSearchData[oneIndex].child[parentIndex] = {
            name: advSearchData[oneIndex].child[parentIndex].child[0].name,
            type: advSearchData[oneIndex].child[parentIndex].child[0].type,
            val: advSearchData[oneIndex].child[parentIndex].child[0].val,
            category: advSearchData[oneIndex].child[parentIndex].category,
            index: advSearchData[oneIndex].child[parentIndex].child[0].index,
            objContent:
              advSearchData[oneIndex].child[parentIndex].child[0].objContent,
            child: advSearchData[oneIndex].child[parentIndex].child[0].child
              ? advSearchData[oneIndex].child[parentIndex].child[0].child
              : [],
          }
        }
      }
      this.advSearch = JSON.parse(JSON.stringify(advSearchData))
      this.$forceUpdate()
    },
    // 导出表格
    openDownload() {
      this.isIndeterminate = false
      this.checkAll = false
      this.dialogDownload = true
      this.checkDownload = []
    },
    // 全选
    handleCheckAllChange(val) {
      this.checkDownload = val ? this.excelTitleDown : []
      this.isIndeterminate = false
    },
    // 取消全选
    handleCheckedCitiesChange(value) {
      let checkedCount = value.length
      this.checkAll = checkedCount === this.excelTitleDown.length
      this.isIndeterminate =
        checkedCount > 0 && checkedCount < this.excelTitleDown.length
    },
    // 导出确认
    subDownload() {
      if (this.checkDownload.length == 0) {
        return false
      }
      this.getDownloadExcel(this.nodeTitle)
    },
    // 获取下载字段
    getDownloadExcel(name) {},
    // 任务类型点击taskname英文名称，请求接口 chinaName中文名称 id项目账号创建的id
    changeAddnew(taskname, chinaName) {
      this.dialogVisible = false
      this.newTaskName = chinaName
      this.nodeCheckName = taskname
      this.checkList = []
      this.reviewStatus222 = true
    },
    //高级搜索,点击菜单发生变化时
    searchForm(item) {
      this.searchFormData = {}
      this.tableHeaderArr.map((key) => {
        if (key.name == this.cascaderMorestr) {
          key.indexSearch = item
          this.searchFormData = JSON.parse(JSON.stringify(key))
        }
      })
    },
    //高级搜索,默认选择第一个子菜单时
    searchFormchange() {
      this.cascaderMorestrSecond = null
      this.searchFormData = {}
      this.searchMoreArr.map((item, index) => {
        if (item.name == this.cascaderMorestr) {
          this.searchMoreArrSecond = item.children
        }
      })
      let twoStr = ''
      let firstStr = this.cascaderMorestr
      this.searchMoreArr.map((key) => {
        if (key.name == firstStr) {
          twoStr = key.children[0].value
        }
      })
    },
    entrySreach(item, obj) {
      let oldname = item.enName.split('+')[0]
      let oldsecond = item.enName.split('+')[1] ? item.enName.split('+')[1] : ''
      this.cascaderMorestr = oldname
      this.searchFormData = {}
      this.searchMoreArr.map((item, index) => {
        if (item.name == this.cascaderMorestr) {
          this.searchMoreArrSecond = item.children
        }
      })
      let twoStr = ''
      let firstStr = this.cascaderMorestr
      this.searchMoreArr.map((key) => {
        if (key.name == firstStr) {
          twoStr = key.children[0].value
        }
      })
      this.cascaderMorestrSecond = oldsecond
      this.tableHeaderArr.map((key) => {
        if (key.name == this.cascaderMorestr) {
          key.indexSearch = this.cascaderMorestrSecond
          this.searchFormData = JSON.parse(JSON.stringify(key))
        }
      })
      this.poHide = true
    },
    searchTagChange(addvitem) {
      let eItemTxt = ''
      let searchObjTagss = []
      if (addvitem.val) {
        this.tableHeaderArr.map((item) => {
          if (item.name == addvitem.name) {
            if (item.enum) {
              eItemTxt += '['
              item.enum.map((eitem) => {
                if (Array.isArray(addvitem.val)) {
                  addvitem.val.map((kkey, kki) => {
                    if (eitem.en == kkey) {
                      eItemTxt += `${eitem.zh},`
                    }
                  })
                } else {
                  if (eitem.en == addvitem.val) {
                    eItemTxt = eitem.zh
                  }
                }
              })
              eItemTxt = eItemTxt.slice(0, eItemTxt.length - 1)
              eItemTxt += ']'
            } else {
              eItemTxt = addvitem.val
            }
            searchObjTagss = {
              name: addvitem.objContent.label,
              value: eItemTxt,
              type: addvitem.type ? addvitem.type : 'select',
              category: addvitem.category,
              child: addvitem.child,
            }
          }
        })
      }
      return searchObjTagss
    },
    refAdvancedSearch(index) {
      this.searchObjTag.splice(index, 1)
      this.lastAdvSearch.splice(index, 1)
      this.getTable(this.nodeTitle)
    },
    refAdvancedSearchLast() {
      this.searchObjTag = []
      this.lastAdvSearch = []
      this.getTable(this.nodeTitle)
    },
    searchShowStatus() {
      this.advSearchIndex = null
      this.poHide = true
      this.advSearch = [
        {
          name: '',
          type: null,
          val: '',
          category: '',
          index: [],
          child: [],
          objContent: {},
        },
      ]
    },
    // 复制名称功能
    _copy(context) {
      // 创建输入框元素
      let oInput = document.createElement('input')
      // 将想要复制的值
      oInput.value = context
      // 页面底部追加输入框
      document.body.appendChild(oInput)
      // 选中输入框
      oInput.select()
      // 执行浏览器复制命令
      document.execCommand('Copy')
      // 弹出复制成功信息
      this.$message.success('复制成功')
      // 复制后移除输入框
      oInput.remove()
    },
    routeChange(obj) {
      if (!this.$route.query.detail_ID) {
        this.drawDetail = false
      }
    },
    // 关闭弹框
    addCloseDraw() {
      this.tableHight = ''
      this.tableDetailStatus = ''
      this.drawDetail = false
      this.paginationSatatus = false
      this.tablePagina.style.textAlign = 'center'
      this.tableRowIndex = null
    },
    // 表格高度
    tableHeight(height, posti, status) {
      this.tableHight = height
      this.tablePagina.style.textAlign = posti
      this.paginationSatatus = status
    },
    resetSearch() {
      this.advSearch = [
        {
          name: '',
          type: null,
          val: '',
          category: '',
          index: [],
          child: [],
          objContent: {},
        },
      ]
    },
    // 展示表格列
    showTH() {
      let checkboxName = `${this.$route.name}checkboxTHBFArr`
      localStorage.setItem(checkboxName, this.checkboxTH)
      if (this.checkboxTH.length == 0) {
        this.checkboxTH = JSON.parse(JSON.stringify(this.checkboxTHBF))
      } else {
        this.checkboxTHBF = JSON.parse(JSON.stringify(this.checkboxTH))
      }
    },
    showTHRefresh() {
      let checkboxName = `${this.$route.name}checkboxTHBFArr`
      localStorage.removeItem(checkboxName)
      this.getTreeData()
    },
    getselectList1(val, obj) {
      if (val) {
        this.getSelList(obj, val)
      }
    },
    getselectList(val) {
      if (val) {
        this.getSelList(this.searchFormData, val)
      }
    },
    getSelList(obj, value) {
      let postData = {}
      let keyV = `${obj.name.split('.')[obj.name.split('.').length - 1]}`
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
      Http.getQueryList(postData).then(
        (response) => {
          if (response.data.data.data) {
            if (obj.DataList) {
              obj.DataList = response.data.data.data
            }
          }
        }
      ).catch((err) => {
        console.log(err);
      });
    },
    // 表格数据
    getTable(name, num) {
      this.tableLoading = true
      if (!name) {
        this.tableLoading = false
        return
      }
      /* if (this.$route.meta.projectNoget) {
        if (!this.checkProjectData.length) {
          this.tableLoading = false
          return
        }
      } */
      if (num) {
        this.pageNum = num
      }
     
      this.tableData = []

      var or = {}

      /* if(this.checkProjectData.length){
        or["project_config"] = {
          name_IN : this.checkProjectData
        }
      } */

      if (this.batchSearchName != '') {
        var showlistObjSearch = []
        if (this.showlistObj.search) {
          showlistObjSearch = JSON.parse(JSON.stringify(this.showlistObj.search))
        } else {
          showlistObjSearch = ['name']
        }
        // var arr = this.batchSearchName.split(" ").filter(item => item)
        showlistObjSearch.map((item) => {
            var searKey1 = item.split('.').reverse()
            var str11 = ''
            searKey1.map((k1, i1) => {
              if (i1 == 0) {
                /* if(arr.length>1){
                  str11 = `{"${k1}_IN":${JSON.stringify(arr)}}`
                }else{
                  } */
               str11 = `{"${k1}_REGEX":"${this.batchSearchName}"}`
              } else {
                str11 = `{"${k1}":${str11}}`
              }
            })
            var getKet = Object.keys(JSON.parse(str11))[0]
            or[getKet] = JSON.parse(str11)[getKet]
        })
      }

      var and = [];
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
                          this.advSearchHandle(advItem3, name)
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
                    arr2[secondType].push(this.advSearchHandle(advItem2, name));
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
                arr1[firstType].push(this.advSearchHandle(advItem, name));
              }
            });
            and.push(arr1);
          } else if (!lastItem[0].category && lastItem[0].val) {
            if (Object.prototype.toString.call(lastItem[0].val).indexOf("Array") != -1) {
              if (lastItem[0].val.length == 0) {
                return;
              }
            }
            let arr1 = this.advSearchHandle(lastItem[0], name);
            and.push(arr1);
          }
        });
      }

      var postData = {
        schema: name,
        page_size:this.pageSize,
        page_num:this.pageNum,
        order:"id DESC",
        where:{}
      }
      if(Object.keys(or).length){
        postData.where.or = or
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
        postData.where.and = arrAnd
      }
      Http.getQueryList(postData).then((response) => {
          this.total = response.data.data.total
          this.tableData = []
          this.tableData = response.data.data.data
          this.tableLoading = false
          this.tableRowIndex = null
          this.tableHeaderArr.map((k1, i1) => {
            if(k1.name!='id'){
              this.tableData.map((k2, i2) => {
                if (this.$route.query.detail_ID) {
                  if (k2.id == this.$route.query.detail_ID) {
                    this.tableRowIndex = i2
                  }
                }
                let arr = k1.name.split('.')
                let tempObj = [JSON.parse(JSON.stringify(k2))]
                for (let i2 of arr) {
                  let tempObjArr = []
                  for (let i3 of tempObj) {
                    let tempValue = ''
                    if (i3[i2] || i3[i2] === 0) {
                      tempValue = i3[i2]
                      if (tempValue instanceof Array) {
                        tempObjArr.push(...tempValue)
                      } else {
                        if (k1.type == 'integer') {
                          tempValue = tempValue == '' ? '0' : tempValue
                        }
                        tempObjArr.push(tempValue)
                      }
                    }
                  }
                  tempObj = tempObjArr
                }
                if (k1.tagsArr) {
                  if (tempObj.length != 0 && tempObj.join() != '') {
                    let tempTagsArr = []
                    tempObj.map((k33, i33) => {
                      let tempTagsArrValue = ''
                      k1.tagsArr.map((k44, i44) => {
                        let tempTagsArrKey = ''
                        if (i44) {
                          tempTagsArrKey = ':'
                        }
                        tempTagsArrValue += `${tempTagsArrKey}${k33[k44]}`
                      })
                      if (tempTagsArrValue.length) {
                        tempTagsArr.push(tempTagsArrValue)
                      }
                    })
                    k2[k1.name.replace(/\./g, '')] = tempTagsArr
                  }
                } else {
                  if ((k1.title.indexOf('时间') != -1 || k1.isTime) && !k1.enum) {
                    if (k1.isTimestamp) {
                      if (tempObj.join(',')) {
                        k2[k1.name.replace(/\./g, '')] = this.$moment(
                          Number(tempObj.join(',')) * 1000
                        ).format('YYYY-MM-DD HH:mm:ss')
                      }
                    } else {
                      if (tempObj.join(',') && tempObj.join(',').length == 10) {
                        k2[k1.name.replace(/\./g, '')] = tempObj.join(',')
                          ? this.$moment(tempObj.join(',')).format('YYYY-MM-DD')
                          : tempObj.join(',')
                      } else {
                        k2[k1.name.replace(/\./g, '')] = tempObj.join(',')
                          ? this.$moment(tempObj.join(',')).format(
                              'YYYY-MM-DD HH:mm:ss'
                            )
                          : tempObj.join(',')
                      }
                    }
                  } else {
                    if(tempObj.length==1){
                      k2[k1.name.replace(/\./g, '')] = tempObj[0]
                    }else{
                      k2[k1.name.replace(/\./g, '')] = tempObj.join(',')
                    }
                  }
                }
              })
            }
          })
          console.log(this.tableData)
          if (this.$route.query.detail_ID) {
            this.tableDetailStatus = sessionStorage.getItem(
              'tableDetailStatusC'
            )
              ? sessionStorage.getItem('tableDetailStatusC')
              : 'halfScreen'
            if (this.tableDetailStatus == 'halfScreen') {
              this.$nextTick(() => {
                this.tableHight = 'up'
                this.paginationSatatus = false
                this.tablePagina.style.textAlign = 'center'
              })
            }
            if (this.tableDetailStatus == 'halfScreenBan') {
              this.tablePagina.style.textAlign = 'left'
              this.paginationSatatus = true
              this.tableHight = ''
            }
            if (this.tableDetailStatus == 'halfAll') {
              this.paginationSatatus = false
              this.tableHight = ''
              this.tablePagina.style.textAlign = 'center'
            }
          }
        }
      ).catch((err) => {
        this.tableLoading = false
      });
      this.poHide = false
      this.advSearch = [
        {
          name: '',
          type: null,
          val: '',
          category: '',
          index: [],
          child: [],
          objContent: {},
        },
      ]
    },
    // 分页
    handleSizeChange(val) {
      this.pageSize = val
      this.pageNum = 1
      this.getTable(this.nodeTitle)
    },
    changehandleCurrent(val) {
      this.pageNum = val
      this.getTable(this.nodeTitle)
    },
    // 状态请求数据
    hintTableSearch(title) {
      this.hintTitle = title
      this.pageNum = 1
      this.getTable(this.nodeTitle)
    },
    // 高级搜索确认
    searchTable(type) {
      this.hintTitle = ''
      if (this.advSearch[0].val || this.advSearch[0].category) {
        if (this.advSearchIndex || this.advSearchIndex === 0) {
          this.lastAdvSearch[this.advSearchIndex] = JSON.parse(
            JSON.stringify(this.advSearch)
          )
        } else {
          this.lastAdvSearch.push(JSON.parse(JSON.stringify(this.advSearch)))
        }
        var searchObjTagLast = JSON.parse(JSON.stringify(this.lastAdvSearch))
        searchObjTagLast.map((lastItem, lastI) => {
          lastItem.map((advItem, advI) => {
            if (advItem.child && advItem.child.length) {
              advItem.child.map((advItem2, adv2) => {
                if (advItem2.child && advItem2.child.length) {
                  advItem2.child.map((advItem3, adv3) => {
                    if (!advItem3.child || !advItem3.child.length) {
                      lastItem[advI].child[adv2].child[adv3] =
                        this.searchTagChange(advItem3)
                    }
                  })
                } else {
                  lastItem[advI].child[adv2] = this.searchTagChange(advItem2)
                }
              })
            } else {
              lastItem[advI] = this.searchTagChange(advItem)
            }
          })
        })
        this.searchObjTag = searchObjTagLast
      }
      if (!type) {
        this.pageNum = 1
      }
      this.getTable(this.nodeTitle)
    },
    // 搜索tag移除、
    handleCloseTag(item) {
      this.searchObj[item.enName] = ''
      this.searchTable()
    },
    // 模糊搜索
    batchSearchNameTable() {
      this.pageNum = 1
      this.pageSize = 10
      this.getTable(this.nodeTitle)
    },
    // 列表选择
    handleSelectionChangeTol(val) {
      this.batchHostarr = val
    },
    batchReasseP(nodeCreate, nodename, title, status) {
      if (this.batchHostarr.length == 0) {
        this.$message({
          showClose: true,
          message: '请选择列表！',
          type: 'warning',
        })
        return
      }
      let arrIdc = []
      for (var k in this.batchHostarr) {
        if (
          arrIdc.indexOf(this.batchHostarr[k].account.id) == -1 &&
          this.batchHostarr[k].account.id
        ) {
          arrIdc.push(this.batchHostarr[k].account.id)
        }
      }
      if (arrIdc.length != 1) {
        this.$message({
          showClose: true,
          message: '请选择同一个账号下列表！',
          type: 'warning',
        })
        return
      }
      if (status) {
        let arrPro = []
        for (var k in this.batchHostarr) {
          if ( arrPro.indexOf(this.batchHostarr[k].project_config.id) == -1 && this.batchHostarr[k].project_config.id) {
            arrPro.push(this.batchHostarr[k].project_config.id)
          }
        }
        if (arrPro.length != 1) {
          this.$message({
            showClose: true,
            message: '请选择同一个项目下列表！',
            type: 'warning',
          })
          return
        }
      }
      this.rebootBtnNoProject(nodeCreate, nodename, this.batchHostarr, title)
    },
    batchReassey(nodeCreate, nodename, title) {
      this.rebootBtnNoProject(nodeCreate, nodename, this.batchHostarr, title)
    },
    batchReassembly(nodeCreate, nodename, title) {
      if (this.batchHostarr.length == 0) {
        this.$message({
          showClose: true,
          message: '请选择列表！',
          type: 'warning',
        })
        return
      }
      var arr = {}
      var arr1 = {}
      for (var k in this.batchHostarr) {
        if (!this.batchHostarr[k].project_config) {
          this.$message({
            showClose: true,
            message:
              this.batchHostarr[k].name + ' 所属项目为空，请联系管理员！',
            type: 'warning',
          })
          return
        }
        arr[this.batchHostarr[k].project_config.id] = ''
        arr1[this.batchHostarr[k].account.alias] = ''
      }
      if (Object.keys(arr).length != 1 || Object.keys(arr1).length != 1) {
        this.$message({
          showClose: true,
          message: '所选列表非同一个项目组或账号！',
          type: 'warning',
        })
        return
      }
      this.rebootBtn(nodeCreate, nodename, this.batchHostarr, title)
    },
    rebootBtn(nodeCreate, nodename, rowArr, title) {
      if (!rowArr[0].project_config) {
        this.$message({
          showClose: true,
          message: '所属项目为空，请联系管理员！',
          type: 'warning',
        })
        return
      }
      this.nodeCheckName = nodeCreate
      this.nodeCheckProductName = nodeCreate
      this.dialogtitle = title
      this.newTaskName = title
      this.childNodename = nodename
      this.checkList = rowArr
      this.reviewStatus222 = true
    },
    rebootBtnNoProject(nodeCreate, nodename, rowArr, title) {
      if (rowArr.length == 0) {
        this.$message({
          showClose: true,
          message: '请选择列表！',
          type: 'warning',
        })
        return
      }
      this.nodeCheckName = nodeCreate
      this.nodeCheckProductName = nodeCreate
      this.dialogtitle = title
      this.newTaskName = title
      this.childNodename = nodename
      this.checkList = rowArr
      this.reviewStatus222 = true
    },
    rebootBtnNew(nodeCreate, title, obj) {
      this.childNodename = ''
      this.checkList = []
      this.nodeCheckName = nodeCreate
      this.newTaskName = title
      this.nodeCheckProductName = nodeCreate
      this.dialogtitle = title
      this.reviewStatus222 = true
    },
    // 发起任务弹框关闭
    addClose(val) {
      this.reviewStatus222 = val
      this.searchTable(false)
      this.historyTableList = {}
    },
    // 标签管理弹框关闭
    addCloseTag(val) {
      this.tagStatus = val
      this.getTreeData()
    },
    resetForm(formName) {
      if (this.$refs[formName]) {
        this.$refs[formName].resetFields()
      }
    },
    tableRowClassName({ row, rowIndex }) {
      if (this.tableRowIndex !== null && rowIndex == this.tableRowIndex) {
        return 'fixed_row'
      }
    },
    // 查看详情
    reviewDetail(obj, title, pobj) {
      this.tableRowIndex = pobj.$index // 获取当前点击行下标　　　 // ... 点击当前行，进行的操作    },
      this.tableDetailStatus = sessionStorage.getItem('tableDetailStatusC')
        ? sessionStorage.getItem('tableDetailStatusC')
        : 'halfScreen'
      this.tablePagina = document.getElementById('tablePagina')
      this.tableHight = ''
      if (this.tableDetailStatus == 'halfScreen') {
        this.tableHight = 'up'
        this.tablePagina.style.textAlign = 'center'
        this.paginationSatatus = false
      } else if (this.tableDetailStatus == 'halfScreenBan') {
        this.$nextTick(() => {
          this.tablePagina.style.textAlign = 'left'
          this.paginationSatatus = true
        })
      } else {
        this.tablePagina.style.textAlign = 'center'
        this.paginationSatatus = false
      }
      sessionStorage.setItem('tableDetailStatusC', this.tableDetailStatus)
      this.jobIds = obj.id
      this.drawDetail = true
      this.$router.replace({
        path: location.pathname,
        query: {
          detail_ID: this.jobIds,
          activeName: this.$route.query.activeName,
        },
      })
    },
    reviewDetailTo(obj, num) {
      this.jobIds = obj.id
      this.drawDetail = true
      this.$router.replace({
        path: location.pathname,
        query: {
          detail_ID: obj.id,
          activeName: num,
        },
      })
    },
    //link跳转打开页面
    linkOpen(linkId, link) {
      window.open(`${link}?detail_ID=${linkId}`)
    },
  },
}

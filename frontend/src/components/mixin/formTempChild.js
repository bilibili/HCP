import Http from "@/components/api/services";
export const formTempChild = {
  //方法集合
  methods: {
    clearableRadioBtn(row, rowindex, result, cobj, iii, formValueT, key) {
      var value = row.paramList ? [] : ''
      formValueT[key] = value
      this.tureOrFalseCoant(value, rowindex, result, cobj, iii)
      let thisInfoForm = JSON.parse(JSON.stringify(this.childObject))
      this.$emit('emitChilddelOne', { a: thisInfoForm, d: this.parentName, e: this.smallName, aiii: this.aiii, waiIndex: this.waiIndex })
    },
    // 服务器多选
    clickServerBatch(formData, forAtt) {
      forAtt[0].map((mitem, mi) => {
        if (mitem.paramremainOperand == mitem.paramreference) {
          this.choseServerList(mi, formData, forAtt, mitem, mitem.parentName, 0, mitem.paramTitle, true)
        }
      })
    },
    // 子表单样式表格
    choseServerBtnBathjoin() {
      let thisInfoForm = JSON.parse(JSON.stringify(this.childObject))
      let thisInfoAttr = JSON.parse(JSON.stringify(this.childData))
      let attrLineone = JSON.parse(JSON.stringify(thisInfoAttr[0]))
      let serverArr = []
      let paramremainOperandStr = attrLineone[0].paramremainOperand
      let showName = ""
      attrLineone.map((k, i) => {
        if (k.paramreference == paramremainOperandStr) {
          showName = k.paramName
        }
      })
      thisInfoForm.map((item, i) => {
        serverArr.push(item[showName])
      })
      this.choseListServer.map((citem, ci) => {
        console.log(citem)
        if (serverArr.indexOf(citem.id) == -1) {
          let obj = {}
          attrLineone.map((oitem, oi) => {
            obj[oitem.paramName] = oitem.paramList ? [] : ''
            if (oitem.paramName == showName) {
              oitem.paramDataList = [{ name: citem.name, value: citem.id }]
              obj[oitem.paramName] = oitem.paramList ? [citem.id] : citem.id
            }
            if (oitem.paramName == "privateIPv4") {
              oitem.paramDataList = [{ name: citem.privateIPv4[0], value: citem.privateIPv4[0] }]
              obj[oitem.paramName] = oitem.paramList ? [citem.id] : citem.privateIPv4[0]
            }
          })
          if (thisInfoForm[0][showName]) {
            thisInfoForm.push(obj)
            thisInfoAttr.push(JSON.parse(JSON.stringify(attrLineone)))
          } else {
            thisInfoForm = [obj]
            thisInfoAttr[0] = JSON.parse(JSON.stringify(attrLineone))
          }
        } else {
          this.$message({
            showClose: true,
            message: "已添加过该对象",
            type: 'warning'
          })
        }
      })
      this.choseServerLog = false
      this.$emit('emitChildcopy', { a: thisInfoForm, c: this.index, b: thisInfoAttr, d: this.parentName, e: this.smallName, aiii: this.aiii, waiIndex: this.waiIndex })
    },
    checkShow(name, i) {
      let thisInfoForm = JSON.parse(JSON.stringify(this.childObject))
      thisInfoForm[i]['_isShowStatus'] = !thisInfoForm[i]['_isShowStatus']
      this.$emit('emitChilddelOne', { a: thisInfoForm, c: i, d: this.parentName, e: this.smallName, aiii: this.aiii, waiIndex: this.waiIndex })
    },
    linkJumpTo(obj, row) {
      let link = obj.paramreferenceLink.replace('{$id}', row.id)
      window.open(`${link}`)
    },
    choseServerList(index, formData, forAtt, obj, title, bodyIndex, serverTitle, type) {
      if (type) {
        this.serverLogStatus = true
      } else {
        this.serverLogStatus = false
      }
      let thisInfoForm = JSON.parse(JSON.stringify(this.childObject))
      let thisInfoAttr = JSON.parse(JSON.stringify(this.childData))
      this.serverTitle = serverTitle
      this.searchFormServe = ''
      this.batchSearchName = ''
      this.serverTypelineIndex = index
      this.serverTypeIndex = bodyIndex
      let that = this
      // 如果是弹框选择列表，找到要定义数据位置
      this.childRowTitle1 = title
      var objData = forAtt[bodyIndex][index]
      this.mainOperandName = objData.paramremainOperand
      this.mainOperandvsrefrename = objData.paramreference
      if (!objData.paramreferenceQuery) {
        objData.paramreferenceQuery = "{}"
      }
      this.searchpageNum = 1
      var strA = objData.paramreferenceQuery;
      var username = this.usernameAll
      var firstStr = '"$';
      var secondStr = '"';
      var newStr11 = ""
      function getNum(str, firstStr, secondStr) {
        if (str == "" || str == null || str == undefined) {
          return "";
        }
        if (str.indexOf(firstStr) < 0) {
          newStr11 = str
          return "";
        }
        var subFirstStr = str.substring(str.indexOf(firstStr) + firstStr.length, str.length);
        var subSecondStr = subFirstStr.substring(0, subFirstStr.indexOf(secondStr));
        var formDataValue = ""
        if (subSecondStr != 'USER') {
          if (formData[subSecondStr]) {
            if (Object.prototype.toString.call(formData[subSecondStr]).indexOf("Array") != -1) {
              let abc = []
              formData[subSecondStr].map((k111, i111) => {
                if (Object.prototype.toString.call(k111).indexOf("Object") == -1) {
                  abc.push(k111)
                } else {
                  abc.push(k111.id ? k111.id : k111.oid)
                }
              })
              formDataValue = JSON.stringify(abc)
            } else {
              formDataValue = formData[subSecondStr] ? formData[subSecondStr] : 'ISvalue_null'
            }
          } else if (that.erinfoClidObj[subSecondStr]) {
            if (Object.prototype.toString.call(that.erinfoClidObj[subSecondStr]).indexOf("Array") != -1) {
              let abc = []
              that.erinfoClidObj[subSecondStr].map((k111, i111) => {
                if (Object.prototype.toString.call(k111).indexOf("Object") == -1) {
                  abc.push(k111)
                } else {
                  abc.push(k111.id ? k111.id : k111.oid)
                }
              })
              formDataValue = JSON.stringify(abc)
            } else {
              formDataValue = that.erinfoClidObj[subSecondStr] ? that.erinfoClidObj[subSecondStr] : 'ISvalue_null'
            }
          } else if (that.oneinfoClidObj[subSecondStr]) {
            if (Object.prototype.toString.call(that.oneinfoClidObj[subSecondStr]).indexOf("Array") != -1) {
              let abc = []
              that.oneinfoClidObj[subSecondStr].map((k111, i111) => {
                if (Object.prototype.toString.call(k111).indexOf("Object") == -1) {
                  abc.push(k111)
                } else {
                  abc.push(k111.id ? k111.id : k111.oid)
                }
              })
              formDataValue = JSON.stringify(abc)
            } else {
              formDataValue = that.oneinfoClidObj[subSecondStr] ? that.oneinfoClidObj[subSecondStr] : 'ISvalue_null'
            }

          } else {
            formDataValue = 'ISvalue_null'
          }
        } else {
          formDataValue = username
        }
        str = str.replace("$" + subSecondStr, formDataValue)
        if (str.indexOf("$") != -1) {
          getNum(str, firstStr, secondStr)
        } else {
          newStr11 = str
        }
      }
      getNum(strA, firstStr, secondStr)
      if (newStr11.indexOf("ISvalue_null") != -1) {
        return false
      }
      newStr11 = newStr11.replace(/\"\[/g, '[').replace(/\]\"/g, ']')
      var queryobj = JSON.parse(newStr11)
      this.choseServerLog = true
      this.searchHostTable = []
      this.searchtotalLen = 0
      this.choseListServer = []
      objData.paramDataList = []
      let referenceShowArrHeader11 = []
      this.referenceShowArrHeader = []
      this.queryObjectSearch = queryobj
      this.nameWhereSearchServe = `${objData.paramreference}Where`
      this.nameWhereSearchServeOption = `${objData.paramreference}Option`
      this.queryNameSearchServe = objData.paramreference
      this.ObjDataRowServer = objData
      let paramreferenceShowNew = []
      this.paramreferenceShowNew1 = []
      let displayShow = obj.paramreferenceDisplay ? JSON.parse(obj.paramreferenceDisplay) : { properties: [] }
      displayShow.properties.map((ditem, dii) => {
        if (ditem.name != 'oid') {
          paramreferenceShowNew.push(ditem.name)
        } else {
          paramreferenceShowNew.push('id')
        }
        if (ditem.index) {
          this.paramreferenceShowNew1.push(ditem.name)
        }
      })
      this.paramreferenceShowNew2 = paramreferenceShowNew
      paramreferenceShowNew.map((k, i) => {
        referenceShowArrHeader11.push({
          title: k,
          value: k
        })
      })
      for (let rekey in displayShow.properties) {
        referenceShowArrHeader11.map(reItem => {
          let title = reItem.title
          if (displayShow.properties[rekey].name == reItem.value && displayShow.properties[rekey].display) {
            this.referenceShowArrHeader.push({
              title: displayShow.properties[rekey].title,
              value: displayShow.properties[rekey].name
            })
          }
        })
      }
      this.referenceShowArrHeader = JSON.parse(JSON.stringify(this.referenceShowArrHeader))
      thisInfoAttr[bodyIndex][index].paramTableTitle = this.referenceShowArrHeader
      this.$emit('emitChildcopy', { a: thisInfoForm, c: this.index, b: thisInfoAttr, d: this.parentName, e: this.smallName, aiii: this.aiii, waiIndex: this.waiIndex })
      //   获取服务器列表接口方法
      this.searchTableServer()
    },
    // 选择服务器列表接口
    searchTableServer(type) {
      if (type) {
        this.searchpageNum = 1
      }
      var postData = {
        schema: this.queryNameSearchServe,
        page_size:this.searchpagesize,
        page_num:this.searchpageNum,
        where:{}
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
          postData.where = {
            "or":or
          }
        }else{
          postData.where = {
            "name_REGEX":this.searchFormServe.trim()
          }
        }
      }
      if (this.batchSearchName.replace(/\s/g, '')) {
        var batcharr = []
        var batchNewarr = []
        batcharr = this.batchSearchName.split('\n')
        batcharr.map((k, i) => {
          if (k.replace(/\s/g, '')) {
            batchNewarr.push(k.replace(/\s/g, ''))
          }
        })
        postData.where = { "name_IN": batchNewarr }
      } 
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
    // 分页
    changehandleCurrent(val) {
      this.searchpageNum = val;
      this.searchTableServer()
    },
    handleSizeChange(val) {
      this.searchpagesize = val
      this.searchpageNum = 1;
      this.searchTableServer(true)
    },
    // 确定选择服务器按钮
    choseServerBtn() {
      let thisInfoForm = JSON.parse(JSON.stringify(this.childObject))
      let thisInfoAttr = JSON.parse(JSON.stringify(this.childData))
      if (!thisInfoForm[this.serverTypeIndex][this.childRowTitle1]) {
        thisInfoForm[this.serverTypeIndex][this.childRowTitle1] = []
      }
      var arr = []
      thisInfoForm.map((k, i) => {
        if (k[this.childRowTitle1] && Object.prototype.toString.call(k[this.childRowTitle1]).indexOf("Array") != -1) {
          k[this.childRowTitle1].map((k2, i2) => {
            if (k2.oid) {
              arr.push(k2.oid)
            } else {
              arr.push(k2.id)
            }
          })
        }
      })
      this.choseListServer.map((k, i) => {
        if (this.mainOperandName == this.mainOperandvsrefrename) {
          if (arr.indexOf(k.id) == -1) {
            thisInfoForm[this.serverTypeIndex][this.childRowTitle1].push(k)
          } else {
            this.$message({
              showClose: true,
              message: "已添加过该对象",
              type: 'warning'
            })
          }
        } else {
          thisInfoForm[this.serverTypeIndex][this.childRowTitle1].push(k)
        }
      })

      if (thisInfoForm[this.serverTypeIndex][this.childRowTitle1].length) {
        thisInfoAttr[this.serverTypeIndex].map((k, i) => {
          if (k.paramName == this.childRowTitle1) {
            k.colorRed = "1"
          }
        })
      }

      this.$emit('emitChildcopy', { a: thisInfoForm, c: this.serverTypelineIndex, b: thisInfoAttr, d: this.parentName, e: this.smallName, aiii: this.aiii, waiIndex: this.waiIndex })
      this.choseServerLog = false
    },
    // 选择服务器
    handleSelectionChangeSearth(val) {
      this.choseListServer = val
    },
    // 是否关联必填 如果关联，会有一个选择  是就展示 不是隐藏
    tureOrFalseCoant(value, rowindex, result, cobj, iii) {
      let childDataNew = JSON.parse(JSON.stringify(this.childData))
      let childObjectNew = JSON.parse(JSON.stringify(this.childObject))
      if (cobj) {
        let ii = Number(this.waiIndex)
        if (value == "") {
          if (this.newTaskName != '资源模板新建' && this.newTaskName != '资源模板编辑') {
            childDataNew[rowindex][iii].colorRed = '5'
            childDataNew[rowindex][iii].customDescription = '不可为空'
          }
        } else {
          childDataNew[rowindex][iii].colorRed = '1'
        }
        for (var Ai = iii + 1; Ai < childDataNew[rowindex].length; Ai++) {
          if (childDataNew[rowindex][Ai].paramreferenceQuery) {
            if (childDataNew[rowindex][Ai].paramreferenceQuery.indexOf(`$${childDataNew[rowindex][iii].paramName}`) != -1) {
              this.focusChildSelect(rowindex, Ai, childObjectNew, childDataNew[rowindex][Ai], childDataNew[rowindex][iii].paramName)
              if (childDataNew[rowindex][Ai].paramList) {
                childDataNew[[rowindex][iii].paramName + "ChildObj3"][rowindex][childDataNew[rowindex][Ai].paramName] = []
              } else {
                childDataNew[[rowindex][iii].paramName + "ChildObj3"][rowindex][childDataNew[rowindex][Ai].paramName] = ''
              }
              childDataNew[rowindex][Ai].paramDataList = []
            }
          }
          if (childDataNew[rowindex][Ai].paramreferenceQuery) {
            if (childDataNew[rowindex][Ai].paramreferenceQuery.indexOf(`\"$${childDataNew.paramName}\"`) != -1) {
              this.focusSelectValue({ row: childDataNew[rowindex][Ai] }, childDataNew[rowindex][Ai].paramName, childObjectNew, Ai)
              childObjectNew[childDataNew[rowindex][Ai].paramName] = ''
            }
          }
        }
      } else {
        let index = this.waiIndex
        if (value == "") {
          if (this.newTaskName != '资源模板新建' && this.newTaskName != '资源模板编辑') {
            childDataNew[rowindex][iii].colorRed = '5'
            childDataNew[rowindex][iii].customDescription = '不可为空'
          }
        } else {
          childDataNew[rowindex][iii].colorRed = '1'
        }
        for (var Ai = index + 1; Ai < childDataNew.length; Ai++) {
          if (childDataNew[Ai].paramreferenceQuery) {
            if (childDataNew[Ai].paramreferenceQuery.indexOf(`$${childDataNew.paramName}`) != -1) {
              this.focusSelectValue({ row: childDataNew[Ai] }, childDataNew[Ai].paramName, childObjectNew, Ai)
              childObjectNew[childDataNew[Ai].paramName] = ''
              childDataNew[Ai].paramDataList = []
            }
          }
        }
      }
      childDataNew.map((achItem, achItemi) => {
        achItem.map((aaaItem, aaaItemi) => {
          if (achItemi == rowindex) {
            if (aaaItem.paramShemas) {
              let shemasObjectch = JSON.parse(aaaItem.paramShemas)
              var bgif = "2"
              var a = 0
              for (let meItemch in shemasObjectch) {
                if (cobj.paramName == meItemch) {
                  if (Object.prototype.toString.call(shemasObjectch[meItemch]).indexOf("Array") != -1) {
                    shemasObjectch[meItemch].map((k222, i222) => {
                      if (Object.prototype.toString.call(result).indexOf("Array") != -1) {
                        if (result.indexOf(k222) != -1) {
                          a += 1
                        }
                      } else {
                        if (result == k222) {
                          a += 1
                        }
                      }

                    })
                  } else {
                    if (shemasObjectch[meItemch] == result) {
                      a += 1
                    }
                  }
                } else {
                  if (childObjectNew[achItemi][meItemch]) {
                    if (shemasObjectch[meItemch].indexOf(childObjectNew[achItemi][meItemch]) != -1) {
                      a += 1
                    }
                  }
                }
              }
              if (a == Object.keys(shemasObjectch).length) {
                bgif = "1"
              }
              if (aaaItem.paramavailableCondition == 'role:admin') {
                if (this.users.indexOf(this.usernameAll) != -1) {
                  aaaItem.thisShowIf = bgif
                } else {
                  aaaItem.thisShowIf = "2"
                }
              } else if (aaaItem.paramavailableCondition && aaaItem.paramavailableCondition.indexOf('groupRole:') != -1) {
                var arrGroup = aaaItem.paramavailableCondition.split("groupRole:")[1].split(",")
                aaaItem.thisShowIf = "2"
                arrGroup.map((gk1, gi1) => {
                  if (this.manageUserGroup.indexOf(gk1) != -1) {
                    aaaItem.thisShowIf = bgif
                  }
                })
              } else {
                aaaItem.thisShowIf = bgif
              }
              if (bgif == 1) {
                if (aaaItem.paramType == "number") {
                  childObjectNew[achItemi][aaaItem.paramName] = childObjectNew[achItemi][aaaItem.paramName] ? childObjectNew[achItemi][aaaItem.paramName] : Number(aaaItem.paramDefault)
                }
              } else {
                if (aaaItem.paramList) {
                  childObjectNew[achItemi][aaaItem.paramName] = []
                } else {
                  childObjectNew[achItemi][aaaItem.paramName] = ''
                }
              }
            }
          }
        })
      })
      let parentNameStr = this.parentName
      childDataNew = JSON.parse(JSON.stringify(childDataNew))
      this.$emit('emitChild', { a: childObjectNew, b: childDataNew, c: rowindex, d: parentNameStr, e: this.smallName, aiii: this.aiii, waiIndex: this.waiIndex })
    },
    // 获取下拉菜单
    getMetaSchemasListval(name, obj, val, formData, pIndex, index) {

      if (!name) {
        return
      }
      var postData = {
        schema: name,
        page_size:100,
        page_num:1,
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
          childDataNew[pIndex][index].paramDataList = arrList
          this.$emit('emitChildvalidating', { attr: childDataNew, aiii: this.aiii, waiIndex: this.waiIndex })
        })
    },
    // 下拉搜索
    filterMethod(val, obj, formData, pIndex, index) {
      this.getstr = val
      this.getMetaSchemasListval(obj.paramreference, obj, val, formData, pIndex, index)
    },
    // 聚焦搜索
    getMetaSchemasListvalFocus(val, obj, formData, pIndex, index) {
      this.getstr = val
      this.getMetaSchemasListval(obj.paramreference, obj, val, formData, pIndex, index)
    },
    deloneline(pinedex, name, index, scindex) {
      let childObjectNew = JSON.parse(JSON.stringify(this.childObject))
      childObjectNew[pinedex][name].splice(scindex.row.$index, 1)
      this.$emit('emitChilddelOne', { a: childObjectNew, c: index, d: this.parentName, e: this.smallName, aiii: this.aiii, waiIndex: this.waiIndex })
    },
    // 下拉选项
    focusChildSelect(bodyIndex, index, formData, ShowTodisplayObj, pname, status) {
      let that = this
      let childDataNew = JSON.parse(JSON.stringify(this.childData))
      let childObjectNew = JSON.parse(JSON.stringify(this.childObject))
      var objData = childDataNew[bodyIndex][index]
      var postData = {
        schema: objData.paramreference,
        page_size:100,
        page_num:1,
      }
      objData.paramDataList = []
      Http.getQueryList(postData).then((response) => {
        objData.paramDataList = []
        response.data.data.data.map((k,i)=>{
          objData.paramDataList.push({
            name: k.name,
            value: k.id
          })
        })
        this.$emit('emitChildcopy', { a: childObjectNew, b: childDataNew, c: index, d: this.parentName, e: this.smallName, aiii: this.aiii, waiIndex: this.waiIndex })       
      })
    },
    // 下拉选择
    changeChildSelect(value, index, name, formData, formAttr, bodyIndex) {
      if (value == "") {
        if (this.newTaskName != '资源模板新建' && this.newTaskName != '资源模板编辑') {
          formAttr[bodyIndex][index].colorRed = '5'
          formAttr[bodyIndex][index].customDescription = '不可为空'
        }
      } else {
        formAttr[bodyIndex][index].colorRed = '1'
      }
      formAttr[bodyIndex].map((k, i) => {
        if (i > index) {
          if (k.paramreferenceQuery) {
            if (k.paramreferenceQuery.indexOf(`$${formAttr[bodyIndex][index].paramName}`) != -1) {
              if (k.paramList) {
                formData[bodyIndex][k.paramName] = []
              } else {
                formData[bodyIndex][k.paramName] = ''
              }
              k.paramDataList = []
              this.focusChildSelect(bodyIndex, i, formData, k, k.paramName)
            }
          }
        }
      })

    },
    focusChildGLArr(obj, pcindex, nrindex, name) {
      let arr = obj.paramprerequisite
      let dataList = []
      this.erinfoChildattr.map((kk, ii) => {
        if (this.aiii.erindex > ii) {
          if (arr.indexOf(kk[0].paramFromName) != -1) {
            //判断自身下标不能push
            dataList.push({
              name: `批次${ii + 1} - ${kk[0].paramFromTitle}`,
              value: `批次${ii + 1} - ${kk[0].paramFromTitle}`
            })
          }
        }
      })
      this.$emit('emitChildfocusChildGLArr', { dataList: dataList, waiIndex: this.waiIndex, aiii: this.aiii, pcindex: pcindex, nrindex: nrindex })
    },
    // 验证状态子表单
    validatingchild(value, rowindex, cindex, minLen, maxLen, pattern, type, formDatathis) {
      let childDataNew = JSON.parse(JSON.stringify(this.childData))
      let reg = new RegExp(pattern)
      if (value == "" && childDataNew[rowindex][cindex].paramRequired) {
        if (this.newTaskName != '资源模板新建' && this.newTaskName != '资源模板编辑') {
          childDataNew[rowindex][cindex].colorRed = '5'
          childDataNew[rowindex][cindex].customDescription = '不可为空'
        }
      } else {
        if (value) {
          childDataNew[rowindex][cindex].colorRed = '1'
        }
        if (minLen != null && minLen != 'undefind') {
          if (type == 'string' && value.length < minLen) {
            childDataNew[rowindex][cindex].colorRed = '2'
            childDataNew[rowindex][cindex].customDescription = '最小长度' + minLen + '个字符'
          } else if (type == 'number' && value < minLen) {
            childDataNew[rowindex][cindex].colorRed = '4'
            childDataNew[rowindex][cindex].customDescription = '数值不得小于' + minLen
          } else {
            childDataNew[rowindex][cindex].colorRed = '1'
          }
        }
        if (maxLen != null && maxLen != 'undefind' && childDataNew[rowindex][cindex].colorRed == '1') {
          if (type == 'string' && value.length > maxLen) {
            childDataNew[rowindex][cindex].colorRed = '2'
            childDataNew[rowindex][cindex].colorRed = '4'
            childDataNew[rowindex][cindex].customDescription = '最大长度' + maxLen + '个字符'
          } else if (type == 'number' && value > maxLen) {
            childDataNew[rowindex][cindex].colorRed = '4'
            childDataNew[rowindex][cindex].colorRed = '4'
            childDataNew[rowindex][cindex].customDescription = '数值不得大于' + maxLen
          } else {
            childDataNew[rowindex][cindex].colorRed = '1'
          }
        }
        if (pattern && childDataNew[rowindex][cindex].colorRed == '1') {
          if (!reg.test(value)) {
            childDataNew[rowindex][cindex].colorRed = '3'
            childDataNew[rowindex][cindex].customDescription = '填写格式有误'
          } else {
            childDataNew[rowindex][cindex].colorRed = '1'
          }
        }
        if (!childDataNew[rowindex][cindex].paramRequired) {
          if (value == '') {
            childDataNew[rowindex][cindex].colorRed = '1'
          }
        }
      }
      if (formDatathis) {
        formDatathis.map((item, i) => {
          if (pattern && childDataNew[rowindex][cindex].colorRed == '1') {
            if (!reg.test(item)) {
              childDataNew[rowindex][cindex].colorRed = '3'
              childDataNew[rowindex][cindex].customDescription = '填写格式有误'
            } else {
              childDataNew[rowindex][cindex].colorRed = '1'

            }
          }
        })
      }
      this.$emit('emitChildvalidating', { attr: childDataNew, aiii: this.aiii, waiIndex: this.waiIndex })
    },
    // 复制表单
    copyCardAttr(obj, name, rowIndex) {
      let thisInfoForm = JSON.parse(JSON.stringify(this.childObject))
      let thisInfoAttr = JSON.parse(JSON.stringify(this.childData))
      let formOne = thisInfoForm[rowIndex]
      let AttrOne = thisInfoAttr[rowIndex]
      thisInfoForm.splice(rowIndex, 0, formOne)
      thisInfoAttr.splice(rowIndex, 0, AttrOne)
      this.$emit('emitChildcopy', { a: thisInfoForm, b: thisInfoAttr, c: rowIndex, d: this.parentName, e: this.smallName, aiii: this.aiii, waiIndex: this.waiIndex })
    },
    deletCardAttr(obj, name, rowIndex) {
      let thisInfoForm = JSON.parse(JSON.stringify(this.childObject))
      let thisInfoAttr = JSON.parse(JSON.stringify(this.childData))
      thisInfoForm.splice(rowIndex, 1)
      thisInfoAttr.splice(rowIndex, 1)
      this.$emit('emitChildcopy', { a: thisInfoForm, b: thisInfoAttr, c: rowIndex, d: this.parentName, e: this.smallName, aiii: this.aiii, waiIndex: this.waiIndex })
    },
    changebatchSearchInput(obj, name, str) {
      var batcharr = []
      var batchNewarr = []
      batcharr = str.split('\n')
      batcharr.map((k, i) => {
        if (k.replace(/\s/g, '')) {
          batchNewarr.push(k.replace(/\s/g, ''))
        }
      })
      if (!batchNewarr.length) return
      obj.map((k, i) => {
        k[name] = batchNewarr[i]
      })
      this.$emit('emitChilddelOne', { a: obj, d: this.parentName, e: this.smallName, aiii: this.aiii, waiIndex: this.waiIndex })
    },
    // 上传文件
    successUpload1(response, file, fileList, rowIndex, prowIndex, name) {
      let thisInfoForm = JSON.parse(JSON.stringify(this.childObject))
      let thisInfoAttr = JSON.parse(JSON.stringify(this.childData))
      if (response.code != 0) {
        this.$message({
          showClose: true,
          message: "文件上传失败：" + response.data,
          type: 'error'
        });
      }
      var arr = []
      fileList.map((k, i) => {
        if (k.response.code == 0) {
          arr.push({
            name: k.name,
            url: k.response.data,
            response: k.response
          })
        }
      })
      thisInfoForm[rowIndex][name] = arr
      thisInfoAttr[rowIndex][prowIndex].colorRed = "1"
      this.$emit('emitChild', { a: thisInfoForm, b: thisInfoAttr, c: rowIndex, d: this.paramName, e: this.smallName, aiii: this.aiii, waiIndex: this.waiIndex })
    },
    handleRemove1(file, fileList, rowIndex, name) {
      let thisInfoForm = JSON.parse(JSON.stringify(this.childObject))
      let thisInfoAttr = JSON.parse(JSON.stringify(this.childData))
      var arr = []
      fileList.map((k, i) => {
        if (k.response.code == 0) {
          arr.push({
            name: k.name,
            url: k.response.data,
            response: k.response
          })
        }
      })
      thisInfoForm[rowIndex][name] = arr
      this.$emit('emitChild', { a: thisInfoForm, b: thisInfoAttr, c: rowIndex, d: this.paramName, e: this.smallName, aiii: this.aiii, waiIndex: this.waiIndex })
    },
  }
};

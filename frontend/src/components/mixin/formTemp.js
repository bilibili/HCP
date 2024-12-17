import Http from "@/components/api/services";
export const formTemp = {
  //方法集合
  methods: {
    //  展示基本信息
    reviewParam(type) {
      this.lookType = type ? true : false
      if (this.infoReviewAttr.length == 0) {
        return false;
      }
      if (!type) {
        this.detailOrLookThis = false
      }
      this.searchHostTable = []
      this.referenceShowArrHeader = []
      this.choseListServer = []
      this.infoReviewAttr.map((k, i) => {
        if (k.paramShemas) {
          k.thisShowIf = "2"
        } else {
          if (k.paramavailableCondition == 'role:admin') {
            if (this.users.indexOf(this.usernameAll) != -1) {
              k.thisShowIf = "1"
            } else {
              k.thisShowIf = "2"
            }
          } else if (k.paramavailableCondition && k.paramavailableCondition.indexOf('groupRole:') != -1) {
            var arrGroup = k.paramavailableCondition.split("groupRole:")[1].split(",")
            k.thisShowIf = "2"
            arrGroup.map((gk1, gi1) => {
              if (this.manageUserGroup.indexOf(gk1) != -1) {
                k.thisShowIf = "1"
              }
            })
          }
          if (this.nodeCheckName == 'FormElasticConfigCloudServer' && (this.newTaskName == '弹性资源模板改配' || this.newTaskName == '弹性资源模板编辑') && k.paramName == 'servers') {
            k.thisShowIf = "2"
          }
        }
        if (k.paramType == 'object') {
          if (k.paramreferenceZ.length && (!k.paramreference || k.paramreference != 'FormTemplate')) {
            k.paramreferenceZ.map((k1, i1) => {
              k.paramDataList.push({
                name: k1, value: k1, newName: k1
              })
            })
          }
          if (!k.paramreferenceQuery) {
            if (k.paramreference == 'FormTemplate') {
              k.paramDataList = []
              var arr1 = []
              k.paramTemplates.map((kss, i) => {
                k.paramDataList.push({
                  name: kss.title,
                  value: kss.id
                })
                this.gaipeiArr.push({
                  name: kss.name,
                  value: kss.id
                })
                arr1.push({
                  name: kss.title,
                  value: kss.id
                })
              })
              if (!type) {
                if (k.paramDataList.length == 1 && !this.checkListLast.length) {
                  this.showChildAttr(k, k.paramName, k.paramDataList[0].value, i, 0, {}, '111')
                }
              }
            } else {
              if(k.paramName != 'formObject' && k.paramreference){
                this.getMetaSchemasList(k.paramreference, k, null, null, null, null, null, null, null, i)
              }
            }
          }
        }
        var a = 0
        if (this.checkListLast.length == 0) {
          if (k.paramDefault) {
            if (k.paramDefault.indexOf("$") != -1) {
              if (a == 0) {
                if (!(type && this.nodeCheckName == 'FormJoinCloudProjectConfig')) {
                  this.focusSelectDefaultValue({ row: k }, k.paramName, this.infoReviewForm, i)
                }
              }
              a += 1
            } else {
              if (k.paramList) {
                if (this.infoReviewForm[k.paramName].length == 0) {
                  this.infoReviewForm[k.paramName] = [k.paramDefault]
                }
              } else {
                if (this.infoReviewForm[k.paramName] === "") {
                  this.infoReviewForm[k.paramName] = k.paramDefault
                }
              }
            }
          } else {
            if (k.paramStyle == 'radio') {
              if (a == 0) {
                this.focusSelectDefaultValue({ row: k }, k.paramName, this.infoReviewForm, i)
              }
              a += 1
            }
          }
        }
        if (this.lookType) {
          if (!k.paramList) {
            this.infoReviewForm[k.paramName] = this.historyTableList[k.paramName]
          }
        }
        if (k.paramShemas) {
          let shemasObject = JSON.parse(k.paramShemas)
          let arrx = 0
          for (let meItem in shemasObject) {
            if (Object.prototype.toString.call(this.infoReviewForm[meItem]).indexOf("Array") != -1) {
              this.infoReviewForm[meItem].map((k1, i1) => {
                if (shemasObject[meItem].indexOf(k1) != -1) {
                  arrx += 1
                }
              })
            } else {
              if (shemasObject[meItem].indexOf(this.infoReviewForm[meItem]) != -1) {
                arrx += 1
              }
            }
          }
          if (arrx == Object.keys(shemasObject).length) {
            if (k.paramavailableCondition == 'role:admin') {
              if (this.users.indexOf(this.usernameAll) != -1) {
                k.thisShowIf = "1"
              } else {
                k.thisShowIf = "2"
              }
            } else if (k.paramavailableCondition && k.paramavailableCondition.indexOf('groupRole:') != -1) {
              var arrGroup = k.paramavailableCondition.split("groupRole:")[1].split(",")
              k.thisShowIf = "2"
              arrGroup.map((gk1, gi1) => {
                if (this.manageUserGroup.indexOf(gk1) != -1) {
                  k.thisShowIf = "1"
                }
              })
            } else {
              k.thisShowIf = "1"
            }
          } else {
            k.thisShowIf = "2"
          }
          if (k.thisShowIf == "2") {
            if (k.paramList) {
              this.infoReviewForm[k.paramName] = []
            } else {
              this.infoReviewForm[k.paramName] = ""
            }
          }
        }
       
        if ((this.nodeCheckName == 'FormRebootCloudServer' || this.nodeCheckName == 'FormReinstallCloudServer' || this.nodeCheckName == 'FormRenameBatchCloudServer'||this.nodeCheckName =='FormConfigCloudProduct'||this.nodeCheckName =='FormDeleteCloudProduct') && !this.lookType && this.checkListLast.length) {
          if (k.paramName == 'project' || k.paramName == 'project_config') {
            k.paramDataList = [{ value: this.checkListLast[0].project_config.id, name: this.checkListLast[0].project_config.name, newName: this.checkListLast[0].project_config.name }]
            this.infoReviewForm[k.paramName] = this.checkListLast[0].project_config.id
          }
          if (k.paramName == 'cloudAccount' || k.paramName == 'account') {
            k.paramDataList = [{ value: this.checkListLast[0].account.id, name: this.checkListLast[0].account.name, newName: this.checkListLast[0].account.name }]
            this.infoReviewForm[k.paramName] = this.checkListLast[0].account.id
          }
          if (k.paramName == 'formObjects') {
            this.infoReviewForm["formObjects"] = this.checkListLast
            var titleProData = JSON.parse(k.paramreferenceDisplay)
            this.listObj[k.paramName] = []
            titleProData.properties.map((k11,i11)=>{
              if(k11.display){
                this.listObj[k.paramName].push({
                  title:k11.title,
                  value:k11.name
                })
              }
            })
          }
          if(this.nodeCheckName == 'FormReinstallCloudServer'){
            if (k.paramName == 'region') {
              k.paramDataList = [{ value: this.checkListLast[0].region.id, name: this.checkListLast[0].region.name, newName: this.checkListLast[0].region.name }]
              this.infoReviewForm[k.paramName] = this.checkListLast[0].region.id
            }
            if (k.paramName == 'image') {
              k.paramDataList = [{ value: this.checkListLast[0].server_image.id, name: this.checkListLast[0].server_image.name, newName: this.checkListLast[0].server_image.name }]
              this.infoReviewForm[k.paramName] = this.checkListLast[0].server_image.id
            }
          }
        }
        if (k.paramStyle == 'upload') {
          this.infoReviewForm[k.paramName] = []
          if (Object.keys(this.historyTableList).length) {
            var arrFile = []
            if (this.historyTableList[k.paramName]) {
              var filedata = Object.prototype.toString.call(this.historyTableList[k.paramName]).indexOf("Array") != -1 ? this.historyTableList[k.paramName] : this.historyTableList[k.paramName].split(",")
              filedata.map((k, i) => {
                arrFile.push({
                  name: k,
                  url: k,
                  response: {
                    code: 0
                  }
                })
              })
            }
            this.historyTableList[k.paramName] = arrFile
          }
        }

        // 编辑回显
        let updateArr = [
          'FormUpdateCloudProvider',
          'FormAdminUpdateCloudAccount',
          'FormAdminUpdateUser',
          'FormUpdateCloudProjectAccountConfig',
          'FormUpdateCloudProjectConfig',
        ]
        if (updateArr.indexOf(this.nodeCheckName) != -1) {
          if (this.lookType) {
            if (k.paramName == 'formObject') {
              this.infoReviewForm[k.paramName] = this.historyTableList.id
              k.paramDataList = [
                {
                  value: this.historyTableList.id,
                  name: this.historyTableList.name,
                  newName: this.historyTableList.name,
                }
              ]
            } else {
              if (k.paramType == 'object') {
                if (this.historyTableList[k.paramName]) {
                  this.infoReviewForm[k.paramName] = this.historyTableList[k.paramName].oid
                  k.paramDataList = [
                    {
                      value: this.historyTableList[k.paramName].oid,
                      name: this.historyTableList[k.paramName].name,
                      newName: this.historyTableList[k.paramName].name,
                    }
                  ]
                }

              } else {
                if (k.paramType == 'boolean') {
                  this.infoReviewForm[k.paramName] = this.historyTableList[k.paramName] ? true : false
                } else {
                  this.infoReviewForm[k.paramName] = this.historyTableList[k.paramName]
                }
              }
            }
          } else {
            if (this.checkListLast.length) {
              if (k.paramName == 'formObject') {
                if(this.nodeCheckName=='FormUpdateCloudProjectAccountConfig'){
                  this.checkListLast[0].name = this.checkListLast[0].project_config.name
                }
                this.infoReviewForm[k.paramName] = k.paramList ? [this.checkListLast[0].id] : this.checkListLast[0].id
                k.paramDataList = [
                  {
                    value: this.checkListLast[0].id,
                    name: this.checkListLast[0].name,
                    newName: this.checkListLast[0].name,
                  }
                ]
              }else{
                if (k.paramType == 'object') {
                  if (k.paramList) {
                    if (k.paramStyle == 'tableStyle') {
                      this.infoReviewForm[k.paramName] = this.checkListLast[0][k.paramName]
                      if (k.paramreferenceDisplay) {
                        this.listObj[k.paramName] = []
                        JSON.parse(k.paramreferenceDisplay).properties.map((k222, i222) => {
                          if (k222.display) {
                            this.listObj[k.paramName].push({ title: k222.title, value: k222.name })
                          }
                        })
                      }
                    } else {
                      k.paramDataList = []
                      if (this.checkListLast[0][k.paramName]) {
                        this.checkListLast[0][k.paramName].map((k1, i1) => {
                          this.infoReviewForm[k.paramName].push(k1.id)
                          k.paramDataList.push({
                            value: k1.id,
                            name: k1.name,
                            newName: k1.name,
                          })
                        })
                      }
                    }
                  } else {
                    if (this.checkListLast[0][k.paramName]) {
                      this.infoReviewForm[k.paramName] = this.checkListLast[0][k.paramName].id
                      k.paramDataList = [
                        {
                          value: this.checkListLast[0][k.paramName].id,
                          name: this.checkListLast[0][k.paramName].name ? this.checkListLast[0][k.paramName].name : (this.checkListLast[0][k.paramName].path ? this.checkListLast[0][k.paramName].path : this.checkListLast[0][k.paramName].code),
                          newName: this.checkListLast[0][k.paramName].name ? this.checkListLast[0][k.paramName].name : (this.checkListLast[0][k.paramName].path ? this.checkListLast[0][k.paramName].path : ''),
                        }
                      ]
                    }
                  }
                } else {
                  if (k.paramType == 'string' && k.paramList) {
                    this.infoReviewForm[k.paramName] = []
                    if (this.checkListLast[0][k.paramName]) {
                      if (Object.prototype.toString.call(this.checkListLast[0][k.paramName]).indexOf("Array") != -1) {
                        this.infoReviewForm[k.paramName] = this.checkListLast[0][k.paramName]
                      } else {
                        this.infoReviewForm[k.paramName] = this.checkListLast[0][k.paramName].split(",")
                      }
                    }
                  } else {
                    if (k.paramType == 'boolean') {
                      this.infoReviewForm[k.paramName] = this.checkListLast[0][k.paramName] == 'true' ? true : false
                    } else {
                      if(k.paramName != 'sync_secret_key' &&  k.paramName != 'operate_secret_key'){
                        this.infoReviewForm[k.paramName] = this.checkListLast[0][k.paramName]
                        if (k.paramStyle == "description") {
                          k.paramDefault = this.checkListLast[0][k.paramName]
                        }
                      } 
                    }
                  }
                }
              }
            } 
          }
        }

        // 删除回显
        let delBatchArr = [
          "FormAdminDeleteUser",
          "FormAdminDeleteCloudAccount",
          "FormDeleteCloudProjectConfig",
        ]
        if (delBatchArr.indexOf(this.nodeCheckName) != -1) {
          if (k.paramName == 'formObject' && !this.lookType) {
            this.infoReviewForm[k.paramName] = []
            this.checkListLast.map((kvv, ivv) => {
              k.paramDataList.push({
                value: kvv.id,
                name: kvv.name,
                newName: kvv.name,
              })
              this.infoReviewForm[k.paramName].push(kvv.id)
            })
          }
        }

      })
      this.resetForm("infoAddForm")
    },
    //  确认创建按钮
    subAddModel(type, status) {
      if (this.jobIdStatus && this.jobIdStatus != '审核中' && !this.historyTableList._copyOrupdate && !type && !status) {
        this.$message({
          showClose: true,
          message: "非审核中的任务禁止更新！",
          type: 'warning'
        });
        return false
      }
      // infoReviewForm
      let inputData = {
      }
      //如果有子表单则定义新对象，以此往里面赋值
      let inputdataapent = {}
      let inputdataapent1 = {}
      //如果有子表单则定义新对象
      let newParArr = []
      let paramNames = ''
      let joinTitle = ''
      let joinTitle1 = ''
      var infoReviewAttrPost = []
      var infoReviewFormPost = JSON.parse(JSON.stringify(this.infoReviewForm))
      var infoReviewAttrPost1 = JSON.parse(JSON.stringify(this.infoReviewAttr))
      for (let key in infoReviewAttrPost1) {
        if (infoReviewAttrPost1[key].paramRequired) {
          if ((infoReviewAttrPost1[key].paramType == 'string' || infoReviewAttrPost1[key].paramType == 'number') && infoReviewAttrPost1[key].paramStyle == 'description') {
            infoReviewFormPost[infoReviewAttrPost1[key].paramName] = infoReviewAttrPost1[key].paramDefault
          } else {
            if (infoReviewFormPost[infoReviewAttrPost1[key].paramName] == "" || infoReviewFormPost[infoReviewAttrPost1[key].paramName] == null) {
              infoReviewAttrPost1[key].colorRed = '5'
              infoReviewAttrPost1[key].customDescription = '不可为空'
            }
          }
        } else {
          if (infoReviewFormPost[infoReviewAttrPost1[key].paramName] == "" || infoReviewFormPost[infoReviewAttrPost1[key].paramName] == null) {
            infoReviewAttrPost1[key].colorRed = '1'
            infoReviewAttrPost1[key].customDescription = ''
          }
        }
        if (infoReviewAttrPost1[key].childAttrArr) {
          let secondName = infoReviewAttrPost1[key].paramName
          for (let inforViewItem1 in infoReviewAttrPost1[key].childAttrArr) {
            let inforViewClidItem1 = infoReviewAttrPost1[key].childAttrArr[inforViewItem1]
            for (let IIItem1 in inforViewClidItem1) {
              if (inforViewClidItem1[IIItem1].paramRequired) {
                if ((inforViewClidItem1[IIItem1].paramType == 'string' || inforViewClidItem1[IIItem1].paramType == 'number') && inforViewClidItem1[IIItem1].paramStyle == 'description') {
                  infoReviewFormPost[inforViewClidItem1[IIItem1].paramName] = inforViewClidItem1[IIItem1].paramDefault
                } else {
                  if (infoReviewFormPost[`${secondName}ChildObj`][inforViewItem1][inforViewClidItem1[IIItem1].paramName] == "" || infoReviewFormPost[`${secondName}ChildObj`][inforViewItem1][inforViewClidItem1[IIItem1].paramName] == null) {
                    infoReviewAttrPost1[key].childAttrArr[inforViewItem1][IIItem1].colorRed = '5'
                    infoReviewAttrPost1[key].childAttrArr[inforViewItem1][IIItem1].customDescription = '不可为空'
                  }
                }
              } else {
                if (infoReviewFormPost[`${secondName}ChildObj`][inforViewItem1][inforViewClidItem1[IIItem1].paramName] == "" || infoReviewFormPost[`${secondName}ChildObj`][inforViewItem1][inforViewClidItem1[IIItem1].paramName] == null) {
                  infoReviewAttrPost1[key].childAttrArr[inforViewItem1][IIItem1].colorRed = '1'
                  infoReviewAttrPost1[key].childAttrArr[inforViewItem1][IIItem1].customDescription = ''
                }
              }
              if (infoReviewAttrPost1[key].childAttrArr[inforViewItem1][IIItem1].childAttrArrThird) {
                let thirdName = infoReviewAttrPost1[key].childAttrArr[inforViewItem1][IIItem1].paramName
                for (let inforViewItem2 in infoReviewAttrPost1[key].childAttrArr[inforViewItem1][IIItem1].childAttrArrThird) {
                  let inforViewClidItem2 = infoReviewAttrPost1[key].childAttrArr[inforViewItem1][IIItem1].childAttrArrThird[inforViewItem2]
                  for (let IIItem2 in inforViewClidItem2) {
                    if (inforViewClidItem2[IIItem2].paramRequired) {
                      if ((inforViewClidItem2[IIItem2].paramType == 'string' || inforViewClidItem2[IIItem2].paramType == 'number') && inforViewClidItem2[IIItem2].paramStyle == 'description') {
                        infoReviewFormPost[`${secondName}ChildObj`][inforViewItem1][`${thirdName}ChildObj3`][inforViewItem2][inforViewClidItem2[IIItem2].paramName] = inforViewClidItem2[IIItem2].paramDefault
                      } else {
                        if (infoReviewFormPost[`${secondName}ChildObj`][inforViewItem1][`${thirdName}ChildObj3`][inforViewItem2][inforViewClidItem2[IIItem2].paramName] == "" || infoReviewFormPost[`${secondName}ChildObj`][inforViewItem1][`${thirdName}ChildObj3`][inforViewItem2][inforViewClidItem2[IIItem2].paramName] == null) {
                          infoReviewAttrPost1[key].childAttrArr[inforViewItem1][IIItem1].childAttrArrThird[inforViewItem2][IIItem2].colorRed = '5'
                          infoReviewAttrPost1[key].childAttrArr[inforViewItem1][IIItem1].childAttrArrThird[inforViewItem2][IIItem2].customDescription = '不可为空'
                        }
                      }
                    } else {
                      if (infoReviewFormPost[`${secondName}ChildObj`][inforViewItem1][`${thirdName}ChildObj3`][inforViewItem2][inforViewClidItem2[IIItem2].paramName] == "" || infoReviewFormPost[`${secondName}ChildObj`][inforViewItem1][`${thirdName}ChildObj3`][inforViewItem2][inforViewClidItem2[IIItem2].paramName] == null) {
                        infoReviewAttrPost1[key].childAttrArr[inforViewItem1][IIItem1].childAttrArrThird[inforViewItem2][IIItem2].colorRed = '1'
                        infoReviewAttrPost1[key].childAttrArr[inforViewItem1][IIItem1].childAttrArrThird[inforViewItem2][IIItem2].customDescription = ''
                      }
                    }
                  }
                }
              }
            }
          }
        }
      }
      this.infoReviewAttr = JSON.parse(JSON.stringify(infoReviewAttrPost1))
      for (let keys in infoReviewAttrPost1) {
        if (infoReviewAttrPost1[keys].childAttrArr) {
          for (let key1 in infoReviewAttrPost1[keys].childAttrArr) {
            let secondName = infoReviewAttrPost1[keys].paramName
            var ArraTTr = []
            for (let key2 in infoReviewAttrPost1[keys].childAttrArr[key1]) {
              if (infoReviewAttrPost1[keys].childAttrArr[key1][key2].thisShowIf == "1") {
                if (!(!infoReviewAttrPost1[keys].childAttrArr[key1][key2].paramRequired && (infoReviewFormPost[infoReviewAttrPost1[keys].paramName + "ChildObj"][key1][infoReviewAttrPost1[keys].childAttrArr[key1][key2].paramName] === '' || infoReviewFormPost[infoReviewAttrPost1[keys].paramName + "ChildObj"][key1][infoReviewAttrPost1[keys].childAttrArr[key1][key2].paramName] === undefined || infoReviewFormPost[infoReviewAttrPost1[keys].paramName + "ChildObj"][key1][infoReviewAttrPost1[keys].childAttrArr[key1][key2].paramName] === null || infoReviewFormPost[infoReviewAttrPost1[keys].paramName + "ChildObj"][key1][infoReviewAttrPost1[keys].childAttrArr[key1][key2].paramName].length == 0 || !infoReviewFormPost[infoReviewAttrPost1[keys].paramName + "ChildObj"][key1][infoReviewAttrPost1[keys].childAttrArr[key1][key2].paramName])) || infoReviewAttrPost1[keys].childAttrArr[key1][key2].paramreference == 'FormTemplate') {
                  ArraTTr.push(infoReviewAttrPost1[keys].childAttrArr[key1][key2])
                } else {
                  delete infoReviewFormPost[infoReviewAttrPost1[keys].paramName + "ChildObj"][key1][infoReviewAttrPost1[keys].childAttrArr[key1][key2].paramName]
                }
              } else {
                delete infoReviewFormPost[infoReviewAttrPost1[keys].paramName + "ChildObj"][key1][infoReviewAttrPost1[keys].childAttrArr[key1][key2].paramName]
              }
              if (infoReviewAttrPost1[keys].childAttrArr[key1][key2].childAttrArrThird) {
                let thirdName = infoReviewAttrPost1[keys].childAttrArr[key1][key2].paramName
                let arrtLine = infoReviewAttrPost1[keys].childAttrArr[key1][key2].childAttrArrThird
                for (let key11 in arrtLine) {
                  var ArraTTr1 = []
                  for (let key22 in arrtLine[key11]) {
                    if (arrtLine[key11][key22].thisShowIf == "1") {
                      if (!((!arrtLine[key11][key22].paramRequired) && (infoReviewFormPost[`${secondName}ChildObj`][key1][`${thirdName}ChildObj3`][key11][arrtLine[key11][key22].paramName] === "" || infoReviewFormPost[`${secondName}ChildObj`][key1][`${thirdName}ChildObj3`][key11][arrtLine[key11][key22].paramName] == undefined || infoReviewFormPost[`${secondName}ChildObj`][key1][`${thirdName}ChildObj3`][key11][arrtLine[key11][key22].paramName] === null || infoReviewFormPost[`${secondName}ChildObj`][key1][`${thirdName}ChildObj3`][key11][arrtLine[key11][key22].paramName].length == 0 || !infoReviewFormPost[`${secondName}ChildObj`][key1][`${thirdName}ChildObj3`][key11][arrtLine[key11][key22].paramName]))) {
                        ArraTTr1.push(arrtLine[key11][key22])
                      } else {
                        delete infoReviewFormPost[`${secondName}ChildObj`][key1][`${thirdName}ChildObj3`][key11][arrtLine[key11][key22].paramName]
                      }
                    } else {
                      infoReviewFormPost[`${secondName}ChildObj`][key1][`${thirdName}ChildObj3`][key11][arrtLine[key11][key22].paramName]
                    }
                  }
                  infoReviewAttrPost1[keys].childAttrArr[key1][key2].childAttrArrThird[key11] = ArraTTr1
                }
              }
            }
            infoReviewAttrPost1[keys].childAttrArr[key1] = ArraTTr
          }
        }
        if (infoReviewAttrPost1[keys].thisShowIf == '1') {
          if (!(!infoReviewAttrPost1[keys].paramRequired && (infoReviewFormPost[infoReviewAttrPost1[keys].paramName] === '' || infoReviewFormPost[infoReviewAttrPost1[keys].paramName] === 0 || !infoReviewFormPost[infoReviewAttrPost1[keys].paramName] || infoReviewFormPost[infoReviewAttrPost1[keys].paramName].length == 0))) {
            infoReviewAttrPost.push(infoReviewAttrPost1[keys])
          } else {
            if (!infoReviewAttrPost1[keys].paramRequired && infoReviewFormPost[infoReviewAttrPost1[keys].paramName] === 0) {
              infoReviewAttrPost.push(infoReviewAttrPost1[keys])
            } else {
              delete infoReviewFormPost[infoReviewAttrPost1[keys].paramName]
            }
          }
        } else {
          delete infoReviewFormPost[infoReviewAttrPost1[keys].paramName]
        }
      }
      for (let key in infoReviewAttrPost) {
        if ((infoReviewAttrPost[key].colorRed === '2' || infoReviewAttrPost[key].colorRed === '3' || infoReviewAttrPost[key].colorRed === '4' || infoReviewAttrPost[key].colorRed === '5')) {
          if (this.newTaskName != '资源模板新建' && this.newTaskName != '资源模板编辑' && this.newTaskName != '弹性资源模板新建' && this.newTaskName != '弹性资源模板编辑' && this.newTaskName != '弹性资源模板改配') {
            this.$message({
              showClose: true,
              message: `请正确填写表单！(错误提示：${infoReviewAttrPost[key].paramTitle})`,
              type: 'warning'
            });
            return false
          }
        }
        if (infoReviewAttrPost[key].paramType == "object") {
          if (infoReviewAttrPost[key].paramList) {
            let arr = []
            if (!infoReviewFormPost[infoReviewAttrPost[key].paramName + "ChildObj"]) {
              infoReviewFormPost[infoReviewAttrPost[key].paramName].map(item => {
                if (Object.prototype.toString.call(item).indexOf('Object')==-1) {
                  arr.push({ 'id': item })
                } else {
                  if (item) {
                    if (this.newTaskName != '资源模板新建' && this.newTaskName != '资源模板编辑' && this.newTaskName != '弹性资源模板新建' && this.newTaskName != '弹性资源模板编辑' && this.newTaskName != '弹性资源模板改配') {
                      if (item.oid) {
                        arr.push({ 'id': item.oid })
                      } else {
                        arr.push({ 'id': item.id })
                      }
                    } else {
                      if (item.oid) {
                        item['id'] = item.oid
                      }
                      arr.push(item)
                    }
                  }
                }
              })
              inputData[infoReviewAttrPost[key].paramName] = arr
            }
          } else {
            inputData[infoReviewAttrPost[key].paramName] = {
              'id': infoReviewFormPost[infoReviewAttrPost[key].paramName]
            }
          }

        } else if (infoReviewAttrPost[key].paramType) {
          if (infoReviewAttrPost[key].paramStyle == 'upload') {
            if (infoReviewAttrPost[key].paramList) {
              inputData[infoReviewAttrPost[key]['paramName']] = []
              infoReviewFormPost[infoReviewAttrPost[key]['paramName']].map((k22, i22) => {
                inputData[infoReviewAttrPost[key]['paramName']].push(k22.url)
              })
            } else {
              inputData[infoReviewAttrPost[key]['paramName']] = infoReviewFormPost[infoReviewAttrPost[key]['paramName']][0].url
            }
          } else {
            inputData[infoReviewAttrPost[key]['paramName']] = infoReviewFormPost[infoReviewAttrPost[key]['paramName']]
          }
        }
        if (infoReviewAttrPost[key].childAttrArr) {
          let secondName = infoReviewAttrPost[key].paramName
          for (let inforViewItem in infoReviewAttrPost[key].childAttrArr) {
            if (infoReviewAttrPost[key].childAttrArr[inforViewItem][0].paramFromName) {
              joinTitle = `joint${infoReviewAttrPost[key].childAttrArr[inforViewItem][0].paramFromName}`
            }
            if (!joinTitle) {
              this.gaipeiArr.map(item => {
                if (item.name == this.typeCheckArr[inforViewItem] || item.value == this.typeCheckArr[inforViewItem]) {
                  joinTitle = `joint${item.name}`
                }
              })
            }
            var arrs = []
            inputdataapent = {}
            let inforViewClidItem = infoReviewAttrPost[key].childAttrArr[inforViewItem]
            for (let IIItem in inforViewClidItem) {
              if (inforViewClidItem[IIItem].colorRed === '2' || inforViewClidItem[IIItem].colorRed === '3' || inforViewClidItem[IIItem].colorRed === '4' || inforViewClidItem[IIItem].colorRed === '5') {
                if (this.newTaskName != '资源模板新建' && this.newTaskName != '资源模板编辑' && this.newTaskName != '弹性资源模板新建' && this.newTaskName != '弹性资源模板编辑' && this.newTaskName != '弹性资源模板改配') {
                  this.$message({
                    showClose: true,
                    message: `请正确填写表单！(错误提示：${inforViewClidItem[IIItem].paramTitle})`,
                    type: 'warning'
                  });
                  return false
                }
              }
              if (inforViewClidItem[IIItem].paramType == "object") {
                if (inforViewClidItem[IIItem].paramList) {
                  arrs = []
                  if (Object.prototype.toString.call(infoReviewFormPost[`${secondName}ChildObj`][inforViewItem][inforViewClidItem[IIItem].paramName]).indexOf("String") != -1) {
                    arrs.push(infoReviewFormPost[`${secondName}ChildObj`][inforViewItem][inforViewClidItem[IIItem].paramName])
                  } else {
                    infoReviewFormPost[`${secondName}ChildObj`][inforViewItem][inforViewClidItem[IIItem].paramName].map(item => {
                      if (Object.prototype.toString.call(item).indexOf("Object") != -1) {
                        if (this.newTaskName != '资源模板新建' && this.newTaskName != '资源模板编辑' && this.newTaskName != '弹性资源模板新建' && this.newTaskName != '弹性资源模板编辑' && this.newTaskName != '弹性资源模板改配') {
                          if (item.oid) {
                            arrs.push({ id: item.oid })
                          } else {
                            arrs.push({ id: item.id })
                          }
                        } else {
                          if (item.oid) {
                            item['id'] = item.oid
                          }
                          arrs.push(item)
                        }
                      } else {
                        arrs.push({ id: item })
                      }
                    })
                  }
                  inputdataapent[inforViewClidItem[IIItem].paramName] = arrs
                } else {
                  inputdataapent[inforViewClidItem[IIItem].paramName] = {
                    'id': infoReviewFormPost[`${secondName}ChildObj`][inforViewItem][inforViewClidItem[IIItem].paramName]
                  }
                }
              } else if (inforViewClidItem[IIItem].paramType) {
                if (inforViewClidItem[IIItem].paramStyle == 'upload') {
                  if (inforViewClidItem[IIItem].paramList) {
                    inputdataapent[inforViewClidItem[IIItem].paramName] = []
                    infoReviewFormPost[`${secondName}ChildObj`][inforViewItem][inforViewClidItem[IIItem]['paramName']].map((k22, i22) => {
                      inputData[infoReviewAttrPost[key]['paramName']].push(k22.url)
                    })
                  } else {
                    inputdataapent[inforViewClidItem[IIItem].paramName] = infoReviewFormPost[`${secondName}ChildObj`][inforViewItem][inforViewClidItem[IIItem]['paramName']][0].url
                  }
                } else {
                  inputdataapent[inforViewClidItem[IIItem].paramName] = infoReviewFormPost[`${secondName}ChildObj`][inforViewItem][inforViewClidItem[IIItem]['paramName']]
                }
              }
              if (inforViewClidItem[IIItem].childAttrArrThird) {
                let thirdName = inforViewClidItem[IIItem].paramName
                inputdataapent[`${thirdName}`] = []
                for (let inforViewItem1 in inforViewClidItem[IIItem].childAttrArrThird) {
                  joinTitle1 = `joint${infoReviewAttrPost[key].childAttrArr[inforViewItem][IIItem].childAttrArrThird[inforViewItem1][0].paramFromName}`
                  var arrs1 = []
                  inputdataapent1 = {}
                  let inforViewClidItem1 = infoReviewAttrPost[key].childAttrArr[inforViewItem][IIItem].childAttrArrThird[inforViewItem1]
                  for (let IIItem1 in inforViewClidItem1) {
                    if (inforViewClidItem1[IIItem1].colorRed === '2' || inforViewClidItem1[IIItem1].colorRed === '3' || inforViewClidItem1[IIItem1].colorRed === '4' || inforViewClidItem1[IIItem1].colorRed === '5') {
                      if (this.newTaskName != '资源模板新建' && this.newTaskName != '资源模板编辑' && this.newTaskName != '弹性资源模板新建' && this.newTaskName != '弹性资源模板编辑' && this.newTaskName != '弹性资源模板改配') {
                        this.$message({
                          showClose: true,
                          message: `请正确填写表单！(错误提示：${inforViewClidItem1[IIItem1].paramTitle})`,
                          type: 'warning'
                        });
                        return false
                      }
                    }
                    if (inforViewClidItem1[IIItem1].paramType == "object") {
                      if (inforViewClidItem1[IIItem1].paramList) {
                        arrs1 = []
                        if (infoReviewFormPost[`${secondName}ChildObj`][inforViewItem][`${thirdName}ChildObj3`][inforViewItem1][inforViewClidItem1[IIItem1].paramName]) {
                          infoReviewFormPost[`${secondName}ChildObj`][inforViewItem][`${thirdName}ChildObj3`][inforViewItem1][inforViewClidItem1[IIItem1].paramName].map(item => {
                            if (Object.prototype.toString.call(item).indexOf("Object") != -1) {
                              if (this.newTaskName != '资源模板新建' && this.newTaskName != '资源模板编辑' && this.newTaskName != '弹性资源模板新建' && this.newTaskName != '弹性资源模板编辑' && this.newTaskName != '弹性资源模板改配') {
                                if (item.oid) {
                                  arrs1.push({ id: item.oid })
                                } else {
                                  arrs1.push({ id: item.id })
                                }
                              } else {
                                if (item.oid) {
                                  item['id'] = item.oid
                                }
                                arrs1.push(item)
                              }
                            } else {
                              arrs1.push({ id: item })
                            }
                          })
                        }
                        inputdataapent1[inforViewClidItem1[IIItem1].paramName] = arrs1
                      } else {
                        inputdataapent1[inforViewClidItem1[IIItem1].paramName] = {
                          'id': infoReviewFormPost[`${secondName}ChildObj`][inforViewItem][`${thirdName}ChildObj3`][inforViewItem1][inforViewClidItem1[IIItem1].paramName]
                        }
                      }
                    } else if (inforViewClidItem1[IIItem1].paramType) {
                      inputdataapent1[inforViewClidItem1[IIItem1].paramName] = infoReviewFormPost[`${secondName}ChildObj`][inforViewItem][`${thirdName}ChildObj3`][inforViewItem1][inforViewClidItem1[IIItem1].paramName]
                    }
                  }
                  if (!joinTitle1) {
                    joinTitle1 = `joint${this.childNodename3}`
                  }
                  if (this.historyTableList._copyOrupdate) {
                    inputdataapent1["formName"] = ''
                  } else {
                    inputdataapent1["formName"] = infoReviewFormPost[`${secondName}ChildObj`][inforViewItem][`${thirdName}ChildObj3`][inforViewItem1]["formName"] ? infoReviewFormPost[`${secondName}ChildObj`][inforViewItem][`${thirdName}ChildObj3`][inforViewItem1]["formName"] : ''
                  }
                  inputdataapent1["formTitle"] = `批次${Number(inforViewItem) + 1}.${Number(inforViewItem1) + 1} - ${inforViewClidItem[IIItem].childAttrArrThird[inforViewItem1][0].paramFromTitle}`
                  let newArrJoin1 = {
                    "formName": infoReviewFormPost[`${secondName}ChildObj`][inforViewItem][`${thirdName}ChildObj3`][inforViewItem1]["formName1"] ? infoReviewFormPost[`${secondName}ChildObj`][inforViewItem][`${thirdName}ChildObj3`][inforViewItem1]["formName1"] : '',
                    [joinTitle1]: inputdataapent1
                  }
                  if (this.historyTableList._copyOrupdate) {
                    newArrJoin1["formName"] = ''
                  }
                  inputdataapent[`${thirdName}`].push(newArrJoin1)
                }
              }
            }
            paramNames = infoReviewAttrPost[key].paramName
            if (!joinTitle) {
              joinTitle = `joint${this.childNodename}`
            }
            if (this.historyTableList._copyOrupdate) {
              inputdataapent["formName"] = ''
            } else {
              inputdataapent["formName"] = infoReviewFormPost[paramNames + "ChildObj"][inforViewItem]["formName"] ? infoReviewFormPost[paramNames + "ChildObj"][inforViewItem]["formName"] : ''
            }
            inputdataapent["formTitle"] = `批次${Number(inforViewItem) + 1} - ${infoReviewAttrPost[key].childAttrArr[inforViewItem][0].paramFromTitle}`
            let newArrJoin = {
              "formName": infoReviewFormPost[paramNames + "ChildObj"][inforViewItem]["formName1"] ? infoReviewFormPost[paramNames + "ChildObj"][inforViewItem]["formName1"] : '',
              [joinTitle]: inputdataapent
            }
            if (this.historyTableList._copyOrupdate) {
              newArrJoin["formName"] = ''
            }
            newParArr.push(newArrJoin)
          }
        }
      }
      if (type) {
        this.rotate()
        return
      }
      if (newParArr.length != 0) {
        inputData[paramNames] = newParArr
      }
      if (this.historyTableList._copyOrupdate) {
        inputData["formName"] = ''
      } else {
        inputData["formName"] = infoReviewFormPost["formName"] ? infoReviewFormPost["formName"] : ''
      }
      let submitDataInput = {
        input: [
          inputData
        ]
      }
      let submitData = {}
      this.potsCreateObj = {}
      if (this.newTaskName == '资源模板新建' || this.newTaskName == '资源模板编辑' || this.newTaskName == '弹性资源模板新建' || this.newTaskName == '弹性资源模板编辑' || this.newTaskName == '弹性资源模板改配') {
        this.infoReviewAttr.map((k, i) => {
          if (k.childAttrArr) {
            k.childAttrArr.map((k1, i1) => {
              k1.map((k2, i2) => {
                if (k2.paramType == 'object') {
                  let nameStr = ""
                  for (var v in inputData[k.paramName][i1]) {
                    if (v.indexOf("joint") == 0) {
                      nameStr = v
                    }
                  }
                  if (k2.childAttrArrThird) {
                    k2.childAttrArrThird.map((k3, i3) => {
                      k3.map((k4, i4) => {
                        if (k4.paramType == 'object' && k2.paramreference == 'FormTemplate') {
                          let nameStr12 = ""
                          for (var v12 in inputData[k.paramName][i1][nameStr][k2.paramName][i3]) {
                            if (v12.indexOf("joint") == 0) {
                              nameStr12 = v12
                            }
                          }
                          if (inputData[k.paramName][i1][nameStr][k2.paramName][i3][nameStr12] && inputData[k.paramName][i1][nameStr][k2.paramName][i3][nameStr12]) {
                            console.log(Object.prototype.toString.call(inputData[k.paramName][i1][nameStr][k2.paramName][i3][nameStr12][k4.paramName]).indexOf("Object") != -1)
                            k4.paramDataList.map((k8, i8) => {
                              if (inputData[k.paramName][i1][nameStr][k2.paramName][i3][nameStr12][k4.paramName]) {
                                if (Object.prototype.toString.call(inputData[k.paramName][i1][nameStr][k2.paramName][i3][nameStr12][k4.paramName]).indexOf("Object") != -1) {
                                  if (k8.value == inputData[k.paramName][i1][nameStr][k2.paramName][i3][nameStr12][k4.paramName].id) {
                                    inputData[k.paramName][i1][nameStr][k2.paramName][i3][i4][nameStr12][k4.paramName] = {
                                      id: k8.value,
                                      oid: k8.value,
                                      name: k8.name,
                                    }
                                  }
                                } else {
                                  inputData[k.paramName][i1][nameStr][k2.paramName][i3][nameStr12][k4.paramName].map((k5, i4) => {
                                    if (k8.value == k5.id) {
                                      k5["oid"] = k8.value
                                      k5["name"] = k8.name
                                    }
                                  })
                                }
                              }
                            })
                          }
                        }
                      })
                    })
                  } else
                    if (inputData[k.paramName][i1][nameStr][k2.paramName]) {
                      k2.paramDataList.map((k3, i3) => {
                        if (Object.prototype.toString.call(inputData[k.paramName][i1][nameStr][k2.paramName]).indexOf("Object") != -1) {
                          if (k3.value == inputData[k.paramName][i1][nameStr][k2.paramName].id) {
                            inputData[k.paramName][i1][nameStr][k2.paramName] = {
                              id: k3.value,
                              oid: k3.value,
                              name: k3.name,
                            }
                          }
                        } else {
                          inputData[k.paramName][i1][nameStr][k2.paramName].map((k4, i4) => {
                            if (k3.value == k4.id) {
                              k4["oid"] = k3.value
                              k4["name"] = k3.name
                            }
                          })
                        }
                      })
                    }
                }
              })
            })
          } else {
            if (k.paramType == 'object') {
              k.paramDataList.map((k1, i1) => {
                if (Object.prototype.toString.call(inputData[k.paramName]).indexOf("Object") != -1) {
                  if (k1.value == inputData[k.paramName].id) {
                    inputData[k.paramName] = {
                      id: k1.value,
                      oid: k1.value,
                      name: k1.name,
                    }
                    if (k1.oldValue) {
                      for (let k44 in k1.oldValue) {
                        if (!inputData[k.paramName][k44]) {
                          inputData[k.paramName][k44] = k1.oldValue[k44]
                        }
                      }
                    }
                  }
                } else {
                  if (inputData[k.paramName]) {
                    inputData[k.paramName].map((k2, i2) => {
                      if (k2.id == k1.value) {
                        k2["oid"] = k1.value
                        k2["name"] = k1.name
                      }
                    })
                  }
                }
              })
            }
          }
        })
        if (this.newTaskName == '资源模板新建' || this.newTaskName == '弹性资源模板新建' || this.newTaskName == '弹性资源模板改配') {
          var mbObj = {
            "input": [
              {
                "cloudProjectConfig": {
                  "where": {
                    "id": inputData.project.id
                  }
                },

                "name": inputData.formTitle,
                "description": "",
                "template": {
                  "where": {
                    "name": this.nodeCheckName
                  }
                },
                "fill": JSON.stringify(inputData)
              }
            ]
          }
          if (this.newTaskName == '弹性资源模板新建' || this.newTaskName == '弹性资源模板改配') {
            mbObj.input[0]['cloudQuotaConfig'] = {
              "where": {
                "id": inputData.cloudQuotaConfig ? inputData.cloudQuotaConfig.id : ''
              }
            }
          }
          submitData = { "query": `mutation ($input: [JobTemplateFillerCreate!]!){createJobTemplateFillers(input: $input) {nodeCreated}}`, "variables": mbObj }
        } else {
          var mbObj = {
            "where": {
              "id": this.historyTableList._id
            },
            "update": {
              "name": inputData.formTitle,
              "description": "",
              "cloudProjectConfig": {
                "where": {
                  "id": inputData.project.id
                }
              },
              "fill": JSON.stringify(inputData)
            }
          }
          if (this.newTaskName == '弹性资源模板编辑') {
            mbObj.update['cloudQuotaConfig'] = {
              "where": {
                "id": inputData.cloudQuotaConfig ? inputData.cloudQuotaConfig.id : ''
              }
            }
          }
          submitData = { "query": `mutation($where:JobTemplateFillerWhere,$update:JobTemplateFillerUpdate){updateJobTemplateFillers(where:$where, update:$update){nodeUpdated}}`, "variables": mbObj }
        }
        this.potsCreateObj = submitData
        if (this.newTaskName == '资源模板新建' || this.newTaskName == '弹性资源模板新建' || this.newTaskName == '弹性资源模板改配') {
          if (this.historyTableList._projectId) {
            this.formPost = {
              name: this.historyTableList._name,
              project: inputData.project.id,
              cloudQuotaConfig: inputData.cloudQuotaConfig ? inputData.cloudQuotaConfig.id : '',
              remark: this.historyTableList._description ? this.historyTableList._description : ''
            }
            if (this.projecArr.length == 0) {
              this.projecArr = [{
                name: inputData.project.name,
                value: inputData.project.id
              }]
            }
            this.cloudQuotaConfigArr = [{
              name: inputData.cloudQuotaConfig ? inputData.cloudQuotaConfig.name : '',
              value: inputData.cloudQuotaConfig ? inputData.cloudQuotaConfig.id : ''
            }]
          } else {
            this.formPost = {
              name: inputData.formTitle,
              project: inputData.project.id,
              cloudQuotaConfig: inputData.cloudQuotaConfig ? inputData.cloudQuotaConfig.id : '',
              remark: ''
            }
            if (this.projecArr.length == 0) {
              this.projecArr = [{
                name: inputData.project.name,
                value: inputData.project.id
              }]
            }
            this.cloudQuotaConfigArr = [{
              name: inputData.cloudQuotaConfig ? inputData.cloudQuotaConfig.name : '',
              value: inputData.cloudQuotaConfig ? inputData.cloudQuotaConfig.id : ''
            }]
          }
        } else if (this.newTaskName == '资源模板编辑' || this.newTaskName == '弹性资源模板编辑') {
          this.formPost = {
            name: this.historyTableList._name,
            project: inputData.project.id,
            cloudQuotaConfig: inputData.cloudQuotaConfig ? inputData.cloudQuotaConfig.id : '',
            remark: this.historyTableList._description ? this.historyTableList._description : ''
          }
          if (this.projecArr.length == 0) {
            this.projecArr = [{
              name: inputData.project.name,
              value: inputData.project.id
            }]
          }
          this.cloudQuotaConfigArr = [{
            name: inputData.cloudQuotaConfig ? inputData.cloudQuotaConfig.name : '',
            value: inputData.cloudQuotaConfig ? inputData.cloudQuotaConfig.id : ''
          }]
        }
        if (!this.formPost.project) {
          this.$message({
            showClose: true,
            message: "模板新建项目不能为空！",
            type: 'warning'
          });
          return false
        }
        this.batchCearteStatus = true
        return false
      } else {
        this.lowerName = this.nodeCheckName.toLowerCase()
        submitData = { "query": `mutation($input: [${this.nodeCheckName}Create!]!){upsert${this.nodeCheckName}s(input:$input){nodeCreated ${this.lowerName}s{id job{id}}}}`, "variables": submitDataInput }
      }
      let createText = '确定提交吗？'
      if (this.nodeCheckName == 'FormCreateCloudProduct') {
        createText = `请业务在资源创建后自行添加监控报警，默认未配置报警；具体配置方法，请咨询`
      }
      if (this.nodeCheckName == 'FormDeleteCloudProduct') {
        createText = '集群内的云服务器请联系运维同学从集群移除后再发起清退申请，直接清退会再行自动创建同规格云服务器'
      }
      this.reviewload = true
      this.$confirm(`${createText}`, {
        confirmButtonText: "确定",
        cancelButtonText: "取消",
        type: "warning",
      }).then(() => {
        var postData = {
          form:inputData,
          "template_id":this.template_id,
        }
        Http.createJob(postData).then((response) => {
            this.$message({
              showClose: true,
              message:"提交成功",
              type: 'success',
              duration: 10000,  
            });
            this.typeCheckArr = []
            this.theadArrShow = []
            this.reviewStatusNew = false
            this.batchCearteStatus = false
            this.reviewload = false
          }).catch(() => {
            this.reviewload = false
          })
      }).catch(() => {
        this.reviewload = false
      })
    },
    // 流转按钮
    rotate() {
      if (this.historyTableList._cometProcess) {
        this.$confirm(`是否确定流转`, {
          confirmButtonText: "确定",
          cancelButtonText: "取消",
          type: "warning",
        }).then(() => {
          console.log(123)
        })
          .catch(() => {
            this.reviewload = false
          })
      }
    },
    // 获取表单信息
    getNodeInfo(type) {
      var postData = {
        schema: "form_template",
        "where":{
          "name":this.nodeCheckName
        }
      }
      Http.getFormTemplate(postData).then(
        (response) => {
          if(response.data.data.data.length){
            var formDataObj = JSON.parse(response.data.data.data[0].data)
            this.template_id = response.data.data.data[0].id
            this.infoReviewAttr = [];
            this.infoReviewForm = {};
            if (formDataObj) {
              this.FormTemplatjobMode = formDataObj.jobMode ? formDataObj.jobMode : "";
              if (!formDataObj.parameters) {
                return false;
              }
              formDataObj.parameters.map((k, i) => {
                var tempArr = [];
                if (k.templates) {
                  k.templates.map((k, i) => {
                    tempArr.push(k.name);
                  });
                }
                var enumQueryArr = [];
                if (k.type != "string") {
                  enumQueryArr = [];
                }
                if (k.enum) {
                  enumQueryArr = [];
                  k.enum.map((k, i) => {
                    enumQueryArr.push({ name: k.zh, value: k.en });
                  });
                }
                this.infoReviewAttr.push({
                  paramType: k.type,
                  paramName: k.name,
                  paramTitle: k.title,
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
                  paramreferenceMutation: k.referenceMutation ? k.referenceMutation : "",
                  paramreferenceSelect: k.referenceSelect ? k.referenceSelect : "",
                  parmreformat: k.format ? k.format : "",
                  paramavailableCondition: k.availableCondition ? k.availableCondition : "",
                  thisShowIf: "1",
                  colorRed: "1",
                  customDescription: "",
                  paramGroup: k.group,
                  paramMutipleof: k.mutipleof,
                  _id: i + 1,
                  ID: k.id,
                  paramTableTitle: [],
                  paramreferenceDisplay: k.referenceDisplay ? k.referenceDisplay : "",
                  paramannotation: k.annotation,
                  paramTemplates: k.templates ? k.templates : "",
                });
                if (k.type == "object") {
                  if (k.list) {
                    this.infoReviewForm[k.name] = [];
                  } else {
                    this.infoReviewForm[k.name] = "";
                  }
                } else {
                  if (k.list) {
                    this.infoReviewForm[k.name] = [];
                  } else {
                    this.infoReviewForm[k.name] = "";
                  }
                }
              });
              this.infoReviewForm = JSON.parse(JSON.stringify(this.infoReviewForm));
              this.infoReviewAttr = JSON.parse(JSON.stringify(this.infoReviewAttr));
              this.reviewParam(type)

              if (this.checkListLast.length) {
                this.infoReviewAttr.map((iAItem, Iaii) => {
                   //首先拿到所有的Attr值判断是否有子表单的元素
                  if (iAItem.paramType == 'object' && iAItem.paramreference == 'FormTemplate') {
                    this.showChildAttr(iAItem, iAItem.paramName, this.childNodename, Iaii, 0, {})
                  }
                })
              }        

            }
          } 
      })
    }
  }
};

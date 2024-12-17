export const advSearchHandle = {
  methods: {
    advSearchHandle(advItem, name, arrTitle) {
      var arrTitleShow = arrTitle || this.showlistObj.properties;
      var objQueryData = {};
      let arrName = [];
      let arrNameNot = [];
      if (advItem.type.indexOf("selectNull") != -1) {
        let strname = advItem.name;
        if (advItem.val == true || advItem.val == "true") {
          let quname = "";
          let qunameArrNewO = {
            OR: [],
          };
          if (strname.indexOf(".") != -1) {
            var searchKey = strname.split(".");
            for (var lenk = 0; lenk < searchKey.length; lenk++) {
              var qunameF = "";
              var arrCopyS = searchKey.slice(0, searchKey.length - lenk).reverse();
              arrCopyS.map((nitem, ni) => {
                if (ni == 0) {
                  qunameF = `{"NOT":{"HAS":${JSON.stringify([nitem])}}}`;
                } else {
                  qunameF = `{"${nitem}":${qunameF}}`;
                }
              });
              let objJsonStr = JSON.parse(qunameF);
              qunameArrNewO["OR"].push(objJsonStr);
            }
            objQueryData = qunameArrNewO;

          } else {
            quname = strname;
            arrNameNot.push(quname);
          }
        } else if (advItem.val == "no") {
          let quname = "";
          if (strname.indexOf(".") != -1) {
            var searchKey = strname.split(".");
            var searKey = searchKey.reverse();
            searKey.map((nitem, ni) => {
              if (ni == 0) {
                quname = `{"HAS":${JSON.stringify([nitem])}}`;
              } else {
                quname = `{"${nitem}":${quname}}`;
              }
            });
            let objJsonStr = JSON.parse(quname);
            objQueryData = objJsonStr;
          } else {
            quname = `${strname}`;
            arrName.push(quname);
          }
        }
      } else {
        arrTitleShow.map((ppitem, ppinx) => {
          if (ppitem.name == advItem.name) {
            if (ppitem.where) {
              let whereQueryStr = JSON.parse(JSON.stringify(ppitem.where));
              if (whereQueryStr.indexOf("$typecond") != -1) {
                let objstr = {};
                advItem.map((olkitem) => {
                  objstr[olkitem] = {};
                });
                let restr = ppitem.where.replace(
                  "{$typecond}",
                  JSON.stringify(objstr)
                );
                objQueryData = restr;
              } else if (whereQueryStr.indexOf("$func") != -1) {
                let objJsonStr = "";
                if (advItem.type == "arr") {
                  var batcharr = [];
                  var batchNewarr = [];
                  batcharr = advItem.val.split("\n");
                  batcharr.map((k, i) => {
                    if (k.replace(/\s/g, "")) {
                      batchNewarr.push(k.replace(/\s/g, ""));
                    }
                  });
                  this.batcharrLength = batchNewarr;
                  if (advItem.name.indexOf(".") != -1) {
                    var searchKey = advItem.name.split(".").reverse();
                    var strTemp = "";
                    searchKey.map((k1, i1) => {
                      let showKey = k1;
                      if (i1 == 0) {
                        strTemp = `{"${showKey}_IN":${JSON.stringify(batchNewarr)}}`;
                      }
                    });
                    objJsonStr = JSON.parse(strTemp);
                  } else {
                    objJsonStr = `{${k + "_IN"}:${batchNewarr}}`;
                  }
                } else {
                  if (advItem.name.indexOf(".") != -1) {
                    var searKey = advItem.name.split(".").reverse();
                    var showJson1 = "";
                    searKey.map((k1, i1) => {
                      let showKey = "";
                      if (advItem.type == "_HGT") {
                        showKey = k1 + "_GT";
                      }
                      if (advItem.type == "_HLT") {
                        showKey = k1 + "_LT";
                      }
                      if (advItem.type == "_DGT") {
                        showKey = k1 + "_GT";
                      }
                      if (advItem.type == "_DLT") {
                        showKey = k1 + "_LT";
                      }
                      if (advItem.type == "select") {
                        showKey = k1 + "_IN";
                      }
                      if (advItem.type != "_HLT" && advItem.type != "_HGT" && advItem.type != "_DLT" && advItem.type != "_DGT" && advItem.type != "select") {
                        showKey = k1 + advItem.type;
                      }
                      if (i1 == 0) {
                        if (Object.prototype.toString.call(advItem.val).indexOf("Array") != -1) {
                          showJson1 = `{"${showKey}":${JSON.stringify(advItem.val)}}`;
                        } else {
                          showJson1 = `{"${showKey}":"${advItem.val}"}`;
                        }
                        showJson1 = `{"${showKey}":"${advItem.val}"}`;
                      } 
                    });
                    objJsonStr = JSON.parse(showJson1);
                  } else {
                    let showKey1 = "";
                    var showJson1 = "";
                    if (advItem.type == "_HGT") {
                      showKey1 = advItem.name + "_GT";
                    }
                    if (advItem.type == "_HLT") {
                      showKey1 = advItem.name + "_LT";
                    }
                    if (advItem.type == "_DGT") {
                      showKey1 = advItem.name + "_GT";
                    }
                    if (advItem.type == "_DLT") {
                      showKey1 = advItem.name + "_LT";
                    }
                    if (advItem.type == "select") {
                      showKey1 = advItem.name + "_IN";
                    }
                    if (advItem.type != "_HLT" && advItem.type != "_HGT" && advItem.type != "_DLT" && advItem.type != "_DGT" && advItem.type != "select") {
                      showKey1 = advItem.name + advItem.type;
                    }
                    showJson1 = `{"${showKey1}":"${advItem.val}"}`;
                    objJsonStr = JSON.parse(showJson1);
                  }
                }
                let reshowJson1 = ppitem.where.replace(
                  '{"name$func":$cond}',
                  JSON.stringify(objJsonStr)
                );
                objQueryData = reshowJson1;
              }
            } else {
              if (advItem.type == "arr") {
                var batcharr = [];
                var batchNewarr = [];
                batcharr = advItem.val.split("\n");
                batcharr.map((k, i) => {
                  if (k.replace(/\s/g, "")) {
                    batchNewarr.push(k.replace(/\s/g, ""));
                  }
                });
                this.batcharrLength = batchNewarr;
                if (advItem.name.indexOf(".") != -1) {
                  var searchKey = advItem.name.split(".").reverse();
                  var strTemp = "";
                  searchKey.map((k1, i1) => {
                    let showKey = k1;
                    if (i1 == 0) {
                      strTemp = `{"${showKey}_IN":${JSON.stringify(batchNewarr)}}`;
                    } else {
                      strTemp = `{"${showKey}":${strTemp}}`;
                    }
                  });
                  var obj1 = JSON.parse(strTemp);
                  objQueryData = obj1;
                } else {
                  objQueryData = {
                    [advItem.name + "_IN"]: batchNewarr,
                  };
                }
              } else if (
                advItem.objContent.tagsArr &&
                advItem.objContent.tagsArr.length
              ) {
                let tgbArrs = advItem.objContent.tagsArr;
                // tag下拉搜索
                let tagValueArr = [];
                let tagArr = [];
                advItem.val.map((akey) => {
                  tagValueArr = akey.split(":");
                  let tagstr = "";
                  tgbArrs.map((tItem, tItemii) => {
                    tagstr += `${tItemii == 0 ? "" : ","}"${tItem}":"${tagValueArr[tItemii]
                      }"`;
                  });
                  let showKey = `{${tagstr}}`;
                  tagArr.push(JSON.parse(showKey));
                });
                objQueryData = {
                  [advItem.name]: {
                    OR: tagArr,
                  },
                };
              } else {
                if (advItem.name.indexOf(".") != -1) {
                  var searKey = advItem.name.split(".").reverse();
                  var showJson1 = "";
                  searKey.map((k1, i1) => {
                    let showKey = "";
                    if (advItem.type == "_HGT") {
                      showKey = k1 + "_GT";
                    }
                    if (advItem.type == "_HLT") {
                      showKey = k1 + "_LT";
                    }
                    if (advItem.type == "_DGT") {
                      showKey = k1 + "_GT";
                    }
                    if (advItem.type == "_DLT") {
                      showKey = k1 + "_LT";
                    }
                    if (advItem.type == "select") {
                      showKey = k1 + "_IN";
                    }
                    if (advItem.type != "_HLT" && advItem.type != "_HGT" && advItem.type != "_DLT" && advItem.type != "_DGT" && advItem.type != "select") {
                      showKey = k1 + advItem.type;
                    }
                    if (i1 == 0) {
                      if (Object.prototype.toString.call(advItem.val).indexOf("Array") != -1) {
                        showJson1 = `{"${showKey}":${JSON.stringify(advItem.val)}}`;
                      } else {
                        showJson1 = `{"${showKey}":"${advItem.val}"}`;
                      }
                    } else {
                      showJson1 = `{"${k1}":${showJson1}}`;
                    }
                  });
                  var obj = null;
                  obj = JSON.parse(showJson1);
                  objQueryData = obj;
                } else {
                  let showKey1 = "";
                  if (advItem.type == "_HGT") {
                    showKey1 = advItem.name + "_GT";
                  }
                  if (advItem.type == "_HLT") {
                    showKey1 = advItem.name + "_LT";
                  }
                  if (advItem.type == "_DGT") {
                    showKey1 = advItem.name + "_GT";
                  }
                  if (advItem.type == "_DLT") {
                    showKey1 = advItem.name + "_LT";
                  }
                  if (advItem.type == "select") {
                    showKey1 = advItem.name + "_IN";
                  }
                  if (advItem.type != "_HLT" && advItem.type != "_HGT" && advItem.type != "_DLT" && advItem.type != "_DGT" && advItem.type != "select") {
                    showKey1 = advItem.name + advItem.type;
                  }
                  if (name == "CloudServerType" && advItem.val == "unknown" && advItem.name == "category") {
                    objQueryData = {
                      NOT: {
                        category_IN: ["通用型", "计算型", "内存型"],
                      },
                    };
                  }
                  objQueryData = {
                    [showKey1]: advItem.val,
                  };
                }
              }
            }
          }
        });
      }
      if (arrName.length) {
        objQueryData = {
          HAS: arrName,
        };
      }
      if (arrNameNot.length) {
        objQueryData = {
          NOT: {
            HAS: arrNameNot,
          },
        };
      }
      return objQueryData;
    },
  },
};

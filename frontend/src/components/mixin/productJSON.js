export const productJSON = {
  methods: {
    getTreeData() {

      this.treedata = {
        schema: "",
        projectNoget: true,
        properties: [],
      };  
      if(this.JSONConfig[this.$route.name]){
        this.treedata = this.JSONConfig[this.$route.name]
      }
      this.searchMoreArr = [];
      this.checkboxTH = [];
      this.nodeTitle = this.treedata.schema;
      this.projectNoget = this.treedata.projectNoget ? true : false;
      this.getTable(this.treedata.schema);
      this.tableHeaderArr = this.treedata.properties;
      this.tableHeaderArr.map((k, i) => {
        if (k.isShowDetail) {
          k.hide = true;
          k.hideIn = true;
        }
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
            let arr = [];

            if (item.index) {
              if (item.DataList || item.TagsList) {
                arr.push({
                  value: "select",
                  label: "选择",
                });
              }
              if (item.index.indexOf("hash") != -1) {
                arr.push(
                  {
                    value: "",
                    label: "精确",
                  },
                  {
                    value: "arr",
                    label: "批量",
                  }
                );
              }
              if (item.index.indexOf("regexp") != -1) {
                arr.push({
                  value: "_REGEX",
                  label: "模糊",
                });
              }
              if (item.index.indexOf("hour") != -1) {
                arr.push(
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
                arr.push(
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
              item.children = arr;
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
      let checkboxName = `${this.$route.name}checkboxTHBFArr`;
      let pArr = localStorage.getItem(checkboxName);
      if (pArr) {
        if (pArr.indexOf(",") == -1) {
          this.checkboxTH = [pArr];
        } else {
          this.checkboxTH = pArr.split(",");
        }
      } else {
        this.checkboxTH = this.checkboxTH;
      }
      this.checkboxTHBF = JSON.parse(JSON.stringify(this.checkboxTH));
      this.searchObj = JSON.parse(JSON.stringify(this.searchObj));
      this.showlistObj = this.treedata;
      if (this.showlistObj.oneSearch) {
        this.showlistObj.oneSearch.map((k, i) => {
          if (k.index && k.index.indexOf("hour") != -1) {
            this.searchOneObj[k.name] = "";
          } else {
            if (k.schema) {
              this.getSelList(k);
            }
            this.searchOneObj[k.name] = [];
          }
        });
        this.searchOneObj = JSON.parse(JSON.stringify(this.searchOneObj));
      }
    },
  },
};

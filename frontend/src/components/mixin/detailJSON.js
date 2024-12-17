export const detailMixin = {
  mounted() {
    this.dateTimeValue = [
      new Date(
        this.$moment().subtract(6, "hours").format("YYYY/MM/DD HH:mm:ss")
      ).getTime(),
      new Date(this.$moment().format("YYYY/MM/DD HH:mm:ss")).getTime(),
    ];
    this.detail_IDQ = this.$route.query.detail_ID;
    this.$EventBus.$on("isCollapse", (value) => {
      this.isCollapseStatus = value ? true : false;
    });
    if (this.JSONConfig[this.$route.name]) {
      var objDataDetail = this.JSONConfig[this.$route.name];
      var objDataOne = {
        title: "基本信息",
        groups: [
          {
            title: "基本信息",
            properties: [],
          },
        ],
      };
      var objArr = {};
      objDataDetail.properties.map((k, i) => {
        if (!k.baseType || k.baseType == "基本信息") {
          objDataOne.groups[0].properties.push(k);
        } else {
          if (!objArr[k.baseType]) {
            objArr[k.baseType] = [];
          }
          objArr[k.baseType].push(k);
        }
      });
      for (var k in objArr) {
        objDataOne.groups.push({
          title: k,
          properties: objArr[k],
        });
      }
      if (objDataDetail.tabs) {
        objDataDetail.tabs.unshift(objDataOne);
      } else {
        objDataDetail["tabs"] = [objDataOne];
      }
      this.showConfigData = JSON.parse(JSON.stringify(objDataDetail));
    }
    this.showConfigData.tabs[0].groups[0].properties.map((k, i) => {
      if (k.name == "status") {
        this.statusObj = k.enum;
      }
    });
    this.getRowObj();
  },
};

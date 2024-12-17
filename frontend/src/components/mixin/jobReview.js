export const jobReview = {
  data() {
    return {
      historyTableList: {},
      jumpArr: [],
      jumorRetrypArr: [],
      searchJumpArr: [],
      dealerArr: [],
      jobIdStatus: "",
    };
  },
  methods: {
    //点击左边获取任务详情
    getJobList(row, type) {
      this.jobIdStatus = this.organizationArr[row.status];
      let searchFormData = {
        Id: row.id,
      };
      this.historyTableList = {};
      this.jumpArr = [];
      this.jumorRetrypArr = [];
      this.searchJumpArr = [];
      this.Http("root", `/api/v1/mars/job/get`, "get", searchFormData).then((response) => {
        if (response.data.code == 0) {
          this.detaileMessage = response.data.data;
          this.dealerArr = response.data.data.dealer;
        } else {
          this.$message({
            showClose: true,
            message: `错误码：${response.data.code}， ${response.data.message}`,
            type: "error",
          });
        }
        //根据jobid获取template拼接参数
        // this.getTemplate(row.id, type);
      })
        .catch((err) => {
          console.log(err);
        });
    },
    reviewHandle(historyTableList) {
      if (Object.keys(historyTableList).length) {
        this.detailOrLook = true;
        this.reviewStatus = true;
      }
    },
  },
  mounted() { },
};

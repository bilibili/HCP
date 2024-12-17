//tabHead 导出表头数据  tabContent  导出表格内容数据  tabMessage 整体表格数据  str表格名字
export const exportPage = (tabHead, tabContent, tabMessage, str) => {
  import("@/components/utils/exportExcel").then((excel) => {
    const tHeader = tabHead
    const filterVal = tabContent
    const list = tabMessage
    const data = formatJson(filterVal, list);
    excel.export_json_to_excel({
      header: tHeader,
      data,
      filename: str,
      autoWidth: true,
      bookType: "xlsx",
    });
  });
}
function formatJson(filterVal, jsonData) {
  return jsonData.map((v) =>
    filterVal.map((j) => {
      if (j === "timestamp") {
        return parseTime(v[j]);
      } else {
        return v[j];
      }
    })
  );
}
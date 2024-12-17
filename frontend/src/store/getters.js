const getters = {
  sidebar: state => state.app.sidebar,
  device: state => state.app.device,
  token: state => state.user.token,
  avatar: state => state.user.avatar,
  name: state => state.user.name,
  permissionRoutes: state => state.user.routes,
  globalarr: state => state.user.globalarr,
  loginUserName: state => state.user.loginUserName,
  loginUserID: state => state.user.loginUserID,
  loginUserInfo: state => state.user.loginUserInfo,
  formPolicyDomain: state => state.user.formPolicyDomain,
  manageUser: state => state.user.manageUser,
  manageUserGroup: state => state.user.manageUserGroup,
}
export default getters

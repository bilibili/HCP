const state = {
  showSettings: "",
  fixedHeader: "",
  sidebarLogo: ""
}

const mutations = {
  CHANGE_SETTING: (state, { key, value }) => {
    state[key] = value
  }
}

const actions = {
  changeSetting({ commit }, data) {
    commit('CHANGE_SETTING', data)
  }
}

export default {
  namespaced: true,
  state,
  mutations,
  actions
}


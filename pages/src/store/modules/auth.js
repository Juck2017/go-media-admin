const app = {
  state: {
    checkAuth: []
  },
  mutations: {
    SET_CHECK_AUTH: (state, checkAuth) => {
      state.checkAuth = checkAuth
    }
  },
  actions: {
    setCheckAuth({ commit }, checkAuth) {
      commit('SET_CHECK_AUTH', checkAuth)
    }
  }
}

export default app

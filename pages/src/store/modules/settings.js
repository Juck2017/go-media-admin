import variables from '@/styles/element-variables.scss'
import defaultSettings from '@/settings'
import { getDictMap } from '@/api/settings'

const { showSettings, tagsView, fixedHeader, sidebarLogo } = defaultSettings

const state = {
  theme: variables.theme,
  showSettings: showSettings,
  tagsView: tagsView,
  fixedHeader: fixedHeader,
  sidebarLogo: sidebarLogo,
  dicMap: {} // 字典map
}

const mutations = {
  CHANGE_SETTING: (state, { key, value }) => {
    // eslint-disable-next-line no-prototype-builtins
    if (state.hasOwnProperty(key)) {
      state[key] = value
    }
  },
  SET_DICT_MAP: (state, dicMap) => {
    state.dicMap = dicMap
  }
}

const actions = {
  changeSetting({ commit }, data) {
    commit('CHANGE_SETTING', data)
  },
  async setDictMap({ commit }) {
    const response = await getDictMap()
    commit('SET_DICT_MAP', response.result)
  }
}

export default {
  namespaced: true,
  state,
  mutations,
  actions
}


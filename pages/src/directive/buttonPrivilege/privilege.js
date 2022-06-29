import store from '../../store'

export default {
  // bind(el, binding) {
  //   if (!binding.value || binding.value === '') {
  //     return
  //   }
  //   if (store.getters.buttonPrivileges.indexOf(binding.value) === -1) {
  //     el.parentNode.removeChild(el)
  //   }
  // },

  inserted(el, binding) {
    if (!binding.value || binding.value === '') {
      return
    }
    if (store.getters.buttonPrivileges.indexOf(binding.value) === -1) {
      el.parentNode.removeChild(el.nextElementSibling) // 移除下一个元素(竖线)
      el.parentNode.removeChild(el)
    }
  },
  update(el, binding) {
    if (!binding.value || binding.value === '') {
      return
    }
    if (store.getters.buttonPrivileges.indexOf(binding.value) === -1 && el.parentNode) {
      el.parentNode.removeChild(el.nextElementSibling) // 移除下一个元素(竖线)
      el.parentNode.removeChild(el)
    }
  }
}

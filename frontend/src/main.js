import Vue from 'vue'
import 'normalize.css/normalize.css' 
import ElementUI from 'element-ui'
import 'element-ui/lib/theme-chalk/index.css'
import locale from 'element-ui/lib/locale/lang/zh-CN'
import '@/assets/style/index.scss'
import App from './App'
import store from './store'
import router from './router'
import '@/permission' 
import { setVueInstance } from "./components/api/request"; 
import VueCookies from 'vue-cookies'
Vue.use(VueCookies)
Vue.prototype.$EventBus = new Vue()

Vue.prototype.get_GraphQl_Str = ((node, name, properties) => {
  var obj = node[properties]
  var getObj = {}
  var getS = ((obj, aaaa) => {
    obj.map((k, i) => {
      aaaa[k[name]] = ""
      if (k[properties]) {
        aaaa[k[name]] = {}
        getS(k[properties], aaaa[k[name]])
      }
    })
  })
  getS(obj, getObj)
  return getObj
})

var logDebug = (location.origin.indexOf("http://localhost") != -1) ? true : false;
if (!logDebug) {
  console.log = (function (oriLogFunc) {
    return function () {
      if (false) {
        oriLogFunc.apply(this, arguments);
      }
    }
  })(console.log);
}

Vue.use(ElementUI, { locale })
import _ from 'lodash'
Vue.prototype._ = _

import moment from "moment"
Vue.prototype.$moment = moment

Vue.filter("filterTimeShow", function (value, status, type) {
  if (!type) {
    var _type = "YYYY-MM-DD HH:mm:ss"
  } else {
    var _type = type
  }
  if (value && value.toString().indexOf('Z') != -1 && status) {
    return moment(value.split('Z')[0]).format(_type)
  } else {
    return moment(value).format(_type)
  }
})

import G6 from '@antv/g6'
Vue.prototype.G6 = G6
Vue.config.productionTip = false

new Vue({
  el: '#app',
  router,
  store,
  render: h => h(App),
  created() {
    setVueInstance(this);
  },
})

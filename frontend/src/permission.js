import router from './router'
import store from './store'
import NProgress from 'nprogress'
import 'nprogress/nprogress.css'

NProgress.configure({ showSpinner: false })
router.beforeEach(async (to, from, next) => {
  NProgress.start()
  const route = store.state.user.addRoutes
  const username = store.state.user.loginUserName
  const hasRouters = route && route.length > 0
  if (to.path == '/login') {
    if (username) {
      next({ path: "/dashboard", replace: true })
    } else {
      try {
        await store.dispatch('user/generateRoutes')
        // 动态添加格式化过的路由
        router.addRoutes(store.state.user.addRoutes)
        const username1 = store.state.user.loginUserName
        if (username1) {
          next({ path: "/dashboard", replace: true })
        } else {
          next()
        }
      } catch (error) {
        next(`/`)
        NProgress.done()
      }
    }
  } else {
    if (hasRouters) {
      if (to.path == '/login') {
        next({ path: "/dashboard", replace: true })
      } else {
        next()
      }
    } else {
      try {
        await store.dispatch('user/generateRoutes')
        // 动态添加格式化过的路由
        router.addRoutes(store.state.user.addRoutes)
        const username2 = store.state.user.loginUserName
        if (!username2) {
          next({ path: "/login", replace: true })
        } else {
          next({ path: location.pathname + location.search, replace: true })
        }
      } catch (error) {
        next(`/`)
        NProgress.done()
      }
    }
  }
})

router.afterEach(() => {
  NProgress.done()
})

import {get, post} from '../axios/http.js'
import GLOBAL from '../common/global.js'

export const getUserInfo = (data) => post(`${GLOBAL.HOME}/api/getUserInfo`, data)
export const login = (data) => get(`${GLOBAL.HOME}/api/user/login`, data)

import {get, post} from '../common/http.js'
import GLOBAL from '../common/global.js'

export const getUserInfo = (data) => post(`${GLOBAL.HOME}/getUserInfo`, data)
export const login = (data) => get(`${GLOBAL.HOME}/user/login`, data)

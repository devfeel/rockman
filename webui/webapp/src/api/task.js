import {get, post} from '../axios/http.js'
import GLOBAL from '../common/global.js'

export const getTaskList = (data) => post(`${GLOBAL.HOME}/task/list`, data)

export const getTaskOnce = (data) => get(`${GLOBAL.HOME}/task/get`, data)

export const taskSave = (data) => post(`${GLOBAL.HOME}/task/save`, data)

export const taskDelete = (data) => get(`${GLOBAL.HOME}/task/delete`, data)

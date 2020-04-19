import {get, post} from '../axios/http.js'
import GLOBAL from '../common/global.js'

export const getTaskList = (data) => get(`${GLOBAL.HOME}/api/task/list`, data)

export const getExecLogList = (data) => get(`${GLOBAL.HOME}/api/task/execlogs`, data)

export const getStateLogList = (data) => get(`${GLOBAL.HOME}/api/task/statelogs`, data)

export const getTaskOnce = (data) => get(`${GLOBAL.HOME}/api/task/get`, data)

export const taskSave = (data) => post(`${GLOBAL.HOME}/api/task/save`, data)

export const taskUpdate = (data) => post(`${GLOBAL.HOME}/api/task/update`, data)

export const taskDelete = (data) => get(`${GLOBAL.HOME}/api/task/delete`, data)

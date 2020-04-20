import {get, post} from '../axios/http.js'
import GLOBAL from '../common/global.js'

export const getTaskList = (data) => get(`${GLOBAL.HOME}/task/list`, data)

export const getExecLogList = (data) => get(`${GLOBAL.HOME}/task/execlogs`, data)

export const getStateLogList = (data) => get(`${GLOBAL.HOME}/task/statelogs`, data)

export const getTaskOnce = (data) => get(`${GLOBAL.HOME}/task/get`, data)

export const taskSave = (data) => post(`${GLOBAL.HOME}/task/save`, data)

export const taskUpdate = (data) => post(`${GLOBAL.HOME}/task/update`, data)

export const taskDelete = (data) => get(`${GLOBAL.HOME}/task/delete`, data)

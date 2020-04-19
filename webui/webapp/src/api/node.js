import {get, post} from '../axios/http.js'
import GLOBAL from '../common/global.js'

export const getNodeList = (data) => get(`${GLOBAL.HOME}/api/node/list`, data)

export const getNodeOnce = (data) => get(`${GLOBAL.HOME}/api/node/get`, data)

export const nodeSave = (data) => post(`${GLOBAL.HOME}/api/node/save`, data)

export const nodeDelete = (data) => get(`${GLOBAL.HOME}/api/node/delete`, data)

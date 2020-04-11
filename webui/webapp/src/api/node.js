import {get, post} from '../axios/http.js'
import GLOBAL from '../common/global.js'

export const getNodeList = (data) => post(`${GLOBAL.HOME}/node/list`, data)

export const getNodeOnce = (data) => get(`${GLOBAL.HOME}/node/get`, data)

export const nodeSave = (data) => post(`${GLOBAL.HOME}/node/save`, data)

export const nodeDelete = (data) => get(`${GLOBAL.HOME}/node/delete`, data)

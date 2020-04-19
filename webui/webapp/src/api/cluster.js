import {get} from '../axios/http.js'
import GLOBAL from '../common/global.js'

export const getClusterInfo = (data) => get(`${GLOBAL.HOME}/api/cluster/info`, data)

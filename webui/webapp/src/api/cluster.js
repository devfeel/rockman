import {get} from '../common/http.js'
import GLOBAL from '../common/global.js'

export const getClusterInfo = (data) => get(`${GLOBAL.HOME}/cluster/info`, data)

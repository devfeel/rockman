import { get, post } from '../axios/http.js';
import GLOBAL from '../common/global.js';

export const getLogList = data => post(`${GLOBAL.HOME}/logs/query`, data);

export const getLogByTaskId = data =>
  get(`${GLOBAL.HOME}/logs/getbytaskid`, data);

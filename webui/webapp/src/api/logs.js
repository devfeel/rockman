import { get } from '../common/http.js';
import GLOBAL from '../common/global.js';

export const getNodeTraceList = data => get(`${GLOBAL.HOME}/log/trace`, data);
export const getTaskExecList = data => get(`${GLOBAL.HOME}/log/exec`, data);
export const getTaskStateList = data => get(`${GLOBAL.HOME}/log/state`, data);
export const getTaskSubmitList = data => get(`${GLOBAL.HOME}/log/submit`, data);

Leader\Worker协议：
master:
1. ElectionLeader
2. SubmitExecutor
3. QueryWorkers

worker:
1. GetLeaderInfo
2. RegisterWorker
3. RegisterExecutor
4. StartExecutor
5. StopExecutor
6. RemoveExecutor
7. QueryExecutors

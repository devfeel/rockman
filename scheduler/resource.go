package scheduler

type (
	ResourceInfo struct {
		EndPoint   string
		CpuRate    int //cpu rate, refresh per 1 minute
		MemoryRate int //memory rate, refresh per 1 minute
		TaskCount  int //job count
		LoadValue  int //load value = cpu * 30 + memory * 30 + jobs * 40
	}

	LoadResources   []*ResourceInfo
	CpuResources    []*ResourceInfo
	MemoryResources []*ResourceInfo
	JobResources    []*ResourceInfo
)

// refreshLoadValue refresh resource's load value
func (r *ResourceInfo) refreshLoadValue() int {
	r.LoadValue = r.CpuRate*30 + r.MemoryRate*30 + r.TaskCount*40
	return r.LoadValue
}

func (rs LoadResources) Len() int {
	return len(rs)
}
func (rs LoadResources) Less(i, j int) bool {
	return rs[i].LoadValue > rs[j].LoadValue
}
func (rs LoadResources) Swap(i, j int) {
	rs[i], rs[j] = rs[j], rs[i]
}

func (rs CpuResources) Len() int {
	return len(rs)
}
func (rs CpuResources) Less(i, j int) bool {
	return rs[i].CpuRate > rs[j].CpuRate
}
func (rs CpuResources) Swap(i, j int) {
	rs[i], rs[j] = rs[j], rs[i]
}

func (rs MemoryResources) Len() int {
	return len(rs)
}
func (rs MemoryResources) Less(i, j int) bool {
	return rs[i].MemoryRate > rs[j].MemoryRate
}
func (rs MemoryResources) Swap(i, j int) {
	rs[i], rs[j] = rs[j], rs[i]
}

func (rs JobResources) Len() int {
	return len(rs)
}
func (rs JobResources) Less(i, j int) bool {
	return rs[i].TaskCount > rs[j].TaskCount
}
func (rs JobResources) Swap(i, j int) {
	rs[i], rs[j] = rs[j], rs[i]
}

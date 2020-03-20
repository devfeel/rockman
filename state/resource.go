package state

type (
	LoadResources   []*ResourceInfo
	CpuResources    []*ResourceInfo
	MemoryResources []*ResourceInfo
	JobResources    []*ResourceInfo
)

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
	return rs[i].JobCount > rs[j].JobCount
}
func (rs JobResources) Swap(i, j int) {
	rs[i], rs[j] = rs[j], rs[i]
}

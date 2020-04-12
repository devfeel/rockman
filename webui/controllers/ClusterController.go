package controllers

import (
	"bytes"
	"encoding/json"
	"html/template"

	"github.com/devfeel/dotweb"
	"github.com/devfeel/rockman/node"
	_const "github.com/devfeel/rockman/webui/const"
)

type ClusterController struct {
}

func (c *ClusterController) ShowClusterInfo(ctx dotweb.Context) error {
	node := getNode(ctx)
	if node == nil {
		return ctx.WriteJson(NewResponse(-1001, "not exists node in app items", nil))
	}
	return ctx.WriteJson(SuccessResponse(node.Cluster.ClusterInfo()))
}

func (c *ClusterController) ShowExecutors(ctx dotweb.Context) error {
	item, isExists := ctx.AppItems().Get(_const.ItemKeyNode)
	if !isExists {
		return ctx.WriteJson(NewResponse(-1001, "not exists node in app items", nil))
	}
	node, isOk := item.(*node.Node)
	if !isOk {
		return ctx.WriteJson(NewResponse(-1002, "not exists correct node in app items", nil))
	}
	return ctx.WriteHtml(FormatJson(NewResponse(0, "", node.Cluster.Executors)))
}

func (c *ClusterController) ShowResources(ctx dotweb.Context) error {
	item, isExists := ctx.AppItems().Get(_const.ItemKeyNode)
	if !isExists {
		return ctx.WriteJson(NewResponse(-1001, "not exists node in app items", nil))
	}
	node, isOk := item.(*node.Node)
	if !isOk {
		return ctx.WriteJson(NewResponse(-1002, "not exists correct node in app items", nil))
	}
	return ctx.WriteHtml(FormatJson(NewResponse(0, "", node.Cluster.Scheduler.Resources())))
}

func FormatJson(data interface{}) string {

	// 格式化Json，添加\t符
	by, _ := json.MarshalIndent(data, "", "\t")
	task := string(by)

	content := struct {
		Task string
	}{
		Task: task,
	}

	// 定义html格式
	const html = `<html>
		<head>
			<meta charset="utf-8" />
		</head>
		<body>
			<div>
				<pre>{{.Task}}</pre>
			</div>
		</body>
	</html>`

	// 使用html渲染json
	var doc bytes.Buffer
	temp, _ := template.New("").Parse(html)
	temp.Execute(&doc, content)

	return doc.String()
}

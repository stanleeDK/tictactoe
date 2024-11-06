package treantchart



// ------------------------------------------------------------------
/*    var simple_chart_config = {
        chart: {
            container: "#OrganiseChart-simple"
        },

        nodeStructure: {
            innerHTML: "</br><table><tr><td>666</td></tr><tr><td>123</td></tr></table>",
            children: [
                {
                    innerHTML: "First Child </br><table><tr><td>333</td></tr><tr><td>123</td></tr></table>",
                    children: [
                        {
                            innerHTML: "</br><table><tr><td>123</td></tr><tr><td>123</td></tr></table>"
                        },
                        {
                            innerHTML: "</br><table><tr><td>123</td></tr><tr><td>123</td></tr></table>"
                        },
                        {
                            innerHTML: "</br><table><tr><td>123</td></tr><tr><td>123</td></tr></table>"
                        },
                        {
                            innerHTML: "</br><table><tr><td>123</td></tr><tr><td>123</td></tr></table>"
                        }
                    ]
                },
                {
                    innerHTML: "level 1 </br><table><tr><td>888</td></tr><tr><td>123</td></tr></table>" 
                }
            ]
        }*/

//     chart: {
//         container: "#OrganiseChart-simple"
//     },
type TreantChart struct {
	Chart TreantChartContainer `json:"chart"`
	NodeStructure TreantNodeStructure `json:"nodeStructure,omitempty"`
}
type TreantChartContainer struct {
	Container string `json:"container"`
    Scrollbar string `json:"scrollbar"`
}

//     nodeStructure: {
//         innerHTML: "</br><table><tr><td>666</td></tr><tr><td>123</td></tr></table>",
//         children: [
//             {    
type TreantNodeStructure struct {
	InnerHTML string 			   `json:"innerHTML,omitempty"`
	Children []TreantNodeStructure `json:"children,omitempty"`
}
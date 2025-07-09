package web;

import(
	Strings "strings"
	JSON    "encoding/json"
);

import _ "embed";



//go:embed base.tpl
var TPL_Base []byte;

//go:embed includes.json
var JSON_Includes []byte;



func (build *Builder) WithIncludes() *Builder {
	// load versions
	versions := make(map[string]string);
	if err := JSON.Unmarshal(JSON_Includes, &versions); err != nil { panic(err); }
	// bootstrap
	build.AddTagURL("URL_BootstrapCSS",           URL_BootstrapCSS,           versions["twbs/bootstrap"]);
	build.AddTagURL("URL_BootstrapJS",            URL_BootstrapJS,            versions["twbs/bootstrap"]);
	build.AddTagURL("URL_BootsIconsCSS",          URL_BootsIconsCSS,          versions["twbs/icons"    ]);
	build.AddTagURL("URL_PopperJS",               URL_PopperJS,      versions["floating-ui/floating-ui"]);
	// jquery
	build.AddTagURL("URL_JQueryJS",               URL_JQueryJS,               versions["jquery/jquery" ]);
	// datatables
	version_datatables := versions["DataTables/Dist-DataTables-Bootstrap5"];
	build.AddTagURL("URL_DataTablesJS",           URL_DataTablesJS,           version_datatables        );
	build.AddTagURL("URL_DataTablesBootstrapJS",  URL_DataTablesBootstrapJS,  version_datatables        );
	build.AddTagURL("URL_DataTablesBootstrapCSS", URL_DataTablesBootstrapCSS, version_datatables        );
	build.AddTagURL("URL_DataTablesScrollerJS",   URL_DataTablesScrollerJS,   version_datatables        );
	build.AddTagURL("URL_DataTablesScrollerCSS",  URL_DataTablesScrollerCSS,  version_datatables        );
	build.AddTagURL("URL_DataTablesPageResizeJS", URL_DataTablesPageResizeJS, version_datatables        );
	// echarts
	build.AddTagURL("URL_EChartsJS",              URL_EChartsJS,              versions["apache/echarts"]);
	return build;
}

func (build *Builder) AddTagURL(key string, url string, version string) {
	if build.IsDev { url = Strings.Replace(url, ".min.", ".", 1); }
	build.Tags[key] = Strings.Replace(url, "{{VERSION}}", version, 1);
}



func (build *Builder) WithBootstrap()  *Builder { build.Tags[Tag_WithBootstrap ] = true; return build;                  }
func (build *Builder) WithBootsIcons() *Builder { build.Tags[Tag_WithBootsIcons] = true; return build.WithBootstrap();  }
func (build *Builder) WithPopper()     *Builder { build.Tags[Tag_WithPopper    ] = true; return build.WithBootstrap();  }
func (build *Builder) WithTooltips()   *Builder { build.Tags[Tag_WithTooltips  ] = true; return build.WithPopper();     }
func (build *Builder) WithJQuery()     *Builder { build.Tags[Tag_WithJQuery    ] = true; return build;                  }
func (build *Builder) WithDataTables() *Builder { build.Tags[Tag_WithDataTables] = true; return build.WithBootstrap();  }
func (build *Builder) WithECharts()    *Builder { build.Tags[Tag_WithECharts   ] = true; return build;                  }

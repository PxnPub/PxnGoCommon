package web;



const LogPrefixWeb   = "[Web] ";
const DefaultBindWeb = "tcp://127.0.0.1:8000";



// mimes
const(
	Mime_Text = "text/plain"
	Mime_HTML = "text/html"
	Mime_JSON = "application/json"
	Mime_SVG  = "image/svg+xml"
);



// template tags
const(
	Tag_Title          = "Title"
	Tag_Page           = "Page"
	Tag_FavIcon        = "FavIcon"
	Tag_FilesCSS       = "FilesCSS"
	Tag_FilesJS        = "FilesJS"
	Tag_RawCSS         = "RawCSS"
	Tag_RawJS          = "RawJS"
	// bootstrap
	Tag_WithBootstrap  = "WithBootstrap"
	Tag_WithBootsIcons = "WithBootsIcons"
	Tag_WithTooltips   = "WithTooltips"
	Tag_WithPopper     = "WithPopper"
	// jquery
	Tag_WithJQuery     = "WithJQuery"
	// datatables
	Tag_WithDataTables = "WithDataTables"
	// echarts
	Tag_WithECharts    = "WithECharts"
	Tag_AppendHead     = "AppendHead"
	Tag_AppendHeader   = "AppendHeader"
	Tag_AppendFooter   = "AppendFooter"
);

// include urls
const(
	// bootstrap
	URL_BootstrapCSS           = "https://cdn.jsdelivr.net/npm/bootstrap@{{VERSION}}/dist/css/bootstrap.min.css";
	URL_BootstrapJS            = "https://cdn.jsdelivr.net/npm/bootstrap@{{VERSION}}/dist/js/bootstrap.bundle.min.js";
	URL_BootsIconsCSS          = "https://cdn.jsdelivr.net/npm/bootstrap-icons@{{VERSION}}/font/bootstrap-icons.min.css";
	URL_PopperJS               = "https://cdn.jsdelivr.net/npm/@popperjs/core@{{VERSION}}/dist/umd/popper.min.js";
	// jquery
	URL_JQueryJS               = "https://cdnjs.cloudflare.com/ajax/libs/jquery/{{VERSION}}/jquery.min.js";
	// datatables
	URL_DataTablesJS           = "https://cdn.datatables.net/{{VERSION}}/js/dataTables.min.js";
	URL_DataTablesBootstrapJS  = "https://cdn.datatables.net/{{VERSION}}/js/dataTables.bootstrap5.min.js";
	URL_DataTablesBootstrapCSS = "https://cdn.datatables.net/{{VERSION}}/css/dataTables.bootstrap5.min.css";
	URL_DataTablesScrollerJS   = "https://cdn.datatables.net/scroller/{{VERSION}}/js/dataTables.scroller.min.js";
	URL_DataTablesScrollerCSS  = "https://cdn.datatables.net/scroller/{{VERSION}}/css/scroller.bootstrap5.min.css";
	URL_DataTablesPageResizeJS = "https://cdn.datatables.net/plug-ins/{{VERSION}}/features/pageResize/dataTables.pageResize.min.js";
	// echarts
	URL_EChartsJS              = "https://cdnjs.cloudflare.com/ajax/libs/echarts/{{VERSION}}/echarts.min.js"
);

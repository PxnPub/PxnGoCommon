package html;

import(
	Fmt  "fmt"
	HTTP "net/http"
);



const(
	MIME_TEXT = "text/plain"
	MIME_HTML = "text/html"
	MIME_JSON = "application/json"
);



type Builder struct {
	IsDev             bool
	Title             string
	FavIcon           string
	AppendHead        string
	AppendHeader      string
	AppendFooter      string
	IsBootstrap       bool
	IsBootstrapIcons  bool
	IsBootstrapPopper bool
	IsEnableTooltips  bool
	IsJQuery          bool
	IsDataTables      bool
	IsECharts         bool
}



func NewBuilder() *Builder {
	return &Builder{};
}



func (build *Builder) Render(contents string) string {
	return build.RenderTop() + "\n" +
		contents +
		"\n" + build.RenderBottom();
}

func (build *Builder) RenderTop() string {
	out := `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8" />
<meta name="viewport" content="width=device-width,initial-scale=1.0" />
<meta http-equiv="Cache-Control" content="max-age=3600, must-revalidate" />
`;
	if build.Title != ""      { out += Fmt.Sprintf("<title>%s</title>\n", build.Title); }
	if build.FavIcon != ""    { out += Fmt.Sprintf(`<link rel="icon" type="image/x-icon" href="%s" />`, build.FavIcon) + "\n"; }
	if build.IsBootstrap      { out += Fmt.Sprintf(`<link rel="stylesheet" href="%s" />`, URL_BootstrapCSS           ) + "\n"; }
	if build.IsBootstrapIcons { out += Fmt.Sprintf(`<link rel="stylesheet" href="%s" />`, URL_BootstrapIconsCSS      ) + "\n"; }
	if build.IsDataTables     { out += Fmt.Sprintf(`<link rel="stylesheet" href="%s" />`, URL_DataTablesBootstrapCSS ) + "\n"; }
	if build.IsDataTables     { out += Fmt.Sprintf(`<link rel="stylesheet" href="%s" />`, URL_DataTablesScrollerCSS  ) + "\n"; }
	if build.AppendHead != "" { out += build.AppendHead; }
	out += "</head>\n<body>\n\n\n";
	if build.AppendHeader != "" { out += build.AppendHeader + "\n\n\n"; }
	return out;
}

func (build *Builder) RenderBottom() string {
	out := "\n\n\n";
	if build.IsJQuery          { out += Fmt.Sprintf(`<script src="%s"></script>`, URL_JQueryJS              ) + "\n"; }
	if build.IsBootstrapPopper { out += Fmt.Sprintf(`<script src="%s"></script>`, URL_BootstrapPopperJS     ) + "\n"; }
	if build.IsBootstrap       { out += Fmt.Sprintf(`<script src="%s"></script>`, URL_BootstrapJS           ) + "\n"; }
	if build.IsDataTables      { out += Fmt.Sprintf(`<script src="%s"></script>`, URL_DataTablesJS          ) + "\n"; }
	if build.IsDataTables      { out += Fmt.Sprintf(`<script src="%s"></script>`, URL_DataTablesBootstrapJS ) + "\n"; }
	if build.IsDataTables      { out += Fmt.Sprintf(`<script src="%s"></script>`, URL_DataTablesScrollerJS  ) + "\n"; }
	if build.IsDataTables      { out += Fmt.Sprintf(`<script src="%s"></script>`, URL_DataTablesPageResizeJS) + "\n"; }
	if build.IsECharts         { out += Fmt.Sprintf(`<script src="%s"></script>`, URL_EChartsJS             ) + "\n"; }
	if build.IsEnableTooltips  { out += "<script>\nvar tooltipTriggerList = "             +
		"[].slice.call(document.querySelectorAll('[data-bs-toggle=\"tooltip\"]'));\n" +
		"var tooltipList = tooltipTriggerList.map(function (tooltipTriggerEl) {\n"    +
		"\treturn new bootstrap.Tooltip(tooltipTriggerEl);\n});\n</script>\n"; }
	if build.AppendFooter != "" { out += build.AppendFooter; }
	if out != "\n\n\n" { out += "\n\n\n"; }
	out += "</body>\n</html>\n";
	return out;
}



func SetContentType(out HTTP.ResponseWriter, mime string) {
	switch mime {
	case "text": mime = MIME_TEXT; break;
	case "html": mime = MIME_HTML; break;
	case "json": mime = MIME_JSON; break;
	default:                       break;
	}
	out.Header().Set("Content-Type", mime);
}



// title
func (build *Builder) SetTitle(title string) *Builder {
	build.Title = title;
	return build;
}

// fav icon
func (build *Builder) SetFavIcon(icon string) *Builder {
	build.FavIcon = icon;
	return build;
}



// css
func (build *Builder) AddCSS(path string) *Builder {
	build.AppendHead += Fmt.Sprintf(`<link rel="stylesheet" href="%s" />`, PubDevURL(build.IsDev, path)) + "\n";
	return build;
}

func (build *Builder) AddRawCSS(css string) *Builder {
	build.AppendHead += `<style type="text/css">` + "\n" + css + "\n</style>\n";
	return build;
}



// js
func (build *Builder) AddTopJS(path string) *Builder {
	build.AppendHead += Fmt.Sprintf(`<script src="%s"></script>`, PubDevURL(build.IsDev, path)) + "\n";
	return build;
}

func (build *Builder) AddBotJS(path string) *Builder {
	build.AppendFooter += Fmt.Sprintf(`<script src="%s"></script>`, PubDevURL(build.IsDev, path)) + "\n";
	return build;
}

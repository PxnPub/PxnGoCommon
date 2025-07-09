<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8" />
<meta name="viewport" content="width=device-width,initial-scale=1.0" />
<meta http-equiv="Cache-Control" content="no-cache, no-store, must-revalidate" />
<meta http-equiv="Pragma" content="no-cache" />
<meta http-equiv="Expires" content="0" />
{{if        .Title         }}<title>{{.Title}}</title>
{{end}}{{if .FavIcon       }}<link rel="icon" type="image/x-icon" href="{{.FavIcon}}" />
{{end}}{{if .WithBootstrap }}<link rel="stylesheet" href="{{.URL_BootstrapCSS     }}" />
{{end}}{{if .WithBootsIcons}}<link rel="stylesheet" href="{{.URL_BootsIconsCSS    }}" />
{{end}}{{if .WithDataTables}}<link rel="stylesheet" href="{{.URL_DataTablesBSCSS  }}" />
	<link rel="stylesheet" href="{{.URL_DataTablesScrollCSS}}" />
{{end}}{{if .FilesCSS      }}{{range .FilesCSS}}<link rel="stylesheet" href="{{.}}" />
{{end}}{{end}}{{if .RawCSS }}<style>
{{.RawCSS}}
</style>
{{end}}{{.AppendHead}}</head>
<body>
{{.AppendHeader}}
<!-- Start Page Content -->
{{template "PageBody" .}}
<!-- End Page Content -->
{{if        .WithJQuery    }}<script src="{{.URL_JQueryJS    }}"></script>
{{end}}{{if .WithPopper    }}<script src="{{.URL_PopperJS    }}"></script>
{{end}}{{if .WithBootstrap }}<script src="{{.URL_BootstrapJS }}"></script>
{{end}}{{if .WithDataTables}}<script src="{{.URL_DataTablesJS}}"></script>
	<script src="{{.URL_DataTablesBootstrapJS }}"></script>
	<script src="{{.URL_DataTablesScrollerJS  }}"></script>
	<script src="{{.URL_DataTablesPageResizeJS}}"></script>
{{end}}{{if .WithECharts        }}<script src="{{.URL_EChartsJS}}"></script>
{{end}}{{if .WithTooltips}}<script>
const tooltipTriggerList = document.querySelectorAll('[data-bs-toggle="tooltip"]');
const tooltipList = [...tooltipTriggerList].map(
	tooltipTriggerEl => new bootstrap.Tooltip(tooltipTriggerEl));
</script>
{{end}}{{if .FilesJS}}{{range .FilesJS}}<script src="{{.}}"></script>
{{end}}{{end}}{{if .RawJS }}
<script>
{{.RawJS}}
</script>
{{end}}{{.AppendFooter}}
</body>
</html>

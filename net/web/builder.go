package web;

import(
	Log   "log"
	Templ "html/template"
);



type Builder struct {
	IsDev bool
	TPL   *Templ.Template
	Tags  map[string]interface{}
}



func NewBuilder() *Builder {
	tpl, err := Templ.New("website").Parse(string(TPL_Base));
	if err != nil { Log.Panicf("%s, in NewFactory()", err); }
	tags := make(map[string]interface{});

//TODO: InitIncludes() ?

	return &Builder{
//TODO
		IsDev: true,
		TPL:  tpl,
		Tags: tags,
	};
}



func (build *Builder) Clone() *Builder {
	tpl, err := build.TPL.Clone();
	if err != nil { Log.Panicf("%s, in Factory->Clone()", err); }
	return &Builder{
		IsDev: build.IsDev,
		TPL:   tpl,
		Tags:  build.CloneTags(),
	};
}

func (build *Builder) CloneTags() map[string]interface{} {
	tags := make(map[string]interface{});
	for k, v := range build.Tags { tags[k] = v; }
	return tags;
}



// tpl
func (build *Builder) AddFilesTPL(files...string) *Builder {
	tpl, err := build.TPL.ParseFiles(files...);
	if err != nil { Log.Panicf("%s, in Factory->AddFiles()", err); }
	build.TPL = tpl;
	return build;
}

func (build *Builder) AddRawTPL(data []byte) *Builder {
	tpl, err := build.TPL.Parse(string(data));
	if err != nil { Log.Panicf("%s, in Factory->AddRawTPL()", err); }
	build.TPL = tpl;
	return build;
}



// css
func (build *Builder) AddFilesCSS(files...string) *Builder {
	array := build.GetTagStringArray(Tag_FilesCSS);
	return build.SetTag(Tag_FilesCSS, append(array, files[:]...));
}

func (build *Builder) AddRawCSS(data []byte) *Builder {
	css := build.GetTagString(Tag_RawCSS);
	return build.SetTag(Tag_RawCSS, css);
}

func (build *Builder) AddStringCSS(lines string) *Builder {
	return build.AddRawCSS([]byte(lines));
}



// js
func (build *Builder) AddFilesJS(files...string) *Builder {
	array := build.GetTagStringArray(Tag_FilesJS);
	return build.SetTag(Tag_FilesJS, append(array, files[:]...));
}

func (build *Builder) AddRawJS(data []byte) *Builder {
	js := build.GetTagString(Tag_RawJS);
	return build.SetTag(Tag_RawJS, js);
}

func (build *Builder) AddStringJS(lines string) *Builder {
	return build.AddRawJS([]byte(lines))
}



func (build *Builder) SetTitle(title string) *Builder {
	return build.SetTag("Title", title)
}

func (build *Builder) SetFavIcon(file string) *Builder {
//TODO: routes?
//router.HandleFunc("/favicon.ico", PxnWeb.NewRedirect("/static/line-chart.ico"));
	return build.SetTag("FavIcon", file);
}



// tags
func (build *Builder) GetTagStringArray(key string) []string {
	val, ok := build.Tags[Tag_RawCSS];
	if ok { if strs, ok := val.([]string); ok { return strs; }}
	return []string{};
}

func (build *Builder) GetTagString(key string) string {
	val, ok := build.Tags[key];
	if ok { if str, ok := val.(string); ok { return str; }}
	return "";
}

func (build *Builder) GetTagBool(key string, def bool) bool {
	val, ok := build.Tags[key];
	if ok { if b, ok := val.(bool); ok { return b; }}
	return def;
}

func (build *Builder) SetTag(key string, value interface{}) *Builder {
	build.Tags[key] = value;
	return build;
}

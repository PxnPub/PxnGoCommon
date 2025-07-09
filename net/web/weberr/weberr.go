package weberr;

import(
	HTTP "net/http"
	JSON "encoding/json"
);



type WebError struct {
	Status int,
	Err    error,
	Trace  string
	IsJSON bool
}



func IntServerErr(err error) *WebError {
	trace := make([]byte, 1024);
	n := Runtime.Stack(trace, true);
	return &WebError{
		Status: HTTP.StatusInternalServerError,
		Trace:  string(trace[:n]),
		Err:    err,
	};
}



func (weberr *WebError) Write(out HTTP.ResponseWriter) {
	if weberr.IsJSON {
		json, err := JSON.Marshal(
			struct{
				Error string
			}{
				Error: weberr.Err.Error(),
			}
		);
		if err != nil { panic(err); }
		HTTP.Error(out, json, weberr.Status);
	} else {
		HTTP.Error(out,
			Fmt.Sprintf("%s\n%s", weberr.Err.Error(), weberr.Trace),
			weberr.Status,
		);
	}
}



func (weberr *WebError) IsJSON() *WebError {
	weberr.IsJSON = true;
	return weberr;
}
